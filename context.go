package cli

import (
	"fmt"
	"strings"
)

type Context struct {
	Arguments map[string][]string
	Options   map[string]string
}

func (ctx *Context) parse(args []string, flags []*Flag) {
	arguments := make([]*Flag, 0)
	options := make(map[string]*Flag, 0)

	for _, flag := range flags {
		if flag.kind == argumentFlag {
			arguments = append(arguments, flag)
		} else {
			options[flag.name] = flag
		}
	}

	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "--"): // We matched and long option
			var value string
			arg := arg[2:]

			if strings.Contains(arg, "=") {
				parts := strings.Split(arg, "=")
				arg = parts[0]
				value = parts[1]
			}

			if _, ok := options[arg]; !ok {
				panic(fmt.Sprintf("The `--%s` option does not exist.", arg))
			}

			println(value)

			break
		default: // We matched an argument
			current := len(ctx.Arguments)
			if hasIndex(len(arguments), current) {
				ctx.Arguments[arguments[current].name] = []string{arg}
			} else if hasIndex(len(arguments), current-1) && arguments[current-1].isArrayArgument() {
				ctx.Arguments[arguments[current-1].name] = append(ctx.Arguments[arguments[current-1].name], arg)
			} else {
				panic("To many arguments!")
			}
		}
	}

	fmt.Println(ctx.Arguments)
}

func hasIndex(size int, i int) bool {
	if size == 0 {
		return false
	}
	if i > -1 && i < size {
		return true
	}

	return false
}

func (ctx *Context) Confirm(message string) bool {
	return true
}

func (ctx *Context) Info(message string) {
	// print something
}
