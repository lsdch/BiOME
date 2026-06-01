package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"sort"
	"text/template"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/packages"
)

type Enum struct {
	EnumType   string
	EnumValues []string
}

type FileData struct {
	Package string
	Enums   []Enum
}

//go:embed enum.go.tmpl
var enumTemplate string

func main() {
	if len(os.Args) != 2 {
		logrus.Fatalf("usage: %s <package>", os.Args[0])
	}

	if err := generate(os.Args[1]); err != nil {
		logrus.Fatal(err)
	}
}

func generate(pkgPattern string) error {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedFiles,
	}

	pkgs, err := packages.Load(cfg, pkgPattern)
	if err != nil {
		return err
	}

	if packages.PrintErrors(pkgs) > 0 {
		return errors.New("package loading failed")
	}

	if len(pkgs) != 1 {
		return fmt.Errorf(
			"expected 1 package, got %d",
			len(pkgs),
		)
	}

	pkg := pkgs[0]

	enums, err := findEnums(pkg)
	if err != nil {
		return err
	}

	sort.Slice(enums, func(i, j int) bool {
		return enums[i].EnumType < enums[j].EnumType
	})

	data := FileData{
		Package: pkg.Name,
		Enums:   enums,
	}

	var buf bytes.Buffer

	tmpl := template.Must(
		template.New("file").Parse(enumTemplate),
	)

	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		src = buf.Bytes()
	}

	pkgDir := filepath.Dir(pkg.GoFiles[0])

	output := filepath.Join(
		pkgDir,
		"enums_gen.go",
	)

	return os.WriteFile(output, src, 0o644)
}

func findEnums(pkg *packages.Package) ([]Enum, error) {
	enumMap := map[string]*Enum{}

	for _, file := range pkg.Syntax {
		filename := pkg.Fset.Position(file.Pos()).Filename

		if filepath.Base(filename) == "enums_gen.go" {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			decl, ok := n.(*ast.GenDecl)
			if !ok || decl.Tok != token.CONST {
				return true
			}

			for _, spec := range decl.Specs {
				vs := spec.(*ast.ValueSpec)

				for _, name := range vs.Names {
					obj := pkg.TypesInfo.Defs[name]

					c, ok := obj.(*types.Const)
					if !ok {
						continue
					}

					named, ok := c.Type().(*types.Named)
					if !ok {
						continue
					}

					enumName := named.Obj().Name()

					enum, exists := enumMap[enumName]
					if !exists {
						enum = &Enum{
							EnumType: enumName,
						}
						enumMap[enumName] = enum
					}

					enum.EnumValues = append(enum.EnumValues, c.Name())
				}
			}

			return true
		})
	}

	result := make([]Enum, 0, len(enumMap))
	for _, enum := range enumMap {
		if len(enum.EnumValues) == 0 {
			continue
		}
		result = append(result, *enum)
	}

	return result, nil
}
