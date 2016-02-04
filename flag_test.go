package cli

import "testing"

func TestIsOptionalArgument(t *testing.T) {
	flag := Flag{
		name:    "Optional Argument Flag",
		options: optional,
	}

	if !flag.isOptionalArgument() {
		t.Errorf("Flag `%s` should be optional argument!", flag)
	}
}

func TestIsArrayArgument(t *testing.T) {
	flag := Flag{
		name:    "Array Argument Flag",
		options: isArray,
	}

	if !flag.isArrayArgument() {
		t.Errorf("Flag `%s` should be array argument!", flag)
	}
}

func TestIsRequiredArgument(t *testing.T) {
	flag := Flag{
		name:    "Required Argument Flag",
		options: required,
	}

	if !flag.isRequiredArgument() {
		t.Errorf("Flag `%s` should be required!", flag)
	}
}

func TestFlagRequiredValue(t *testing.T) {
	flag := Flag{
		name:    "RequiredFlag",
		options: valueRequired | valueArray,
	}

	if !flag.isValueRequired() {
		t.Errorf("Flag `%s` should have a required value!", flag)
	}
}

func TestFlagArrayValue(t *testing.T) {
	flag := Flag{
		name:    "Array Value Flag",
		options: valueRequired | valueArray,
	}

	if !flag.isValueArray() {
		t.Errorf("Flag `%s` should have an array value!", flag)
	}
}

func TestFlagOptionalValue(t *testing.T) {
	flag := Flag{
		name:    "RequiredFlag",
		options: valueRequired | valueArray | valueOptional,
	}

	if !flag.isValueOptional() {
		t.Errorf("Flag `%s` should have an optional value!", flag)
	}
}

func TestAcceptValue(t *testing.T) {
	flag := Flag{
		name:    "RequiredFlag",
		options: valueRequired | valueArray,
	}

	if !flag.acceptValue() {
		t.Errorf("Flag `%s` should have a required value!", flag)
	}
}
