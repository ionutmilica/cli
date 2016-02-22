package cli

const (
	argumentFlag = iota
	optionFlag
	longOptionFlag
)

const (
	// Flags for arguments
	optional = 1
	isArray  = 2
	required = 4

	// Flags for options
	valueNone     = 8
	valueRequired = 16
	valueOptional = 32
	valueArray    = 64
)

/** Option flags **/

type Flag struct {
	kind        int8
	name        string
	options     int8
	description string
	value       string
}

// Check if the flag is an argument
func (f Flag) isArgument() bool {
	return f.kind == argumentFlag
}

// Check if the flag is an option
func (f Flag) isLongOption() bool {
	return f.kind == longOptionFlag
}

// Check if a given flag is an option
func (f Flag) isOption() bool {
	return f.kind == optionFlag
}

// Check if option accepts a value
func (f Flag) acceptValue() bool {
	if f.options&valueNone == valueNone {
		return false
	}
	return f.isOptional() || f.isRequired()
}

// Check if argument or option accepts more than one value
func (f Flag) isArray() bool {
	if f.isArgument() {
		return f.options&isArray == isArray
	}
	return f.options&valueArray == valueArray
}

// Check if argument is optional or value for option is optional (as the option itself is always optional)
func (f Flag) isOptional() bool {
	if f.isArgument() {
		return f.options&optional == optional
	}
	return f.options&valueOptional == valueOptional
}

// Check if argument is required or value for option is required
func (f Flag) isRequired() bool {
	if f.isArgument() {
		return f.options&required == required
	}
	return f.options&valueRequired == valueRequired
}

// Get the name of the argument/option
func (f Flag) String() string {
	return f.name
}

/** Flag list **/

type FlagList []*Flag

// Append flag to the flag list
func (fl *FlagList) Append(flags ...*Flag) {
	*fl = append(*fl, flags...)
}

// Get arguments that are required
func (fl *FlagList) requiredArgs() []string {
	flags := []string{}

	for _, arg := range *fl {
		if arg.isArgument() && arg.isRequired() {
			flags = append(flags, arg.name)
		}
	}

	return flags
}

// Find the argument number `pos` from the list of the flags. Options will be skipped
func (fl *FlagList) argument(pos int) *Flag {
	current := 0
	for _, arg := range *fl {
		if arg.isArgument() {
			if current == pos {
				return arg
			}
			current++
		} else {
			continue
		}
	}
	return nil
}

// Find option by name. Arguments will be skipped
func (fl *FlagList) option(opt string) *Flag {
	for _, flag := range *fl {
		if !flag.isArgument() && flag.name == opt {
			return flag
		}
	}
	return nil
}
