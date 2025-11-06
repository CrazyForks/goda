// Package goda provides date and time types without timezone information,
// following the ThreeTen/JSR-310 model (java.time package).
//
// This package implements four main types:
//
//   - LocalDate: A date without time or timezone (e.g., 2024-03-15)
//   - LocalTime: A time without date or timezone (e.g., 14:30:45.123456789)
//   - LocalDateTime: A date-time without timezone (e.g., 2024-03-15T14:30:45.123456789)
//   - Year, Month, DayOfWeek: Supporting types for date/time operations
//
// All types implement standard interfaces for serialization:
//   - encoding.TextMarshaler and encoding.TextUnmarshaler (ISO 8601 format)
//   - encoding.json.Marshaler and encoding.json.Unmarshaler
//   - database/sql.Scanner and database/sql/driver.Valuer
//
// # Quick Start
//
// Creating dates and times:
//
//	// Specific date and time
//	date := goda.MustNewLocalDate(2024, goda.March, 15)
//	time := goda.MustNewLocalTime(14, 30, 45, 0)
//	datetime := goda.NewLocalDateTime(date, time)
//
//	// Current date and time
//	today := goda.LocalDateNow()
//	now := goda.LocalTimeNow()
//	currentDateTime := goda.LocalDateTimeNow()
//
//	// Parse from string
//	date = goda.MustParseLocalDate("2024-03-15")
//	time = goda.MustParseLocalTime("14:30:45.123456789")
//	datetime = goda.MustParseLocalDateTime("2024-03-15T14:30:45.123456789")
//
// LocalDate arithmetic:
//
//	tomorrow := today.PlusDays(1)
//	nextMonth := today.PlusMonths(1)
//	nextYear := today.PlusYears(1)
//
// Comparisons:
//
//	if date1.IsBefore(date2) {
//	    fmt.Println("date1 is earlier")
//	}
//
// Serialization:
//
//	// JSON
//	jsonBytes, _ := json.Marshal(date)  // "2024-03-15"
//	json.Unmarshal(jsonBytes, &date)
//
//	// String format (ISO 8601)
//	str := date.String()  // "2024-03-15"
//	str := time.String()  // "14:30:45.123456789"
//
// # Design Philosophy
//
// This package follows the ThreeTen/JSR-310 design:
//
//   - Immutable: All operations return new instances
//   - Type-safe: Strong typing prevents mixing dates and times
//   - Zero-value safe: Zero values are clearly invalid and IsZero() returns true
//   - Timezone-free: These types represent local date/time without timezone information
//
// # Comparison with time.Time
//
// Go's time.Time combines date, time, and location. This package separates concerns:
//
//   - Use LocalDate when you only need a date (e.g., birthdays, deadlines)
//   - Use LocalTime when you only need a time (e.g., office hours, schedules)
//   - Convert to/from time.Time when timezone information is needed
//
// # ISO 8601 Format
//
// All types use ISO 8601 standard formats:
//
//   - LocalDate: YYYY-MM-DD (e.g., "2024-03-15")
//   - LocalTime: HH:MM:SS[.nnnnnnnnn] (e.g., "14:30:45.123456789")
//
// Fractional seconds in LocalTime automatically trim trailing zeros while
// maintaining precision (following Java LocalTime behavior).
package goda
