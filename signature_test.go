package cli

import (
	"testing"
)

func makeCmd(signature string) Command {
	return Command{
		Signature: signature,
	}
}

func toFlags(signature string) []*Flag {
	cmd := makeCmd(signature)
	cmd.parse()
	return cmd.Flags
}

/** Arguments **/

func TestPatternParseArgumentFlag(t *testing.T) {
	flags := toFlags("{user} {other}")

	expectedFlags := []string{"user", "other"}

	if len(flags) != len(expectedFlags) {
		t.Errorf("Expected `%d` value flags but got `%d`!", len(expectedFlags), len(flags))
		return
	}

	for i, name := range expectedFlags {
		if !flags[i].isArgument() || flags[i].name != name {
			t.Errorf("Expected Name: `%s`, Kind: `%d`, but got: Name: `%s`, Kind: `%d`", name, argumentFlag, flags[i].name, flags[i].kind)
		}
	}
}

func TestPatternParseArgumentWithDescription(t *testing.T) {
	flags := toFlags("{file : This argument it's used for amazing stuff!}")

	if len(flags) < 1 {
		t.Errorf("Expected 1 value flag but got `%d`!", len(flags))
		return
	}

	if !flags[0].isArgument() {
		t.Errorf("Argument kind it's not `Argument` but got: %d", flags[0].kind)
	}

	if flags[0].description != "This argument it's used for amazing stuff!" {
		t.Errorf("Argument description is wrong! Got: %s", flags[0].description)
	}
}

func TestPatternWithOptionalArgument(t *testing.T) {
	flags := toFlags("{user?}")

	if len(flags) < 1 {
		t.Errorf("Expected 1 value flag but got `%d`!", len(flags))
		return
	}

	if !flags[0].isOptional() {
		t.Errorf("Argument `user` should have optional option but got: %d", flags[0].options)
	}
}

func TestPatternWithArrayArgument(t *testing.T) {
	flags := toFlags("{user*}")

	if len(flags) < 1 {
		t.Errorf("Expected 1 value flag but got `%d`!", len(flags))
		return
	}

	if !flags[0].isArray() || !flags[0].isRequired() {
		t.Errorf("Argument `user` should be array and required but got: %d", flags[0].options)
	}
}

func TestPatternWithOptionalArrayArgument(t *testing.T) {
	flags := toFlags("{user?*}")

	if len(flags) < 1 {
		t.Errorf("Expected 1 value flag but got `%d`!", len(flags))
		return
	}

	if !flags[0].isArray() || !flags[0].isOptional() {
		t.Errorf("Argument `user` should be optional array but got: %d", flags[0].options)
	}
}

func TestDefaultValueArgument(t *testing.T) {
	flags := toFlags("{user=ionut}")

	if len(flags) < 1 {
		t.Errorf("Expected 1 value flag but got `%d`!", len(flags))
		return
	}

	if !flags[0].isOptional() || flags[0].value != "ionut" {
		t.Errorf("Argument `user` should be optional and have default value=ionut got: %d, val=%s", flags[0].options, flags[0].value)
	}
}

/** LONG OPTIONS and Short options **/

func TestOptionParse(t *testing.T) {
	flags := toFlags("{--test}")

	if len(flags) < 1 {
		t.Errorf("Expected 1 value flag but got `%d`!", len(flags))
		return
	}
	if !flags[0].isLongOption() {
		t.Errorf("Argument `test` should be longOptionFlag but got: %d", flags[0].kind)
	}

	// Short option
	flags = toFlags("{-t}")

	if len(flags) < 1 {
		t.Errorf("Expected 1 value flag but got `%d`!", len(flags))
		return
	}
	if !flags[0].isOption() {
		t.Errorf("Argument `t` should be optionFlag but got: %d", flags[0].kind)
	}
}

func TestOptionWithOptionalValue(t *testing.T) {
	flags := toFlags("{--test=}")
	if len(flags) < 1 {
		t.Errorf("Expected 1 value flag but got `%d`!", len(flags))
		return
	}

	if !flags[0].isLongOption() || !flags[0].isOptional() {
		t.Errorf("Argument `test` should be longOptionFlag and have optional value but got: %d, %d", flags[0].kind, flags[0].options)
	}

	// Short option
	flags = toFlags("{-t=}")
	if len(flags) < 1 {
		t.Errorf("Expected 1 value flag but got `%d`!", len(flags))
		return
	}
	if !flags[0].isOption() || !flags[0].isOptional() {
		t.Errorf("Argument `t` should be optionFlag and have optional value but got: %d, %d", flags[0].kind, flags[0].options)
	}
}

func TestOptionWithOptionalArrayValue(t *testing.T) {
	flags := toFlags("{--test=*}")

	if len(flags) < 1 {
		t.Errorf("Expected 1 value flag but got `%d`!", len(flags))
		return
	}

	if !flags[0].isLongOption() || !flags[0].isOptional() || !flags[0].isArray() {
		t.Errorf("Argument `test` should be longOptionFlag and have optional value but got: %d, %d", flags[0].kind, flags[0].options)
	}

	// Short option
	flags = toFlags("{-t=*}")
	if len(flags) < 1 {
		t.Errorf("Expected 1 value flag but got `%d`!", len(flags))
		return
	}
	if !flags[0].isOption() || !flags[0].isOptional() || !flags[0].isArray() {
		t.Errorf("Argument `t` should be optionFlag and have optional value but got: %d, %d", flags[0].kind, flags[0].options)
	}
}

func TestOptionWithRequiredArrayValue(t *testing.T) {
	flags := toFlags("{--test=+}")
	if len(flags) < 1 {
		t.Errorf("Expected 1 value flag but got `%d`!", len(flags))
		return
	}
	if !flags[0].isLongOption() || !flags[0].isRequired() || !flags[0].isArray() {
		t.Errorf("Argument `test` should be longOptionFlag and have required value but got: %d, %d", flags[0].kind, flags[0].options)
	}

	// Short option
	flags = toFlags("{-t=+}")
	if len(flags) < 1 {
		t.Errorf("Expected 1 value flag but got `%d`!", len(flags))
		return
	}
	if !flags[0].isOption() || !flags[0].isRequired() || !flags[0].isArray() {
		t.Errorf("Argument `t` should be optionFlag and have required value but got: %d, %d", flags[0].kind, flags[0].options)
	}
}

func TestOptionWithDefaultValue(t *testing.T) {
	flags := toFlags("{--test=ionut}")
	if len(flags) < 1 {
		t.Errorf("Expected 1 value flag but got `%d`!", len(flags))
		return
	}
	if !flags[0].isLongOption() || !flags[0].isOptional() || flags[0].value != "ionut" {
		t.Errorf("Argument `test` should be longOptionFlag and have optional value with default=ionut but got: %d, %d, val=%s", flags[0].kind, flags[0].options, flags[0].value)
	}

	// Short option
	flags = toFlags("{-t=ionut}")
	if len(flags) < 1 {
		t.Errorf("Expected 1 value flag but got `%d`!", len(flags))
		return
	}
	if !flags[0].isOption() || !flags[0].isOptional() || flags[0].value != "ionut" {
		t.Errorf("Argument `t` should be longOptionFlag and have optional value with default=ionut but got: %d, %d, val=%s", flags[0].kind, flags[0].options, flags[0].value)
	}
}

func TestExtractDescriptionFunction(t *testing.T) {
	name, description := extractDescription("ion : Hello world!")

	if name != "ion" || description != "Hello world!" {
		t.Errorf("Expected [%s, %s] but got [%s, %s]", "ion", "Hello world!", name, description)
	}
}
