package cli

import (
	"errors"
	"fmt"
	"strings"
)

type Context struct {
	Arguments map[string][]string
	Options   map[string][]string

	args   []string
	cursor int
}

func (ctx *Context) parse(args []string, mgr *FlagMgr) {
	ctx.args = args
	ctx.cursor = 0

	for ctx.cursor < len(ctx.args) {
		arg := args[ctx.cursor]
		switch {
		case strings.HasPrefix(arg, "--"): // We matched and long option
			ctx.parseLongOption(mgr, arg)
			break
		default: // We matched an argument
			ctx.parseArgument(mgr, arg)
		}

		ctx.cursor++
	}
}

func (ctx *Context) parseLongOption(mgr *FlagMgr, arg string) {
	var value string
	arg = arg[2:]

	if strings.Contains(arg, "=") {
		parts := strings.Split(arg, "=")
		arg = parts[0]
		value = parts[1]
	}

	if !mgr.hasOption(arg) {
		panic(fmt.Sprintf("The `--%s` option does not exist.", arg))
	}

	option := mgr.option(arg)

	if value != "" && !option.acceptValue() {
		panic(fmt.Sprintf("The `--%s` option does not accept a value!", arg))
	}

	if value == "" && option.acceptValue() && hasIndex(len(ctx.args), ctx.cursor+1) {
		next := ctx.args[ctx.cursor+1]
		if len(next) > 0 && next[0] != '-' {
			value = next
			ctx.cursor += 2
		}
	}

	if value == "" {
		if option.isValueRequired() {
			panic(fmt.Sprintf("The `--%s` option requres a value!", arg))
		}

		if !option.isValueArray() && option.isValueOptional() {
			value = option.value
		}
	}
	ctx.Options[arg] = []string{value}
}

func (ctx *Context) parseArgument(mgr *FlagMgr, arg string) error {
	current := len(ctx.Arguments)

	if mgr.hasArgument(current) {
		ctx.Arguments[mgr.argument(current).name] = []string{arg}
	} else if mgr.hasArgument(current-1) && mgr.argument(current-1).isArrayArgument() {
		ctx.Arguments[mgr.argument(current-1).name] = append(ctx.Arguments[mgr.argument(current-1).name], arg)
	} else {
		return errors.New("To many arguments!")
	}

	return nil
}

func (ctx *Context) Confirm(message string) bool {
	return true
}

func (ctx *Context) Info(message string) {
	// print something
}
