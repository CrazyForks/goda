# goda

[![CI](https://github.com/iseki0/goda/workflows/CI/badge.svg)](https://github.com/iseki0/goda/actions?query=workflow%3ACI)
[![Go Reference](https://pkg.go.dev/badge/github.com/iseki0/goda.svg)](https://pkg.go.dev/github.com/iseki0/goda)
[![Go Report Card](https://goreportcard.com/badge/github.com/iseki0/goda)](https://goreportcard.com/report/github.com/iseki0/goda)
[![codecov](https://codecov.io/gh/iseki0/goda/branch/main/graph/badge.svg)](https://codecov.io/gh/iseki0/goda)

ThreeTen model in Go - Date and time types without timezone information.

## Features

- **LocalDate**: Date without time or timezone (e.g., `2024-03-15`)
- **LocalTime**: Time without date or timezone (e.g., `14:30:45.123456789`)
- **LocalDateTime**: Date-time without timezone (e.g., `2024-03-15T14:30:45.123456789`)
- **ISO 8601** format support
- Full **JSON** and **SQL** database integration
- Comprehensive **date arithmetic** (add/subtract days, months, years)
- **Zero-copy** text marshaling with `encoding.TextAppender`

## Installation

```bash
go get github.com/iseki0/goda
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/iseki0/goda"
)

func main() {
    // Create dates and times
    date := goda.MustNewLocalDate(2024, goda.March, 15)
    time := goda.MustNewLocalTime(14, 30, 45, 0)
    datetime := goda.NewLocalDateTime(date, time)
    
    fmt.Println(date)     // 2024-03-15
    fmt.Println(time)     // 14:30:45
    fmt.Println(datetime) // 2024-03-15T14:30:45
    
    // Parse from strings
    date = goda.MustParseLocalDate("2024-03-15")
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

## Documentation

Full API documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/iseki0/goda).

## Design Philosophy

This package follows the **ThreeTen/JSR-310** model (Java's `java.time` package), providing date and time types that are:

- **Immutable**: All operations return new values
- **Type-safe**: Distinct types for date, time, and datetime
- **Timezone-free**: No confusion about timezones
- **Standard**: ISO 8601 format support
- **Database-friendly**: Direct SQL integration

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
