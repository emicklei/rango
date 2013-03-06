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

type SourceHolder struct {
	LoopCount int
	Type      int
	Source    string
	Hidden    bool
	// type data
	ImportNames   []string
	VariableNames []string
}

func (s *SourceHolder) Hide() {
	s.Hidden = true
}

func NewImport(loopcount int, source string) SourceHolder {
	return SourceHolder{LoopCount: loopcount, Type: Import, Source: source}
}

func NewStatement(loopcount int, source string) SourceHolder {
	return SourceHolder{LoopCount: loopcount, Type: Statement, Source: source}
}

func NewVariableDecl(loopcount int, source string) SourceHolder {
	return SourceHolder{
		LoopCount:     loopcount,
		Type:          VariableDecl,
		Source:        source,
		VariableNames: ParseVariableNames(source)}
}

func NewPrint(loopcount int, source string) SourceHolder {
	return SourceHolder{LoopCount: loopcount, Type: Print, Source: source, Hidden: true}
}

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

func (s SourceHolder) IsVariable(entry string) bool {
	return s.Type == VariableDecl && s.VariableNames[0] == entry
}

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

func CollectVariables(holders []SourceHolder) []string {
	names := []string{}
	for _, each := range holders {
		if VariableDecl == each.Type {
			names = append(names, each.VariableNames[0])
		}
	}
	return names
}
