package evaluator

import (
	"github.com/tomocy/monkey/ast"
	"github.com/tomocy/monkey/object"
)

func DefineMacro(program *ast.Program, env *object.Environment) {
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
