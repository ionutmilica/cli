package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

type Handler func(*Context)

// Context store the arguments and options and have attached helpers methods
// to deal with console operations
type Context struct {
	Arguments map[string]*Result
	Options   map[string]*Result
	Reader    io.Reader
	Writer    io.Writer

	handlers []Handler
	cursor   int
}

// Creates a new context
// It needs: reader, writer, arguments map and option map
func newContext(reader io.Reader, writer io.Writer, args map[string]*Result, opts map[string]*Result) *Context {
	return &Context{
		Arguments: args,
		Options:   opts,
		Reader:    reader,
		Writer:    writer,
	}
}

// Add new handlers to the end of the chain
func (ctx *Context) AppendHandler(handlers ...Handler) {
	ctx.handlers = append(ctx.handlers, handlers...)
}

// Triggers the next handler in the chain
func (ctx *Context) Next() {
	ctx.cursor++
	ctx.Run()
}

// Run the handler as position `cursor`
// This will execute handlers one by one until a handler with no ctx.Next() is found
func (ctx *Context) Run() {
	if len(ctx.handlers) > ctx.cursor {
		ctx.handlers[ctx.cursor](ctx)
	}
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

// Display a message then waits for an answer
func (ctx *Context) Ask(msg string) string {
	// Display the message
	fmt.Fprint(ctx.Writer, msg)

	// Create a new reader from ctx.Reader and wait for the response
	reader := bufio.NewReader(ctx.Reader)
	text, err := reader.ReadString('\n')

	if err != nil {
		return ""
	}

	return text
}
