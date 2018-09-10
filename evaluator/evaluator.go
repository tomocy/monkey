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
	rightVal := Eval(exp)
	switch operator {
	case "!":
		return evalBang(rightVal)
	case "-":
		return evalMinusPrefix(rightVal)
	default:
		return nullObj
	}
}

func evalInfix(leftExp ast.Expression, operator string, rightExp ast.Expression) object.Object {
	leftVal := Eval(leftExp)
	rightVal := Eval(rightExp)
	switch {
	case leftVal.Type() == object.Integer && rightVal.Type() == object.Integer:
		leftObj := leftVal.(*object.IntegerObject)
		rightObj := rightVal.(*object.IntegerObject)
		return evalInfixOfInteger(leftObj, operator, rightObj)
	default:
		return nullObj
	}
}

func evalInfixOfInteger(leftObj *object.IntegerObject, operator string, rightObj *object.IntegerObject) object.Object {
	switch operator {
	case "+":
		return &object.IntegerObject{Value: leftObj.Value + rightObj.Value}
	case "-":
		return &object.IntegerObject{Value: leftObj.Value - rightObj.Value}
	case "*":
		return &object.IntegerObject{Value: leftObj.Value * rightObj.Value}
	case "/":
		return &object.IntegerObject{Value: leftObj.Value / rightObj.Value}
	default:
		return nullObj
	}
}

func evalBang(rightVal object.Object) object.Object {
	switch rightVal {
	case trueObj:
		return falseObj
	case falseObj, nullObj:
		return trueObj
	default:
		return falseObj
	}
}

func evalMinusPrefix(rightVal object.Object) object.Object {
	if rightVal.Type() != object.Integer {
		return nullObj
	}

	val := rightVal.(*object.IntegerObject).Value
	return &object.IntegerObject{Value: -val}
}

func convertToBooleanObject(b bool) object.Object {
	if b {
		return trueObj
	}

	return falseObj
}
