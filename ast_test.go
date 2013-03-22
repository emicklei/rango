package main

import (
	"fmt"
	"go/ast"
	"testing"
)

var lines = []struct {
	line     string
	assigned []string
	declared []string
}{
	{"a := 1", []string{"a"}, []string{}},
	{"b = 2", []string{"b"}, []string{}},
	{"c,d = 3,4", []string{"c", "d"}, []string{}},
	{"var ( a = 'a' ; b int ) ", []string{}, []string{"a", "b"}},
	{"fmt.Println(\"here\")", []string{}, []string{}},
}

func TestParseVariables(t *testing.T) {
	for i, each := range lines {
		node, err := ParseStatement(each.line)
		if err != nil {
			panic(err)
		}
		av := new(AstVisitor)
		ast.Walk(av, node)
		fmt.Printf("assigned:%v\n", av.VariablesAssigned)
		fmt.Printf("declared:%v\n", av.VariablesDeclared)

		if !equal(each.assigned, av.VariablesAssigned) {
			t.Fatal("i=", i)
		}
		if !equal(each.declared, av.VariablesDeclared) {
			t.Fatal("i=", i)
		}
	}
}

func TestParseImports(t *testing.T) {
	e := "import ( \"bufio\" ; \"bytes\" )"
	node, err := ParseImport(e)
	if err != nil {
		panic(err)
	}
	av := new(AstVisitor)
	ast.Walk(av, node)
	if len(av.Imports) != 2 {
		t.Fail()
	}
}

func equal(one []string, other []string) bool {
	if len(one) != len(other) {
		return false
	}
	for i, left := range one {
		right := other[i]
		if left != right {
			return false
		}
	}
	return true
}
