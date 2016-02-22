package cli

import (
	_ "fmt"
	"testing"
)

func TestFlagList(t *testing.T) {
	flagList := flags("{a} {b} {--ion}")

	if len(flagList) != 3 {
		t.Errorf("Expected `%d` flags in the list but got: %d", 3, len(flagList))
	}

	// Test required flags
	required := flagList.requiredArgs()

	if len(required) != 2 {
		t.Errorf("Expected 2 reguired flags but got %d", len(required))
	}

	if arg := flagList.argument(0); arg == nil || arg.name != "a" {
		t.Errorf("First argument expected to be 'a' but got %s", arg)
	}

	if opt := flagList.option("ion"); opt == nil || opt.name != "ion" {
		t.Errorf("Find option by key `%s` expected `%s` value but got `%s`", "ion", "ion", opt)
	}
}

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
