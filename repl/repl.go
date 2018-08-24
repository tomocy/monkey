package repl

import (
	"bufio"
	"fmt"

	"io"

	"github.com/tomocy/monkey/lexer"
	"github.com/tomocy/monkey/token"
)

const prompt = ">> "

func Start(in io.Reader, w io.Writer) {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		sourceCode := scanner.Text()
		fmt.Print(tokens(sourceCode))
		fmt.Print(prompt)
	}
}

func tokens(input string) string {
	b := make([]byte, 0, 10)
	lexer := lexer.New(input)
	for t := lexer.NextToken(); t.Type != token.EOF; t = lexer.NextToken() {
		b = append(b, fmt.Sprintf("%+v\n", t)...)
	}

	return string(b)
}
