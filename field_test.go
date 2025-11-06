package goda

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestField_String(t *testing.T) {
	tests := []struct {
		field    Field
		expected string
	}{
		// Time fields
		{NanoOfSecond, "NanoOfSecond"},
		{NanoOfDay, "NanoOfDay"},
		{MicroOfSecond, "MicroOfSecond"},
		{MicroOfDay, "MicroOfDay"},
		{MilliOfSecond, "MilliOfSecond"},
		{MilliOfDay, "MilliOfDay"},
		{SecondOfMinute, "SecondOfMinute"},
		{SecondOfDay, "SecondOfDay"},
		{MinuteOfHour, "MinuteOfHour"},
		{MinuteOfDay, "MinuteOfDay"},
		{HourOfAmPm, "HourOfAmPm"},
		{ClockHourOfAmPm, "ClockHourOfAmPm"},
		{HourOfDay, "HourOfDay"},
		{ClockHourOfDay, "ClockHourOfDay"},
		{AmPmOfDay, "AmPmOfDay"},

		// Date fields
		{DayOfWeekField, "DayOfWeek"},
		{AlignedDayOfWeekInMonth, "AlignedDayOfWeekInMonth"},
		{AlignedDayOfWeekInYear, "AlignedDayOfWeekInYear"},
		{DayOfMonth, "DayOfMonth"},
		{DayOfYear, "DayOfYear"},
		{EpochDay, "EpochDay"},
		{AlignedWeekOfMonth, "AlignedWeekOfMonth"},
		{AlignedWeekOfYear, "AlignedWeekOfYear"},
		{MonthOfYear, "MonthOfYear"},
		{ProlepticMonth, "ProlepticMonth"},
		{YearOfEra, "YearOfEra"},
		{YearField, "Year"},
		{Era, "Era"},

		// Other fields
		{InstantSeconds, "InstantSeconds"},
		{OffsetSeconds, "OffsetSeconds"},

		// Zero value
		{Field(0), ""},

		// Unknown field
		{Field(9999), ""},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.field.String())
		})
	}
}

func TestLocalDate_IsSupportedField(t *testing.T) {
	date := MustNewLocalDate(2024, March, 15)

	// Supported fields
	assert.True(t, date.IsSupportedField(DayOfWeekField))
	assert.True(t, date.IsSupportedField(DayOfMonth))
	assert.True(t, date.IsSupportedField(DayOfYear))
	assert.True(t, date.IsSupportedField(EpochDay))
	assert.True(t, date.IsSupportedField(MonthOfYear))
	assert.True(t, date.IsSupportedField(ProlepticMonth))
	assert.True(t, date.IsSupportedField(YearOfEra))
	assert.True(t, date.IsSupportedField(YearField))
	assert.True(t, date.IsSupportedField(Era))

	// Unsupported fields (time fields)
	assert.False(t, date.IsSupportedField(HourOfDay))
	assert.False(t, date.IsSupportedField(MinuteOfHour))
	assert.False(t, date.IsSupportedField(SecondOfMinute))
	assert.False(t, date.IsSupportedField(NanoOfSecond))
}

func TestLocalDate_GetFieldInt64(t *testing.T) {
	date := MustNewLocalDate(2024, March, 15) // Friday

	// All fields
	assert.Equal(t, int64(5), date.GetFieldInt64(DayOfWeekField)) // Friday = 5
	assert.Equal(t, int64(15), date.GetFieldInt64(DayOfMonth))
	assert.Equal(t, int64(75), date.GetFieldInt64(DayOfYear)) // 31+29+15 = 75 (2024 is leap year)
	assert.Equal(t, int64(3), date.GetFieldInt64(MonthOfYear))
	assert.Equal(t, int64(2024), date.GetFieldInt64(YearField))
	assert.Equal(t, int64(2024), date.GetFieldInt64(YearOfEra))
	assert.Equal(t, int64(1), date.GetFieldInt64(Era)) // CE

	// EpochDay
	epochDays := date.GetFieldInt64(EpochDay)
	assert.Greater(t, epochDays, int64(0))

	// ProlepticMonth: 2024*12 + 3 - 1 = 24291
	prolepticMonth := date.GetFieldInt64(ProlepticMonth)
	assert.Equal(t, int64(2024*12+3-1), prolepticMonth)

	// Zero date
	var zeroDate LocalDate
	assert.Equal(t, int64(0), zeroDate.GetFieldInt64(DayOfMonth))
}

func TestLocalTime_IsSupportedField(t *testing.T) {
	time := MustNewLocalTime(14, 30, 45, 123456789)

	// Supported fields
	assert.True(t, time.IsSupportedField(NanoOfSecond))
	assert.True(t, time.IsSupportedField(NanoOfDay))
	assert.True(t, time.IsSupportedField(MicroOfSecond))
	assert.True(t, time.IsSupportedField(MicroOfDay))
	assert.True(t, time.IsSupportedField(MilliOfSecond))
	assert.True(t, time.IsSupportedField(MilliOfDay))
	assert.True(t, time.IsSupportedField(SecondOfMinute))
	assert.True(t, time.IsSupportedField(SecondOfDay))
	assert.True(t, time.IsSupportedField(MinuteOfHour))
	assert.True(t, time.IsSupportedField(MinuteOfDay))
	assert.True(t, time.IsSupportedField(HourOfAmPm))
	assert.True(t, time.IsSupportedField(ClockHourOfAmPm))
	assert.True(t, time.IsSupportedField(HourOfDay))
	assert.True(t, time.IsSupportedField(ClockHourOfDay))
	assert.True(t, time.IsSupportedField(AmPmOfDay))

	// Unsupported fields (date fields)
	assert.False(t, time.IsSupportedField(DayOfMonth))
	assert.False(t, time.IsSupportedField(MonthOfYear))
	assert.False(t, time.IsSupportedField(YearField))
}

func TestLocalTime_GetFieldInt64(t *testing.T) {
	// Test 14:30:45.123456789 (PM)
	time := MustNewLocalTime(14, 30, 45, 123456789)

	assert.Equal(t, int64(123456789), time.GetFieldInt64(NanoOfSecond))
	assert.Equal(t, int64(123456), time.GetFieldInt64(MicroOfSecond))
	assert.Equal(t, int64(123), time.GetFieldInt64(MilliOfSecond))
	assert.Equal(t, int64(45), time.GetFieldInt64(SecondOfMinute))
	assert.Equal(t, int64(30), time.GetFieldInt64(MinuteOfHour))
	assert.Equal(t, int64(14), time.GetFieldInt64(HourOfDay))
	assert.Equal(t, int64(14), time.GetFieldInt64(ClockHourOfDay))
	assert.Equal(t, int64(2), time.GetFieldInt64(HourOfAmPm))      // 14 % 12 = 2
	assert.Equal(t, int64(2), time.GetFieldInt64(ClockHourOfAmPm)) // same as HourOfAmPm for 14
	assert.Equal(t, int64(1), time.GetFieldInt64(AmPmOfDay))       // PM

	// Calculate expected values for *OfDay fields
	nanoOfDay := int64(14)*int64(3600000000000) + int64(30)*int64(60000000000) + int64(45)*int64(1000000000) + 123456789
	assert.Equal(t, nanoOfDay, time.GetFieldInt64(NanoOfDay))
	assert.Equal(t, nanoOfDay/1000, time.GetFieldInt64(MicroOfDay))
	assert.Equal(t, nanoOfDay/1000000, time.GetFieldInt64(MilliOfDay))
	assert.Equal(t, int64(14*3600+30*60+45), time.GetFieldInt64(SecondOfDay))
	assert.Equal(t, int64(14*60+30), time.GetFieldInt64(MinuteOfDay))

	// Test midnight (00:00:00)
	midnight := MustNewLocalTime(0, 0, 0, 0)
	assert.Equal(t, int64(0), midnight.GetFieldInt64(HourOfDay))
	assert.Equal(t, int64(24), midnight.GetFieldInt64(ClockHourOfDay)) // 24 for midnight
	assert.Equal(t, int64(0), midnight.GetFieldInt64(HourOfAmPm))
	assert.Equal(t, int64(12), midnight.GetFieldInt64(ClockHourOfAmPm)) // 12 for midnight
	assert.Equal(t, int64(0), midnight.GetFieldInt64(AmPmOfDay))        // AM

	// Test noon (12:00:00)
	noon := MustNewLocalTime(12, 0, 0, 0)
	assert.Equal(t, int64(12), noon.GetFieldInt64(HourOfDay))
	assert.Equal(t, int64(12), noon.GetFieldInt64(ClockHourOfDay))
	assert.Equal(t, int64(0), noon.GetFieldInt64(HourOfAmPm))       // 12 % 12 = 0
	assert.Equal(t, int64(12), noon.GetFieldInt64(ClockHourOfAmPm)) // 12 for noon
	assert.Equal(t, int64(1), noon.GetFieldInt64(AmPmOfDay))        // PM

	// Zero time
	var zeroTime LocalTime
	assert.Equal(t, int64(0), zeroTime.GetFieldInt64(NanoOfDay))
}

func TestLocalDateTime_IsSupportedField(t *testing.T) {
	dt := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 123456789)

	// Should support both date and time fields
	assert.True(t, dt.IsSupportedField(DayOfMonth))
	assert.True(t, dt.IsSupportedField(MonthOfYear))
	assert.True(t, dt.IsSupportedField(YearField))
	assert.True(t, dt.IsSupportedField(HourOfDay))
	assert.True(t, dt.IsSupportedField(MinuteOfHour))
	assert.True(t, dt.IsSupportedField(SecondOfMinute))
	assert.True(t, dt.IsSupportedField(NanoOfSecond))

	// Unsupported fields
	assert.False(t, dt.IsSupportedField(InstantSeconds))
	assert.False(t, dt.IsSupportedField(OffsetSeconds))
}

func TestLocalDateTime_GetFieldInt64(t *testing.T) {
	dt := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 123456789)

	// Date fields
	assert.Equal(t, int64(15), dt.GetFieldInt64(DayOfMonth))
	assert.Equal(t, int64(3), dt.GetFieldInt64(MonthOfYear))
	assert.Equal(t, int64(2024), dt.GetFieldInt64(YearField))

	epochDays := dt.GetFieldInt64(EpochDay)
	assert.Greater(t, epochDays, int64(0))

	// Time fields
	assert.Equal(t, int64(14), dt.GetFieldInt64(HourOfDay))
	assert.Equal(t, int64(30), dt.GetFieldInt64(MinuteOfHour))
	assert.Equal(t, int64(45), dt.GetFieldInt64(SecondOfMinute))
	assert.Equal(t, int64(123456789), dt.GetFieldInt64(NanoOfSecond))

	nanoOfDay := dt.GetFieldInt64(NanoOfDay)
	assert.Greater(t, nanoOfDay, int64(0))
}
