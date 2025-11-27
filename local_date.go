package goda

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"math"
	"time"
)

// LocalDate represents a date without a time zone in the Gregorian calendar system,
// such as 2024-03-15. It stores the year, month, and day-of-month.
//
// LocalDate is comparable and can be used as a map key.
// The zero value represents an unset date and IsZero returns true for it.
//
// LocalDate implements sql.Scanner and driver.Valuer for database operations,
// encoding.TextMarshaler and encoding.TextUnmarshaler for text serialization,
// and json.Marshaler and json.Unmarshaler for JSON serialization.
//
// Format: yyyy-MM-dd (e.g., "2024-03-15"). Uses ISO 8601 basic calendar date format.
// Week dates (YYYY-Www-D) and ordinal dates (YYYY-DDD) are not supported.
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
		*d = LocalDateOfGoTime(v)
		return nil
	default:
		return sqlScannerDefaultBranch(v)
	}
}

// Value implements the driver.Valuer interface.
// It returns nil for zero values, otherwise returns the date as a string in yyyy-MM-dd format.
func (d LocalDate) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return d.String(), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It accepts JSON strings in yyyy-MM-dd format or JSON null.
func (d *LocalDate) UnmarshalJSON(bytes []byte) error {
	if len(bytes) == 4 && string(bytes) == "null" {
		*d = LocalDate{}
		return nil
	}
	return unmarshalJsonImpl(d, bytes)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It parses dates in yyyy-MM-dd format.
// Empty input is treated as zero value.
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
	dd, e := LocalDateOf(Year(y), Month(m), dom)
	if e != nil {
		return
	}
	*d = dd
	return
}

// MarshalText implements the encoding.TextMarshaler interface.
// It returns the date in yyyy-MM-dd format, or empty for zero value.
func (d LocalDate) MarshalText() (text []byte, err error) {
	return marshalTextImpl(d)
}

// MarshalJSON implements the json.Marshaler interface.
// It returns the date as a JSON string in yyyy-MM-dd format, or empty string for zero value.
func (d LocalDate) MarshalJSON() ([]byte, error) {
	return marshalJsonImpl(d)
}

// String returns the date in yyyy-MM-dd format, or empty string for zero value.
func (d LocalDate) String() string {
	return stringImpl(d)
}

// AppendText implements the encoding.TextAppender interface.
// It appends the date in yyyy-MM-dd format to b and returns the extended buffer.
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
	case FieldDayOfWeek, FieldDayOfMonth, FieldDayOfYear, FieldEpochDay, FieldMonthOfYear, FieldProlepticMonth, FieldYearOfEra, FieldYear, FieldEra:
		return true
	default:
		return false
	}
}

// GetField returns the value of the specified field as a TemporalValue.
// This method queries the date for the value of the specified field.
// The returned value may be unsupported if the field is not supported by LocalDate.
//
// If the date is zero (IsZero() returns true), an unsupported TemporalValue is returned.
// For fields not supported by LocalDate (such as time fields), an unsupported TemporalValue is returned.
//
// Supported fields include:
//   - FieldDayOfWeek: returns the day of week (1=Monday, 7=Sunday)
//   - FieldDayOfMonth: returns the day of month (1-31)
//   - FieldDayOfYear: returns the day of year (1-366)
//   - FieldMonthOfYear: returns the month (1=January, 12=December)
//   - FieldYear: returns the proleptic year
//   - FieldYearOfEra: returns the year within the era (same as FieldYear for CE dates)
//   - FieldEra: returns the era (0=BCE, 1=CE)
//   - FieldEpochDay: returns the number of days since Unix epoch (1970-01-01)
//   - FieldProlepticMonth: returns the number of months since year 0
//
// Overflow Analysis:
// None of the supported fields can overflow int64 in practice:
//   - FieldDayOfWeek: range 1-7, cannot overflow
//   - FieldDayOfMonth: range 1-31, cannot overflow
//   - FieldDayOfYear: range 1-366, cannot overflow
//   - FieldMonthOfYear: range 1-12, cannot overflow
//   - FieldYear/FieldYearOfEra: Year is int64, direct cast, cannot overflow
//   - FieldEra: values 0 or 1, cannot overflow
//   - FieldEpochDay: int64, calculated from year/month/day which are bounded by LocalDate's internal representation
//   - FieldProlepticMonth: Year * 12 + Month.
//     Year is int64, so max value is approximately int64_max * 12,
//     which would overflow.
//     However, LocalDate stores Year in the upper 48 bits of a 64-bit value,
//     limiting the practical range to approximately ±140 trillion years, making overflow impossible
//     in any realistic scenario.
func (d LocalDate) GetField(field Field) TemporalValue {
	if d.IsZero() {
		return TemporalValue{v: 0, unsupported: true}
	}
	var v int64
	switch field {
	case FieldDayOfWeek:
		// Range: 1-7, no overflow possible
		v = int64(d.DayOfWeek())
	case FieldDayOfMonth:
		// Range: 1-31, no overflow possible
		v = int64(d.DayOfMonth())
	case FieldDayOfYear:
		// Range: 1-366, no overflow possible
		v = int64(d.DayOfYear())
	case FieldMonthOfYear:
		// Range: 1-12, no overflow possible
		v = int64(d.Month())
	case FieldYear:
		// Year is already int64, direct cast, no overflow possible
		v = int64(d.Year())
	case FieldYearOfEra:
		v = int64(d.Year())
		if v < 0 {
			v = -v + 1
		}
	case FieldEra:
		// Values: 0 or 1, no overflow possible
		if d.Year() >= 1 {
			v = 1 // CE (Common FieldEra)
		} else {
			v = 0 // BCE (Before Common FieldEra)
		}
	case FieldEpochDay:
		// Returns int64, computed from bounded year/month/day values, no overflow in practice
		v = d.UnixEpochDays()
	case FieldProlepticMonth:
		// Calculate proleptic month (months since year 0)
		// Year is stored in 48 bits internally, limiting range to ±140 trillion years
		// Year * 12 cannot overflow int64 in this constrained range
		v = int64(d.Year())*12 + int64(d.Month()) - 1
	default:
		return TemporalValue{unsupported: true}
	}
	return TemporalValue{v: v}
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
// Negative values subtract days.
// Returns zero value if called on zero value.
func (d LocalDate) PlusDays(days int) LocalDate {
	if d.IsZero() {
		return d
	}
	return LocalDateOfUnixEpochDays(d.UnixEpochDays() + int64(days))
}

// MinusDays returns a copy of this date with the specified number of days subtracted.
// Equivalent to PlusDays(-days).
func (d LocalDate) MinusDays(days int) LocalDate {
	return d.PlusDays(-days)
}

// PlusMonths returns a copy of this date with the specified number of months added.
// Negative values subtract months.
// If the resulting day-of-month is invalid,
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
	return MustLocalDateOf(Year(y), Month(m), min(d.DayOfMonth(), Month(m).Length(Year(y).IsLeapYear())))
}

// MinusMonths returns a copy of this date with the specified number of months subtracted.
// Equivalent to PlusMonths(-months).
func (d LocalDate) MinusMonths(months int) LocalDate {
	return d.PlusMonths(-months)
}

// PlusYears returns a copy of this date with the specified number of years added.
// Negative values subtract years.
// If the resulting day-of-month is invalid
// (e.g., Feb 29 in a non-leap year), it is adjusted to the last valid day of the month.
// Returns zero value if called on zero value.
func (d LocalDate) PlusYears(years int) LocalDate {
	if d.IsZero() {
		return d
	}
	var year = Year(d.Year().Int64() + int64(years))
	return MustLocalDateOf(year, d.Month(), min(d.DayOfMonth(), d.Month().Length(year.IsLeapYear())))
}

// MinusYears returns a copy of this date with the specified number of years subtracted.
// Equivalent to PlusYears(-years).
func (d LocalDate) MinusYears(years int) LocalDate {
	return d.PlusYears(-years)
}

// PlusWeeks returns a copy of this date with the specified number of weeks added.
// Negative values subtract weeks.
// This is equivalent to PlusDays(weeks * 7).
// Returns zero value if called on zero value.
func (d LocalDate) PlusWeeks(weeks int) LocalDate {
	return d.PlusDays(7 * weeks)
}

// MinusWeeks returns a copy of this date with the specified number of weeks subtracted.
// Equivalent to PlusWeeks(-weeks).
func (d LocalDate) MinusWeeks(weeks int) LocalDate {
	return d.MinusDays(7 * weeks)
}

// WithTemporal returns a copy of this LocalDate with the specified field replaced.
// Zero values return zero immediately.
//
// Supported fields mirror Java's LocalDate#with(TemporalField, long):
//   - FieldDayOfMonth: sets the day-of-month while keeping year and month.
//   - FieldDayOfYear: sets the day-of-year while keeping year.
//   - FieldMonthOfYear: sets the month-of-year while keeping year and day-of-month (adjusted if necessary).
//   - FieldYear: sets the year while keeping month and day-of-month (adjusted if necessary).
//   - FieldYearOfEra: sets the year within the current era (same as FieldYear for CE dates).
//   - FieldEra: switches between BCE/CE eras while preserving year, month, and day.
//   - FieldEpochDay: sets the date based on days since Unix epoch (1970-01-01).
//   - FieldProlepticMonth: sets the date based on months since year 0.
//
// Fields outside this list return an error. Range violations propagate the validation error.
func (d LocalDate) WithTemporal(field Field, value TemporalValue) (LocalDate, error) {
	if d.IsZero() {
		return d, nil
	}

	v := value.Int64()

	switch field {
	case FieldDayOfMonth:
		e := checkTemporalInRange(FieldDayOfMonth, 1, 31, value, nil)
		if e != nil {
			return d, e
		}
		return d.WithDayOfMonth(int(v))
	case FieldDayOfYear:
		e := checkTemporalInRange(FieldDayOfYear, 1, 366, value, nil)
		if e != nil {
			return d, e
		}
		return d.WithDayOfYear(int(v))
	case FieldMonthOfYear:
		e := checkTemporalInRange(FieldMonthOfYear, 1, 12, value, nil)
		if e != nil {
			return d, e
		}
		return d.WithMonth(Month(v))
	case FieldYear:
		// Year can be any int value, no range check needed for the field itself
		// LocalDateOf will validate the resulting date
		return d.WithYear(Year(v))
	case FieldYearOfEra:
		// For CE (positive years), YearOfEra equals Year
		// For BCE (negative years), YearOfEra = -Year + 1
		e := checkTemporalInRange(FieldYearOfEra, 1, math.MaxInt, value, nil)
		if e != nil {
			return d, e
		}
		var y int64
		if d.Year() >= 1 {
			// CE
			y = v
		} else {
			// BCE: convert YearOfEra back to negative year
			y = -(v - 1)
		}
		result, err := LocalDateOf(Year(y), d.Month(), d.DayOfMonth())
		return result, err
	case FieldEra:
		e := checkTemporalInRange(FieldEra, 0, 1, value, nil)
		if e != nil {
			return d, e
		}
		var y = int64(d.Year())
		if v == 0 && d.Year() > 0 {
			// Switch from CE to BCE
			y = -int64(d.Year()) + 1
		} else if v == 1 && d.Year() < 0 {
			// Switch from BCE to CE
			y = -int64(d.Year()) + 1
		}
		// If already in the correct era, do nothing
		result, err := LocalDateOf(Year(y), d.Month(), d.DayOfMonth())
		return result, err
	case FieldEpochDay:
		// Convert epoch days back to LocalDate
		return LocalDateOfUnixEpochDays(v), nil
	case FieldProlepticMonth:
		// Convert proleptic month back to year and month
		// Proleptic month = year * 12 + (month - 1)
		year := v / 12
		month := v%12 + 1
		// Ensure year fits in our Year type
		if year < math.MinInt || year > math.MaxInt {
			return d, newError("proleptic month %d results in year out of range", v)
		}
		result, err := LocalDateOf(Year(year), Month(month), d.DayOfMonth())
		return result, err
	default:
		return d, newError("unsupported field: %v", field)
	}
}

// WithDayOfMonth returns a copy of this date with the day-of-month altered.
// If the day-of-month is invalid for the month and year, an error is returned.
// Returns zero value without error if called on zero value.
func (d LocalDate) WithDayOfMonth(dayOfMonth int) (r LocalDate, e error) {
	if d.IsZero() {
		return
	}
	r, e = LocalDateOf(d.Year(), d.Month(), dayOfMonth)
	if e != nil {
		e = newError("dayOfMonth %d out of range", dayOfMonth)
	}
	return
}

// MustWithDayOfMonth returns a copy of this date with the day-of-month altered.
// Panics if the day-of-month is invalid.
// Use WithDayOfMonth for error handling.
func (d LocalDate) MustWithDayOfMonth(dayOfMonth int) LocalDate {
	return mustValue(d.WithDayOfMonth(dayOfMonth))
}

// WithDayOfYear returns a copy of this date with the day-of-year altered.
// The day-of-year must be valid for the year (1-365 for non-leap years, 1-366 for leap years).
// Returns zero value without error if called on zero value.
func (d LocalDate) WithDayOfYear(dayOfYear int) (r LocalDate, e error) {
	if d.IsZero() {
		return
	}
	for m := December; m > 0; m-- {
		if dayOfYear >= m.FirstDayOfYear(d.Year().IsLeapYear()) {
			return LocalDateOf(d.Year(), m, dayOfYear-m.FirstDayOfYear(d.Year().IsLeapYear())+1)
		}
	}
	e = newError("dayOfYear %d out of range", dayOfYear)
	return
}

// MustWithDayOfYear returns a copy of this date with the day-of-year altered.
// Panics if the day-of-year is invalid.
// Use WithDayOfYear for error handling.
func (d LocalDate) MustWithDayOfYear(dayOfYear int) LocalDate {
	return mustValue(d.WithDayOfYear(dayOfYear))
}

// WithMonth returns a copy of this date with the month altered.
// If the resulting day-of-month is invalid for the new month,
// it is adjusted to the last valid day of the month.
// For example, January 31 with month set to February becomes February 28/29.
// Returns zero value without error if called on zero value.
func (d LocalDate) WithMonth(month Month) (r LocalDate, e error) {
	if d.IsZero() {
		return
	}
	if month < January || month > December {
		e = newError("month %d out of range", month)
		return
	}
	return LocalDateOf(d.Year(), month, min(d.DayOfMonth(), month.Length(d.Year().IsLeapYear())))
}

// MustWithMonth returns a copy of this date with the month altered.
// Panics if the month is invalid.
// Use WithMonth for error handling.
func (d LocalDate) MustWithMonth(month Month) LocalDate {
	return mustValue(d.WithMonth(month))
}

// WithYear returns a copy of this date with the year altered.
// If the resulting day-of-month is invalid for the new year
// (e.g., February 29 in a non-leap year), it is adjusted to the last valid day of the month.
// Returns zero value without error if called on zero value.
func (d LocalDate) WithYear(year Year) (r LocalDate, e error) {
	if d.IsZero() {
		return
	}
	return LocalDateOf(year, d.Month(), min(d.DayOfMonth(), d.Month().Length(year.IsLeapYear())))
}

// MustWithYear returns a copy of this date with the year altered.
// Panics if the resulting date is invalid.
// Use WithYear for error handling.
func (d LocalDate) MustWithYear(year Year) LocalDate {
	return mustValue(d.WithYear(year))
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

// AtTime combines this date with a time to create a LocalDateTime.
func (d LocalDate) AtTime(time LocalTime) LocalDateTime {
	return LocalDateTime{
		date: d,
		time: time,
	}
}

// LengthOfMonth returns the number of days in the month of this date.
// Returns 28, 29, 30, or 31 depending on the month and whether it's a leap year.
// Returns 0 for zero value.
func (d LocalDate) LengthOfMonth() int {
	if d.IsZero() {
		return 0
	}
	return d.Month().Length(d.Year().IsLeapYear())
}

// LengthOfYear returns the number of days in the year of this date.
// Returns 365 for non-leap years and 366 for leap years.
// Returns 0 for zero value.
func (d LocalDate) LengthOfYear() int {
	if d.IsZero() {
		return 0
	}
	return d.Year().Length()
}

// IsZero returns true if this is the zero value of LocalDate.
func (d LocalDate) IsZero() bool {
	return d.v == 0
}

// LocalDateOf creates a new LocalDate from the specified year, month, and day-of-month.
// Returns an error if the date is invalid (e.g., month out of range 1-12,
// day out of range for the month, or February 29 in a non-leap year).
func LocalDateOf(year Year, month Month, dayOfMonth int) (d LocalDate, e error) {
	if year > 1<<47-1 || year < -(1<<47) {
		e = newError("year %d out of range", year)
		return
	}
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

// MustLocalDateOf creates a new LocalDate from the specified year, month, and day-of-month.
// Panics if the date is invalid.
// Use LocalDateOf for error handling.
func MustLocalDateOf(year Year, month Month, dayOfMonth int) LocalDate {
	return mustValue(LocalDateOf(year, month, dayOfMonth))
}

// LocalDateOfGoTime creates a LocalDate from a time.Time.
// The time zone and time-of-day components are ignored.
// Returns zero value if t.IsZero().
func LocalDateOfGoTime(t time.Time) LocalDate {
	if t.IsZero() {
		return LocalDate{}
	}
	return MustLocalDateOf(Year(t.Year()), Month(t.Month()), t.Day())
}

// LocalDateNow returns the current date in the system's local time zone.
// This is equivalent to LocalDateOfGoTime(time.Now()).
// For UTC time, use LocalDateNowUTC.
// For a specific timezone, use LocalDateNowIn.
func LocalDateNow() LocalDate {
	return LocalDateOfGoTime(time.Now())
}

// LocalDateNowIn returns the current date in the specified time zone.
// This is equivalent to LocalDateOfGoTime(time.Now().In(loc)).
func LocalDateNowIn(loc *time.Location) LocalDate {
	return LocalDateOfGoTime(time.Now().In(loc))
}

// LocalDateNowUTC returns the current date in UTC.
// This is equivalent to LocalDateOfGoTime(time.Now().UTC()).
func LocalDateNowUTC() LocalDate {
	return LocalDateOfGoTime(time.Now().UTC())
}

// LocalDateParse parses a date string in yyyy-MM-dd format.
// Returns an error if the string is invalid or represents an invalid date.
//
// Supported format: yyyy-MM-dd (e.g., "2024-03-15")
// Week dates and ordinal dates are not supported.
//
// Example:
//
//	date, err := LocalDateParse("2024-03-15")
//	if err != nil {
//	    // handle error
//	}
func LocalDateParse(s string) (LocalDate, error) {
	var d LocalDate
	err := d.UnmarshalText([]byte(s))
	return d, err
}

// MustLocalDateParse parses a date string in yyyy-MM-dd format.
// Panics if the string is invalid.
// Use LocalDateParse for error handling.
//
// Example:
//
//	date := MustLocalDateParse("2024-03-15")
func MustLocalDateParse(s string) LocalDate {
	return mustValue(LocalDateParse(s))
}

// LocalDateOfUnixEpochDays creates a LocalDate from the number of days since Unix epoch (1970-01-01).
// Positive values represent dates after the epoch, negative before.
func LocalDateOfUnixEpochDays(days int64) LocalDate {
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
	return MustLocalDateOf(Year(yearEst), month, dom)
}

var (
	_ encoding.TextAppender    = (*LocalDate)(nil)
	_ fmt.Stringer             = (*LocalDate)(nil)
	_ encoding.TextMarshaler   = (*LocalDate)(nil)
	_ encoding.TextUnmarshaler = (*LocalDate)(nil)
	_ json.Marshaler           = (*LocalDate)(nil)
	_ json.Unmarshaler         = (*LocalDate)(nil)
	_ driver.Valuer            = (*LocalDate)(nil)
	_ sql.Scanner              = (*LocalDate)(nil)
)

// Compile-time check that LocalDate is comparable
func _assertLocalDateIsComparable[T comparable](t T) {}

var _ = _assertLocalDateIsComparable[LocalDate]
