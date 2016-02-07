package cli

type Context struct {
	// @todo: Replace string with some structure that will allow type conversion, i.e Arguments["arg"].ToInt64()
	Arguments map[string][]string
	Options   map[string][]string
}

func newContext() *Context {
	return &Context{
		Arguments: make(map[string][]string, 0),
		Options:   make(map[string][]string, 0),
	}
}

// Set argument with values
func (ctx *Context) SetArgument(key string, values ...string) {
	ctx.Arguments[key] = values
}

// Append values to the argument
func (ctx *Context) AppendToArgument(key string, value string) {
	ctx.Arguments[key] = append(ctx.Arguments[key], value)
}

// Set option with values
func (ctx *Context) SetOption(key string, values ...string) {
	ctx.Options[key] = values
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
