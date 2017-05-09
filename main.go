package main

import (
	"fmt"
	"monkey/repl"
	"monkey/script"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	if len(os.Args) > 2 {
		fmt.Printf("Scripting mode.\n")
		script.Run(os.Stdout, os.Args[1:])
	} else {
		fmt.Printf("Hello %s! This is the monkey programming language!\n", user.Username)
		fmt.Printf("Feel free to type in commands\n")

		repl.Start(os.Stdin, os.Stdout)
	}
}
