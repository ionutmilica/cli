package cli

import (
	"errors"
	"fmt"
	"strings"
)

type matcher struct {
	ctx *Context
	mgr *FlagMgr

	//
	args   []string
	cursor int
}

func newMatcher(args []string, flags []*Flag) *matcher {
	matcher := &matcher{
		ctx:    newContext(),
		mgr:    newFlagMgr(flags),
		args:   args,
		cursor: 0,
	}

	return matcher
}

// Have one more arg?
func (m *matcher) hasNext() bool {
	if m.cursor >= len(m.args) {
		return false
	}
	return true
}

// Go to the next arg
func (m *matcher) next(steps ...int) {
	if len(steps) > 0 {
		m.cursor += steps[0]
	} else {
		m.cursor += 1
	}
}

func (m *matcher) peek() string {
	return m.args[m.cursor+1]
}

func (m *matcher) match() error {
	m.cursor = 0

	for m.hasNext() {
		arg := m.args[m.cursor]
		switch {
		case strings.HasPrefix(arg, "--"): // We matched and long option
			if err := m.matchLongOption(arg); err != nil {
				return err
			}
			break
		default: // We matched an argument
			if err := m.matchArgument(arg); err != nil {
				return err
			}
		}

		m.next()
	}

	return m.validate()
}

// Validate arguments so the matcher will return error if requiredArgs != foundArgs
func (m *matcher) validate() error {
	requiredArgs := m.mgr.requiredArgs()

	if len(requiredArgs) == len(m.ctx.Arguments) {
		return nil
	}

	var missing []string
	for _, arg := range requiredArgs {
		if _, ok := m.ctx.Arguments[arg]; !ok {
			missing = append(missing, arg)
		}
	}
	return errors.New(fmt.Sprintf("Not enough arguments (missing: %s).", strings.Join(missing, ", ")))
}

// Parses options like --opt, --opt=val --opt val according to the defined flags
func (m *matcher) matchLongOption(arg string) error {
	var value string
	arg = arg[2:]

	if strings.Contains(arg, "=") {
		parts := strings.Split(arg, "=")
		arg = parts[0]
		value = parts[1]
	}

	if !m.mgr.hasOption(arg) {
		return errors.New(fmt.Sprintf("The `--%s` option does not exist.", arg))
	}

	option := m.mgr.option(arg)

	if value != "" && !option.acceptValue() {
		return errors.New(fmt.Sprintf("The `--%s` option does not accept a value!", arg))
	}

	if value == "" && option.acceptValue() && m.hasNext() {
		peek := m.peek()
		if len(peek) > 0 && peek[0] != '-' {
			value = peek
			m.next()
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
	m.ctx.Options[arg] = []string{value}

	return nil
}

// Parse strings that are not starting with - as arguments and group them according to the signature
func (m *matcher) matchArgument(arg string) error {
	current := len(m.ctx.Arguments)

	if m.mgr.hasArgument(current) {
		m.ctx.Arguments[m.mgr.argument(current).name] = []string{arg}
	} else if m.mgr.hasArgument(current-1) && m.mgr.argument(current-1).isArrayArgument() {
		m.ctx.Arguments[m.mgr.argument(current-1).name] = append(m.ctx.Arguments[m.mgr.argument(current-1).name], arg)
	} else {
		return errors.New("To many arguments!")
	}

	return nil
}
