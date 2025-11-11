# goda

[![CI](https://github.com/iseki0/goda/workflows/CI/badge.svg)](https://github.com/iseki0/goda/actions?query=workflow%3ACI)
[![Go Reference](https://pkg.go.dev/badge/github.com/iseki0/goda.svg)](https://pkg.go.dev/github.com/iseki0/goda)
[![Go Report Card](https://goreportcard.com/badge/github.com/iseki0/goda)](https://goreportcard.com/report/github.com/iseki0/goda)
[![codecov](https://codecov.io/gh/iseki0/goda/graph/badge.svg?token=TBHUZUY561)](https://codecov.io/gh/iseki0/goda)

> **ThreeTen/JSR-310** model in Go

> [English Version](README.md)

受 Java 的 `java.time` 包（JSR-310）启发的 Go 实现，提供**类型安全**且**易用**的不可变日期时间类型。

## 特性

### 核心类型

- 📅 **LocalDate**：无时间的日期（例如：`2024-03-15`）
- ⏰ **LocalTime**：无日期的时间（例如：`14:30:45.123456789`）
- 📆 **LocalDateTime**：日期时间（例如：`2024-03-15T14:30:45.123456789`）
- 🌍 **OffsetDateTime**：带有 UTC 偏移的日期时间（例如：`2024-03-15T14:30:45.123456789+09:00`）
- 🔢 **Field**：日期时间字段枚举（类似 Java 的 `ChronoField`）

### 主要特性

- ✅ **ISO 8601 基本格式**支持（yyyy-MM-dd、HH:mm:ss[.nnnnnnnnn]，使用 'T' 组合）
- ✅ **Java.time 兼容格式**：小数秒对齐到 3 位边界（毫秒、微秒、纳秒）
- ✅ **完整的 JSON 和 SQL** 数据库集成
- ✅ **日期算术**：添加/减去天、月、年，支持溢出处理
- ✅ **字段访问**：获取任何字段值（年、月、小时、一天中的纳秒等）
- ✅ **零拷贝文本序列化**，使用 `encoding.TextAppender`
- ✅ **不可变**：所有操作返回新值
- ✅ **类型安全**：使用不同类型提供编译时安全
- ✅ **零值友好**：正确处理零值

## 安装

```bash
go get github.com/iseki0/goda
```

## 快速开始

### 基本用法

```go
package main

import (
    "fmt"
    "github.com/iseki0/goda"
)

func main() {
    // 创建日期和时间
    date := goda.MustNewLocalDate(2024, goda.March, 15)
    time := goda.MustNewLocalTime(14, 30, 45, 123456789)
    datetime := goda.NewLocalDateTime(date, time)

    fmt.Println(date)     // 2024-03-15
    fmt.Println(time)     // 14:30:45.123456789
    fmt.Println(datetime) // 2024-03-15T14:30:45.123456789

    // 创建偏移日期时间
    offset := goda.MustNewZoneOffsetHours(9)
    offsetDateTime := goda.NewOffsetDateTime(datetime, offset)
    fmt.Println(offsetDateTime) // 2024-03-15T14:30:45.123456789+09:00

    // 从字符串解析
    date, _ = goda.ParseLocalDate("2024-03-15")
    time = goda.MustParseLocalTime("14:30:45.123456789")
    datetime = goda.MustParseLocalDateTime("2024-03-15T14:30:45")
    offsetDateTime = goda.MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")

    // 获取当前日期/时间
    today := goda.LocalDateNow()
    now := goda.LocalTimeNow()
    currentDateTime := goda.LocalDateTimeNow()
    currentOffsetDateTime := goda.OffsetDateTimeNow()

    // 日期算术
    tomorrow := today.PlusDays(1)
    nextMonth := today.PlusMonths(1)
    nextYear := today.PlusYears(1)

    // 比较
    if tomorrow.IsAfter(today) {
        fmt.Println("明天在今天之后！")
    }
}
```

### 组合日期和时间

您可以将 LocalDate 和 LocalTime 组合创建 LocalDateTime：

```go
date := goda.MustNewLocalDate(2024, goda.March, 15)
time := goda.MustNewLocalTime(14, 30, 45, 123456789)

// 将日期与时间组合
dateTime := date.AtTime(time)
fmt.Println(dateTime) // 2024-03-15T14:30:45.123456789

// 将时间与日期组合
dateTime2 := time.AtDate(date)
fmt.Println(dateTime2) // 2024-03-15T14:30:45.123456789
```

### 字段访问

使用 `Field` 枚举访问单个日期时间字段：

```go
date := goda.MustNewLocalDate(2024, goda.March, 15)

// 检查字段支持
fmt.Println(date.IsSupportedField(goda.DayOfMonth))  // true
fmt.Println(date.IsSupportedField(goda.HourOfDay))   // false

// 获取字段值
year := date.GetFieldInt64(goda.YearField)           // 2024
dayOfWeek := date.GetFieldInt64(goda.DayOfWeekField) // 5 (Friday)
dayOfYear := date.GetFieldInt64(goda.DayOfYear)      // 75
epochDays := date.GetFieldInt64(goda.EpochDay)       // 自 Unix 纪元以来的天数

time := goda.MustNewLocalTime(14, 30, 45, 123456789)
hour := time.GetFieldInt64(goda.HourOfDay)           // 14
nanoOfDay := time.GetFieldInt64(goda.NanoOfDay)      // 自午夜以来的总纳秒数
ampm := time.GetFieldInt64(goda.AmPmOfDay)           // 1 (PM)
```

### 使用 OffsetDateTime

```go
// 使用偏移创建
offset := goda.MustNewZoneOffsetHours(9) // 东京：UTC+9
odt := goda.MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")

// 在偏移之间转换
utc := odt.ToUTC()
fmt.Println(utc) // 2024-03-15T05:30:45Z

// 更改偏移，保持相同瞬间
newYork := goda.MustNewZoneOffsetHours(-5)
odtNY := odt.WithOffsetSameInstant(newYork)
fmt.Println(odtNY) // 2024-03-15T00:30:45-05:00

// 比较瞬间（忽略偏移差异）
odt1 := goda.MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")
odt2 := goda.MustParseOffsetDateTime("2024-03-15T05:30:45Z")
fmt.Println(odt1.IsEqual(odt2)) // true (相同瞬间)

// 时间算术（跨越日期边界调整）
later := odt.PlusHours(10) // 为瞬间添加 10 小时
```

### JSON 序列化

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

### 数据库集成

```go
type Record struct {
    ID        int64
    CreatedAt goda.LocalDateTime
    Date      goda.LocalDate
    Timestamp goda.OffsetDateTime
}

// 与 database/sql 配合使用 - 实现了 sql.Scanner 和 driver.Valuer
db.QueryRow("SELECT id, created_at, date, timestamp FROM records WHERE id = ?", 1).Scan(
    &record.ID, &record.CreatedAt, &record.Date, &record.Timestamp,
)
```

## API 概览

### 核心类型

| 类型 | 描述 | 示例 |
|------|------|------|
| `LocalDate` | 无时间的日期 | `2024-03-15` |
| `LocalTime` | 无日期的时间 | `14:30:45.123456789` |
| `LocalDateTime` | 日期时间 | `2024-03-15T14:30:45` |
| `OffsetDateTime` | 带有 UTC 偏移的日期时间 | `2024-03-15T14:30:45+09:00` |
| `ZoneOffset` | UTC 偏移 | `+09:00`, `Z` |
| `Month` | 年中的月份（1-12） | `March` |
| `Year` | 年 | `2024` |
| `DayOfWeek` | 星期中的日期（1=星期一，7=星期日） | `Friday` |
| `Field` | 日期时间字段枚举 | `HourOfDay`, `DayOfMonth` |

### 时间格式

时间值使用带有 **Java.time 兼容** 小数秒对齐的 ISO 8601 格式：

| 精度 | 位数 | 示例 |
|------|------|------|
| 整秒 | 0 | `14:30:45` |
| 毫秒 | 3 | `14:30:45.100`, `14:30:45.123` |
| 微秒 | 6 | `14:30:45.123400`, `14:30:45.123456` |
| 纳秒 | 9 | `14:30:45.000000001`, `14:30:45.123456789` |

小数秒自动对齐到 3 位边界（毫秒、微秒、纳秒），匹配 Java 的 `LocalTime` 行为。解析接受任何长度的小数秒。

### 字段常量（30 个字段）

**时间字段**：`NanoOfSecond`, `NanoOfDay`, `MicroOfSecond`, `MicroOfDay`, `MilliOfSecond`, `MilliOfDay`, `SecondOfMinute`, `SecondOfDay`, `MinuteOfHour`, `MinuteOfDay`, `HourOfAmPm`, `ClockHourOfAmPm`, `HourOfDay`, `ClockHourOfDay`, `AmPmOfDay`

**日期字段**：`DayOfWeekField`, `DayOfMonth`, `DayOfYear`, `EpochDay`, `AlignedDayOfWeekInMonth`, `AlignedDayOfWeekInYear`, `AlignedWeekOfMonth`, `AlignedWeekOfYear`, `MonthOfYear`, `ProlepticMonth`, `YearOfEra`, `YearField`, `Era`

**其他字段**：`InstantSeconds`, `OffsetSeconds`

### 实现的接口

所有类型都实现了：
- `fmt.Stringer`
- `encoding.TextMarshaler` / `encoding.TextUnmarshaler`
- `encoding.TextAppender`（零拷贝文本序列化）
- `json.Marshaler` / `json.Unmarshaler`
- `sql.Scanner` / `driver.Valuer`

## 设计理念

此包遵循 **ThreeTen/JSR-310** 模型（Java 的 `java.time` 包），提供的日期和时间类型：

- **不可变**：所有操作返回新值
- **类型安全**：为日期、时间和日期时间使用不同类型
- **简单格式**：使用 ISO 8601 基本格式（而不是完整的复杂规范）
- **数据库友好**：直接 SQL 集成
- **基于字段的访问**：通过 `GetFieldInt64` 的通用字段访问模式
- **零值安全**：在所有地方正确处理零值

### 何时使用每种类型

**LocalDate、LocalTime、LocalDateTime**

当您只需要没有时区信息的日期/时间时使用本地类型：
- **生日**："3月15日"在任何地方都是3月15日
- **营业时间**：当地语境中的"上午9:00 - 下午5:00"
- **日程安排**：没有时区顾虑的"下午2:30开会"
- **日历日期**：历史日期、重复事件

**OffsetDateTime**

当您需要使用特定 UTC 偏移表示瞬间时使用 OffsetDateTime：
- **API 响应**：带有时区信息的时戳
- **预定事件**：发生在特定瞬间的事件（例如："2024-03-15T14:30:00+09:00"）
- **数据库时戳**：存储带有偏移信息的时戳时
- **国际协调**：当您需要同时知道本地时间和 UTC 偏移时

对于具有 DST 处理的完整时区感知操作，请使用 `ZonedDateTime`（即将推出）。

## 文档

完整的 API 文档可在 [pkg.go.dev](https://pkg.go.dev/github.com/iseki0/goda) 获取。

## 贡献

欢迎贡献！请随时提交 Pull Request。

## 许可证

此项目根据 MIT 许可证授权 - 请查看 LICENSE 文件了解详情。
