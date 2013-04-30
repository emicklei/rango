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
func generate_compile_run(imageName string, sourceLines []SourceHolder) (string, error, int) {
	// generate
	gosource := fmt.Sprintf("%s.go", imageName)
	err := generate(gosource, sourceLines)
	if err != nil {
		return fmt.Sprintf("[rango] generate Go source failed"), err, GenerationError
	}
	// build
	command := fmt.Sprintf("go build %s.go", imageName)
	output, err := execCommand(imageName, command)
	if !*DEBUG {
		defer os.Remove(gosource)
	}
	if err != nil {
		return output, err, CompilationError
	}
	// run
	command = fmt.Sprintf("./%s", imageName)
	defer os.Remove(imageName)
	output, err = execCommand(imageName, command)
	if err != nil {
		return output, err, ExecutionError
	}
	// success
	return output, nil, NoError
}

// generate produces a Go source file from a list of Go code sourceLines
func generate(goSourceFile string, sourceLines []SourceHolder) error {
	t := template.Must(template.New("image").Parse(imageSourceTemplate()))
	var sourceBuffer bytes.Buffer
	t.Execute(&sourceBuffer, buildTemplateVars(sourceLines))
	return ioutil.WriteFile(goSourceFile, sourceBuffer.Bytes(), 0644)
}

func execCommand(imageName, command string) (string, error) {
	logName := fmt.Sprintf("%s.exec.log", imageName)

	// In order to capture the output of command, a temporary script is generated and executed by the shell
	// Note: using the script contents in a Command, Run/Start it while capturing Stdout & Stdout via Pipes does not work.
	script := fmt.Sprintf("%s > %s 2>&1", command, logName)
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
			imageVars.Imports = append(imageVars.Imports, &sourceLines[i])
		case Print:
			// only preserve prints of the last entry
			if i == len(sourceLines)-1 {
				imageVars.Statements = append(imageVars.Statements, &sourceLines[i])
			}
		default:
			imageVars.Statements = append(imageVars.Statements, &sourceLines[i])
		}
	}
	// assign line numbers
	lineNumber := 3
	for _, each := range imageVars.Imports {
		each.LineNumber = lineNumber
		lineNumber++
	}
	lineNumber += 3
	for _, each := range imageVars.Statements {
		each.LineNumber = lineNumber
		lineNumber++
	}
	return *imageVars
}

// templateVars holds the template variables for the Go source to evaluate
type templateVars struct {
	Imports    []*SourceHolder
	Statements []*SourceHolder
}

// imageSourceTemplate returns a Go program template that requires templateVars to produce Go source
func imageSourceTemplate() string {
	return `package main
import "fmt"
{{range .Imports}}{{.Source}} 			// {{.LineNumber}}
{{end}}
func rango_first(value ...interface{}) (interface{}) {
	return value[0]
}
func main() {
fmt.Print("")
{{range .Statements}}{{.Source}} 		// {{.LineNumber}}
{{end}}
}
`
}
