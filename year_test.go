package goda

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYear_IsLeapYear(t *testing.T) {
	tests := []struct {
		year   Year
		isLeap bool
		reason string
	}{
		// Basic leap year tests
		{2024, true, "divisible by 4, not by 100"},
		{2023, false, "not divisible by 4"},
		{2022, false, "not divisible by 4"},
		{2020, true, "divisible by 4, not by 100"},

		// Century years (divisible by 100)
		{2000, true, "divisible by 400"},
		{1900, false, "divisible by 100, not by 400"},
		{1800, false, "divisible by 100, not by 400"},
		{2100, false, "divisible by 100, not by 400"},
		{2200, false, "divisible by 100, not by 400"},
		{2400, true, "divisible by 400"},

		// Edge cases
		{1600, true, "divisible by 400"},
		{1700, false, "divisible by 100, not by 400"},
		{4, true, "small leap year"},
		{1, false, "year 1"},
		{0, true, "year 0 (divisible by 400)"},

		// Negative years
		{-4, true, "negative leap year"},
		{-1, false, "negative non-leap year"},
		{-100, false, "negative century, not leap"},
		{-400, true, "negative, divisible by 400"},

		// Large years
		{10000, true, "large year divisible by 400"},
		{9999, false, "large year not divisible by 4"},
	}

	for _, tt := range tests {
		t.Run(tt.reason, func(t *testing.T) {
			result := tt.year.IsLeapYear()
			assert.Equal(t, tt.isLeap, result, "Year %d: %s", tt.year, tt.reason)
		})
	}
}

func TestYear_Length(t *testing.T) {
	tests := []struct {
		year   Year
		length int
	}{
		{2024, 366}, // leap year
		{2023, 365}, // non-leap year
		{2000, 366}, // century leap year
		{1900, 365}, // century non-leap year
		{2020, 366}, // leap year
		{2021, 365}, // non-leap year
	}

	for _, tt := range tests {
		result := tt.year.Length()
		assert.Equal(t, tt.length, result, "Year %d should have %d days", tt.year, tt.length)
	}
}

func TestYear_Int(t *testing.T) {
	tests := []Year{0, 1, 2024, -1, -2024, 9999, -9999}

	for _, year := range tests {
		result := year.Int()
		assert.Equal(t, int(year), result, "Year %d", year)
	}
}

func TestYear_Int64(t *testing.T) {
	tests := []Year{0, 1, 2024, -1, -2024, 9999, -9999}

	for _, year := range tests {
		result := year.Int64()
		assert.Equal(t, int64(year), result, "Year %d", year)
	}
}

func TestYear_String(t *testing.T) {
	tests := []struct {
		year     Year
		expected string
	}{
		{2024, "2024"},
		{2000, "2000"},
		{1999, "1999"},
		{1, "0001"},
		{0, "0000"},
		{-1, "-0001"},
		{-2024, "-2024"},
		{9999, "9999"},
		{10000, "10000"},
		{-9999, "-9999"},
		{-10000, "-10000"},
	}

	for _, tt := range tests {
		result := tt.year.String()
		assert.Equal(t, tt.expected, result, "Year %d", tt.year)
	}
}

func TestYear_AppendText(t *testing.T) {
	tests := []struct {
		year     Year
		prefix   string
		expected string
	}{
		{2024, "Year: ", "Year: 2024"},
		{1, "", "0001"},
		{0, "", "0000"},
		{-1, "", "-0001"},
		{9999, "", "9999"},
		{10000, "Big: ", "Big: 10000"},
		{-10000, "Negative: ", "Negative: -10000"},
	}

	for _, tt := range tests {
		buf := []byte(tt.prefix)
		buf, err := tt.year.AppendText(buf)
		require.NoError(t, err)
		assert.Equal(t, tt.expected, string(buf), "Year %d with prefix '%s'", tt.year, tt.prefix)
	}
}

func TestYear_LeapYearPattern(t *testing.T) {
	// Test the 400-year cycle pattern
	// In 400 years, there are 97 leap years:
	// - 100 years divisible by 4 = 100 leap years
	// - minus 4 century years = 96 leap years
	// - plus 1 year divisible by 400 = 97 leap years

	leapCount := 0
	for y := Year(0); y < 400; y++ {
		if y.IsLeapYear() {
			leapCount++
		}
	}
	assert.Equal(t, 97, leapCount, "400-year cycle should have exactly 97 leap years")

	// Test that the pattern repeats
	year1 := Year(100)
	year2 := Year(100 + 400)
	assert.Equal(t, year1.IsLeapYear(), year2.IsLeapYear(), "Leap year pattern should repeat every 400 years")
}
