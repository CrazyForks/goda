package goda

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalDate_Epoch(t *testing.T) {
	var check = func(i int64, tt time.Time) bool {
		var st = MustNewLocalDate(Year(tt.Year()), Month(tt.Month()), tt.Day())
		if !assert.Equal(t, i, st.UnixEpochDays(), tt) {
			return false
		}
		if !assert.Equal(t, st, NewLocalDateByUnixEpochDays(i), tt) {
			return false
		}
		if !assert.Equal(t, st.DayOfWeek().GoWeekday(), tt.Weekday(), tt) {
			return false
		}
		return assert.Equal(t, tt.Unix()/(24*60*60), st.UnixEpochDays(), tt)
	}
	var begin = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := int64(0); i < 100_0000; i++ {
		if !check(i, begin) {
			break
		}
		begin = begin.AddDate(0, 0, 1)
	}
	begin = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	// negative
	for i := int64(0); i > -100_0000; i-- {
		if !check(i, begin) {
			break
		}
		begin = begin.AddDate(0, 0, -1)
	}
}

func TestNewLocalDate(t *testing.T) {
	t.Run("valid dates", func(t *testing.T) {
		d, err := NewLocalDate(2024, January, 1)
		require.NoError(t, err)
		assert.Equal(t, Year(2024), d.Year())
		assert.Equal(t, January, d.Month())
		assert.Equal(t, 1, d.DayOfMonth())

		d, err = NewLocalDate(2024, February, 29) // leap year
		require.NoError(t, err)
		assert.Equal(t, 29, d.DayOfMonth())
	})

	t.Run("invalid day of month", func(t *testing.T) {
		_, err := NewLocalDate(2024, January, 32)
		assert.Error(t, err)

		_, err = NewLocalDate(2023, February, 29) // not a leap year
		assert.Error(t, err)

		_, err = NewLocalDate(2024, February, 30)
		assert.Error(t, err)

		_, err = NewLocalDate(2024, April, 31)
		assert.Error(t, err)
	})

	t.Run("invalid month", func(t *testing.T) {
		_, err := NewLocalDate(2024, Month(0), 1)
		assert.Error(t, err)

		_, err = NewLocalDate(2024, Month(13), 1)
		assert.Error(t, err)
	})
}

func TestMustNewLocalDate(t *testing.T) {
	t.Run("valid date", func(t *testing.T) {
		assert.NotPanics(t, func() {
			d := MustNewLocalDate(2024, March, 15)
			assert.Equal(t, Year(2024), d.Year())
		})
	})

	t.Run("invalid date panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustNewLocalDate(2024, January, 32)
		})
	})
}

func TestLocalDate_IsZero(t *testing.T) {
	var zero LocalDate
	assert.True(t, zero.IsZero())

	d := MustNewLocalDate(2024, January, 1)
	assert.False(t, d.IsZero())
}

func TestLocalDate_IsLeapYear(t *testing.T) {
	tests := []struct {
		year   Year
		isLeap bool
	}{
		{2000, true},  // divisible by 400
		{2004, true},  // divisible by 4
		{2100, false}, // divisible by 100 but not 400
		{2023, false}, // not divisible by 4
		{2024, true},  // divisible by 4
		{1900, false}, // divisible by 100 but not 400
	}

	for _, tt := range tests {
		d := MustNewLocalDate(tt.year, January, 1)
		assert.Equal(t, tt.isLeap, d.IsLeapYear(), "year %d", tt.year)
	}
}

func TestLocalDate_DayOfYear(t *testing.T) {
	tests := []struct {
		date      LocalDate
		dayOfYear int
	}{
		{MustNewLocalDate(2024, January, 1), 1},
		{MustNewLocalDate(2024, January, 31), 31},
		{MustNewLocalDate(2024, February, 1), 32},
		{MustNewLocalDate(2024, February, 29), 60}, // leap year
		{MustNewLocalDate(2023, February, 28), 59}, // non-leap year
		{MustNewLocalDate(2024, March, 1), 61},     // leap year
		{MustNewLocalDate(2023, March, 1), 60},     // non-leap year
		{MustNewLocalDate(2024, December, 31), 366},
		{MustNewLocalDate(2023, December, 31), 365},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.dayOfYear, tt.date.DayOfYear(), "date: %s", tt.date.String())
	}

	var zero LocalDate
	assert.Equal(t, 0, zero.DayOfYear())
}

func TestLocalDate_DayOfWeek(t *testing.T) {
	tests := []struct {
		date      LocalDate
		dayOfWeek DayOfWeek
	}{
		{MustNewLocalDate(2024, 11, 5), Tuesday},  // Known date
		{MustNewLocalDate(2024, 1, 1), Monday},    // New Year 2024
		{MustNewLocalDate(2023, 1, 1), Sunday},    // New Year 2023
		{MustNewLocalDate(2000, 1, 1), Saturday},  // Y2K
		{MustNewLocalDate(1970, 1, 1), Thursday},  // Unix epoch
		{MustNewLocalDate(2024, 2, 29), Thursday}, // Leap day 2024
	}

	for _, tt := range tests {
		assert.Equal(t, tt.dayOfWeek, tt.date.DayOfWeek(), "date: %s", tt.date.String())
	}

	var zero LocalDate
	assert.Equal(t, DayOfWeek(0), zero.DayOfWeek())
}

func TestLocalDate_PlusDays(t *testing.T) {
	d := MustNewLocalDate(2024, January, 15)

	tests := []struct {
		days     int
		expected LocalDate
	}{
		{0, MustNewLocalDate(2024, January, 15)},
		{1, MustNewLocalDate(2024, January, 16)},
		{16, MustNewLocalDate(2024, January, 31)},
		{17, MustNewLocalDate(2024, February, 1)},
		{365, MustNewLocalDate(2025, January, 14)}, // 2024 is leap year
		{-1, MustNewLocalDate(2024, January, 14)},
		{-15, MustNewLocalDate(2023, December, 31)},
	}

	for _, tt := range tests {
		result := d.PlusDays(tt.days)
		assert.Equal(t, tt.expected, result, "days: %d", tt.days)
	}

	var zero LocalDate
	assert.Equal(t, zero, zero.PlusDays(10))
}

func TestLocalDate_MinusDays(t *testing.T) {
	d := MustNewLocalDate(2024, February, 15)

	result := d.MinusDays(10)
	assert.Equal(t, MustNewLocalDate(2024, February, 5), result)

	result = d.MinusDays(15)
	assert.Equal(t, MustNewLocalDate(2024, January, 31), result)
}

func TestLocalDate_PlusMonths(t *testing.T) {
	tests := []struct {
		start    LocalDate
		months   int
		expected LocalDate
	}{
		{MustNewLocalDate(2024, January, 15), 1, MustNewLocalDate(2024, February, 15)},
		{MustNewLocalDate(2024, January, 15), 12, MustNewLocalDate(2025, January, 15)},
		{MustNewLocalDate(2024, January, 15), 13, MustNewLocalDate(2025, February, 15)},
		{MustNewLocalDate(2024, January, 15), -1, MustNewLocalDate(2023, December, 15)},
		{MustNewLocalDate(2024, January, 15), -13, MustNewLocalDate(2022, December, 15)},
		// Test day clamping
		{MustNewLocalDate(2024, January, 31), 1, MustNewLocalDate(2024, February, 29)}, // leap year
		{MustNewLocalDate(2023, January, 31), 1, MustNewLocalDate(2023, February, 28)}, // non-leap year
		{MustNewLocalDate(2024, March, 31), 1, MustNewLocalDate(2024, April, 30)},
		{MustNewLocalDate(2024, May, 31), 1, MustNewLocalDate(2024, June, 30)},
		// Large month additions
		{MustNewLocalDate(2024, January, 1), 24, MustNewLocalDate(2026, January, 1)},
		{MustNewLocalDate(2024, January, 1), -24, MustNewLocalDate(2022, January, 1)},
	}

	for _, tt := range tests {
		result := tt.start.PlusMonths(tt.months)
		assert.Equal(t, tt.expected, result, "start: %s, months: %d", tt.start.String(), tt.months)
	}

	var zero LocalDate
	assert.Equal(t, zero, zero.PlusMonths(10))
}

func TestLocalDate_MinusMonths(t *testing.T) {
	d := MustNewLocalDate(2024, March, 15)

	result := d.MinusMonths(2)
	assert.Equal(t, MustNewLocalDate(2024, January, 15), result)

	result = d.MinusMonths(14)
	assert.Equal(t, MustNewLocalDate(2023, January, 15), result)
}

func TestLocalDate_PlusYears(t *testing.T) {
	tests := []struct {
		start    LocalDate
		years    int
		expected LocalDate
	}{
		{MustNewLocalDate(2024, March, 15), 1, MustNewLocalDate(2025, March, 15)},
		{MustNewLocalDate(2024, March, 15), -1, MustNewLocalDate(2023, March, 15)},
		{MustNewLocalDate(2024, March, 15), 10, MustNewLocalDate(2034, March, 15)},
		// Leap year edge case
		{MustNewLocalDate(2024, February, 29), 1, MustNewLocalDate(2025, February, 28)},
		{MustNewLocalDate(2024, February, 29), 4, MustNewLocalDate(2028, February, 29)},
		{MustNewLocalDate(2024, February, 29), -4, MustNewLocalDate(2020, February, 29)},
	}

	for _, tt := range tests {
		result := tt.start.PlusYears(tt.years)
		assert.Equal(t, tt.expected, result, "start: %s, years: %d", tt.start.String(), tt.years)
	}

	var zero LocalDate
	assert.Equal(t, zero, zero.PlusYears(10))
}

func TestLocalDate_MinusYears(t *testing.T) {
	d := MustNewLocalDate(2024, March, 15)

	result := d.MinusYears(1)
	assert.Equal(t, MustNewLocalDate(2023, March, 15), result)

	result = d.MinusYears(10)
	assert.Equal(t, MustNewLocalDate(2014, March, 15), result)
}

func TestLocalDate_Compare(t *testing.T) {
	d1 := MustNewLocalDate(2024, March, 15)
	d2 := MustNewLocalDate(2024, March, 15)
	d3 := MustNewLocalDate(2024, March, 16)
	d4 := MustNewLocalDate(2024, February, 15)
	d5 := MustNewLocalDate(2023, March, 15)

	assert.Equal(t, 0, d1.Compare(d2))
	assert.Equal(t, -1, d1.Compare(d3))
	assert.Equal(t, 1, d3.Compare(d1))
	assert.Equal(t, 1, d1.Compare(d4))
	assert.Equal(t, 1, d1.Compare(d5))
}

func TestLocalDate_IsBefore(t *testing.T) {
	d1 := MustNewLocalDate(2024, March, 15)
	d2 := MustNewLocalDate(2024, March, 16)
	d3 := MustNewLocalDate(2024, March, 15)

	assert.True(t, d1.IsBefore(d2))
	assert.False(t, d2.IsBefore(d1))
	assert.False(t, d1.IsBefore(d3))
}

func TestLocalDate_IsAfter(t *testing.T) {
	d1 := MustNewLocalDate(2024, March, 15)
	d2 := MustNewLocalDate(2024, March, 16)
	d3 := MustNewLocalDate(2024, March, 15)

	assert.False(t, d1.IsAfter(d2))
	assert.True(t, d2.IsAfter(d1))
	assert.False(t, d1.IsAfter(d3))
}

func TestLocalDate_GoTime(t *testing.T) {
	d := MustNewLocalDate(2024, March, 15)
	goTime := d.GoTime()

	assert.Equal(t, 2024, goTime.Year())
	assert.Equal(t, time.March, goTime.Month())
	assert.Equal(t, 15, goTime.Day())
	assert.Equal(t, 0, goTime.Hour())
	assert.Equal(t, 0, goTime.Minute())
	assert.Equal(t, 0, goTime.Second())
	assert.Equal(t, time.UTC, goTime.Location())

	var zero LocalDate
	assert.True(t, zero.GoTime().IsZero())
}

func TestNewLocalDateByGoTime(t *testing.T) {
	goTime := time.Date(2024, time.March, 15, 14, 30, 45, 0, time.UTC)
	d := NewLocalDateByGoTime(goTime)

	assert.Equal(t, Year(2024), d.Year())
	assert.Equal(t, March, d.Month())
	assert.Equal(t, 15, d.DayOfMonth())

	// Test with zero time
	d = NewLocalDateByGoTime(time.Time{})
	assert.True(t, d.IsZero())
}

func TestLocalDate_String(t *testing.T) {
	tests := []struct {
		date     LocalDate
		expected string
	}{
		{MustNewLocalDate(2024, March, 15), "2024-03-15"},
		{MustNewLocalDate(2024, January, 1), "2024-01-01"},
		{MustNewLocalDate(2024, December, 31), "2024-12-31"},
		{MustNewLocalDate(1999, June, 5), "1999-06-05"},
		{LocalDate{}, ""},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.date.String())
	}
}

func TestLocalDate_MarshalText(t *testing.T) {
	d := MustNewLocalDate(2024, March, 15)
	text, err := d.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "2024-03-15", string(text))

	var zero LocalDate
	text, err = zero.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "", string(text))
}

func TestLocalDate_UnmarshalText(t *testing.T) {
	t.Run("valid dates", func(t *testing.T) {
		var d LocalDate
		err := d.UnmarshalText([]byte("2024-03-15"))
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalDate(2024, March, 15), d)

		err = d.UnmarshalText([]byte("1999-12-31"))
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalDate(1999, December, 31), d)
	})

	t.Run("empty string", func(t *testing.T) {
		var d LocalDate
		err := d.UnmarshalText([]byte(""))
		require.NoError(t, err)
		assert.True(t, d.IsZero())
	})

	t.Run("invalid format", func(t *testing.T) {
		var d LocalDate
		err := d.UnmarshalText([]byte("2024/03/15"))
		assert.Error(t, err)

		err = d.UnmarshalText([]byte("2024-3-15"))
		assert.Error(t, err)

		err = d.UnmarshalText([]byte("not-a-date"))
		assert.Error(t, err)
	})

	t.Run("invalid date", func(t *testing.T) {
		var d LocalDate
		err := d.UnmarshalText([]byte("2024-02-30"))
		assert.Error(t, err)

		err = d.UnmarshalText([]byte("2024-13-01"))
		assert.Error(t, err)
	})
}

func TestLocalDate_MarshalJSON(t *testing.T) {
	d := MustNewLocalDate(2024, March, 15)
	data, err := json.Marshal(d)
	require.NoError(t, err)
	assert.Equal(t, `"2024-03-15"`, string(data))

	var zero LocalDate
	data, err = json.Marshal(zero)
	require.NoError(t, err)
	assert.Equal(t, `""`, string(data))
}

func TestLocalDate_UnmarshalJSON(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		var d LocalDate
		err := json.Unmarshal([]byte(`"2024-03-15"`), &d)
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalDate(2024, March, 15), d)
	})

	t.Run("empty string", func(t *testing.T) {
		var d LocalDate
		err := json.Unmarshal([]byte(`""`), &d)
		require.NoError(t, err)
		assert.True(t, d.IsZero())
	})

	t.Run("null", func(t *testing.T) {
		var d LocalDate
		err := json.Unmarshal([]byte(`null`), &d)
		require.NoError(t, err)
		assert.True(t, d.IsZero())
	})

	t.Run("invalid JSON", func(t *testing.T) {
		var d LocalDate
		err := json.Unmarshal([]byte(`"invalid-date"`), &d)
		assert.Error(t, err)
	})
}

func TestLocalDate_Scan(t *testing.T) {
	t.Run("nil value", func(t *testing.T) {
		var d LocalDate
		err := d.Scan(nil)
		require.NoError(t, err)
		assert.True(t, d.IsZero())
	})

	t.Run("string value", func(t *testing.T) {
		var d LocalDate
		err := d.Scan("2024-03-15")
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalDate(2024, March, 15), d)
	})

	t.Run("byte slice value", func(t *testing.T) {
		var d LocalDate
		err := d.Scan([]byte("2024-03-15"))
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalDate(2024, March, 15), d)
	})

	t.Run("time.LocalTime value", func(t *testing.T) {
		var d LocalDate
		goTime := time.Date(2024, time.March, 15, 14, 30, 0, 0, time.UTC)
		err := d.Scan(goTime)
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalDate(2024, March, 15), d)
	})
}

func TestLocalDate_Value(t *testing.T) {
	d := MustNewLocalDate(2024, March, 15)
	val, err := d.Value()
	require.NoError(t, err)
	assert.Equal(t, "2024-03-15", val)

	var zero LocalDate
	val, err = zero.Value()
	require.NoError(t, err)
	assert.Nil(t, val)
}

func TestLocalDate_AppendText(t *testing.T) {
	d := MustNewLocalDate(2024, March, 15)
	buf := []byte("LocalDate: ")
	buf, err := d.AppendText(buf)
	require.NoError(t, err)
	assert.Equal(t, "LocalDate: 2024-03-15", string(buf))

	var zero LocalDate
	buf = []byte("LocalDate: ")
	buf, err = zero.AppendText(buf)
	require.NoError(t, err)
	assert.Equal(t, "LocalDate: ", string(buf))
}

func TestLocalDate_SpecialDates(t *testing.T) {
	t.Run("leap year Feb 29", func(t *testing.T) {
		d := MustNewLocalDate(2024, February, 29)
		assert.Equal(t, 29, d.DayOfMonth())
		assert.Equal(t, 60, d.DayOfYear())

		// Add years to non-leap year
		next := d.PlusYears(1)
		assert.Equal(t, MustNewLocalDate(2025, February, 28), next)
	})

	t.Run("year boundaries", func(t *testing.T) {
		d := MustNewLocalDate(2024, December, 31)
		assert.Equal(t, 366, d.DayOfYear()) // leap year

		next := d.PlusDays(1)
		assert.Equal(t, MustNewLocalDate(2025, January, 1), next)
		assert.Equal(t, 1, next.DayOfYear())
	})

	t.Run("negative years", func(t *testing.T) {
		// Test that the system can handle negative years
		d := MustNewLocalDate(-1, January, 1)
		assert.Equal(t, Year(-1), d.Year())

		d = MustNewLocalDate(-100, December, 31)
		assert.Equal(t, Year(-100), d.Year())
		assert.Equal(t, December, d.Month())
	})
}

func TestLocalDate_MonthBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		date     LocalDate
		addDays  int
		expected LocalDate
	}{
		{"Jan to Feb", MustNewLocalDate(2024, January, 31), 1, MustNewLocalDate(2024, February, 1)},
		{"Feb to Mar (leap)", MustNewLocalDate(2024, February, 29), 1, MustNewLocalDate(2024, March, 1)},
		{"Feb to Mar (non-leap)", MustNewLocalDate(2023, February, 28), 1, MustNewLocalDate(2023, March, 1)},
		{"Apr to May", MustNewLocalDate(2024, April, 30), 1, MustNewLocalDate(2024, May, 1)},
		{"Dec to Jan", MustNewLocalDate(2024, December, 31), 1, MustNewLocalDate(2025, January, 1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.date.PlusDays(tt.addDays)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLocalDateNow(t *testing.T) {
	// Test that LocalDateNow() returns a valid date
	today := LocalDateNow()
	assert.False(t, today.IsZero(), "LocalDateNow should not be zero")

	// Test that it matches time.Now()
	now := time.Now()
	expected := NewLocalDateByGoTime(now)

	// Allow for the possibility that the date changed between calls
	// (very unlikely but possible at midnight)
	diff := today.Compare(expected)
	assert.True(t, diff >= -1 && diff <= 1, "LocalDateNow should be within 1 day of current time")
}

func TestLocalDateNowUTC(t *testing.T) {
	todayUTC := LocalDateNowUTC()
	assert.False(t, todayUTC.IsZero(), "LocalDateNowUTC should not be zero")

	// Test that it matches time.Now().UTC()
	now := time.Now().UTC()
	expected := NewLocalDateByGoTime(now)

	diff := todayUTC.Compare(expected)
	assert.True(t, diff >= -1 && diff <= 1, "LocalDateNowUTC should be within 1 day of current UTC time")
}

func TestParseLocalDate(t *testing.T) {
	t.Run("valid dates", func(t *testing.T) {
		date, err := ParseLocalDate("2024-03-15")
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalDate(2024, March, 15), date)

		date, err = ParseLocalDate("1999-12-31")
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalDate(1999, December, 31), date)

		date, err = ParseLocalDate("2000-02-29") // leap year
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalDate(2000, February, 29), date)
	})

	t.Run("invalid format", func(t *testing.T) {
		_, err := ParseLocalDate("2024/03/15")
		assert.Error(t, err)

		_, err = ParseLocalDate("2024-3-15")
		assert.Error(t, err)

		_, err = ParseLocalDate("not-a-date")
		assert.Error(t, err)
	})

	t.Run("invalid date", func(t *testing.T) {
		_, err := ParseLocalDate("2024-02-30")
		assert.Error(t, err)

		_, err = ParseLocalDate("2024-13-01")
		assert.Error(t, err)
	})

	t.Run("empty string", func(t *testing.T) {
		date, err := ParseLocalDate("")
		require.NoError(t, err)
		assert.True(t, date.IsZero())
	})
}

func TestMustParseLocalDate(t *testing.T) {
	t.Run("valid date", func(t *testing.T) {
		assert.NotPanics(t, func() {
			date := MustParseLocalDate("2024-03-15")
			assert.Equal(t, MustNewLocalDate(2024, March, 15), date)
		})
	})

	t.Run("invalid date panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustParseLocalDate("2024-02-30")
		})

		assert.Panics(t, func() {
			MustParseLocalDate("invalid")
		})
	})
}

func TestLocalDateNowIn(t *testing.T) {
	// Test with different time zones
	locations := []struct {
		name string
		loc  *time.Location
	}{
		{"UTC", time.UTC},
		{"Local", time.Local},
	}

	for _, tt := range locations {
		t.Run(tt.name, func(t *testing.T) {
			todayIn := LocalDateNowIn(tt.loc)
			assert.False(t, todayIn.IsZero(), "LocalDateNowIn should not be zero")

			now := time.Now().In(tt.loc)
			expected := NewLocalDateByGoTime(now)

			diff := todayIn.Compare(expected)
			assert.True(t, diff >= -1 && diff <= 1, "LocalDateNowIn should be within 1 day of current time in specified zone")
		})
	}
}
