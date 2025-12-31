package goda

import (
	"errors"
	"fmt"
)

// ErrUnsupported indicates that the specified operation or field is not supported.
var ErrUnsupported = errors.ErrUnsupported

// ErrOutOfRange indicates that the specified value is out of range.
var ErrOutOfRange = errors.New("out of range")

// ErrArithmeticOverflow indicates that the result of an arithmetic operation overflows.
var ErrArithmeticOverflow = errors.New("arithmetic overflow")

// newError creates a new Error with the given format and arguments.
// All error messages are prefixed with "goda: ".
func newError(format string, a ...any) error {
	return &Error{message: fmt.Sprintf(format, a...)}
}

// sqlScannerDefaultBranch creates an error for unsupported SQL scan types.
func sqlScannerDefaultBranch(value any) error {
	return newError("cannot scan value of type %T", value)
}

const (
	errReasonInvalidField = iota + 1
	errReasonUnsupportedField
	errReasonOutOfRange
	errReasonArithmeticOverflow
	errReasonParseFailed
	errReasonInvalidZoneId
)

// Error is the error type used by this package.
// It wraps error messages with the "goda: " prefix.
type Error struct {
	field      Field
	int64Value int64
	message    string
	cause      error
	typeNameId int8
	funcNameId int8
	reason     int8
}

// Error implements the error interface.
func (e Error) Error() string {
	var text string
	var cause = e.cause
	switch e.reason {
	case errReasonInvalidField:
		text = fmt.Sprintf("goda: invalid field (value=%d)", int64(e.field))
	case errReasonUnsupportedField:
		text = fmt.Sprintf("goda: unsupported field %s", e.field)
	case errReasonOutOfRange:
		fr := e.field.fieldRange()
		text = fmt.Sprintf("goda: invalid value of %s (valid range %d - %d): %d", e.field, fr.Min, fr.Max, e.int64Value)
	case errReasonArithmeticOverflow:
		text = "goda: arithmetic overflow"
	case errReasonParseFailed:
		text = "goda: parse user input failed"
		if cause != nil {
			text += ", " + cause.Error()
			cause = nil
		}
	case errReasonInvalidZoneId:
		text = "goda: invalid zone id"
	default:
		text = "goda: " + e.message
	}
	if e.typeNameId != 0 {
		text += " at " + tyNames[e.typeNameId] + "/" + fnNames[e.funcNameId]
	}
	if cause != nil {
		text += ", caused by: " + cause.Error()
	}
	return text
}

func (e Error) Unwrap() error {
	if e.cause == nil {
		switch e.reason {
		case errReasonOutOfRange:
			return ErrOutOfRange
		case errReasonArithmeticOverflow:
			return ErrArithmeticOverflow
		case errReasonUnsupportedField:
			return ErrUnsupported
		}
	}
	return e.cause
}

func overflowError() error {
	return &Error{reason: errReasonArithmeticOverflow}
}

func fieldOutOfRangeError(field Field, value int64) error {
	return &Error{reason: errReasonOutOfRange, field: field, int64Value: value}
}

func unsupportedField(field Field) error {
	return &Error{reason: errReasonUnsupportedField, field: field}
}

func invalidFieldError(field Field) error {
	return &Error{reason: errReasonInvalidField, field: field}
}

func parseFailedError(userInput []byte) error {
	return &Error{reason: errReasonParseFailed}
}

func parseFailedErrorWithCause(userInput []byte, cause error) error {
	return &Error{reason: errReasonParseFailed, cause: cause}
}

func deferOpInParse(userInput []byte, e *error) {
	if *e == nil {
		return
	}
	//goland:noinspection GoTypeAssertionOnErrors
	if _, ok := (*e).(*Error); !ok {
		*e = parseFailedErrorWithCause(userInput, *e)
	}
}
