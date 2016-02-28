package cli

import (
	"errors"
	"fmt"
	"strings"
)

type matcher struct {
	arguments map[string]*Result
	options   map[string]*Result
	flags     FlagList

	//
	args   []string
	cursor int
}

func newMatcher(args []string, flags FlagList) *matcher {
	matcher := &matcher{
		arguments: make(map[string]*Result, 0),
		options:   make(map[string]*Result, 0),
		flags:     flags,
		args:      args,
		cursor:    0,
	}

	return matcher
}

// Set argument with values
func (m *matcher) setArgument(key string, values ...string) {
	if m.arguments[key] == nil {
		m.arguments[key] = &Result{}
	}

	m.arguments[key].Append(values...)
}

// Set option with values
func (m *matcher) setOption(key string, values ...string) {
	if m.options[key] == nil {
		m.options[key] = &Result{}
	}
	m.options[key].Append(values...)
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
	if !m.hasNext() {
		return
	}
	if len(steps) > 0 {
		m.cursor += steps[0]
	} else {
		m.cursor += 1
	}
}

// Get the current item
func (m *matcher) current() string {
	return m.args[m.cursor]
}

// Look ahead and get the next element without moving the curosr
func (m *matcher) peek() (string, error) {
	if m.cursor+1 >= len(m.args) {
		return "", errors.New("Cannot peek if we don't have any more items!")
	}
	return m.args[m.cursor+1], nil
}

// start the matching process
func (m *matcher) match() error {
	m.cursor = 0

	for m.hasNext() {
		arg := m.current()
		switch {
		case strings.HasPrefix(arg, "-"): // We matched an option
			if err := m.matchOption(arg); err != nil {
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

	// Set arguments with default value
	for _, flag := range m.flags {
		if !flag.isArgument() {
			continue
		}
		if _, ok := m.arguments[flag.name]; !ok && flag.value != "" {
			m.setArgument(flag.name, flag.value)
		}
	}

	// Set flags with default value
	for _, flag := range m.flags {
		if flag.isArgument() {
			continue
		}
		if _, ok := m.options[flag.name]; !ok && flag.value != "" {
			m.setOption(flag.name, flag.value)
		}
	}
	requiredArgs := m.flags.requiredArgs()

	if len(requiredArgs) <= len(m.arguments) {
		return nil
	}

	var missing []string
	for _, arg := range requiredArgs {
		if _, ok := m.arguments[arg]; !ok {
			missing = append(missing, arg)
		}
	}
	return m.fail("Not enough arguments (missing: %s).", strings.Join(missing, ", "))
}

// Parses options like --opt, --opt=val --opt val according to the defined flags
func (m *matcher) matchOption(arg string) error {
	value := ""
	isShort := false

	if strings.HasPrefix(arg, "--") {
		arg = arg[2:]
	} else {
		arg = arg[1:]
		isShort = true
	}

	if strings.Contains(arg, "=") {
		parts := strings.Split(arg, "=")
		arg = parts[0]
		value = parts[1]
	}

	if isShort {
		if len(arg) > 1 {
			// -abc=something is strange
			if value != "" {
				return m.fail("The `-%s` options cannot accept values!", arg)
			}

			// -abc should contains 3 options
			for _, c := range strings.Split(arg, "") {
				err := m.matchOption("-" + c)
				if err != nil {
					return err
				}
			}
			return nil
		}
	}

	option := m.flags.option(arg)

	if option == nil {
		return m.fail("The `--%s` option does not exist.", arg)
	}

	if value != "" && !option.acceptValue() {
		return m.fail("The `--%s` option does not accept a value!", arg)
	}

	if value == "" && option.acceptValue() && m.hasNext() {
		peek, err := m.peek()
		if err == nil && len(peek) > 0 && peek[0] != '-' {
			value = peek
			m.next()
		}
	}

	if value == "" {
		if option.isRequired() {
			return m.fail("The `--%s` option requires a value!", arg)
		}

		if !option.isArray() && option.isOptional() {
			value = option.value
		}
	}

	if value != "" {
		if _, ok := m.options[arg]; !ok {
			m.setOption(arg, value)
			return nil
		}
		if !option.isArray() {
			return m.fail("The `--%s` option does not accept an array of values!", arg)
		}
		// Append to option
		m.setOption(arg, value)
	} else {
		m.setOption(arg)
	}

	return nil
}

// Parse strings that are not starting with - as arguments and group them according to the signature
func (m *matcher) matchArgument(argName string) error {
	current := len(m.arguments)

	if arg := m.flags.argument(current); arg != nil {
		m.setArgument(arg.name, argName)
	} else if arg := m.flags.argument(current - 1); arg != nil && arg.isArray() {
		m.setArgument(arg.name, argName)
	} else {
		return m.fail("To many arguments!")
	}

	return nil
}

func (m *matcher) reset() {
	m.cursor = 0
	m.arguments = map[string]*Result{}
	m.options = map[string]*Result{}
}

// Clean the context and return the error
func (m *matcher) fail(msg string, args ...interface{}) error {
	m.reset()
	return fmt.Errorf(msg, args...)
}
