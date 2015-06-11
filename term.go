package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/GeertJohan/go.linenoise"
)

var lastHistoryEntry string

func loop() {
	linenoise.LoadHistory(".rango-history")
	for {
		entered, err := linenoise.Line("> ")
		if err != nil {
			if err == linenoise.KillSignalError {
				os.Exit(0)
			}
			fmt.Println("Unexpected error: %s", err)
			os.Exit(0)
		}
		entry := strings.TrimLeft(entered, "\t ") // without tabs,spaces
		var output string
		if entry != lastHistoryEntry {
			err = linenoise.AddHistory(entry)
			if err != nil {
				fmt.Printf("error: %s\n", entry)
			}
			lastHistoryEntry = entry
			linenoise.SaveHistory(".rango-history")
		}
		output = dispatch(entry)
		if len(output) > 0 {
			fmt.Println(output)
		}
	}
}
