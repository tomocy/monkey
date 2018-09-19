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
