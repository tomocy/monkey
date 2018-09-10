package evaluator

import (
	"github.com/tomocy/monkey/ast"
	"github.com/tomocy/monkey/object"
)

var (
	nullObj  = &object.NullObject{}
	trueObj  = &object.BooleanObject{Value: true}
	falseObj = &object.BooleanObject{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Value)
	case *ast.Prefix:
		return evalPrefix(node.Operator, node.RightValue)
	case *ast.Infix:
		return evalInfix(node.LeftValue, node.Operator, node.RightValue)
	case *ast.Integer:
		return &object.IntegerObject{Value: node.Value}
	case *ast.Boolean:
		return convertToBooleanObject(node.Value)
	}

	return nullObj
}

func evalStatements(stmts []ast.Statement) object.Object {
	var obj object.Object
	for _, stmt := range stmts {
		obj = Eval(stmt)
	}

	return obj
}

func evalPrefix(operator string, exp ast.Expression) object.Object {
	rightObj := Eval(exp)
	switch operator {
	case "!":
		return evalBang(rightObj)
	case "-":
		return evalMinusPrefix(rightObj)
	default:
		return nullObj
	}
}

func evalBang(rightObj object.Object) object.Object {
	switch rightObj {
	case trueObj:
		return falseObj
	case falseObj, nullObj:
		return trueObj
	default:
		return falseObj
	}
}

func evalMinusPrefix(rightObj object.Object) object.Object {
	if rightObj.Type() != object.Integer {
		return nullObj
	}

	val := rightObj.(*object.IntegerObject).Value
	return &object.IntegerObject{Value: -val}
}

func evalInfix(leftExp ast.Expression, operator string, rightExp ast.Expression) object.Object {
	leftObj := Eval(leftExp)
	rightObj := Eval(rightExp)
	switch {
	case leftObj.Type() == object.Integer && rightObj.Type() == object.Integer:
		return evalInfixOfInteger(leftObj, operator, rightObj)
	default:
		return nullObj
	}
}

func evalInfixOfInteger(leftObj object.Object, operator string, rightObj object.Object) object.Object {
	if leftObj.Type() != object.Integer || rightObj.Type() != object.Integer {
		return nullObj
	}
	
	leftVal := leftObj.(*object.IntegerObject).Value
	rightVal := rightObj.(*object.IntegerObject).Value
	switch operator {
	case "+":
		return &object.IntegerObject{Value: leftVal + rightVal}
	case "-":
		return &object.IntegerObject{Value: leftVal - rightVal}
	case "*":
		return &object.IntegerObject{Value: leftVal * rightVal}
	case "/":
		return &object.IntegerObject{Value: leftVal / rightVal}
	default:
		return nullObj
	}
}

func convertToBooleanObject(b bool) object.Object {
	if b {
		return trueObj
	}

	return falseObj
}
