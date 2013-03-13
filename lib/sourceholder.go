package lib

import (
	"fmt"
	"strings"
)

const (
	Import = iota
	Statement
	VariableDecl
	Print
)

// SourceHolder is basically one line of Go source code with meta data
type SourceHolder struct {
	LoopCount int    // In which REPL count was this created
	Type      int    // one of the constants Import,...
	Source    string // Go code entered or hidden code produced by rango
	Hidden    bool   // If false then hide this from source listing
	// type data
	ImportNames   []string // If the holder is of type Import then store the packages here
	VariableNames []string // If the holder is of type VariableDecl then store the variable names here
}

// Hide marks a SourceHolder as a hidden line ; they will not show up in source listing
func (s *SourceHolder) Hide() {
	s.Hidden = true
}

// NewImport creates a new SourceHolder of type Import
func NewImport(loopcount int, source string) SourceHolder {
	return SourceHolder{LoopCount: loopcount, Type: Import, Source: source}
}

// NewStatement creates a new SourceHolder of type Statement
func NewStatement(loopcount int, source string) SourceHolder {
	return SourceHolder{LoopCount: loopcount, Type: Statement, Source: source}
}

// NewVariableDecl creates a new SourceHolder of type VariableDecl
func NewVariableDecl(loopcount int, source string) SourceHolder {
	return SourceHolder{
		LoopCount:     loopcount,
		Type:          VariableDecl,
		Source:        source,
		VariableNames: ParseVariableNames(source)}
}

// NewPrint creates a new SourceHolder of type Print
func NewPrint(loopcount int, source string) SourceHolder {
	return SourceHolder{LoopCount: loopcount, Type: Print, Source: source, Hidden: true}
}

// AppendTo adds a new SourceHolder to the collection of entries
// As a side effect, it may produce additional SourceHolders based on its type.
func (s SourceHolder) AppendTo(holders []SourceHolder) []SourceHolder {
	extended := append(holders, s)
	if VariableDecl == s.Type {
		uselessSource := fmt.Sprintf("nop(%s)", s.VariableNames[0])
		useless := NewStatement(s.LoopCount, uselessSource)
		(&useless).Hide()
		extended = append(extended, useless)
	}
	return extended
}

// IsVariable says whether the receiver is known as declared Variable name.
func (s SourceHolder) IsVariable(entry string) bool {
	return s.Type == VariableDecl && s.VariableNames[0] == entry
}

// IsVariableDeclaration says whether the receiver is of such type.
func IsVariableDeclaration(source string) bool {
	assignmentIndex := strings.Index(source, ":=")
	return assignmentIndex > 0
}

// TODO read more than 1
func ParseVariableNames(source string) []string {
	assignmentIndex := strings.Index(source, ":=")
	varName := strings.Trim(source[0:assignmentIndex], " ")
	return []string{varName}
}

// CollectVariables returns the list of declared variable names entered by the user.
func CollectVariables(holders []SourceHolder) []string {
	names := []string{}
	for _, each := range holders {
		if VariableDecl == each.Type {
			names = append(names, each.VariableNames[0])
		}
	}
	return names
}
