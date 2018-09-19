package evaluator

import (
	"fmt"

	"github.com/tomocy/monkey/ast"
	"github.com/tomocy/monkey/object"
	"github.com/tomocy/monkey/token"
)

var (
	nullObj  = &object.NullObject{}
	trueObj  = &object.BooleanObject{Value: true}
	falseObj = &object.BooleanObject{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
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
		return evalPrefix(node, env)
	case *ast.Infix:
		return evalInfix(node, env)
	case *ast.Function:
		return evalFunction(node, env)
	case *ast.FunctionCall:
		if isQuoteCall(node) {
			return evalQuote(node.Arguments[0], env)
		}
		return evalFunctionCall(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.Integer:
		return evalInteger(node)
	case *ast.Boolean:
		return evalBoolean(node)
	case *ast.String:
		return evalString(node)
	case *ast.Array:
		return evalArray(node, env)
	case *ast.Hash:
		return evalHash(node, env)
	case *ast.Subscript:
		return evalSubscript(node, env)
	}

	return nullObj
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
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

func evalBlockStatements(blockStmt *ast.BlockStatement, env *object.Environment) object.Object {
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

func evalLetStatement(node *ast.LetStatement, env *object.Environment) object.Object {
	obj := Eval(node.Value, env)
	if obj.Type() == object.Error {
		return obj
	}

	env.Set(node.Ident.Value, obj)

	return obj
}

func evalReturnStatement(node *ast.ReturnStatement, env *object.Environment) object.Object {
	obj := Eval(node.Value, env)
	if obj.Type() == object.Error {
		return obj
	}

	return &object.ReturnObject{Value: obj}
}

func evalIf(node *ast.If, env *object.Environment) object.Object {
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

func evalPrefix(node *ast.Prefix, env *object.Environment) object.Object {
	rightObj := Eval(node.RightValue, env)
	if rightObj.Type() == object.Error {
		return rightObj
	}

	switch node.Operator {
	case "!":
		return evalBang(rightObj)
	case "-":
		return evalMinusPrefix(rightObj)
	default:
		return newError("unknown operation: %s%s", node.Operator, rightObj.Type())
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

func evalInfix(node *ast.Infix, env *object.Environment) object.Object {
	leftObj := Eval(node.LeftValue, env)
	if leftObj.Type() == object.Error {
		return leftObj
	}

	rightObj := Eval(node.RightValue, env)
	if rightObj.Type() == object.Error {
		return rightObj
	}

	switch {
	case leftObj.Type() == object.Integer && rightObj.Type() == object.Integer:
		return evalInfixOfInteger(leftObj, node.Operator, rightObj)
	case leftObj.Type() == object.String && rightObj.Type() == object.String:
		return evalInfixOfString(leftObj, node.Operator, rightObj)
	case node.Operator == "==":
		return convertToBooleanObject(leftObj == rightObj)
	case node.Operator == "!=":
		return convertToBooleanObject(leftObj != rightObj)
	default:
		return newError("unknown operation: %s %s %s", leftObj.Type(), node.Operator, rightObj.Type())
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

func evalInfixOfString(leftObj object.Object, operator string, rightObj object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operation: %s %s %s", leftObj.Type(), operator, rightObj.Type())
	}

	leftVal := leftObj.(*object.StringObject).Value
	rightVal := rightObj.(*object.StringObject).Value

	return &object.StringObject{Value: leftVal + rightVal}
}

func evalFunction(node *ast.Function, env *object.Environment) object.Object {
	return &object.FunctionObject{
		Parameters: node.Parameters,
		Body:       node.Body,
		Env:        env,
	}
}

func isQuoteCall(node ast.Node) bool {
	funcCall, ok := node.(*ast.FunctionCall)
	if !ok {
		return false
	}
	if len(funcCall.Arguments) != 1 {
		return false
	}

	return funcCall.Function.TokenLiteral() == "quote"
}

func evalQuote(node ast.Node, env *object.Environment) object.Object {
	return &object.QuoteObject{
		Value: evalUnquotes(node, env),
	}
}

func evalUnquotes(node ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(node, func(node ast.Node) ast.Node {
		if !isUnquote(node) {
			return node
		}

		unquote := node.(*ast.FunctionCall)
		obj := Eval(unquote.Arguments[0], env)

		return convertObjectToASTNode(obj)
	})
}

func isUnquote(node ast.Node) bool {
	funcCall, ok := node.(*ast.FunctionCall)
	if !ok {
		return false
	}
	if len(funcCall.Arguments) != 1 {
		return false
	}

	return funcCall.Function.TokenLiteral() == "unquote"
}

func convertObjectToASTNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.IntegerObject:
		return convertIntegerObjectToASTNode(obj)
	default:
		return nil
	}
}

func convertIntegerObjectToASTNode(obj *object.IntegerObject) ast.Node {
	return &ast.Integer{
		Token: token.Token{
			Type:    token.Int,
			Literal: fmt.Sprintf("%d", obj.Value),
		},
		Value: obj.Value,
	}
}

func evalFunctionCall(node *ast.FunctionCall, env *object.Environment) object.Object {
	functionObj := Eval(node.Function, env)
	if functionObj.Type() == object.Error {
		return functionObj
	}

	argObjs := evalExpressions(node.Arguments, env)
	if len(argObjs) == 1 && argObjs[0].Type() == object.Error {
		return argObjs[0]
	}

	return applyFunction(functionObj, argObjs)
}

func applyFunction(functionObj object.Object, argObjs []object.Object) object.Object {
	switch function := functionObj.(type) {
	case *object.FunctionObject:
		return applyUserDefinedFunction(function, argObjs)
	case *object.BuiltinFunctionObject:
		return function.Function(argObjs...)
	default:
		return newError("unknown object: %T", functionObj)
	}
}

func applyUserDefinedFunction(functionObj *object.FunctionObject, argObjs []object.Object) object.Object {
	extendedEnv := extendFunctionEnvironment(functionObj, argObjs)
	obj := Eval(functionObj.Body, extendedEnv)

	if obj.Type() == object.Return {
		return obj.(*object.ReturnObject).Value
	}

	return obj
}

func extendFunctionEnvironment(functionObj *object.FunctionObject, argObjs []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(functionObj.Env)
	for i, param := range functionObj.Parameters {
		env.Set(param.Value, argObjs[i])
	}

	return env
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if obj, ok := env.Get(node.Value); ok {
		return obj
	}

	if builtinFn, ok := builtinFns[node.Value]; ok {
		return builtinFn
	}

	return newError("unknown identifier: %s", node.Value)
}

func evalInteger(node *ast.Integer) object.Object {
	return &object.IntegerObject{Value: node.Value}
}

func evalBoolean(node *ast.Boolean) object.Object {
	return convertToBooleanObject(node.Value)
}

func convertToBooleanObject(b bool) object.Object {
	if b {
		return trueObj
	}

	return falseObj
}

func evalString(node *ast.String) object.Object {
	return &object.StringObject{Value: node.Value}
}

func evalArray(node *ast.Array, env *object.Environment) object.Object {
	objs := evalExpressions(node.Elements, env)
	if len(objs) == 1 && objs[0].Type() == object.Error {
		return objs[0]
	}

	return &object.ArrayObject{Elements: objs}
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	objs := make([]object.Object, len(exps))
	for i, exp := range exps {
		obj := Eval(exp, env)
		if obj.Type() == object.Error {
			return []object.Object{obj}
		}

		objs[i] = obj
	}

	return objs
}

func evalHash(node *ast.Hash, env *object.Environment) object.Object {
	values := make(map[object.HashKey]object.HashValue)
	for key, value := range node.Values {
		keyObj := Eval(key, env)
		if keyObj.Type() == object.Error {
			return keyObj
		}
		hashKey, ok := keyObj.(object.HashKeyable)
		if !ok {
			return newError("unusable as hash key: %s", keyObj.Type())
		}

		valueObj := Eval(value, env)
		if valueObj.Type() == object.Error {
			return valueObj
		}

		values[hashKey.HashKey()] = object.HashValue{
			Key:   keyObj,
			Value: valueObj,
		}
	}

	return &object.HashObject{Values: values}
}

func evalSubscript(node *ast.Subscript, env *object.Environment) object.Object {
	leftObj := Eval(node.LeftValue, env)
	if leftObj.Type() == object.Error {
		return leftObj
	}
	index := Eval(node.Index, env)
	if index.Type() == object.Error {
		return index
	}

	switch {
	case leftObj.Type() == object.Array && index.Type() == object.Integer:
		return evalSubscriptToArray(leftObj.(*object.ArrayObject), index.(*object.IntegerObject))
	case leftObj.Type() == object.Hash:
		return evalSubscriptToHash(leftObj.(*object.HashObject), index)
	default:
		return newError("unknown operation: %s[%s]", leftObj.Type(), index.Type())
	}
}

func evalSubscriptToArray(arrayObj *object.ArrayObject, index *object.IntegerObject) object.Object {
	maxIndex := int64(len(arrayObj.Elements) - 1)
	if index.Value < 0 || maxIndex < index.Value {
		return nullObj
	}

	return arrayObj.Elements[index.Value]
}

func evalSubscriptToHash(hashObj *object.HashObject, index object.Object) object.Object {
	hashKey, ok := index.(object.HashKeyable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	hashValue, ok := hashObj.Values[hashKey.HashKey()]
	if !ok {
		return nullObj
	}

	return hashValue.Value
}

func newError(format string, a ...interface{}) object.Object {
	return &object.ErrorObject{Message: fmt.Sprintf(format, a...)}
}
