package cli

import (
	"errors"
	"fmt"
	"strconv"
)

// Every option or argument will have a result object that will contains all the matched data
// i.e: -i => [file1, file2, fil3]
type Result []string

// Allows adding of new items into the result list
func (r *Result) Append(item ...string) {
	*r = append(*r, item...)
}

// Check if there is an item at position `i`
func (r Result) Has(i int) bool {
	if i < 0 || i > len(r)-1 {
		return false
	}
	return true
}

// Get the first(or by pos) item as string
func (r Result) Str(i ...int) (string, error) {
	pos := getPos(i)

	if !r.Has(pos) {
		return "", errors.New("Item not found!")
	}

	return r[pos], nil
}

// Convert the first (or specified) item from string to int
func (r Result) Int(i ...int) (int, error) {
	pos := getPos(i)

	if !r.Has(pos) {
		return -1, errors.New("Item not found!")
	}

	return strconv.Atoi(r[pos])
}

// Returns the content of the Result as string slice
func (r Result) StrSlice() []string {
	return r
}

// Override slice String method
func (r *Result) String() string {
	return fmt.Sprintf("%s", *r)
}

// Get item pos by array or args
func getPos(i []int) int {
	if len(i) > 0 {
		return i[0]
	}
	return 0
}
