package goda

import (
	"encoding"
	"fmt"
	"strconv"
)

type Year int64

func (y Year) String() string {
	return stringImpl(y)
}

func (y Year) AppendText(b []byte) ([]byte, error) {
	if y >= 0 && y <= 9999 {
		return append(b, '0'+byte(y/1000), '0'+byte((y/100)%10), '0'+byte((y/10)%10), '0'+byte(y%10)), nil
	} else if y < 0 && y >= -9999 {
		return append(b, '-', '0'+byte((-y)/1000), '0'+byte(((-y)/100)%10), '0'+byte(((-y)/10)%10), '0'+byte((-y)%10)), nil
	}
	b = append(b, strconv.FormatInt(y.Int64(), 10)...)
	return b, nil
}

func (y Year) Int() int {
	return int(y)
}

func (y Year) Int64() int64 {
	return int64(y)
}

func (y Year) IsLeapYear() bool {
	return (y%4 == 0 && y%100 != 0) || (y%400 == 0)
}

func (y Year) Length() int {
	if y.IsLeapYear() {
		return 366
	}
	return 365
}

var _ encoding.TextAppender = Year(0)
var _ fmt.Stringer = Year(0)
