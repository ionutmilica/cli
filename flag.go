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

type Flag struct {
	kind        int8
	options     int8
	name        string
	description string
	value       string // Default value for flag
}

func (f Flag) isOptionalArgument() bool {
	return f.options&optional == optional
}

func (f Flag) isArrayArgument() bool {
	return f.options&isArray == isArray
}

func (f Flag) isRequiredArgument() bool {
	return f.options&required == required
}

func (f Flag) acceptValue() bool {
	return f.isValueOptional() || f.isValueRequired()
}

func (f Flag) isValueArray() bool {
	return f.options&valueArray == valueArray
}

func (f Flag) isValueOptional() bool {
	return f.options&valueOptional == valueOptional
}

func (f Flag) isValueRequired() bool {
	return f.options&valueRequired == valueRequired
}

func (f Flag) String() string {
	return f.name
}
