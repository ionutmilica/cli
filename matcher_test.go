package cli

import "testing"

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
