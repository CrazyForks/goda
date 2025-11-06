package goda

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"time"
)

// LocalDate represents a date without a time zone in the ISO-8601 calendar system,
// such as 2024-03-15. It stores the year, month, and day-of-month.
//
// LocalDate is comparable and can be used as a map key.
// The zero value represents an unset date and IsZero returns true for it.
//
// LocalDate implements sql.Scanner and driver.Valuer for database operations,
// encoding.TextMarshaler and encoding.TextUnmarshaler for text serialization,
// and json.Marshaler and json.Unmarshaler for JSON serialization.
// The text format is YYYY-MM-DD (ISO 8601).
type LocalDate struct {
	v int64
}

// Scan implements the sql.Scanner interface.
// It supports scanning from nil, string, []byte, and time.Time.
// Nil values are converted to the zero value of LocalDate.
func (d *LocalDate) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*d = LocalDate{}
		return nil
	case string:
		return d.UnmarshalText([]byte(v))
	case []byte:
		return d.UnmarshalText(v)
	case time.Time:
		*d = NewLocalDateByGoTime(v)
		return nil
	default:
		return sqlScannerDefaultBranch(v)
	}
}

// Value implements the driver.Valuer interface.
// It returns nil for zero values, otherwise returns the date as a string in YYYY-MM-DD format.
func (d *LocalDate) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return d.String(), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It accepts JSON strings in YYYY-MM-DD format or JSON null.
func (d *LocalDate) UnmarshalJSON(bytes []byte) error {
	if len(bytes) == 4 && string(bytes) == "null" {
		*d = LocalDate{}
		return nil
	}
	return unmarshalJsonImpl(d, bytes)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It parses dates in YYYY-MM-DD format. Empty input is treated as zero value.
func (d *LocalDate) UnmarshalText(text []byte) (e error) {
	if len(text) == 0 {
		*d = LocalDate{}
		return nil
	}
	if text[len(text)-3] != '-' || text[len(text)-6] != '-' {
		return unmarshalError(text)
	}
	var y int64
	var m, dom int
	dom, e = parseInt(text[len(text)-2:])
	if e != nil {
		return
	}
	m, e = parseInt(text[len(text)-5 : len(text)-3])
	if e != nil {
		return
	}
	y, e = parseInt64(text[:len(text)-6])
	if e != nil {
		return
	}
	dd, e := NewLocalDate(Year(y), Month(m), dom)
	if e != nil {
		return
	}
	*d = dd
	return
}

// MarshalText implements the encoding.TextMarshaler interface.
// It returns the date in YYYY-MM-DD format, or empty for zero value.
func (d LocalDate) MarshalText() (text []byte, err error) {
	return marshalTextImpl(d)
}

// MarshalJSON implements the json.Marshaler interface.
// It returns the date as a JSON string in YYYY-MM-DD format, or empty string for zero value.
func (d LocalDate) MarshalJSON() ([]byte, error) {
	return marshalJsonImpl(d)
}

// String returns the date in YYYY-MM-DD format, or empty string for zero value.
func (d LocalDate) String() string {
	return stringImpl(d)
}

// AppendText implements the encoding.TextAppender interface.
// It appends the date in YYYY-MM-DD format to b and returns the extended buffer.
func (d LocalDate) AppendText(b []byte) ([]byte, error) {
	if d.IsZero() {
		return b, nil
	}
	b, _ = d.Year().AppendText(b)
	b = append(b, '-', byte('0'+d.Month()/10), byte('0'+d.Month()%10), '-', byte('0'+d.DayOfMonth()/10), byte('0'+d.DayOfMonth()%10))
	return b, nil
}

// Year returns the year component of this date.
func (d LocalDate) Year() Year {
	return Year(d.v >> 16)
}

// IsLeapYear returns true if the year of this date is a leap year.
// A leap year is divisible by 4, unless it's divisible by 100 (but not 400).
func (d LocalDate) IsLeapYear() bool {
	return d.Year().IsLeapYear()
}

// IsSupportedField returns true if the field is supported by LocalDate.
func (d LocalDate) IsSupportedField(field Field) bool {
	switch field {
	case DayOfWeekField, DayOfMonth, DayOfYear, EpochDay,
		MonthOfYear, ProlepticMonth, YearOfEra, YearField, Era:
		return true
	default:
		return false
	}
}

// GetFieldInt64 returns the value of the specified field as an int64.
// Returns 0 if the field is not supported or the date is zero.
func (d LocalDate) GetFieldInt64(field Field) int64 {
	if d.IsZero() {
		return 0
	}
	switch field {
	case DayOfWeekField:
		return int64(d.DayOfWeek())
	case DayOfMonth:
		return int64(d.DayOfMonth())
	case DayOfYear:
		return int64(d.DayOfYear())
	case MonthOfYear:
		return int64(d.Month())
	case YearField, YearOfEra:
		return int64(d.Year())
	case Era:
		if d.Year() >= 1 {
			return 1 // CE (Common Era)
		}
		return 0 // BCE (Before Common Era)
	case EpochDay:
		return d.UnixEpochDays()
	case ProlepticMonth:
		// Calculate proleptic month (months since year 0)
		return int64(d.Year())*12 + int64(d.Month()) - 1
	default:
		return 0
	}
}

// Month returns the month component of this date (1-12).
func (d LocalDate) Month() Month {
	return Month(d.v >> 8 & 0xff)
}

// DayOfMonth returns the day-of-month component of this date (1-31).
func (d LocalDate) DayOfMonth() int {
	return int(d.v & 0xff)
}

// DayOfWeek returns the day-of-week for this date.
// Returns 0 for zero value, otherwise Monday=1 through Sunday=7.
func (d LocalDate) DayOfWeek() DayOfWeek {
	if d.IsZero() {
		return 0
	}
	return DayOfWeek(floorMod(d.UnixEpochDays()+3, 7) + 1)
}

// DayOfYear returns the day-of-year for this date (1-366).
// Returns 0 for zero value.
func (d LocalDate) DayOfYear() int {
	if d.IsZero() {
		return 0
	}
	return d.Month().FirstDayOfYear(d.IsLeapYear()) - 1 + d.DayOfMonth()
}

// PlusDays returns a copy of this date with the specified number of days added.
// Negative values subtract days. Returns zero value if called on zero value.
func (d LocalDate) PlusDays(days int) LocalDate {
	if d.IsZero() {
		return d
	}
	return NewLocalDateByUnixEpochDays(d.UnixEpochDays() + int64(days))
}

// MinusDays returns a copy of this date with the specified number of days subtracted.
// Equivalent to PlusDays(-days).
func (d LocalDate) MinusDays(days int) LocalDate {
	return d.PlusDays(-days)
}

// PlusMonths returns a copy of this date with the specified number of months added.
// Negative values subtract months. If the resulting day-of-month is invalid,
// it is adjusted to the last valid day of the month.
// For example, 2024-01-31 plus 1 month becomes 2024-02-29 (leap year).
// Returns zero value if called on zero value.
func (d LocalDate) PlusMonths(months int) LocalDate {
	if d.IsZero() {
		return d
	}
	var m = int(d.Month()) + months
	var y = d.Year().Int64()
	if m > 12 {
		y += int64((m - 1) / 12)
		m = (m-1)%12 + 1
	} else if m < 1 {
		y += int64((m - 12) / 12)
		m = (m-1)%12 + 1
		if m < 1 {
			m += 12
		}
	}
	return MustNewLocalDate(Year(y), Month(m), min(d.DayOfMonth(), Month(m).Length(Year(y).IsLeapYear())))
}

// MinusMonths returns a copy of this date with the specified number of months subtracted.
// Equivalent to PlusMonths(-months).
func (d LocalDate) MinusMonths(months int) LocalDate {
	return d.PlusMonths(-months)
}

// PlusYears returns a copy of this date with the specified number of years added.
// Negative values subtract years. If the resulting day-of-month is invalid
// (e.g., Feb 29 in a non-leap year), it is adjusted to the last valid day of the month.
// Returns zero value if called on zero value.
func (d LocalDate) PlusYears(years int) LocalDate {
	if d.IsZero() {
		return d
	}
	var year = Year(d.Year().Int64() + int64(years))
	return MustNewLocalDate(year, d.Month(), min(d.DayOfMonth(), d.Month().Length(year.IsLeapYear())))
}

// MinusYears returns a copy of this date with the specified number of years subtracted.
// Equivalent to PlusYears(-years).
func (d LocalDate) MinusYears(years int) LocalDate {
	return d.PlusYears(-years)
}

// Compare compares this date with another date.
// Returns -1 if this date is before other, 0 if equal, and 1 if after.
// Zero values are considered less than non-zero values.
func (d LocalDate) Compare(other LocalDate) int {
	return doCompare(d, other, compareZero, comparing(LocalDate.Year), comparing(LocalDate.Month), comparing(LocalDate.DayOfMonth))
}

// IsBefore returns true if this date is before the specified date.
func (d LocalDate) IsBefore(other LocalDate) bool {
	return d.Compare(other) < 0
}

// IsAfter returns true if this date is after the specified date.
func (d LocalDate) IsAfter(other LocalDate) bool {
	return d.Compare(other) > 0
}

// GoTime converts this date to a time.Time at midnight UTC.
// Returns time.Time{} (zero) for zero value.
func (d LocalDate) GoTime() time.Time {
	if d.IsZero() {
		return time.Time{}
	}
	return time.Date(int(d.Year()), time.Month(d.Month()), d.DayOfMonth(), 0, 0, 0, 0, time.UTC)
}

// UnixEpochDays returns the number of days since Unix epoch (1970-01-01).
// Positive values represent dates after the epoch, negative before.
// Returns 0 for zero value.
func (d LocalDate) UnixEpochDays() int64 {
	if d.IsZero() {
		return 0
	}
	const DaysPerCycle = 365*400 + 97
	const Days0000To1970 = (DaysPerCycle * 5) - (30*365 + 7)

	y := d.Year().Int64()
	m := int64(d.Month())
	day := int64(d.DayOfMonth())
	total := int64(0)

	// Calculate year contribution
	total += 365 * y
	if y >= 0 {
		total += (y+3)/4 - (y+99)/100 + (y+399)/400
	} else {
		total -= y/-4 - y/-100 + y/-400
	}

	// Calculate month contribution
	total += (367*m - 362) / 12

	// Calculate day contribution
	total += day - 1

	// Adjust for leap year if month > February
	if m > 2 {
		total--
		if !d.Year().IsLeapYear() {
			total--
		}
	}

	return total - Days0000To1970
}

// IsZero returns true if this is the zero value of LocalDate.
func (d LocalDate) IsZero() bool {
	return d.v == 0
}

var _ encoding.TextAppender = (*LocalDate)(nil)
var _ fmt.Stringer = (*LocalDate)(nil)
var _ encoding.TextMarshaler = (*LocalDate)(nil)
var _ encoding.TextUnmarshaler = (*LocalDate)(nil)
var _ json.Marshaler = (*LocalDate)(nil)
var _ json.Unmarshaler = (*LocalDate)(nil)
var _ driver.Valuer = (*LocalDate)(nil)
var _ sql.Scanner = (*LocalDate)(nil)

// NewLocalDate creates a new LocalDate from the specified year, month, and day-of-month.
// Returns an error if the date is invalid (e.g., month out of range 1-12,
// day out of range for the month, or February 29 in a non-leap year).
func NewLocalDate(year Year, month Month, dayOfMonth int) (d LocalDate, e error) {
	if month < January || month > December {
		e = newError("month %d out of range", month)
		return
	}
	if dayOfMonth < 1 || dayOfMonth > month.Length(year.IsLeapYear()) {
		e = newError("day %d of month out of range", dayOfMonth)
		return
	}
	d = LocalDate{
		v: int64(year)<<16 | int64(month)<<8 | int64(dayOfMonth),
	}
	return
}

// MustNewLocalDate creates a new LocalDate from the specified year, month, and day-of-month.
// Panics if the date is invalid. Use NewLocalDate for error handling.
func MustNewLocalDate(year Year, month Month, dayOfMonth int) LocalDate {
	nld, e := NewLocalDate(year, month, dayOfMonth)
	if e != nil {
		panic(e)
	}
	return nld
}

// NewLocalDateByGoTime creates a LocalDate from a time.Time.
// The time zone and time-of-day components are ignored.
// Returns zero value if t.IsZero().
func NewLocalDateByGoTime(t time.Time) LocalDate {
	if t.IsZero() {
		return LocalDate{}
	}
	return MustNewLocalDate(Year(t.Year()), Month(t.Month()), t.Day())
}

// LocalDateNow returns the current date in the system's local time zone.
// This is equivalent to NewLocalDateByGoTime(time.Now()).
// For UTC time, use LocalDateNowUTC. For a specific timezone, use LocalDateNowIn.
func LocalDateNow() LocalDate {
	return NewLocalDateByGoTime(time.Now())
}

// LocalDateNowIn returns the current date in the specified time zone.
// This is equivalent to NewLocalDateByGoTime(time.Now().In(loc)).
func LocalDateNowIn(loc *time.Location) LocalDate {
	return NewLocalDateByGoTime(time.Now().In(loc))
}

// LocalDateNowUTC returns the current date in UTC.
// This is equivalent to NewLocalDateByGoTime(time.Now().UTC()).
func LocalDateNowUTC() LocalDate {
	return NewLocalDateByGoTime(time.Now().UTC())
}

// ParseLocalDate parses a date string in ISO 8601 format (YYYY-MM-DD).
// Returns an error if the string is invalid or represents an invalid date.
// For lenient parsing that returns zero value on error, see ParseLocalDateOrZero.
//
// Example:
//
//	date, err := ParseLocalDate("2024-03-15")
//	if err != nil {
//	    // handle error
//	}
func ParseLocalDate(s string) (LocalDate, error) {
	var d LocalDate
	err := d.UnmarshalText([]byte(s))
	return d, err
}

// MustParseLocalDate parses a date string in ISO 8601 format (YYYY-MM-DD).
// Panics if the string is invalid. Use ParseLocalDate for error handling.
//
// Example:
//
//	date := MustParseLocalDate("2024-03-15")
func MustParseLocalDate(s string) LocalDate {
	d, err := ParseLocalDate(s)
	if err != nil {
		panic(err)
	}
	return d
}

// NewLocalDateByUnixEpochDays creates a LocalDate from the number of days since Unix epoch (1970-01-01).
// Positive values represent dates after the epoch, negative before.
func NewLocalDateByUnixEpochDays(days int64) LocalDate {
	const DaysPerCycle = 365*400 + 97
	const Days0000To1970 = (DaysPerCycle * 5) - (30*365 + 7)
	zeroDay := days + Days0000To1970

	// Adjust to March-based year (makes leap day fall at end of 4-year cycle)
	zeroDay -= 60 // Shift to 0000-03-01 as base

	var adjust int64
	if zeroDay < 0 {
		// For negative years, adjust to positive for calculation
		adjustCycles := (zeroDay+1)/DaysPerCycle - 1
		adjust = adjustCycles * 400
		zeroDay += -adjustCycles * DaysPerCycle
	}

	// Estimate the year
	yearEst := (400*zeroDay + 591) / DaysPerCycle
	// Calculate day-of-year for estimated year
	doyEst := zeroDay - (365*yearEst + yearEst/4 - yearEst/100 + yearEst/400)
	if doyEst < 0 {
		// Correct the estimate
		yearEst--
		doyEst = zeroDay - (365*yearEst + yearEst/4 - yearEst/100 + yearEst/400)
	}

	yearEst += adjust // Restore negative year adjustment

	// Convert from March-based day-of-year
	marchDoy0 := int(doyEst)

	// Convert back to January-based month and day
	marchMonth0 := (marchDoy0*5 + 2) / 153
	month := Month(marchMonth0 + 3)
	if month > 12 {
		month -= 12
	}
	dom := marchDoy0 - (marchMonth0*306+5)/10 + 1
	if marchDoy0 >= 306 {
		yearEst++
	}
	return MustNewLocalDate(Year(yearEst), month, dom)
}
