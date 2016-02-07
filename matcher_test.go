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
			t.Errorf("Test %d expected to fail but no error got!", i+1)
		}

		if !reflect.DeepEqual(test.arguments, m.ctx.Arguments) {
			t.Errorf("Failed on test %d, got arguments: %s but expected: %s!", i+1, m.ctx.Arguments, test.arguments)
		}

		if !reflect.DeepEqual(test.options, m.ctx.Options) {
			t.Errorf("Failed on test %d, got options: %s but expected: %s!", i+1, m.ctx.Options, test.options)
		}

	}
}

func TestMatchArgument(t *testing.T) {
	tests := []Test{
		// One argument
		Test{
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
			flags:     flags("{a?}"),
			args:      args(),
			fail:      false,
			arguments: map[string][]string{},
			options:   map[string][]string{},
		},

		// Array argument
		Test{
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
			flags:     flags("{a*}"),
			args:      args(),
			fail:      true,
			arguments: map[string][]string{},
			options:   map[string][]string{},
		},

		// Optional Array argument with 0 received

		Test{
			flags:     flags("{a?*}"),
			args:      args(),
			fail:      false,
			arguments: map[string][]string{},
			options:   map[string][]string{},
		},

		// Optional Array argument
		Test{
			flags: flags("{a?*}"),
			args:  args("a", "b", "c", "d"),
			fail:  false,
			arguments: map[string][]string{
				"a": []string{"a", "b", "c", "d"},
			},
			options: map[string][]string{},
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

		// Option requires a value
		Test{
			flags:     flags("{--file=+}"),
			args:      args("--file"),
			fail:      true,
			arguments: map[string][]string{},
			options:   map[string][]string{},
		},

		// Option default value
		/*
			@todo: Implement this in signature parser
			Test{
				flags:     flags("{--file=ion}"),
				args:      args("--file"),
				fail:      false,
				arguments: map[string][]string{},
				options: map[string][]string{
					"file": []string{"ion"},
				},
			},*/

	}

	test(t, tests)
}
