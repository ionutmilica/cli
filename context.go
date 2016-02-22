package cli

import (
	"bufio"
	"errors"
	"io"
)

// Context store the arguments and options and have attached helpers methods
// to deal with console operations
type Context struct {
	Arguments map[string]*Result
	Options   map[string]*Result
	Reader    io.Reader
	Writer    io.Writer
}

func newContext(reader io.Reader, writer io.Writer) *Context {
	return &Context{
		Arguments: make(map[string]*Result, 0),
		Options:   make(map[string]*Result, 0),
		Reader:    reader,
		Writer:    writer,
	}
}

func (ctx *Context) reset() {
	ctx.Arguments = map[string]*Result{}
	ctx.Options = map[string]*Result{}
}

// Set argument with values
func (ctx *Context) SetArgument(key string, values ...string) {
	if ctx.Arguments[key] == nil {
		ctx.Arguments[key] = &Result{}
	}

	ctx.Arguments[key].Append(values...)
}

// Set option with values
func (ctx *Context) SetOption(key string, values ...string) {
	if ctx.Options[key] == nil {
		ctx.Options[key] = &Result{}
	}
	ctx.Options[key].Append(values...)
}

// Check if context has a specific option
func (ctx *Context) HasOption(key string) bool {
	if _, ok := ctx.Options[key]; ok {
		return true
	}
	return false
}

// Check if context has a specific argument
func (ctx *Context) HasArgument(key string) bool {
	if _, ok := ctx.Arguments[key]; ok {
		return true
	}
	return false
}

// Get option from the context
func (ctx *Context) Option(key string) (*Result, error) {
	if _, ok := ctx.Options[key]; ok {
		return ctx.Options[key], nil
	}

	return nil, errors.New("Option not present!")
}

// Get argument from the context
func (ctx *Context) Argument(key string) (*Result, error) {
	if _, ok := ctx.Arguments[key]; ok {
		return ctx.Arguments[key], nil
	}

	return nil, errors.New("Argument not present!")
}

func (ctx *Context) Ask(msg string) string {
	reader := bufio.NewReader(ctx.Reader)
	text, err := reader.ReadString('\n')

	if err != nil {
		return ""
	}

	return text
}
