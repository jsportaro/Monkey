package evaluator

import "monkey/object"
import "fmt"

var builtins = map[string]*object.BuiltIn{
	"len": &object.BuiltIn{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": &object.BuiltIn{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. wanted 1 got %d", len(args))
			}

			if args[0].Type() != object.ArrayObj {
				return newError("arguments to `first` must be ARRAY")
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},
	"last": &object.BuiltIn{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. wanted 1 got %d", len(args))
			}

			if args[0].Type() != object.ArrayObj {
				return newError("arguments to `last` must be ARRAY")
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}

			return NULL
		},
	},
	"rest": &object.BuiltIn{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. wanted 1 got %d", len(args))
			}

			if args[0].Type() != object.ArrayObj {
				return newError("arguments to `rest` must be ARRAY")
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])

				return &object.Array{Elements: newElements}
			}

			return NULL
		},
	},
	"push": &object.BuiltIn{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. wanted 2 got %d", len(args))
			}

			if args[0].Type() != object.ArrayObj {
				return newError("first argument to `push` must be ARRAY")
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]

			return &object.Array{Elements: newElements}
		},
	},
	"puts": &object.BuiltIn{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL
		},
	},
}
