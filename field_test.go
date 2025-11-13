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

func TestLocalDate_GetField(t *testing.T) {
	date := MustNewLocalDate(2024, March, 15) // Friday

	t.Run("supported fields", func(t *testing.T) {
		// Test DayOfWeekField
		dayOfWeek := date.GetField(DayOfWeekField)
		assert.True(t, dayOfWeek.Valid())
		assert.False(t, dayOfWeek.Unsupported())
		assert.Equal(t, int64(5), dayOfWeek.Int64()) // Friday = 5

		// Test DayOfMonth
		dayOfMonth := date.GetField(DayOfMonth)
		assert.True(t, dayOfMonth.Valid())
		assert.Equal(t, 15, dayOfMonth.Int())

		// Test DayOfYear
		dayOfYear := date.GetField(DayOfYear)
		assert.True(t, dayOfYear.Valid())
		assert.Equal(t, int64(75), dayOfYear.Int64()) // 31+29+15 = 75 (2024 is leap year)

		// Test MonthOfYear
		month := date.GetField(MonthOfYear)
		assert.True(t, month.Valid())
		assert.Equal(t, int64(3), month.Int64())

		// Test YearField
		year := date.GetField(YearField)
		assert.True(t, year.Valid())
		assert.Equal(t, int64(2024), year.Int64())

		// Test YearOfEra
		yearOfEra := date.GetField(YearOfEra)
		assert.True(t, yearOfEra.Valid())
		assert.Equal(t, int64(2024), yearOfEra.Int64())

		// Test Era
		era := date.GetField(Era)
		assert.True(t, era.Valid())
		assert.Equal(t, int64(1), era.Int64()) // CE

		// Test EpochDay
		epochDay := date.GetField(EpochDay)
		assert.True(t, epochDay.Valid())
		assert.Greater(t, epochDay.Int64(), int64(0))

		// Test ProlepticMonth
		prolepticMonth := date.GetField(ProlepticMonth)
		assert.True(t, prolepticMonth.Valid())
		assert.Equal(t, int64(2024*12+3-1), prolepticMonth.Int64())
	})

	t.Run("unsupported fields", func(t *testing.T) {
		// Time fields are not supported
		hourOfDay := date.GetField(HourOfDay)
		assert.False(t, hourOfDay.Valid())
		assert.True(t, hourOfDay.Unsupported())

		nanoOfSecond := date.GetField(NanoOfSecond)
		assert.False(t, nanoOfSecond.Valid())
		assert.True(t, nanoOfSecond.Unsupported())
	})

	t.Run("zero date", func(t *testing.T) {
		var zeroDate LocalDate
		field := zeroDate.GetField(DayOfMonth)
		assert.False(t, field.Valid())
		assert.True(t, field.Unsupported())
	})

	t.Run("BCE date", func(t *testing.T) {
		bceDate := MustNewLocalDate(-100, January, 1)
		era := bceDate.GetField(Era)
		assert.True(t, era.Valid())
		assert.Equal(t, int64(0), era.Int64()) // BCE
	})
}

func TestLocalTime_GetField(t *testing.T) {
	// Test 14:30:45.123456789 (PM)
	time := MustNewLocalTime(14, 30, 45, 123456789)

	t.Run("supported fields", func(t *testing.T) {
		// Test NanoOfSecond
		nanoOfSecond := time.GetField(NanoOfSecond)
		assert.True(t, nanoOfSecond.Valid())
		assert.False(t, nanoOfSecond.Unsupported())
		assert.Equal(t, int64(123456789), nanoOfSecond.Int64())

		// Test MicroOfSecond
		microOfSecond := time.GetField(MicroOfSecond)
		assert.True(t, microOfSecond.Valid())
		assert.Equal(t, int64(123456), microOfSecond.Int64())

		// Test MilliOfSecond
		milliOfSecond := time.GetField(MilliOfSecond)
		assert.True(t, milliOfSecond.Valid())
		assert.Equal(t, int64(123), milliOfSecond.Int64())

		// Test SecondOfMinute
		secondOfMinute := time.GetField(SecondOfMinute)
		assert.True(t, secondOfMinute.Valid())
		assert.Equal(t, int64(45), secondOfMinute.Int64())

		// Test MinuteOfHour
		minuteOfHour := time.GetField(MinuteOfHour)
		assert.True(t, minuteOfHour.Valid())
		assert.Equal(t, int64(30), minuteOfHour.Int64())

		// Test HourOfDay
		hourOfDay := time.GetField(HourOfDay)
		assert.True(t, hourOfDay.Valid())
		assert.Equal(t, int64(14), hourOfDay.Int64())

		// Test ClockHourOfDay
		clockHourOfDay := time.GetField(ClockHourOfDay)
		assert.True(t, clockHourOfDay.Valid())
		assert.Equal(t, int64(14), clockHourOfDay.Int64())

		// Test HourOfAmPm
		hourOfAmPm := time.GetField(HourOfAmPm)
		assert.True(t, hourOfAmPm.Valid())
		assert.Equal(t, int64(2), hourOfAmPm.Int64()) // 14 % 12 = 2

		// Test ClockHourOfAmPm
		clockHourOfAmPm := time.GetField(ClockHourOfAmPm)
		assert.True(t, clockHourOfAmPm.Valid())
		assert.Equal(t, int64(2), clockHourOfAmPm.Int64())

		// Test AmPmOfDay
		amPmOfDay := time.GetField(AmPmOfDay)
		assert.True(t, amPmOfDay.Valid())
		assert.Equal(t, int64(1), amPmOfDay.Int64()) // PM

		// Test NanoOfDay
		nanoOfDay := time.GetField(NanoOfDay)
		assert.True(t, nanoOfDay.Valid())
		expectedNanoOfDay := int64(14)*int64(3600000000000) + int64(30)*int64(60000000000) + int64(45)*int64(1000000000) + 123456789
		assert.Equal(t, expectedNanoOfDay, nanoOfDay.Int64())

		// Test MicroOfDay
		microOfDay := time.GetField(MicroOfDay)
		assert.True(t, microOfDay.Valid())
		assert.Equal(t, expectedNanoOfDay/1000, microOfDay.Int64())

		// Test MilliOfDay
		milliOfDay := time.GetField(MilliOfDay)
		assert.True(t, milliOfDay.Valid())
		assert.Equal(t, expectedNanoOfDay/1000000, milliOfDay.Int64())

		// Test SecondOfDay
		secondOfDay := time.GetField(SecondOfDay)
		assert.True(t, secondOfDay.Valid())
		assert.Equal(t, int64(14*3600+30*60+45), secondOfDay.Int64())

		// Test MinuteOfDay
		minuteOfDay := time.GetField(MinuteOfDay)
		assert.True(t, minuteOfDay.Valid())
		assert.Equal(t, int64(14*60+30), minuteOfDay.Int64())
	})

	t.Run("midnight special cases", func(t *testing.T) {
		midnight := MustNewLocalTime(0, 0, 0, 0)

		// ClockHourOfDay at midnight should be 24
		clockHourOfDay := midnight.GetField(ClockHourOfDay)
		assert.True(t, clockHourOfDay.Valid())
		assert.Equal(t, int64(24), clockHourOfDay.Int64())

		// ClockHourOfAmPm at midnight should be 12
		clockHourOfAmPm := midnight.GetField(ClockHourOfAmPm)
		assert.True(t, clockHourOfAmPm.Valid())
		assert.Equal(t, int64(12), clockHourOfAmPm.Int64())

		// AmPmOfDay at midnight should be 0 (AM)
		amPmOfDay := midnight.GetField(AmPmOfDay)
		assert.True(t, amPmOfDay.Valid())
		assert.Equal(t, int64(0), amPmOfDay.Int64())
	})

	t.Run("noon special cases", func(t *testing.T) {
		noon := MustNewLocalTime(12, 0, 0, 0)

		// ClockHourOfDay at noon should be 12
		clockHourOfDay := noon.GetField(ClockHourOfDay)
		assert.True(t, clockHourOfDay.Valid())
		assert.Equal(t, int64(12), clockHourOfDay.Int64())

		// ClockHourOfAmPm at noon should be 12
		clockHourOfAmPm := noon.GetField(ClockHourOfAmPm)
		assert.True(t, clockHourOfAmPm.Valid())
		assert.Equal(t, int64(12), clockHourOfAmPm.Int64())

		// AmPmOfDay at noon should be 1 (PM)
		amPmOfDay := noon.GetField(AmPmOfDay)
		assert.True(t, amPmOfDay.Valid())
		assert.Equal(t, int64(1), amPmOfDay.Int64())
	})

	t.Run("unsupported fields", func(t *testing.T) {
		// Date fields are not supported
		dayOfMonth := time.GetField(DayOfMonth)
		assert.False(t, dayOfMonth.Valid())
		assert.True(t, dayOfMonth.Unsupported())

		monthOfYear := time.GetField(MonthOfYear)
		assert.False(t, monthOfYear.Valid())
		assert.True(t, monthOfYear.Unsupported())
	})

	t.Run("zero time", func(t *testing.T) {
		var zeroTime LocalTime
		field := zeroTime.GetField(HourOfDay)
		assert.False(t, field.Valid())
		assert.True(t, field.Unsupported())
	})
}

func TestLocalDateTime_GetField(t *testing.T) {
	dt := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 123456789)

	t.Run("date fields", func(t *testing.T) {
		// Test DayOfWeekField
		dayOfWeek := dt.GetField(DayOfWeekField)
		assert.True(t, dayOfWeek.Valid())
		assert.Equal(t, int64(5), dayOfWeek.Int64()) // Friday

		// Test DayOfMonth
		dayOfMonth := dt.GetField(DayOfMonth)
		assert.True(t, dayOfMonth.Valid())
		assert.Equal(t, int64(15), dayOfMonth.Int64())

		// Test DayOfYear
		dayOfYear := dt.GetField(DayOfYear)
		assert.True(t, dayOfYear.Valid())
		assert.Equal(t, int64(75), dayOfYear.Int64())

		// Test MonthOfYear
		month := dt.GetField(MonthOfYear)
		assert.True(t, month.Valid())
		assert.Equal(t, int64(3), month.Int64())

		// Test YearField
		year := dt.GetField(YearField)
		assert.True(t, year.Valid())
		assert.Equal(t, int64(2024), year.Int64())

		// Test Era
		era := dt.GetField(Era)
		assert.True(t, era.Valid())
		assert.Equal(t, int64(1), era.Int64()) // CE

		// Test EpochDay
		epochDay := dt.GetField(EpochDay)
		assert.True(t, epochDay.Valid())
		assert.Greater(t, epochDay.Int64(), int64(0))

		// Test ProlepticMonth
		prolepticMonth := dt.GetField(ProlepticMonth)
		assert.True(t, prolepticMonth.Valid())
		assert.Equal(t, int64(2024*12+3-1), prolepticMonth.Int64())
	})

	t.Run("time fields", func(t *testing.T) {
		// Test NanoOfSecond
		nanoOfSecond := dt.GetField(NanoOfSecond)
		assert.True(t, nanoOfSecond.Valid())
		assert.Equal(t, int64(123456789), nanoOfSecond.Int64())

		// Test MicroOfSecond
		microOfSecond := dt.GetField(MicroOfSecond)
		assert.True(t, microOfSecond.Valid())
		assert.Equal(t, int64(123456), microOfSecond.Int64())

		// Test MilliOfSecond
		milliOfSecond := dt.GetField(MilliOfSecond)
		assert.True(t, milliOfSecond.Valid())
		assert.Equal(t, int64(123), milliOfSecond.Int64())

		// Test SecondOfMinute
		secondOfMinute := dt.GetField(SecondOfMinute)
		assert.True(t, secondOfMinute.Valid())
		assert.Equal(t, int64(45), secondOfMinute.Int64())

		// Test MinuteOfHour
		minuteOfHour := dt.GetField(MinuteOfHour)
		assert.True(t, minuteOfHour.Valid())
		assert.Equal(t, int64(30), minuteOfHour.Int64())

		// Test HourOfDay
		hourOfDay := dt.GetField(HourOfDay)
		assert.True(t, hourOfDay.Valid())
		assert.Equal(t, int64(14), hourOfDay.Int64())

		// Test ClockHourOfDay
		clockHourOfDay := dt.GetField(ClockHourOfDay)
		assert.True(t, clockHourOfDay.Valid())
		assert.Equal(t, int64(14), clockHourOfDay.Int64())

		// Test HourOfAmPm
		hourOfAmPm := dt.GetField(HourOfAmPm)
		assert.True(t, hourOfAmPm.Valid())
		assert.Equal(t, int64(2), hourOfAmPm.Int64())

		// Test AmPmOfDay
		amPmOfDay := dt.GetField(AmPmOfDay)
		assert.True(t, amPmOfDay.Valid())
		assert.Equal(t, int64(1), amPmOfDay.Int64()) // PM

		// Test NanoOfDay
		nanoOfDay := dt.GetField(NanoOfDay)
		assert.True(t, nanoOfDay.Valid())
		expectedNanoOfDay := int64(14)*int64(3600000000000) + int64(30)*int64(60000000000) + int64(45)*int64(1000000000) + 123456789
		assert.Equal(t, expectedNanoOfDay, nanoOfDay.Int64())
	})

	t.Run("unsupported fields", func(t *testing.T) {
		// InstantSeconds and OffsetSeconds are not supported
		instantSeconds := dt.GetField(InstantSeconds)
		assert.False(t, instantSeconds.Valid())
		assert.True(t, instantSeconds.Unsupported())

		offsetSeconds := dt.GetField(OffsetSeconds)
		assert.False(t, offsetSeconds.Valid())
		assert.True(t, offsetSeconds.Unsupported())
	})

	t.Run("zero datetime", func(t *testing.T) {
		var zeroDT LocalDateTime
		field := zeroDT.GetField(HourOfDay)
		assert.False(t, field.Valid())
		assert.True(t, field.Unsupported())

		field = zeroDT.GetField(DayOfMonth)
		assert.False(t, field.Valid())
		assert.True(t, field.Unsupported())
	})

	t.Run("delegation to date and time", func(t *testing.T) {
		// Verify that date fields are delegated to LocalDate
		dateField := dt.GetField(YearField)
		dateDirectField := dt.LocalDate().GetField(YearField)
		assert.Equal(t, dateDirectField.Int64(), dateField.Int64())

		// Verify that time fields are delegated to LocalTime
		timeField := dt.GetField(HourOfDay)
		timeDirectField := dt.LocalTime().GetField(HourOfDay)
		assert.Equal(t, timeDirectField.Int64(), timeField.Int64())
	})
}
