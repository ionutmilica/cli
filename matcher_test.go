package cli

import (
	"reflect"
	"testing"
)

// Helpers

// Convert signature to flags for compact tests
func flags(signature string) []*Flag {
	cmd := Command{
		Signature: signature,
	}
	cmd.parse()
	return cmd.Flags
}

// args("argument", "--option")
func args(args ...string) []string {
	return args
}

func TestNewMatcher(t *testing.T) {
	m := newMatcher(args("--ion"), flags("{--ion}"))

	if m == nil {
		t.Error("Matched not created!")
	}

	if m.ctx == nil {
		t.Error("Context not created!")
	}
}

func TestHasNext(t *testing.T) {
	m := newMatcher(args("--ion"), flags("{--ion}"))

	if !m.hasNext() {
		t.Errorf("Matcher got hasNext but that's false!")
	}

	m = newMatcher(args(), flags(""))

	if m.hasNext() {
		t.Errorf("Matcher expected false value for hasNext but got true!")
	}
}

func TestCurrentMethod(t *testing.T) {
	m := newMatcher(args("--ion"), flags("{--ion}"))

	if m.current() != "--ion" {
		t.Errorf("Expected current value `--ion` but got `%s`!", m.current())
	}
}

func TestPeekMethod(t *testing.T) {
	m := newMatcher(args("--ion"), flags("{--ion}"))

	if _, err := m.peek(); err == nil {
		t.Errorf("Peek should have failed!")
	}
}

func TestNextMethod(t *testing.T) {
	m := newMatcher(args("--ion"), flags("{--ion}"))

	if m.cursor != 0 {
		t.Errorf("Matcher should have cursor = 0 in initial state!")
	}

	m.next()

	if m.cursor != 1 {
		t.Errorf("Matcher should have value `1` after one advancement but got `%d`!", m.cursor)
	}

	m.next(2)

	if m.cursor != 1 {
		t.Errorf("Matcher should have value `1` after another advance with 2 steps with no more items available, but got `%d`!", m.cursor)
	}
}

func TestValidateMethod(t *testing.T) {
	m := newMatcher(args("file1"), flags("{file} {file2}"))
	err := m.match()

	if err == nil || err.Error() != "Not enough arguments (missing: file2)." {
		t.Errorf("Received wrong error response for not enough provided args!")
	}

	m = newMatcher(args("file1", "file2"), flags("{file} {file2}"))
	err = m.match()

	if err != nil {
		t.Errorf("Received error despite there should be none! (%s", err.Error())
	}
}

func TestToManyArguments(t *testing.T) {
	m := newMatcher(args("file1", "file2"), flags("{file2}"))
	err := m.match()

	if err == nil || err.Error() != "To many arguments!" {
		msg := "nil"
		if err != nil {
			msg = err.Error()
		}
		t.Errorf("Expected to many arguments error but got `%s`!", msg)
	}
}

type Test struct {
	name      string
	args      []string
	flags     []*Flag
	fail      bool
	arguments map[string][]string
	options   map[string][]string
}

func test(t *testing.T, tests []Test) {
	for i, test := range tests {
		m := newMatcher(test.args, test.flags)
		err := m.match()

		if test.fail && err == nil {
			t.Errorf("Test #%d(%s) expected to fail but no error got!", i+1, test.name)
		}

		if !reflect.DeepEqual(test.arguments, m.ctx.Arguments) {
			t.Errorf("Failed on test #%d(%s), got arguments: %s but expected: %s!", i+1, test.name, m.ctx.Arguments, test.arguments)
		}

		if !reflect.DeepEqual(test.options, m.ctx.Options) {
			t.Errorf("Failed on test #%d(%s), got options: %s but expected: %s!", i+1, test.name, m.ctx.Options, test.options)
		}

	}
}

func TestMatchArgument(t *testing.T) {
	tests := []Test{
		// One argument
		Test{
			name:  "Match one argument",
			flags: flags("{file}"),
			args:  args("file"),
			fail:  false,
			arguments: map[string][]string{
				"file": []string{"file"},
			},
			options: map[string][]string{},
		},
		// Two arguments
		Test{
			name:  "Match two arguments",
			flags: flags("{a} {b}"),
			args:  args("ion", "maria"),
			fail:  false,
			arguments: map[string][]string{
				"a": []string{"ion"},
				"b": []string{"maria"},
			},
			options: map[string][]string{},
		},

		// Optional argument
		Test{
			name:  "Match optional argument",
			flags: flags("{a?}"),
			args:  args("test"),
			fail:  false,
			arguments: map[string][]string{
				"a": []string{"test"},
			},
			options: map[string][]string{},
		},

		// Optional argument with 0 provided
		Test{
			name:      "Match optional argument with 0 provided",
			flags:     flags("{a?}"),
			args:      args(),
			fail:      false,
			arguments: map[string][]string{},
			options:   map[string][]string{},
		},

		// Array argument
		Test{
			name:  "Match array argument",
			flags: flags("{a*}"),
			args:  args("ion", "maria"),
			fail:  false,
			arguments: map[string][]string{
				"a": []string{"ion", "maria"},
			},
			options: map[string][]string{},
		},

		// Array argument with 0 received - Fail
		Test{
			name:      "Match argument with 0 provided",
			flags:     flags("{a*}"),
			args:      args(),
			fail:      true,
			arguments: map[string][]string{},
			options:   map[string][]string{},
		},

		// Optional Array argument with 0 received

		Test{
			name:      "Match array argument with 0 provided",
			flags:     flags("{a?*}"),
			args:      args(),
			fail:      false,
			arguments: map[string][]string{},
			options:   map[string][]string{},
		},

		// Optional Array argument
		Test{
			name:  "Match optional array argument",
			flags: flags("{a?*}"),
			args:  args("a", "b", "c", "d"),
			fail:  false,
			arguments: map[string][]string{
				"a": []string{"a", "b", "c", "d"},
			},
			options: map[string][]string{},
		},

		// Argument default value
		Test{
			name:  "Match argument with default value",
			flags: flags("{a==test}"),
			args:  args(),
			fail:  false,
			arguments: map[string][]string{
				"a": []string{"=test"},
			},
			options: map[string][]string{},
		},
	}

	test(t, tests)
}

func TestMatchOption(t *testing.T) {
	tests := []Test{
		// No match
		Test{
			name:      "Match option with 0 provided",
			flags:     flags("{-f}"),
			args:      args("file"),
			fail:      false,
			arguments: map[string][]string{},
			options:   map[string][]string{},
		},

		// One match
		Test{
			name:      "Match one option",
			flags:     flags("{-f}"),
			args:      args("-f"),
			fail:      false,
			arguments: map[string][]string{},
			options: map[string][]string{
				"f": []string{},
			},
		},

		// Match only the flag
		Test{
			name:      "Match the option without the flag",
			flags:     flags("{-f}"),
			args:      args("-f", "youpi"),
			fail:      false,
			arguments: map[string][]string{},
			options: map[string][]string{
				"f": []string{},
			},
		},

		// Match optional value option
		Test{
			name:      "Match optional value option",
			flags:     flags("{-f=}"),
			args:      args("-f"),
			fail:      false,
			arguments: map[string][]string{},
			options: map[string][]string{
				"f": []string{},
			},
		},

		Test{
			name:      "Match optional value option with value provided",
			flags:     flags("{-f=}"),
			args:      args("-f=dada"),
			fail:      false,
			arguments: map[string][]string{},
			options: map[string][]string{
				"f": []string{"dada"},
			},
		},

		Test{
			name:      "Match optional value option with value provided by arg",
			flags:     flags("{-f=}"),
			args:      args("-f", "dada"),
			fail:      false,
			arguments: map[string][]string{},
			options: map[string][]string{
				"f": []string{"dada"},
			},
		},

		// Option does not exist
		Test{
			name:      "Option does not exist",
			flags:     flags(""),
			args:      args("-f"),
			fail:      true,
			arguments: map[string][]string{},
			options:   map[string][]string{},
		},

		// Option does not accept a value
		// Option does not accept a value
		Test{
			name:      "Options does not accept value",
			flags:     flags("{-f}"),
			args:      args("-f=ion.so"),
			fail:      true,
			arguments: map[string][]string{},
			options:   map[string][]string{},
		},

		Test{
			name:      "Options does not accept value",
			flags:     flags("{-f}"),
			args:      args("-f", "ion.so"),
			fail:      false,
			arguments: map[string][]string{},
			options: map[string][]string{
				"f": []string{},
			},
		},

		// Option requires a value
		Test{
			name:      "Option requires value",
			flags:     flags("{-f=+}"),
			args:      args("-f"),
			fail:      true,
			arguments: map[string][]string{},
			options:   map[string][]string{},
		},
		Test{
			flags:     flags("{-f=+}"),
			args:      args("-f", "22", "-f", "something"),
			fail:      false,
			arguments: map[string][]string{},
			options: map[string][]string{
				"f": []string{"22", "something"},
			},
		},

		// Option default value
		Test{
			name:      "Option default value",
			flags:     flags("{-f=ion}"),
			args:      args("--f"),
			fail:      false,
			arguments: map[string][]string{},
			options: map[string][]string{
				"f": []string{"ion"},
			},
		},

		// Array of option values
		Test{
			name:      "Array of option values",
			flags:     flags("{--f=*}"),
			args:      args("-f", "ionut", "-f", "ion"),
			fail:      false,
			arguments: map[string][]string{},
			options: map[string][]string{
				"f": []string{"ionut", "ion"},
			},
		},

		Test{
			flags:     flags("{-f=}"),
			args:      args("-f", "ionut", "-f", "ion"),
			fail:      true,
			arguments: map[string][]string{},
			options:   map[string][]string{},
		},
	}

	test(t, tests)
}

func TestMatchLongOption(t *testing.T) {
	tests := []Test{
		// No match
		Test{
			flags:     flags("{--file}"),
			args:      args("file"),
			fail:      false,
			arguments: map[string][]string{},
			options:   map[string][]string{},
		},
		// Match
		Test{
			flags:     flags("{--file}"),
			args:      args("--file"),
			fail:      false,
			arguments: map[string][]string{},
			options: map[string][]string{
				"file": []string{},
			},
		},

		// Match only the flag
		Test{
			flags:     flags("{--file}"),
			args:      args("--file", "youpi"),
			fail:      false,
			arguments: map[string][]string{},
			options: map[string][]string{
				"file": []string{},
			},
		},

		// Match optional value option
		Test{
			flags:     flags("{--file=}"),
			args:      args("--file"),
			fail:      false,
			arguments: map[string][]string{},
			options: map[string][]string{
				"file": []string{},
			},
		},

		Test{
			flags:     flags("{--file=}"),
			args:      args("--file=dada"),
			fail:      false,
			arguments: map[string][]string{},
			options: map[string][]string{
				"file": []string{"dada"},
			},
		},

		Test{
			flags:     flags("{--file=}"),
			args:      args("--file", "dada"),
			fail:      false,
			arguments: map[string][]string{},
			options: map[string][]string{
				"file": []string{"dada"},
			},
		},

		// Option does not exist
		Test{
			flags:     flags(""),
			args:      args("--file"),
			fail:      true,
			arguments: map[string][]string{},
			options:   map[string][]string{},
		},

		// Option does not accept a value
		Test{
			flags:     flags("{--file}"),
			args:      args("--file=ion.so"),
			fail:      true,
			arguments: map[string][]string{},
			options:   map[string][]string{},
		},

		// Option does not accept a value
		Test{
			flags:     flags("{-file}"),
			args:      args("--file", "ion.so"),
			fail:      false,
			arguments: map[string][]string{},
			options: map[string][]string{
				"file": []string{},
			},
		},

		// Option requires a value
		Test{
			flags:     flags("{--file=+}"),
			args:      args("--file"),
			fail:      true,
			arguments: map[string][]string{},
			options:   map[string][]string{},
		},

		Test{
			flags:     flags("{--file=+}"),
			args:      args("--file=22", "--file", "something"),
			fail:      false,
			arguments: map[string][]string{},
			options: map[string][]string{
				"file": []string{"22", "something"},
			},
		},

		// Option default value
		Test{
			flags:     flags("{--file=ion}"),
			args:      args("--file"),
			fail:      false,
			arguments: map[string][]string{},
			options: map[string][]string{
				"file": []string{"ion"},
			},
		},

		// Array of option values
		Test{
			flags:     flags("{--file=*}"),
			args:      args("--file=ionut", "--file=ion"),
			fail:      false,
			arguments: map[string][]string{},
			options: map[string][]string{
				"file": []string{"ionut", "ion"},
			},
		},

		Test{
			flags:     flags("{--file=}"),
			args:      args("--file=ionut", "--file=ion"),
			fail:      true,
			arguments: map[string][]string{},
			options:   map[string][]string{},
		},
	}

	test(t, tests)
}
