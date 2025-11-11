package goda

import (
	"encoding"
	"encoding/json"
	"fmt"
)

// ZoneOffset represents an offset from UTC/Greenwich, measured in seconds.
// It ranges from -18:00 to +18:00.
//
// ZoneOffset is comparable and can be used as a map key.
// The zero value represents UTC (offset 0).
type ZoneOffset struct {
	// totalSeconds stores the offset in seconds (-64800 to +64800)
	// Range: -18 hours to +18 hours
	totalSeconds int32
}

const (
	// MaxOffsetSeconds is the maximum offset in seconds (+18 hours)
	MaxOffsetSeconds = 18 * 60 * 60
	// MinOffsetSeconds is the minimum offset in seconds (-18 hours)
	MinOffsetSeconds = -18 * 60 * 60
)

// ZoneOffsetUTC represents the UTC offset (zero offset).
var ZoneOffsetUTC = ZoneOffset{totalSeconds: 0}

// NewZoneOffset creates a ZoneOffset from hours, minutes, and seconds.
// The total offset must be within -18:00:00 to +18:00:00.
// Returns an error if the offset is out of range.
func NewZoneOffset(hours, minutes, seconds int) (ZoneOffset, error) {
	totalSeconds := hours*3600 + minutes*60 + seconds
	if totalSeconds < MinOffsetSeconds || totalSeconds > MaxOffsetSeconds {
		return ZoneOffset{}, newError("offset out of range: %d:%02d:%02d", hours, minutes, seconds)
	}
	return ZoneOffset{totalSeconds: int32(totalSeconds)}, nil
}

// MustNewZoneOffset creates a ZoneOffset from hours, minutes, and seconds.
// Panics if the offset is out of range.
func MustNewZoneOffset(hours, minutes, seconds int) ZoneOffset {
	zo, err := NewZoneOffset(hours, minutes, seconds)
	if err != nil {
		panic(err)
	}
	return zo
}

// NewZoneOffsetHours creates a ZoneOffset from hours only.
// Returns an error if the offset is out of range.
func NewZoneOffsetHours(hours int) (ZoneOffset, error) {
	return NewZoneOffset(hours, 0, 0)
}

// MustNewZoneOffsetHours creates a ZoneOffset from hours only.
// Panics if the offset is out of range.
func MustNewZoneOffsetHours(hours int) ZoneOffset {
	return MustNewZoneOffset(hours, 0, 0)
}

// NewZoneOffsetSeconds creates a ZoneOffset from total seconds.
// Returns an error if the offset is out of range.
func NewZoneOffsetSeconds(seconds int) (ZoneOffset, error) {
	if seconds < MinOffsetSeconds || seconds > MaxOffsetSeconds {
		return ZoneOffset{}, newError("offset seconds out of range: %d", seconds)
	}
	return ZoneOffset{totalSeconds: int32(seconds)}, nil
}

// MustNewZoneOffsetSeconds creates a ZoneOffset from total seconds.
// Panics if the offset is out of range.
func MustNewZoneOffsetSeconds(seconds int) ZoneOffset {
	zo, err := NewZoneOffsetSeconds(seconds)
	if err != nil {
		panic(err)
	}
	return zo
}

// TotalSeconds returns the total offset in seconds.
func (zo ZoneOffset) TotalSeconds() int {
	return int(zo.totalSeconds)
}

// IsZero returns true if this is the zero value (UTC).
func (zo ZoneOffset) IsZero() bool {
	return zo.totalSeconds == 0
}

// String returns the offset in ISO 8601 format (±HH:mm or ±HH:mm:ss).
// UTC is represented as "Z".
func (zo ZoneOffset) String() string {
	return stringImpl(zo)
}

// AppendText implements encoding.TextAppender.
func (zo ZoneOffset) AppendText(b []byte) ([]byte, error) {
	seconds := zo.totalSeconds
	if seconds == 0 {
		return append(b, 'Z'), nil
	}

	if seconds < 0 {
		b = append(b, '-')
		seconds = -seconds
	} else {
		b = append(b, '+')
	}

	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	b = append(b,
		byte('0'+hours/10),
		byte('0'+hours%10),
		':',
		byte('0'+minutes/10),
		byte('0'+minutes%10),
	)

	if secs != 0 {
		b = append(b,
			':',
			byte('0'+secs/10),
			byte('0'+secs%10),
		)
	}

	return b, nil
}

// MarshalText implements encoding.TextMarshaler.
func (zo ZoneOffset) MarshalText() ([]byte, error) {
	return marshalTextImpl(zo)
}

// UnmarshalText implements encoding.TextUnmarshaler.
// Accepts ISO 8601 format: Z, ±HH:mm, ±HH:mm:ss, ±HHmm, ±HH
func (zo *ZoneOffset) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*zo = ZoneOffsetUTC
		return nil
	}

	if len(text) == 1 && (text[0] == 'Z' || text[0] == 'z') {
		*zo = ZoneOffsetUTC
		return nil
	}

	if len(text) < 3 {
		return newError("invalid offset format: %q", string(text))
	}

	sign := 1
	if text[0] == '-' {
		sign = -1
	} else if text[0] != '+' {
		return newError("invalid offset format: expected '+' or '-', got %q", string(text[0:1]))
	}

	text = text[1:] // skip sign

	var hours, minutes, seconds int
	var err error

	// Parse hours
	if len(text) < 2 {
		return newError("invalid offset format: missing hours")
	}
	hours, err = parseInt(text[0:2])
	if err != nil {
		return err
	}
	text = text[2:]

	// Parse minutes
	if len(text) > 0 {
		if text[0] == ':' {
			text = text[1:]
		}
		if len(text) >= 2 {
			minutes, err = parseInt(text[0:2])
			if err != nil {
				return err
			}
			text = text[2:]
		}
	}

	// Parse seconds
	if len(text) > 0 {
		if text[0] == ':' {
			text = text[1:]
		}
		if len(text) >= 2 {
			seconds, err = parseInt(text[0:2])
			if err != nil {
				return err
			}
		}
	}

	totalSeconds := sign * (hours*3600 + minutes*60 + seconds)
	offset, err := NewZoneOffsetSeconds(totalSeconds)
	if err != nil {
		return err
	}
	*zo = offset
	return nil
}

// MarshalJSON implements json.Marshaler.
func (zo ZoneOffset) MarshalJSON() ([]byte, error) {
	return marshalJsonImpl(zo)
}

// UnmarshalJSON implements json.Unmarshaler.
func (zo *ZoneOffset) UnmarshalJSON(data []byte) error {
	if len(data) == 4 && string(data) == "null" {
		*zo = ZoneOffsetUTC
		return nil
	}
	return unmarshalJsonImpl(zo, data)
}

// ParseZoneOffset parses an offset string in ISO 8601 format.
// Accepts: Z, ±HH:mm, ±HH:mm:ss, ±HHmm, ±HH
func ParseZoneOffset(s string) (ZoneOffset, error) {
	var zo ZoneOffset
	err := zo.UnmarshalText([]byte(s))
	return zo, err
}

// MustParseZoneOffset parses an offset string in ISO 8601 format.
// Panics if the string is invalid.
func MustParseZoneOffset(s string) ZoneOffset {
	zo, err := ParseZoneOffset(s)
	if err != nil {
		panic(err)
	}
	return zo
}

// Compile-time interface checks
var (
	_ encoding.TextAppender    = (*ZoneOffset)(nil)
	_ fmt.Stringer             = (*ZoneOffset)(nil)
	_ encoding.TextMarshaler   = (*ZoneOffset)(nil)
	_ encoding.TextUnmarshaler = (*ZoneOffset)(nil)
	_ json.Marshaler           = (*ZoneOffset)(nil)
	_ json.Unmarshaler         = (*ZoneOffset)(nil)
)
