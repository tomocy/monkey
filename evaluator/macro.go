package evaluator

import (
	"fmt"

	"github.com/tomocy/monkey/ast"
	"github.com/tomocy/monkey/object"
	"github.com/tomocy/monkey/token"
)

func evalQuote(node ast.Node, env *object.Environment) object.Object {
	return &object.QuoteObject{
		Value: evalUnquotes(node, env),
	}
}

func isQuote(node ast.Node) bool {
	funcCall, ok := node.(*ast.FunctionCall)
	if !ok {
		return false
	}
	if len(funcCall.Arguments) != 1 {
		return false
	}

	return funcCall.Function.TokenLiteral() == "quote"
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
	case *object.BooleanObject:
		return convertBooleanObjectToASTNode(obj)
	case *object.StringObject:
		return convertStringObjectToASTNode(obj)
	case *object.ArrayObject:
		return convertArrayObjectToASTNode(obj)
	case *object.QuoteObject:
		return obj.Value
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

func convertBooleanObjectToASTNode(obj *object.BooleanObject) ast.Node {
	var t token.Token
	if obj.Value {
		t = token.Token{
			Type:    token.True,
			Literal: "true",
		}
	} else {
		t = token.Token{
			Type:    token.False,
			Literal: "false",
		}
	}

	return &ast.Boolean{
		Token: t,
		Value: obj.Value,
	}
}

func convertStringObjectToASTNode(obj *object.StringObject) ast.Node {
	return &ast.String{
		Token: token.Token{
			Type:    token.String,
			Literal: obj.Value,
		},
		Value: obj.Value,
	}
}

func convertArrayObjectToASTNode(obj *object.ArrayObject) ast.Node {
	exps := make([]ast.Expression, len(obj.Elements))
	for i, obj := range obj.Elements {
		exps[i] = convertObjectToASTNode(obj).(ast.Expression)
	}
	return &ast.Array{
		Token: token.Token{
			Type:    token.LBracket,
			Literal: obj.Inspect(),
		},
		Elements: exps,
	}
}

func DefineMacros(program *ast.Program, env *object.Environment) {
	setMacroDefinitionsInEnv(program, env)
	removeMacroDefinitionsFromProgram(program)
}

func setMacroDefinitionsInEnv(program *ast.Program, env *object.Environment) {
	for _, stmt := range program.Statements {
		if !isMacroDefinition(stmt) {
			continue
		}

		letStmt := stmt.(*ast.LetStatement)
		macro := letStmt.Value.(*ast.Macro)
		macroObj := &object.MacroObject{
			Parameters: macro.Parameters,
			Body:       macro.Body,
			Env:        env,
		}

		env.Set(letStmt.Ident.Value, macroObj)
	}
}

func removeMacroDefinitionsFromProgram(program *ast.Program) {
	macroDefIndexes := findMacroDefinitions(*program)
	for i := len(macroDefIndexes) - 1; 0 <= i; i-- {
		index := macroDefIndexes[i]
		program.Statements = append(program.Statements[:index], program.Statements[index+1:]...)
	}
}

func findMacroDefinitions(program ast.Program) []int {
	indexes := make([]int, 0)
	for i, stmt := range program.Statements {
		if !isMacroDefinition(stmt) {
			continue
		}
		indexes = append(indexes, i)
	}

	return indexes
}

func isMacroDefinition(stmt ast.Statement) bool {
	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		return false
	}

	_, ok = letStmt.Value.(*ast.Macro)
	return ok
}

func ExpandMacros(program *ast.Program, env *object.Environment) ast.Node {
	return ast.Modify(program, func(node ast.Node) ast.Node {
		funcCall, ok := node.(*ast.FunctionCall)
		if !ok {
			return node
		}

		macro, ok := getMacroObject(funcCall, *env)
		if !ok {
			return node
		}

		quotedArgObjs := quoteExpressions(funcCall.Arguments)

		return applyMacro(*macro, quotedArgObjs)
	})
}

func getMacroObject(node *ast.FunctionCall, env object.Environment) (*object.MacroObject, bool) {
	ident, ok := node.Function.(*ast.Identifier)
	if !ok {
		return nil, false
	}

	obj, ok := env.Get(ident.Value)
	if !ok {
		return nil, false
	}

	macro, ok := obj.(*object.MacroObject)

	return macro, ok
}

func quoteExpressions(nodes []ast.Expression) []*object.QuoteObject {
	exps := make([]*object.QuoteObject, len(nodes))
	for i, exp := range nodes {
		exps[i] = &object.QuoteObject{Value: exp}
	}

	return exps
}

func applyMacro(macro object.MacroObject, argObjs []*object.QuoteObject) ast.Node {
	extendedEnv := extendMacroEnvironment(macro, argObjs)
	obj := Eval(macro.Body, extendedEnv)
	quote, ok := obj.(*object.QuoteObject)
	if !ok {
		panic("invalid macro definition: macro should return quoted value")
	}

	return quote.Value
}

func extendMacroEnvironment(macro object.MacroObject, argObjs []*object.QuoteObject) *object.Environment {
	extendedEnv := object.NewEnclosedEnvironment(macro.Env)
	for i, param := range macro.Parameters {
		extendedEnv.Set(param.Value, argObjs[i])
	}

	return extendedEnv
}
