package object

import "Monkey/ast"
import "bytes"
import "strings"

//Function function
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

//Type type
func (f *Function) Type() ObjectType {
	return FunctionObj
}

//Inspect inspect
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}

	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString("} {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
