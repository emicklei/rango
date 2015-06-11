// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	GenerateCompileRun = iota
	UpdateSourceOnly

	//  gen-com-run errors
	GenerationError
	CompilationError
	ExecutionError
	NoError

	ShowLineNumbers = true
)

var (
	Stdin       *bufio.Reader
	imageName   = "generated_by_rango"
	sourceLines []SourceHolder
	entryCount  int
	logChanges  = false
	// debug option
	DEBUG = flag.Bool("debug", false, "produce more output")
)

func init() {
	flag.Parse()
	Stdin = bufio.NewReader(os.Stdin)
	sourceLines = []SourceHolder{}
}

func main() {
	welcome()
	if len(os.Args) > 1 { // interpret the last arg as projectname, unless it was an option
		imageName = os.Args[len(os.Args)-1]
		if !strings.HasPrefix(imageName, "-") {
			processChanges()
			logChanges = true
		}
	}
	loop()
}

func welcome() {
	fmt.Println(handleHelp())
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
		return handlePrintSource(ShowLineNumbers)
	case strings.HasPrefix(entry, ".u"):
		return handleUndo()
	case strings.HasPrefix(entry, ".?"):
		return handleHelp()
	case strings.HasPrefix(entry, "="):
		return handlePrintExpressionValue(entry[1:])
	case strings.HasPrefix(entry, "!"):
		wantsLog := logChanges
		logChanges = false
		out := handleSource(entry[1:], GenerateCompileRun)
		undo(entryCount)
		logChanges = wantsLog
		// TODO if handleSource failed then undo is done twice =>  Need (string,error) return values!
		return out
	case strings.HasPrefix(entry, "."):
		return handleUnknownCommand(entry)
	}
	return handleSource(entry, GenerateCompileRun)
}

func handleHelp() string {
	return "[rango] .q = quit, !<source> = eval once , =<source> = print once, .v = variables, .s = source, .u = undo, .? = help"
}

func handleUndo() string {
	undo(entryCount)
	if logChanges {
		dumpChanges()
	}
	return handlePrintSource(ShowLineNumbers)
}

func handleSource(entry string, mode int) string {
	if strings.HasPrefix(entry, "import") {
		return handleImport(entry)
	}
	if isVariable(entry) {
		if GenerateCompileRun == mode {
			return handlePrintExpressionValue(entry)
		} else {
			return ""
		}
	}
	assigned, declared, err := ParseVariables(entry)
	if err != nil { // error is already printed
		return ""
	}
	entryCount++
	if len(assigned) > 0 {
		handleVariableAssignments(assigned, entry)
	}
	if len(declared) > 0 {
		handleVariableDeclarations(declared, entry)
	}
	// TODO handle mix in one entry
	if len(assigned) == 0 && len(declared) == 0 {
		addEntry(NewStatement(entryCount, entry))
	}
	if UpdateSourceOnly == mode {
		return ""
	}
	output, err, kind := generate_compile_run(imageName, sourceLines)
	if err != nil {
		// output has reason for failure
		undo(entryCount)
		// if compiler error then parse it to produce better output
		if CompilationError == kind {
			output = prepareCompilerErrorOutput(output)
		}
	} else {
		if logChanges {
			dumpChanges()
		}
	}
	return output
}

func handleVariableAssignments(names []string, entry string) {
	// detect if assign+decl
	areAssignmentsOnly := true
	for _, each := range names {
		if !isVariable(each) {
			areAssignmentsOnly = false
			break
		}
	}
	if areAssignmentsOnly {
		addEntry(NewVariableAssign(entryCount, entry, names))
	} else {
		// it is a combi, handle as decl
		addEntry(NewVariableDecl(entryCount, entry, names))
	}
	handlePrintVariableValues(names)
}
func handleVariableDeclarations(names []string, entry string) {
	addEntry(NewVariableDecl(entryCount, entry, names))
	handlePrintVariableValues(names)
}

func handlePrintVariableValues(names []string) {
	// fmt.Printf( "%v,%v,%v", a , b ,c )
	var buf bytes.Buffer
	buf.WriteString("fmt.Printf(\"")
	for i := 0; i < len(names); i++ {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString("%v")
	}
	buf.WriteString("\"")
	for _, each := range names {
		buf.WriteString(",")
		buf.WriteString(each)
	}
	buf.WriteString(")")
	addEntry(NewPrint(entryCount, string(buf.Bytes())))
}

// handlePrintSource list the statements into a single string.
// Line numbers are optional
func handlePrintSource(withLineNumbers bool) string {
	var buf bytes.Buffer
	line := 1
	// First imports then functions then main statements
	for _, each := range sourceLines {
		if (Import == each.Type) && !each.Hidden {
			if line > 1 {
				buf.WriteString("\n")
			}
			if withLineNumbers {
				buf.WriteString(fmt.Sprintf("%  d:\t%s", line, each.Source))
			} else {
				buf.WriteString(each.Source)
			}
			line++
		}
	}
	for _, each := range sourceLines {
		if (Statement == each.Type ||
			VariableDecl == each.Type ||
			VariableAssign == each.Type) && !each.Hidden {
			if line > 1 {
				buf.WriteString("\n")
			}
			if withLineNumbers {
				buf.WriteString(fmt.Sprintf("%  d:\t%s", line, each.Source))
			} else {
				buf.WriteString(each.Source)
			}
			line++
		}
	}
	return string(buf.Bytes())
}

// handleUnknownCommand is called when the entry did not match a known command
func handleUnknownCommand(entry string) string {
	return fmt.Sprintf("[rango] \"%s\": command not found", entry)
}

// handlePrintExpressionValue adds a print statement to display the value of an expression
func handlePrintExpressionValue(expression string) string {
	printEntry := fmt.Sprintf("fmt.Printf(\"%%v\",rango_first(%s))", expression)
	addEntry(NewPrint(entryCount, printEntry))
	output, _, _ := generate_compile_run(imageName, sourceLines)
	// no need to rollback entry
	return output
}

// handleImport adds a non-existing import package.
// Source will be updated on the next statement.
func handleImport(entry string) string {
	names, err := ParseImports(entry)
	if err != nil { // error is already printed
		return ""
	}
	entryCount++
	sourceLines = NewImport(entryCount, entry, names).AppendTo(sourceLines)
	return ""
}

func prepareCompilerErrorOutput(output string) string {
	var buf bytes.Buffer
	lines := strings.Split(output, "\n")
	written := false
	for _, each := range lines {
		if written {
			buf.WriteString("\n")
		}
		if len(each) > 0 && !strings.HasPrefix(each, "#") {
			written = true
			// ./generated_by_rango.go:9: undefined: b
			// TODO scan for the line number and find the matching sourceLine (first)
			buf.WriteString(each)
		}
	}
	return string(buf.Bytes())
}

func addEntry(holder SourceHolder) {
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

func log(what string, err error) {
	fmt.Printf("[rango] %s : %v\n", what, err)
}
