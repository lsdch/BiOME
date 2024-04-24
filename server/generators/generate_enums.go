//go:build ignore
// +build ignore

// Code generator to wrap gin handlers and provide additional arguments

package main

import (
	_ "embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	// Walk through all directories and subdirectories
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			parseDir(path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error walking through directories:", err)
	}
}

type EnumDecl struct {
	Decl    *ast.GenDecl
	Ordered bool
}

func parseDir(path string) {
	// Parse current directory and its subdirectories for Go files
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(
		fset,
		path,
		func(info os.FileInfo) bool {
			return !info.IsDir() && strings.HasSuffix(info.Name(), ".go")
		},
		parser.ParseComments,
	)

	if err != nil {
		fmt.Println("Error parsing directory:", err)
		return
	}

	for _, pkg := range pkgs {
		for fileName, file := range pkg.Files {
			var foundEnums []EnumDecl // Stores enum value declarations
			for _, decl := range file.Decls {
				if genEnum, ok := decl.(*ast.GenDecl); ok && genEnum.Doc != nil {
					for _, comment := range genEnum.Doc.List {
						if strings.Contains(comment.Text, "generate:enum") {
							foundEnums = append(foundEnums, EnumDecl{Decl: genEnum})
						} else if strings.Contains(comment.Text, "generate:order-enum") {
							foundEnums = append(foundEnums, EnumDecl{Decl: genEnum, Ordered: true})
						}
					}
				}
			}
			if len(foundEnums) > 0 {
				generateEnumCode(pkg, fileName, foundEnums)
			}
		}
	}
}

type EnumData struct {
	EnumType   string
	EnumValues []string
	Ordered    bool
}

type EnumTemplateData struct {
	Pkg   string
	Enums []EnumData
}

//go:embed enum.go.tmpl
var enumTemplate string

func generateFileName(fileName string) string {
	fileExtension := filepath.Ext(fileName)
	strippedPath := strings.TrimSuffix(fileName, fileExtension)
	return fmt.Sprintf("%s_gen%s", strippedPath, fileExtension)
}

func generateEnumTemplateData(decls []EnumDecl) []EnumData {
	var enums []EnumData

	for _, decl := range decls {
		var enumType string
		var values []string

		for _, spec := range decl.Decl.Specs {
			if valueSpec, ok := spec.(*ast.ValueSpec); ok {
				if t, ok := valueSpec.Type.(*ast.Ident); ok {
					enumType = t.Name
				}
				for _, enumValue := range valueSpec.Names {
					values = append(values, enumValue.String())
				}
			}
		}
		if enumType != "" && len(values) > 0 {
			enums = append(enums, EnumData{
				EnumType:   enumType,
				EnumValues: values,
				Ordered:    decl.Ordered,
			})
		}
	}
	return enums
}

func generateEnumCode(pkg *ast.Package, fileName string, decls []EnumDecl) {

	generatedFileName := generateFileName(fileName)

	data := EnumTemplateData{
		Pkg:   pkg.Name,
		Enums: generateEnumTemplateData(decls),
	}

	tmpl := template.Must(template.New("wrapper").Parse(enumTemplate))

	file, err := os.Create(generatedFileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Generated enum boilerplate: %s\n", filepath.Base(generatedFileName))
}
