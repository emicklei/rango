package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"text/template"
)

func Generate_compile_run(imageName string, entries []SourceHolder) string {
	err := generate(imageName, entries)
	if err != nil {
		return fmt.Sprintf("generate failed:%v", err)
	}
	output, err := gorun(imageName)
	if err != nil {
		return fmt.Sprintf("run failed:%v", err)
	}
	return output
}

func generate(imageName string, entries []SourceHolder) error {
	t := template.Must(template.New("image").Parse(imageSourceTemplate()))
	var sourceBuffer bytes.Buffer
	t.Execute(&sourceBuffer, buildTemplateVars(entries))
	err := ioutil.WriteFile(imageName, sourceBuffer.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}

func gorun(imageName string) (string, error) {
	cmd := exec.Command("go", "run", imageName)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}
	if err = cmd.Start(); err != nil {
		return "", err
	}
	output, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}
	failure, err := ioutil.ReadAll(stderr)
	if err != nil {
		return "", err
	}
	if err := cmd.Wait(); err != nil {
		return "", err
	}
	if len(failure) > 0 {
		return string(failure), nil
	}
	return string(output), nil
}

func buildTemplateVars(entries []SourceHolder) templateVars {
	imageVars := new(templateVars)
	for i, each := range entries {
		switch each.Type {
		case Import:
			imageVars.Imports = append(imageVars.Imports, each.Source)
		case Print:
			// only preserve the last
			if i == len(entries)-1 {
				imageVars.Statements = append(imageVars.Statements, each.Source)
			}
		default:
			imageVars.Statements = append(imageVars.Statements, each.Source)
		}
	}
	return *imageVars
}

type templateVars struct {
	Imports    []string
	Statements []string
}

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
