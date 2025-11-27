package goda

// TemporalValue represents the value of a temporal field, with support for
// overflow and unsupported field states.
//
// TemporalValue is returned by TemporalAccessor.GetField() and used by
// temporal types that support field-based mutation. It provides a safe way
// to handle field values that may be out of range or unsupported.
//
// The value can be in one of three states:
//   - Valid: Contains a normal field value
//   - Overflow: The value is too large to fit in an int64
//   - Unsupported: The field is not supported by the temporal type
//
// Use Valid(), Overflow(), and Unsupported() methods to check the state.
// Use Int() or Int64() to get the value when Valid() returns true.
type TemporalValue struct {
	v           int64
	overflow    bool
	unsupported bool
}

// Int64 returns the value as an int64.
// Only call this when Valid() returns true.
func (t TemporalValue) Int64() int64 {
	return t.v
}

// Int returns the value as an int.
// Only call this when Valid() returns true and the value fits in an int.
func (t TemporalValue) Int() int {
	return int(t.v)
}

// Overflow returns true if the field value overflowed and cannot be represented.
// When this returns true, the value is undefined.
func (t TemporalValue) Overflow() bool {
	return t.overflow
}

// Unsupported returns true if the field is not supported by the temporal type.
// When this returns true, no meaningful value is available.
func (t TemporalValue) Unsupported() bool {
	return t.unsupported
}

// Valid returns true if the value is valid and can be safely accessed.
// Valid() is equivalent to !Unsupported() && !Overflow().
func (t TemporalValue) Valid() bool {
	return !t.unsupported && !t.overflow
}

// TemporalValueOf creates a new TemporalValue from an integer type.
// The resulting value will be valid (neither overflow nor unsupported).
func TemporalValueOf[T interface {
	int | int64 | int32 | int16 | int8
}](v T) TemporalValue {
	return TemporalValue{v: int64(v)}
}
