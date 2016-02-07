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
