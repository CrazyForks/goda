package goda

import "time"

// DayOfWeek represents a day-of-week in the ISO-8601 calendar system,
// where Monday=1 and Sunday=7. This differs from time.Weekday where Sunday=0.
type DayOfWeek int

const (
	Monday    DayOfWeek = iota + 1 // Monday (day 1)
	Tuesday                        // Tuesday (day 2)
	Wednesday                      // Wednesday (day 3)
	Thursday                       // Thursday (day 4)
	Friday                         // Friday (day 5)
	Saturday                       // Saturday (day 6)
	Sunday                         // Sunday (day 7)
)

// String returns the English name of the day (e.g., "Monday", "Sunday").
// Returns empty string for zero value.
func (d DayOfWeek) String() string {
	if d.IsZero() {
		return ""
	}
	return d.GoWeekday().String()
}

// IsZero returns true if this is the zero value (not a valid day-of-week).
func (d DayOfWeek) IsZero() bool {
	return d == 0
}

// GoWeekday converts this day-of-week to time.Weekday.
// Note that DayOfWeek uses ISO-8601 (Monday=1, Sunday=7) while
// time.Weekday uses Sunday=0, Monday=1, etc.
func (d DayOfWeek) GoWeekday() time.Weekday {
	if d == 0 || d == Sunday {
		return time.Sunday
	}
	return time.Weekday(d)
}

// DayOfWeekFromGoWeekday converts a time.Weekday to DayOfWeek.
func DayOfWeekFromGoWeekday(w time.Weekday) DayOfWeek {
	if w == time.Sunday {
		return Sunday
	}
	return DayOfWeek(w)
}
