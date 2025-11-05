package goda

import (
	"fmt"
)

func newError(format string, a ...any) error {
	return &Error{"goda: " + fmt.Sprintf(format, a...)}
}

func unmarshalError(userInput []byte) error {
	return newError("unable to unmarshal user input: %q", string(userInput))
}

func sqlScannerDefaultBranch(value any) error {
	return newError("cannot scan value of type %T", value)
}

type Error struct {
	message string
}

func (e Error) Error() string {
	return e.message
}
