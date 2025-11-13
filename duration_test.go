package goda

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDurationOfSeconds(t *testing.T) {
	t.Run("basic creation", func(t *testing.T) {
		d := NewDurationOfSeconds(10, 500000000)
		assert.Equal(t, int64(10), d.seconds)
		assert.Equal(t, int32(500000000), d.nanos)
	})

	t.Run("nano overflow positive", func(t *testing.T) {
		// 1 second + 2 billion nanos = 3 seconds
		d := NewDurationOfSeconds(1, 2000000000)
		assert.Equal(t, int64(3), d.seconds)
		assert.Equal(t, int32(0), d.nanos)
	})

	t.Run("nano overflow negative", func(t *testing.T) {
		// 5 seconds + (-1.5 seconds in nanos) = 3.5 seconds
		d := NewDurationOfSeconds(5, -1500000000)
		assert.Equal(t, int64(3), d.seconds)
		assert.Equal(t, int32(500000000), d.nanos)
	})

	t.Run("nano underflow", func(t *testing.T) {
		// -1 second + (-500 million nanos) = -1.5 seconds
		d := NewDurationOfSeconds(-1, -500000000)
		assert.Equal(t, int64(-2), d.seconds)
		assert.Equal(t, int32(500000000), d.nanos)
	})

	t.Run("zero", func(t *testing.T) {
		d := NewDurationOfSeconds(0, 0)
		assert.True(t, d.IsZero())
		assert.False(t, d.IsPositive())
		assert.False(t, d.IsNegative())
	})

	t.Run("wrapping behavior - large positive nano adjustment", func(t *testing.T) {
		// Test that large nano adjustments wrap correctly
		// 10 billion nanos should be 10 seconds
		d := NewDurationOfSeconds(0, 10000000000)
		assert.Equal(t, int64(10), d.seconds)
		assert.Equal(t, int32(0), d.nanos)
	})

	t.Run("wrapping behavior - large negative nano adjustment", func(t *testing.T) {
		// -10 billion nanos should be -10 seconds
		d := NewDurationOfSeconds(0, -10000000000)
		assert.Equal(t, int64(-10), d.seconds)
		assert.Equal(t, int32(0), d.nanos)
	})

	t.Run("wrapping with mixed signs", func(t *testing.T) {
		// 10 seconds + (-5.5 billion nanos) = 4.5 seconds
		d := NewDurationOfSeconds(10, -5500000000)
		assert.Equal(t, int64(4), d.seconds)
		assert.Equal(t, int32(500000000), d.nanos)
	})
}

func TestNewDurationByGoDuration(t *testing.T) {
	t.Run("positive duration", func(t *testing.T) {
		goDur := 5*time.Second + 500*time.Millisecond
		d := NewDurationByGoDuration(goDur)
		assert.Equal(t, int64(5), d.seconds)
		assert.Equal(t, int32(500000000), d.nanos)
	})

	t.Run("negative duration", func(t *testing.T) {
		goDur := -(3*time.Second + 250*time.Millisecond)
		d := NewDurationByGoDuration(goDur)
		assert.Equal(t, int64(-4), d.seconds)
		assert.Equal(t, int32(750000000), d.nanos) // wraps around
	})

	t.Run("zero duration", func(t *testing.T) {
		d := NewDurationByGoDuration(0)
		assert.True(t, d.IsZero())
	})
}

func TestDuration_IsZero(t *testing.T) {
	assert.True(t, Duration{}.IsZero())
	assert.True(t, NewDurationOfSeconds(0, 0).IsZero())
	assert.False(t, NewDurationOfSeconds(1, 0).IsZero())
	assert.False(t, NewDurationOfSeconds(0, 1).IsZero())
	assert.False(t, NewDurationOfSeconds(-1, 0).IsZero())
}

func TestDuration_IsPositive(t *testing.T) {
	assert.True(t, NewDurationOfSeconds(1, 0).IsPositive())
	assert.True(t, NewDurationOfSeconds(0, 1).IsPositive())
	assert.False(t, NewDurationOfSeconds(0, 0).IsPositive())
	assert.False(t, NewDurationOfSeconds(-1, 0).IsPositive())
}

func TestDuration_IsNegative(t *testing.T) {
	assert.True(t, NewDurationOfSeconds(-1, 0).IsNegative())
	assert.False(t, NewDurationOfSeconds(0, 0).IsNegative())
	assert.False(t, NewDurationOfSeconds(1, 0).IsNegative())
	assert.False(t, NewDurationOfSeconds(0, 1).IsNegative())
}

func TestDuration_Plus(t *testing.T) {
	t.Run("simple addition", func(t *testing.T) {
		d1 := NewDurationOfSeconds(5, 0)
		d2 := NewDurationOfSeconds(3, 0)
		result := d1.Plus(d2)
		assert.Equal(t, int64(8), result.seconds)
		assert.Equal(t, int32(0), result.nanos)
	})

	t.Run("addition with nano overflow", func(t *testing.T) {
		d1 := NewDurationOfSeconds(5, 700000000)
		d2 := NewDurationOfSeconds(3, 500000000)
		result := d1.Plus(d2)
		assert.Equal(t, int64(9), result.seconds)
		assert.Equal(t, int32(200000000), result.nanos)
	})

	t.Run("addition with negative", func(t *testing.T) {
		d1 := NewDurationOfSeconds(5, 0)
		d2 := NewDurationOfSeconds(-3, 0)
		result := d1.Plus(d2)
		assert.Equal(t, int64(2), result.seconds)
		assert.Equal(t, int32(0), result.nanos)
	})

	t.Run("wrapping behavior - nano overflow", func(t *testing.T) {
		d1 := NewDurationOfSeconds(0, 999999999)
		d2 := NewDurationOfSeconds(0, 2)
		result := d1.Plus(d2)
		assert.Equal(t, int64(1), result.seconds)
		assert.Equal(t, int32(1), result.nanos)
	})

	t.Run("wrapping behavior - crossing zero", func(t *testing.T) {
		d1 := NewDurationOfSeconds(3, 200000000)
		d2 := NewDurationOfSeconds(-3, -500000000)
		result := d1.Plus(d2)
		assert.Equal(t, int64(-1), result.seconds)
		assert.Equal(t, int32(700000000), result.nanos)
	})
}

func TestDuration_Minus(t *testing.T) {
	t.Run("simple subtraction", func(t *testing.T) {
		d1 := NewDurationOfSeconds(5, 0)
		d2 := NewDurationOfSeconds(3, 0)
		result := d1.Minus(d2)
		assert.Equal(t, int64(2), result.seconds)
		assert.Equal(t, int32(0), result.nanos)
	})

	t.Run("subtraction with nano underflow", func(t *testing.T) {
		d1 := NewDurationOfSeconds(5, 200000000)
		d2 := NewDurationOfSeconds(3, 500000000)
		result := d1.Minus(d2)
		assert.Equal(t, int64(1), result.seconds)
		assert.Equal(t, int32(700000000), result.nanos)
	})

	t.Run("wrapping behavior - result becomes negative", func(t *testing.T) {
		d1 := NewDurationOfSeconds(3, 0)
		d2 := NewDurationOfSeconds(5, 0)
		result := d1.Minus(d2)
		assert.Equal(t, int64(-2), result.seconds)
		assert.Equal(t, int32(0), result.nanos)
	})
}

func TestDuration_Negated(t *testing.T) {
	t.Run("negate positive", func(t *testing.T) {
		d := NewDurationOfSeconds(5, 500000000)
		neg := d.Negated()
		assert.Equal(t, int64(-6), neg.seconds)
		assert.Equal(t, int32(500000000), neg.nanos)
	})

	t.Run("negate negative", func(t *testing.T) {
		d := NewDurationOfSeconds(-5, 500000000)
		neg := d.Negated()
		assert.Equal(t, int64(4), neg.seconds)
		assert.Equal(t, int32(500000000), neg.nanos)
	})

	t.Run("negate zero", func(t *testing.T) {
		d := NewDurationOfSeconds(0, 0)
		neg := d.Negated()
		assert.True(t, neg.IsZero())
	})

	t.Run("double negation", func(t *testing.T) {
		d := NewDurationOfSeconds(3, 250000000)
		neg := d.Negated().Negated()
		assert.Equal(t, d.seconds, neg.seconds)
		assert.Equal(t, d.nanos, neg.nanos)
	})
}

func TestDuration_Abs(t *testing.T) {
	t.Run("abs of positive", func(t *testing.T) {
		d := NewDurationOfSeconds(5, 0)
		abs := d.Abs()
		assert.Equal(t, d, abs)
	})

	t.Run("abs of negative", func(t *testing.T) {
		d := NewDurationOfSeconds(-5, 500000000)
		abs := d.Abs()
		assert.Equal(t, int64(4), abs.seconds)
		assert.Equal(t, int32(500000000), abs.nanos)
		assert.True(t, abs.IsPositive())
	})

	t.Run("abs of zero", func(t *testing.T) {
		d := NewDurationOfSeconds(0, 0)
		abs := d.Abs()
		assert.True(t, abs.IsZero())
	})
}

func TestDuration_Compare(t *testing.T) {
	t.Run("equal durations", func(t *testing.T) {
		d1 := NewDurationOfSeconds(5, 500000000)
		d2 := NewDurationOfSeconds(5, 500000000)
		assert.Equal(t, 0, d1.Compare(d2))
	})

	t.Run("compare seconds", func(t *testing.T) {
		d1 := NewDurationOfSeconds(5, 0)
		d2 := NewDurationOfSeconds(3, 0)
		assert.Greater(t, d1.Compare(d2), 0)
		assert.Less(t, d2.Compare(d1), 0)
	})

	t.Run("compare nanos", func(t *testing.T) {
		d1 := NewDurationOfSeconds(5, 600000000)
		d2 := NewDurationOfSeconds(5, 500000000)
		assert.Greater(t, d1.Compare(d2), 0)
		assert.Less(t, d2.Compare(d1), 0)
	})

	t.Run("compare with negatives", func(t *testing.T) {
		d1 := NewDurationOfSeconds(-5, 0)
		d2 := NewDurationOfSeconds(3, 0)
		assert.Less(t, d1.Compare(d2), 0)
		assert.Greater(t, d2.Compare(d1), 0)
	})
}

func TestParseDuration(t *testing.T) {
	t.Run("valid formats", func(t *testing.T) {
		tests := []struct {
			input   string
			seconds int64
			nanos   int32
		}{
			{"PT0S", 0, 0},
			{"PT1S", 1, 0},
			{"PT1M", 60, 0},
			{"PT1H", 3600, 0},
			{"PT1H30M", 5400, 0},
			{"PT1H30M45S", 5445, 0},
			{"PT1H30M45.5S", 5445, 500000000},
			{"PT1.5S", 1, 500000000},
			{"PT0.123456789S", 0, 123456789},
			{"PT-1S", -1, 0},
			{"PT-1H30M", -5400, 0},
			{"PT8H6M12.345S", 29172, 345000000},
		}

		for _, tt := range tests {
			t.Run(tt.input, func(t *testing.T) {
				d, err := ParseDuration(tt.input)
				require.NoError(t, err)
				assert.Equal(t, tt.seconds, d.seconds, "seconds mismatch")
				assert.Equal(t, tt.nanos, d.nanos, "nanos mismatch")
			})
		}
	})

	t.Run("invalid formats", func(t *testing.T) {
		tests := []string{
			"",
			"P1D",      // days not supported in time duration
			"1H",       // missing PT prefix
			"T1H",      // missing P prefix
			"PT",       // empty after PT
			"PTX",      // invalid character
			"PT1H2X",   // invalid unit
			"PT1.2.3S", // multiple decimal points
		}

		for _, input := range tests {
			t.Run(input, func(t *testing.T) {
				_, err := ParseDuration(input)
				assert.Error(t, err)
			})
		}
	})

	t.Run("round trip", func(t *testing.T) {
		original := NewDurationOfSeconds(29172, 345000000)
		str := original.String()
		parsed, err := ParseDuration(str)
		require.NoError(t, err)
		assert.Equal(t, original, parsed)
	})
}

func TestMustParseDuration(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		assert.NotPanics(t, func() {
			d := MustParseDuration("PT1H30M")
			assert.Equal(t, int64(5400), d.seconds)
		})
	})

	t.Run("invalid input panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustParseDuration("invalid")
		})
	})
}

func TestDuration_String(t *testing.T) {
	tests := []struct {
		name     string
		duration Duration
		expected string
	}{
		{"zero", NewDurationOfSeconds(0, 0), "PT0S"},
		{"1 second", NewDurationOfSeconds(1, 0), "PT1S"},
		{"1 minute", NewDurationOfSeconds(60, 0), "PT1M"},
		{"1 hour", NewDurationOfSeconds(3600, 0), "PT1H"},
		{"complex", NewDurationOfSeconds(5445, 500000000), "PT1H30M45.5S"},
		{"only nanos", NewDurationOfSeconds(0, 123456789), "PT0.123456789S"},
		{"negative", NewDurationOfSeconds(-5445, 0), "PT-1H30M45S"},
		{"8h6m12.345s", NewDurationOfSeconds(29172, 345000000), "PT8H6M12.345S"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.duration.String())
		})
	}
}

func TestDuration_AppendText(t *testing.T) {
	t.Run("append to empty buffer", func(t *testing.T) {
		d := NewDurationOfSeconds(3661, 500000000)
		buf, err := d.AppendText(nil)
		require.NoError(t, err)
		assert.Equal(t, "PT1H1M1.5S", string(buf))
	})

	t.Run("append to existing buffer", func(t *testing.T) {
		d := NewDurationOfSeconds(60, 0)
		buf := []byte("Duration: ")
		buf, err := d.AppendText(buf)
		require.NoError(t, err)
		assert.Equal(t, "Duration: PT1M", string(buf))
	})
}

func TestDuration_MarshalText(t *testing.T) {
	d := NewDurationOfSeconds(3661, 500000000)
	text, err := d.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "PT1H1M1.5S", string(text))
}

func TestDuration_UnmarshalText(t *testing.T) {
	t.Run("valid text", func(t *testing.T) {
		var d Duration
		err := d.UnmarshalText([]byte("PT1H30M45S"))
		require.NoError(t, err)
		assert.Equal(t, int64(5445), d.seconds)
	})

	t.Run("invalid text", func(t *testing.T) {
		var d Duration
		err := d.UnmarshalText([]byte("invalid"))
		assert.Error(t, err)
	})
}

func TestDuration_MarshalJSON(t *testing.T) {
	d := NewDurationOfSeconds(3661, 500000000)
	jsonBytes, err := json.Marshal(d)
	require.NoError(t, err)
	assert.Equal(t, `"PT1H1M1.5S"`, string(jsonBytes))
}

func TestDuration_UnmarshalJSON(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		var d Duration
		err := json.Unmarshal([]byte(`"PT1H30M45S"`), &d)
		require.NoError(t, err)
		assert.Equal(t, int64(5445), d.seconds)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		var d Duration
		err := json.Unmarshal([]byte(`"invalid"`), &d)
		assert.Error(t, err)
	})

	t.Run("not a string", func(t *testing.T) {
		var d Duration
		err := json.Unmarshal([]byte(`123`), &d)
		assert.Error(t, err)
	})
}

func TestDuration_Scan(t *testing.T) {
	t.Run("scan string", func(t *testing.T) {
		var d Duration
		err := d.Scan("PT1H30M")
		require.NoError(t, err)
		assert.Equal(t, int64(5400), d.seconds)
	})

	t.Run("scan byte slice", func(t *testing.T) {
		var d Duration
		err := d.Scan([]byte("PT1H30M"))
		require.NoError(t, err)
		assert.Equal(t, int64(5400), d.seconds)
	})

	t.Run("scan int64 nanoseconds", func(t *testing.T) {
		var d Duration
		nanos := int64(5 * time.Second)
		err := d.Scan(nanos)
		require.NoError(t, err)
		assert.Equal(t, int64(5), d.seconds)
	})

	t.Run("scan nil", func(t *testing.T) {
		var d Duration
		err := d.Scan(nil)
		require.NoError(t, err)
		assert.True(t, d.IsZero())
	})

	t.Run("scan unsupported type", func(t *testing.T) {
		var d Duration
		err := d.Scan(123.45)
		assert.Error(t, err)
	})
}

func TestDuration_Value(t *testing.T) {
	d := NewDurationOfSeconds(3661, 500000000)
	val, err := d.Value()
	require.NoError(t, err)
	assert.Equal(t, "PT1H1M1.5S", val)
}

// Test overflow wrapping scenarios - this is the key difference from Java's java.time
func TestDuration_OverflowWrapping(t *testing.T) {
	t.Run("nano overflow wraps to seconds", func(t *testing.T) {
		// Adding 1.9 billion nanos to 0.5 seconds should wrap to 2.4 seconds
		d := NewDurationOfSeconds(0, 2400000000)
		assert.Equal(t, int64(2), d.seconds)
		assert.Equal(t, int32(400000000), d.nanos)
	})

	t.Run("negative nano wraps correctly", func(t *testing.T) {
		// 5 seconds - 2.5 billion nanos = 2.5 seconds
		d := NewDurationOfSeconds(5, -2500000000)
		assert.Equal(t, int64(2), d.seconds)
		assert.Equal(t, int32(500000000), d.nanos)
	})

	t.Run("extreme positive nano wrapping", func(t *testing.T) {
		// 100 billion nanos = 100 seconds
		d := NewDurationOfSeconds(0, 100000000000)
		assert.Equal(t, int64(100), d.seconds)
		assert.Equal(t, int32(0), d.nanos)
	})

	t.Run("extreme negative nano wrapping", func(t *testing.T) {
		// -100 billion nanos = -100 seconds
		d := NewDurationOfSeconds(0, -100000000000)
		assert.Equal(t, int64(-100), d.seconds)
		assert.Equal(t, int32(0), d.nanos)
	})

	t.Run("addition causing overflow wrap", func(t *testing.T) {
		d1 := NewDurationOfSeconds(10, 800000000)
		d2 := NewDurationOfSeconds(5, 300000000)
		result := d1.Plus(d2)
		// 10.8 + 5.3 = 16.1 seconds
		assert.Equal(t, int64(16), result.seconds)
		assert.Equal(t, int32(100000000), result.nanos)
	})

	t.Run("subtraction causing underflow wrap", func(t *testing.T) {
		d1 := NewDurationOfSeconds(10, 200000000)
		d2 := NewDurationOfSeconds(5, 700000000)
		result := d1.Minus(d2)
		// 10.2 - 5.7 = 4.5 seconds
		assert.Equal(t, int64(4), result.seconds)
		assert.Equal(t, int32(500000000), result.nanos)
	})

	t.Run("mixed sign operations", func(t *testing.T) {
		// Positive seconds + negative nanos
		d := NewDurationOfSeconds(10, -3000000000)
		// 10 - 3 = 7 seconds
		assert.Equal(t, int64(7), d.seconds)
		assert.Equal(t, int32(0), d.nanos)
	})

	t.Run("negative seconds + positive nanos", func(t *testing.T) {
		// -10 seconds + 3 billion nanos = -7 seconds
		d := NewDurationOfSeconds(-10, 3000000000)
		assert.Equal(t, int64(-7), d.seconds)
		assert.Equal(t, int32(0), d.nanos)
	})

	t.Run("wrapping preserves correctness", func(t *testing.T) {
		// Create duration with overflow
		d1 := NewDurationOfSeconds(5, 2000000000) // 7 seconds
		d2 := NewDurationOfSeconds(7, 0)          // 7 seconds
		// Both should be equal after wrapping
		assert.Equal(t, d1, d2)
	})

	t.Run("multiple wraps", func(t *testing.T) {
		// 50 billion nanos = 50 seconds
		d := NewDurationOfSeconds(10, 50000000000)
		// 10 + 50 = 60 seconds
		assert.Equal(t, int64(60), d.seconds)
		assert.Equal(t, int32(0), d.nanos)
	})

	t.Run("fractional negative wrapping", func(t *testing.T) {
		// 0 seconds - 0.5 billion nanos = -0.5 seconds
		d := NewDurationOfSeconds(0, -500000000)
		assert.Equal(t, int64(-1), d.seconds)
		assert.Equal(t, int32(500000000), d.nanos)
	})
}

// Test the Between methods for LocalTime and LocalDateTime
func TestLocalTime_Between(t *testing.T) {
	t.Run("positive difference", func(t *testing.T) {
		t1 := MustNewLocalTime(10, 0, 0, 0)
		t2 := MustNewLocalTime(10, 30, 0, 0)
		d := t1.Between(t2)
		assert.Equal(t, int64(1800), d.seconds) // 30 minutes
		assert.Equal(t, int32(0), d.nanos)
	})

	t.Run("negative difference", func(t *testing.T) {
		t1 := MustNewLocalTime(10, 30, 0, 0)
		t2 := MustNewLocalTime(10, 0, 0, 0)
		d := t1.Between(t2)
		assert.Equal(t, int64(-1800), d.seconds) // -30 minutes
		assert.Equal(t, int32(0), d.nanos)
	})

	t.Run("with nanoseconds", func(t *testing.T) {
		t1 := MustNewLocalTime(10, 0, 0, 500000000)
		t2 := MustNewLocalTime(10, 0, 1, 250000000)
		d := t1.Between(t2)
		assert.Equal(t, int64(0), d.seconds)
		assert.Equal(t, int32(750000000), d.nanos)
	})

	t.Run("zero time returns zero duration", func(t *testing.T) {
		t1 := LocalTime{}
		t2 := MustNewLocalTime(10, 0, 0, 0)
		d := t1.Between(t2)
		assert.True(t, d.IsZero())
	})

	t.Run("across day boundary concept", func(t *testing.T) {
		// LocalTime doesn't wrap around midnight, just calculates difference
		t1 := MustNewLocalTime(23, 30, 0, 0)
		t2 := MustNewLocalTime(0, 30, 0, 0)
		d := t1.Between(t2)
		// This gives negative duration as t2 is "earlier" in the day
		assert.True(t, d.IsNegative())
	})
}

func TestLocalDateTime_Between(t *testing.T) {
	t.Run("same day", func(t *testing.T) {
		dt1 := MustNewLocalDateTimeFromComponents(2024, March, 15, 10, 0, 0, 0)
		dt2 := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 0, 0)
		d := dt1.Between(dt2)
		assert.Equal(t, int64(16200), d.seconds) // 4.5 hours
		assert.Equal(t, int32(0), d.nanos)
	})

	t.Run("different days", func(t *testing.T) {
		dt1 := MustNewLocalDateTimeFromComponents(2024, March, 15, 10, 0, 0, 0)
		dt2 := MustNewLocalDateTimeFromComponents(2024, March, 16, 10, 0, 0, 0)
		d := dt1.Between(dt2)
		assert.Equal(t, int64(86400), d.seconds) // 24 hours
		assert.Equal(t, int32(0), d.nanos)
	})

	t.Run("different days with time", func(t *testing.T) {
		dt1 := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 0, 0)
		dt2 := MustNewLocalDateTimeFromComponents(2024, March, 16, 10, 0, 0, 0)
		d := dt1.Between(dt2)
		// From 14:30 to next day 10:00 = 19.5 hours = 70200 seconds
		assert.Equal(t, int64(70200), d.seconds)
		assert.Equal(t, int32(0), d.nanos)
	})

	t.Run("negative duration", func(t *testing.T) {
		dt1 := MustNewLocalDateTimeFromComponents(2024, March, 16, 10, 0, 0, 0)
		dt2 := MustNewLocalDateTimeFromComponents(2024, March, 15, 10, 0, 0, 0)
		d := dt1.Between(dt2)
		assert.Equal(t, int64(-86400), d.seconds) // -24 hours
		assert.Equal(t, int32(0), d.nanos)
	})

	t.Run("with nanoseconds", func(t *testing.T) {
		dt1 := MustNewLocalDateTimeFromComponents(2024, March, 15, 10, 0, 0, 500000000)
		dt2 := MustNewLocalDateTimeFromComponents(2024, March, 15, 10, 0, 1, 250000000)
		d := dt1.Between(dt2)
		assert.Equal(t, int64(0), d.seconds)
		assert.Equal(t, int32(750000000), d.nanos)
	})

	t.Run("zero datetime returns zero duration", func(t *testing.T) {
		dt1 := LocalDateTime{}
		dt2 := MustNewLocalDateTimeFromComponents(2024, March, 15, 10, 0, 0, 0)
		d := dt1.Between(dt2)
		assert.True(t, d.IsZero())
	})

	t.Run("across month boundary", func(t *testing.T) {
		dt1 := MustNewLocalDateTimeFromComponents(2024, February, 28, 12, 0, 0, 0)
		dt2 := MustNewLocalDateTimeFromComponents(2024, March, 1, 12, 0, 0, 0)
		d := dt1.Between(dt2)
		// 2024 is leap year, so Feb has 29 days: Feb 28->29->Mar 1 = 2 days
		assert.Equal(t, int64(172800), d.seconds) // 48 hours
		assert.Equal(t, int32(0), d.nanos)
	})
}

// Benchmark tests
func BenchmarkDuration_Plus(b *testing.B) {
	d1 := NewDurationOfSeconds(100, 500000000)
	d2 := NewDurationOfSeconds(50, 300000000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d1.Plus(d2)
	}
}

func BenchmarkDuration_Parse(b *testing.B) {
	s := "PT8H6M12.345S"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseDuration(s)
	}
}

func BenchmarkDuration_String(b *testing.B) {
	d := NewDurationOfSeconds(29172, 345000000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.String()
	}
}
