// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

const (
	Import = iota
	Statement
	VariableDecl
	Print
)

// SourceHolder is basically one line of Go source code with meta data
type SourceHolder struct {
	EntryCount int    // In which REPL count was this created
	Type       int    // one of the constants Import,...
	Source     string // Go code entered or hidden code produced by rango
	Hidden     bool   // If false then hide this from source listing
	LineNumber int    // The exact line number in the generated Go source ; used for compiler error reporting
	// type data
	PackageNames  []string // If the holder is of type Import then store the packages here
	VariableNames []string // If the holder is of type VariableDecl then store the variable names here
}

// Hide marks a SourceHolder as a hidden line ; they will not show up in source listing
func (s *SourceHolder) Hide() {
	s.Hidden = true
}

// NewImport creates a new SourceHolder of type Import
func NewImport(entryCount int, source string, packageNames []string) SourceHolder {
	return SourceHolder{EntryCount: entryCount, Type: Import, Source: source, PackageNames: packageNames}
}

// NewStatement creates a new SourceHolder of type Statement
func NewStatement(entryCount int, source string) SourceHolder {
	return SourceHolder{EntryCount: entryCount, Type: Statement, Source: source}
}

// NewVariableDecl creates a new SourceHolder of type VariableDecl
func NewVariableDecl(entryCount int, source string, names []string) SourceHolder {
	return SourceHolder{
		EntryCount:    entryCount,
		Type:          VariableDecl,
		Source:        source,
		VariableNames: names}
}

// NewPrint creates a new SourceHolder of type Print
func NewPrint(entryCount int, source string) SourceHolder {
	return SourceHolder{EntryCount: entryCount, Type: Print, Source: source, Hidden: true}
}

// AppendTo adds a new SourceHolder to the collection of entries.
// As a side effect, it may produce additional SourceHolders based on its type.
func (s SourceHolder) AppendTo(sourceLines []SourceHolder) []SourceHolder {
	extended := append(sourceLines, s)
	if VariableDecl == s.Type {
		for _, each := range s.VariableNames {
			uselessSource := fmt.Sprintf("_ = %s", each)
			useless := NewStatement(s.EntryCount, uselessSource)
			(&useless).Hide()
			extended = append(extended, useless)
		}
	}
	return extended
}

// IsVariable says whether the receiver is known as declared Variable name.
func (s SourceHolder) IsVariable(entry string) bool {
	if s.Type != VariableDecl {
		return false
	}
	for _, each := range s.VariableNames {
		if each == entry {
			return true
		}
	}
	return false
}

// CollectVariables returns the list of declared variable names entered by the user.
func CollectVariables(sourceLines []SourceHolder) []string {
	names := []string{}
	for _, each := range sourceLines {
		if VariableDecl == each.Type {
			names = append(names, each.VariableNames[0])
		}
	}
	return names
}
