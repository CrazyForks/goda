package goda

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalDate_IsSupportedField(t *testing.T) {
	date := MustLocalDateOf(2024, March, 15)

	// Supported fields
	assert.True(t, date.IsSupportedField(FieldDayOfWeek))
	assert.True(t, date.IsSupportedField(FieldDayOfMonth))
	assert.True(t, date.IsSupportedField(FieldDayOfYear))
	assert.True(t, date.IsSupportedField(FieldEpochDay))
	assert.True(t, date.IsSupportedField(FieldMonthOfYear))
	assert.True(t, date.IsSupportedField(FieldProlepticMonth))
	assert.True(t, date.IsSupportedField(FieldYearOfEra))
	assert.True(t, date.IsSupportedField(FieldYear))
	assert.True(t, date.IsSupportedField(FieldEra))

	// Unsupported fields (time fields)
	assert.False(t, date.IsSupportedField(FieldHourOfDay))
	assert.False(t, date.IsSupportedField(FieldMinuteOfHour))
	assert.False(t, date.IsSupportedField(FieldSecondOfMinute))
	assert.False(t, date.IsSupportedField(FieldNanoOfSecond))
}

func TestLocalTime_IsSupportedField(t *testing.T) {
	time := MustLocalTimeOf(14, 30, 45, 123456789)

	// Supported fields
	assert.True(t, time.IsSupportedField(FieldNanoOfSecond))
	assert.True(t, time.IsSupportedField(FieldNanoOfDay))
	assert.True(t, time.IsSupportedField(FieldMicroOfSecond))
	assert.True(t, time.IsSupportedField(FieldMicroOfDay))
	assert.True(t, time.IsSupportedField(FieldMilliOfSecond))
	assert.True(t, time.IsSupportedField(FieldMilliOfDay))
	assert.True(t, time.IsSupportedField(FieldSecondOfMinute))
	assert.True(t, time.IsSupportedField(FieldSecondOfDay))
	assert.True(t, time.IsSupportedField(FieldMinuteOfHour))
	assert.True(t, time.IsSupportedField(FieldMinuteOfDay))
	assert.True(t, time.IsSupportedField(FieldHourOfAmPm))
	assert.True(t, time.IsSupportedField(FieldClockHourOfAmPm))
	assert.True(t, time.IsSupportedField(FieldHourOfDay))
	assert.True(t, time.IsSupportedField(FieldClockHourOfDay))
	assert.True(t, time.IsSupportedField(FieldAmPmOfDay))

	// Unsupported fields (date fields)
	assert.False(t, time.IsSupportedField(FieldDayOfMonth))
	assert.False(t, time.IsSupportedField(FieldMonthOfYear))
	assert.False(t, time.IsSupportedField(FieldYear))
}

func TestLocalDateTime_IsSupportedField(t *testing.T) {
	dt := MustLocalDateTimeOf(2024, March, 15, 14, 30, 45, 123456789)

	// Should support both date and time fields
	assert.True(t, dt.IsSupportedField(FieldDayOfMonth))
	assert.True(t, dt.IsSupportedField(FieldMonthOfYear))
	assert.True(t, dt.IsSupportedField(FieldYear))
	assert.True(t, dt.IsSupportedField(FieldHourOfDay))
	assert.True(t, dt.IsSupportedField(FieldMinuteOfHour))
	assert.True(t, dt.IsSupportedField(FieldSecondOfMinute))
	assert.True(t, dt.IsSupportedField(FieldNanoOfSecond))

	// Unsupported fields
	assert.False(t, dt.IsSupportedField(FieldInstantSeconds))
	assert.False(t, dt.IsSupportedField(FieldOffsetSeconds))
}

func TestLocalDate_GetField(t *testing.T) {
	date := MustLocalDateOf(2024, March, 15) // Friday

	t.Run("supported fields", func(t *testing.T) {
		// Test FieldDayOfWeek
		dayOfWeek := date.GetField(FieldDayOfWeek)
		assert.True(t, dayOfWeek.Valid())
		assert.False(t, dayOfWeek.Unsupported())
		assert.Equal(t, int64(5), dayOfWeek.Int64()) // Friday = 5

		// Test FieldDayOfMonth
		dayOfMonth := date.GetField(FieldDayOfMonth)
		assert.True(t, dayOfMonth.Valid())
		assert.Equal(t, 15, dayOfMonth.Int())

		// Test FieldDayOfYear
		dayOfYear := date.GetField(FieldDayOfYear)
		assert.True(t, dayOfYear.Valid())
		assert.Equal(t, int64(75), dayOfYear.Int64()) // 31+29+15 = 75 (2024 is leap year)

		// Test FieldMonthOfYear
		month := date.GetField(FieldMonthOfYear)
		assert.True(t, month.Valid())
		assert.Equal(t, int64(3), month.Int64())

		// Test FieldYear
		year := date.GetField(FieldYear)
		assert.True(t, year.Valid())
		assert.Equal(t, int64(2024), year.Int64())

		// Test FieldYearOfEra
		yearOfEra := date.GetField(FieldYearOfEra)
		assert.True(t, yearOfEra.Valid())
		assert.Equal(t, int64(2024), yearOfEra.Int64())

		// Test FieldEra
		era := date.GetField(FieldEra)
		assert.True(t, era.Valid())
		assert.Equal(t, int64(1), era.Int64()) // CE

		// Test FieldEpochDay
		epochDay := date.GetField(FieldEpochDay)
		assert.True(t, epochDay.Valid())
		assert.Greater(t, epochDay.Int64(), int64(0))

		// Test FieldProlepticMonth
		prolepticMonth := date.GetField(FieldProlepticMonth)
		assert.True(t, prolepticMonth.Valid())
		assert.Equal(t, int64(2024*12+3-1), prolepticMonth.Int64())
	})

	t.Run("unsupported fields", func(t *testing.T) {
		// Time fields are not supported
		hourOfDay := date.GetField(FieldHourOfDay)
		assert.False(t, hourOfDay.Valid())
		assert.True(t, hourOfDay.Unsupported())

		nanoOfSecond := date.GetField(FieldNanoOfSecond)
		assert.False(t, nanoOfSecond.Valid())
		assert.True(t, nanoOfSecond.Unsupported())
	})

	t.Run("zero date", func(t *testing.T) {
		var zeroDate LocalDate
		field := zeroDate.GetField(FieldDayOfMonth)
		assert.False(t, field.Valid())
		assert.True(t, field.Unsupported())
	})

	t.Run("BCE date", func(t *testing.T) {
		bceDate := MustLocalDateOf(-100, January, 1)
		era := bceDate.GetField(FieldEra)
		assert.True(t, era.Valid())
		assert.Equal(t, int64(0), era.Int64()) // BCE
	})
}

func TestLocalTime_GetField(t *testing.T) {
	// Test 14:30:45.123456789 (PM)
	time := MustLocalTimeOf(14, 30, 45, 123456789)

	t.Run("supported fields", func(t *testing.T) {
		// Test FieldNanoOfSecond
		nanoOfSecond := time.GetField(FieldNanoOfSecond)
		assert.True(t, nanoOfSecond.Valid())
		assert.False(t, nanoOfSecond.Unsupported())
		assert.Equal(t, int64(123456789), nanoOfSecond.Int64())

		// Test FieldMicroOfSecond
		microOfSecond := time.GetField(FieldMicroOfSecond)
		assert.True(t, microOfSecond.Valid())
		assert.Equal(t, int64(123456), microOfSecond.Int64())

		// Test FieldMilliOfSecond
		milliOfSecond := time.GetField(FieldMilliOfSecond)
		assert.True(t, milliOfSecond.Valid())
		assert.Equal(t, int64(123), milliOfSecond.Int64())

		// Test FieldSecondOfMinute
		secondOfMinute := time.GetField(FieldSecondOfMinute)
		assert.True(t, secondOfMinute.Valid())
		assert.Equal(t, int64(45), secondOfMinute.Int64())

		// Test FieldMinuteOfHour
		minuteOfHour := time.GetField(FieldMinuteOfHour)
		assert.True(t, minuteOfHour.Valid())
		assert.Equal(t, int64(30), minuteOfHour.Int64())

		// Test FieldHourOfDay
		hourOfDay := time.GetField(FieldHourOfDay)
		assert.True(t, hourOfDay.Valid())
		assert.Equal(t, int64(14), hourOfDay.Int64())

		// Test FieldClockHourOfDay
		clockHourOfDay := time.GetField(FieldClockHourOfDay)
		assert.True(t, clockHourOfDay.Valid())
		assert.Equal(t, int64(14), clockHourOfDay.Int64())

		// Test FieldHourOfAmPm
		hourOfAmPm := time.GetField(FieldHourOfAmPm)
		assert.True(t, hourOfAmPm.Valid())
		assert.Equal(t, int64(2), hourOfAmPm.Int64()) // 14 % 12 = 2

		// Test FieldClockHourOfAmPm
		clockHourOfAmPm := time.GetField(FieldClockHourOfAmPm)
		assert.True(t, clockHourOfAmPm.Valid())
		assert.Equal(t, int64(2), clockHourOfAmPm.Int64())

		// Test FieldAmPmOfDay
		amPmOfDay := time.GetField(FieldAmPmOfDay)
		assert.True(t, amPmOfDay.Valid())
		assert.Equal(t, int64(1), amPmOfDay.Int64()) // PM

		// Test FieldNanoOfDay
		nanoOfDay := time.GetField(FieldNanoOfDay)
		assert.True(t, nanoOfDay.Valid())
		expectedNanoOfDay := int64(14)*int64(3600000000000) + int64(30)*int64(60000000000) + int64(45)*int64(1000000000) + 123456789
		assert.Equal(t, expectedNanoOfDay, nanoOfDay.Int64())

		// Test FieldMicroOfDay
		microOfDay := time.GetField(FieldMicroOfDay)
		assert.True(t, microOfDay.Valid())
		assert.Equal(t, expectedNanoOfDay/1000, microOfDay.Int64())

		// Test FieldMilliOfDay
		milliOfDay := time.GetField(FieldMilliOfDay)
		assert.True(t, milliOfDay.Valid())
		assert.Equal(t, expectedNanoOfDay/1000000, milliOfDay.Int64())

		// Test FieldSecondOfDay
		secondOfDay := time.GetField(FieldSecondOfDay)
		assert.True(t, secondOfDay.Valid())
		assert.Equal(t, int64(14*3600+30*60+45), secondOfDay.Int64())

		// Test FieldMinuteOfDay
		minuteOfDay := time.GetField(FieldMinuteOfDay)
		assert.True(t, minuteOfDay.Valid())
		assert.Equal(t, int64(14*60+30), minuteOfDay.Int64())
	})

	t.Run("midnight special cases", func(t *testing.T) {
		midnight := MustLocalTimeOf(0, 0, 0, 0)

		// FieldClockHourOfDay at midnight should be 24
		clockHourOfDay := midnight.GetField(FieldClockHourOfDay)
		assert.True(t, clockHourOfDay.Valid())
		assert.Equal(t, int64(24), clockHourOfDay.Int64())

		// FieldClockHourOfAmPm at midnight should be 12
		clockHourOfAmPm := midnight.GetField(FieldClockHourOfAmPm)
		assert.True(t, clockHourOfAmPm.Valid())
		assert.Equal(t, int64(12), clockHourOfAmPm.Int64())

		// FieldAmPmOfDay at midnight should be 0 (AM)
		amPmOfDay := midnight.GetField(FieldAmPmOfDay)
		assert.True(t, amPmOfDay.Valid())
		assert.Equal(t, int64(0), amPmOfDay.Int64())
	})

	t.Run("noon special cases", func(t *testing.T) {
		noon := MustLocalTimeOf(12, 0, 0, 0)

		// FieldClockHourOfDay at noon should be 12
		clockHourOfDay := noon.GetField(FieldClockHourOfDay)
		assert.True(t, clockHourOfDay.Valid())
		assert.Equal(t, int64(12), clockHourOfDay.Int64())

		// FieldClockHourOfAmPm at noon should be 12
		clockHourOfAmPm := noon.GetField(FieldClockHourOfAmPm)
		assert.True(t, clockHourOfAmPm.Valid())
		assert.Equal(t, int64(12), clockHourOfAmPm.Int64())

		// FieldAmPmOfDay at noon should be 1 (PM)
		amPmOfDay := noon.GetField(FieldAmPmOfDay)
		assert.True(t, amPmOfDay.Valid())
		assert.Equal(t, int64(1), amPmOfDay.Int64())
	})

	t.Run("unsupported fields", func(t *testing.T) {
		// Date fields are not supported
		dayOfMonth := time.GetField(FieldDayOfMonth)
		assert.False(t, dayOfMonth.Valid())
		assert.True(t, dayOfMonth.Unsupported())

		monthOfYear := time.GetField(FieldMonthOfYear)
		assert.False(t, monthOfYear.Valid())
		assert.True(t, monthOfYear.Unsupported())
	})

	t.Run("zero time", func(t *testing.T) {
		var zeroTime LocalTime
		field := zeroTime.GetField(FieldHourOfDay)
		assert.False(t, field.Valid())
		assert.True(t, field.Unsupported())
	})
}

func TestLocalDateTime_GetField(t *testing.T) {
	dt := MustLocalDateTimeOf(2024, March, 15, 14, 30, 45, 123456789)

	t.Run("date fields", func(t *testing.T) {
		// Test FieldDayOfWeek
		dayOfWeek := dt.GetField(FieldDayOfWeek)
		assert.True(t, dayOfWeek.Valid())
		assert.Equal(t, int64(5), dayOfWeek.Int64()) // Friday

		// Test FieldDayOfMonth
		dayOfMonth := dt.GetField(FieldDayOfMonth)
		assert.True(t, dayOfMonth.Valid())
		assert.Equal(t, int64(15), dayOfMonth.Int64())

		// Test FieldDayOfYear
		dayOfYear := dt.GetField(FieldDayOfYear)
		assert.True(t, dayOfYear.Valid())
		assert.Equal(t, int64(75), dayOfYear.Int64())

		// Test FieldMonthOfYear
		month := dt.GetField(FieldMonthOfYear)
		assert.True(t, month.Valid())
		assert.Equal(t, int64(3), month.Int64())

		// Test FieldYear
		year := dt.GetField(FieldYear)
		assert.True(t, year.Valid())
		assert.Equal(t, int64(2024), year.Int64())

		// Test FieldEra
		era := dt.GetField(FieldEra)
		assert.True(t, era.Valid())
		assert.Equal(t, int64(1), era.Int64()) // CE

		// Test FieldEpochDay
		epochDay := dt.GetField(FieldEpochDay)
		assert.True(t, epochDay.Valid())
		assert.Greater(t, epochDay.Int64(), int64(0))

		// Test FieldProlepticMonth
		prolepticMonth := dt.GetField(FieldProlepticMonth)
		assert.True(t, prolepticMonth.Valid())
		assert.Equal(t, int64(2024*12+3-1), prolepticMonth.Int64())
	})

	t.Run("time fields", func(t *testing.T) {
		// Test FieldNanoOfSecond
		nanoOfSecond := dt.GetField(FieldNanoOfSecond)
		assert.True(t, nanoOfSecond.Valid())
		assert.Equal(t, int64(123456789), nanoOfSecond.Int64())

		// Test FieldMicroOfSecond
		microOfSecond := dt.GetField(FieldMicroOfSecond)
		assert.True(t, microOfSecond.Valid())
		assert.Equal(t, int64(123456), microOfSecond.Int64())

		// Test FieldMilliOfSecond
		milliOfSecond := dt.GetField(FieldMilliOfSecond)
		assert.True(t, milliOfSecond.Valid())
		assert.Equal(t, int64(123), milliOfSecond.Int64())

		// Test FieldSecondOfMinute
		secondOfMinute := dt.GetField(FieldSecondOfMinute)
		assert.True(t, secondOfMinute.Valid())
		assert.Equal(t, int64(45), secondOfMinute.Int64())

		// Test FieldMinuteOfHour
		minuteOfHour := dt.GetField(FieldMinuteOfHour)
		assert.True(t, minuteOfHour.Valid())
		assert.Equal(t, int64(30), minuteOfHour.Int64())

		// Test FieldHourOfDay
		hourOfDay := dt.GetField(FieldHourOfDay)
		assert.True(t, hourOfDay.Valid())
		assert.Equal(t, int64(14), hourOfDay.Int64())

		// Test FieldClockHourOfDay
		clockHourOfDay := dt.GetField(FieldClockHourOfDay)
		assert.True(t, clockHourOfDay.Valid())
		assert.Equal(t, int64(14), clockHourOfDay.Int64())

		// Test FieldHourOfAmPm
		hourOfAmPm := dt.GetField(FieldHourOfAmPm)
		assert.True(t, hourOfAmPm.Valid())
		assert.Equal(t, int64(2), hourOfAmPm.Int64())

		// Test FieldAmPmOfDay
		amPmOfDay := dt.GetField(FieldAmPmOfDay)
		assert.True(t, amPmOfDay.Valid())
		assert.Equal(t, int64(1), amPmOfDay.Int64()) // PM

		// Test FieldNanoOfDay
		nanoOfDay := dt.GetField(FieldNanoOfDay)
		assert.True(t, nanoOfDay.Valid())
		expectedNanoOfDay := int64(14)*int64(3600000000000) + int64(30)*int64(60000000000) + int64(45)*int64(1000000000) + 123456789
		assert.Equal(t, expectedNanoOfDay, nanoOfDay.Int64())
	})

	t.Run("unsupported fields", func(t *testing.T) {
		// FieldInstantSeconds and FieldOffsetSeconds are not supported
		instantSeconds := dt.GetField(FieldInstantSeconds)
		assert.False(t, instantSeconds.Valid())
		assert.True(t, instantSeconds.Unsupported())

		offsetSeconds := dt.GetField(FieldOffsetSeconds)
		assert.False(t, offsetSeconds.Valid())
		assert.True(t, offsetSeconds.Unsupported())
	})

	t.Run("zero datetime", func(t *testing.T) {
		var zeroDT LocalDateTime
		field := zeroDT.GetField(FieldHourOfDay)
		assert.False(t, field.Valid())
		assert.True(t, field.Unsupported())

		field = zeroDT.GetField(FieldDayOfMonth)
		assert.False(t, field.Valid())
		assert.True(t, field.Unsupported())
	})

	t.Run("delegation to date and time", func(t *testing.T) {
		// Verify that date fields are delegated to LocalDate
		dateField := dt.GetField(FieldYear)
		dateDirectField := dt.LocalDate().GetField(FieldYear)
		assert.Equal(t, dateDirectField.Int64(), dateField.Int64())

		// Verify that time fields are delegated to LocalTime
		timeField := dt.GetField(FieldHourOfDay)
		timeDirectField := dt.LocalTime().GetField(FieldHourOfDay)
		assert.Equal(t, timeDirectField.Int64(), timeField.Int64())
	})
}
