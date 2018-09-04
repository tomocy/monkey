package repl

import (
	"bufio"
	"fmt"
	"strings"

	"io"

	"github.com/tomocy/monkey/lexer"
	"github.com/tomocy/monkey/parser"
)

const prompt = ">> "

func Start(in io.Reader, w io.Writer) {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		sourceCode := scanner.Text()
		fmt.Println(parsedProgramOrErrorMessages(sourceCode))
		fmt.Print(prompt)
	}
}

func parsedProgramOrErrorMessages(in string) string {
	parser := parser.New(lexer.New(in))
	program := parser.ParseProgram()
	if len(parser.Errors()) != 0 {
		return strings.Join(parser.Errors(), "\n")
	}

	return program.String()
}
