package main

// Other implementation
// go-repl
// https://gitorious.org/golang/go-repl/blobs/master/main.go

import (
	"bufio"
	"bytes"
	"fmt"
	// "github.com/emicklei/hopwatch"
	. "github.com/emicklei/rango/lib"
	"io"
	"os"
	"strings"
)

const (
	GenerateCompileRun = iota
	UpdateSourceOnly
)

var (
	Stdin         *bufio.Reader
	imageName     = "chameleon"
	entries       []SourceHolder
	lastLoopCount int
	loopCount     int
)

func init() {
	Stdin = bufio.NewReader(os.Stdin)
	entries = []SourceHolder{}
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
	fmt.Println("[rango] (.q = quit, .v = variables, .s = source, .u = undo, #<statement> = do not remember statement, .? = more help)")
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
			loopCount++
		}
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("error reading file ", err)
			break
		}
	}
	fmt.Printf("[rango] processed %d lines from %s\n", loopCount+1, changesName)
}

func loop() {
	for {
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
		loopCount++
	}
}

func dispatch(entry string) string {
	if len(entry) == 0 {
		return entry
	}
	switch {
	case strings.HasPrefix(entry, ".v"):
		fmt.Printf("%v\n", CollectVariables(entries))
		return ""
	case strings.HasPrefix(entry, ".q"):
		os.Exit(0)
	case strings.HasPrefix(entry, ".s"):
		return handlePrintSource()
	case strings.HasPrefix(entry, ".u"):
		return handleUndo()
	case strings.HasPrefix(entry, "#"):
		// forget sources for current loopcount
		return handleSource(entry[1:], GenerateCompileRun)
	case strings.HasPrefix(entry, "."):
		return handleUnknownCommand(entry)
	case isVariable(entry):
		return handlePrintVariable(entry)
	}
	return handleSource(entry, GenerateCompileRun)
}

func handleUndo() string {
	undo(loopCount - 2) // loop already incremented
	loopCount -= 2
	return handlePrintSource()
}

func handleSource(entry string, mode int) string {
	if len(entry) == 0 {
		return entry
	}
	if strings.HasPrefix(entry, "import ") {
		return handleImport(entry)
	}
	if IsVariableDeclaration(entry) {
		vardecl := NewVariableDecl(loopCount, entry)
		addEntry(vardecl)
		// copied from PrintVariable
		printEntry := fmt.Sprintf("fmt.Printf(\"%%v\",%s)", vardecl.VariableNames[0])
		addEntry(NewPrint(loopCount, printEntry))
	} else {
		addEntry(NewStatement(loopCount, entry))
	}
	if UpdateSourceOnly == mode {
		return ""
	}
	return Generate_compile_run(fmt.Sprintf("%s.go", imageName), entries)
}

func handlePrintSource() string {
	var buf bytes.Buffer
	line := 1
	for _, each := range entries {
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
	addEntry(NewPrint(loopCount, printEntry))
	return Generate_compile_run(fmt.Sprintf("%s.go", imageName), entries)
}

// handleImport adds a non-existing import package.
// Source will be updated on the next statement.
func handleImport(entry string) string {
	entries = NewImport(loopCount, entry).AppendTo(entries)
	return ""
}

func addEntry(holder SourceHolder) {
	//fmt.Printf("%#v\n", holder)
	entries = holder.AppendTo(entries)
}

func isVariable(entry string) bool {
	for _, each := range entries {
		if each.IsVariable(entry) {
			return true
		}
	}
	return false
}

// undo removes entries appended
func undo(until int) {
	for {
		if len(entries) == 0 {
			break
		}
		last := entries[len(entries)-1]
		if until == last.LoopCount {
			break
		}
		// fmt.Printf("removed:%s\n", last.Source)
		entries = entries[:len(entries)-1]
	}
}
