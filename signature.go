package cli

import (
	_ "fmt"
	"regexp"
	"strings"
)

func (cmd *Command) parse() {
	re := regexp.MustCompile("{([^{}]*)}")
	matches := re.FindAllStringSubmatch(cmd.Signature, -1)

	//
	hadArrayArg := false

	for _, m := range matches {
		var f *Flag
		flag := m[1]

		if len(flag) == 0 {
			panic("Flag cannot be empty! Syntax like {} is not acceptable!")
		}

		if flag[0] == '-' {
			f = cmd.parseOption(flag)
		} else {
			if hadArrayArg {
				panic("Command cannot have argument after an array argument!")
			}

			f = cmd.parseArgument(flag)

			if f.options&isArray == isArray {
				hadArrayArg = true
			}
		}
	}
}

func (cmd *Command) parseOption(opt string) *Flag {
	var description string
	var kind int8

	if strings.HasPrefix(opt, "--") {
		kind = longOptionFlag
	} else {
		kind = optionFlag
	}

	if strings.Contains(opt, " : ") {
		parts := strings.Split(opt, " : ")
		opt = parts[0]
		description = parts[1]
	}

	flag := &Flag{
		kind:        kind,
		name:        opt,
		description: description,
	}
	cmd.Flags = append(cmd.Flags, flag)

	return flag
}

func (cmd *Command) parseArgument(arg string) *Flag {
	var description string
	var options int8

	if strings.Contains(arg, " : ") {
		parts := strings.Split(arg, " : ")
		arg = parts[0]
		description = parts[1]
	}

	switch {
	case strings.HasSuffix(arg, "?*"):
		options = isArray | optional
		arg = strings.TrimSuffix(arg, "?*")
		break
	case strings.HasSuffix(arg, "*"):
		options = isArray | required
		arg = strings.TrimSuffix(arg, "*")
		break
	case strings.HasSuffix(arg, "?"):
		options = optional
		arg = strings.TrimSuffix(arg, "?")
		break
	default:
		options = 0
	}

	flag := &Flag{
		kind:        argumentFlag,
		options:     options,
		name:        arg,
		description: description,
	}
	cmd.Flags = append(cmd.Flags, flag)

	return flag
}
