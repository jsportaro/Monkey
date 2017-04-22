// repl/repl.go

package repl

import (
	"Monkey/lexer"
	"Monkey/parser"
	"Monkey/token"
	"bufio"
	"fmt"
	"io"
)

const prompt = ">>"

//Start It begins here
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(prompt)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()

		l := lexer.New(line)
		p := parser.New(l)
		p.ParseProgram()
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
