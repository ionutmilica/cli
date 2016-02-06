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

func (f Flag) acceptValue() bool {
	return f.isOptional() || f.isRequired()
}

func (f Flag) isArray() bool {
	if f.kind == argumentFlag {
		return f.options&isArray == isArray
	}
	return f.options&valueArray == valueArray
}

func (f Flag) isOptional() bool {
	if f.kind == argumentFlag {
		return f.options&optional == optional
	}
	return f.options&valueOptional == valueOptional
}

func (f Flag) isRequired() bool {
	if f.kind == argumentFlag {
		return f.options&required == required
	}
	return f.options&valueRequired == valueRequired
}

func (f Flag) String() string {
	return f.name
}
