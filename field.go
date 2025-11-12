package goda

// Field represents a date-time field, such as month-of-year or hour-of-day.
// This is similar to Java's ChronoField.
type Field int

// Field constants representing date and time components.
const (
	// NanoOfSecond represents the nano-of-second field (0-999,999,999).
	NanoOfSecond Field = iota + 1

	// NanoOfDay represents the nano-of-day field (0-86,399,999,999,999).
	NanoOfDay

	// MicroOfSecond represents the micro-of-second field (0-999,999).
	MicroOfSecond

	// MicroOfDay represents the micro-of-day field (0-86,399,999,999).
	MicroOfDay

	// MilliOfSecond represents the milli-of-second field (0-999).
	MilliOfSecond

	// MilliOfDay represents the milli-of-day field (0-86,399,999).
	MilliOfDay

	// SecondOfMinute represents the second-of-minute field (0-59).
	SecondOfMinute

	// SecondOfDay represents the second-of-day field (0-86,399).
	SecondOfDay

	// MinuteOfHour represents the minute-of-hour field (0-59).
	MinuteOfHour

	// MinuteOfDay represents the minute-of-day field (0-1,439).
	MinuteOfDay

	// HourOfAmPm represents the hour-of-am-pm field (0-11).
	HourOfAmPm

	// ClockHourOfAmPm represents the clock-hour-of-am-pm field (1-12).
	ClockHourOfAmPm

	// HourOfDay represents the hour-of-day field (0-23).
	HourOfDay

	// ClockHourOfDay represents the clock-hour-of-day field (1-24).
	ClockHourOfDay

	// AmPmOfDay represents the am-pm-of-day field (0=AM, 1=PM).
	AmPmOfDay

	// DayOfWeekField represents the day-of-week field (1=Monday, 7=Sunday).
	DayOfWeekField

	// AlignedDayOfWeekInMonth represents the aligned day-of-week within a month.
	AlignedDayOfWeekInMonth

	// AlignedDayOfWeekInYear represents the aligned day-of-week within a year.
	AlignedDayOfWeekInYear

	// DayOfMonth represents the day-of-month field (1-31).
	DayOfMonth

	// DayOfYear represents the day-of-year field (1-366).
	DayOfYear

	// EpochDay represents the epoch-day field, based on the Unix epoch of 1970-01-01.
	EpochDay

	// AlignedWeekOfMonth represents the aligned week within a month.
	AlignedWeekOfMonth

	// AlignedWeekOfYear represents the aligned week within a year.
	AlignedWeekOfYear

	// MonthOfYear represents the month-of-year field (1=January, 12=December).
	MonthOfYear

	// ProlepticMonth represents the proleptic-month, counting months sequentially from year 0.
	ProlepticMonth

	// YearOfEra represents the year within the era.
	YearOfEra

	// YearField represents the proleptic year, such as 2024.
	YearField

	// Era represents the era field.
	Era

	// InstantSeconds represents the instant epoch-seconds.
	InstantSeconds

	// OffsetSeconds represents the offset from UTC/Greenwich in seconds.
	OffsetSeconds
)

// String returns the name of the field.
func (f Field) String() string {
	switch f {
	case NanoOfSecond:
		return "NanoOfSecond"
	case NanoOfDay:
		return "NanoOfDay"
	case MicroOfSecond:
		return "MicroOfSecond"
	case MicroOfDay:
		return "MicroOfDay"
	case MilliOfSecond:
		return "MilliOfSecond"
	case MilliOfDay:
		return "MilliOfDay"
	case SecondOfMinute:
		return "SecondOfMinute"
	case SecondOfDay:
		return "SecondOfDay"
	case MinuteOfHour:
		return "MinuteOfHour"
	case MinuteOfDay:
		return "MinuteOfDay"
	case HourOfAmPm:
		return "HourOfAmPm"
	case ClockHourOfAmPm:
		return "ClockHourOfAmPm"
	case HourOfDay:
		return "HourOfDay"
	case ClockHourOfDay:
		return "ClockHourOfDay"
	case AmPmOfDay:
		return "AmPmOfDay"
	case DayOfWeekField:
		return "DayOfWeek"
	case AlignedDayOfWeekInMonth:
		return "AlignedDayOfWeekInMonth"
	case AlignedDayOfWeekInYear:
		return "AlignedDayOfWeekInYear"
	case DayOfMonth:
		return "DayOfMonth"
	case DayOfYear:
		return "DayOfYear"
	case EpochDay:
		return "EpochDay"
	case AlignedWeekOfMonth:
		return "AlignedWeekOfMonth"
	case AlignedWeekOfYear:
		return "AlignedWeekOfYear"
	case MonthOfYear:
		return "MonthOfYear"
	case ProlepticMonth:
		return "ProlepticMonth"
	case YearOfEra:
		return "YearOfEra"
	case YearField:
		return "Year"
	case Era:
		return "Era"
	case InstantSeconds:
		return "InstantSeconds"
	case OffsetSeconds:
		return "OffsetSeconds"
	default:
		return ""
	}
}

// IsDateBased checks if this field represents a component of a date.
//
// A field is date-based if it can be derived from EpochDay.
// Note that it is valid for both IsDateBased() and IsTimeBased() to return false,
// such as when representing a field like minute-of-week.
//
// Returns true if this field is a component of a date.
func (f Field) IsDateBased() bool {
	switch f {
	case DayOfWeekField,
		AlignedDayOfWeekInMonth,
		AlignedDayOfWeekInYear,
		DayOfMonth,
		DayOfYear,
		EpochDay,
		AlignedWeekOfMonth,
		AlignedWeekOfYear,
		MonthOfYear,
		ProlepticMonth,
		YearOfEra,
		YearField,
		Era:
		return true
	default:
		return false
	}
}

// IsTimeBased checks if this field represents a component of a time.
//
// A field is time-based if it can be derived from NanoOfDay.
// Note that it is valid for both IsDateBased() and IsTimeBased() to return false,
// such as when representing a field like minute-of-week.
//
// Returns true if this field is a component of a time.
func (f Field) IsTimeBased() bool {
	switch f {
	case NanoOfSecond,
		NanoOfDay,
		MicroOfSecond,
		MicroOfDay,
		MilliOfSecond,
		MilliOfDay,
		SecondOfMinute,
		SecondOfDay,
		MinuteOfHour,
		MinuteOfDay,
		HourOfAmPm,
		ClockHourOfAmPm,
		HourOfDay,
		ClockHourOfDay,
		AmPmOfDay:
		return true
	default:
		return false
	}
}
