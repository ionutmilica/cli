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

func TestPatternParseArgumentFlag(t *testing.T) {
	flags := toFlags("{user} {other}")

	expectedFlags := []string{"user", "other"}

	if len(flags) != len(expectedFlags) {
		t.Errorf("Expected `%d` value flags but got `%d`!", len(expectedFlags), len(flags))
		return
	}

	for i, name := range expectedFlags {
		if flags[i].kind != argumentFlag || flags[i].name != name {
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

	if flags[0].kind != argumentFlag {
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

	if flags[0].options != optional {
		t.Errorf("Argument `user` should have optional option but got: %d", flags[0].options)
	}
}

func TestPatternWithArrayArgument(t *testing.T) {
	flags := toFlags("{user*}")

	if len(flags) < 1 {
		t.Errorf("Expected 1 value flag but got `%d`!", len(flags))
		return
	}

	if flags[0].options&isArray != isArray || flags[0].options&required != required {
		t.Errorf("Argument `user` should be array and required but got: %d", flags[0].options)
	}
}

func TestPatternWithOptionalArrayArgument(t *testing.T) {
	flags := toFlags("{user?*}")

	if len(flags) < 1 {
		t.Errorf("Expected 1 value flag but got `%d`!", len(flags))
		return
	}

	if flags[0].options&isArray != isArray || flags[0].options&optional != optional {
		t.Errorf("Argument `user` should be optional array but got: %d", flags[0].options)
	}
}

func TestPatternParse(t *testing.T) {
	//newPattern("{user} {--lib=test} {--Q|queue")
	//fmt.Println(pattern)
}
