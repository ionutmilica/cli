package cli

type Context struct {
}

func (ctx *Context) Confirm(message string) bool {
	return true
}

func (ctx *Context) Info(message string) {
	// print something
}
