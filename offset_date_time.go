package goda

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"time"
)

// OffsetDateTime represents a date-time with an offset from UTC/Greenwich,
// such as 2024-03-15T14:30:45.123456789+09:00.
// It combines LocalDateTime with a ZoneOffset.
//
// OffsetDateTime is comparable and can be used as a map key.
// The zero value represents an unset date-time and IsZero returns true for it.
//
// OffsetDateTime implements sql.Scanner and driver.Valuer for database operations,
// encoding.TextMarshaler and encoding.TextUnmarshaler for text serialization,
// and json.Marshaler and json.Unmarshaler for JSON serialization.
//
// Format: yyyy-MM-ddTHH:mm:ss[.nnnnnnnnn]±HH:mm (e.g., "2024-03-15T14:30:45.123456789+09:00").
// The offset is always included and represents the difference from UTC.
type OffsetDateTime struct {
	dateTime LocalDateTime
	offset   ZoneOffset
}

// NewOffsetDateTime creates a new OffsetDateTime from a LocalDateTime and ZoneOffset.
func NewOffsetDateTime(dateTime LocalDateTime, offset ZoneOffset) OffsetDateTime {
	return OffsetDateTime{
		dateTime: dateTime,
		offset:   offset,
	}
}

// NewOffsetDateTimeFromComponents creates a new OffsetDateTime from individual components.
// Returns an error if any component is invalid.
func NewOffsetDateTimeFromComponents(year Year, month Month, day, hour, minute, second, nanosecond int, offset ZoneOffset) (OffsetDateTime, error) {
	dt, err := NewLocalDateTimeFromComponents(year, month, day, hour, minute, second, nanosecond)
	if err != nil {
		return OffsetDateTime{}, err
	}
	return NewOffsetDateTime(dt, offset), nil
}

// MustNewOffsetDateTimeFromComponents creates a new OffsetDateTime from individual components.
// Panics if any component is invalid.
func MustNewOffsetDateTimeFromComponents(year Year, month Month, day, hour, minute, second, nanosecond int, offset ZoneOffset) OffsetDateTime {
	odt, err := NewOffsetDateTimeFromComponents(year, month, day, hour, minute, second, nanosecond, offset)
	if err != nil {
		panic(err)
	}
	return odt
}

// OffsetDateTimeNow returns the current date-time with the system's offset.
func OffsetDateTimeNow() OffsetDateTime {
	return NewOffsetDateTimeByGoTime(time.Now())
}

// OffsetDateTimeNowIn returns the current date-time in the specified time zone.
func OffsetDateTimeNowIn(loc *time.Location) OffsetDateTime {
	return NewOffsetDateTimeByGoTime(time.Now().In(loc))
}

// OffsetDateTimeNowUTC returns the current date-time in UTC.
func OffsetDateTimeNowUTC() OffsetDateTime {
	return NewOffsetDateTimeByGoTime(time.Now().UTC())
}

// NewOffsetDateTimeByGoTime creates an OffsetDateTime from a time.Time.
// Returns zero value if t.IsZero().
func NewOffsetDateTimeByGoTime(t time.Time) OffsetDateTime {
	if t.IsZero() {
		return OffsetDateTime{}
	}
	_, offset := t.Zone()
	return OffsetDateTime{
		dateTime: NewLocalDateTimeByGoTime(t),
		offset:   MustNewZoneOffsetSeconds(offset),
	}
}

// ParseOffsetDateTime parses a date-time string with offset in ISO 8601 format.
// The date must be in yyyy-MM-dd form, time in HH:mm:ss[.nnnnnnnnn] form,
// and offset in ±HH:mm or Z form.
//
// Examples:
//
//	odt, err := ParseOffsetDateTime("2024-03-15T14:30:45.123456789+09:00")
//	odt, err := ParseOffsetDateTime("2024-03-15T14:30:45Z")
//	odt, err := ParseOffsetDateTime("2024-03-15T14:30:45+09:00")
func ParseOffsetDateTime(s string) (OffsetDateTime, error) {
	var odt OffsetDateTime
	err := odt.UnmarshalText([]byte(s))
	return odt, err
}

// MustParseOffsetDateTime parses a date-time string with offset in ISO 8601 format.
// Panics if the string is invalid.
func MustParseOffsetDateTime(s string) OffsetDateTime {
	odt, err := ParseOffsetDateTime(s)
	if err != nil {
		panic(err)
	}
	return odt
}

// LocalDateTime returns the local date-time part.
func (odt OffsetDateTime) LocalDateTime() LocalDateTime {
	return odt.dateTime
}

// Offset returns the zone offset.
func (odt OffsetDateTime) Offset() ZoneOffset {
	return odt.offset
}

// LocalDate returns the date part.
func (odt OffsetDateTime) LocalDate() LocalDate {
	return odt.dateTime.LocalDate()
}

// LocalTime returns the time part.
func (odt OffsetDateTime) LocalTime() LocalTime {
	return odt.dateTime.LocalTime()
}

// Year returns the year component.
func (odt OffsetDateTime) Year() Year {
	return odt.dateTime.Year()
}

// Month returns the month component.
func (odt OffsetDateTime) Month() Month {
	return odt.dateTime.Month()
}

// DayOfMonth returns the day-of-month component.
func (odt OffsetDateTime) DayOfMonth() int {
	return odt.dateTime.DayOfMonth()
}

// DayOfWeek returns the day-of-week.
func (odt OffsetDateTime) DayOfWeek() DayOfWeek {
	return odt.dateTime.DayOfWeek()
}

// DayOfYear returns the day-of-year.
func (odt OffsetDateTime) DayOfYear() int {
	return odt.dateTime.DayOfYear()
}

// Hour returns the hour component (0-23).
func (odt OffsetDateTime) Hour() int {
	return odt.dateTime.Hour()
}

// Minute returns the minute component (0-59).
func (odt OffsetDateTime) Minute() int {
	return odt.dateTime.Minute()
}

// Second returns the second component (0-59).
func (odt OffsetDateTime) Second() int {
	return odt.dateTime.Second()
}

// Millisecond returns the millisecond component (0-999).
func (odt OffsetDateTime) Millisecond() int {
	return odt.dateTime.Millisecond()
}

// Nanosecond returns the nanosecond component (0-999999999).
func (odt OffsetDateTime) Nanosecond() int {
	return odt.dateTime.Nanosecond()
}

// IsZero returns true if this is the zero value.
func (odt OffsetDateTime) IsZero() bool {
	return odt.dateTime.IsZero()
}

// IsLeapYear returns true if the year is a leap year.
func (odt OffsetDateTime) IsLeapYear() bool {
	return odt.dateTime.IsLeapYear()
}

// IsSupportedField returns true if the field is supported by OffsetDateTime.
// OffsetDateTime supports all fields from LocalDateTime plus OffsetSeconds and InstantSeconds.
func (odt OffsetDateTime) IsSupportedField(field Field) bool {
	return odt.dateTime.IsSupportedField(field) || field == OffsetSeconds || field == InstantSeconds
}

// GetFieldInt64 returns the value of the specified field as an int64.
// Returns 0 if the field is not supported or the datetime is zero.
func (odt OffsetDateTime) GetFieldInt64(field Field) int64 {
	if odt.IsZero() {
		return 0
	}
	switch field {
	case OffsetSeconds:
		return int64(odt.offset.TotalSeconds())
	case InstantSeconds:
		return odt.ToInstant().Unix()
	default:
		return odt.dateTime.GetFieldInt64(field)
	}
}

// ToInstant converts this OffsetDateTime to a time.Time (instant).
// Returns time.Time{} (zero) for zero value.
func (odt OffsetDateTime) ToInstant() time.Time {
	if odt.IsZero() {
		return time.Time{}
	}
	// Create a time.Time with the local date-time and offset
	t := time.Date(
		int(odt.Year()),
		time.Month(odt.Month()),
		odt.DayOfMonth(),
		odt.Hour(),
		odt.Minute(),
		odt.Second(),
		odt.Nanosecond(),
		time.UTC,
	)
	// Adjust for offset
	return t.Add(-time.Duration(odt.offset.TotalSeconds()) * time.Second)
}

// GoTime is an alias for ToInstant for consistency with other types.
func (odt OffsetDateTime) GoTime() time.Time {
	return odt.ToInstant()
}

// WithOffsetSameLocal returns a copy with a different offset, keeping the local date-time the same.
// This changes the instant represented by this OffsetDateTime.
func (odt OffsetDateTime) WithOffsetSameLocal(offset ZoneOffset) OffsetDateTime {
	return NewOffsetDateTime(odt.dateTime, offset)
}

// WithOffsetSameInstant returns a copy with a different offset, adjusting the local date-time
// to represent the same instant.
func (odt OffsetDateTime) WithOffsetSameInstant(offset ZoneOffset) OffsetDateTime {
	if odt.IsZero() {
		return odt
	}
	// Adjust the local time
	instant := odt.ToInstant()
	newTime := instant.Add(time.Duration(offset.TotalSeconds()) * time.Second)
	return OffsetDateTime{
		dateTime: NewLocalDateTimeByGoTime(newTime),
		offset:   offset,
	}
}

// ToUTC converts this OffsetDateTime to UTC, adjusting the local date-time.
func (odt OffsetDateTime) ToUTC() OffsetDateTime {
	return odt.WithOffsetSameInstant(ZoneOffsetUTC)
}

// Compare compares this date-time with another based on the instant (epoch seconds).
// Returns -1 if this is before other, 0 if equal, 1 if after.
func (odt OffsetDateTime) Compare(other OffsetDateTime) int {
	if odt.IsZero() && other.IsZero() {
		return 0
	}
	if odt.IsZero() {
		return -1
	}
	if other.IsZero() {
		return 1
	}
	// Compare instants
	t1 := odt.ToInstant()
	t2 := other.ToInstant()
	if t1.Before(t2) {
		return -1
	}
	if t1.After(t2) {
		return 1
	}
	return 0
}

// IsBefore returns true if this date-time instant is before the specified date-time instant.
func (odt OffsetDateTime) IsBefore(other OffsetDateTime) bool {
	return odt.Compare(other) < 0
}

// IsAfter returns true if this date-time instant is after the specified date-time instant.
func (odt OffsetDateTime) IsAfter(other OffsetDateTime) bool {
	return odt.Compare(other) > 0
}

// IsEqual returns true if this date-time represents the same instant as the specified date-time.
func (odt OffsetDateTime) IsEqual(other OffsetDateTime) bool {
	return odt.Compare(other) == 0
}

// PlusDays returns a copy with the specified number of days added.
func (odt OffsetDateTime) PlusDays(days int) OffsetDateTime {
	return NewOffsetDateTime(odt.dateTime.PlusDays(days), odt.offset)
}

// MinusDays returns a copy with the specified number of days subtracted.
func (odt OffsetDateTime) MinusDays(days int) OffsetDateTime {
	return odt.PlusDays(-days)
}

// PlusMonths returns a copy with the specified number of months added.
func (odt OffsetDateTime) PlusMonths(months int) OffsetDateTime {
	return NewOffsetDateTime(odt.dateTime.PlusMonths(months), odt.offset)
}

// MinusMonths returns a copy with the specified number of months subtracted.
func (odt OffsetDateTime) MinusMonths(months int) OffsetDateTime {
	return odt.PlusMonths(-months)
}

// PlusYears returns a copy with the specified number of years added.
func (odt OffsetDateTime) PlusYears(years int) OffsetDateTime {
	return NewOffsetDateTime(odt.dateTime.PlusYears(years), odt.offset)
}

// MinusYears returns a copy with the specified number of years subtracted.
func (odt OffsetDateTime) MinusYears(years int) OffsetDateTime {
	return odt.PlusYears(-years)
}

// PlusHours returns a copy with the specified number of hours added.
func (odt OffsetDateTime) PlusHours(hours int) OffsetDateTime {
	if odt.IsZero() {
		return odt
	}
	instant := odt.ToInstant().Add(time.Duration(hours) * time.Hour)
	return NewOffsetDateTimeByGoTimeWithOffset(instant, odt.offset)
}

// MinusHours returns a copy with the specified number of hours subtracted.
func (odt OffsetDateTime) MinusHours(hours int) OffsetDateTime {
	return odt.PlusHours(-hours)
}

// PlusMinutes returns a copy with the specified number of minutes added.
func (odt OffsetDateTime) PlusMinutes(minutes int) OffsetDateTime {
	if odt.IsZero() {
		return odt
	}
	instant := odt.ToInstant().Add(time.Duration(minutes) * time.Minute)
	return NewOffsetDateTimeByGoTimeWithOffset(instant, odt.offset)
}

// MinusMinutes returns a copy with the specified number of minutes subtracted.
func (odt OffsetDateTime) MinusMinutes(minutes int) OffsetDateTime {
	return odt.PlusMinutes(-minutes)
}

// PlusSeconds returns a copy with the specified number of seconds added.
func (odt OffsetDateTime) PlusSeconds(seconds int) OffsetDateTime {
	if odt.IsZero() {
		return odt
	}
	instant := odt.ToInstant().Add(time.Duration(seconds) * time.Second)
	return NewOffsetDateTimeByGoTimeWithOffset(instant, odt.offset)
}

// MinusSeconds returns a copy with the specified number of seconds subtracted.
func (odt OffsetDateTime) MinusSeconds(seconds int) OffsetDateTime {
	return odt.PlusSeconds(-seconds)
}

// PlusNanoseconds returns a copy with the specified number of nanoseconds added.
func (odt OffsetDateTime) PlusNanoseconds(nanos int64) OffsetDateTime {
	if odt.IsZero() {
		return odt
	}
	instant := odt.ToInstant().Add(time.Duration(nanos))
	return NewOffsetDateTimeByGoTimeWithOffset(instant, odt.offset)
}

// MinusNanoseconds returns a copy with the specified number of nanoseconds subtracted.
func (odt OffsetDateTime) MinusNanoseconds(nanos int64) OffsetDateTime {
	return odt.PlusNanoseconds(-nanos)
}

// NewOffsetDateTimeByGoTimeWithOffset creates an OffsetDateTime from a time.Time and a specific offset.
// The instant is preserved, but the local date-time is adjusted to the specified offset.
func NewOffsetDateTimeByGoTimeWithOffset(t time.Time, offset ZoneOffset) OffsetDateTime {
	if t.IsZero() {
		return OffsetDateTime{}
	}
	// Adjust time to the specified offset
	adjusted := t.Add(time.Duration(offset.TotalSeconds()) * time.Second)
	return OffsetDateTime{
		dateTime: NewLocalDateTimeByGoTime(adjusted),
		offset:   offset,
	}
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
	b, _ = odt.dateTime.AppendText(b)
	b, _ = odt.offset.AppendText(b)
	return b, nil
}

// MarshalText implements encoding.TextMarshaler.
func (odt OffsetDateTime) MarshalText() ([]byte, error) {
	return marshalTextImpl(odt)
}

// UnmarshalText implements encoding.TextUnmarshaler.
// Accepts ISO 8601 format: yyyy-MM-ddTHH:mm:ss[.nnnnnnnnn]±HH:mm or Z
func (odt *OffsetDateTime) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*odt = OffsetDateTime{}
		return nil
	}

	// Find the offset part (last '+', '-', or 'Z')
	offsetIdx := -1
	for i := len(text) - 1; i >= 0; i-- {
		if text[i] == '+' || text[i] == '-' || text[i] == 'Z' || text[i] == 'z' {
			offsetIdx = i
			break
		}
	}

	if offsetIdx < 0 {
		return newError("invalid offset date-time format: missing offset")
	}

	// Parse the datetime part
	var dt LocalDateTime
	if err := dt.UnmarshalText(text[:offsetIdx]); err != nil {
		return err
	}

	// Parse the offset part
	var offset ZoneOffset
	if err := offset.UnmarshalText(text[offsetIdx:]); err != nil {
		return err
	}

	*odt = NewOffsetDateTime(dt, offset)
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
		*odt = NewOffsetDateTimeByGoTime(v)
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
