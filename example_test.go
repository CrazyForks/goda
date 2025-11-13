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

// ExampleLocalDate_GetField demonstrates querying date fields with TemporalValue.
func ExampleLocalDate_GetField() {
	date := goda.MustNewLocalDate(2024, goda.March, 15) // Friday

	// Query various date fields
	dayOfWeek := date.GetField(goda.DayOfWeekField)
	if dayOfWeek.Valid() {
		fmt.Printf("Day of week: %d (1=Monday, 7=Sunday)\n", dayOfWeek.Int())
	}

	month := date.GetField(goda.MonthOfYear)
	if month.Valid() {
		fmt.Printf("Month: %d\n", month.Int())
	}

	year := date.GetField(goda.YearField)
	if year.Valid() {
		fmt.Printf("Year: %d\n", year.Int())
	}

	// Query unsupported field (time field on date)
	hour := date.GetField(goda.HourOfDay)
	if hour.Unsupported() {
		fmt.Println("Hour field is not supported for LocalDate")
	}

	// Output:
	// Day of week: 5 (1=Monday, 7=Sunday)
	// Month: 3
	// Year: 2024
	// Hour field is not supported for LocalDate
}

// ExampleLocalDate_GetField_advancedFields demonstrates advanced date field queries.
func ExampleLocalDate_GetField_advancedFields() {
	date := goda.MustNewLocalDate(2024, goda.March, 15)

	// Day of year (1-366)
	dayOfYear := date.GetField(goda.DayOfYear)
	fmt.Printf("Day of year: %d\n", dayOfYear.Int())

	// Epoch days (days since 1970-01-01)
	epochDay := date.GetField(goda.EpochDay)
	fmt.Printf("Days since Unix epoch: %d\n", epochDay.Int64())

	// Proleptic month (months since year 0)
	prolepticMonth := date.GetField(goda.ProlepticMonth)
	fmt.Printf("Proleptic month: %d\n", prolepticMonth.Int64())

	// Era (0=BCE, 1=CE)
	era := date.GetField(goda.Era)
	if era.Int() == 1 {
		fmt.Println("Era: CE (Common Era)")
	}

	// Output:
	// Day of year: 75
	// Days since Unix epoch: 19797
	// Proleptic month: 24290
	// Era: CE (Common Era)
}

// ExampleLocalTime_GetField demonstrates querying time fields with TemporalValue.
func ExampleLocalTime_GetField() {
	t := goda.MustNewLocalTime(14, 30, 45, 123456789)

	// Query various time fields
	hour := t.GetField(goda.HourOfDay)
	if hour.Valid() {
		fmt.Printf("Hour of day (0-23): %d\n", hour.Int())
	}

	minute := t.GetField(goda.MinuteOfHour)
	if minute.Valid() {
		fmt.Printf("Minute of hour: %d\n", minute.Int())
	}

	second := t.GetField(goda.SecondOfMinute)
	if second.Valid() {
		fmt.Printf("Second of minute: %d\n", second.Int())
	}

	nanos := t.GetField(goda.NanoOfSecond)
	if nanos.Valid() {
		fmt.Printf("Nanoseconds: %d\n", nanos.Int())
	}

	// Query unsupported field (date field on time)
	dayOfMonth := t.GetField(goda.DayOfMonth)
	if dayOfMonth.Unsupported() {
		fmt.Println("DayOfMonth field is not supported for LocalTime")
	}

	// Output:
	// Hour of day (0-23): 14
	// Minute of hour: 30
	// Second of minute: 45
	// Nanoseconds: 123456789
	// DayOfMonth field is not supported for LocalTime
}

// ExampleLocalTime_GetField_clockHours demonstrates 12-hour clock field queries.
func ExampleLocalTime_GetField_clockHours() {
	// Afternoon time (2:30 PM)
	afternoon := goda.MustNewLocalTime(14, 30, 0, 0)

	// 24-hour format
	hourOfDay := afternoon.GetField(goda.HourOfDay)
	fmt.Printf("24-hour format: %d:30\n", hourOfDay.Int())

	// 12-hour format components
	hourOfAmPm := afternoon.GetField(goda.HourOfAmPm)
	amPm := afternoon.GetField(goda.AmPmOfDay)
	amPmStr := "AM"
	if amPm.Int() == 1 {
		amPmStr = "PM"
	}
	fmt.Printf("12-hour format: %d:30 %s\n", hourOfAmPm.Int(), amPmStr)

	// Clock hour (1-12 instead of 0-11)
	clockHour := afternoon.GetField(goda.ClockHourOfAmPm)
	fmt.Printf("Clock hour (1-12): %d:30 %s\n", clockHour.Int(), amPmStr)

	// Midnight special case
	midnight := goda.MustNewLocalTime(0, 0, 0, 0)
	midnightClock := midnight.GetField(goda.ClockHourOfDay)
	fmt.Printf("Midnight clock hour: %d:00\n", midnightClock.Int())

	// Output:
	// 24-hour format: 14:30
	// 12-hour format: 2:30 PM
	// Clock hour (1-12): 2:30 PM
	// Midnight clock hour: 24:00
}

// ExampleLocalTime_GetField_ofDayFields demonstrates querying cumulative daily values.
func ExampleLocalTime_GetField_ofDayFields() {
	t := goda.MustNewLocalTime(14, 30, 45, 500000000) // 2:30:45.5 PM

	// Total seconds elapsed since midnight
	secondOfDay := t.GetField(goda.SecondOfDay)
	fmt.Printf("Seconds since midnight: %d\n", secondOfDay.Int())

	// Total minutes elapsed since midnight
	minuteOfDay := t.GetField(goda.MinuteOfDay)
	fmt.Printf("Minutes since midnight: %d\n", minuteOfDay.Int())

	// Total milliseconds elapsed since midnight
	milliOfDay := t.GetField(goda.MilliOfDay)
	fmt.Printf("Milliseconds since midnight: %d\n", milliOfDay.Int64())

	// Total nanoseconds elapsed since midnight
	nanoOfDay := t.GetField(goda.NanoOfDay)
	fmt.Printf("Nanoseconds since midnight: %d\n", nanoOfDay.Int64())

	// Output:
	// Seconds since midnight: 52245
	// Minutes since midnight: 870
	// Milliseconds since midnight: 52245500
	// Nanoseconds since midnight: 52245500000000
}

// ExampleLocalDateTime_GetField demonstrates querying fields from a date-time.
func ExampleLocalDateTime_GetField() {
	dt := goda.MustNewLocalDateTimeFromComponents(2024, goda.March, 15, 14, 30, 45, 123456789)

	// Query date fields
	year := dt.GetField(goda.YearField)
	month := dt.GetField(goda.MonthOfYear)
	day := dt.GetField(goda.DayOfMonth)

	if year.Valid() && month.Valid() && day.Valid() {
		fmt.Printf("Date: %04d-%02d-%02d\n", year.Int(), month.Int(), day.Int())
	}

	// Query time fields
	hour := dt.GetField(goda.HourOfDay)
	minute := dt.GetField(goda.MinuteOfHour)
	second := dt.GetField(goda.SecondOfMinute)

	if hour.Valid() && minute.Valid() && second.Valid() {
		fmt.Printf("Time: %02d:%02d:%02d\n", hour.Int(), minute.Int(), second.Int())
	}

	// Query day of week
	dayOfWeek := dt.GetField(goda.DayOfWeekField)
	if dayOfWeek.Valid() {
		fmt.Printf("Day of week: %d (Friday)\n", dayOfWeek.Int())
	}

	// Output:
	// Date: 2024-03-15
	// Time: 14:30:45
	// Day of week: 5 (Friday)
}

// ExampleLocalDateTime_GetField_delegation demonstrates field delegation.
func ExampleLocalDateTime_GetField_delegation() {
	dt := goda.MustNewLocalDateTimeFromComponents(2024, goda.March, 15, 14, 30, 45, 0)

	// LocalDateTime delegates date fields to LocalDate
	dayOfYear := dt.GetField(goda.DayOfYear)
	fmt.Printf("Day of year: %d\n", dayOfYear.Int())

	// LocalDateTime delegates time fields to LocalTime
	nanoOfDay := dt.GetField(goda.NanoOfDay)
	fmt.Printf("Nanoseconds of day: %d\n", nanoOfDay.Int64())

	// Unsupported fields return unsupported TemporalValue
	offsetSeconds := dt.GetField(goda.OffsetSeconds)
	if offsetSeconds.Unsupported() {
		fmt.Println("OffsetSeconds is not supported for LocalDateTime")
	}

	// Output:
	// Day of year: 75
	// Nanoseconds of day: 52245000000000
	// OffsetSeconds is not supported for LocalDateTime
}

// ExampleNewDurationOfSeconds demonstrates creating a Duration from seconds and nanoseconds.
func ExampleNewDurationOfSeconds() {
	// Simple duration: 10.5 seconds
	d := goda.NewDurationOfSeconds(10, 500000000)
	fmt.Println(d)

	// Zero duration
	zero := goda.NewDurationOfSeconds(0, 0)
	fmt.Println("Zero:", zero)

	// Negative duration
	neg := goda.NewDurationOfSeconds(-5, 0)
	fmt.Println("Negative:", neg)

	// Output:
	// PT10.5S
	// Zero: PT0S
	// Negative: PT-5S
}

// ExampleNewDurationOfSeconds_overflow demonstrates overflow wrapping behavior.
// Unlike Java's java.time which throws exceptions, this implementation wraps around
// following Go's convention.
func ExampleNewDurationOfSeconds_overflow() {
	// Nanoseconds overflow: 2 billion nanos = 2 seconds, wraps to seconds
	d1 := goda.NewDurationOfSeconds(5, 2000000000)
	fmt.Println("5s + 2bn nanos:", d1)

	// Negative nanoseconds: 5s - 1.5bn nanos = 3.5s
	d2 := goda.NewDurationOfSeconds(5, -1500000000)
	fmt.Println("5s - 1.5bn nanos:", d2)

	// Large overflow: 10 billion nanos wraps to 10 seconds
	d3 := goda.NewDurationOfSeconds(0, 10000000000)
	fmt.Println("10bn nanos:", d3)

	// Output:
	// 5s + 2bn nanos: PT7S
	// 5s - 1.5bn nanos: PT3.5S
	// 10bn nanos: PT10S
}

// ExampleNewDurationByGoDuration demonstrates converting from time.Duration.
func ExampleNewDurationByGoDuration() {
	// From Go's time.Duration
	goDuration := 5*time.Second + 500*time.Millisecond
	d := goda.NewDurationByGoDuration(goDuration)
	fmt.Println(d)

	// From minutes
	minutes := 90 * time.Minute
	d2 := goda.NewDurationByGoDuration(minutes)
	fmt.Println(d2)

	// Output:
	// PT5.5S
	// PT1H30M
}

// ExampleParseDuration demonstrates parsing ISO-8601 duration strings.
func ExampleParseDuration() {
	// Parse various formats
	d1, _ := goda.ParseDuration("PT1H30M")
	fmt.Println("1h 30m:", d1)

	d2, _ := goda.ParseDuration("PT45.5S")
	fmt.Println("45.5s:", d2)

	d3, _ := goda.ParseDuration("PT8H6M12.345S")
	fmt.Println("Complex:", d3)

	// Negative duration
	d4, _ := goda.ParseDuration("PT-2H30M")
	fmt.Println("Negative:", d4)

	// Output:
	// 1h 30m: PT1H30M
	// 45.5s: PT45.5S
	// Complex: PT8H6M12.345S
	// Negative: PT-2H30M
}

// ExampleMustParseDuration demonstrates parsing that panics on error.
func ExampleMustParseDuration() {
	d := goda.MustParseDuration("PT2H30M45S")
	fmt.Println(d)

	// Output:
	// PT2H30M45S
}

// ExampleDuration_Plus demonstrates adding durations with overflow wrapping.
func ExampleDuration_Plus() {
	d1 := goda.NewDurationOfSeconds(3600, 0) // 1 hour
	d2 := goda.NewDurationOfSeconds(1800, 0) // 30 minutes

	result := d1.Plus(d2)
	fmt.Println("1h + 30m:", result)

	// Adding with nanosecond overflow
	d3 := goda.NewDurationOfSeconds(5, 700000000) // 5.7s
	d4 := goda.NewDurationOfSeconds(3, 500000000) // 3.5s
	result2 := d3.Plus(d4)
	fmt.Println("5.7s + 3.5s:", result2)

	// Output:
	// 1h + 30m: PT1H30M
	// 5.7s + 3.5s: PT9.2S
}

// ExampleDuration_Minus demonstrates subtracting durations.
func ExampleDuration_Minus() {
	d1 := goda.NewDurationOfSeconds(3600, 0) // 1 hour
	d2 := goda.NewDurationOfSeconds(1800, 0) // 30 minutes

	result := d1.Minus(d2)
	fmt.Println("1h - 30m:", result)

	// Subtraction resulting in negative
	d3 := goda.NewDurationOfSeconds(1000, 0)
	d4 := goda.NewDurationOfSeconds(2000, 0)
	result2 := d3.Minus(d4)
	fmt.Println("1000s - 2000s:", result2)

	// Output:
	// 1h - 30m: PT30M
	// 1000s - 2000s: PT-16M40S
}

// ExampleDuration_Negated demonstrates negating a duration.
func ExampleDuration_Negated() {
	d := goda.NewDurationOfSeconds(3600, 500000000) // 1h 0.5s
	neg := d.Negated()
	fmt.Println("Original:", d)
	fmt.Println("Negated:", neg)

	// Double negation
	doubleNeg := neg.Negated()
	fmt.Println("Double negated:", doubleNeg)

	// Output:
	// Original: PT1H0.5S
	// Negated: PT-1H0.5S
	// Double negated: PT1H0.5S
}

// ExampleDuration_Abs demonstrates getting absolute value.
func ExampleDuration_Abs() {
	positive := goda.NewDurationOfSeconds(5, 0)
	negative := goda.NewDurationOfSeconds(-5, 0)

	fmt.Println("Abs of positive:", positive.Abs())
	fmt.Println("Abs of negative:", negative.Abs())

	// Output:
	// Abs of positive: PT5S
	// Abs of negative: PT5S
}

// ExampleDuration_IsZero demonstrates checking for zero duration.
func ExampleDuration_IsZero() {
	zero := goda.NewDurationOfSeconds(0, 0)
	nonZero := goda.NewDurationOfSeconds(1, 0)

	fmt.Println("Zero duration:", zero.IsZero())
	fmt.Println("Non-zero duration:", nonZero.IsZero())

	// Output:
	// Zero duration: true
	// Non-zero duration: false
}

// ExampleDuration_IsPositive demonstrates checking for positive duration.
func ExampleDuration_IsPositive() {
	positive := goda.NewDurationOfSeconds(1, 0)
	negative := goda.NewDurationOfSeconds(-1, 0)
	zero := goda.NewDurationOfSeconds(0, 0)

	fmt.Println("Positive:", positive.IsPositive())
	fmt.Println("Negative:", negative.IsPositive())
	fmt.Println("Zero:", zero.IsPositive())

	// Output:
	// Positive: true
	// Negative: false
	// Zero: false
}

// ExampleDuration_IsNegative demonstrates checking for negative duration.
func ExampleDuration_IsNegative() {
	positive := goda.NewDurationOfSeconds(1, 0)
	negative := goda.NewDurationOfSeconds(-1, 0)
	zero := goda.NewDurationOfSeconds(0, 0)

	fmt.Println("Positive:", positive.IsNegative())
	fmt.Println("Negative:", negative.IsNegative())
	fmt.Println("Zero:", zero.IsNegative())

	// Output:
	// Positive: false
	// Negative: true
	// Zero: false
}

// ExampleDuration_Compare demonstrates comparing durations.
func ExampleDuration_Compare() {
	d1 := goda.NewDurationOfSeconds(3600, 0) // 1 hour
	d2 := goda.NewDurationOfSeconds(1800, 0) // 30 minutes
	d3 := goda.NewDurationOfSeconds(3600, 0) // 1 hour

	fmt.Println("d1 > d2:", d1.Compare(d2) > 0)
	fmt.Println("d1 < d2:", d1.Compare(d2) < 0)
	fmt.Println("d1 == d3:", d1.Compare(d3) == 0)

	// Output:
	// d1 > d2: true
	// d1 < d2: false
	// d1 == d3: true
}

// ExampleDuration_MarshalJSON demonstrates JSON serialization.
func ExampleDuration_MarshalJSON() {
	d := goda.NewDurationOfSeconds(3661, 500000000)
	jsonBytes, _ := json.Marshal(d)
	fmt.Println(string(jsonBytes))

	// Output:
	// "PT1H1M1.5S"
}

// ExampleDuration_UnmarshalJSON demonstrates JSON deserialization.
func ExampleDuration_UnmarshalJSON() {
	var d goda.Duration
	jsonData := []byte(`"PT1H30M45S"`)
	err := json.Unmarshal(jsonData, &d)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(d)

	// Output:
	// PT1H30M45S
}

// ExampleLocalTime_Between demonstrates calculating duration between times.
func ExampleLocalTime_Between() {
	t1 := goda.MustNewLocalTime(10, 0, 0, 0)
	t2 := goda.MustNewLocalTime(14, 30, 0, 0)

	duration := t1.Between(t2)
	fmt.Println("From 10:00 to 14:30:", duration)

	// Negative duration (going backwards)
	duration2 := t2.Between(t1)
	fmt.Println("From 14:30 to 10:00:", duration2)

	// With nanoseconds
	t3 := goda.MustNewLocalTime(10, 0, 0, 500000000)
	t4 := goda.MustNewLocalTime(10, 0, 2, 250000000)
	duration3 := t3.Between(t4)
	fmt.Println("Precise duration:", duration3)

	// Output:
	// From 10:00 to 14:30: PT4H30M
	// From 14:30 to 10:00: PT-4H30M
	// Precise duration: PT1.75S
}

// ExampleLocalDateTime_Between demonstrates calculating duration between date-times.
func ExampleLocalDateTime_Between() {
	dt1 := goda.MustNewLocalDateTimeFromComponents(2024, goda.March, 15, 10, 0, 0, 0)
	dt2 := goda.MustNewLocalDateTimeFromComponents(2024, goda.March, 15, 14, 30, 0, 0)

	// Same day
	duration := dt1.Between(dt2)
	fmt.Println("Same day, 4.5 hours:", duration)

	// Different days
	dt3 := goda.MustNewLocalDateTimeFromComponents(2024, goda.March, 16, 10, 0, 0, 0)
	duration2 := dt1.Between(dt3)
	fmt.Println("Next day, 24 hours:", duration2)

	// Complex example: 1 day and 5.5 hours
	dt4 := goda.MustNewLocalDateTimeFromComponents(2024, goda.March, 15, 10, 0, 0, 0)
	dt5 := goda.MustNewLocalDateTimeFromComponents(2024, goda.March, 16, 15, 30, 0, 0)
	duration3 := dt4.Between(dt5)
	fmt.Println("1 day + 5.5 hours:", duration3)

	// Output:
	// Same day, 4.5 hours: PT4H30M
	// Next day, 24 hours: PT24H
	// 1 day + 5.5 hours: PT29H30M
}
