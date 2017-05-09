package script

import (
	"Monkey/executor"
	"fmt"
	"io"
	"io/ioutil"
	"monkey/object"
)

//Run runs all the files
func Run(output io.Writer, files []string) {
	env := object.NewEnvironment()
	sources := []string{}
	for _, file := range files {
		srcBytes, _ := ioutil.ReadFile(file)
		fmt.Printf("Evaluating file %s\n", file)

		src := string(srcBytes[:])
		sources = append(sources, src)
	}

	executor.Execute(sources, env, output)
}
