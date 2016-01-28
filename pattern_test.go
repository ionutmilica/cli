package cli

import (
	"testing"
)

func TestPatternParseValueFlag(t *testing.T) {
	flags := newPattern("{user} {other}").flags
	expectedFlags := []string{"user", "other"}

	if len(flags) != len(expectedFlags) {
		t.Errorf("Expected `%d` value flags but got `%d`!", len(expectedFlags), len(flags))
		return
	}

	for i, name := range expectedFlags {
		if flags[i].kind != valueFlag || flags[i].name != name {
			t.Errorf("Expected Name: `%s`, Kind: `%d`, but got: Name: `%s`, Kind: `%d`", name, valueFlag, flags[i].name, flags[i].kind)
		}
	}
}

func TestPatternParseOptionFlag(t *testing.T) {
	flags := newPattern("{-v} {-o}").flags
	expectedFlags := []string{"v", "o"}

	if len(flags) != len(expectedFlags) {
		t.Errorf("Expected `%d` value flags but got `%d`!", len(expectedFlags), len(flags))
		return
	}

	for i, name := range expectedFlags {
		if flags[i].kind != optionFlag || flags[i].name != name {
			t.Errorf("Expected Name: `%s`, Kind: `%d`, but got: Name: `%s`, Kind: `%d`", name, valueFlag, flags[i].name, flags[i].kind)
		}
	}
}

func TestPatternParse(t *testing.T) {
	//newPattern("{user} {--lib=test} {--Q|queue")
	//fmt.Println(pattern)
}
