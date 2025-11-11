# goda

[![CI](https://github.com/iseki0/goda/workflows/CI/badge.svg)](https://github.com/iseki0/goda/actions?query=workflow%3ACI)
[![Go Reference](https://pkg.go.dev/badge/github.com/iseki0/goda.svg)](https://pkg.go.dev/github.com/iseki0/goda)
[![Go Report Card](https://goreportcard.com/badge/github.com/iseki0/goda)](https://goreportcard.com/report/github.com/iseki0/goda)
[![codecov](https://codecov.io/gh/iseki0/goda/graph/badge.svg?token=TBHUZUY561)](https://codecov.io/gh/iseki0/goda)

> **ThreeTen/JSR-310** model in Go

> [中文版本](README.zh-CN.md)

A Go implementation inspired by Java's `java.time` package (JSR-310), providing immutable date and time types that are **type-safe** and **easy to use**.

## Features

### Core Types

- 📅 **LocalDate**: Date without time (e.g., `2024-03-15`)
- ⏰ **LocalTime**: Time without date (e.g., `14:30:45.123456789`)
- 📆 **LocalDateTime**: Date-time (e.g., `2024-03-15T14:30:45.123456789`)
- 🌍 **OffsetDateTime**: Date-time with UTC offset (e.g., `2024-03-15T14:30:45.123456789+09:00`)
- 🔢 **Field**: Enumeration of date-time fields (like Java's `ChronoField`)

### Key Features

- ✅ **ISO 8601 basic format** support (yyyy-MM-dd, HH:mm:ss[.nnnnnnnnn], combined with 'T')
- ✅ **Java.time compatible formatting**: Fractional seconds aligned to 3-digit boundaries (milliseconds, microseconds, nanoseconds)
- ✅ **Full JSON and SQL** database integration
- ✅ **Date arithmetic**: Add/subtract days, months, years with overflow handling
- ✅ **Field access**: Get any field value (year, month, hour, nano-of-day, etc.)
- ✅ **Zero-copy text marshaling** with `encoding.TextAppender`
- ✅ **Immutable**: All operations return new values
- ✅ **Type-safe**: Compile-time safety with distinct types
- ✅ **Zero-value friendly**: Zero values are properly handled

## Installation

```bash
go get github.com/iseki0/goda
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/iseki0/goda"
)

func main() {
    // Create dates and times
    date := goda.MustNewLocalDate(2024, goda.March, 15)
    time := goda.MustNewLocalTime(14, 30, 45, 123456789)
    datetime := goda.NewLocalDateTime(date, time)
    
    fmt.Println(date)     // 2024-03-15
    fmt.Println(time)     // 14:30:45.123456789
    fmt.Println(datetime) // 2024-03-15T14:30:45.123456789
    
    // Create offset datetime
    offset := goda.MustNewZoneOffsetHours(9)
    offsetDateTime := goda.NewOffsetDateTime(datetime, offset)
    fmt.Println(offsetDateTime) // 2024-03-15T14:30:45.123456789+09:00
    
    // Parse from strings
    date, _ = goda.ParseLocalDate("2024-03-15")
    time = goda.MustParseLocalTime("14:30:45.123456789")
    datetime = goda.MustParseLocalDateTime("2024-03-15T14:30:45")
    offsetDateTime = goda.MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")
    
    // Get current date/time
    today := goda.LocalDateNow()
    now := goda.LocalTimeNow()
    currentDateTime := goda.LocalDateTimeNow()
    currentOffsetDateTime := goda.OffsetDateTimeNow()
    
    // Date arithmetic
    tomorrow := today.PlusDays(1)
    nextMonth := today.PlusMonths(1)
    nextYear := today.PlusYears(1)
    
    // Comparisons
    if tomorrow.IsAfter(today) {
        fmt.Println("Tomorrow is after today!")
    }
}
```

### Combining Date and Time

You can combine LocalDate and LocalTime to create LocalDateTime:

```go
date := goda.MustNewLocalDate(2024, goda.March, 15)
time := goda.MustNewLocalTime(14, 30, 45, 123456789)

// Combine date with time
dateTime := date.AtTime(time)
fmt.Println(dateTime) // 2024-03-15T14:30:45.123456789

// Combine time with date
dateTime2 := time.AtDate(date)
fmt.Println(dateTime2) // 2024-03-15T14:30:45.123456789
```

### Field Access

Access individual date-time fields using the `Field` enumeration:

```go
date := goda.MustNewLocalDate(2024, goda.March, 15)

// Check field support
fmt.Println(date.IsSupportedField(goda.DayOfMonth))  // true
fmt.Println(date.IsSupportedField(goda.HourOfDay))   // false

// Get field values
year := date.GetFieldInt64(goda.YearField)           // 2024
dayOfWeek := date.GetFieldInt64(goda.DayOfWeekField) // 5 (Friday)
dayOfYear := date.GetFieldInt64(goda.DayOfYear)      // 75
epochDays := date.GetFieldInt64(goda.EpochDay)       // Days since Unix epoch

time := goda.MustNewLocalTime(14, 30, 45, 123456789)
hour := time.GetFieldInt64(goda.HourOfDay)           // 14
nanoOfDay := time.GetFieldInt64(goda.NanoOfDay)      // Total nanoseconds since midnight
ampm := time.GetFieldInt64(goda.AmPmOfDay)           // 1 (PM)
```

### Working with OffsetDateTime

```go
// Create with offset
offset := goda.MustNewZoneOffsetHours(9) // Tokyo: UTC+9
odt := goda.MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")

// Convert between offsets
utc := odt.ToUTC()
fmt.Println(utc) // 2024-03-15T05:30:45Z

// Change offset, keeping the same instant
newYork := goda.MustNewZoneOffsetHours(-5)
odtNY := odt.WithOffsetSameInstant(newYork)
fmt.Println(odtNY) // 2024-03-15T00:30:45-05:00

// Compare instants (ignores offset difference)
odt1 := goda.MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")
odt2 := goda.MustParseOffsetDateTime("2024-03-15T05:30:45Z")
fmt.Println(odt1.IsEqual(odt2)) // true (same instant)

// Time arithmetic (adjusts across day boundaries)
later := odt.PlusHours(10) // Adds 10 hours to the instant
```

### JSON Serialization

```go
type Event struct {
    Name      string                `json:"name"`
    Date      goda.LocalDate        `json:"date"`
    Time      goda.LocalTime        `json:"time"`
    CreatedAt goda.LocalDateTime    `json:"created_at"`
    Scheduled goda.OffsetDateTime   `json:"scheduled"`
}

event := Event{
    Name:      "Meeting",
    Date:      goda.MustNewLocalDate(2024, goda.March, 15),
    Time:      goda.MustNewLocalTime(14, 30, 0, 0),
    CreatedAt: goda.MustParseLocalDateTime("2024-03-15T14:30:00"),
    Scheduled: goda.MustParseOffsetDateTime("2024-03-15T14:30:00+09:00"),
}

jsonData, _ := json.Marshal(event)
// {"name":"Meeting","date":"2024-03-15","time":"14:30:00","created_at":"2024-03-15T14:30:00","scheduled":"2024-03-15T14:30:00+09:00"}
```

### Database Integration

```go
type Record struct {
    ID        int64
    CreatedAt goda.LocalDateTime
    Date      goda.LocalDate
    Timestamp goda.OffsetDateTime
}

// Works with database/sql - implements sql.Scanner and driver.Valuer
db.QueryRow("SELECT id, created_at, date, timestamp FROM records WHERE id = ?", 1).Scan(
    &record.ID, &record.CreatedAt, &record.Date, &record.Timestamp,
)
```

## API Overview

### Core Types

| Type | Description | Example |
|------|-------------|---------|
| `LocalDate` | Date without time | `2024-03-15` |
| `LocalTime` | Time without date | `14:30:45.123456789` |
| `LocalDateTime` | Date-time | `2024-03-15T14:30:45` |
| `OffsetDateTime` | Date-time with UTC offset | `2024-03-15T14:30:45+09:00` |
| `ZoneOffset` | UTC offset | `+09:00`, `Z` |
| `Month` | Month of year (1-12) | `March` |
| `Year` | Year | `2024` |
| `DayOfWeek` | Day of week (1=Monday, 7=Sunday) | `Friday` |
| `Field` | Date-time field enumeration | `HourOfDay`, `DayOfMonth` |

### Time Formatting

Time values use ISO 8601 format with **Java.time compatible** fractional second alignment:

| Precision | Digits | Example |
|-----------|--------|---------|
| Whole seconds | 0 | `14:30:45` |
| Milliseconds | 3 | `14:30:45.100`, `14:30:45.123` |
| Microseconds | 6 | `14:30:45.123400`, `14:30:45.123456` |
| Nanoseconds | 9 | `14:30:45.000000001`, `14:30:45.123456789` |

Fractional seconds are automatically aligned to 3-digit boundaries (milliseconds, microseconds, nanoseconds), matching Java's `LocalTime` behavior. Parsing accepts any length of fractional seconds.

### Field Constants (30 fields)

**Time Fields**: `NanoOfSecond`, `NanoOfDay`, `MicroOfSecond`, `MicroOfDay`, `MilliOfSecond`, `MilliOfDay`, `SecondOfMinute`, `SecondOfDay`, `MinuteOfHour`, `MinuteOfDay`, `HourOfAmPm`, `ClockHourOfAmPm`, `HourOfDay`, `ClockHourOfDay`, `AmPmOfDay`

**Date Fields**: `DayOfWeekField`, `DayOfMonth`, `DayOfYear`, `EpochDay`, `AlignedDayOfWeekInMonth`, `AlignedDayOfWeekInYear`, `AlignedWeekOfMonth`, `AlignedWeekOfYear`, `MonthOfYear`, `ProlepticMonth`, `YearOfEra`, `YearField`, `Era`

**Other Fields**: `InstantSeconds`, `OffsetSeconds`

### Implemented Interfaces

All types implement:
- `fmt.Stringer`
- `encoding.TextMarshaler` / `encoding.TextUnmarshaler`
- `encoding.TextAppender` (zero-copy text marshaling)
- `json.Marshaler` / `json.Unmarshaler`
- `sql.Scanner` / `driver.Valuer`

## Design Philosophy

This package follows the **ThreeTen/JSR-310** model (Java's `java.time` package), providing date and time types that are:

- **Immutable**: All operations return new values
- **Type-safe**: Distinct types for date, time, and datetime
- **Simple formats**: Uses ISO 8601 basic formats (not the full complex specification)
- **Database-friendly**: Direct SQL integration
- **Field-based access**: Universal field access pattern via `GetFieldInt64`
- **Zero-value safe**: Zero values are properly handled throughout

### When to Use Each Type

**LocalDate, LocalTime, LocalDateTime**

Use the local types when you only need the date/time without timezone information:
- **Birthdays**: "March 15" means March 15 everywhere
- **Business hours**: "9:00 AM - 5:00 PM" in local context
- **Schedules**: "Meeting at 2:30 PM" without timezone concerns
- **Calendar dates**: Historical dates, recurring events

**OffsetDateTime**

Use OffsetDateTime when you need to represent an instant with a specific UTC offset:
- **API responses**: Timestamps with timezone information
- **Scheduled events**: Events that occur at a specific instant (e.g., "2024-03-15T14:30:00+09:00")
- **Database timestamps**: When storing timestamps with offset information
- **International coordination**: When you need to know both the local time and UTC offset

For full timezone-aware operations with DST handling, use `ZonedDateTime` (coming soon).

## Documentation

Full API documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/iseki0/goda).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
