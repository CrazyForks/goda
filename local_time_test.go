package goda

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLocalTime(t *testing.T) {
	t.Run("valid times", func(t *testing.T) {
		lt, err := NewLocalTime(0, 0, 0, 0)
		require.NoError(t, err)
		assert.Equal(t, 0, lt.Hour())
		assert.Equal(t, 0, lt.Minute())
		assert.Equal(t, 0, lt.Second())
		assert.Equal(t, 0, lt.Nanosecond())

		lt, err = NewLocalTime(23, 59, 59, 999999999)
		require.NoError(t, err)
		assert.Equal(t, 23, lt.Hour())
		assert.Equal(t, 59, lt.Minute())
		assert.Equal(t, 59, lt.Second())
		assert.Equal(t, 999999999, lt.Nanosecond())

		lt, err = NewLocalTime(12, 30, 45, 123456789)
		require.NoError(t, err)
		assert.Equal(t, 12, lt.Hour())
		assert.Equal(t, 30, lt.Minute())
		assert.Equal(t, 45, lt.Second())
		assert.Equal(t, 123456789, lt.Nanosecond())
	})

	t.Run("invalid hour", func(t *testing.T) {
		_, err := NewLocalTime(24, 0, 0, 0)
		assert.Error(t, err)

		_, err = NewLocalTime(-1, 0, 0, 0)
		assert.Error(t, err)

		_, err = NewLocalTime(25, 0, 0, 0)
		assert.Error(t, err)
	})

	t.Run("invalid minute", func(t *testing.T) {
		_, err := NewLocalTime(0, 60, 0, 0)
		assert.Error(t, err)

		_, err = NewLocalTime(0, -1, 0, 0)
		assert.Error(t, err)

		_, err = NewLocalTime(0, 61, 0, 0)
		assert.Error(t, err)
	})

	t.Run("invalid second", func(t *testing.T) {
		_, err := NewLocalTime(0, 0, 60, 0)
		assert.Error(t, err)

		_, err = NewLocalTime(0, 0, -1, 0)
		assert.Error(t, err)

		_, err = NewLocalTime(0, 0, 61, 0)
		assert.Error(t, err)
	})

	t.Run("invalid nanosecond", func(t *testing.T) {
		_, err := NewLocalTime(0, 0, 0, 1000000000)
		assert.Error(t, err)

		_, err = NewLocalTime(0, 0, 0, -1)
		assert.Error(t, err)
	})
}

func TestMustNewLocalTime(t *testing.T) {
	t.Run("valid time", func(t *testing.T) {
		assert.NotPanics(t, func() {
			lt := MustNewLocalTime(14, 30, 45, 123456789)
			assert.Equal(t, 14, lt.Hour())
			assert.Equal(t, 30, lt.Minute())
			assert.Equal(t, 45, lt.Second())
			assert.Equal(t, 123456789, lt.Nanosecond())
		})
	})

	t.Run("invalid time panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustNewLocalTime(24, 0, 0, 0)
		})

		assert.Panics(t, func() {
			MustNewLocalTime(0, 60, 0, 0)
		})

		assert.Panics(t, func() {
			MustNewLocalTime(0, 0, 60, 0)
		})

		assert.Panics(t, func() {
			MustNewLocalTime(0, 0, 0, 1000000000)
		})
	})
}

func TestLocalTime_IsZero(t *testing.T) {
	var zero LocalTime
	assert.True(t, zero.IsZero())

	lt := MustNewLocalTime(0, 0, 0, 0)
	assert.False(t, lt.IsZero())

	lt = MustNewLocalTime(12, 30, 45, 0)
	assert.False(t, lt.IsZero())
}

func TestLocalTime_Components(t *testing.T) {
	tests := []struct {
		name        string
		hour        int
		minute      int
		second      int
		nanosecond  int
		millisecond int
	}{
		{"midnight", 0, 0, 0, 0, 0},
		{"noon", 12, 0, 0, 0, 0},
		{"end of day", 23, 59, 59, 999999999, 999},
		{"with milliseconds", 14, 30, 45, 123000000, 123},
		{"with nanoseconds", 9, 15, 30, 123456789, 123},
		{"1 second before midnight", 23, 59, 59, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lt := MustNewLocalTime(tt.hour, tt.minute, tt.second, tt.nanosecond)
			assert.Equal(t, tt.hour, lt.Hour(), "Hour")
			assert.Equal(t, tt.minute, lt.Minute(), "Minute")
			assert.Equal(t, tt.second, lt.Second(), "Second")
			assert.Equal(t, tt.nanosecond, lt.Nanosecond(), "Nanosecond")
			assert.Equal(t, tt.millisecond, lt.Millisecond(), "Millisecond")
		})
	}

	var zero LocalTime
	assert.Equal(t, 0, zero.Hour())
	assert.Equal(t, 0, zero.Minute())
	assert.Equal(t, 0, zero.Second())
	assert.Equal(t, 0, zero.Nanosecond())
	assert.Equal(t, 0, zero.Millisecond())
}

func TestLocalTime_Compare(t *testing.T) {
	t1 := MustNewLocalTime(12, 30, 45, 123456789)
	t2 := MustNewLocalTime(12, 30, 45, 123456789)
	t3 := MustNewLocalTime(12, 30, 46, 0)
	t4 := MustNewLocalTime(12, 31, 0, 0)
	t5 := MustNewLocalTime(13, 0, 0, 0)
	t6 := MustNewLocalTime(12, 30, 45, 123456788)

	assert.Equal(t, 0, t1.Compare(t2), "same time")
	assert.Equal(t, -1, t1.Compare(t3), "earlier by second")
	assert.Equal(t, 1, t3.Compare(t1), "later by second")
	assert.Equal(t, -1, t1.Compare(t4), "earlier by minute")
	assert.Equal(t, -1, t1.Compare(t5), "earlier by hour")
	assert.Equal(t, 1, t1.Compare(t6), "later by nanosecond")

	var zero LocalTime
	assert.Equal(t, 0, zero.Compare(LocalTime{}), "zero equals zero")
	assert.Equal(t, -1, zero.Compare(t1), "zero is before non-zero")
	assert.Equal(t, 1, t1.Compare(zero), "non-zero is after zero")
}

func TestLocalTime_IsBefore(t *testing.T) {
	t1 := MustNewLocalTime(10, 0, 0, 0)
	t2 := MustNewLocalTime(11, 0, 0, 0)
	t3 := MustNewLocalTime(10, 0, 0, 0)

	assert.True(t, t1.IsBefore(t2))
	assert.False(t, t2.IsBefore(t1))
	assert.False(t, t1.IsBefore(t3))
}

func TestLocalTime_IsAfter(t *testing.T) {
	t1 := MustNewLocalTime(10, 0, 0, 0)
	t2 := MustNewLocalTime(11, 0, 0, 0)
	t3 := MustNewLocalTime(10, 0, 0, 0)

	assert.False(t, t1.IsAfter(t2))
	assert.True(t, t2.IsAfter(t1))
	assert.False(t, t1.IsAfter(t3))
}

func TestLocalTime_GoTime(t *testing.T) {
	lt := MustNewLocalTime(14, 30, 45, 123456789)
	goTime := lt.GoTime()

	assert.Equal(t, 14, goTime.Hour())
	assert.Equal(t, 30, goTime.Minute())
	assert.Equal(t, 45, goTime.Second())
	assert.Equal(t, 123456789, goTime.Nanosecond())
	assert.Equal(t, time.UTC, goTime.Location())

	// Check that date is set to epoch
	assert.Equal(t, 1970, goTime.Year())
	assert.Equal(t, time.January, goTime.Month())
	assert.Equal(t, 1, goTime.Day())

	var zero LocalTime
	assert.True(t, zero.GoTime().IsZero())
}

func TestNewLocalTimeByGoTime(t *testing.T) {
	goTime := time.Date(2024, time.March, 15, 14, 30, 45, 123456789, time.UTC)
	lt := LocalTimeOfGoTime(goTime)

	assert.Equal(t, 14, lt.Hour())
	assert.Equal(t, 30, lt.Minute())
	assert.Equal(t, 45, lt.Second())
	assert.Equal(t, 123456789, lt.Nanosecond())

	// Test with different time zone (should ignore timezone)
	loc, _ := time.LoadLocation("America/New_York")
	goTime = time.Date(2024, time.March, 15, 14, 30, 45, 123456789, loc)
	lt = LocalTimeOfGoTime(goTime)

	assert.Equal(t, 14, lt.Hour())
	assert.Equal(t, 30, lt.Minute())
	assert.Equal(t, 45, lt.Second())
	assert.Equal(t, 123456789, lt.Nanosecond())

	// Test with zero time
	lt = LocalTimeOfGoTime(time.Time{})
	assert.True(t, lt.IsZero())
}

func TestLocalTime_String(t *testing.T) {
	tests := []struct {
		time     LocalTime
		expected string
	}{
		{MustNewLocalTime(0, 0, 0, 0), "00:00:00"},
		{MustNewLocalTime(12, 30, 45, 0), "12:30:45"},
		{MustNewLocalTime(23, 59, 59, 0), "23:59:59"},
		{MustNewLocalTime(9, 5, 7, 0), "09:05:07"},
		{MustNewLocalTime(14, 30, 45, 123000000), "14:30:45.123"},
		{MustNewLocalTime(14, 30, 45, 123456000), "14:30:45.123456"},
		{MustNewLocalTime(14, 30, 45, 123456789), "14:30:45.123456789"},
		{MustNewLocalTime(14, 30, 45, 100000000), "14:30:45.100"},
		{MustNewLocalTime(14, 30, 45, 120000000), "14:30:45.120"},
		{MustNewLocalTime(14, 30, 45, 1), "14:30:45.000000001"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.time.String())
		})
	}
}

func TestLocalTime_MarshalText(t *testing.T) {
	lt := MustNewLocalTime(14, 30, 45, 123456789)
	text, err := lt.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "14:30:45.123456789", string(text))

	lt = MustNewLocalTime(9, 5, 7, 0)
	text, err = lt.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "09:05:07", string(text))

	var zero LocalTime
	text, err = zero.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "", string(text))
}

func TestLocalTime_UnmarshalText(t *testing.T) {
	t.Run("valid times", func(t *testing.T) {
		var lt LocalTime
		err := lt.UnmarshalText([]byte("14:30:45"))
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalTime(14, 30, 45, 0), lt)

		err = lt.UnmarshalText([]byte("09:05:07"))
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalTime(9, 5, 7, 0), lt)

		err = lt.UnmarshalText([]byte("00:00:00"))
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalTime(0, 0, 0, 0), lt)

		err = lt.UnmarshalText([]byte("23:59:59"))
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalTime(23, 59, 59, 0), lt)
	})

	t.Run("with fractional seconds", func(t *testing.T) {
		var lt LocalTime

		err := lt.UnmarshalText([]byte("14:30:45.123"))
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalTime(14, 30, 45, 123000000), lt)

		err = lt.UnmarshalText([]byte("14:30:45.123456"))
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalTime(14, 30, 45, 123456000), lt)

		err = lt.UnmarshalText([]byte("14:30:45.123456789"))
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalTime(14, 30, 45, 123456789), lt)

		err = lt.UnmarshalText([]byte("14:30:45.1"))
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalTime(14, 30, 45, 100000000), lt)

		err = lt.UnmarshalText([]byte("14:30:45.000000001"))
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalTime(14, 30, 45, 1), lt)
	})

	t.Run("empty string", func(t *testing.T) {
		var lt LocalTime
		err := lt.UnmarshalText([]byte(""))
		require.NoError(t, err)
		assert.True(t, lt.IsZero())
	})

	t.Run("invalid format", func(t *testing.T) {
		var lt LocalTime
		err := lt.UnmarshalText([]byte("14-30-45"))
		assert.Error(t, err)

		err = lt.UnmarshalText([]byte("14:30"))
		assert.Error(t, err)

		err = lt.UnmarshalText([]byte("not-a-time"))
		assert.Error(t, err)

		err = lt.UnmarshalText([]byte("1:2:3"))
		assert.Error(t, err)
	})

	t.Run("invalid values", func(t *testing.T) {
		var lt LocalTime
		err := lt.UnmarshalText([]byte("24:00:00"))
		assert.Error(t, err)

		err = lt.UnmarshalText([]byte("23:60:00"))
		assert.Error(t, err)

		err = lt.UnmarshalText([]byte("23:59:60"))
		assert.Error(t, err)

		err = lt.UnmarshalText([]byte("25:00:00"))
		assert.Error(t, err)

		err = lt.UnmarshalText([]byte("12:61:00"))
		assert.Error(t, err)

		err = lt.UnmarshalText([]byte("12:30:61"))
		assert.Error(t, err)
	})
}

func TestLocalTime_MarshalJSON(t *testing.T) {
	lt := MustNewLocalTime(14, 30, 45, 123456789)
	data, err := json.Marshal(lt)
	require.NoError(t, err)
	assert.Equal(t, `"14:30:45.123456789"`, string(data))

	lt = MustNewLocalTime(9, 5, 7, 0)
	data, err = json.Marshal(lt)
	require.NoError(t, err)
	assert.Equal(t, `"09:05:07"`, string(data))

	var zero LocalTime
	data, err = json.Marshal(zero)
	require.NoError(t, err)
	assert.Equal(t, `""`, string(data))
}

func TestLocalTime_UnmarshalJSON(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		var lt LocalTime
		err := json.Unmarshal([]byte(`"14:30:45"`), &lt)
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalTime(14, 30, 45, 0), lt)

		err = json.Unmarshal([]byte(`"14:30:45.123456789"`), &lt)
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalTime(14, 30, 45, 123456789), lt)
	})

	t.Run("empty string", func(t *testing.T) {
		var lt LocalTime
		err := json.Unmarshal([]byte(`""`), &lt)
		require.NoError(t, err)
		assert.True(t, lt.IsZero())
	})

	t.Run("null", func(t *testing.T) {
		var lt LocalTime
		err := json.Unmarshal([]byte(`null`), &lt)
		require.NoError(t, err)
		assert.True(t, lt.IsZero())
	})

	t.Run("invalid JSON", func(t *testing.T) {
		var lt LocalTime
		err := json.Unmarshal([]byte(`"invalid-time"`), &lt)
		assert.Error(t, err)

		err = json.Unmarshal([]byte(`"24:00:00"`), &lt)
		assert.Error(t, err)
	})
}

func TestLocalTime_Scan(t *testing.T) {
	t.Run("nil value", func(t *testing.T) {
		var lt LocalTime
		err := lt.Scan(nil)
		require.NoError(t, err)
		assert.True(t, lt.IsZero())
	})

	t.Run("string value", func(t *testing.T) {
		var lt LocalTime
		err := lt.Scan("14:30:45")
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalTime(14, 30, 45, 0), lt)

		err = lt.Scan("14:30:45.123456789")
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalTime(14, 30, 45, 123456789), lt)
	})

	t.Run("time.LocalTime value", func(t *testing.T) {
		var lt LocalTime
		goTime := time.Date(2024, time.March, 15, 14, 30, 45, 123456789, time.UTC)
		err := lt.Scan(goTime)
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalTime(14, 30, 45, 123456789), lt)
	})

	t.Run("unsupported type", func(t *testing.T) {
		var lt LocalTime
		err := lt.Scan(12345)
		assert.Error(t, err)
	})
}

func TestLocalTime_Value(t *testing.T) {
	lt := MustNewLocalTime(14, 30, 45, 123456789)
	val, err := lt.Value()
	require.NoError(t, err)
	assert.Equal(t, "14:30:45.123456789", val)

	lt = MustNewLocalTime(9, 5, 7, 0)
	val, err = lt.Value()
	require.NoError(t, err)
	assert.Equal(t, "09:05:07", val)

	var zero LocalTime
	val, err = zero.Value()
	require.NoError(t, err)
	assert.Nil(t, val)
}

func TestLocalTime_AppendText(t *testing.T) {
	lt := MustNewLocalTime(14, 30, 45, 123456789)
	buf := []byte("LocalTime: ")
	buf, err := lt.AppendText(buf)
	require.NoError(t, err)
	assert.Equal(t, "LocalTime: 14:30:45.123456789", string(buf))

	lt = MustNewLocalTime(9, 5, 7, 0)
	buf = []byte("LocalTime: ")
	buf, err = lt.AppendText(buf)
	require.NoError(t, err)
	assert.Equal(t, "LocalTime: 09:05:07", string(buf))

	var zero LocalTime
	buf = []byte("LocalTime: ")
	buf, err = zero.AppendText(buf)
	require.NoError(t, err)
	assert.Equal(t, "LocalTime: ", string(buf))
}

func TestLocalTime_SpecialCases(t *testing.T) {
	t.Run("midnight", func(t *testing.T) {
		lt := MustNewLocalTime(0, 0, 0, 0)
		assert.Equal(t, "00:00:00", lt.String())
		assert.Equal(t, 0, lt.Hour())
		assert.Equal(t, 0, lt.Minute())
		assert.Equal(t, 0, lt.Second())
		assert.Equal(t, 0, lt.Nanosecond())
	})

	t.Run("one nanosecond before midnight", func(t *testing.T) {
		lt := MustNewLocalTime(23, 59, 59, 999999999)
		assert.Equal(t, 23, lt.Hour())
		assert.Equal(t, 59, lt.Minute())
		assert.Equal(t, 59, lt.Second())
		assert.Equal(t, 999999999, lt.Nanosecond())
		assert.Equal(t, 999, lt.Millisecond())
	})

	t.Run("noon", func(t *testing.T) {
		lt := MustNewLocalTime(12, 0, 0, 0)
		assert.Equal(t, "12:00:00", lt.String())
		assert.Equal(t, 12, lt.Hour())
	})

	t.Run("fractional seconds precision", func(t *testing.T) {
		// Millisecond precision
		lt := MustNewLocalTime(12, 0, 0, 123000000)
		assert.Equal(t, "12:00:00.123", lt.String())
		assert.Equal(t, 123, lt.Millisecond())

		// Microsecond precision
		lt = MustNewLocalTime(12, 0, 0, 123456000)
		assert.Equal(t, "12:00:00.123456", lt.String())

		// Nanosecond precision
		lt = MustNewLocalTime(12, 0, 0, 123456789)
		assert.Equal(t, "12:00:00.123456789", lt.String())

		// Single digit fractional second (aligned to milliseconds)
		lt = MustNewLocalTime(12, 0, 0, 100000000)
		assert.Equal(t, "12:00:00.100", lt.String())

		// Trailing zeros aligned to 3-digit boundaries
		lt = MustNewLocalTime(12, 0, 0, 120000000)
		assert.Equal(t, "12:00:00.120", lt.String())

		// More trailing zero examples (aligned to microseconds)
		lt = MustNewLocalTime(12, 0, 0, 123400000)
		assert.Equal(t, "12:00:00.123400", lt.String())
	})
}

func TestLocalTime_BoundaryValues(t *testing.T) {
	tests := []struct {
		name       string
		hour       int
		minute     int
		second     int
		nanosecond int
	}{
		{"min values", 0, 0, 0, 0},
		{"max values", 23, 59, 59, 999999999},
		{"max hour", 23, 0, 0, 0},
		{"max minute", 0, 59, 0, 0},
		{"max second", 0, 0, 59, 0},
		{"max nanosecond", 0, 0, 0, 999999999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lt := MustNewLocalTime(tt.hour, tt.minute, tt.second, tt.nanosecond)
			assert.Equal(t, tt.hour, lt.Hour())
			assert.Equal(t, tt.minute, lt.Minute())
			assert.Equal(t, tt.second, lt.Second())
			assert.Equal(t, tt.nanosecond, lt.Nanosecond())

			// Round-trip through string
			str := lt.String()
			var lt2 LocalTime
			err := lt2.UnmarshalText([]byte(str))
			require.NoError(t, err)
			assert.Equal(t, lt, lt2)
		})
	}
}

func TestLocalTime_Serialization(t *testing.T) {
	tests := []struct {
		name       string
		time       LocalTime
		textFormat string
	}{
		{"midnight", MustNewLocalTime(0, 0, 0, 0), "00:00:00"},
		{"noon", MustNewLocalTime(12, 0, 0, 0), "12:00:00"},
		{"with milliseconds", MustNewLocalTime(14, 30, 45, 123000000), "14:30:45.123"},
		{"with microseconds", MustNewLocalTime(14, 30, 45, 123456000), "14:30:45.123456"},
		{"with nanoseconds", MustNewLocalTime(14, 30, 45, 123456789), "14:30:45.123456789"},
		{"end of day", MustNewLocalTime(23, 59, 59, 999999999), "23:59:59.999999999"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test String
			assert.Equal(t, tt.textFormat, tt.time.String())

			// Test MarshalText
			text, err := tt.time.MarshalText()
			require.NoError(t, err)
			assert.Equal(t, tt.textFormat, string(text))

			// Test UnmarshalText
			var lt LocalTime
			err = lt.UnmarshalText([]byte(tt.textFormat))
			require.NoError(t, err)
			assert.Equal(t, tt.time, lt)

			// Test MarshalJSON
			jsonData, err := json.Marshal(tt.time)
			require.NoError(t, err)
			assert.Equal(t, `"`+tt.textFormat+`"`, string(jsonData))

			// Test UnmarshalJSON
			var lt2 LocalTime
			err = json.Unmarshal(jsonData, &lt2)
			require.NoError(t, err)
			assert.Equal(t, tt.time, lt2)
		})
	}
}

func TestLocalTime_CompareConsistency(t *testing.T) {
	times := []LocalTime{
		MustNewLocalTime(0, 0, 0, 0),
		MustNewLocalTime(6, 0, 0, 0),
		MustNewLocalTime(12, 0, 0, 0),
		MustNewLocalTime(12, 30, 0, 0),
		MustNewLocalTime(12, 30, 30, 0),
		MustNewLocalTime(12, 30, 30, 500000000),
		MustNewLocalTime(18, 0, 0, 0),
		MustNewLocalTime(23, 59, 59, 999999999),
	}

	// Test that times are ordered correctly
	for i := 0; i < len(times)-1; i++ {
		assert.True(t, times[i].IsBefore(times[i+1]), "times[%d] should be before times[%d]", i, i+1)
		assert.False(t, times[i].IsAfter(times[i+1]), "times[%d] should not be after times[%d]", i, i+1)
		assert.Equal(t, -1, times[i].Compare(times[i+1]), "times[%d] should compare as -1 to times[%d]", i, i+1)
	}

	// Test equality
	for i, lt := range times {
		assert.Equal(t, 0, lt.Compare(lt), "time should equal itself")
		assert.False(t, lt.IsBefore(lt), "time should not be before itself")
		assert.False(t, lt.IsAfter(lt), "time should not be after itself")

		// Create copy and test
		copy := MustNewLocalTime(lt.Hour(), lt.Minute(), lt.Second(), lt.Nanosecond())
		assert.Equal(t, 0, lt.Compare(copy), "times[%d] should equal its copy", i)
	}
}

func TestLocalTimeNow(t *testing.T) {
	// Test that LocalTimeNow() returns a valid time
	now := LocalTimeNow()
	assert.False(t, now.IsZero(), "LocalTimeNow should not be zero")

	// Test that it's reasonable (between midnight and end of day)
	assert.True(t, now.Hour() >= 0 && now.Hour() < 24, "Hour should be valid")
	assert.True(t, now.Minute() >= 0 && now.Minute() < 60, "Minute should be valid")
	assert.True(t, now.Second() >= 0 && now.Second() < 60, "Second should be valid")
}

func TestLocalTimeNowUTC(t *testing.T) {
	nowUTC := LocalTimeNowUTC()
	assert.False(t, nowUTC.IsZero(), "LocalTimeNowUTC should not be zero")

	// Test that it's reasonable
	assert.True(t, nowUTC.Hour() >= 0 && nowUTC.Hour() < 24, "Hour should be valid")
	assert.True(t, nowUTC.Minute() >= 0 && nowUTC.Minute() < 60, "Minute should be valid")
	assert.True(t, nowUTC.Second() >= 0 && nowUTC.Second() < 60, "Second should be valid")
}

func TestParseLocalTime(t *testing.T) {
	t.Run("valid times", func(t *testing.T) {
		time, err := ParseLocalTime("14:30:45")
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalTime(14, 30, 45, 0), time)

		time, err = ParseLocalTime("09:05:07")
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalTime(9, 5, 7, 0), time)

		time, err = ParseLocalTime("00:00:00")
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalTime(0, 0, 0, 0), time)

		time, err = ParseLocalTime("23:59:59")
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalTime(23, 59, 59, 0), time)
	})

	t.Run("with fractional seconds", func(t *testing.T) {
		time, err := ParseLocalTime("14:30:45.123")
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalTime(14, 30, 45, 123000000), time)

		time, err = ParseLocalTime("14:30:45.123456")
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalTime(14, 30, 45, 123456000), time)

		time, err = ParseLocalTime("14:30:45.123456789")
		require.NoError(t, err)
		assert.Equal(t, MustNewLocalTime(14, 30, 45, 123456789), time)
	})

	t.Run("invalid format", func(t *testing.T) {
		_, err := ParseLocalTime("14-30-45")
		assert.Error(t, err)

		_, err = ParseLocalTime("14:30")
		assert.Error(t, err)

		_, err = ParseLocalTime("not-a-time")
		assert.Error(t, err)
	})

	t.Run("invalid values", func(t *testing.T) {
		_, err := ParseLocalTime("24:00:00")
		assert.Error(t, err)

		_, err = ParseLocalTime("23:60:00")
		assert.Error(t, err)

		_, err = ParseLocalTime("23:59:60")
		assert.Error(t, err)
	})

	t.Run("empty string", func(t *testing.T) {
		time, err := ParseLocalTime("")
		require.NoError(t, err)
		assert.True(t, time.IsZero())
	})
}

func TestMustParseLocalTime(t *testing.T) {
	t.Run("valid time", func(t *testing.T) {
		assert.NotPanics(t, func() {
			time := MustParseLocalTime("14:30:45.123456789")
			assert.Equal(t, MustNewLocalTime(14, 30, 45, 123456789), time)
		})
	})

	t.Run("invalid time panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustParseLocalTime("24:00:00")
		})

		assert.Panics(t, func() {
			MustParseLocalTime("invalid")
		})
	})
}

func TestLocalTimeNowIn(t *testing.T) {
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
			nowIn := LocalTimeNowIn(tt.loc)
			assert.False(t, nowIn.IsZero(), "LocalTimeNowIn should not be zero")

			// Test that it's reasonable
			assert.True(t, nowIn.Hour() >= 0 && nowIn.Hour() < 24, "Hour should be valid")
			assert.True(t, nowIn.Minute() >= 0 && nowIn.Minute() < 60, "Minute should be valid")
			assert.True(t, nowIn.Second() >= 0 && nowIn.Second() < 60, "Second should be valid")
		})
	}
}

func TestLocalTime_ValuePostgres(t *testing.T) {
	var pg = GetPG(t)
	t.Run("normal", func(t *testing.T) {
		var expected = MustParseLocalTime("12:00:00")
		var actual LocalTime
		var expectedTrue bool
		var e = pg.QueryRow("SELECT $1::time without time zone, $1::time without time zone = '12:00:00'", expected).Scan(&actual, &expectedTrue)
		assert.NoError(t, e)
		assert.Equal(t, expected, actual)
		assert.True(t, expectedTrue)
	})
	t.Run("null_value", func(t *testing.T) {
		var actual LocalTime
		var expectedTrue bool
		var e = pg.QueryRow("SELECT NULL::time without time zone, $1::time without time zone is null", actual).Scan(&actual, &expectedTrue)
		assert.NoError(t, e)
		assert.True(t, actual.IsZero())
		assert.True(t, expectedTrue)
	})
}

func TestLocalTime_ValueMySQL(t *testing.T) {
	var mysql = GetMySQL(t)
	t.Run("normal", func(t *testing.T) {
		var expected = MustParseLocalTime("08:00:00")
		var actual LocalTime
		var expectedTrue bool
		// I don't want to understand why MySQL has this bug
		var e = mysql.QueryRow("SELECT cast(cast(? as char) as time), cast(cast(? as char) as time) = '08:00:00'", expected, expected).Scan(&actual, &expectedTrue)
		assert.NoError(t, e)
		t.Log(expected.Value())
		assert.Equal(t, expected, actual)
		assert.True(t, expectedTrue)
	})
	t.Run("null_value", func(t *testing.T) {
		var actual LocalTime
		var expectedTrue bool
		// I don't want to understand why MySQL has this bug
		var e = mysql.QueryRow("SELECT cast(cast(? as char) as time), cast(cast(? as char) as time) is null", actual, actual).Scan(&actual, &expectedTrue)
		assert.NoError(t, e)
		assert.True(t, actual.IsZero())
		assert.True(t, expectedTrue)
	})
}

func TestLocalTime_PlusHours(t *testing.T) {
	t.Run("add positive hours", func(t *testing.T) {
		// Normal addition
		lt := MustNewLocalTime(10, 30, 45, 123456789)
		result := lt.PlusHours(2)
		assert.Equal(t, MustNewLocalTime(12, 30, 45, 123456789), result)

		// Add 1 hour
		lt = MustNewLocalTime(14, 0, 0, 0)
		result = lt.PlusHours(1)
		assert.Equal(t, MustNewLocalTime(15, 0, 0, 0), result)
	})

	t.Run("wrap around midnight", func(t *testing.T) {
		// 23:00 + 2 hours = 01:00
		lt := MustNewLocalTime(23, 0, 0, 0)
		result := lt.PlusHours(2)
		assert.Equal(t, MustNewLocalTime(1, 0, 0, 0), result)

		// 22:30 + 3 hours = 01:30
		lt = MustNewLocalTime(22, 30, 45, 123)
		result = lt.PlusHours(3)
		assert.Equal(t, MustNewLocalTime(1, 30, 45, 123), result)

		// Add exactly 24 hours (full day wrap)
		lt = MustNewLocalTime(10, 30, 0, 0)
		result = lt.PlusHours(24)
		assert.Equal(t, lt, result)
	})

	t.Run("add negative hours", func(t *testing.T) {
		// Subtract hours
		lt := MustNewLocalTime(10, 30, 45, 123456789)
		result := lt.PlusHours(-2)
		assert.Equal(t, MustNewLocalTime(8, 30, 45, 123456789), result)

		// Wrap backwards past midnight
		lt = MustNewLocalTime(1, 0, 0, 0)
		result = lt.PlusHours(-2)
		assert.Equal(t, MustNewLocalTime(23, 0, 0, 0), result)
	})

	t.Run("large values", func(t *testing.T) {
		// Add many hours (multiple days)
		lt := MustNewLocalTime(10, 0, 0, 0)
		result := lt.PlusHours(50) // 2 days + 2 hours
		assert.Equal(t, MustNewLocalTime(12, 0, 0, 0), result)

		// Subtract many hours
		lt = MustNewLocalTime(10, 0, 0, 0)
		result = lt.PlusHours(-50) // Go back 2 days + 2 hours
		assert.Equal(t, MustNewLocalTime(8, 0, 0, 0), result)
	})

	t.Run("zero value", func(t *testing.T) {
		var zero LocalTime
		result := zero.PlusHours(5)
		assert.True(t, result.IsZero())
	})
}

func TestLocalTime_PlusMinutes(t *testing.T) {
	t.Run("add positive minutes", func(t *testing.T) {
		// Normal addition
		lt := MustNewLocalTime(10, 30, 45, 123456789)
		result := lt.PlusMinutes(15)
		assert.Equal(t, MustNewLocalTime(10, 45, 45, 123456789), result)

		// Add minutes that roll over to next hour
		lt = MustNewLocalTime(14, 50, 0, 0)
		result = lt.PlusMinutes(20)
		assert.Equal(t, MustNewLocalTime(15, 10, 0, 0), result)
	})

	t.Run("wrap around midnight", func(t *testing.T) {
		// 23:50 + 20 minutes = 00:10
		lt := MustNewLocalTime(23, 50, 0, 0)
		result := lt.PlusMinutes(20)
		assert.Equal(t, MustNewLocalTime(0, 10, 0, 0), result)

		// 23:59 + 1 minute = 00:00
		lt = MustNewLocalTime(23, 59, 30, 500)
		result = lt.PlusMinutes(1)
		assert.Equal(t, MustNewLocalTime(0, 0, 30, 500), result)
	})

	t.Run("add negative minutes", func(t *testing.T) {
		// Subtract minutes
		lt := MustNewLocalTime(10, 30, 45, 123456789)
		result := lt.PlusMinutes(-15)
		assert.Equal(t, MustNewLocalTime(10, 15, 45, 123456789), result)

		// Wrap backwards past midnight
		lt = MustNewLocalTime(0, 10, 0, 0)
		result = lt.PlusMinutes(-20)
		assert.Equal(t, MustNewLocalTime(23, 50, 0, 0), result)
	})

	t.Run("large values", func(t *testing.T) {
		// Add many minutes (multiple days)
		lt := MustNewLocalTime(10, 0, 0, 0)
		result := lt.PlusMinutes(1500) // 25 hours = 1 day + 1 hour
		assert.Equal(t, MustNewLocalTime(11, 0, 0, 0), result)
	})

	t.Run("zero value", func(t *testing.T) {
		var zero LocalTime
		result := zero.PlusMinutes(30)
		assert.True(t, result.IsZero())
	})
}

func TestLocalTime_PlusSeconds(t *testing.T) {
	t.Run("add positive seconds", func(t *testing.T) {
		// Normal addition
		lt := MustNewLocalTime(10, 30, 45, 123456789)
		result := lt.PlusSeconds(10)
		assert.Equal(t, MustNewLocalTime(10, 30, 55, 123456789), result)

		// Add seconds that roll over to next minute
		lt = MustNewLocalTime(14, 30, 50, 0)
		result = lt.PlusSeconds(20)
		assert.Equal(t, MustNewLocalTime(14, 31, 10, 0), result)
	})

	t.Run("wrap around midnight", func(t *testing.T) {
		// 23:59:50 + 20 seconds = 00:00:10
		lt := MustNewLocalTime(23, 59, 50, 0)
		result := lt.PlusSeconds(20)
		assert.Equal(t, MustNewLocalTime(0, 0, 10, 0), result)

		// 23:59:59 + 1 second = 00:00:00
		lt = MustNewLocalTime(23, 59, 59, 999)
		result = lt.PlusSeconds(1)
		assert.Equal(t, MustNewLocalTime(0, 0, 0, 999), result)
	})

	t.Run("add negative seconds", func(t *testing.T) {
		// Subtract seconds
		lt := MustNewLocalTime(10, 30, 45, 123456789)
		result := lt.PlusSeconds(-10)
		assert.Equal(t, MustNewLocalTime(10, 30, 35, 123456789), result)

		// Wrap backwards past midnight
		lt = MustNewLocalTime(0, 0, 10, 0)
		result = lt.PlusSeconds(-20)
		assert.Equal(t, MustNewLocalTime(23, 59, 50, 0), result)
	})

	t.Run("large values", func(t *testing.T) {
		// Add many seconds (multiple days)
		lt := MustNewLocalTime(10, 0, 0, 0)
		result := lt.PlusSeconds(90000) // 25 hours in seconds
		assert.Equal(t, MustNewLocalTime(11, 0, 0, 0), result)
	})

	t.Run("zero value", func(t *testing.T) {
		var zero LocalTime
		result := zero.PlusSeconds(30)
		assert.True(t, result.IsZero())
	})
}

func TestLocalTime_PlusNanos(t *testing.T) {
	t.Run("add positive nanoseconds", func(t *testing.T) {
		// Normal addition
		lt := MustNewLocalTime(10, 30, 45, 123456789)
		result := lt.PlusNanos(100)
		assert.Equal(t, MustNewLocalTime(10, 30, 45, 123456889), result)

		// Add nanoseconds that roll over to next second
		lt = MustNewLocalTime(14, 30, 50, 999999999)
		result = lt.PlusNanos(1)
		assert.Equal(t, MustNewLocalTime(14, 30, 51, 0), result)
	})

	t.Run("wrap around midnight", func(t *testing.T) {
		// Add nanoseconds that wrap to next day
		lt := MustNewLocalTime(23, 59, 59, 999999999)
		result := lt.PlusNanos(1)
		assert.Equal(t, MustNewLocalTime(0, 0, 0, 0), result)
	})

	t.Run("add negative nanoseconds", func(t *testing.T) {
		// Subtract nanoseconds
		lt := MustNewLocalTime(10, 30, 45, 123456789)
		result := lt.PlusNanos(-100)
		assert.Equal(t, MustNewLocalTime(10, 30, 45, 123456689), result)

		// Wrap backwards past midnight
		lt = MustNewLocalTime(0, 0, 0, 0)
		result = lt.PlusNanos(-1)
		assert.Equal(t, MustNewLocalTime(23, 59, 59, 999999999), result)
	})

	t.Run("add milliseconds as nanoseconds", func(t *testing.T) {
		lt := MustNewLocalTime(10, 0, 0, 0)
		result := lt.PlusNanos(int64(500 * time.Millisecond))
		assert.Equal(t, MustNewLocalTime(10, 0, 0, 500000000), result)
	})

	t.Run("large values", func(t *testing.T) {
		// Add nanoseconds equivalent to multiple days
		lt := MustNewLocalTime(10, 0, 0, 0)
		result := lt.PlusNanos(int64(25 * time.Hour))
		assert.Equal(t, MustNewLocalTime(11, 0, 0, 0), result)
	})

	t.Run("zero value", func(t *testing.T) {
		var zero LocalTime
		result := zero.PlusNanos(1000000)
		assert.True(t, result.IsZero())
	})
}

func TestLocalTime_PlusMethodsConsistency(t *testing.T) {
	t.Run("hours and minutes equivalence", func(t *testing.T) {
		lt := MustNewLocalTime(10, 0, 0, 0)

		// 2 hours should equal 120 minutes
		resultHours := lt.PlusHours(2)
		resultMinutes := lt.PlusMinutes(120)
		assert.Equal(t, resultHours, resultMinutes)
	})

	t.Run("minutes and seconds equivalence", func(t *testing.T) {
		lt := MustNewLocalTime(10, 0, 0, 0)

		// 5 minutes should equal 300 seconds
		resultMinutes := lt.PlusMinutes(5)
		resultSeconds := lt.PlusSeconds(300)
		assert.Equal(t, resultMinutes, resultSeconds)
	})

	t.Run("seconds and nanoseconds equivalence", func(t *testing.T) {
		lt := MustNewLocalTime(10, 0, 0, 0)

		// 3 seconds should equal 3,000,000,000 nanoseconds
		resultSeconds := lt.PlusSeconds(3)
		resultNanos := lt.PlusNanos(3000000000)
		assert.Equal(t, resultSeconds, resultNanos)
	})

	t.Run("chaining operations", func(t *testing.T) {
		lt := MustNewLocalTime(10, 30, 45, 123456789)

		// Chain multiple operations
		result := lt.PlusHours(2).PlusMinutes(30).PlusSeconds(15).PlusNanos(100)
		expected := MustNewLocalTime(13, 1, 0, 123456889)
		assert.Equal(t, expected, result)
	})
}

func TestLocalTime_MinusHours(t *testing.T) {
	t.Run("subtract positive hours", func(t *testing.T) {
		// Normal subtraction
		lt := MustNewLocalTime(12, 30, 45, 123456789)
		result := lt.MinusHours(2)
		assert.Equal(t, MustNewLocalTime(10, 30, 45, 123456789), result)

		// Subtract 1 hour
		lt = MustNewLocalTime(14, 0, 0, 0)
		result = lt.MinusHours(1)
		assert.Equal(t, MustNewLocalTime(13, 0, 0, 0), result)
	})

	t.Run("wrap backwards past midnight", func(t *testing.T) {
		// 01:00 - 2 hours = 23:00
		lt := MustNewLocalTime(1, 0, 0, 0)
		result := lt.MinusHours(2)
		assert.Equal(t, MustNewLocalTime(23, 0, 0, 0), result)

		// 01:30 - 3 hours = 22:30
		lt = MustNewLocalTime(1, 30, 45, 123)
		result = lt.MinusHours(3)
		assert.Equal(t, MustNewLocalTime(22, 30, 45, 123), result)
	})

	t.Run("subtract negative hours (adds)", func(t *testing.T) {
		lt := MustNewLocalTime(10, 30, 45, 123456789)
		result := lt.MinusHours(-2)
		assert.Equal(t, MustNewLocalTime(12, 30, 45, 123456789), result)
	})

	t.Run("equivalence with PlusHours", func(t *testing.T) {
		lt := MustNewLocalTime(15, 30, 0, 0)

		// MinusHours(5) should equal PlusHours(-5)
		result1 := lt.MinusHours(5)
		result2 := lt.PlusHours(-5)
		assert.Equal(t, result1, result2)
	})

	t.Run("zero value", func(t *testing.T) {
		var zero LocalTime
		result := zero.MinusHours(5)
		assert.True(t, result.IsZero())
	})
}

func TestLocalTime_MinusMinutes(t *testing.T) {
	t.Run("subtract positive minutes", func(t *testing.T) {
		// Normal subtraction
		lt := MustNewLocalTime(10, 45, 45, 123456789)
		result := lt.MinusMinutes(15)
		assert.Equal(t, MustNewLocalTime(10, 30, 45, 123456789), result)

		// Subtract minutes that roll back to previous hour
		lt = MustNewLocalTime(15, 10, 0, 0)
		result = lt.MinusMinutes(20)
		assert.Equal(t, MustNewLocalTime(14, 50, 0, 0), result)
	})

	t.Run("wrap backwards past midnight", func(t *testing.T) {
		// 00:10 - 20 minutes = 23:50
		lt := MustNewLocalTime(0, 10, 0, 0)
		result := lt.MinusMinutes(20)
		assert.Equal(t, MustNewLocalTime(23, 50, 0, 0), result)
	})

	t.Run("equivalence with PlusMinutes", func(t *testing.T) {
		lt := MustNewLocalTime(10, 30, 45, 123456789)

		// MinusMinutes(30) should equal PlusMinutes(-30)
		result1 := lt.MinusMinutes(30)
		result2 := lt.PlusMinutes(-30)
		assert.Equal(t, result1, result2)
	})

	t.Run("zero value", func(t *testing.T) {
		var zero LocalTime
		result := zero.MinusMinutes(30)
		assert.True(t, result.IsZero())
	})
}

func TestLocalTime_MinusSeconds(t *testing.T) {
	t.Run("subtract positive seconds", func(t *testing.T) {
		// Normal subtraction
		lt := MustNewLocalTime(10, 30, 55, 123456789)
		result := lt.MinusSeconds(10)
		assert.Equal(t, MustNewLocalTime(10, 30, 45, 123456789), result)

		// Subtract seconds that roll back to previous minute
		lt = MustNewLocalTime(14, 31, 10, 0)
		result = lt.MinusSeconds(20)
		assert.Equal(t, MustNewLocalTime(14, 30, 50, 0), result)
	})

	t.Run("wrap backwards past midnight", func(t *testing.T) {
		// 00:00:10 - 20 seconds = 23:59:50
		lt := MustNewLocalTime(0, 0, 10, 0)
		result := lt.MinusSeconds(20)
		assert.Equal(t, MustNewLocalTime(23, 59, 50, 0), result)
	})

	t.Run("equivalence with PlusSeconds", func(t *testing.T) {
		lt := MustNewLocalTime(10, 30, 45, 123456789)

		// MinusSeconds(45) should equal PlusSeconds(-45)
		result1 := lt.MinusSeconds(45)
		result2 := lt.PlusSeconds(-45)
		assert.Equal(t, result1, result2)
	})

	t.Run("zero value", func(t *testing.T) {
		var zero LocalTime
		result := zero.MinusSeconds(30)
		assert.True(t, result.IsZero())
	})
}

func TestLocalTime_MinusNanos(t *testing.T) {
	t.Run("subtract positive nanoseconds", func(t *testing.T) {
		// Normal subtraction
		lt := MustNewLocalTime(10, 30, 45, 123456889)
		result := lt.MinusNanos(100)
		assert.Equal(t, MustNewLocalTime(10, 30, 45, 123456789), result)

		// Subtract nanoseconds that roll back to previous second
		lt = MustNewLocalTime(14, 30, 51, 0)
		result = lt.MinusNanos(1)
		assert.Equal(t, MustNewLocalTime(14, 30, 50, 999999999), result)
	})

	t.Run("wrap backwards past midnight", func(t *testing.T) {
		// 00:00:00.0 - 1 nanosecond = 23:59:59.999999999
		lt := MustNewLocalTime(0, 0, 0, 0)
		result := lt.MinusNanos(1)
		assert.Equal(t, MustNewLocalTime(23, 59, 59, 999999999), result)
	})

	t.Run("equivalence with PlusNanos", func(t *testing.T) {
		lt := MustNewLocalTime(10, 30, 45, 123456789)

		// MinusNanos(1000) should equal PlusNanos(-1000)
		result1 := lt.MinusNanos(1000)
		result2 := lt.PlusNanos(-1000)
		assert.Equal(t, result1, result2)
	})

	t.Run("zero value", func(t *testing.T) {
		var zero LocalTime
		result := zero.MinusNanos(1000000)
		assert.True(t, result.IsZero())
	})
}

func TestLocalTime_PlusMinusCombinations(t *testing.T) {
	t.Run("plus and minus cancel out", func(t *testing.T) {
		lt := MustNewLocalTime(10, 30, 45, 123456789)

		// Add then subtract same amount - should return to original
		result := lt.PlusHours(5).MinusHours(5)
		assert.Equal(t, lt, result)

		result = lt.PlusMinutes(100).MinusMinutes(100)
		assert.Equal(t, lt, result)

		result = lt.PlusSeconds(500).MinusSeconds(500)
		assert.Equal(t, lt, result)

		result = lt.PlusNanos(1000000000).MinusNanos(1000000000)
		assert.Equal(t, lt, result)
	})

	t.Run("complex chaining", func(t *testing.T) {
		lt := MustNewLocalTime(12, 0, 0, 0)

		// Complex chain: +3h -1h +30m -15m +45s -30s
		result := lt.PlusHours(3).MinusHours(1).PlusMinutes(30).MinusMinutes(15).PlusSeconds(45).MinusSeconds(30)
		// Net: +2h +15m +15s = 14:15:15
		expected := MustNewLocalTime(14, 15, 15, 0)
		assert.Equal(t, expected, result)
	})
}
