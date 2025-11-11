package goda

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOffsetDateTime_NewOffsetDateTime(t *testing.T) {
	dt := MustNewLocalDateTimeFromComponents(2024, March, 15, 14, 30, 45, 123456789)
	offset := MustNewZoneOffsetHours(9)
	odt := NewOffsetDateTime(dt, offset)

	assert.Equal(t, dt, odt.LocalDateTime())
	assert.Equal(t, offset, odt.Offset())
	assert.Equal(t, Year(2024), odt.Year())
	assert.Equal(t, March, odt.Month())
	assert.Equal(t, 15, odt.DayOfMonth())
	assert.Equal(t, 14, odt.Hour())
	assert.Equal(t, 30, odt.Minute())
	assert.Equal(t, 45, odt.Second())
	assert.Equal(t, 123456789, odt.Nanosecond())
}

func TestOffsetDateTime_String(t *testing.T) {
	tests := []struct {
		name string
		odt  OffsetDateTime
		want string
	}{
		{
			"With positive offset",
			MustNewOffsetDateTimeFromComponents(2024, March, 15, 14, 30, 45, 123456789, MustNewZoneOffsetHours(9)),
			"2024-03-15T14:30:45.123456789+09:00",
		},
		{
			"With negative offset",
			MustNewOffsetDateTimeFromComponents(2024, March, 15, 14, 30, 45, 0, MustNewZoneOffsetHours(-5)),
			"2024-03-15T14:30:45-05:00",
		},
		{
			"UTC",
			MustNewOffsetDateTimeFromComponents(2024, March, 15, 14, 30, 45, 0, ZoneOffsetUTC),
			"2024-03-15T14:30:45Z",
		},
		{
			"With offset minutes",
			MustNewOffsetDateTimeFromComponents(2024, March, 15, 14, 30, 45, 0, MustNewZoneOffset(5, 30, 0)),
			"2024-03-15T14:30:45+05:30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.odt.String())
		})
	}
}

func TestOffsetDateTime_ParseOffsetDateTime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"With nanoseconds", "2024-03-15T14:30:45.123456789+09:00", false},
		{"Without fractional", "2024-03-15T14:30:45+09:00", false},
		{"UTC", "2024-03-15T14:30:45Z", false},
		{"Negative offset", "2024-03-15T14:30:45-05:00", false},
		{"With offset minutes", "2024-03-15T14:30:45+05:30", false},
		{"Invalid", "invalid", true},
		{"Missing offset", "2024-03-15T14:30:45", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			odt, err := ParseOffsetDateTime(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.False(t, odt.IsZero())
				// Round-trip test
				assert.Equal(t, tt.input, odt.String())
			}
		})
	}
}

func TestOffsetDateTime_MarshalJSON(t *testing.T) {
	odt := MustNewOffsetDateTimeFromComponents(2024, March, 15, 14, 30, 45, 123456789, MustNewZoneOffsetHours(9))
	data, err := json.Marshal(odt)
	require.NoError(t, err)
	assert.Equal(t, `"2024-03-15T14:30:45.123456789+09:00"`, string(data))

	// Test zero value
	var zero OffsetDateTime
	data, err = json.Marshal(zero)
	require.NoError(t, err)
	assert.Equal(t, `""`, string(data))
}

func TestOffsetDateTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid", `"2024-03-15T14:30:45.123456789+09:00"`, false},
		{"UTC", `"2024-03-15T14:30:45Z"`, false},
		{"null", `null`, false},
		{"Empty string", `""`, false},
		{"Invalid JSON", `invalid`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var odt OffsetDateTime
			err := json.Unmarshal([]byte(tt.input), &odt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOffsetDateTime_Compare(t *testing.T) {
	// Same instant, different offsets
	odt1 := MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")
	odt2 := MustParseOffsetDateTime("2024-03-15T05:30:45Z")
	assert.Equal(t, 0, odt1.Compare(odt2), "Same instant should be equal")
	assert.True(t, odt1.IsEqual(odt2))

	// Different instants
	odt3 := MustParseOffsetDateTime("2024-03-15T14:30:46+09:00")
	assert.Equal(t, -1, odt1.Compare(odt3))
	assert.True(t, odt1.IsBefore(odt3))
	assert.False(t, odt1.IsAfter(odt3))

	odt4 := MustParseOffsetDateTime("2024-03-15T14:30:44+09:00")
	assert.Equal(t, 1, odt1.Compare(odt4))
	assert.False(t, odt1.IsBefore(odt4))
	assert.True(t, odt1.IsAfter(odt4))

	// Zero values
	var zero OffsetDateTime
	assert.Equal(t, -1, zero.Compare(odt1))
	assert.Equal(t, 1, odt1.Compare(zero))
	assert.Equal(t, 0, zero.Compare(zero))
}

func TestOffsetDateTime_WithOffsetSameLocal(t *testing.T) {
	odt := MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")
	newOffset := MustNewZoneOffsetHours(-5)
	newOdt := odt.WithOffsetSameLocal(newOffset)

	assert.Equal(t, odt.LocalDateTime(), newOdt.LocalDateTime(), "Local datetime should be same")
	assert.Equal(t, newOffset, newOdt.Offset())
	assert.NotEqual(t, odt.ToInstant(), newOdt.ToInstant(), "Instants should be different")
}

func TestOffsetDateTime_WithOffsetSameInstant(t *testing.T) {
	odt := MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")
	newOffset := MustNewZoneOffsetHours(-5)
	newOdt := odt.WithOffsetSameInstant(newOffset)

	assert.Equal(t, odt.ToInstant(), newOdt.ToInstant(), "Instants should be same")
	assert.Equal(t, newOffset, newOdt.Offset())
	assert.NotEqual(t, odt.LocalDateTime(), newOdt.LocalDateTime(), "Local datetime should be different")
}

func TestOffsetDateTime_ToUTC(t *testing.T) {
	odt := MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")
	utc := odt.ToUTC()

	assert.Equal(t, ZoneOffsetUTC, utc.Offset())
	assert.Equal(t, odt.ToInstant(), utc.ToInstant(), "Instant should be preserved")
	assert.Equal(t, "2024-03-15T05:30:45Z", utc.String())
}

func TestOffsetDateTime_PlusMethods(t *testing.T) {
	odt := MustParseOffsetDateTime("2024-01-15T10:30:45+09:00")

	// PlusDays
	assert.Equal(t, "2024-01-25T10:30:45+09:00", odt.PlusDays(10).String())
	assert.Equal(t, "2024-01-05T10:30:45+09:00", odt.MinusDays(10).String())

	// PlusMonths
	assert.Equal(t, "2024-02-15T10:30:45+09:00", odt.PlusMonths(1).String())
	assert.Equal(t, "2023-12-15T10:30:45+09:00", odt.MinusMonths(1).String())

	// PlusYears
	assert.Equal(t, "2025-01-15T10:30:45+09:00", odt.PlusYears(1).String())
	assert.Equal(t, "2023-01-15T10:30:45+09:00", odt.MinusYears(1).String())

	// PlusHours
	plus := odt.PlusHours(2)
	assert.Equal(t, 12, plus.Hour())
	minus := odt.MinusHours(2)
	assert.Equal(t, 8, minus.Hour())

	// PlusMinutes
	plus = odt.PlusMinutes(30)
	assert.Equal(t, 11, plus.Hour())
	assert.Equal(t, 0, plus.Minute())

	// PlusSeconds
	plus = odt.PlusSeconds(15)
	assert.Equal(t, 0, plus.Second())
	assert.Equal(t, 31, plus.Minute())

	// PlusNanoseconds
	plus = odt.PlusNanoseconds(1000000000) // 1 second
	assert.Equal(t, 46, plus.Second())
}

func TestOffsetDateTime_NewOffsetDateTimeByGoTime(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)
	goTime := time.Date(2024, 3, 15, 14, 30, 45, 123456789, loc)
	odt := NewOffsetDateTimeByGoTime(goTime)

	assert.Equal(t, Year(2024), odt.Year())
	assert.Equal(t, March, odt.Month())
	assert.Equal(t, 15, odt.DayOfMonth())
	assert.Equal(t, 14, odt.Hour())
	assert.Equal(t, 30, odt.Minute())
	assert.Equal(t, 45, odt.Second())
	assert.Equal(t, 123456789, odt.Nanosecond())

	// Offset should match Tokyo's offset
	_, offset := goTime.Zone()
	assert.Equal(t, offset, odt.Offset().TotalSeconds())
}

func TestOffsetDateTime_ToInstant(t *testing.T) {
	odt := MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")
	instant := odt.ToInstant()

	// Convert back with UTC offset should give us UTC time
	utc := NewOffsetDateTimeByGoTimeWithOffset(instant, ZoneOffsetUTC)
	assert.Equal(t, "2024-03-15T05:30:45Z", utc.String())

	// Same instant should compare equal
	odt2 := MustParseOffsetDateTime("2024-03-15T05:30:45Z")
	assert.True(t, odt.IsEqual(odt2))
}

func TestOffsetDateTime_IsZero(t *testing.T) {
	var zero OffsetDateTime
	assert.True(t, zero.IsZero())

	odt := MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")
	assert.False(t, odt.IsZero())
}

func TestOffsetDateTime_IsSupportedField(t *testing.T) {
	odt := MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")

	// LocalDateTime fields
	assert.True(t, odt.IsSupportedField(YearField))
	assert.True(t, odt.IsSupportedField(MonthOfYear))
	assert.True(t, odt.IsSupportedField(DayOfMonth))
	assert.True(t, odt.IsSupportedField(HourOfDay))
	assert.True(t, odt.IsSupportedField(MinuteOfHour))
	assert.True(t, odt.IsSupportedField(SecondOfMinute))
	assert.True(t, odt.IsSupportedField(NanoOfSecond))

	// OffsetDateTime specific fields
	assert.True(t, odt.IsSupportedField(OffsetSeconds))
	assert.True(t, odt.IsSupportedField(InstantSeconds))
}

func TestOffsetDateTime_GetFieldInt64(t *testing.T) {
	odt := MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")

	assert.Equal(t, int64(2024), odt.GetFieldInt64(YearField))
	assert.Equal(t, int64(3), odt.GetFieldInt64(MonthOfYear))
	assert.Equal(t, int64(15), odt.GetFieldInt64(DayOfMonth))
	assert.Equal(t, int64(14), odt.GetFieldInt64(HourOfDay))
	assert.Equal(t, int64(30), odt.GetFieldInt64(MinuteOfHour))
	assert.Equal(t, int64(45), odt.GetFieldInt64(SecondOfMinute))
	assert.Equal(t, int64(9*3600), odt.GetFieldInt64(OffsetSeconds))

	instant := odt.ToInstant()
	assert.Equal(t, instant.Unix(), odt.GetFieldInt64(InstantSeconds))
}

func TestOffsetDateTime_Scan(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{"nil", nil, false},
		{"string", "2024-03-15T14:30:45+09:00", false},
		{"bytes", []byte("2024-03-15T14:30:45+09:00"), false},
		{"time.Time", time.Now(), false},
		{"invalid type", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var odt OffsetDateTime
			err := odt.Scan(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOffsetDateTime_Value(t *testing.T) {
	odt := MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")
	val, err := odt.Value()
	require.NoError(t, err)
	assert.Equal(t, "2024-03-15T14:30:45+09:00", val)

	// Zero value
	var zero OffsetDateTime
	val, err = zero.Value()
	require.NoError(t, err)
	assert.Nil(t, val)
}

func TestOffsetDateTimeNow(t *testing.T) {
	now := OffsetDateTimeNow()
	assert.False(t, now.IsZero())
	assert.NotEqual(t, Year(0), now.Year())

	utc := OffsetDateTimeNowUTC()
	assert.False(t, utc.IsZero())
	assert.Equal(t, ZoneOffsetUTC, utc.Offset())
}

func TestOffsetDateTime_DayOfWeek(t *testing.T) {
	// 2024-03-15 is Friday
	odt := MustParseOffsetDateTime("2024-03-15T14:30:45+09:00")
	assert.Equal(t, Friday, odt.DayOfWeek())
}

func TestOffsetDateTime_IsLeapYear(t *testing.T) {
	leap := MustParseOffsetDateTime("2024-02-29T12:00:00Z")
	assert.True(t, leap.IsLeapYear())

	notLeap := MustParseOffsetDateTime("2023-02-28T12:00:00Z")
	assert.False(t, notLeap.IsLeapYear())
}

func TestOffsetDateTime_RoundTrip(t *testing.T) {
	original := "2024-03-15T14:30:45.123456789+09:00"
	odt, err := ParseOffsetDateTime(original)
	require.NoError(t, err)

	// String round-trip
	assert.Equal(t, original, odt.String())

	// JSON round-trip
	data, err := json.Marshal(odt)
	require.NoError(t, err)
	var decoded OffsetDateTime
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, odt, decoded)

	// Text round-trip
	text, err := odt.MarshalText()
	require.NoError(t, err)
	var decoded2 OffsetDateTime
	err = decoded2.UnmarshalText(text)
	require.NoError(t, err)
	assert.Equal(t, odt, decoded2)
}

func TestOffsetDateTime_MonthArithmetic(t *testing.T) {
	// Test month overflow with offset
	odt := MustParseOffsetDateTime("2024-01-31T12:00:00+09:00")
	next := odt.PlusMonths(1)
	assert.Equal(t, "2024-02-29T12:00:00+09:00", next.String())

	// Test month underflow
	prev := odt.MinusMonths(1)
	assert.Equal(t, "2023-12-31T12:00:00+09:00", prev.String())
}

func TestOffsetDateTime_CrossDayBoundary(t *testing.T) {
	// Adding hours that cross day boundary
	odt := MustParseOffsetDateTime("2024-03-15T23:00:00+09:00")
	next := odt.PlusHours(2)
	assert.Equal(t, 16, next.DayOfMonth())
	assert.Equal(t, 1, next.Hour())

	// Subtracting hours that cross day boundary
	prev := odt.MinusHours(24)
	assert.Equal(t, 14, prev.DayOfMonth())
	assert.Equal(t, 23, prev.Hour())
}
