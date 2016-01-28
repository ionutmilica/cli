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
	re := regexp.MustCompile("{([^{}]+)}")
	matches := re.FindAllStringSubmatch(p.raw, -1)

	for _, m := range matches {
		flag := m[1]

		if len(flag) == 0 {
			panic("Flag cannot be empty! Syntax like {} is not acceptable!")
		}

		var kind uint
		var name string

		if strings.HasPrefix(flag, "--") {
			name = flag[2:]
			kind = longOptionFlag
		} else if flag[0] == '-' {
			name = flag[1:]
			kind = optionFlag
		} else {
			name = flag
			kind = valueFlag
		}

		p.flags = append(p.flags, &Flag{
			kind: kind,
			name: name,
		})

		// Value flag
	}
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
