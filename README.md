# goda

[![CI](https://github.com/iseki0/goda/workflows/CI/badge.svg)](https://github.com/iseki0/goda/actions?query=workflow%3ACI)
[![Go Reference](https://pkg.go.dev/badge/github.com/iseki0/goda.svg)](https://pkg.go.dev/github.com/iseki0/goda)
[![Go Report Card](https://goreportcard.com/badge/github.com/iseki0/goda)](https://goreportcard.com/report/github.com/iseki0/goda)
[![codecov](https://codecov.io/gh/iseki0/goda/branch/main/graph/badge.svg)](https://codecov.io/gh/iseki0/goda)

> **ThreeTen/JSR-310** model in Go - Date and time types without timezone information.

A Go implementation inspired by Java's `java.time` package (JSR-310), providing immutable date and time types that are **timezone-free**, **type-safe**, and **easy to use**.

## Features

### Core Types

- 📅 **LocalDate**: Date without time or timezone (e.g., `2024-03-15`)
- ⏰ **LocalTime**: Time without date or timezone (e.g., `14:30:45.123456789`)
- 📆 **LocalDateTime**: Date-time without timezone (e.g., `2024-03-15T14:30:45.123456789`)
- 🔢 **Field**: Enumeration of date-time fields (like Java's `ChronoField`)

### Key Features

- ✅ **ISO 8601 basic format** support (yyyy-MM-dd, HH:mm:ss[.nnnnnnnnn], combined with 'T')
- ✅ **Smart fractional seconds**: Automatically trims trailing zeros (14:30:45.1 instead of 14:30:45.100)
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
    
    // Parse from strings
    date, _ = goda.ParseLocalDate("2024-03-15")
    time = goda.MustParseLocalTime("14:30:45.123456789")
    datetime = goda.MustParseLocalDateTime("2024-03-15T14:30:45")
    
    // Get current date/time
    today := goda.LocalDateNow()
    now := goda.LocalTimeNow()
    currentDateTime := goda.LocalDateTimeNow()
    
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

### JSON Serialization

```go
type Event struct {
    Name string             `json:"name"`
    Date goda.LocalDate     `json:"date"`
    Time goda.LocalTime     `json:"time"`
}

event := Event{
    Name: "Meeting",
    Date: goda.MustNewLocalDate(2024, goda.March, 15),
    Time: goda.MustNewLocalTime(14, 30, 0, 0),
}

jsonData, _ := json.Marshal(event)
// {"name":"Meeting","date":"2024-03-15","time":"14:30:00"}
```

### Database Integration

```go
type Record struct {
    ID        int64
    CreatedAt goda.LocalDateTime
    Date      goda.LocalDate
}

// Works with database/sql - implements sql.Scanner and driver.Valuer
db.QueryRow("SELECT id, created_at, date FROM records WHERE id = ?", 1).Scan(
    &record.ID, &record.CreatedAt, &record.Date,
)
```

## API Overview

### Core Types

| Type | Description | Example |
|------|-------------|---------|
| `LocalDate` | Date without time/timezone | `2024-03-15` |
| `LocalTime` | Time without date/timezone | `14:30:45.123456789` |
| `LocalDateTime` | Date-time without timezone | `2024-03-15T14:30:45` |
| `Month` | Month of year (1-12) | `March` |
| `Year` | Year | `2024` |
| `DayOfWeek` | Day of week (1=Monday, 7=Sunday) | `Friday` |
| `Field` | Date-time field enumeration | `HourOfDay`, `DayOfMonth` |

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
- **Timezone-free**: No confusion about timezones
- **Simple formats**: Uses ISO 8601 basic formats (not the full complex specification)
- **Database-friendly**: Direct SQL integration
- **Field-based access**: Universal field access pattern via `GetFieldInt64`
- **Zero-value safe**: Zero values are properly handled throughout

### Why Timezone-Free?

Many applications deal with dates and times that are independent of timezones:
- **Birthdays**: "March 15" means March 15 everywhere
- **Business hours**: "9:00 AM - 5:00 PM" in local context
- **Schedules**: "Meeting at 2:30 PM" without timezone concerns
- **Calendar dates**: Historical dates, recurring events

## Documentation

Full API documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/iseki0/goda).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
