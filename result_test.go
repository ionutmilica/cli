package cli

import "testing"

func TestNewResult(t *testing.T) {
	r := Result{}

	if len(r) != 0 {
		t.Errorf("Result should have 0 items but got %d!", len(r))
	}
}

func TestHasMethod(t *testing.T) {
	r := Result{}
	if r.Has(0) {
		t.Errorf("t.Has(0) found element but expected 0!")
	}

	r2 := Result{"hello"}

	if !r2.Has(0) {
		t.Error("t.Has(0) expected true but got false!")
	}
}

func TestAppendAndGetResult(t *testing.T) {
	r := Result{}
	r.Append("My item")
	r.Append("My dream")

	if len(r) != 2 {
		t.Errorf("Expected result length `%d` but got `%d`!", 2, len(r))
	}

	// Test first item
	item, err := r.Str()

	if err != nil {
		t.Errorf("Got error: %s!", err.Error())
	}

	if item != "My item" {
		t.Errorf("r.Get() expected `%s` but got `%s`!", "My item", item)
	}

	// Test second item
	item, err = r.Str(1)

	if err != nil {
		t.Errorf("Got error: %s!", err.Error())
	}

	if item != "My dream" {
		t.Errorf("r.Get() expected `%s` but got `%s`!", "My dream", item)
	}

	// Test out of bounds item
	item, err = r.Str(2)

	if err == nil {
		t.Errorf("Got no error but expected one. Item: `%s`!", item)
	}
}

func TestIntResult(t *testing.T) {
	r := Result{}
	r.Append("-1")
	r.Append("Lama")

	if len(r) != 2 {
		t.Errorf("Expected result length `%d` but got `%d`!", 2, len(r))
	}

	// Test first item
	item, err := r.Int()

	if err != nil {
		t.Errorf("Got error: %s!", err.Error())
	}

	if item != -1 {
		t.Errorf("r.Get() expected `%d` but got `%d`!", -1, item)
	}

	// Test second item
	item, err = r.Int(1)

	if err == nil {
		t.Error("Got no error but expected int parse error!")
	}

	// Test out of bounds item
	item, err = r.Int(2)

	if err == nil {
		t.Errorf("Got no error but expected one. Item: `%d`!", item)
	}
}
