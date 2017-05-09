// repl/repl.go

package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/executor"
	"monkey/object"
)

const prompt = ">>"

//Start It begins here
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(prompt)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		executor.Execute([]string{line}, env, out)
	}
}
