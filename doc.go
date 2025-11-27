// Package goda provides date and time types following the ThreeTen/JSR-310 model (java.time package).
//
// This package implements the following main types:
//
//   - LocalDate: A date without time (e.g., 2024-03-15)
//   - LocalTime: A time without date (e.g., 14:30:45.123456789)
//   - LocalDateTime: A date-time without timezone (e.g., 2024-03-15T14:30:45.123456789)
//   - ZoneOffset: A time-zone offset from UTC (e.g., +08:00, -05:00, Z)
//   - OffsetDateTime: A date-time with offset from UTC (e.g., 2024-03-15T14:30:45+08:00)
//   - Year, Month, DayOfWeek: Supporting types for date/time operations
//
// All types implement standard interfaces for serialization:
//   - encoding.TextMarshaler and encoding.TextUnmarshaler (ISO 8601 basic format)
//   - encoding.json.Marshaler and encoding.json.Unmarshaler
//   - database/sql.Scanner and database/sql/driver.Valuer
//
// Note: This package uses ISO 8601 basic formats only (yyyy-MM-dd, HH:mm:ss[.nnnnnnnnn]),
// not the full complex ISO 8601 specification (no week dates, ordinal dates, or timezone offsets).
//
// # Quick Start
//
// See the Example function for comprehensive usage examples.
//
// # Design Philosophy
//
// This package follows the ThreeTen/JSR-310 design:
//
//   - Immutable: All operations return new instances
//   - Type-safe: Strong typing prevents mixing dates and times
//   - Zero-value safe: Zero values are clearly invalid and IsZero() returns true
//
// # Comparison with time.Time
//
// Go's time.Time combines date, time, and location. This package separates concerns:
//
//   - Use LocalDate when you only need a date (e.g., birthdays, deadlines)
//   - Use LocalTime when you only need a time (e.g., office hours, schedules)
//   - Convert to/from time.Time when timezone information is needed
//
// # Format Specification
//
// This package uses ISO 8601 basic calendar date and time formats (not the full specification):
//
//   - LocalDate: yyyy-MM-dd (e.g., "2024-03-15")
//     Only Gregorian calendar dates. No week dates (YYYY-Www-D) or ordinal dates (YYYY-DDD).
//
//   - LocalTime: HH:mm:ss[.nnnnnnnnn] (e.g., "14:30:45.123456789")
//     24-hour format. Fractional seconds up to nanoseconds.
//     Fractional seconds are aligned to 3-digit boundaries (milliseconds, microseconds, nanoseconds)
//     for Java.time compatibility: 100ms → "14:30:45.100", 123.4ms → "14:30:45.123400".
//     Parsing accepts any length of fractional seconds (e.g., "14:30:45.1" → 100ms).
//
//   - LocalDateTime: yyyy-MM-ddTHH:mm:ss[.nnnnnnnnn] (e.g., "2024-03-15T14:30:45.123456789")
//     Combined with 'T' separator (lowercase 't' accepted when parsing).
//
//   - ZoneOffset: ±HH:mm[:ss] or Z for UTC (e.g., "+08:00", "-05:30", "Z")
//     Hours must be in range [-18, 18], minutes and seconds in [0, 59].
//     Compact formats (±HH, ±HHMM, ±HHMMSS) are also supported.
//
//   - OffsetDateTime: yyyy-MM-ddTHH:mm:ss[.nnnnnnnnn]±HH:mm[:ss] (e.g., "2024-03-15T14:30:45+08:00")
//     Combines LocalDateTime and ZoneOffset. 'Z' is accepted as UTC offset.
package goda
