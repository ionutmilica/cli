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

			if f.isArray() {
				hadArrayArg = true
			}
		}
	}
}

func (cmd *Command) parseOption(opt string) *Flag {
	var description string
	var kind int8
	var options int8

	if strings.HasPrefix(opt, "--") {
		kind = longOptionFlag
		opt = opt[2:]
	} else {
		kind = optionFlag
		opt = opt[1:]
	}

	if strings.Contains(opt, " : ") {
		parts := strings.Split(opt, " : ")
		opt = parts[0]
		description = parts[1]
	}

	switch {
	case strings.HasSuffix(opt, "="):
		options = valueOptional
		opt = strings.TrimSuffix(opt, "=")
		break
	case strings.HasSuffix(opt, "=*"):
		options = valueOptional | valueArray
		opt = strings.TrimSuffix(opt, "=*")
		break
	case strings.HasSuffix(opt, "=+"):
		options = valueRequired | valueArray
		opt = strings.TrimSuffix(opt, "=+")
		break
	default:
		options = valueNone
	}

	flag := &Flag{
		kind:        kind,
		name:        opt,
		description: description,
		options:     options,
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
		options = required
	}

	flag := &Flag{
		name:        arg,
		kind:        argumentFlag,
		options:     options,
		description: description,
	}
	cmd.Flags = append(cmd.Flags, flag)

	return flag
}
