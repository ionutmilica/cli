package cli

import (
	"reflect"
	"testing"
)

func makeCtx() *Context {
	return &Context{
		Arguments: make(map[string][]string, 0),
		Options:   make(map[string][]string),
	}
}

type ArgumentParseTest struct {
	name       string
	input      string
	values     []string
	shouldFail bool
}

func argumentTest(t *testing.T, flagMgr *FlagMgr, steps []ArgumentParseTest, ctx *Context) {
	for _, step := range steps {
		err := ctx.parseArgument(flagMgr, step.input)

		if step.shouldFail && err == nil {
			t.Errorf("Test should fail for `%s`", step.input)
		}

		if !step.shouldFail && err != nil {
			t.Errorf("Got error: `%s` for input `%s`", err, step.input)
		}

		if step.shouldFail {
			continue
		}

		if _, ok := ctx.Arguments[step.name]; !ok {
			t.Errorf("Argument `%s` not found!", step.name)
		}

		if !reflect.DeepEqual(step.values, ctx.Arguments[step.name]) {
			t.Errorf("Expected values [%s] but got [%s]", step.values, ctx.Arguments[step.name])
		}
		ctx.cursor++
	}
}

func longOptionTest(t *testing.T, flagMgr *FlagMgr, steps []ArgumentParseTest, ctx *Context) {
	for _, step := range steps {
		err := ctx.parseLongOption(flagMgr, step.input)

		if step.shouldFail && err == nil {
			t.Errorf("Test should fail for `%s`", step.input)
		}

		if !step.shouldFail && err != nil {
			t.Errorf("Got error: `%s` for input `%s`", err, step.input)
		}

		if step.shouldFail {
			continue
		}

		if _, ok := ctx.Options[step.name]; !ok {
			t.Errorf("Option `%s` not found!", step.name)
		}

		if !reflect.DeepEqual(step.values, ctx.Options[step.name]) && len(step.values) > 0 && len(ctx.Options[step.name]) > 0 {
			t.Errorf("Expected values [%s] but got [%s]", step.values, ctx.Options[step.name])
		}
		ctx.cursor++
	}
}

func TestParseArgument(t *testing.T) {
	flags := []*Flag{
		&Flag{
			name:    "foo",
			kind:    argumentFlag,
			options: required,
		},
		&Flag{
			name:    "bar",
			kind:    argumentFlag,
			options: required,
		},
	}
	flagMgr := newFlagMgr(flags)

	ctx := makeCtx()
	ctx.cursor = 0

	steps := []ArgumentParseTest{
		ArgumentParseTest{
			name:       "foo",
			input:      "val1",
			values:     []string{"val1"},
			shouldFail: false,
		},
		ArgumentParseTest{
			name:       "bar",
			input:      "val2",
			values:     []string{"val2"},
			shouldFail: false,
		},
		ArgumentParseTest{
			name:       "fuu",
			input:      "val1",
			values:     []string{"val1"},
			shouldFail: true,
		},
	}

	argumentTest(t, flagMgr, steps, ctx)
}

func TestParseArgumentArray(t *testing.T) {
	flags := []*Flag{
		&Flag{
			name:    "foo",
			kind:    argumentFlag,
			options: required,
		},
		&Flag{
			name:    "bar",
			kind:    argumentFlag,
			options: required | isArray,
		},
	}
	flagMgr := newFlagMgr(flags)

	ctx := makeCtx()
	ctx.cursor = 0

	steps := []ArgumentParseTest{
		ArgumentParseTest{
			name:       "foo",
			input:      "val1",
			values:     []string{"val1"},
			shouldFail: false,
		},
		ArgumentParseTest{
			name:       "bar",
			input:      "val2",
			values:     []string{"val2"},
			shouldFail: false,
		},
		ArgumentParseTest{
			name:       "bar",
			input:      "val3",
			values:     []string{"val2", "val3"},
			shouldFail: false,
		},
	}

	argumentTest(t, flagMgr, steps, ctx)
}

func TestParseArgumentWithMissingArguments(t *testing.T) {
	flags := []*Flag{
		&Flag{
			name:    "foo",
			kind:    argumentFlag,
			options: required,
		},
		&Flag{
			name:    "bar",
			kind:    argumentFlag,
			options: required,
		},
	}
	flagMgr := newFlagMgr(flags)

	ctx := makeCtx()
	ctx.cursor = 0
	err := ctx.parse([]string{}, flagMgr)

	if err == nil || err.Error() != "Not enough arguments (missing: `foo`)." {
		t.Errorf("Parse with 1 missing argument failed with wrong error: %s", err)
	}
}

func TestParseLongOption(t *testing.T) {
	flags := []*Flag{
		&Flag{
			name:    "file",
			kind:    longOptionFlag,
			options: valueRequired,
		},
		&Flag{
			name:    "use-something",
			kind:    longOptionFlag,
			options: valueNone,
		},
		&Flag{
			name:    "output",
			kind:    longOptionFlag,
			options: valueOptional,
		},
	}
	flagMgr := newFlagMgr(flags)

	ctx := makeCtx()
	ctx.cursor = 0
	ctx.args = []string{"--file", "test.txt", "--use-something", "--output=file.out"}

	steps := []ArgumentParseTest{
		ArgumentParseTest{
			name:       "file",
			input:      "--file",
			values:     []string{"test.txt"},
			shouldFail: false,
		},
		ArgumentParseTest{
			name:       "use-something",
			input:      "--use-something",
			values:     []string{},
			shouldFail: false,
		},
		ArgumentParseTest{
			name:       "output",
			input:      "--output=file.out",
			values:     []string{"file.out"},
			shouldFail: false,
		},
	}

	longOptionTest(t, flagMgr, steps, ctx)
}

func TestParseLongOptionWithDefaultValue(t *testing.T) {
	flags := []*Flag{
		&Flag{
			name:    "file",
			kind:    longOptionFlag,
			options: valueOptional,
			value:   "default",
		},
	}
	flagMgr := newFlagMgr(flags)

	ctx := makeCtx()
	ctx.cursor = 0
	ctx.args = []string{"--file"}

	steps := []ArgumentParseTest{
		ArgumentParseTest{
			name:       "file",
			input:      "--file",
			values:     []string{"default"},
			shouldFail: false,
		},
	}

	longOptionTest(t, flagMgr, steps, ctx)
}

func TestParseOptionWithMissingValue(t *testing.T) {
	flags := []*Flag{
		&Flag{
			name:    "foo",
			kind:    longOptionFlag,
			options: valueRequired,
		},
	}
	flagMgr := newFlagMgr(flags)

	ctx := makeCtx()
	ctx.cursor = 0
	err := ctx.parse([]string{"--foo"}, flagMgr)

	if err == nil || err.Error() != "The `--foo` option requres a value!" {
		t.Errorf("Parse with missing option value failed with wrong error: %s", err)
	}
}

func TestParseOptionInvalidOptionProvided(t *testing.T) {
	flags := []*Flag{
		&Flag{
			name:    "foo",
			kind:    longOptionFlag,
			options: valueOptional,
		},
	}
	flagMgr := newFlagMgr(flags)

	ctx := makeCtx()
	ctx.cursor = 0
	err := ctx.parse([]string{"--boo"}, flagMgr)

	if err == nil || err.Error() != "The `--boo` option does not exist." {
		t.Errorf("Parse with invalid option failed with wrong error: %s", err)
	}
}

func TestParseOptionWithNoWantedValue(t *testing.T) {
	flags := []*Flag{
		&Flag{
			name:    "foo",
			kind:    longOptionFlag,
			options: valueNone,
		},
	}
	flagMgr := newFlagMgr(flags)

	ctx := makeCtx()
	ctx.cursor = 0
	err := ctx.parse([]string{"--foo=\"bar\""}, flagMgr)

	if err == nil || err.Error() != "The `--foo` option does not accept a value!" {
		t.Errorf("Parse with invalid option failed with wrong error: %s", err)
	}
}
