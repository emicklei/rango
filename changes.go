package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func processChanges() {
	changesName := fmt.Sprintf("%s.changes", imageName)
	file, err := os.Open(changesName)
	if err != nil {
		// ignore missing changes file
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
			log("error reading changes file ", err)
			break
		}
	}
	fmt.Println(handlePrintSource(ShowLineNumbers))
}

func dumpChanges() {
	changesName := fmt.Sprintf("%s.changes", imageName)
	file, err := os.Create(changesName)
	if err != nil {
		log("error creating changes file ", err)
		return
	}
	defer file.Close()
	out := bufio.NewWriter(file)
	_, err = out.WriteString(handlePrintSource(!ShowLineNumbers))
	if err != nil {
		log("error writing changes file ", err)
		return
	}
	out.Flush()
}
