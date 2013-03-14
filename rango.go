// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	GenerateCompileRun = iota
	UpdateSourceOnly
)

var (
	Stdin       *bufio.Reader
	imageName   = "generated_by_rango"
	sourceLines []SourceHolder
	entryCount  int
)

func init() {
	Stdin = bufio.NewReader(os.Stdin)
	sourceLines = []SourceHolder{}
}

func main() {
	welcome()
	if len(os.Args) > 1 {
		imageName = os.Args[1]
		processChanges()
	}
	loop()
}

func welcome() {
	//	fmt.Println("[rango] .q = quit, .v = variables, .s = source, .u = undo, #<statement> = execute once, .? = more help")
	fmt.Println("[rango] .q = quit, .v = variables, .s = source")
}

func processChanges() {
	changesName := fmt.Sprintf("%s.changes", imageName)
	file, err := os.Open(changesName)
	if err != nil {
		return
	}
	defer file.Close()
	in := bufio.NewReader(file)
	for {
		entered, err := in.ReadString('\n')
		if len(entered) > 0 {
			handleSource(strings.TrimRight(entered, "\n"), UpdateSourceOnly) // without newline
		}
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("error reading file ", err)
			break
		}
	}
	fmt.Printf("[rango] processed %d lines from %s\n", entryCount+1, changesName)
}

func loop() {
	for {
		//fmt.Printf("current entryCount:%d\n", entryCount)
		fmt.Print("> ")
		in := bufio.NewReader(os.Stdin)
		entered, err := in.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		output := dispatch(entered[:len(entered)-1]) // without newline
		if len(output) > 0 {
			fmt.Println(output)
		}
	}
}

func dispatch(entry string) string {
	if len(entry) == 0 {
		return entry
	}
	switch {
	case strings.HasPrefix(entry, ".v"):
		fmt.Printf("%v\n", CollectVariables(sourceLines))
		return ""
	case strings.HasPrefix(entry, ".q"):
		os.Exit(0)
	case strings.HasPrefix(entry, ".s"):
		return handlePrintSource()
	case strings.HasPrefix(entry, ".u"):
		return handleUndo()
	case strings.HasPrefix(entry, "#"):
		// TODO forget sources for current entryCount
		return handleSource(entry[1:], GenerateCompileRun)
	case strings.HasPrefix(entry, "."):
		return handleUnknownCommand(entry)
	case isVariable(entry):
		return handlePrintVariable(entry)
	}
	return handleSource(entry, GenerateCompileRun)
}

func handleUndo() string {
	undo(entryCount)
	return handlePrintSource()
}

func handleSource(entry string, mode int) string {
	if len(entry) == 0 {
		return entry
	}
	entryCount++
	if strings.HasPrefix(entry, "import ") {
		return handleImport(entry)
	}
	if IsVariableDeclaration(entry) {
		vardecl := NewVariableDecl(entryCount, entry)
		addEntry(vardecl)
		// copied from PrintVariable
		printEntry := fmt.Sprintf("fmt.Printf(\"%%v\",%s)", vardecl.VariableNames[0])
		addEntry(NewPrint(entryCount, printEntry))
	} else {
		addEntry(NewStatement(entryCount, entry))
	}
	if UpdateSourceOnly == mode {
		return ""
	}
	return Generate_compile_run(fmt.Sprintf("%s.go", imageName), sourceLines)
}

func handlePrintSource() string {
	var buf bytes.Buffer
	line := 1
	// First imports then functions then main statements
	for _, each := range sourceLines {
		if (Import == each.Type) && !each.Hidden {
			if line > 1 {
				buf.WriteString("\n")
			}
			buf.WriteString(fmt.Sprintf("%  d:\t%s", line, each.Source))
			line++
		}
	}
	for _, each := range sourceLines {
		if (Statement == each.Type || VariableDecl == each.Type) && !each.Hidden {
			if line > 1 {
				buf.WriteString("\n")
			}
			buf.WriteString(fmt.Sprintf("%  d:\t%s", line, each.Source))
			line++
		}
	}
	return string(buf.Bytes())
}

func handleUnknownCommand(entry string) string {
	return fmt.Sprintf("[rango] \"%s\": command not found", entry)
}

func handlePrintVariable(varname string) string {
	printEntry := fmt.Sprintf("fmt.Printf(\"%%v\",%s)", varname)
	addEntry(NewPrint(entryCount, printEntry))
	return Generate_compile_run(fmt.Sprintf("%s.go", imageName), sourceLines)
}

// handleImport adds a non-existing import package.
// Source will be updated on the next statement.
func handleImport(entry string) string {
	sourceLines = NewImport(entryCount, entry).AppendTo(sourceLines)
	return ""
}

func addEntry(holder SourceHolder) {
	//fmt.Printf("%#v\n", holder)
	sourceLines = holder.AppendTo(sourceLines)
}

func isVariable(entry string) bool {
	for _, each := range sourceLines {
		if each.IsVariable(entry) {
			return true
		}
	}
	return false
}

// undo removes sourceLines appended
func undo(until int) {
	for {
		if len(sourceLines) == 0 {
			fmt.Println("(no go source)")
			break
		}
		last := sourceLines[len(sourceLines)-1]
		if last.EntryCount < until {
			// set new entry count
			entryCount = last.EntryCount
			break
		}
		sourceLines = sourceLines[:len(sourceLines)-1]
	}
}
