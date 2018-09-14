package evaluator

import (
	"fmt"

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
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Value)
	case *ast.BlockStatement:
		return evalBlockStatements(node)
	case *ast.ReturnStatement:
		return evalReturnStatement(node)
	case *ast.If:
		return evalIf(node)
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

func evalProgram(program *ast.Program) object.Object {
	var obj object.Object
	for _, stmt := range program.Statements {
		obj = Eval(stmt)
		if obj.Type() == object.Error {
			return obj
		}
		if obj.Type() == object.Return {
			return obj.(*object.ReturnObject).Value
		}
	}

	return obj
}

func evalBlockStatements(blockStmt *ast.BlockStatement) object.Object {
	var obj object.Object
	for _, stmt := range blockStmt.Statements {
		obj = Eval(stmt)
		if obj.Type() == object.Return {
			return obj
		}
		if obj.Type() == object.Error {
			return obj
		}
	}

	return obj
}

func evalReturnStatement(exp *ast.ReturnStatement) object.Object {
	obj := Eval(exp.Value)
	if obj.Type() == object.Error {
		return obj
	}

	return &object.ReturnObject{Value: obj}
}

func evalIf(ifExp *ast.If) object.Object {
	condition := Eval(ifExp.Condition)
	if condition.Type() == object.Error {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ifExp.Consequence)
	}

	if ifExp.Alternative != nil {
		return Eval(ifExp.Alternative)
	}

	return nullObj
}

func isTruthy(obj object.Object) bool {
	return obj != falseObj && obj != nullObj
}

func evalPrefix(operator string, exp ast.Expression) object.Object {
	rightObj := Eval(exp)
	if rightObj.Type() == object.Error {
		return rightObj
	}

	switch operator {
	case "!":
		return evalBang(rightObj)
	case "-":
		return evalMinusPrefix(rightObj)
	default:
		return newError("unknown operation: %s%s", operator, rightObj.Type())
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
		return newError("unknown operation: -%s", rightObj.Type())
	}

	rightVal := rightObj.(*object.IntegerObject).Value
	return &object.IntegerObject{Value: -rightVal}
}

func evalInfix(leftExp ast.Expression, operator string, rightExp ast.Expression) object.Object {
	leftObj := Eval(leftExp)
	if leftObj.Type() == object.Error {
		return leftObj
	}

	rightObj := Eval(rightExp)
	if rightObj.Type() == object.Error {
		return rightObj
	}

	switch {
	case leftObj.Type() == object.Integer && rightObj.Type() == object.Integer:
		return evalInfixOfInteger(leftObj, operator, rightObj)
	case operator == "==":
		return convertToBooleanObject(leftObj == rightObj)
	case operator == "!=":
		return convertToBooleanObject(leftObj != rightObj)
	default:
		return newError("unknown operation: %s %s %s", leftObj.Type(), operator, rightObj.Type())
	}
}

func newError(format string, a ...interface{}) object.Object {
	return &object.ErrorObject{Message: fmt.Sprintf(format, a...)}
}

func evalInfixOfInteger(leftObj object.Object, operator string, rightObj object.Object) object.Object {
	if leftObj.Type() != object.Integer || rightObj.Type() != object.Integer {
		return newError("unknown operation: %s %s %s", leftObj.Type(), operator, rightObj.Type())
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
	case "<":
		return convertToBooleanObject(leftVal < rightVal)
	case ">":
		return convertToBooleanObject(leftVal > rightVal)
	case "==":
		return convertToBooleanObject(leftVal == rightVal)
	case "!=":
		return convertToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operation: %s %s %s", leftObj.Type(), operator, rightObj.Type())
	}
}

func convertToBooleanObject(b bool) object.Object {
	if b {
		return trueObj
	}

	return falseObj
}
