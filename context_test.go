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
