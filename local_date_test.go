package goda

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLocalDate_Epoch(t *testing.T) {
	var check = func(i int64, tt time.Time) bool {
		var st = MustNewLocalDate(Year(tt.Year()), Month(tt.Month()), tt.Day())
		if !assert.Equal(t, i, st.UnixEpochDays(), tt) {
			return false
		}
		if !assert.Equal(t, st, NewLocalDateByUnixEpochDays(i), tt) {
			return false
		}
		if !assert.Equal(t, st.DayOfWeek().GoWeekday(), tt.Weekday(), tt) {
			return false
		}
		return assert.Equal(t, tt.Unix()/(24*60*60), st.UnixEpochDays(), tt)
	}
	var begin = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := int64(0); i < 100_0000; i++ {
		if !check(i, begin) {
			break
		}
		begin = begin.AddDate(0, 0, 1)
	}
	begin = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	// negative
	for i := int64(0); i > -100_0000; i-- {
		if !check(i, begin) {
			break
		}
		begin = begin.AddDate(0, 0, -1)
	}
}
