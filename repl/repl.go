package repl

import (
	"bufio"
	"fmt"
	"strings"

	"io"

	"github.com/tomocy/monkey/evaluator"
	"github.com/tomocy/monkey/lexer"
	"github.com/tomocy/monkey/object"
	"github.com/tomocy/monkey/parser"
)

const prompt = ">> "

var env = object.NewEnvironment()
var macroEnv = object.NewEnvironment()

func Start(in io.Reader, w io.Writer) {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		sourceCode := scanner.Text()
		fmt.Println(evaluatedProgramOrErrorMessages(sourceCode))
		fmt.Print(prompt)
	}
}

func evaluatedProgramOrErrorMessages(in string) string {
	parser := parser.New(lexer.New(in))
	program := parser.ParseProgram()
	if len(parser.Errors()) != 0 {
		return strings.Join(parser.Errors(), "\n")
	}

	evaluator.DefineMacros(program, macroEnv)
	expandedProgram := evaluator.ExpandMacros(program, macroEnv)

	evaluatedProgram := evaluator.Eval(expandedProgram, env)
	if evaluatedProgram == nil {
		return ""
	}

	return evaluatedProgram.Inspect()
}
