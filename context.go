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

// Parse and match os arguments against the registered flags
func (ctx *Context) parse(args []string, mgr *FlagMgr) error {
	ctx.args = args
	ctx.cursor = 0

	for ctx.cursor < len(ctx.args) {
		arg := args[ctx.cursor]
		switch {
		case strings.HasPrefix(arg, "--"): // We matched and long option
			if err := ctx.parseLongOption(mgr, arg); err != nil {
				return err
			}
			break
		default: // We matched an argument
			if err := ctx.parseArgument(mgr, arg); err != nil {
				return err
			}
		}

		ctx.cursor++
	}

	// Validate arguments
	requiredArgs := mgr.requiredArgs()

	if len(requiredArgs) > len(ctx.Arguments) {
		missing := []string{}
		for _, arg := range requiredArgs {
			if _, ok := ctx.Arguments[arg]; !ok {
				missing = append(missing, arg)
			}
			return errors.New(fmt.Sprintf("Not enough arguments (missing: `%s`).", strings.Join(missing, ", ")))
		}
	}

	return nil
}

// Parses options like --opt, --opt=val --opt val according to the defined flags
func (ctx *Context) parseLongOption(mgr *FlagMgr, arg string) error {
	var value string
	arg = arg[2:]

	if strings.Contains(arg, "=") {
		parts := strings.Split(arg, "=")
		arg = parts[0]
		value = parts[1]
	}

	if !mgr.hasOption(arg) {
		return errors.New(fmt.Sprintf("The `--%s` option does not exist.", arg))
	}

	option := mgr.option(arg)

	if value != "" && !option.acceptValue() {
		return errors.New(fmt.Sprintf("The `--%s` option does not accept a value!", arg))
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
			return errors.New(fmt.Sprintf("The `--%s` option requres a value!", arg))
		}

		if !option.isValueArray() && option.isValueOptional() {
			value = option.value
		}
	}
	ctx.Options[arg] = []string{value}

	return nil
}

// Parse strings that are not starting with - as arguments and group them according to the signature
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
