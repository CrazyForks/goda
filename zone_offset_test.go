package goda

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZoneOffset_NewZoneOffset(t *testing.T) {
	tests := []struct {
		name    string
		hours   int
		minutes int
		seconds int
		want    int
		wantErr bool
	}{
		{"UTC", 0, 0, 0, 0, false},
		{"Plus 9 hours", 9, 0, 0, 9 * 3600, false},
		{"Minus 5 hours", -5, 0, 0, -5 * 3600, false},
		{"Plus 5:30", 5, 30, 0, 5*3600 + 30*60, false},
		{"Plus 9:30:45", 9, 30, 45, 9*3600 + 30*60 + 45, false},
		{"Max offset", 18, 0, 0, 18 * 3600, false},
		{"Min offset", -18, 0, 0, -18 * 3600, false},
		{"Over max", 19, 0, 0, 0, true},
		{"Under min", -19, 0, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zo, err := NewZoneOffset(tt.hours, tt.minutes, tt.seconds)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, zo.TotalSeconds())
			}
		})
	}
}

func TestZoneOffset_String(t *testing.T) {
	tests := []struct {
		name   string
		offset ZoneOffset
		want   string
	}{
		{"UTC", ZoneOffsetUTC, "Z"},
		{"Plus 9", MustNewZoneOffsetHours(9), "+09:00"},
		{"Minus 5", MustNewZoneOffsetHours(-5), "-05:00"},
		{"Plus 5:30", MustNewZoneOffset(5, 30, 0), "+05:30"},
		{"Plus 9:30:45", MustNewZoneOffset(9, 30, 45), "+09:30:45"},
		{"Minus 8:30", MustNewZoneOffset(-8, -30, 0), "-08:30"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.offset.String())
		})
	}
}

func TestZoneOffset_ParseZoneOffset(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{"Z", "Z", 0, false},
		{"z lowercase", "z", 0, false},
		{"+09:00", "+09:00", 9 * 3600, false},
		{"-05:00", "-05:00", -5 * 3600, false},
		{"+05:30", "+05:30", 5*3600 + 30*60, false},
		{"+0530 no colon", "+0530", 5*3600 + 30*60, false},
		{"+09", "+09", 9 * 3600, false},
		{"+09:30:45", "+09:30:45", 9*3600 + 30*60 + 45, false},
		{"invalid", "invalid", 0, true},
		{"empty", "", 0, false}, // empty is treated as UTC
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zo, err := ParseZoneOffset(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, zo.TotalSeconds())
			}
		})
	}
}

func TestZoneOffset_MarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		offset ZoneOffset
		want   string
	}{
		{"UTC", ZoneOffsetUTC, `"Z"`},
		{"Plus 9", MustNewZoneOffsetHours(9), `"+09:00"`},
		{"Minus 5", MustNewZoneOffsetHours(-5), `"-05:00"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.offset)
			require.NoError(t, err)
			assert.Equal(t, tt.want, string(data))
		})
	}
}

func TestZoneOffset_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{"Z", `"Z"`, 0, false},
		{"+09:00", `"+09:00"`, 9 * 3600, false},
		{"null", `null`, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var zo ZoneOffset
			err := json.Unmarshal([]byte(tt.input), &zo)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, zo.TotalSeconds())
			}
		})
	}
}

func TestZoneOffset_EdgeCases(t *testing.T) {
	// Test maximum offset
	max, err := NewZoneOffsetSeconds(MaxOffsetSeconds)
	require.NoError(t, err)
	assert.Equal(t, "+18:00", max.String())

	// Test minimum offset
	min, err := NewZoneOffsetSeconds(MinOffsetSeconds)
	require.NoError(t, err)
	assert.Equal(t, "-18:00", min.String())

	// Test just over maximum
	_, err = NewZoneOffsetSeconds(MaxOffsetSeconds + 1)
	assert.Error(t, err)

	// Test just under minimum
	_, err = NewZoneOffsetSeconds(MinOffsetSeconds - 1)
	assert.Error(t, err)
}
