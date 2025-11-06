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
	fmt.Println("Date:", date)

	// Create a specific time
	timeOfDay := goda.MustNewLocalTime(14, 30, 45, 123456789)
	fmt.Println("Time:", timeOfDay)

	// Get current date and time
	today := goda.LocalDateNow()
	now := goda.LocalTimeNow()
	fmt.Printf("Type of today: %T\n", today)
	fmt.Printf("Type of now: %T\n", now)

	// Output:
	// Date: 2024-03-15
	// Time: 14:30:45.123456789
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
	// Convert from time.Time
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
	// Convert from time.Time
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
	// Time without fractional seconds
	t1 := goda.MustNewLocalTime(14, 30, 45, 0)
	fmt.Println(t1)

	// Time with milliseconds
	t2 := goda.MustNewLocalTime(14, 30, 45, 123000000)
	fmt.Println(t2)

	// Time with microseconds
	t3 := goda.MustNewLocalTime(14, 30, 45, 123456000)
	fmt.Println(t3)

	// Time with nanoseconds
	t4 := goda.MustNewLocalTime(14, 30, 45, 123456789)
	fmt.Println(t4)

	// Time with trailing zeros removed
	t5 := goda.MustNewLocalTime(14, 30, 45, 100000000)
	fmt.Println(t5)

	// Output:
	// 14:30:45
	// 14:30:45.123
	// 14:30:45.123456
	// 14:30:45.123456789
	// 14:30:45.1
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
