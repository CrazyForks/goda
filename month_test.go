package goda

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMonth_IsZero(t *testing.T) {
	var zero Month
	assert.True(t, zero.IsZero())

	assert.False(t, January.IsZero())
	assert.False(t, February.IsZero())
	assert.False(t, December.IsZero())
}

func TestMonth_MaxDays(t *testing.T) {
	tests := []struct {
		month   Month
		maxDays int
	}{
		{January, 31},
		{February, 29}, // Note: MaxDays returns 29 for February regardless of leap year
		{March, 31},
		{April, 30},
		{May, 31},
		{June, 30},
		{July, 31},
		{August, 31},
		{September, 30},
		{October, 31},
		{November, 30},
		{December, 31},
	}

	for _, tt := range tests {
		result := tt.month.MaxDays()
		assert.Equal(t, tt.maxDays, result, "Month %s should have max %d days", tt.month, tt.maxDays)
	}
}

func TestMonth_MaxDays_Zero(t *testing.T) {
	var zero Month
	assert.Equal(t, 0, zero.MaxDays())
}

func TestMonth_Length(t *testing.T) {
	tests := []struct {
		month  Month
		isLeap bool
		length int
	}{
		// January - always 31 days
		{January, false, 31},
		{January, true, 31},

		// February - depends on leap year
		{February, false, 28},
		{February, true, 29},

		// March - always 31 days
		{March, false, 31},
		{March, true, 31},

		// April - always 30 days
		{April, false, 30},
		{April, true, 30},

		// May - always 31 days
		{May, false, 31},
		{May, true, 31},

		// June - always 30 days
		{June, false, 30},
		{June, true, 30},

		// July - always 31 days
		{July, false, 31},
		{July, true, 31},

		// August - always 31 days
		{August, false, 31},
		{August, true, 31},

		// September - always 30 days
		{September, false, 30},
		{September, true, 30},

		// October - always 31 days
		{October, false, 31},
		{October, true, 31},

		// November - always 30 days
		{November, false, 30},
		{November, true, 30},

		// December - always 31 days
		{December, false, 31},
		{December, true, 31},
	}

	for _, tt := range tests {
		result := tt.month.Length(tt.isLeap)
		leapStr := "non-leap"
		if tt.isLeap {
			leapStr = "leap"
		}
		assert.Equal(t, tt.length, result, "Month %s in %s year should have %d days", tt.month, leapStr, tt.length)
	}
}

func TestMonth_FirstDayOfYear(t *testing.T) {
	tests := []struct {
		month    Month
		isLeap   bool
		firstDay int
	}{
		// Non-leap year
		{January, false, 1},
		{February, false, 32},
		{March, false, 60},
		{April, false, 91},
		{May, false, 121},
		{June, false, 152},
		{July, false, 182},
		{August, false, 213},
		{September, false, 244},
		{October, false, 274},
		{November, false, 305},
		{December, false, 335},

		// Leap year (March onwards should be +1)
		{January, true, 1},
		{February, true, 32},
		{March, true, 61},
		{April, true, 92},
		{May, true, 122},
		{June, true, 153},
		{July, true, 183},
		{August, true, 214},
		{September, true, 245},
		{October, true, 275},
		{November, true, 306},
		{December, true, 336},
	}

	for _, tt := range tests {
		result := tt.month.FirstDayOfYear(tt.isLeap)
		leapStr := "non-leap"
		if tt.isLeap {
			leapStr = "leap"
		}
		assert.Equal(t, tt.firstDay, result, "Month %s in %s year should start on day %d", tt.month, leapStr, tt.firstDay)
	}
}

func TestMonth_FirstDayOfYear_Zero(t *testing.T) {
	var zero Month
	assert.Equal(t, 0, zero.FirstDayOfYear(false))
	assert.Equal(t, 0, zero.FirstDayOfYear(true))
}

func TestMonth_String(t *testing.T) {
	tests := []struct {
		month    Month
		expected string
	}{
		{January, "January"},
		{February, "February"},
		{March, "March"},
		{April, "April"},
		{May, "May"},
		{June, "June"},
		{July, "July"},
		{August, "August"},
		{September, "September"},
		{October, "October"},
		{November, "November"},
		{December, "December"},
	}

	for _, tt := range tests {
		result := tt.month.String()
		assert.Equal(t, tt.expected, result, "Month %d", tt.month)
	}
}

func TestMonth_String_Zero(t *testing.T) {
	var zero Month
	assert.Equal(t, "", zero.String())
}

func TestMonth_ConsistencyWithGoTimeMonth(t *testing.T) {
	// Verify that our Month values match Go's time.Month
	for m := January; m <= December; m++ {
		goMonth := time.Month(m)
		assert.Equal(t, goMonth.String(), m.String(), "Month %d", m)
	}
}

func TestMonth_AllMonthsCovered(t *testing.T) {
	// Ensure all 12 months are defined
	months := []Month{
		January, February, March, April, May, June,
		July, August, September, October, November, December,
	}

	assert.Len(t, months, 12, "Should have exactly 12 months")

	// Verify they are sequential
	for i, month := range months {
		assert.Equal(t, Month(i+1), month, "Month at index %d should be %d", i, i+1)
	}
}

func TestMonth_LeapYearFebruaryLength(t *testing.T) {
	// Specific test for February in different years
	tests := []struct {
		year   Year
		length int
	}{
		{2024, 29}, // leap year
		{2023, 28}, // non-leap year
		{2000, 29}, // century leap year
		{1900, 28}, // century non-leap year
		{2020, 29}, // leap year
		{2100, 28}, // non-leap century year
	}

	for _, tt := range tests {
		result := February.Length(tt.year.IsLeapYear())
		assert.Equal(t, tt.length, result, "February in year %d should have %d days", tt.year, tt.length)
	}
}

func TestMonth_YearDaysConsistency(t *testing.T) {
	// Test that the sum of all month lengths equals the year length
	t.Run("non-leap year", func(t *testing.T) {
		totalDays := 0
		for m := January; m <= December; m++ {
			totalDays += m.Length(false)
		}
		assert.Equal(t, 365, totalDays, "Non-leap year should have 365 days")
	})

	t.Run("leap year", func(t *testing.T) {
		totalDays := 0
		for m := January; m <= December; m++ {
			totalDays += m.Length(true)
		}
		assert.Equal(t, 366, totalDays, "Leap year should have 366 days")
	})
}

func TestMonth_FirstDayConsistency(t *testing.T) {
	// Test that FirstDayOfYear values are consistent with month lengths
	t.Run("non-leap year", func(t *testing.T) {
		expectedDay := 1
		for m := January; m <= December; m++ {
			assert.Equal(t, expectedDay, m.FirstDayOfYear(false), "Month %s", m)
			expectedDay += m.Length(false)
		}
	})

	t.Run("leap year", func(t *testing.T) {
		expectedDay := 1
		for m := January; m <= December; m++ {
			assert.Equal(t, expectedDay, m.FirstDayOfYear(true), "Month %s", m)
			expectedDay += m.Length(true)
		}
	})
}
