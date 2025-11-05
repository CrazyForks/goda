package goda

import (
	"fmt"
)

// newError creates a new Error with the given format and arguments.
// All error messages are prefixed with "goda: ".
func newError(format string, a ...any) error {
	return &Error{"goda: " + fmt.Sprintf(format, a...)}
}

// unmarshalError creates an error for invalid unmarshaling input.
func unmarshalError(userInput []byte) error {
	return newError("unable to unmarshal user input: %q", string(userInput))
}

// sqlScannerDefaultBranch creates an error for unsupported SQL scan types.
func sqlScannerDefaultBranch(value any) error {
	return newError("cannot scan value of type %T", value)
}

// Error is the error type used by this package.
// It wraps error messages with the "goda: " prefix.
type Error struct {
	message string
}

// Error implements the error interface.
func (e Error) Error() string {
	return e.message
}
