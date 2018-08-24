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
		lexer := lexer.New(sourceCode)
		for t := lexer.NextToken(); t.Type != token.EOF; t = lexer.NextToken() {
			fmt.Printf("%+v\n", t)
		}
		fmt.Print(prompt)
	}
}
