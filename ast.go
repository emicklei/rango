package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func IsExpressionStatement(line string) bool {
	node, err := ParseStatement(line)
	if err != nil {
		log("parsing variables failed", err)
		return false
	}
	av := new(AstVisitor)
	ast.Walk(av, node)
	return av.IsExpression
}

// ParseVariables parse the names of variables assigned or declared in a line
func ParseVariables(line string) (assigned []string, declared []string, err error) {
	node, err := ParseStatement(line)
	if err != nil {
		log("parsing variables failed", err)
		return assigned, declared, err
	}
	av := new(AstVisitor)
	ast.Walk(av, node)
	return av.VariablesAssigned, av.VariablesDeclared, nil
}

// ParseImports parse the name of the packages from the import declaration in a line
func ParseImports(line string) ([]string, error) {
	node, err := ParseImport(line)
	if err != nil {
		log("parsing imports failed", err)
		return nil, err
	}
	av := new(AstVisitor)
	ast.Walk(av, node)
	return av.Imports, nil
}

// AstVisitor implements a ast.Visitor and collect variable and import info
type AstVisitor struct {
	VariablesAssigned []string
	VariablesDeclared []string
	Imports           []string
	IsExpression      bool
}

// Visit inspects the type of a Node to detect a Assignment, Declaration or Import
func (av *AstVisitor) Visit(node ast.Node) ast.Visitor {
	switch node.(type) {
	case *ast.AssignStmt:
		for _, each := range node.(*ast.AssignStmt).Lhs {
			av.VariablesAssigned = append(av.VariablesAssigned, each.(*ast.Ident).Name)
		}
	case *ast.DeclStmt:
		for _, each := range node.(*ast.DeclStmt).Decl.(*ast.GenDecl).Specs {
			valueSpec, ok := each.(*ast.ValueSpec)
			if ok {
				for _, other := range valueSpec.Names {
					av.VariablesDeclared = append(av.VariablesDeclared, other.Name)
				}
			}
		}
	case *ast.ImportSpec:
		av.Imports = append(av.Imports, node.(*ast.ImportSpec).Path.Value)
	case *ast.ExprStmt:
		av.IsExpression = true
	}
	return av
}

// ParseStatement is a modified version of go/parser.ParseExpr
func ParseStatement(x string) (ast.Stmt, error) {
	// parse x within the context of a complete package for correct scopes;
	// put x alone on a separate line (handles line comments), followed by a ';'
	// to force an error if the expression is incomplete
	file, err := parser.ParseFile(token.NewFileSet(), "", "package p;func _(){\n"+x+"\n;}", 0)
	if err != nil {
		return nil, err
	}
	return file.Decls[0].(*ast.FuncDecl).Body.List[0], nil
}

// ParseImport is a modified version of go/parser.ParseExpr
func ParseImport(x string) (ast.Node, error) {
	// parse x within the context of a complete package for correct scopes;
	// put x alone on a separate line (handles line comments), followed by a ';'
	// to force an error if the expression is incomplete
	file, err := parser.ParseFile(token.NewFileSet(), "", "package p\n"+x+"\n", 0)
	if err != nil {
		return nil, err
	}
	return file.Decls[0], nil
}
