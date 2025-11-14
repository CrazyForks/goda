package goda

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLocalDateTime(t *testing.T) {
	date := MustNewLocalDate(2024, March, 15)
	time := MustNewLocalTime(14, 30, 45, 123456789)
	dt := NewLocalDateTime(date, time)

	assert.Equal(t, date, dt.LocalDate())
	assert.Equal(t, time, dt.LocalTime())
	assert.False(t, dt.IsZero())
}

func TestNewLocalDateTimeFromComponents(t *testing.T) {
	t.Run("valid components", func(t *testing.T) {
		dt, err := NewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 123456789)
		require.NoError(t, err)
		assert.Equal(t, Year(2024), dt.Year())
		assert.Equal(t, March, dt.Month())
		assert.Equal(t, 15, dt.DayOfMonth())
		assert.Equal(t, 14, dt.Hour())
		assert.Equal(t, 30, dt.Minute())
		assert.Equal(t, 45, dt.Second())
		assert.Equal(t, 123456789, dt.Nanosecond())
	})

	t.Run("invalid date", func(t *testing.T) {
		_, err := NewLocalDateTimeFromComponents(2024, February, 30, 14, 30, 45, 0)
		assert.Error(t, err)
	})

	t.Run("invalid time", func(t *testing.T) {
		_, err := NewLocalDateTimeFromComponents(2024, March, 15, 25, 30, 45, 0)
		assert.Error(t, err)
	})
}

func TestMustNewLocalDateTimeFromComponents(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		assert.NotPanics(t, func() {
			dt := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 123456789)
			assert.Equal(t, Year(2024), dt.Year())
		})
	})

	t.Run("invalid panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustNewLocalDateTimeFromComponents(2024, February, 30, 14, 30, 45, 0)
		})
	})
}

func TestNewLocalDateTimeByGoTime(t *testing.T) {
	t.Run("valid time", func(t *testing.T) {
		goTime := time.Date(2024, 3, 15, 14, 30, 45, 123456789, time.UTC)
		dt := NewLocalDateTimeByGoTime(goTime)

		assert.Equal(t, Year(2024), dt.Year())
		assert.Equal(t, March, dt.Month())
		assert.Equal(t, 15, dt.DayOfMonth())
		assert.Equal(t, 14, dt.Hour())
		assert.Equal(t, 30, dt.Minute())
		assert.Equal(t, 45, dt.Second())
		assert.Equal(t, 123456789, dt.Nanosecond())
	})

	t.Run("zero time", func(t *testing.T) {
		dt := NewLocalDateTimeByGoTime(time.Time{})
		assert.True(t, dt.IsZero())
	})
}

func TestLocalDateTime_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		var dt LocalDateTime
		assert.True(t, dt.IsZero())
	})

	t.Run("non-zero value", func(t *testing.T) {
		dt := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 0)
		assert.False(t, dt.IsZero())
	})

	t.Run("midnight is not zero", func(t *testing.T) {
		dt := MustNewLocalDateTimeFromComponents(2024, March, 15, 0, 0, 0, 0)
		assert.False(t, dt.IsZero())
	})
}

func TestLocalDateTime_ComponentAccessors(t *testing.T) {
	dt := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 123456789)

	// LocalDate components
	assert.Equal(t, Year(2024), dt.Year())
	assert.Equal(t, March, dt.Month())
	assert.Equal(t, 15, dt.DayOfMonth())
	assert.Equal(t, Friday, dt.DayOfWeek())
	assert.Equal(t, 75, dt.DayOfYear())

	// LocalTime components
	assert.Equal(t, 14, dt.Hour())
	assert.Equal(t, 30, dt.Minute())
	assert.Equal(t, 45, dt.Second())
	assert.Equal(t, 123, dt.Millisecond())
	assert.Equal(t, 123456789, dt.Nanosecond())
}

func TestLocalDateTime_GoTime(t *testing.T) {
	t.Run("non-zero", func(t *testing.T) {
		dt := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 123456789)
		goTime := dt.GoTime()

		assert.Equal(t, 2024, goTime.Year())
		assert.Equal(t, time.March, goTime.Month())
		assert.Equal(t, 15, goTime.Day())
		assert.Equal(t, 14, goTime.Hour())
		assert.Equal(t, 30, goTime.Minute())
		assert.Equal(t, 45, goTime.Second())
		assert.Equal(t, 123456789, goTime.Nanosecond())
		assert.Equal(t, time.UTC, goTime.Location())
	})

	t.Run("zero", func(t *testing.T) {
		var dt LocalDateTime
		goTime := dt.GoTime()
		assert.True(t, goTime.IsZero())
	})
}

func TestLocalDateTime_Compare(t *testing.T) {
	dt1 := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 0)
	dt2 := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 0)
	dt3 := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 46, 0)
	dt4 := MustNewLocalDateTimeFromComponents(2024, March, 16, 14, 30, 45, 0)
	dt5 := MustNewLocalDateTimeFromComponents(2024, March, 15, 15, 30, 45, 0)

	assert.Equal(t, 0, dt1.Compare(dt2))
	assert.Equal(t, -1, dt1.Compare(dt3))
	assert.Equal(t, 1, dt3.Compare(dt1))
	assert.Equal(t, -1, dt1.Compare(dt4))
	assert.Equal(t, 1, dt4.Compare(dt1))
	assert.Equal(t, -1, dt1.Compare(dt5))
	assert.Equal(t, 1, dt5.Compare(dt1))
}

func TestLocalDateTime_IsBefore_IsAfter(t *testing.T) {
	dt1 := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 0)
	dt2 := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 46, 0)
	dt3 := MustNewLocalDateTimeFromComponents(2024, March, 16, 14, 30, 45, 0)

	assert.True(t, dt1.IsBefore(dt2))
	assert.False(t, dt2.IsBefore(dt1))
	assert.False(t, dt1.IsBefore(dt1))

	assert.True(t, dt2.IsAfter(dt1))
	assert.False(t, dt1.IsAfter(dt2))
	assert.False(t, dt1.IsAfter(dt1))

	assert.True(t, dt1.IsBefore(dt3))
	assert.True(t, dt3.IsAfter(dt1))
}

func TestLocalDateTime_PlusDays(t *testing.T) {
	dt := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 0)

	dt2 := dt.PlusDays(10)
	assert.Equal(t, Year(2024), dt2.Year())
	assert.Equal(t, March, dt2.Month())
	assert.Equal(t, 25, dt2.DayOfMonth())
	assert.Equal(t, 14, dt2.Hour())
	assert.Equal(t, 30, dt2.Minute())

	dt3 := dt.PlusDays(-10)
	assert.Equal(t, 5, dt3.DayOfMonth())
}

func TestLocalDateTime_MinusDays(t *testing.T) {
	dt := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 0)

	dt2 := dt.MinusDays(10)
	assert.Equal(t, 5, dt2.DayOfMonth())
	assert.Equal(t, 14, dt2.Hour())
}

func TestLocalDateTime_PlusMonths(t *testing.T) {
	dt := MustNewLocalDateTimeFromComponents(2024, January, 31, 14, 30, 45, 0)

	dt2 := dt.PlusMonths(1)
	assert.Equal(t, February, dt2.Month())
	assert.Equal(t, 29, dt2.DayOfMonth()) // 2024 is leap year
	assert.Equal(t, 14, dt2.Hour())
}

func TestLocalDateTime_MinusMonths(t *testing.T) {
	dt := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 0)

	dt2 := dt.MinusMonths(1)
	assert.Equal(t, February, dt2.Month())
	assert.Equal(t, 15, dt2.DayOfMonth())
}

func TestLocalDateTime_PlusYears(t *testing.T) {
	dt := MustNewLocalDateTimeFromComponents(2024, February, 29, 14, 30, 45, 0)

	dt2 := dt.PlusYears(1)
	assert.Equal(t, Year(2025), dt2.Year())
	assert.Equal(t, February, dt2.Month())
	assert.Equal(t, 28, dt2.DayOfMonth()) // 2025 is not leap year
}

func TestLocalDateTime_MinusYears(t *testing.T) {
	dt := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 0)

	dt2 := dt.MinusYears(1)
	assert.Equal(t, Year(2023), dt2.Year())
}

func TestLocalDateTime_String(t *testing.T) {
	tests := []struct {
		name     string
		dt       LocalDateTime
		expected string
	}{
		{
			name:     "full nanoseconds",
			dt:       MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 123456789),
			expected: "2024-03-15T14:30:45.123456789",
		},
		{
			name:     "milliseconds",
			dt:       MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 123000000),
			expected: "2024-03-15T14:30:45.123",
		},
		{
			name:     "no fractional seconds",
			dt:       MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 0),
			expected: "2024-03-15T14:30:45",
		},
		{
			name:     "midnight",
			dt:       MustNewLocalDateTimeFromComponents(2024, March, 15, 0, 0, 0, 0),
			expected: "2024-03-15T00:00:00",
		},
		{
			name:     "zero value",
			dt:       LocalDateTime{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.dt.String())
		})
	}
}

func TestLocalDateTime_MarshalText(t *testing.T) {
	dt := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 123456789)
	text, err := dt.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "2024-03-15T14:30:45.123456789", string(text))
}

func TestLocalDateTime_UnmarshalText(t *testing.T) {
	t.Run("valid datetime", func(t *testing.T) {
		var dt LocalDateTime
		err := dt.UnmarshalText([]byte("2024-03-15T14:30:45.123456789"))
		require.NoError(t, err)
		assert.Equal(t, Year(2024), dt.Year())
		assert.Equal(t, March, dt.Month())
		assert.Equal(t, 15, dt.DayOfMonth())
		assert.Equal(t, 14, dt.Hour())
		assert.Equal(t, 30, dt.Minute())
		assert.Equal(t, 45, dt.Second())
		assert.Equal(t, 123456789, dt.Nanosecond())
	})

	t.Run("lowercase t separator", func(t *testing.T) {
		var dt LocalDateTime
		err := dt.UnmarshalText([]byte("2024-03-15t14:30:45"))
		require.NoError(t, err)
		assert.Equal(t, Year(2024), dt.Year())
		assert.Equal(t, 14, dt.Hour())
	})

	t.Run("empty string", func(t *testing.T) {
		var dt LocalDateTime
		err := dt.UnmarshalText([]byte(""))
		require.NoError(t, err)
		assert.True(t, dt.IsZero())
	})

	t.Run("invalid date", func(t *testing.T) {
		var dt LocalDateTime
		err := dt.UnmarshalText([]byte("2024-02-30T14:30:45"))
		assert.Error(t, err)
	})

	t.Run("invalid time", func(t *testing.T) {
		var dt LocalDateTime
		err := dt.UnmarshalText([]byte("2024-03-15T25:30:45"))
		assert.Error(t, err)
	})
}

func TestLocalDateTime_MarshalJSON(t *testing.T) {
	dt := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 123456789)
	jsonBytes, err := json.Marshal(dt)
	require.NoError(t, err)
	assert.Equal(t, `"2024-03-15T14:30:45.123456789"`, string(jsonBytes))
}

func TestLocalDateTime_UnmarshalJSON(t *testing.T) {
	t.Run("valid json", func(t *testing.T) {
		var dt LocalDateTime
		err := json.Unmarshal([]byte(`"2024-03-15T14:30:45.123456789"`), &dt)
		require.NoError(t, err)
		assert.Equal(t, Year(2024), dt.Year())
		assert.Equal(t, 14, dt.Hour())
	})

	t.Run("null", func(t *testing.T) {
		var dt LocalDateTime
		err := json.Unmarshal([]byte(`null`), &dt)
		require.NoError(t, err)
		assert.True(t, dt.IsZero())
	})
}

func TestLocalDateTime_Scan(t *testing.T) {
	t.Run("from string", func(t *testing.T) {
		var dt LocalDateTime
		err := dt.Scan("2024-03-15T14:30:45.123456789")
		require.NoError(t, err)
		assert.Equal(t, Year(2024), dt.Year())
		assert.Equal(t, 14, dt.Hour())
	})

	t.Run("from bytes", func(t *testing.T) {
		var dt LocalDateTime
		err := dt.Scan([]byte("2024-03-15T14:30:45"))
		require.NoError(t, err)
		assert.Equal(t, Year(2024), dt.Year())
	})

	t.Run("from time.LocalTime", func(t *testing.T) {
		var dt LocalDateTime
		goTime := time.Date(2024, 3, 15, 14, 30, 45, 123456789, time.UTC)
		err := dt.Scan(goTime)
		require.NoError(t, err)
		assert.Equal(t, Year(2024), dt.Year())
		assert.Equal(t, 14, dt.Hour())
	})

	t.Run("from nil", func(t *testing.T) {
		var dt LocalDateTime
		err := dt.Scan(nil)
		require.NoError(t, err)
		assert.True(t, dt.IsZero())
	})
}

func TestLocalDateTime_Value(t *testing.T) {
	t.Run("non-zero", func(t *testing.T) {
		dt := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 123456789)
		val, err := dt.Value()
		require.NoError(t, err)
		assert.Equal(t, "2024-03-15T14:30:45.123456789", val)
	})

	t.Run("zero", func(t *testing.T) {
		var dt LocalDateTime
		val, err := dt.Value()
		require.NoError(t, err)
		assert.Nil(t, val)
	})
}

func TestParseLocalDateTime(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		dt, err := ParseLocalDateTime("2024-03-15T14:30:45.123456789")
		require.NoError(t, err)
		assert.Equal(t, Year(2024), dt.Year())
		assert.Equal(t, 14, dt.Hour())

		dt, err = ParseLocalDateTime("2024-03-15T14:30:45")
		require.NoError(t, err)
		assert.Equal(t, Year(2024), dt.Year())
	})

	t.Run("empty", func(t *testing.T) {
		dt, err := ParseLocalDateTime("")
		require.NoError(t, err)
		assert.True(t, dt.IsZero())
	})
}

func TestMustParseLocalDateTime(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		assert.NotPanics(t, func() {
			dt := MustParseLocalDateTime("2024-03-15T14:30:45.123456789")
			assert.Equal(t, Year(2024), dt.Year())
		})
	})

	t.Run("invalid panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustParseLocalDateTime("invalid")
		})
	})
}

func TestLocalDateTimeNow(t *testing.T) {
	now := LocalDateTimeNow()
	assert.False(t, now.IsZero())

	// Verify components are in valid ranges
	assert.True(t, now.Year() >= 1970 && now.Year() <= 9999)
	assert.True(t, now.Hour() >= 0 && now.Hour() < 24)
}

func TestLocalDateTimeNowUTC(t *testing.T) {
	nowUTC := LocalDateTimeNowUTC()
	assert.False(t, nowUTC.IsZero())
}

func TestLocalDateTimeNowIn(t *testing.T) {
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	nowTokyo := LocalDateTimeNowIn(tokyo)
	assert.False(t, nowTokyo.IsZero())
}

func TestLocalDateTime_IsLeapYear(t *testing.T) {
	dt2024 := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 0)
	dt2023 := MustNewLocalDateTimeFromComponents(2023, March, 15, 14, 30, 45, 0)

	assert.True(t, dt2024.IsLeapYear())
	assert.False(t, dt2023.IsLeapYear())
}

//go:embed TestFieldGetter.txt
var TestFieldGetterData string

func TestFieldGetter(t *testing.T) {
	var fields = []Field{FieldNanoOfSecond, FieldNanoOfDay, FieldMicroOfSecond, FieldMicroOfDay, FieldMilliOfSecond, FieldMilliOfDay, FieldSecondOfMinute, FieldSecondOfDay, FieldMinuteOfHour, FieldMinuteOfDay, FieldHourOfAmPm, FieldClockHourOfAmPm, FieldHourOfDay, FieldClockHourOfDay, FieldAmPmOfDay, FieldDayOfWeek, FieldAlignedDayOfWeekInMonth, FieldAlignedDayOfWeekInYear, FieldDayOfMonth, FieldDayOfYear, FieldEpochDay, FieldAlignedWeekOfMonth, FieldAlignedWeekOfYear, FieldMonthOfYear, FieldProlepticMonth, FieldYearOfEra, FieldYear, FieldEra}
	for _, line := range strings.Split(TestFieldGetterData, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		cols := strings.Split(line, ",")
		dt, e := ParseLocalDateTime(cols[0])
		if e != nil {
			t.Fatal(e)
		}
		for i, v := range cols[1:] {
			if dt.GetField(fields[i]).Unsupported() {
				continue
			}
			if !assert.Equal(t, v, fmt.Sprint(dt.GetField(fields[i]).Int64())) {
				t.Log(fields[i].String(), line)
				t.Failed()
				return
			}
		}
	}
}
