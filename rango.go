package main

// Other implementation
// go-repl
// https://gitorious.org/golang/go-repl/blobs/master/main.go

import (
	"bufio"
	"bytes"
	"fmt"
	. "github.com/emicklei/rango/lib"
	"os"
	"strings"
)

var (
	Stdin         *bufio.Reader
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
	loop()
}

func welcome() {
	fmt.Println("[rango] Go ahead (.q = quit, .v = variables, .s = source, .u = undo, #<statement> = do not remember statement, .? = more help)")
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
	case strings.HasPrefix(entry, "import "):
		return handleImport(entry)
	case strings.HasPrefix(entry, ".s"):
		return handlePrintSource()
	case isVariable(entry):
		return handlePrintVariable(entry)
	case strings.HasPrefix(entry, "#"):
		// forget sources for current loopcount
		return dispatch(entry[1:])
	case strings.HasPrefix(entry, "."):
		return handleUnknown(entry)
	}
	return handleStatement(entry)
}

func handlePrintSource() string {
	var buf bytes.Buffer
	for i, each := range entries {
		if (Statement == each.Type || VariableDecl == each.Type) && !each.Hidden {
			if i > 0 {
				buf.WriteString("\n")
			}
			buf.WriteString(fmt.Sprintf("%  d:\t%s", i, each.Source))
		}
	}
	return string(buf.Bytes())
}

func handleUnknown(entry string) string {
	return fmt.Sprintf("[rango] \"%s\": command not found", entry)
}

func handlePrintVariable(varname string) string {
	printEntry := fmt.Sprintf("fmt.Printf(\"%%v\",%s)", varname)
	addEntry(NewPrint(loopCount, printEntry))
	return Generate_compile_run("image.go", entries)
}

// handleImport adds a non-existing import package.
// Source will be updated on the next statement.
func handleImport(entry string) string {
	entries = NewImport(loopCount, entry).AppendTo(entries)
	return ""
}
func handleStatement(entry string) string {
	if IsVariableDeclaration(entry) {
		addEntry(NewVariableDecl(loopCount, entry))
	} else {
		addEntry(NewStatement(loopCount, entry))
	}
	return Generate_compile_run("image.go", entries)
}

func addEntry(holder SourceHolder) {
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
		entries = entries[:len(entries)-1]
	}
}
