// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"text/template"
)

var temporaryShellScriptName = "rangorun"

// Generate_compile_run takes a list of Go source lines, create a Go program from it, compiles that source and runs that program.
// Return the captured output from the compilation or the execution of the Go program.
func Generate_compile_run(goSourceFile string, sourceLines []SourceHolder) string {
	err := generate(goSourceFile, sourceLines)
	if err != nil {
		return fmt.Sprintf("[rango] generate Go source failed:%v", err)
	}
	output, _ := gorun(goSourceFile)
	return output
}

// generate produces a Go source file from a list of Go code sourceLines
func generate(goSourceFile string, sourceLines []SourceHolder) error {
	t := template.Must(template.New("image").Parse(imageSourceTemplate()))
	var sourceBuffer bytes.Buffer
	t.Execute(&sourceBuffer, buildTemplateVars(sourceLines))
	return ioutil.WriteFile(goSourceFile, sourceBuffer.Bytes(), 0644)
}

// gorun compiles and executes the source found in goSourceFile and reports its output
// The output can be compiler or runtime error info, otherwise the output is what the program spits out.
func gorun(goSourceFile string) (string, error) {
	logName := fmt.Sprintf("%s.build.log", goSourceFile)

	// In order to capture the output of "go run", a temporary script is generated and executed by the shell
	// Note: using the script contents in a Command, Run/Start it while capturing Stdout & Stdout via Pipes does not work.
	script := fmt.Sprintf("go run %s > %s  2>&1", goSourceFile, logName)
	buf := new(bytes.Buffer)
	buf.WriteString(script)
	ioutil.WriteFile(temporaryShellScriptName, buf.Bytes(), os.ModePerm)
	// clean up afterwards
	defer os.Remove(temporaryShellScriptName)

	cmd := exec.Command("sh", temporaryShellScriptName)
	runError := cmd.Run()

	// The captured output has been written to a temporary log file
	file, err := os.Open(logName)
	if err != nil {
		return "[rango] open output failed", err
	}
	// make sure it is closed before return (and removal)
	defer file.Close()
	// clean up afterwards
	defer os.Remove(logName)

	logBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "[rango] read output failed", err
	}
	return string(logBytes), runError
}

// buildTemplateVars creates a templateVars struct from the list of code sourceLines
func buildTemplateVars(sourceLines []SourceHolder) templateVars {
	imageVars := new(templateVars)
	for i, each := range sourceLines {
		switch each.Type {
		case Import:
			imageVars.Imports = append(imageVars.Imports, each.Source)
		case Print:
			// only preserve the last
			if i == len(sourceLines)-1 {
				imageVars.Statements = append(imageVars.Statements, each.Source)
			}
		default:
			imageVars.Statements = append(imageVars.Statements, each.Source)
		}
	}
	return *imageVars
}

// templateVars holds the template variables for the Go source to evaluate
type templateVars struct {
	Imports    []string
	Statements []string
}

// imageSourceTemplate returns a Go program template that requires templateVars to produce Go source
func imageSourceTemplate() string {
	return `package main
import "fmt"
{{range .Imports}}{{.}}
{{end}}
func nop(v interface{}){}
func main() {
fmt.Print("")
{{range .Statements}}{{.}}
{{end}}
}
`
}
