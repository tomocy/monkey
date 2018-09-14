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

func Eval(node ast.Node, env *Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Value, env)
	case *ast.BlockStatement:
		return evalBlockStatements(node, env)
	case *ast.LetStatement:
		return evalLetStatement(node, env)
	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)
	case *ast.If:
		return evalIf(node, env)
	case *ast.Prefix:
		return evalPrefix(node.Operator, node.RightValue, env)
	case *ast.Infix:
		return evalInfix(node.LeftValue, node.Operator, node.RightValue, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.Integer:
		return &object.IntegerObject{Value: node.Value}
	case *ast.Boolean:
		return convertToBooleanObject(node.Value)
	}

	return nullObj
}

func evalProgram(program *ast.Program, env *Environment) object.Object {
	var obj object.Object
	for _, stmt := range program.Statements {
		obj = Eval(stmt, env)
		if obj.Type() == object.Error {
			return obj
		}
		if obj.Type() == object.Return {
			return obj.(*object.ReturnObject).Value
		}
	}

	return obj
}

func evalBlockStatements(blockStmt *ast.BlockStatement, env *Environment) object.Object {
	var obj object.Object
	for _, stmt := range blockStmt.Statements {
		obj = Eval(stmt, env)
		if obj.Type() == object.Return {
			return obj
		}
		if obj.Type() == object.Error {
			return obj
		}
	}

	return obj
}

func evalLetStatement(node *ast.LetStatement, env *Environment) object.Object {
	obj := Eval(node.Value, env)
	if obj.Type() == object.Error {
		return obj
	}

	env.Set(node.Ident.Value, obj)

	return obj
}

func evalReturnStatement(node *ast.ReturnStatement, env *Environment) object.Object {
	obj := Eval(node.Value, env)
	if obj.Type() == object.Error {
		return obj
	}

	return &object.ReturnObject{Value: obj}
}

func evalIf(node *ast.If, env *Environment) object.Object {
	condition := Eval(node.Condition, env)
	if condition.Type() == object.Error {
		return condition
	}

	if isTruthy(condition) {
		return Eval(node.Consequence, env)
	}

	if node.Alternative != nil {
		return Eval(node.Alternative, env)
	}

	return nullObj
}

func isTruthy(obj object.Object) bool {
	return obj != falseObj && obj != nullObj
}

func evalPrefix(operator string, exp ast.Expression, env *Environment) object.Object {
	rightObj := Eval(exp, env)
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

func evalInfix(leftExp ast.Expression, operator string, rightExp ast.Expression, env *Environment) object.Object {
	leftObj := Eval(leftExp, env)
	if leftObj.Type() == object.Error {
		return leftObj
	}

	rightObj := Eval(rightExp, env)
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

func evalIdentifier(node *ast.Identifier, env *Environment) object.Object {
	obj, ok := env.Get(node.Value)
	if !ok {
		return newError("unknown identifier: %s", node.Value)
	}

	return obj
}

func convertToBooleanObject(b bool) object.Object {
	if b {
		return trueObj
	}

	return falseObj
}

func newError(format string, a ...interface{}) object.Object {
	return &object.ErrorObject{Message: fmt.Sprintf(format, a...)}
}
