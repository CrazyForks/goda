package goda

import (
	"encoding"
	"fmt"
	"strconv"
)

// Year represents a year in the ISO-8601 calendar system.
// It can represent any year from math.MinInt64 to math.MaxInt64.
type Year int64

// String returns the string representation of this year.
// Years 0-9999 are formatted as 4 digits with leading zeros (e.g., "0001", "2024").
// Years outside this range are formatted without padding.
func (y Year) String() string {
	return stringImpl(y)
}

// AppendText implements the encoding.TextAppender interface.
// It appends the year representation to b and returns the extended buffer.
func (y Year) AppendText(b []byte) ([]byte, error) {
	if y >= 0 && y <= 9999 {
		return append(b, '0'+byte(y/1000), '0'+byte((y/100)%10), '0'+byte((y/10)%10), '0'+byte(y%10)), nil
	} else if y < 0 && y >= -9999 {
		return append(b, '-', '0'+byte((-y)/1000), '0'+byte(((-y)/100)%10), '0'+byte(((-y)/10)%10), '0'+byte((-y)%10)), nil
	}
	b = append(b, strconv.FormatInt(y.Int64(), 10)...)
	return b, nil
}

// Int returns this year as an int.
func (y Year) Int() int {
	return int(y)
}

// Int64 returns this year as an int64.
func (y Year) Int64() int64 {
	return int64(y)
}

// IsLeapYear returns true if this year is a leap year.
// A leap year is divisible by 4, unless it's divisible by 100 (but not 400).
// For example: 2024 is a leap year, 1900 is not, but 2000 is.
func (y Year) IsLeapYear() bool {
	return (y%4 == 0 && y%100 != 0) || (y%400 == 0)
}

// Length returns the number of days in this year (365 or 366).
func (y Year) Length() int {
	if y.IsLeapYear() {
		return 366
	}
	return 365
}

var _ encoding.TextAppender = Year(0)
var _ fmt.Stringer = Year(0)
