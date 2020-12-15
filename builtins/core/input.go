package core

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Input reads a line from standard input optionally printing prompt.
// input([prompt]) prints the prompt.
func Input(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("input", args,
		typing.RangeOfArgs(0, 1),
		typing.WithTypes(object.STRING),
	); err != nil {
		return object.NewError(err.Error())
	}

	if len(args) == 1 {
		prompt := args[0].(*object.String).Value
		fmt.Fprintf(os.Stdout, prompt)
	}

	buffer := bufio.NewReader(os.Stdin)

	line, _, err := buffer.ReadLine()
	if err != nil && err != io.EOF {
		return object.NewError(fmt.Sprintf("error reading input from stdin: %s", err))
	}

	return &object.String{Value: string(line)}
}
