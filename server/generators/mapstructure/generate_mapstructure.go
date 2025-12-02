//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var jsonTagRegexp = regexp.MustCompile(`json:"([^",]+)`)

// an edit to apply to the original source: replace [start:end) with newText
type edit struct {
	start   int // byte offset
	end     int // byte offset
	newText string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run gen_mapstructure_preserve.go <file-or-dir>")
		os.Exit(1)
	}

	root := os.Args[1]
	info, err := os.Stat(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "stat error: %v\n", err)
		os.Exit(1)
	}

	var files []string
	if info.IsDir() {
		err := filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if fi.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "walk error: %v\n", err)
			os.Exit(1)
		}
	} else {
		files = append(files, root)
	}

	if len(files) == 0 {
		fmt.Println("No .go files found")
		return
	}

	for _, f := range files {
		if err := processFile(f); err != nil {
			fmt.Fprintf(os.Stderr, "error processing %s: %v\n", f, err)
		}
	}
}

func processFile(filename string) error {
	srcBytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	src := string(srcBytes)

	fset := token.NewFileSet()
	// parse with comments
	fileAST, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return err
	}

	var edits []edit
	modified := false

	// iterate declarations so we can read the GenDecl.Doc (preceding comments)
	for _, decl := range fileAST.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}

		// check for the @mapstructure annotation in the GenDecl (the comment block above type)
		hasAnnotation := false
		if gd.Doc != nil {
			for _, c := range gd.Doc.List {
				if strings.Contains(c.Text, "@mapstructure") {
					hasAnnotation = true
					break
				}
			}
		}
		if !hasAnnotation {
			// also check line comments attached to the decl (rare), and specs' docs below
			// we'll check per TypeSpec below as well, so continue.
		}

		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok || st.Fields == nil {
				continue
			}

			// allow annotation either on the GenDecl or the TypeSpec's doc/comment
			typeHasAnnotation := hasAnnotation
			if !typeHasAnnotation {
				if ts.Doc != nil {
					for _, c := range ts.Doc.List {
						if strings.Contains(c.Text, "@mapstructure") {
							typeHasAnnotation = true
							break
						}
					}
				}
				if !typeHasAnnotation && ts.Comment != nil {
					for _, c := range ts.Comment.List {
						if strings.Contains(c.Text, "@mapstructure") {
							typeHasAnnotation = true
							break
						}
					}
				}
			}
			if !typeHasAnnotation {
				// skip this type
				continue
			}

			// process each field inside the struct
			for _, field := range st.Fields.List {
				if field.Tag == nil {
					// no tag to update
					continue
				}

				// compute byte offsets of the tag in the original file using fset.Position
				start := fset.Position(field.Tag.Pos()).Offset
				end := fset.Position(field.Tag.End()).Offset
				if start < 0 || end > len(srcBytes) || start >= end {
					// fallback - skip if positions aren't sane
					continue
				}

				rawTag := src[start:end] // includes backticks
				tagInner := strings.Trim(rawTag, "`")

				origTagInner := tagInner // remember for comparison

				// If json inline is present and mapstructure doesn't have squash, add it
				if strings.Contains(tagInner, "inline") && !strings.Contains(tagInner, "mapstructure") {
					// only add squash if not already present
					tagInner += ` mapstructure:",squash"`
				}

				// Parse json tag name and add mapstructure:"name" if missing
				jsonMatches := jsonTagRegexp.FindStringSubmatch(tagInner)
				if len(jsonMatches) > 1 {
					jsonName := jsonMatches[1]
					if jsonName != "-" && !strings.Contains(tagInner, "mapstructure:") {
						// add corresponding mapstructure tag
						tagInner += fmt.Sprintf(` mapstructure:"%s"`, jsonName)
					}
				}

				// if tagInner changed, schedule an edit
				if tagInner != origTagInner {
					newTag := "`" + tagInner + "`"
					edits = append(edits, edit{
						start:   start,
						end:     end,
						newText: newTag,
					})
					modified = true
				}
			}
		}
	}

	if !modified {
		// nothing to do
		return nil
	}

	// apply edits in reverse order of start offset so indices remain valid
	sort.Slice(edits, func(i, j int) bool { return edits[i].start > edits[j].start })

	out := []byte(src)
	for _, e := range edits {
		if e.start < 0 || e.end > len(out) || e.start > e.end {
			return fmt.Errorf("invalid edit range %d-%d (file %s)", e.start, e.end, filename)
		}
		newOut := make([]byte, 0, len(out)+len(e.newText)-(e.end-e.start))
		newOut = append(newOut, out[:e.start]...)
		newOut = append(newOut, []byte(e.newText)...)
		newOut = append(newOut, out[e.end:]...)
		out = newOut
	}

	// overwrite file
	if err := os.WriteFile(filename, out, 0644); err != nil {
		return err
	}
	fmt.Printf("Updated %s (%d edits)\n", filename, len(edits))
	return nil
}
