package goda

import "time"

type DayOfWeek int

const (
	Monday DayOfWeek = iota + 1
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func (d DayOfWeek) String() string {
	if d.IsZero() {
		return ""
	}
	return d.GoWeekday().String()
}

func (d DayOfWeek) IsZero() bool {
	return d == 0
}

func (d DayOfWeek) GoWeekday() time.Weekday {
	if d == 0 || d == Sunday {
		return time.Sunday
	}
	return time.Weekday(d)
}

func DayOfWeekFromGoWeekday(w time.Weekday) DayOfWeek {
	if w == time.Sunday {
		return Sunday
	}
	return DayOfWeek(w)
}
