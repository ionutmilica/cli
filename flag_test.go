package cli

import "testing"

func TestIsOptionalArgument(t *testing.T) {
	flag := Flag{
		kind:    argumentFlag,
		name:    "Optional Argument Flag",
		options: optional,
	}

	if !flag.isOptional() {
		t.Errorf("Flag `%s` should be optional argument!", flag)
	}
}

func TestIsArrayArgument(t *testing.T) {
	flag := Flag{
		kind:    argumentFlag,
		name:    "Array Argument Flag",
		options: isArray,
	}

	if !flag.isArray() {
		t.Errorf("Flag `%s` should be array argument!", flag)
	}
}

func TestIsRequiredArgument(t *testing.T) {
	flag := Flag{
		kind:    argumentFlag,
		name:    "Required Argument Flag",
		options: required,
	}

	if !flag.isRequired() {
		t.Errorf("Flag `%s` should be required!", flag)
	}
}

func TestFlagRequiredValue(t *testing.T) {
	flag := Flag{
		kind:    optionFlag,
		name:    "RequiredFlag",
		options: valueRequired | valueArray,
	}

	if !flag.isRequired() {
		t.Errorf("Flag `%s` should have a required value!", flag)
	}
}

func TestFlagArrayValue(t *testing.T) {
	flag := Flag{
		kind:    optionFlag,
		name:    "Array Value Flag",
		options: valueRequired | valueArray,
	}

	if !flag.isArray() {
		t.Errorf("Flag `%s` should have an array value!", flag)
	}
}

func TestFlagOptionalValue(t *testing.T) {
	flag := Flag{
		kind:    optionFlag,
		name:    "RequiredFlag",
		options: valueRequired | valueArray | valueOptional,
	}

	if !flag.isOptional() {
		t.Errorf("Flag `%s` should have an optional value!", flag)
	}
}

func TestAcceptValue(t *testing.T) {
	flag := Flag{
		kind:    optionFlag,
		name:    "RequiredFlag",
		options: valueRequired | valueArray,
	}

	if !flag.acceptValue() {
		t.Errorf("Flag `%s` should accept value!", flag)
	}
}
