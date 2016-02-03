package cli

import (
	_ "fmt"
	"regexp"
	"strings"
)

type pattern struct {
	raw   string
	flags []*Flag
}

func (p *pattern) parse() {
	re := regexp.MustCompile("{([^{}]*)}")
	matches := re.FindAllStringSubmatch(p.raw, -1)

	for _, m := range matches {
		flag := m[1]

		if len(flag) == 0 {
			panic("Flag cannot be empty! Syntax like {} is not acceptable!")
		}

		if flag[0] == '-' {
			p.parseOption(flag)
		} else {
			p.parseArgument(flag)
		}
	}
}

func (p *pattern) parseOption(opt string) {
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

	p.flags = append(p.flags, &Flag{
		kind:        kind,
		name:        opt,
		description: description,
	})
}

func (p *pattern) parseArgument(arg string) {
	var description string
	var options int8

	if strings.Contains(arg, " : ") {
		parts := strings.Split(arg, " : ")
		arg = parts[0]
		description = parts[1]
	}

	if strings.HasSuffix(arg, "?*") {
		options = isArray | optional
		arg = strings.TrimSuffix(arg, "?*")
	}

	if strings.HasSuffix(arg, "*") {
		options = isArray | required
		arg = strings.TrimSuffix(arg, "*")
	}

	if strings.HasSuffix(arg, "?") {
		options = optional
		arg = strings.TrimSuffix(arg, "?")
	}

	p.flags = append(p.flags, &Flag{
		kind:        argumentFlag,
		options:     options,
		name:        arg,
		description: description,
	})
}

func (p *pattern) match(args []string) {

}

// Takes a pattern and parse it in pieces
func newPattern(signature string) *pattern {
	pattern := &pattern{
		raw:   signature,
		flags: make([]*Flag, 0),
	}
	pattern.parse()

	return pattern
}
