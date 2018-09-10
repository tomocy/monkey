package evaluator

import (
	"github.com/tomocy/monkey/ast"
	"github.com/tomocy/monkey/object"
)

var (
	trueObj  = &object.Boolean{Value: true}
	falseObj = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Value)
	case *ast.Integer:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return convertToBooleanObject(node.Value)
	}

	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var obj object.Object
	for _, stmt := range stmts {
		obj = Eval(stmt)
	}

	return obj
}

func convertToBooleanObject(b bool) object.Object {
	if b {
		return trueObj
	}

	return falseObj
}
