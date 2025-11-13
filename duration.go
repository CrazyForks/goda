package goda

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Duration represents a time-based amount of time, such as '34.5 seconds'.
// This is similar to time.Duration but following the JSR-310 API design.
//
// Duration stores time as seconds + nanoseconds, allowing for a much larger
// range than a simple nanosecond count while maintaining nanosecond precision.
//
// Duration is immutable and thread-safe.
type Duration struct {
	seconds int64 // seconds component
	nanos   int32 // nanoseconds component (0-999,999,999)
}

func NewDurationOfSeconds(seconds int64, nanoAdjustment int64) Duration {
	return Duration{
		seconds: seconds + floorDiv(nanoAdjustment, int64(time.Second)),
		nanos:   int32(floorMod(nanoAdjustment, int64(time.Second))),
	}
}

// NewDurationByGoDuration creates a Duration from a time.Duration.
func NewDurationByGoDuration(d time.Duration) Duration {
	return NewDurationOfSeconds(int64(d/time.Second), int64(d%time.Second))
}

// IsZero returns true if this duration is zero length.
func (d Duration) IsZero() bool {
	return d.seconds == 0 && d.nanos == 0
}

// IsNegative returns true if this duration is negative, excluding zero.
func (d Duration) IsNegative() bool {
	return d.seconds < 0
}

// IsPositive returns true if this duration is positive, excluding zero.
func (d Duration) IsPositive() bool {
	return d.seconds > 0 || (d.seconds == 0 && d.nanos > 0)
}

// Plus returns a copy of this duration with the specified duration added.
func (d Duration) Plus(other Duration) Duration {
	return NewDurationOfSeconds(d.seconds+other.seconds, int64(d.nanos)+int64(other.nanos))
}

// Minus returns a copy of this duration with the specified duration subtracted.
func (d Duration) Minus(other Duration) Duration {
	return d.Plus(other.Negated())
}

// Negated returns a copy of this duration with the length negated.
func (d Duration) Negated() Duration {
	return NewDurationOfSeconds(-d.seconds, -int64(d.nanos))
}

// Abs returns a copy of this duration with a positive length.
func (d Duration) Abs() Duration {
	if d.IsNegative() {
		return d.Negated()
	}
	return d
}

// Compare compares this duration to another duration.
// Returns:
//   - negative value if this < other
//   - zero if this == other
//   - positive value if this > other
func (d Duration) Compare(other Duration) int {
	if d.seconds < other.seconds {
		return -1
	}
	if d.seconds > other.seconds {
		return 1
	}
	if d.nanos < other.nanos {
		return -1
	}
	if d.nanos > other.nanos {
		return 1
	}
	return 0
}

// String returns a string representation of this duration using ISO-8601 duration format,
// such as PT8H6M12.345S.
//
// The format follows the pattern: PTnHnMn.nS where:
//   - P is the duration designator (for period) placed at the start
//   - T is the time designator that precedes the time components
//   - H is the hour designator that follows the value for hours
//   - M is the minute designator that follows the value for minutes
//   - S is the second designator that follows the value for seconds
//
// Examples:
//   - PT0S for zero duration
//   - PT8H for 8 hours
//   - PT6M for 6 minutes
//   - PT12.345S for 12.345 seconds
//   - PT8H6M12.345S for 8 hours, 6 minutes, and 12.345 seconds
//   - PT-6H3M for -6 hours and 3 minutes
func (d Duration) String() string {
	return stringImpl(d)
}

// ParseDuration parses a string to produce a Duration.
// The string must represent a valid ISO-8601 duration format, such as PT8H6M12.345S.
//
// The format follows the pattern: PTnHnMn.nS where:
//   - P is the duration designator (for period) placed at the start
//   - T is the time designator that precedes the time components
//   - H is the hour designator that follows the value for hours
//   - M is the minute designator that follows the value for minutes
//   - S is the second designator that follows the value for seconds
func ParseDuration(s string) (Duration, error) {
	if s == "" {
		return Duration{}, newError("empty duration string")
	}

	// Must start with PT
	if !strings.HasPrefix(s, "PT") {
		return Duration{}, newError("duration must start with PT")
	}

	s = s[2:] // Remove PT prefix

	if s == "" {
		return Duration{}, newError("invalid duration format")
	}

	// Handle negative sign
	negative := false
	if s[0] == '-' {
		negative = true
		s = s[1:]
	}

	var totalSeconds int64
	var totalNanos int64
	i := 0

	// Parse hours
	if idx := strings.IndexByte(s[i:], 'H'); idx != -1 {
		hours, err := strconv.ParseInt(s[i:i+idx], 10, 64)
		if err != nil {
			return Duration{}, newError("invalid hours: %v", err)
		}
		totalSeconds += hours * 3600
		i += idx + 1
	}

	// Parse minutes
	if i < len(s) {
		if idx := strings.IndexByte(s[i:], 'M'); idx != -1 {
			minutes, err := strconv.ParseInt(s[i:i+idx], 10, 64)
			if err != nil {
				return Duration{}, newError("invalid minutes: %v", err)
			}
			totalSeconds += minutes * 60
			i += idx + 1
		}
	}

	// Parse seconds (may include fractional seconds)
	if i < len(s) {
		if idx := strings.IndexByte(s[i:], 'S'); idx != -1 {
			secondsStr := s[i : i+idx]
			if dotIdx := strings.IndexByte(secondsStr, '.'); dotIdx != -1 {
				// Has fractional seconds
				wholePart := secondsStr[:dotIdx]
				fracPart := secondsStr[dotIdx+1:]

				seconds, err := strconv.ParseInt(wholePart, 10, 64)
				if err != nil {
					return Duration{}, newError("invalid seconds: %v", err)
				}

				// Pad or truncate fractional part to 9 digits (nanoseconds)
				for len(fracPart) < 9 {
					fracPart += "0"
				}
				if len(fracPart) > 9 {
					fracPart = fracPart[:9]
				}

				nanos, err := strconv.ParseInt(fracPart, 10, 64)
				if err != nil {
					return Duration{}, newError("invalid fractional seconds: %v", err)
				}

				totalSeconds += seconds
				totalNanos = nanos
			} else {
				// No fractional seconds
				seconds, err := strconv.ParseInt(secondsStr, 10, 64)
				if err != nil {
					return Duration{}, newError("invalid seconds: %v", err)
				}
				totalSeconds += seconds
			}
			i += idx + 1
		}
	}

	// Check if we parsed everything
	if i < len(s) {
		return Duration{}, newError("invalid duration format: unparsed remainder")
	}

	if totalSeconds == 0 && totalNanos == 0 && (strings.Contains(s, "H") || strings.Contains(s, "M") || strings.Contains(s, "S")) {
		// Valid parse but resulted in zero (like PT0S)
	} else if totalSeconds == 0 && totalNanos == 0 && i == 0 {
		return Duration{}, newError("invalid duration format")
	}

	if negative {
		totalSeconds = -totalSeconds
		totalNanos = -totalNanos
	}

	return NewDurationOfSeconds(totalSeconds, totalNanos), nil
}

// MustParseDuration is like ParseDuration but panics on error.
func MustParseDuration(s string) Duration {
	d, err := ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return d
}

// AppendText implements encoding.TextAppender.
// It appends the duration in ISO-8601 format to b and returns the extended buffer.
func (d Duration) AppendText(b []byte) ([]byte, error) {
	if d.IsZero() {
		return append(b, "PT0S"...), nil
	}

	b = append(b, 'P', 'T')

	seconds := d.seconds
	nanos := int64(d.nanos)

	// Handle negative durations
	// Due to our wrapping behavior, negative durations are stored with negative seconds
	// and positive nanos (e.g., -0.5s is stored as seconds=-1, nanos=500000000)
	isNegative := seconds < 0
	if isNegative {
		b = append(b, '-')
		// Convert to total nanoseconds to handle the wrapping correctly
		totalNanos := seconds*int64(time.Second) + nanos
		totalNanos = -totalNanos
		seconds = totalNanos / int64(time.Second)
		nanos = totalNanos % int64(time.Second)
	}

	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60
	seconds %= 60

	if hours != 0 {
		b = strconv.AppendInt(b, hours, 10)
		b = append(b, 'H')
	}
	if minutes != 0 {
		b = strconv.AppendInt(b, minutes, 10)
		b = append(b, 'M')
	}
	if seconds != 0 || nanos != 0 {
		b = strconv.AppendInt(b, seconds, 10)
		if nanos != 0 {
			b = append(b, '.')
			// Format nanoseconds, removing trailing zeros
			nanosStr := fmt.Sprintf("%09d", nanos)
			nanosStr = strings.TrimRight(nanosStr, "0")
			b = append(b, nanosStr...)
		}
		b = append(b, 'S')
	}

	return b, nil
}

// MarshalText implements encoding.TextMarshaler.
func (d Duration) MarshalText() ([]byte, error) {
	return marshalTextImpl(d)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (d *Duration) UnmarshalText(text []byte) error {
	parsed, err := ParseDuration(string(text))
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

// MarshalJSON implements json.Marshaler.
func (d Duration) MarshalJSON() ([]byte, error) {
	return marshalJsonImpl(d)
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *Duration) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := ParseDuration(s)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

// Value implements driver.Valuer for database serialization.
// Stores the duration as a string in ISO-8601 format.
func (d Duration) Value() (driver.Value, error) {
	return d.String(), nil
}

// Scan implements sql.Scanner for database deserialization.
// Accepts string in ISO-8601 format or int64 as nanoseconds.
func (d *Duration) Scan(value interface{}) error {
	if value == nil {
		*d = Duration{}
		return nil
	}

	switch v := value.(type) {
	case string:
		parsed, err := ParseDuration(v)
		if err != nil {
			return err
		}
		*d = parsed
		return nil
	case []byte:
		parsed, err := ParseDuration(string(v))
		if err != nil {
			return err
		}
		*d = parsed
		return nil
	case int64:
		*d = NewDurationByGoDuration(time.Duration(v))
		return nil
	default:
		return fmt.Errorf("unsupported type for Duration: %T", value)
	}
}

var _ encoding.TextAppender = (*Duration)(nil)
var _ fmt.Stringer = (*Duration)(nil)
var _ encoding.TextMarshaler = (*Duration)(nil)
var _ encoding.TextUnmarshaler = (*Duration)(nil)
var _ json.Marshaler = (*Duration)(nil)
var _ json.Unmarshaler = (*Duration)(nil)
var _ driver.Valuer = (*Duration)(nil)
var _ sql.Scanner = (*Duration)(nil)
