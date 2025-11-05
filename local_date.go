package goda

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"time"
)

type LocalDate struct {
	v int64
}

func (d *LocalDate) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*d = LocalDate{}
		return nil
	case string:
		return d.UnmarshalText([]byte(v))
	case []byte:
		return d.UnmarshalText(v)
	case time.Time:
		*d = NewLocalDateByGoTime(v)
		return nil
	default:
		return sqlScannerDefaultBranch(v)
	}
}

func (d *LocalDate) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return d.String(), nil
}

func (d *LocalDate) UnmarshalJSON(bytes []byte) error {
	return unmarshalJsonImpl(d, bytes)
}

func (d *LocalDate) UnmarshalText(text []byte) (e error) {
	if len(text) == 0 {
		*d = LocalDate{}
		return nil
	}
	if text[len(text)-3] != '-' || text[len(text)-6] != '-' {
		return unmarshalError(text)
	}
	var y int64
	var m, dom int
	dom, e = parseInt(text[len(text)-2:])
	if e != nil {
		return
	}
	m, e = parseInt(text[len(text)-5 : len(text)-3])
	if e != nil {
		return
	}
	y, e = parseInt64(text[:len(text)-6])
	if e != nil {
		return
	}
	dd, e := NewLocalDate(Year(y), Month(m), dom)
	if e != nil {
		return
	}
	*d = dd
	return
}

func (d LocalDate) MarshalText() (text []byte, err error) {
	return marshalTextImpl(d)
}

func (d LocalDate) MarshalJSON() ([]byte, error) {
	return marshalJsonImpl(d)
}

func (d LocalDate) String() string {
	return stringImpl(d)
}

func (d LocalDate) AppendText(b []byte) ([]byte, error) {
	if d.IsZero() {
		return b, nil
	}
	b, _ = d.Year().AppendText(b)
	b = append(b, '-', byte('0'+d.Month()/10), byte('0'+d.Month()%10), '-', byte('0'+d.DayOfMonth()/10), byte('0'+d.DayOfMonth()%10))
	return b, nil
}

func (d LocalDate) Year() Year {
	return Year(d.v >> 16)
}

func (d LocalDate) IsLeapYear() bool {
	return d.Year().IsLeapYear()
}

func (d LocalDate) Month() Month {
	return Month(d.v >> 8 & 0xff)
}

func (d LocalDate) DayOfMonth() int {
	return int(d.v & 0xff)
}

func (d LocalDate) DayOfWeek() DayOfWeek {
	if d.IsZero() {
		return 0
	}
	return DayOfWeek(floorMod(d.UnixEpochDays()+3, 7) + 1)
}

func (d LocalDate) DayOfYear() int {
	if d.IsZero() {
		return 0
	}
	return d.Month().FirstDayOfYear(d.IsLeapYear()) - 1 + d.DayOfMonth()
}

func (d LocalDate) PlusDays(days int) LocalDate {
	if d.IsZero() {
		return d
	}
	return NewLocalDateByUnixEpochDays(d.UnixEpochDays() + int64(days))
}

func (d LocalDate) MinusDays(days int) LocalDate {
	return d.PlusDays(-days)
}

func (d LocalDate) PlusMonths(months int) LocalDate {
	if d.IsZero() {
		return d
	}
	var m = int(d.Month()) + months
	var y = d.Year().Int64()
	if m > 12 && months > 0 {
		y += int64(m / 12)
		m = m%12 + 1
	} else if m < 1 && months < 0 {
		y += int64(m/12) - 1
		m = m%12 + 12
	}
	return MustNewLocalDate(Year(y), Month(m), min(d.DayOfMonth(), Month(m).Length(Year(y).IsLeapYear())))
}

func (d LocalDate) MinusMonths(months int) LocalDate {
	return d.PlusMonths(-months)
}

func (d LocalDate) PlusYears(years int) LocalDate {
	if d.IsZero() {
		return d
	}
	var year = Year(d.Year().Int64() + int64(years))
	return MustNewLocalDate(year, d.Month(), min(d.DayOfMonth(), d.Month().Length(year.IsLeapYear())))
}

func (d LocalDate) MinusYears(years int) LocalDate {
	return d.PlusYears(-years)
}

func (d LocalDate) Compare(other LocalDate) int {
	return doCompare(d, other, compareZero, comparing(LocalDate.Year), comparing(LocalDate.Month), comparing(LocalDate.DayOfMonth))
}

func (d LocalDate) IsBefore(other LocalDate) bool {
	return d.Compare(other) < 0
}

func (d LocalDate) IsAfter(other LocalDate) bool {
	return d.Compare(other) > 0
}

func (d LocalDate) GoTime() time.Time {
	if d.IsZero() {
		return time.Time{}
	}
	return time.Date(int(d.Year()), time.Month(d.Month()), d.DayOfMonth(), 0, 0, 0, 0, time.UTC)
}

func (d LocalDate) UnixEpochDays() int64 {
	if d.IsZero() {
		return 0
	}
	const DaysPerCycle = 365*400 + 97
	const Days0000To1970 = (DaysPerCycle * 5) - (30*365 + 7)

	y := d.Year().Int64()
	m := int64(d.Month())
	day := int64(d.DayOfMonth())
	total := int64(0)

	// 计算年
	total += 365 * y
	if y >= 0 {
		total += (y+3)/4 - (y+99)/100 + (y+399)/400
	} else {
		total -= y/-4 - y/-100 + y/-400
	}

	// 计算月
	total += (367*m - 362) / 12

	// 计算天数
	total += day - 1

	// 如果月份大于2，则进行闰年修正
	if m > 2 {
		total--
		if !d.Year().IsLeapYear() {
			total--
		}
	}

	return total - Days0000To1970
}

func (d LocalDate) IsZero() bool {
	return d.v == 0
}

var _ encoding.TextAppender = (*LocalDate)(nil)
var _ fmt.Stringer = (*LocalDate)(nil)
var _ encoding.TextMarshaler = (*LocalDate)(nil)
var _ encoding.TextUnmarshaler = (*LocalDate)(nil)
var _ json.Marshaler = (*LocalDate)(nil)
var _ json.Unmarshaler = (*LocalDate)(nil)
var _ driver.Valuer = (*LocalDate)(nil)
var _ sql.Scanner = (*LocalDate)(nil)

func NewLocalDate(year Year, month Month, dayOfMonth int) (d LocalDate, e error) {
	if dayOfMonth < 1 || dayOfMonth > month.Length(year.IsLeapYear()) {
		e = newError("day %d of month out of range", dayOfMonth)
		return
	}
	if month < January || month > December {
		e = newError("month %d out of range", month)
		return
	}
	d = LocalDate{
		v: int64(year)<<16 | int64(month)<<8 | int64(dayOfMonth),
	}
	return
}

func MustNewLocalDate(year Year, month Month, dayOfMonth int) LocalDate {
	nld, e := NewLocalDate(year, month, dayOfMonth)
	if e != nil {
		panic(e)
	}
	return nld
}

func NewLocalDateByGoTime(t time.Time) LocalDate {
	if t.IsZero() {
		return LocalDate{}
	}
	return MustNewLocalDate(Year(t.Year()), Month(t.Month()), t.Day())
}

func NewLocalDateByUnixEpochDays(days int64) LocalDate {
	const DaysPerCycle = 365*400 + 97
	const Days0000To1970 = (DaysPerCycle * 5) - (30*365 + 7)
	zeroDay := days + Days0000To1970

	// 调整为基于3月的年
	zeroDay -= 60 // 将日期调整到 0000-03-01，保证闰日位于四年周期的末尾

	var adjust int64
	if zeroDay < 0 {
		// 如果是负年份，调整为正年份进行计算
		adjustCycles := (zeroDay+1)/DaysPerCycle - 1
		adjust = adjustCycles * 400
		zeroDay += -adjustCycles * DaysPerCycle
	}

	// 估算年
	yearEst := (400*zeroDay + 591) / DaysPerCycle
	// 计算该年的天数
	doyEst := zeroDay - (365*yearEst + yearEst/4 - yearEst/100 + yearEst/400)
	if doyEst < 0 {
		// 修正估算
		yearEst--
		doyEst = zeroDay - (365*yearEst + yearEst/4 - yearEst/100 + yearEst/400)
	}

	yearEst += adjust // 复位负年份

	// 转换为三月为基准的日期
	marchDoy0 := int(doyEst)

	// 转换回基于一月的月份和日期
	marchMonth0 := (marchDoy0*5 + 2) / 153
	month := Month(marchMonth0 + 3)
	if month > 12 {
		month -= 12
	}
	dom := marchDoy0 - (marchMonth0*306+5)/10 + 1
	if marchDoy0 >= 306 {
		yearEst++
	}
	return MustNewLocalDate(Year(yearEst), month, dom)
}
