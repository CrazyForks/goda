package goda

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"time"
)

// LocalTime represents a time without a date or time zone,
// such as 14:30:45.123456789. It stores the hour, minute, second, and nanosecond.
//
// LocalTime is comparable and can be used as a map key.
// The zero value represents an unset time and IsZero returns true for it.
// Note: 00:00:00 (midnight) is a valid time and is different from the zero value.
//
// LocalTime implements sql.Scanner and driver.Valuer for database operations,
// encoding.TextMarshaler and encoding.TextUnmarshaler for text serialization,
// and json.Marshaler and json.Unmarshaler for JSON serialization.
//
// Format: HH:mm:ss[.nnnnnnnnn] (e.g., "14:30:45.123456789"). Uses 24-hour format.
// Fractional seconds support nanosecond precision and are aligned to 3-digit boundaries
// (milliseconds, microseconds, nanoseconds) for Java.time compatibility.
// Timezone offsets are not supported.
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
func (t LocalTime) Value() (driver.Value, error) {
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
	case []byte:
		return t.UnmarshalText(v)
	case string:
		return t.UnmarshalText([]byte(v))
	case time.Time:
		*t = LocalTimeOfGoTime(v)
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

	// Trim trailing zeros in groups of 3 (align to milliseconds, microseconds, nanoseconds)
	// This follows Java LocalTime behavior: 100000000ns -> ".100" not ".1"
	l := 18
	for l > 9 && buf[l-1] == '0' {
		l--
	}
	// Align to 3-digit boundaries (9, 12, 15, 18 total length)
	remainder := (l - 9) % 3
	if remainder != 0 {
		l += 3 - remainder
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

// IsSupportedField returns true if the field is supported by LocalTime.
func (t LocalTime) IsSupportedField(field Field) bool {
	switch field {
	case FieldNanoOfSecond, FieldNanoOfDay, FieldMicroOfSecond, FieldMicroOfDay, FieldMilliOfSecond, FieldMilliOfDay, FieldSecondOfMinute, FieldSecondOfDay, FieldMinuteOfHour, FieldMinuteOfDay, FieldHourOfAmPm, FieldClockHourOfAmPm, FieldHourOfDay, FieldClockHourOfDay, FieldAmPmOfDay:
		return true
	default:
		return false
	}
}

// GetField returns the value of the specified field as a TemporalValue.
// This method queries the time for the value of the specified field.
// The returned value may be unsupported if the field is not supported by LocalTime.
//
// If the time is zero (IsZero() returns true), an unsupported TemporalValue is returned.
// For fields not supported by LocalTime (such as date fields), an unsupported TemporalValue is returned.
//
// Supported fields include:
//   - FieldNanoOfSecond: returns the nanosecond within the second (0-999,999,999)
//   - FieldNanoOfDay: returns the nanosecond of day (0-86,399,999,999,999)
//   - FieldMicroOfSecond: returns the microsecond within the second (0-999,999)
//   - FieldMicroOfDay: returns the microsecond of day (0-86,399,999,999)
//   - FieldMilliOfSecond: returns the millisecond within the second (0-999)
//   - FieldMilliOfDay: returns the millisecond of day (0-86,399,999)
//   - FieldSecondOfMinute: returns the second within the minute (0-59)
//   - FieldSecondOfDay: returns the second of day (0-86,399)
//   - FieldMinuteOfHour: returns the minute within the hour (0-59)
//   - FieldMinuteOfDay: returns the minute of day (0-1,439)
//   - FieldHourOfDay: returns the hour of day (0-23)
//   - FieldClockHourOfDay: returns the clock hour of day (1-24)
//   - FieldHourOfAmPm: returns the hour within AM/PM (0-11)
//   - FieldClockHourOfAmPm: returns the clock hour within AM/PM (1-12)
//   - FieldAmPmOfDay: returns AM/PM indicator (0=AM, 1=PM)
//
// Overflow Analysis:
// None of the supported fields can overflow int64:
//   - FieldNanoOfSecond: range 0-999,999,999, cannot overflow
//   - FieldNanoOfDay: max 86,399,999,999,999 (< 2^47), cannot overflow
//   - FieldMicroOfSecond: range 0-999,999, cannot overflow
//   - FieldMicroOfDay: max 86,399,999,999 (< 2^37), cannot overflow
//   - FieldMilliOfSecond: range 0-999, cannot overflow
//   - FieldMilliOfDay: max 86,399,999 (< 2^27), cannot overflow
//   - FieldSecondOfMinute: range 0-59, cannot overflow
//   - FieldSecondOfDay: max 86,399 (< 2^17), cannot overflow
//   - FieldMinuteOfHour: range 0-59, cannot overflow
//   - FieldMinuteOfDay: max 1,439 (< 2^11), cannot overflow
//   - FieldHourOfDay: range 0-23, cannot overflow
//   - FieldClockHourOfDay: range 1-24, cannot overflow
//   - FieldHourOfAmPm: range 0-11, cannot overflow
//   - FieldClockHourOfAmPm: range 1-12, cannot overflow
//   - FieldAmPmOfDay: values 0 or 1, cannot overflow
func (t LocalTime) GetField(field Field) TemporalValue {
	if t.IsZero() {
		return TemporalValue{v: 0, unsupported: true}
	}
	var v int64
	ns := t.v & localTimeValueMask
	switch field {
	case FieldNanoOfSecond:
		// Range: 0-999,999,999, no overflow possible
		v = ns % int64(time.Second)
	case FieldNanoOfDay:
		// Max: 86,399,999,999,999 (< 2^47), no overflow possible
		v = ns
	case FieldMicroOfSecond:
		// Range: 0-999,999, no overflow possible
		v = ns % int64(time.Second) / int64(time.Microsecond)
	case FieldMicroOfDay:
		// Max: 86,399,999,999 (< 2^37), no overflow possible
		v = ns / int64(time.Microsecond)
	case FieldMilliOfSecond:
		// Range: 0-999, no overflow possible
		v = ns % int64(time.Second) / int64(time.Millisecond)
	case FieldMilliOfDay:
		// Max: 86,399,999 (< 2^27), no overflow possible
		v = ns / int64(time.Millisecond)
	case FieldSecondOfMinute:
		// Range: 0-59, no overflow possible
		v = ns / int64(time.Second) % 60
	case FieldSecondOfDay:
		// Max: 86,399 (< 2^17), no overflow possible
		v = ns / int64(time.Second)
	case FieldMinuteOfHour:
		// Range: 0-59, no overflow possible
		v = ns / int64(time.Minute) % 60
	case FieldMinuteOfDay:
		// Max: 1,439 (< 2^11), no overflow possible
		v = ns / int64(time.Minute)
	case FieldHourOfDay:
		// Range: 0-23, no overflow possible
		v = ns / int64(time.Hour)
	case FieldClockHourOfDay:
		// Range: 1-24, no overflow possible
		h := ns / int64(time.Hour)
		if h == 0 {
			v = 24
		} else {
			v = h
		}
	case FieldHourOfAmPm:
		// Range: 0-11, no overflow possible
		v = ns / int64(time.Hour) % 12
	case FieldClockHourOfAmPm:
		// Range: 1-12, no overflow possible
		h := ns / int64(time.Hour) % 12
		if h == 0 {
			v = 12
		} else {
			v = h
		}
	case FieldAmPmOfDay:
		// Values: 0 or 1, no overflow possible
		if ns/int64(time.Hour) < 12 {
			v = 0 // AM
		} else {
			v = 1 // PM
		}
	default:
		return TemporalValue{v: 0, unsupported: true}
	}
	return TemporalValue{v: v}
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

// AtDate combines this time with a date to create a LocalDateTime.
func (t LocalTime) AtDate(date LocalDate) LocalDateTime {
	return LocalDateTime{
		date: date,
		time: t,
	}
}

// PlusHours returns a copy of this LocalTime with the specified number of hours added.
// The calculation wraps around midnight. For example, 23:00 + 2 hours = 01:00.
// Negative values subtract hours. Returns zero value if called on zero value.
func (t LocalTime) PlusHours(hours int) LocalTime {
	if t.IsZero() {
		return t
	}
	return t.plusNanos(int64(hours) * int64(time.Hour))
}

// MinusHours returns a copy of this LocalTime with the specified number of hours subtracted.
// Equivalent to PlusHours(-hours).
func (t LocalTime) MinusHours(hours int) LocalTime {
	return t.PlusHours(-hours)
}

// PlusMinutes returns a copy of this LocalTime with the specified number of minutes added.
// The calculation wraps around midnight. For example, 23:50 + 20 minutes = 00:10.
// Negative values subtract minutes. Returns zero value if called on zero value.
func (t LocalTime) PlusMinutes(minutes int) LocalTime {
	if t.IsZero() {
		return t
	}
	return t.plusNanos(int64(minutes) * int64(time.Minute))
}

// MinusMinutes returns a copy of this LocalTime with the specified number of minutes subtracted.
// Equivalent to PlusMinutes(-minutes).
func (t LocalTime) MinusMinutes(minutes int) LocalTime {
	return t.PlusMinutes(-minutes)
}

// PlusSeconds returns a copy of this LocalTime with the specified number of seconds added.
// The calculation wraps around midnight. For example, 23:59:50 + 20 seconds = 00:00:10.
// Negative values subtract seconds. Returns zero value if called on zero value.
func (t LocalTime) PlusSeconds(seconds int) LocalTime {
	if t.IsZero() {
		return t
	}
	return t.plusNanos(int64(seconds) * int64(time.Second))
}

// MinusSeconds returns a copy of this LocalTime with the specified number of seconds subtracted.
// Equivalent to PlusSeconds(-seconds).
func (t LocalTime) MinusSeconds(seconds int) LocalTime {
	return t.PlusSeconds(-seconds)
}

// PlusNanos returns a copy of this LocalTime with the specified number of nanoseconds added.
// The calculation wraps around midnight.
// Negative values subtract nanoseconds. Returns zero value if called on zero value.
func (t LocalTime) PlusNanos(nanos int64) LocalTime {
	if t.IsZero() {
		return t
	}
	return t.plusNanos(nanos)
}

// MinusNanos returns a copy of this LocalTime with the specified number of nanoseconds subtracted.
// Equivalent to PlusNanos(-nanos).
func (t LocalTime) MinusNanos(nanos int64) LocalTime {
	return t.PlusNanos(-nanos)
}

// plusNanos is the internal implementation for adding nanoseconds.
// It handles wrapping around midnight (modulo operation with nanoseconds per day).
func (t LocalTime) plusNanos(nanosToAdd int64) LocalTime {
	const nanosPerDay = 24 * int64(time.Hour)
	currentNanos := t.v & localTimeValueMask

	// Add the nanoseconds
	newNanos := currentNanos + nanosToAdd

	// Wrap around midnight (handle both positive and negative overflow)
	newNanos = newNanos % nanosPerDay
	if newNanos < 0 {
		newNanos += nanosPerDay
	}

	return LocalTime{v: newNanos | localTimeValidBit}
}

func (t LocalTime) WithHour(hour int) (LocalTime, error) {
	return t.WithTemporal(FieldHourOfDay, TemporalValue{v: int64(hour)})
}

func (t LocalTime) WithMinute(minute int) (LocalTime, error) {
	return t.WithTemporal(FieldMinuteOfHour, TemporalValue{v: int64(minute)})
}

func (t LocalTime) WithSecond(second int) (LocalTime, error) {
	return t.WithTemporal(FieldSecondOfMinute, TemporalValue{v: int64(second)})
}

func (t LocalTime) WithNanosecond(nano int) (LocalTime, error) {
	return t.WithTemporal(FieldNanoOfSecond, TemporalValue{v: int64(nano)})
}

func (t LocalTime) MustWithHour(hour int) LocalTime {
	return mustValue(t.WithHour(hour))
}

func (t LocalTime) MustWithMinute(minute int) LocalTime {
	return mustValue(t.WithMinute(minute))
}

func (t LocalTime) MustWithSecond(second int) LocalTime {
	return mustValue(t.WithSecond(second))
}

func (t LocalTime) MustWithNanosecond(nano int) LocalTime {
	return mustValue(t.WithNanosecond(nano))
}

// WithTemporal returns a copy of this LocalTime with the specified field replaced.
// Zero values return zero immediately.
//
// Supported fields mirror Java's LocalTime#with(TemporalField, long):
//   - FieldNanoOfSecond: sets the nano-of-second while keeping hour, minute, and second.
//   - FieldNanoOfDay: replaces the entire time using the provided nano-of-day (equivalent to LocalTimeOfNanoOfDay).
//   - FieldMicroOfSecond: replaces the nano-of-second with micro-of-second × 1,000; hour, minute, and second stay the same.
//   - FieldMicroOfDay: replaces the entire time using micro-of-day × 1,000 (equivalent to LocalTimeOfNanoOfDay).
//   - FieldMilliOfSecond: replaces the nano-of-second with milli-of-second × 1,000,000; hour, minute, and second stay the same.
//   - FieldMilliOfDay: replaces the entire time using milli-of-day × 1,000,000 (equivalent to LocalTimeOfNanoOfDay).
//   - FieldSecondOfMinute: sets the second-of-minute while leaving hour, minute, and nano-of-second untouched.
//   - FieldSecondOfDay: replaces hour, minute, and second based on the provided second-of-day while keeping nano-of-second.
//   - FieldMinuteOfHour: sets the minute-of-hour; hour, second, and nano-of-second stay unchanged.
//   - FieldMinuteOfDay: replaces hour and minute based on the minute-of-day; second and nano-of-second stay unchanged.
//   - FieldHourOfAmPm: sets the hour within AM/PM while keeping the current half of day, minute, second, and nano-of-second.
//   - FieldClockHourOfAmPm: sets the 1-12 clock hour within AM/PM while keeping the current half of day, minute, second, and nano-of-second.
//   - FieldHourOfDay: sets the hour-of-day while leaving minute, second, and nano-of-second untouched.
//   - FieldClockHourOfDay: sets the 1-24 clock hour-of-day (24 → 0) while leaving minute, second, and nano-of-second untouched.
//   - FieldAmPmOfDay: toggles AM/PM while preserving the hour-of-am-pm, minute, second, and nano-of-second.
//
// Fields outside this list return an error. Range violations propagate the validation error.
func (t LocalTime) WithTemporal(field Field, value TemporalValue) (r LocalTime, e error) {
	if t.IsZero() {
		return t, nil
	}

	h := t.GetField(FieldHourOfDay).Int64()
	m := t.GetField(FieldMinuteOfHour).Int64()
	s := t.GetField(FieldSecondOfMinute).Int64()
	n := t.GetField(FieldNanoOfSecond).Int64()
	v := value.Int64()

	switch field {
	case FieldNanoOfSecond:
		e = checkTemporalInRange(FieldNanoOfSecond, 0, 999_999_999, value, e)
		n = v
	case FieldNanoOfDay:
		e = checkTemporalInRange(FieldNanoOfDay, 0, 86_399_999_999_999, value, e)
		return LocalTimeOfNanoOfDay(v)
	case FieldMicroOfSecond:
		e = checkTemporalInRange(FieldMicroOfSecond, 0, 999_999, value, e)
		n = v * 1_000
	case FieldMicroOfDay:
		e = checkTemporalInRange(FieldMicroOfDay, 0, 86_399_999_999, value, e)
		return LocalTimeOfNanoOfDay(v * 1_000)
	case FieldMilliOfSecond:
		e = checkTemporalInRange(FieldMilliOfSecond, 0, 999, value, e)
		n = v * 1_000_000
	case FieldMilliOfDay:
		e = checkTemporalInRange(FieldMilliOfDay, 0, 86_399_999, value, e)
		return LocalTimeOfNanoOfDay(v * 1_000_000)
	case FieldSecondOfMinute:
		e = checkTemporalInRange(FieldSecondOfMinute, 0, 59, value, e)
		s = v
	case FieldSecondOfDay:
		e = checkTemporalInRange(FieldSecondOfDay, 0, 86_399, value, e)
		h = v / 3600
		m = (v / 60) % 60
		s = v % 60
		// nano unchanged
	case FieldMinuteOfHour:
		e = checkTemporalInRange(FieldMinuteOfHour, 0, 59, value, e)
		m = v
	case FieldMinuteOfDay:
		e = checkTemporalInRange(FieldMinuteOfDay, 0, 1_439, value, e)
		h = v / 60
		m = v % 60
	case FieldHourOfAmPm:
		e = checkTemporalInRange(FieldHourOfAmPm, 0, 11, value, e)
		if h >= 12 { // PM
			h = v + 12
		} else {
			h = v
		}
	case FieldClockHourOfAmPm:
		e = checkTemporalInRange(FieldClockHourOfAmPm, 1, 12, value, e)
		tmp := v % 12 // 12 → 0
		if h >= 12 {  // PM
			h = tmp + 12
		} else { // AM
			h = tmp
		}
	case FieldHourOfDay:
		e = checkTemporalInRange(FieldHourOfDay, 0, 23, value, e)
		h = v
	case FieldClockHourOfDay:
		e = checkTemporalInRange(FieldClockHourOfDay, 1, 24, value, e)
		if v == 24 {
			h = 0
		} else {
			h = v
		}
	case FieldAmPmOfDay:
		e = checkTemporalInRange(FieldAmPmOfDay, 0, 1, value, e)
		if v == 0 && h >= 12 { // to AM
			h -= 12
		}
		if v == 1 && h < 12 { // to PM
			h += 12
		}
	default:
		return t, newError("unsupported field: %v", field)
	}
	if e != nil {
		return
	}
	return localTimeOf(h, m, s, n)
}

// LocalTimeOf creates a new LocalTime from the specified hour, minute, second, and nanosecond.
// Returns an error if any component is out of range:
// - hour must be 0-23
// - minute must be 0-59
// - second must be 0-59
// - nanosecond must be 0-999999999
func LocalTimeOf(hour, minute, second, nanosecond int) (LocalTime, error) {
	return localTimeOf(int64(hour), int64(minute), int64(second), int64(nanosecond))
}

func localTimeOf(hour, minute, second, nanosecond int64) (r LocalTime, e error) {
	if hour < 0 || hour >= 24 {
		e = newError("hour %d out of range", hour)
	}
	if minute < 0 || minute >= 60 {
		e = newError("minute %d out of range", minute)
	}
	if second < 0 || second >= 60 {
		e = newError("second %d out of range", second)
	}
	if nanosecond < 0 || nanosecond >= 1000000000 {
		e = newError("nanosecond %d out of range", nanosecond)
	}
	if e != nil {
		return
	}
	nanos := hour*int64(time.Hour) +
		minute*int64(time.Minute) +
		second*int64(time.Second) +
		nanosecond
	return LocalTime{
		v: nanos | localTimeValidBit,
	}, nil
}

// MustLocalTimeOf creates a new LocalTime from the specified hour, minute, second, and nanosecond.
// Panics if any component is out of range. Use LocalTimeOf for error handling.
func MustLocalTimeOf(hour, minute, second, nanosecond int) LocalTime {
	return mustValue(LocalTimeOf(hour, minute, second, nanosecond))
}

// LocalTimeOfGoTime creates a LocalTime from a time.Time.
// The date and time zone components are ignored.
// Returns zero value if t.IsZero().
func LocalTimeOfGoTime(t time.Time) LocalTime {
	if t.IsZero() {
		return LocalTime{}
	}
	nanos := time.Date(1970, 1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC).UnixNano()
	return LocalTime{
		v: nanos | localTimeValidBit,
	}
}

// LocalTimeOfNanoOfDay creates a LocalTime from the nanosecond-of-day value.
// Returns an error if nanoOfDay is out of the valid range (0 to 86,399,999,999,999).
// The valid range represents 00:00:00.000000000 to 23:59:59.999999999.
//
// Example:
//
//	// Create time for 12:00:00
//	lt, err := LocalTimeOfNanoOfDay(12 * 60 * 60 * 1000000000)
func LocalTimeOfNanoOfDay(nanoOfDay int64) (LocalTime, error) {
	const nanosPerDay = 24 * int64(time.Hour)
	if nanoOfDay < 0 || nanoOfDay >= nanosPerDay {
		return LocalTime{}, newError("nanoOfDay %d out of range [0, %d)", nanoOfDay, nanosPerDay)
	}
	return LocalTime{
		v: nanoOfDay | localTimeValidBit,
	}, nil
}

// MustLocalTimeOfNanoOfDay creates a LocalTime from the nanosecond-of-day value.
// Panics if nanoOfDay is out of range. Use LocalTimeOfNanoOfDay for error handling.
//
// Example:
//
//	// Create time for 12:00:00
//	lt := MustLocalTimeOfNanoOfDay(12 * 60 * 60 * 1000000000)
func MustLocalTimeOfNanoOfDay(nanoOfDay int64) LocalTime {
	return mustValue(LocalTimeOfNanoOfDay(nanoOfDay))
}

// LocalTimeOfSecondOfDay creates a LocalTime from the second-of-day value.
// Returns an error if secondOfDay is out of the valid range (0 to 86,399).
// The valid range represents 00:00:00 to 23:59:59.
// The nanosecond component will be set to zero.
//
// Example:
//
//	// Create time for 12:00:00
//	lt, err := LocalTimeOfSecondOfDay(12 * 60 * 60)
func LocalTimeOfSecondOfDay(secondOfDay int) (LocalTime, error) {
	const secondsPerDay = 24 * 60 * 60
	if secondOfDay < 0 || secondOfDay >= secondsPerDay {
		return LocalTime{}, newError("secondOfDay %d out of range [0, %d)", secondOfDay, secondsPerDay)
	}
	nanos := int64(secondOfDay) * int64(time.Second)
	return LocalTime{
		v: nanos | localTimeValidBit,
	}, nil
}

// MustLocalTimeOfSecondOfDay creates a LocalTime from the second-of-day value.
// Panics if secondOfDay is out of range. Use LocalTimeOfSecondOfDay for error handling.
//
// Example:
//
//	// Create time for 12:00:00
//	lt := MustLocalTimeOfSecondOfDay(12 * 60 * 60)
func MustLocalTimeOfSecondOfDay(secondOfDay int) LocalTime {
	return mustValue(LocalTimeOfSecondOfDay(secondOfDay))
}

// LocalTimeNow returns the current time in the system's local time zone.
// This is equivalent to LocalTimeOfGoTime(time.Now()).
// For UTC time, use LocalTimeNowUTC. For a specific timezone, use LocalTimeNowIn.
func LocalTimeNow() LocalTime {
	return LocalTimeOfGoTime(time.Now())
}

// LocalTimeNowIn returns the current time in the specified time zone.
// This is equivalent to LocalTimeOfGoTime(time.Now().In(loc)).
func LocalTimeNowIn(loc *time.Location) LocalTime {
	return LocalTimeOfGoTime(time.Now().In(loc))
}

// LocalTimeNowUTC returns the current time in UTC.
// This is equivalent to LocalTimeOfGoTime(time.Now().UTC()).
func LocalTimeNowUTC() LocalTime {
	return LocalTimeOfGoTime(time.Now().UTC())
}

// LocalTimeParse parses a time string in HH:mm:ss[.nnnnnnnnn] format (24-hour).
// Returns an error if the string is invalid or represents an invalid time.
//
// Accepts fractional seconds of any length (1-9 digits):
//   - HH:mm:ss (e.g., "14:30:45")
//   - HH:mm:ss.f (e.g., "14:30:45.1" → 100 milliseconds)
//   - HH:mm:ss.fff (e.g., "14:30:45.123" → 123 milliseconds)
//   - HH:mm:ss.ffffff (e.g., "14:30:45.123456" → 123.456 milliseconds)
//   - HH:mm:ss.nnnnnnnnn (e.g., "14:30:45.123456789" → full nanosecond precision)
//
// Note: Output formatting aligns to 3-digit boundaries (milliseconds, microseconds, nanoseconds).
// Timezone offsets are not supported.
//
// Example:
//
//	time, err := LocalTimeParse("14:30:45.123")
//	if err != nil {
//	    // handle error
//	}
func LocalTimeParse(s string) (LocalTime, error) {
	var t LocalTime
	err := t.UnmarshalText([]byte(s))
	return t, err
}

// MustLocalTimeParse parses a time string in HH:mm:ss[.nnnnnnnnn] format (24-hour).
// Panics if the string is invalid. Use LocalTimeParse for error handling.
//
// Example:
//
//	time := MustLocalTimeParse("14:30:45.123456789")
func MustLocalTimeParse(s string) LocalTime {
	return mustValue(LocalTimeParse(s))
}

var (
	_ encoding.TextAppender    = (*LocalTime)(nil)
	_ fmt.Stringer             = (*LocalTime)(nil)
	_ encoding.TextMarshaler   = (*LocalTime)(nil)
	_ encoding.TextUnmarshaler = (*LocalTime)(nil)
	_ json.Marshaler           = (*LocalTime)(nil)
	_ json.Unmarshaler         = (*LocalTime)(nil)
	_ driver.Valuer            = (*LocalTime)(nil)
	_ sql.Scanner              = (*LocalTime)(nil)
)

// Compile-time check that LocalTime is comparable
func _assertLocalTimeIsComparable[T comparable](t T) {}

var _ = _assertLocalTimeIsComparable[LocalTime]
