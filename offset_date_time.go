package goda

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"time"
)

// OffsetDateTime represents a date-time with a time-zone offset from UTC,
// such as 2024-03-15T14:30:45.123456789+01:00.
//
// OffsetDateTime stores an offset from UTC, but does not store or use a time-zone.
// Use ZonedDateTime when you need full time-zone support including DST transitions.
//
// OffsetDateTime is comparable and can be used as a map key.
// The zero value represents an unset date-time and IsZero returns true for it.
//
// OffsetDateTime implements sql.Scanner and driver.Valuer for database operations,
// encoding.TextMarshaler and encoding.TextUnmarshaler for text serialization,
// and json.Marshaler and json.Unmarshaler for JSON serialization.
//
// Format: yyyy-MM-ddTHH:mm:ss[.nnnnnnnnn]±HH:mm (e.g., "2024-03-15T14:30:45.123456789+01:00").
type OffsetDateTime struct {
	datetime LocalDateTime
	offset   ZoneOffset
}

// LocalDateTime returns the local date-time part without the offset.
func (odt OffsetDateTime) LocalDateTime() LocalDateTime {
	return odt.datetime
}

// LocalDate returns the date part of this date-time.
func (odt OffsetDateTime) LocalDate() LocalDate {
	return odt.datetime.LocalDate()
}

// LocalTime returns the time part of this date-time.
func (odt OffsetDateTime) LocalTime() LocalTime {
	return odt.datetime.LocalTime()
}

// Offset returns the zone offset.
func (odt OffsetDateTime) Offset() ZoneOffset {
	return odt.offset
}

// WithOffsetSameLocal returns a copy with the specified offset.
// The local date-time is not changed, only the offset is changed.
func (odt OffsetDateTime) WithOffsetSameLocal(offset ZoneOffset) OffsetDateTime {
	if odt.IsZero() {
		return OffsetDateTime{}
	}
	return OffsetDateTime{
		datetime: odt.datetime,
		offset:   offset,
	}
}

// WithOffsetSameInstant returns a copy with the specified offset.
// The instant represented by this date-time is preserved.
// The local date-time is adjusted to maintain the same instant.
func (odt OffsetDateTime) WithOffsetSameInstant(offset ZoneOffset) OffsetDateTime {
	if odt.IsZero() {
		return OffsetDateTime{}
	}
	diff := offset.TotalSeconds() - odt.offset.TotalSeconds()
	newOdt := odt.PlusSeconds(int64(diff))
	// Update the offset to the new one
	newOdt.offset = offset
	return newOdt
}

// Year returns the year component.
func (odt OffsetDateTime) Year() Year {
	return odt.datetime.Year()
}

// Month returns the month component.
func (odt OffsetDateTime) Month() Month {
	return odt.datetime.Month()
}

// DayOfMonth returns the day-of-month component.
func (odt OffsetDateTime) DayOfMonth() int {
	return odt.datetime.DayOfMonth()
}

// DayOfWeek returns the day-of-week.
func (odt OffsetDateTime) DayOfWeek() DayOfWeek {
	return odt.datetime.DayOfWeek()
}

// DayOfYear returns the day-of-year.
func (odt OffsetDateTime) DayOfYear() int {
	return odt.datetime.DayOfYear()
}

// Hour returns the hour component (0-23).
func (odt OffsetDateTime) Hour() int {
	return odt.datetime.Hour()
}

// Minute returns the minute component (0-59).
func (odt OffsetDateTime) Minute() int {
	return odt.datetime.Minute()
}

// Second returns the second component (0-59).
func (odt OffsetDateTime) Second() int {
	return odt.datetime.Second()
}

// Millisecond returns the millisecond component (0-999).
func (odt OffsetDateTime) Millisecond() int {
	return odt.datetime.Millisecond()
}

// Nanosecond returns the nanosecond component (0-999999999).
func (odt OffsetDateTime) Nanosecond() int {
	return odt.datetime.Nanosecond()
}

// IsZero returns true if this is the zero value.
func (odt OffsetDateTime) IsZero() bool {
	return odt.datetime.IsZero() && odt.offset.IsZero()
}

// IsLeapYear returns true if the year is a leap year.
func (odt OffsetDateTime) IsLeapYear() bool {
	return odt.datetime.IsLeapYear()
}

// IsSupportedField returns true if the field is supported by OffsetDateTime.
// OffsetDateTime supports all fields from LocalDateTime plus FieldOffsetSeconds and FieldInstantSeconds.
func (odt OffsetDateTime) IsSupportedField(field Field) bool {
	return odt.datetime.IsSupportedField(field) || odt.offset.IsSupportedField(field) || field == FieldInstantSeconds
}

// GetField returns the value of the specified field as a TemporalValue.
// This method queries the date-time for the value of the specified field.
// The returned value may be unsupported if the field is not supported by OffsetDateTime.
//
// If the date-time is zero (IsZero() returns true), an unsupported TemporalValue is returned.
//
// OffsetDateTime supports all fields from LocalDateTime plus:
//   - FieldOffsetSeconds: the offset in seconds
//   - FieldInstantSeconds: the epoch seconds (Unix timestamp)
func (odt OffsetDateTime) GetField(field Field) TemporalValue {
	if odt.IsZero() {
		return TemporalValue{v: 0, unsupported: true}
	}
	if odt.offset.IsSupportedField(field) {
		return odt.offset.GetField(field)
	}

	// Handle offset-specific fields
	switch field {
	case FieldInstantSeconds:
		return TemporalValue{v: odt.ToEpochSecond() - int64(odt.offset.TotalSeconds())}
	}

	// Delegate to LocalDateTime for date and time fields
	return odt.datetime.GetField(field)
}

// GoTime converts this offset date-time to a time.Time.
// Returns time.Time{} (zero) for zero value.
func (odt OffsetDateTime) GoTime() time.Time {
	if odt.IsZero() {
		return time.Time{}
	}
	// Create a fixed location with the offset
	loc := time.FixedZone("", odt.offset.TotalSeconds())
	return time.Date(int(odt.Year()), time.Month(odt.Month()), odt.DayOfMonth(), odt.Hour(), odt.Minute(), odt.Second(), odt.Nanosecond(), loc)
}

// ToEpochSecond returns the number of seconds since Unix epoch (1970-01-01T00:00:00Z).
func (odt OffsetDateTime) ToEpochSecond() int64 {
	if odt.IsZero() {
		return 0
	}
	epochDay := odt.datetime.LocalDate().UnixEpochDays()
	secondsOfDay := odt.datetime.LocalTime().GetField(FieldSecondOfDay).Int64()
	return epochDay*86400 + secondsOfDay - int64(odt.offset.TotalSeconds())
}

// Compare compares this offset date-time with another.
// The comparison is based on the instant then on the local date-time.
// Returns -1 if this is before other, 0 if equal, 1 if after.
func (odt OffsetDateTime) Compare(other OffsetDateTime) int {
	return doCompare(odt, other, comparing(OffsetDateTime.ToEpochSecond), comparing(OffsetDateTime.Nanosecond), comparing1(OffsetDateTime.LocalDateTime))
}

// IsBefore returns true if this offset date-time is before the specified offset date-time.
func (odt OffsetDateTime) IsBefore(other OffsetDateTime) bool {
	return odt.Compare(other) < 0
}

// IsAfter returns true if this offset date-time is after the specified offset date-time.
func (odt OffsetDateTime) IsAfter(other OffsetDateTime) bool {
	return odt.Compare(other) > 0
}

// PlusYears returns a copy with the specified number of years added.
func (odt OffsetDateTime) PlusYears(years int) OffsetDateTime {
	if odt.IsZero() {
		return OffsetDateTime{}
	}
	return OffsetDateTime{
		datetime: odt.datetime.PlusYears(years),
		offset:   odt.offset,
	}
}

// MinusYears returns a copy with the specified number of years subtracted.
func (odt OffsetDateTime) MinusYears(years int) OffsetDateTime {
	return odt.PlusYears(-years)
}

// PlusMonths returns a copy with the specified number of months added.
func (odt OffsetDateTime) PlusMonths(months int) OffsetDateTime {
	if odt.IsZero() {
		return OffsetDateTime{}
	}
	return OffsetDateTime{
		datetime: odt.datetime.PlusMonths(months),
		offset:   odt.offset,
	}
}

// MinusMonths returns a copy with the specified number of months subtracted.
func (odt OffsetDateTime) MinusMonths(months int) OffsetDateTime {
	return odt.PlusMonths(-months)
}

// PlusDays returns a copy with the specified number of days added.
func (odt OffsetDateTime) PlusDays(days int) OffsetDateTime {
	if odt.IsZero() {
		return OffsetDateTime{}
	}
	return OffsetDateTime{
		datetime: odt.datetime.PlusDays(days),
		offset:   odt.offset,
	}
}

// MinusDays returns a copy with the specified number of days subtracted.
func (odt OffsetDateTime) MinusDays(days int) OffsetDateTime {
	return odt.PlusDays(-days)
}

// PlusHours returns a copy with the specified number of hours added.
func (odt OffsetDateTime) PlusHours(hours int64) OffsetDateTime {
	return odt.PlusSeconds(hours * 3600)
}

// MinusHours returns a copy with the specified number of hours subtracted.
func (odt OffsetDateTime) MinusHours(hours int64) OffsetDateTime {
	return odt.PlusHours(-hours)
}

// PlusMinutes returns a copy with the specified number of minutes added.
func (odt OffsetDateTime) PlusMinutes(minutes int64) OffsetDateTime {
	return odt.PlusSeconds(minutes * 60)
}

// MinusMinutes returns a copy with the specified number of minutes subtracted.
func (odt OffsetDateTime) MinusMinutes(minutes int64) OffsetDateTime {
	return odt.PlusMinutes(-minutes)
}

// PlusSeconds returns a copy with the specified number of seconds added.
func (odt OffsetDateTime) PlusSeconds(seconds int64) OffsetDateTime {
	if odt.IsZero() || seconds == 0 {
		return odt
	}
	return odt.PlusNanos(seconds * 1_000_000_000)
}

// MinusSeconds returns a copy with the specified number of seconds subtracted.
func (odt OffsetDateTime) MinusSeconds(seconds int64) OffsetDateTime {
	return odt.PlusSeconds(-seconds)
}

// PlusNanos returns a copy with the specified number of nanoseconds added.
func (odt OffsetDateTime) PlusNanos(nanos int64) OffsetDateTime {
	if odt.IsZero() || nanos == 0 {
		return odt
	}

	// Convert current date-time to total nanoseconds since epoch day
	secondsOfDay := int64(odt.datetime.LocalTime().GetField(FieldSecondOfDay).Int64())
	nanosOfDay := secondsOfDay*1_000_000_000 + int64(odt.Nanosecond())

	// Add the nanos
	totalNanos := nanosOfDay + nanos

	// Calculate days overflow
	daysToAdd := totalNanos / 86400_000_000_000
	newNanosOfDay := totalNanos % 86400_000_000_000

	// Handle negative overflow
	if newNanosOfDay < 0 {
		newNanosOfDay += 86400_000_000_000
		daysToAdd--
	}

	// Create new date and time
	newDate := odt.datetime.LocalDate().PlusDays(int(daysToAdd))
	newTime := MustLocalTimeOfNanoOfDay(newNanosOfDay)

	return OffsetDateTime{
		datetime: newDate.AtTime(newTime),
		offset:   odt.offset,
	}
}

// MinusNanos returns a copy with the specified number of nanoseconds subtracted.
func (odt OffsetDateTime) MinusNanos(nanos int64) OffsetDateTime {
	return odt.PlusNanos(-nanos)
}

// String returns the ISO 8601 string representation (yyyy-MM-ddTHH:mm:ss[.nnnnnnnnn]±HH:mm).
func (odt OffsetDateTime) String() string {
	return stringImpl(odt)
}

// AppendText implements encoding.TextAppender.
func (odt OffsetDateTime) AppendText(b []byte) ([]byte, error) {
	if odt.IsZero() {
		return b, nil
	}
	b, _ = odt.datetime.AppendText(b)
	b, _ = odt.offset.AppendText(b)
	return b, nil
}

// MarshalText implements encoding.TextMarshaler.
func (odt OffsetDateTime) MarshalText() ([]byte, error) {
	return marshalTextImpl(odt)
}

// UnmarshalText implements encoding.TextUnmarshaler.
// Accepts ISO 8601 format: yyyy-MM-ddTHH:mm:ss[.nnnnnnnnn]±HH:mm[:ss] or Z for UTC.
func (odt *OffsetDateTime) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*odt = OffsetDateTime{}
		return nil
	}

	// Find the offset part (starts with +, -, or Z)
	offsetIdx := -1
	for i := len(text) - 1; i >= 0; i-- {
		ch := text[i]
		if ch == '+' || ch == '-' {
			offsetIdx = i
			break
		}
		if ch == 'Z' || ch == 'z' {
			offsetIdx = i
			break
		}
	}

	if offsetIdx < 0 {
		return newError("invalid offset date-time format: missing offset")
	}

	// Parse date-time part
	var dt LocalDateTime
	if err := dt.UnmarshalText(text[:offsetIdx]); err != nil {
		return err
	}

	// Parse offset part
	var offset ZoneOffset
	if err := offset.UnmarshalText(text[offsetIdx:]); err != nil {
		return err
	}

	*odt = OffsetDateTime{
		datetime: dt,
		offset:   offset,
	}
	return nil
}

// MarshalJSON implements json.Marshaler.
func (odt OffsetDateTime) MarshalJSON() ([]byte, error) {
	return marshalJsonImpl(odt)
}

// UnmarshalJSON implements json.Unmarshaler.
func (odt *OffsetDateTime) UnmarshalJSON(data []byte) error {
	if len(data) == 4 && string(data) == "null" {
		*odt = OffsetDateTime{}
		return nil
	}
	return unmarshalJsonImpl(odt, data)
}

// Scan implements sql.Scanner.
func (odt *OffsetDateTime) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*odt = OffsetDateTime{}
		return nil
	case string:
		return odt.UnmarshalText([]byte(v))
	case []byte:
		return odt.UnmarshalText(v)
	case time.Time:
		*odt = OffsetDateTimeOfGoTime(v)
		return nil
	default:
		return sqlScannerDefaultBranch(v)
	}
}

// Value implements driver.Valuer.
func (odt OffsetDateTime) Value() (driver.Value, error) {
	if odt.IsZero() {
		return nil, nil
	}
	return odt.String(), nil
}

// OffsetDateTimeOf creates a new OffsetDateTime from individual components.
// Returns an error if any component is invalid.
func OffsetDateTimeOf(year Year, month Month, day int, hour int, minute int, second int, nanosecond int, offset ZoneOffset) (OffsetDateTime, error) {
	dt, err := LocalDateTimeOf(year, month, day, hour, minute, second, nanosecond)
	if err != nil {
		return OffsetDateTime{}, err
	}
	return OffsetDateTime{
		datetime: dt,
		offset:   offset,
	}, nil
}

// MustOffsetDateTimeOf creates a new OffsetDateTime from individual components.
// Panics if any component is invalid.
func MustOffsetDateTimeOf(year Year, month Month, day int, hour int, minute int, second int, nanosecond int, offset ZoneOffset) OffsetDateTime {
	return mustValue(OffsetDateTimeOf(year, month, day, hour, minute, second, nanosecond, offset))
}

// OffsetDateTimeNow returns the current date-time with offset in the system's local time zone.
func OffsetDateTimeNow() OffsetDateTime {
	return OffsetDateTimeOfGoTime(time.Now())
}

// OffsetDateTimeNowUTC returns the current date-time with offset in UTC.
func OffsetDateTimeNowUTC() OffsetDateTime {
	return OffsetDateTimeOfGoTime(time.Now().UTC())
}

// OffsetDateTimeOfGoTime creates an OffsetDateTime from a time.Time.
// Returns zero value if t.IsZero().
func OffsetDateTimeOfGoTime(t time.Time) OffsetDateTime {
	if t.IsZero() {
		return OffsetDateTime{}
	}
	_, offsetSeconds := t.Zone()
	offset := MustZoneOffsetOfSeconds(offsetSeconds)
	return LocalDateTimeOfGoTime(t).AtOffset(offset)
}

// OffsetDateTimeParse parses a date-time string in RFC3339-compatible format.
// The format is yyyy-MM-ddTHH:mm:ss[.nnnnnnnnn]±HH:mm[:ss] or with 'Z' for UTC.
//
// Examples:
//
//	odt, err := OffsetDateTimeParse("2024-03-15T14:30:45.123456789+01:00")
//	odt, err := OffsetDateTimeParse("2024-03-15T14:30:45Z")
//	odt, err := OffsetDateTimeParse("2024-03-15T14:30:45+05:30")
func OffsetDateTimeParse(s string) (OffsetDateTime, error) {
	var odt OffsetDateTime
	err := odt.UnmarshalText([]byte(s))
	return odt, err
}

// MustOffsetDateTimeParse parses a date-time string with offset.
// Panics if the string is invalid.
func MustOffsetDateTimeParse(s string) OffsetDateTime {
	odt, err := OffsetDateTimeParse(s)
	if err != nil {
		panic(err)
	}
	return odt
}

// Compile-time interface checks
var (
	_ encoding.TextAppender    = (*OffsetDateTime)(nil)
	_ fmt.Stringer             = (*OffsetDateTime)(nil)
	_ encoding.TextMarshaler   = (*OffsetDateTime)(nil)
	_ encoding.TextUnmarshaler = (*OffsetDateTime)(nil)
	_ json.Marshaler           = (*OffsetDateTime)(nil)
	_ json.Unmarshaler         = (*OffsetDateTime)(nil)
	_ driver.Valuer            = (*OffsetDateTime)(nil)
	_ sql.Scanner              = (*OffsetDateTime)(nil)
)

// Compile-time check that OffsetDateTime is comparable
func _assertOffsetDateTimeIsComparable[T comparable](t T) {}

var _ = _assertOffsetDateTimeIsComparable[OffsetDateTime]
