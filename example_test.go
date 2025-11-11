package goda_test

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/iseki0/goda"
)

// Example demonstrates basic usage of the goda package.
func Example() {
	// Create a specific date
	date := goda.MustNewLocalDate(2024, goda.March, 15)
	fmt.Println("LocalDate:", date)

	// Create a specific time
	timeOfDay := goda.MustNewLocalTime(14, 30, 45, 123456789)
	fmt.Println("LocalTime:", timeOfDay)

	// Get current date and time
	today := goda.LocalDateNow()
	now := goda.LocalTimeNow()
	fmt.Printf("Type of today: %T\n", today)
	fmt.Printf("Type of now: %T\n", now)

	// Output:
	// LocalDate: 2024-03-15
	// LocalTime: 14:30:45.123456789
	// Type of today: goda.LocalDate
	// Type of now: goda.LocalTime
}

// ExampleLocalDateNow demonstrates how to get the current date.
func ExampleLocalDateNow() {
	// Get current date in local timezone
	today := goda.LocalDateNow()

	// Check that we got a valid date
	fmt.Printf("Got valid date: %v\n", !today.IsZero())
	fmt.Printf("Has year component: %v\n", today.Year() != 0)

	// Output:
	// Got valid date: true
	// Has year component: true
}

// ExampleLocalDateNowUTC demonstrates how to get the current date in UTC.
func ExampleLocalDateNowUTC() {
	// Get current date in UTC
	todayUTC := goda.LocalDateNowUTC()

	// Verify it's a valid date
	fmt.Printf("Valid: %v\n", !todayUTC.IsZero())

	// Output:
	// Valid: true
}

// ExampleLocalDateNowIn demonstrates how to get the current date in a specific timezone.
func ExampleLocalDateNowIn() {
	// Get current date in Tokyo timezone
	tokyo, _ := time.LoadLocation("Asia/Tokyo")
	todayTokyo := goda.LocalDateNowIn(tokyo)

	// Verify it's valid
	fmt.Printf("Valid: %v\n", !todayTokyo.IsZero())

	// Output:
	// Valid: true
}

// ExampleParseLocalDate demonstrates parsing a date from a string.
func ExampleParseLocalDate() {
	// Parse a date string
	date, err := goda.ParseLocalDate("2024-03-15")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(date)

	// Output:
	// 2024-03-15
}

// ExampleMustParseLocalDate demonstrates parsing a date that panics on error.
func ExampleMustParseLocalDate() {
	// Parse a date string (panics if invalid)
	date := goda.MustParseLocalDate("2024-03-15")
	fmt.Println(date)

	// Output:
	// 2024-03-15
}

// ExampleNewLocalDate demonstrates how to create a date.
func ExampleNewLocalDate() {
	// Create a valid date
	date, err := goda.NewLocalDate(2024, goda.January, 15)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(date)

	// Try to create an invalid date
	_, err = goda.NewLocalDate(2024, goda.February, 30)
	fmt.Println("Error:", err)

	// Output:
	// 2024-01-15
	// Error: goda: day 30 of month out of range
}

// ExampleMustNewLocalDate demonstrates how to create a date that panics on error.
func ExampleMustNewLocalDate() {
	// Create a date (panics if invalid)
	date := goda.MustNewLocalDate(2024, goda.March, 15)
	fmt.Println(date)

	// Output:
	// 2024-03-15
}

// ExampleNewLocalDateByGoTime demonstrates converting from time.Time to LocalDate.
func ExampleNewLocalDateByGoTime() {
	// Convert from time.LocalTime
	goTime := time.Date(2024, time.March, 15, 14, 30, 0, 0, time.UTC)
	date := goda.NewLocalDateByGoTime(goTime)
	fmt.Println(date)

	// Output:
	// 2024-03-15
}

// ExampleLocalDate_PlusDays demonstrates adding days to a date.
func ExampleLocalDate_PlusDays() {
	date := goda.MustNewLocalDate(2024, goda.January, 15)
	fmt.Println("Original:", date)
	fmt.Println("Plus 10 days:", date.PlusDays(10))
	fmt.Println("Minus 10 days:", date.PlusDays(-10))

	// Output:
	// Original: 2024-01-15
	// Plus 10 days: 2024-01-25
	// Minus 10 days: 2024-01-05
}

// ExampleLocalDate_PlusMonths demonstrates adding months to a date.
func ExampleLocalDate_PlusMonths() {
	date := goda.MustNewLocalDate(2024, goda.January, 31)
	fmt.Println("Original:", date)
	fmt.Println("Plus 1 month:", date.PlusMonths(1))
	fmt.Println("Plus 2 months:", date.PlusMonths(2))

	// Output:
	// Original: 2024-01-31
	// Plus 1 month: 2024-02-29
	// Plus 2 months: 2024-03-31
}

// ExampleLocalDate_Compare demonstrates comparing dates.
func ExampleLocalDate_Compare() {
	date1 := goda.MustNewLocalDate(2024, goda.March, 15)
	date2 := goda.MustNewLocalDate(2024, goda.March, 20)
	date3 := goda.MustNewLocalDate(2024, goda.March, 15)

	fmt.Println("date1 < date2:", date1.IsBefore(date2))
	fmt.Println("date1 > date2:", date1.IsAfter(date2))
	fmt.Println("date1 == date3:", date1.Compare(date3) == 0)

	// Output:
	// date1 < date2: true
	// date1 > date2: false
	// date1 == date3: true
}

// ExampleLocalDate_DayOfWeek demonstrates getting the day of week.
func ExampleLocalDate_DayOfWeek() {
	date := goda.MustNewLocalDate(2024, goda.March, 15)
	fmt.Println("Day of week:", date.DayOfWeek())
	fmt.Println("Is Friday?", date.DayOfWeek() == goda.Friday)

	// Output:
	// Day of week: Friday
	// Is Friday? true
}

// ExampleLocalDate_AtTime demonstrates combining a date with a time to create a datetime.
func ExampleLocalDate_AtTime() {
	date := goda.MustNewLocalDate(2024, goda.March, 15)
	time := goda.MustNewLocalTime(14, 30, 45, 123456789)

	dateTime := date.AtTime(time)
	fmt.Println("Date:", date)
	fmt.Println("Time:", time)
	fmt.Println("DateTime:", dateTime)

	// Output:
	// Date: 2024-03-15
	// Time: 14:30:45.123456789
	// DateTime: 2024-03-15T14:30:45.123456789
}

// ExampleLocalTimeNow demonstrates how to get the current time.
func ExampleLocalTimeNow() {
	// Get current time in local timezone
	now := goda.LocalTimeNow()

	// Verify it's valid and within range
	fmt.Printf("Valid: %v\n", !now.IsZero())
	fmt.Printf("Hour in range: %v\n", now.Hour() >= 0 && now.Hour() < 24)

	// Output:
	// Valid: true
	// Hour in range: true
}

// ExampleLocalTimeNowUTC demonstrates how to get the current time in UTC.
func ExampleLocalTimeNowUTC() {
	// Get current time in UTC
	nowUTC := goda.LocalTimeNowUTC()

	// Verify it's valid
	fmt.Printf("Valid: %v\n", !nowUTC.IsZero())

	// Output:
	// Valid: true
}

// ExampleLocalTimeNowIn demonstrates how to get the current time in a specific timezone.
func ExampleLocalTimeNowIn() {
	// Get current time in Tokyo timezone
	tokyo, _ := time.LoadLocation("Asia/Tokyo")
	nowTokyo := goda.LocalTimeNowIn(tokyo)

	// Verify it's valid
	fmt.Printf("Valid: %v\n", !nowTokyo.IsZero())

	// Output:
	// Valid: true
}

// ExampleParseLocalTime demonstrates parsing a time from a string.
func ExampleParseLocalTime() {
	// Parse a time string
	t, err := goda.ParseLocalTime("14:30:45.123456789")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(t)

	// Output:
	// 14:30:45.123456789
}

// ExampleMustParseLocalTime demonstrates parsing a time that panics on error.
func ExampleMustParseLocalTime() {
	// Parse a time string (panics if invalid)
	t := goda.MustParseLocalTime("14:30:45")
	fmt.Println(t)

	// Output:
	// 14:30:45
}

// ExampleNewLocalTime demonstrates how to create a time.
func ExampleNewLocalTime() {
	// Create a valid time
	t, err := goda.NewLocalTime(14, 30, 45, 123456789)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(t)

	// Try to create an invalid time
	_, err = goda.NewLocalTime(25, 0, 0, 0)
	fmt.Println("Error:", err)

	// Output:
	// 14:30:45.123456789
	// Error: goda: hour 25 out of range
}

// ExampleMustNewLocalTime demonstrates how to create a time that panics on error.
func ExampleMustNewLocalTime() {
	// Create a time (panics if invalid)
	t := goda.MustNewLocalTime(14, 30, 45, 0)
	fmt.Println(t)

	// Midnight
	midnight := goda.MustNewLocalTime(0, 0, 0, 0)
	fmt.Println("Midnight:", midnight)

	// Output:
	// 14:30:45
	// Midnight: 00:00:00
}

// ExampleNewLocalTimeByGoTime demonstrates converting from time.Time to LocalTime.
func ExampleNewLocalTimeByGoTime() {
	// Convert from time.LocalTime
	goTime := time.Date(2024, time.March, 15, 14, 30, 45, 123456789, time.UTC)
	t := goda.NewLocalTimeByGoTime(goTime)
	fmt.Println(t)

	// Output:
	// 14:30:45.123456789
}

// ExampleLocalTime_hour demonstrates accessing time components.
func ExampleLocalTime_hour() {
	t := goda.MustNewLocalTime(14, 30, 45, 123456789)
	fmt.Println("Hour:", t.Hour())
	fmt.Println("Minute:", t.Minute())
	fmt.Println("Second:", t.Second())
	fmt.Println("Millisecond:", t.Millisecond())
	fmt.Println("Nanosecond:", t.Nanosecond())

	// Output:
	// Hour: 14
	// Minute: 30
	// Second: 45
	// Millisecond: 123
	// Nanosecond: 123456789
}

// ExampleLocalTime_Compare demonstrates comparing times.
func ExampleLocalTime_Compare() {
	t1 := goda.MustNewLocalTime(14, 30, 0, 0)
	t2 := goda.MustNewLocalTime(15, 0, 0, 0)
	t3 := goda.MustNewLocalTime(14, 30, 0, 0)

	fmt.Println("t1 < t2:", t1.IsBefore(t2))
	fmt.Println("t1 > t2:", t1.IsAfter(t2))
	fmt.Println("t1 == t3:", t1.Compare(t3) == 0)

	// Output:
	// t1 < t2: true
	// t1 > t2: false
	// t1 == t3: true
}

// ExampleLocalTime_AtDate demonstrates combining a time with a date to create a datetime.
func ExampleLocalTime_AtDate() {
	date := goda.MustNewLocalDate(2024, goda.March, 15)
	time := goda.MustNewLocalTime(14, 30, 45, 123456789)

	dateTime := time.AtDate(date)
	fmt.Println("Date:", date)
	fmt.Println("Time:", time)
	fmt.Println("DateTime:", dateTime)

	// Output:
	// Date: 2024-03-15
	// Time: 14:30:45.123456789
	// DateTime: 2024-03-15T14:30:45.123456789
}

// ExampleLocalTime_String demonstrates the string format with fractional seconds.
func ExampleLocalTime_String() {
	// LocalTime without fractional seconds
	t1 := goda.MustNewLocalTime(14, 30, 45, 0)
	fmt.Println(t1)

	// LocalTime with milliseconds
	t2 := goda.MustNewLocalTime(14, 30, 45, 123000000)
	fmt.Println(t2)

	// LocalTime with microseconds
	t3 := goda.MustNewLocalTime(14, 30, 45, 123456000)
	fmt.Println(t3)

	// LocalTime with nanoseconds
	t4 := goda.MustNewLocalTime(14, 30, 45, 123456789)
	fmt.Println(t4)

	// LocalTime with trailing zeros aligned to 3-digit boundaries
	t5 := goda.MustNewLocalTime(14, 30, 45, 100000000)
	fmt.Println(t5)

	// Output:
	// 14:30:45
	// 14:30:45.123
	// 14:30:45.123456
	// 14:30:45.123456789
	// 14:30:45.100
}

// ExampleMonth demonstrates working with months.
func ExampleMonth() {
	// Months are constants from January (1) to December (12)
	fmt.Println("January:", goda.January)
	fmt.Println("December:", goda.December)

	// Get month length
	fmt.Println("Days in February (non-leap):", goda.February.Length(false))
	fmt.Println("Days in February (leap):", goda.February.Length(true))

	// Output:
	// January: January
	// December: December
	// Days in February (non-leap): 28
	// Days in February (leap): 29
}

// ExampleYear_IsLeapYear demonstrates checking for leap years.
func ExampleYear_IsLeapYear() {
	fmt.Println("2024 is leap:", goda.Year(2024).IsLeapYear())
	fmt.Println("2023 is leap:", goda.Year(2023).IsLeapYear())
	fmt.Println("2000 is leap:", goda.Year(2000).IsLeapYear())
	fmt.Println("1900 is leap:", goda.Year(1900).IsLeapYear())

	// Output:
	// 2024 is leap: true
	// 2023 is leap: false
	// 2000 is leap: true
	// 1900 is leap: false
}

// ExampleLocalDate_MarshalJSON demonstrates JSON serialization.
func ExampleLocalDate_MarshalJSON() {
	date := goda.MustNewLocalDate(2024, goda.March, 15)
	jsonBytes, _ := json.Marshal(date)
	fmt.Println(string(jsonBytes))

	// Output:
	// "2024-03-15"
}

// ExampleLocalDate_UnmarshalJSON demonstrates JSON deserialization.
func ExampleLocalDate_UnmarshalJSON() {
	var date goda.LocalDate
	jsonData := []byte(`"2024-03-15"`)
	err := json.Unmarshal(jsonData, &date)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(date)

	// Output:
	// 2024-03-15
}

// ExampleLocalTime_MarshalJSON demonstrates JSON serialization.
func ExampleLocalTime_MarshalJSON() {
	t := goda.MustNewLocalTime(14, 30, 45, 123456789)
	jsonBytes, _ := json.Marshal(t)
	fmt.Println(string(jsonBytes))

	// Output:
	// "14:30:45.123456789"
}

// ExampleLocalTime_UnmarshalJSON demonstrates JSON deserialization.
func ExampleLocalTime_UnmarshalJSON() {
	var t goda.LocalTime
	jsonData := []byte(`"14:30:45.123456789"`)
	err := json.Unmarshal(jsonData, &t)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(t)

	// Output:
	// 14:30:45.123456789
}

// ExampleLocalDateTime demonstrates basic LocalDateTime usage.
func ExampleLocalDateTime() {
	// Create from components
	dt := goda.MustNewLocalDateTimeFromComponents(2024, goda.March, 15, 14, 30, 45, 123456789)
	fmt.Println(dt)

	// Access date and time parts
	fmt.Printf("LocalDate: %s\n", dt.LocalDate())
	fmt.Printf("LocalTime: %s\n", dt.LocalTime())

	// Access individual components
	fmt.Printf("Year: %d, Hour: %d\n", dt.Year(), dt.Hour())

	// Output:
	// 2024-03-15T14:30:45.123456789
	// LocalDate: 2024-03-15
	// LocalTime: 14:30:45.123456789
	// Year: 2024, Hour: 14
}

// ExampleParseLocalDateTime demonstrates parsing a datetime from a string.
func ExampleParseLocalDateTime() {
	dt, err := goda.ParseLocalDateTime("2024-03-15T14:30:45.123456789")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(dt)

	// Output:
	// 2024-03-15T14:30:45.123456789
}

// ExampleMustParseLocalDateTime demonstrates parsing that panics on error.
func ExampleMustParseLocalDateTime() {
	dt := goda.MustParseLocalDateTime("2024-03-15T14:30:45")
	fmt.Println(dt)

	// Output:
	// 2024-03-15T14:30:45
}

// ExampleNewLocalDateTime demonstrates creating a datetime from date and time.
func ExampleNewLocalDateTime() {
	date := goda.MustNewLocalDate(2024, goda.March, 15)
	time := goda.MustNewLocalTime(14, 30, 45, 123456789)
	dt := goda.NewLocalDateTime(date, time)
	fmt.Println(dt)

	// Output:
	// 2024-03-15T14:30:45.123456789
}

// ExampleLocalDateTimeNow demonstrates getting the current datetime.
func ExampleLocalDateTimeNow() {
	now := goda.LocalDateTimeNow()

	// Verify it's valid
	fmt.Printf("Valid: %v\n", !now.IsZero())
	fmt.Printf("Has components: %v\n", now.Year() != 0 && now.Hour() >= 0)

	// Output:
	// Valid: true
	// Has components: true
}

// ExampleLocalDateTime_Compare demonstrates comparing datetimes.
func ExampleLocalDateTime_Compare() {
	dt1 := goda.MustParseLocalDateTime("2024-03-15T14:30:45")
	dt2 := goda.MustParseLocalDateTime("2024-03-15T14:30:45")
	dt3 := goda.MustParseLocalDateTime("2024-03-15T15:30:45")

	fmt.Printf("dt1 == dt2: %v\n", dt1.Compare(dt2) == 0)
	fmt.Printf("dt1 < dt3: %v\n", dt1.IsBefore(dt3))
	fmt.Printf("dt3 > dt1: %v\n", dt3.IsAfter(dt1))

	// Output:
	// dt1 == dt2: true
	// dt1 < dt3: true
	// dt3 > dt1: true
}

// ExampleLocalDateTime_PlusDays demonstrates adding days.
func ExampleLocalDateTime_PlusDays() {
	dt := goda.MustParseLocalDateTime("2024-03-15T14:30:45")
	future := dt.PlusDays(10)
	fmt.Println(future)

	// Output:
	// 2024-03-25T14:30:45
}

// ExampleLocalDateTime_MarshalJSON demonstrates JSON serialization.
func ExampleLocalDateTime_MarshalJSON() {
	dt := goda.MustParseLocalDateTime("2024-03-15T14:30:45.123456789")
	jsonBytes, _ := json.Marshal(dt)
	fmt.Println(string(jsonBytes))

	// Output:
	// "2024-03-15T14:30:45.123456789"
}

// ExampleLocalDateTime_UnmarshalJSON demonstrates JSON deserialization.
func ExampleLocalDateTime_UnmarshalJSON() {
	var dt goda.LocalDateTime
	jsonData := []byte(`"2024-03-15T14:30:45.123456789"`)
	err := json.Unmarshal(jsonData, &dt)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(dt)

	// Output:
	// 2024-03-15T14:30:45.123456789
}

// ExampleField demonstrates working with Field constants.
func ExampleField() {
	// Time fields
	fmt.Println(goda.HourOfDay)
	fmt.Println(goda.MinuteOfHour)
	fmt.Println(goda.SecondOfMinute)
	fmt.Println(goda.NanoOfSecond)

	// Date fields
	fmt.Println(goda.YearField)
	fmt.Println(goda.MonthOfYear)
	fmt.Println(goda.DayOfMonth)
	fmt.Println(goda.DayOfWeekField)

	// Output:
	// HourOfDay
	// MinuteOfHour
	// SecondOfMinute
	// NanoOfSecond
	// Year
	// MonthOfYear
	// DayOfMonth
	// DayOfWeek
}

// ExampleLocalDateTime_IsSupportedField demonstrates checking field support.
func ExampleLocalDateTime_IsSupportedField() {
	dt := goda.MustNewLocalDateTimeFromComponents(2024, goda.March, 15, 14, 30, 45, 0)

	fmt.Printf("Supports HourOfDay: %v\n", dt.IsSupportedField(goda.HourOfDay))
	fmt.Printf("Supports DayOfMonth: %v\n", dt.IsSupportedField(goda.DayOfMonth))
	fmt.Printf("Supports OffsetSeconds: %v\n", dt.IsSupportedField(goda.OffsetSeconds))

	// Output:
	// Supports HourOfDay: true
	// Supports DayOfMonth: true
	// Supports OffsetSeconds: false
}

// ExampleLocalDate_GetFieldInt64 demonstrates getting field values from a date.
func ExampleLocalDate_GetFieldInt64() {
	date := goda.MustNewLocalDate(2024, goda.March, 15) // Friday

	fmt.Printf("Year: %d\n", date.GetFieldInt64(goda.YearField))
	fmt.Printf("Month: %d\n", date.GetFieldInt64(goda.MonthOfYear))
	fmt.Printf("Day: %d\n", date.GetFieldInt64(goda.DayOfMonth))
	fmt.Printf("Day of week: %d\n", date.GetFieldInt64(goda.DayOfWeekField))
	fmt.Printf("Day of year: %d\n", date.GetFieldInt64(goda.DayOfYear))

	// Output:
	// Year: 2024
	// Month: 3
	// Day: 15
	// Day of week: 5
	// Day of year: 75
}

// ExampleLocalTime_GetFieldInt64 demonstrates getting field values from a time.
func ExampleLocalTime_GetFieldInt64() {
	t := goda.MustNewLocalTime(14, 30, 45, 123456789)

	fmt.Printf("Hour: %d\n", t.GetFieldInt64(goda.HourOfDay))
	fmt.Printf("Minute: %d\n", t.GetFieldInt64(goda.MinuteOfHour))
	fmt.Printf("Second: %d\n", t.GetFieldInt64(goda.SecondOfMinute))
	fmt.Printf("Millisecond: %d\n", t.GetFieldInt64(goda.MilliOfSecond))
	fmt.Printf("AM/PM: %d\n", t.GetFieldInt64(goda.AmPmOfDay))

	// Output:
	// Hour: 14
	// Minute: 30
	// Second: 45
	// Millisecond: 123
	// AM/PM: 1
}

// ExampleOffsetDateTime demonstrates basic OffsetDateTime usage.
func ExampleOffsetDateTime() {
	// Create from components with offset
	odt := goda.MustNewOffsetDateTimeFromComponents(2024, goda.March, 15, 14, 30, 45, 123456789, goda.MustNewZoneOffsetHours(9))
	fmt.Println(odt)

	// Access parts
	fmt.Printf("LocalDateTime: %s\n", odt.LocalDateTime())
	fmt.Printf("Offset: %s\n", odt.Offset())
	fmt.Printf("Hour: %d\n", odt.Hour())

	// Output:
	// 2024-03-15T14:30:45.123456789+09:00
	// LocalDateTime: 2024-03-15T14:30:45.123456789
	// Offset: +09:00
	// Hour: 14
}

// ExampleParseOffsetDateTime demonstrates parsing an offset datetime from a string.
func ExampleParseOffsetDateTime() {
	// Parse with positive offset
	odt1, err := goda.ParseOffsetDateTime("2024-03-15T14:30:45.123456789+09:00")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(odt1)

	// Parse with UTC (Z)
	odt2, err := goda.ParseOffsetDateTime("2024-03-15T14:30:45Z")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(odt2)

	// Output:
	// 2024-03-15T14:30:45.123456789+09:00
	// 2024-03-15T14:30:45Z
}

// ExampleMustParseOffsetDateTime demonstrates parsing that panics on error.
func ExampleMustParseOffsetDateTime() {
	odt := goda.MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")
	fmt.Println(odt)

	// Output:
	// 2024-03-15T14:30:45+09:00
}

// ExampleOffsetDateTimeNow demonstrates getting the current offset datetime.
func ExampleOffsetDateTimeNow() {
	now := goda.OffsetDateTimeNow()

	// Verify it's valid
	fmt.Printf("Valid: %v\n", !now.IsZero())
	fmt.Printf("Has offset: %v\n", true)

	// Output:
	// Valid: true
	// Has offset: true
}

// ExampleOffsetDateTimeNowUTC demonstrates getting the current UTC datetime.
func ExampleOffsetDateTimeNowUTC() {
	utc := goda.OffsetDateTimeNowUTC()

	// Verify it's UTC
	fmt.Printf("Valid: %v\n", !utc.IsZero())
	fmt.Printf("Is UTC: %v\n", utc.Offset().IsZero())

	// Output:
	// Valid: true
	// Is UTC: true
}

// ExampleNewOffsetDateTime demonstrates creating an offset datetime.
func ExampleNewOffsetDateTime() {
	date := goda.MustNewLocalDate(2024, goda.March, 15)
	time := goda.MustNewLocalTime(14, 30, 45, 123456789)
	dateTime := goda.NewLocalDateTime(date, time)
	offset := goda.MustNewZoneOffsetHours(9)

	odt := goda.NewOffsetDateTime(dateTime, offset)
	fmt.Println(odt)

	// Output:
	// 2024-03-15T14:30:45.123456789+09:00
}

// ExampleOffsetDateTime_Compare demonstrates comparing offset datetimes.
func ExampleOffsetDateTime_Compare() {
	// Same instant, different offsets
	odt1 := goda.MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")
	odt2 := goda.MustParseOffsetDateTime("2024-03-15T05:30:45Z")

	fmt.Printf("Same instant: %v\n", odt1.IsEqual(odt2))

	// Different instants
	odt3 := goda.MustParseOffsetDateTime("2024-03-15T14:30:46+09:00")
	fmt.Printf("odt1 < odt3: %v\n", odt1.IsBefore(odt3))
	fmt.Printf("odt1 > odt3: %v\n", odt1.IsAfter(odt3))

	// Output:
	// Same instant: true
	// odt1 < odt3: true
	// odt1 > odt3: false
}

// ExampleOffsetDateTime_ToUTC demonstrates converting to UTC.
func ExampleOffsetDateTime_ToUTC() {
	odt := goda.MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")
	utc := odt.ToUTC()

	fmt.Println("Original:", odt)
	fmt.Println("UTC:", utc)
	fmt.Printf("Same instant: %v\n", odt.IsEqual(utc))

	// Output:
	// Original: 2024-03-15T14:30:45+09:00
	// UTC: 2024-03-15T05:30:45Z
	// Same instant: true
}

// ExampleOffsetDateTime_WithOffsetSameLocal demonstrates changing offset keeping local time.
func ExampleOffsetDateTime_WithOffsetSameLocal() {
	odt := goda.MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")
	newOffset := goda.MustNewZoneOffsetHours(-5)
	newOdt := odt.WithOffsetSameLocal(newOffset)

	fmt.Println("Original:", odt)
	fmt.Println("New offset, same local:", newOdt)
	fmt.Printf("Same local time: %v\n", odt.LocalDateTime() == newOdt.LocalDateTime())
	fmt.Printf("Same instant: %v\n", odt.IsEqual(newOdt))

	// Output:
	// Original: 2024-03-15T14:30:45+09:00
	// New offset, same local: 2024-03-15T14:30:45-05:00
	// Same local time: true
	// Same instant: false
}

// ExampleOffsetDateTime_WithOffsetSameInstant demonstrates changing offset keeping the instant.
func ExampleOffsetDateTime_WithOffsetSameInstant() {
	odt := goda.MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")
	newOffset := goda.MustNewZoneOffsetHours(-5)
	newOdt := odt.WithOffsetSameInstant(newOffset)

	fmt.Println("Original:", odt)
	fmt.Println("New offset, same instant:", newOdt)
	fmt.Printf("Same local time: %v\n", odt.LocalDateTime() == newOdt.LocalDateTime())
	fmt.Printf("Same instant: %v\n", odt.IsEqual(newOdt))

	// Output:
	// Original: 2024-03-15T14:30:45+09:00
	// New offset, same instant: 2024-03-15T00:30:45-05:00
	// Same local time: false
	// Same instant: true
}

// ExampleOffsetDateTime_PlusDays demonstrates adding days.
func ExampleOffsetDateTime_PlusDays() {
	odt := goda.MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")
	future := odt.PlusDays(10)
	past := odt.MinusDays(10)

	fmt.Println("Original:", odt)
	fmt.Println("Plus 10 days:", future)
	fmt.Println("Minus 10 days:", past)

	// Output:
	// Original: 2024-03-15T14:30:45+09:00
	// Plus 10 days: 2024-03-25T14:30:45+09:00
	// Minus 10 days: 2024-03-05T14:30:45+09:00
}

// ExampleOffsetDateTime_PlusHours demonstrates adding hours.
func ExampleOffsetDateTime_PlusHours() {
	odt := goda.MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")
	future := odt.PlusHours(5)

	fmt.Println("Original:", odt)
	fmt.Println("Plus 5 hours:", future)

	// Output:
	// Original: 2024-03-15T14:30:45+09:00
	// Plus 5 hours: 2024-03-15T19:30:45+09:00
}

// ExampleOffsetDateTime_MarshalJSON demonstrates JSON serialization.
func ExampleOffsetDateTime_MarshalJSON() {
	odt := goda.MustParseOffsetDateTime("2024-03-15T14:30:45.123456789+09:00")
	jsonBytes, _ := json.Marshal(odt)
	fmt.Println(string(jsonBytes))

	// Output:
	// "2024-03-15T14:30:45.123456789+09:00"
}

// ExampleOffsetDateTime_UnmarshalJSON demonstrates JSON deserialization.
func ExampleOffsetDateTime_UnmarshalJSON() {
	var odt goda.OffsetDateTime
	jsonData := []byte(`"2024-03-15T14:30:45.123456789+09:00"`)
	err := json.Unmarshal(jsonData, &odt)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(odt)

	// Output:
	// 2024-03-15T14:30:45.123456789+09:00
}

// ExampleZoneOffset demonstrates working with zone offsets.
func ExampleZoneOffset() {
	// Create offsets
	utc := goda.ZoneOffsetUTC
	tokyo := goda.MustNewZoneOffsetHours(9)
	india := goda.MustNewZoneOffset(5, 30, 0)
	newYork := goda.MustNewZoneOffsetHours(-5)

	fmt.Println("UTC:", utc)
	fmt.Println("Tokyo:", tokyo)
	fmt.Println("India:", india)
	fmt.Println("New York:", newYork)

	// Output:
	// UTC: Z
	// Tokyo: +09:00
	// India: +05:30
	// New York: -05:00
}

// ExampleParseZoneOffset demonstrates parsing zone offsets.
func ExampleParseZoneOffset() {
	// Parse various formats
	zo1, _ := goda.ParseZoneOffset("Z")
	zo2, _ := goda.ParseZoneOffset("+09:00")
	zo3, _ := goda.ParseZoneOffset("-05:00")
	zo4, _ := goda.ParseZoneOffset("+0530")

	fmt.Println(zo1)
	fmt.Println(zo2)
	fmt.Println(zo3)
	fmt.Println(zo4)

	// Output:
	// Z
	// +09:00
	// -05:00
	// +05:30
}

// ExampleOffsetDateTime_IsSupportedField demonstrates checking field support.
func ExampleOffsetDateTime_IsSupportedField() {
	odt := goda.MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")

	fmt.Printf("Supports HourOfDay: %v\n", odt.IsSupportedField(goda.HourOfDay))
	fmt.Printf("Supports OffsetSeconds: %v\n", odt.IsSupportedField(goda.OffsetSeconds))
	fmt.Printf("Supports InstantSeconds: %v\n", odt.IsSupportedField(goda.InstantSeconds))

	// Output:
	// Supports HourOfDay: true
	// Supports OffsetSeconds: true
	// Supports InstantSeconds: true
}

// ExampleOffsetDateTime_GetFieldInt64 demonstrates getting field values.
func ExampleOffsetDateTime_GetFieldInt64() {
	odt := goda.MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")

	fmt.Printf("Year: %d\n", odt.GetFieldInt64(goda.YearField))
	fmt.Printf("Hour: %d\n", odt.GetFieldInt64(goda.HourOfDay))
	fmt.Printf("Offset seconds: %d\n", odt.GetFieldInt64(goda.OffsetSeconds))

	// Output:
	// Year: 2024
	// Hour: 14
	// Offset seconds: 32400
}
