//go:build ignore
// +build ignore

// Code generator to wrap gin handlers and provide additional arguments

package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"os"
	"text/template"

	"github.com/sirupsen/logrus"
)

//go:embed route_methods.go.tmpl
var wrapperTemplate string

const OUTPUT_FILE = "wrapped_methods.go"

type methodInfo struct {
	MethodName        string
	WrapperMethodName string
}

func main() {
	methods := []string{"GET", "POST", "PUT", "PATCH", "OPTIONS", "DELETE"}

	file, err := os.Create(OUTPUT_FILE)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	generateWrapperCode(file, methods)
}

func generateWrapperCode(file *os.File, methods []string) {

	tmpl, err := template.New("wrapper").Parse(wrapperTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, methods)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	// Format the generated code
	formattedCode, err := format.Source(buffer.Bytes())
	if err != nil {
		fmt.Println("Error formatting code:", err)
		return
	}
	logrus.Debug(formattedCode)

	// Write the formatted code back to the file
	file.Truncate(0)
	file.Seek(0, 0)
	file.Write(formattedCode)
}
