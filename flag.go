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
	data        []string
}

func (f Flag) ToStr(i ...int) string {
	index := getIndex(i)
	return f.data[index]
}

// Check if the flag is an argument
func (f Flag) isArgument() bool {
	return f.kind == argumentFlag
}

// Check if the flag is an option
func (f Flag) isLongOption() bool {
	return f.kind == longOptionFlag
}

// Check if option accepts a value
func (f Flag) acceptValue() bool {
	return f.isOptional() || f.isRequired()
}

// Check if argument or option accepts more than one value
func (f Flag) isArray() bool {
	if f.kind == argumentFlag {
		return f.options&isArray == isArray
	}
	return f.options&valueArray == valueArray
}

// Check if argument is optional or value for option is optional (as the option itself is always optional)
func (f Flag) isOptional() bool {
	if f.kind == argumentFlag {
		return f.options&optional == optional
	}
	return f.options&valueOptional == valueOptional
}

// Check if argument is required or value for option is required
func (f Flag) isRequired() bool {
	if f.kind == argumentFlag {
		return f.options&required == required
	}
	return f.options&valueRequired == valueRequired
}

// Get the name of the argument/option
func (f Flag) String() string {
	return f.name
}

//
func getIndex(i []int) int {
	index := 0
	if len(i) > 0 {
		index = i[0]
	}
	return index
}
