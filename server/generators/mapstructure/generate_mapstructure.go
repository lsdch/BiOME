//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run gen_mapstructure_tags.go <file_or_dir>")
		os.Exit(1)
	}
	root := os.Args[1]

	info, err := os.Stat(root)
	if err != nil {
		panic(err)
	}

	var files []string
	if info.IsDir() {
		// Walk directory recursively
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, ".go") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	} else {
		files = append(files, root)
	}

	for _, filename := range files {
		err := processFile(filename)
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", filename, err)
		}
	}
}

func processFile(filename string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	modified := false

	jsonTagRegexp := regexp.MustCompile(`json:"([^",]+)`)

	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}

			// check for preceding comment with @mapstructure
			hasTagComment := false
			if gd.Doc != nil {
				for _, c := range gd.Doc.List {
					if strings.Contains(c.Text, "@mapstructure") {
						hasTagComment = true
						break
					}
				}
			}
			if !hasTagComment && st.Fields != nil {
				// also check field-level comment
				for _, f := range st.Fields.List {
					if f.Doc != nil {
						for _, c := range f.Doc.List {
							if strings.Contains(c.Text, "@mapstructure") {
								hasTagComment = true
								break
							}
						}
					}
				}
			}

			if !hasTagComment {
				continue // skip this struct
			}

			// modify struct fields
			for _, field := range st.Fields.List {
				if field.Tag == nil {
					continue
				}

				tagValue := strings.Trim(field.Tag.Value, "`") // remove backticks

				// Parse JSON tag name
				jsonMatches := jsonTagRegexp.FindStringSubmatch(tagValue)
				var jsonName string
				if len(jsonMatches) > 1 {
					jsonName = jsonMatches[1]
				}

				// Handle inline embedded struct
				if strings.Contains(tagValue, "inline") && !strings.Contains(tagValue, "mapstructure") {
					tagValue += ` mapstructure:",squash"`
				} else if jsonName != "" && !strings.Contains(tagValue, "mapstructure") {
					tagValue += fmt.Sprintf(` mapstructure:"%s"`, jsonName)
				}

				field.Tag.Value = "`" + tagValue + "`"
				modified = true
			}
		}
	}

	if !modified {
		// nothing changed
		return nil
	}

	var buf bytes.Buffer
	printer.Fprint(&buf, fset, f)

	err = os.WriteFile(filename, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Updated struct(s) in %s\n", filename)
	return nil
}
