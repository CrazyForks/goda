package goda

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"time"
)

// LocalTime represents a time without a date or time zone in the ISO-8601 calendar system,
// such as 14:30:45.123456789. It stores the hour, minute, second, and nanosecond.
//
// LocalTime is comparable and can be used as a map key.
// The zero value represents an unset time and IsZero returns true for it.
// Note: 00:00:00 (midnight) is a valid time and is different from the zero value.
//
// LocalTime implements sql.Scanner and driver.Valuer for database operations,
// encoding.TextMarshaler and encoding.TextUnmarshaler for text serialization,
// and json.Marshaler and json.Unmarshaler for JSON serialization.
// The text format is HH:MM:SS[.nnnnnnnnn] (ISO 8601).
//
// Internally, v uses bit 63 (the highest bit) as a validity flag.
// If bit 63 is set, the time is valid and bits 0-62 contain nanoseconds since midnight.
// If bit 63 is clear (v == 0), the time is invalid/zero.
// Max nanoseconds in a day: 86,399,999,999,999 (< 2^47), so we have plenty of space.
type LocalTime struct {
	v int64 // bit 63: valid flag; bits 0-62: nanoseconds since midnight
}

const (
	localTimeValidBit  int64 = -0x8000000000000000 // bit 63 set (sign bit in two's complement)
	localTimeValueMask int64 = 0x7FFFFFFFFFFFFFFF  // bits 0-62
)

// Value implements driver.Valuer for database serialization.
func (t *LocalTime) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.String(), nil
}

// Scan implements sql.Scanner for database deserialization.
func (t *LocalTime) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*t = LocalTime{}
		return nil
	case string:
		return t.UnmarshalText([]byte(v))
	case time.Time:
		*t = NewLocalTimeByGoTime(v)
		return nil
	default:
		return sqlScannerDefaultBranch(src)
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *LocalTime) UnmarshalJSON(bytes []byte) error {
	if len(bytes) == 4 && string(bytes) == "null" {
		*t = LocalTime{}
		return nil
	}
	return unmarshalJsonImpl(t, bytes)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (t *LocalTime) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*t = LocalTime{}
		return nil
	}

	if len(text) < 8 {
		return newError("invalid format")
	}

	hour, err := parseInt(text[0:2])
	if err != nil {
		return err
	}
	if text[2] != ':' {
		return newError("expect ':'")
	}
	minute, err := parseInt(text[3:5])
	if err != nil {
		return err
	}
	if text[5] != ':' {
		return newError("expect ':'")
	}
	second, err := parseInt(text[6:8])
	if err != nil {
		return err
	}

	var nano int64 = 0
	var e error
	if len(text) > 8 {
		if text[8] != '.' {
			return newError("expect '.'")
		}
	}
	if len(text) > 9 {
		var nanoBuf [9]byte
		copy(nanoBuf[:], text[9:min(len(text), 18)])
		for i := 8; i >= 0 && nanoBuf[i] == 0; i-- {
			nanoBuf[i] = '0'
		}
		nano, e = parseInt64(nanoBuf[:])
		if e != nil {
			return e
		}
	}

	if hour >= 24 || minute >= 60 || second >= 60 {
		return newError("invalid time value")
	}

	nanos := int64(hour)*int64(time.Hour) + int64(minute)*int64(time.Minute) + int64(second)*int64(time.Second) + nano
	t.v = nanos | localTimeValidBit
	return nil
}

// AppendText in ISO-8601 format
func (t LocalTime) AppendText(b []byte) ([]byte, error) {
	if t.IsZero() {
		return b, nil
	}

	var buf [18]byte
	hour := t.Hour()
	minute := t.Minute()
	second := t.Second()
	nano := t.Nanosecond()

	// Write hours, minutes, seconds
	buf[0] = byte('0' + hour/10)
	buf[1] = byte('0' + hour%10)
	buf[2] = ':'
	buf[3] = byte('0' + minute/10)
	buf[4] = byte('0' + minute%10)
	buf[5] = ':'
	buf[6] = byte('0' + second/10)
	buf[7] = byte('0' + second%10)

	if nano == 0 {
		return append(b, buf[:8]...), nil
	}

	buf[8] = '.'
	// Format nanoseconds
	ns := nano
	for i := 17; i >= 9; i-- {
		buf[i] = byte('0' + ns%10)
		ns /= 10
	}

	// Trim trailing zeros (following Java LocalTime behavior and ISO 8601)
	// This means 100000000ns -> ".1" not ".100"
	l := 18
	for l > 9 && buf[l-1] == '0' {
		l--
	}
	return append(b, buf[:l]...), nil
}

// MarshalText implements encoding.TextMarshaler.
func (t LocalTime) MarshalText() (text []byte, err error) {
	return marshalTextImpl(t)
}

// MarshalJSON implements json.Marshaler.
func (t LocalTime) MarshalJSON() ([]byte, error) {
	return marshalJsonImpl(t)
}

func (t LocalTime) String() string {
	return stringImpl(t)
}

// Hour returns the hour component (0-23).
func (t LocalTime) Hour() int {
	ns := t.v & localTimeValueMask
	return int(ns / int64(time.Hour))
}

// Minute returns the minute component (0-59).
func (t LocalTime) Minute() int {
	ns := t.v & localTimeValueMask
	return int(ns / int64(time.Minute) % 60)
}

// Second returns the second component (0-59).
func (t LocalTime) Second() int {
	ns := t.v & localTimeValueMask
	return int(ns / int64(time.Second) % 60)
}

// Millisecond returns the millisecond component (0-999).
func (t LocalTime) Millisecond() int {
	ns := t.v & localTimeValueMask
	return int(ns / int64(time.Millisecond) % 1000)
}

// Nanosecond returns the nanosecond component (0-999999999).
func (t LocalTime) Nanosecond() int {
	ns := t.v & localTimeValueMask
	return int(ns % int64(time.Second))
}

// NewLocalTime creates a new LocalTime from the specified hour, minute, second, and nanosecond.
// Returns an error if any component is out of range:
// - hour must be 0-23
// - minute must be 0-59
// - second must be 0-59
// - nanosecond must be 0-999999999
func NewLocalTime(hour, minute, second, nanosecond int) (LocalTime, error) {
	if hour < 0 || hour >= 24 {
		return LocalTime{}, newError("hour %d out of range", hour)
	}
	if minute < 0 || minute >= 60 {
		return LocalTime{}, newError("minute %d out of range", minute)
	}
	if second < 0 || second >= 60 {
		return LocalTime{}, newError("second %d out of range", second)
	}
	if nanosecond < 0 || nanosecond >= 1000000000 {
		return LocalTime{}, newError("nanosecond %d out of range", nanosecond)
	}
	nanos := int64(hour)*int64(time.Hour) +
		int64(minute)*int64(time.Minute) +
		int64(second)*int64(time.Second) +
		int64(nanosecond)
	return LocalTime{
		v: nanos | localTimeValidBit,
	}, nil
}

// MustNewLocalTime creates a new LocalTime from the specified hour, minute, second, and nanosecond.
// Panics if any component is out of range. Use NewLocalTime for error handling.
func MustNewLocalTime(hour, minute, second, nanosecond int) LocalTime {
	lt, err := NewLocalTime(hour, minute, second, nanosecond)
	if err != nil {
		panic(err)
	}
	return lt
}

// NewLocalTimeByGoTime creates a LocalTime from a time.Time.
// The date and time zone components are ignored.
// Returns zero value if t.IsZero().
func NewLocalTimeByGoTime(t time.Time) LocalTime {
	if t.IsZero() {
		return LocalTime{}
	}
	nanos := time.Date(1970, 1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC).UnixNano()
	return LocalTime{
		v: nanos | localTimeValidBit,
	}
}

// LocalTimeNow returns the current time in the system's local time zone.
// This is equivalent to NewLocalTimeByGoTime(time.Now()).
// For UTC time, use LocalTimeNowUTC. For a specific timezone, use LocalTimeNowIn.
func LocalTimeNow() LocalTime {
	return NewLocalTimeByGoTime(time.Now())
}

// LocalTimeNowIn returns the current time in the specified time zone.
// This is equivalent to NewLocalTimeByGoTime(time.Now().In(loc)).
func LocalTimeNowIn(loc *time.Location) LocalTime {
	return NewLocalTimeByGoTime(time.Now().In(loc))
}

// LocalTimeNowUTC returns the current time in UTC.
// This is equivalent to NewLocalTimeByGoTime(time.Now().UTC()).
func LocalTimeNowUTC() LocalTime {
	return NewLocalTimeByGoTime(time.Now().UTC())
}

// ParseLocalTime parses a time string in ISO 8601 format (HH:MM:SS[.nnnnnnnnn]).
// Returns an error if the string is invalid or represents an invalid time.
//
// Supported formats:
//   - HH:MM:SS (e.g., "14:30:45")
//   - HH:MM:SS.fff (milliseconds, e.g., "14:30:45.123")
//   - HH:MM:SS.ffffff (microseconds, e.g., "14:30:45.123456")
//   - HH:MM:SS.nnnnnnnnn (nanoseconds, e.g., "14:30:45.123456789")
//
// Example:
//
//	time, err := ParseLocalTime("14:30:45.123")
//	if err != nil {
//	    // handle error
//	}
func ParseLocalTime(s string) (LocalTime, error) {
	var t LocalTime
	err := t.UnmarshalText([]byte(s))
	return t, err
}

// MustParseLocalTime parses a time string in ISO 8601 format (HH:MM:SS[.nnnnnnnnn]).
// Panics if the string is invalid. Use ParseLocalTime for error handling.
//
// Example:
//
//	time := MustParseLocalTime("14:30:45.123456789")
func MustParseLocalTime(s string) LocalTime {
	t, err := ParseLocalTime(s)
	if err != nil {
		panic(err)
	}
	return t
}

// IsZero returns true if this is the zero value of LocalTime.
func (t LocalTime) IsZero() bool {
	return t.v == 0
}

// GoTime converts this time to a time.Time at the Unix epoch date (1970-01-01) in UTC.
// Returns time.Time{} (zero) for zero value.
func (t LocalTime) GoTime() time.Time {
	if t.IsZero() {
		return time.Time{}
	}
	nanos := t.v & localTimeValueMask
	return time.Unix(0, nanos).UTC()
}

// Compare compares this time with another time.
// Returns -1 if this time is before other, 0 if equal, and 1 if after.
// Zero values are considered less than non-zero values.
func (t LocalTime) Compare(other LocalTime) int {
	return doCompare(t, other, compareZero, comparing(LocalTime.Hour), comparing(LocalTime.Minute), comparing(LocalTime.Second), comparing(LocalTime.Nanosecond))
}

// IsBefore returns true if this time is before the specified time.
func (t LocalTime) IsBefore(other LocalTime) bool {
	return t.Compare(other) < 0
}

// IsAfter returns true if this time is after the specified time.
func (t LocalTime) IsAfter(other LocalTime) bool {
	return t.Compare(other) > 0
}

var _ encoding.TextAppender = (*LocalTime)(nil)
var _ fmt.Stringer = (*LocalTime)(nil)
var _ encoding.TextMarshaler = (*LocalTime)(nil)
var _ encoding.TextUnmarshaler = (*LocalTime)(nil)
var _ json.Marshaler = (*LocalTime)(nil)
var _ json.Unmarshaler = (*LocalTime)(nil)
var _ driver.Valuer = (*LocalTime)(nil)
var _ sql.Scanner = (*LocalTime)(nil)
