package goda

import (
	"strconv"
	"time"
)

type Month time.Month

const (
	January Month = iota + 1
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

func (m Month) IsZero() bool {
	return m == 0
}

func (m Month) FirstDayOfYear(isLeap bool) int {
	var d int
	switch m {
	case 0:
		return 0
	case January:
		d = 0 + 1
	case February:
		d = 31 + 1
	case March:
		d = 59 + 1
	case April:
		d = 90 + 1
	case May:
		d = 120 + 1
	case June:
		d = 151 + 1
	case July:
		d = 181 + 1
	case August:
		d = 212 + 1
	case September:
		d = 243 + 1
	case October:
		d = 273 + 1
	case November:
		d = 304 + 1
	case December:
		d = 334 + 1
	default:
		panic("invalid month: " + strconv.Itoa(int(m)))
	}
	if isLeap && m >= March {
		d += 1
	}
	return d
}

// MaxDays returns the maximum days in the month. For February, it's 29.
func (m Month) MaxDays() int {
	switch m {
	case 0:
		return 0
	case January, March, May, July, August, October, December:
		return 31
	case April, June, September, November:
		return 30
	case February:
		return 29
	default:
		panic("invalid month: " + strconv.Itoa(int(m)))
	}
}

func (m Month) Length(isLeap bool) int {
	if m == February {
		if isLeap {
			return 29
		}
		return 28
	}
	return m.MaxDays()
}

func (m Month) String() string {
	if m.IsZero() {
		return ""
	}
	return time.Month(m).String()
}
