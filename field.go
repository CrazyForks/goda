package goda

// Field represents a date-time field, such as month-of-year or hour-of-day.
// This is similar to Java's ChronoField.
type Field int

// Field constants representing date and time components.
const (
	// FieldNanoOfSecond represents the nano-of-second field (0-999,999,999).
	FieldNanoOfSecond Field = iota + 1

	// FieldNanoOfDay represents the nano-of-day field (0-86,399,999,999,999).
	FieldNanoOfDay

	// FieldMicroOfSecond represents the micro-of-second field (0-999,999).
	FieldMicroOfSecond

	// FieldMicroOfDay represents the micro-of-day field (0-86,399,999,999).
	FieldMicroOfDay

	// FieldMilliOfSecond represents the milli-of-second field (0-999).
	FieldMilliOfSecond

	// FieldMilliOfDay represents the milli-of-day field (0-86,399,999).
	FieldMilliOfDay

	// FieldSecondOfMinute represents the second-of-minute field (0-59).
	FieldSecondOfMinute

	// FieldSecondOfDay represents the second-of-day field (0-86,399).
	FieldSecondOfDay

	// FieldMinuteOfHour represents the minute-of-hour field (0-59).
	FieldMinuteOfHour

	// FieldMinuteOfDay represents the minute-of-day field (0-1,439).
	FieldMinuteOfDay

	// FieldHourOfAmPm represents the hour-of-am-pm field (0-11).
	FieldHourOfAmPm

	// FieldClockHourOfAmPm represents the clock-hour-of-am-pm field (1-12).
	FieldClockHourOfAmPm

	// FieldHourOfDay represents the hour-of-day field (0-23).
	FieldHourOfDay

	// FieldClockHourOfDay represents the clock-hour-of-day field (1-24).
	FieldClockHourOfDay

	// FieldAmPmOfDay represents the am-pm-of-day field (0=AM, 1=PM).
	FieldAmPmOfDay

	// FieldDayOfWeek represents the day-of-week field (1=Monday, 7=Sunday).
	FieldDayOfWeek

	// FieldAlignedDayOfWeekInMonth represents the aligned day-of-week within a month.
	FieldAlignedDayOfWeekInMonth

	// FieldAlignedDayOfWeekInYear represents the aligned day-of-week within a year.
	FieldAlignedDayOfWeekInYear

	// FieldDayOfMonth represents the day-of-month field (1-31).
	FieldDayOfMonth

	// FieldDayOfYear represents the day-of-year field (1-366).
	FieldDayOfYear

	// FieldEpochDay represents the epoch-day field, based on the Unix epoch of 1970-01-01.
	FieldEpochDay

	// FieldAlignedWeekOfMonth represents the aligned week within a month.
	FieldAlignedWeekOfMonth

	// FieldAlignedWeekOfYear represents the aligned week within a year.
	FieldAlignedWeekOfYear

	// FieldMonthOfYear represents the month-of-year field (1=January, 12=December).
	FieldMonthOfYear

	// FieldProlepticMonth represents the proleptic-month, counting months sequentially from year 0.
	FieldProlepticMonth

	// FieldYearOfEra represents the year within the era.
	FieldYearOfEra

	// FieldYear represents the proleptic year, such as 2024.
	FieldYear

	// FieldEra represents the era field.
	FieldEra

	// FieldInstantSeconds represents the instant epoch-seconds.
	FieldInstantSeconds

	// FieldOffsetSeconds represents the offset from UTC/Greenwich in seconds.
	FieldOffsetSeconds
)

// String returns the name of the field.
func (f Field) String() string {
	switch f {
	case FieldNanoOfSecond:
		return "NanoOfSecond"
	case FieldNanoOfDay:
		return "NanoOfDay"
	case FieldMicroOfSecond:
		return "MicroOfSecond"
	case FieldMicroOfDay:
		return "MicroOfDay"
	case FieldMilliOfSecond:
		return "MilliOfSecond"
	case FieldMilliOfDay:
		return "MilliOfDay"
	case FieldSecondOfMinute:
		return "SecondOfMinute"
	case FieldSecondOfDay:
		return "SecondOfDay"
	case FieldMinuteOfHour:
		return "MinuteOfHour"
	case FieldMinuteOfDay:
		return "MinuteOfDay"
	case FieldHourOfAmPm:
		return "HourOfAmPm"
	case FieldClockHourOfAmPm:
		return "ClockHourOfAmPm"
	case FieldHourOfDay:
		return "HourOfDay"
	case FieldClockHourOfDay:
		return "ClockHourOfDay"
	case FieldAmPmOfDay:
		return "AmPmOfDay"
	case FieldDayOfWeek:
		return "DayOfWeek"
	case FieldAlignedDayOfWeekInMonth:
		return "AlignedDayOfWeekInMonth"
	case FieldAlignedDayOfWeekInYear:
		return "AlignedDayOfWeekInYear"
	case FieldDayOfMonth:
		return "DayOfMonth"
	case FieldDayOfYear:
		return "DayOfYear"
	case FieldEpochDay:
		return "EpochDay"
	case FieldAlignedWeekOfMonth:
		return "AlignedWeekOfMonth"
	case FieldAlignedWeekOfYear:
		return "AlignedWeekOfYear"
	case FieldMonthOfYear:
		return "MonthOfYear"
	case FieldProlepticMonth:
		return "ProlepticMonth"
	case FieldYearOfEra:
		return "YearOfEra"
	case FieldYear:
		return "Year"
	case FieldEra:
		return "Era"
	case FieldInstantSeconds:
		return "InstantSeconds"
	case FieldOffsetSeconds:
		return "OffsetSeconds"
	default:
		return ""
	}
}

// IsDateBased checks if this field represents a component of a date.
//
// A field is date-based if it can be derived from FieldEpochDay.
// Note that it is valid for both IsDateBased() and IsTimeBased() to return false,
// such as when representing a field like minute-of-week.
//
// Returns true if this field is a component of a date.
func (f Field) IsDateBased() bool {
	switch f {
	case FieldDayOfWeek, FieldAlignedDayOfWeekInMonth, FieldAlignedDayOfWeekInYear, FieldDayOfMonth, FieldDayOfYear, FieldEpochDay, FieldAlignedWeekOfMonth, FieldAlignedWeekOfYear, FieldMonthOfYear, FieldProlepticMonth, FieldYearOfEra, FieldYear, FieldEra:
		return true
	default:
		return false
	}
}

// IsTimeBased checks if this field represents a component of a time.
//
// A field is time-based if it can be derived from FieldNanoOfDay.
// Note that it is valid for both IsDateBased() and IsTimeBased() to return false,
// such as when representing a field like minute-of-week.
//
// Returns true if this field is a component of a time.
func (f Field) IsTimeBased() bool {
	switch f {
	case FieldNanoOfSecond, FieldNanoOfDay, FieldMicroOfSecond, FieldMicroOfDay, FieldMilliOfSecond, FieldMilliOfDay, FieldSecondOfMinute, FieldSecondOfDay, FieldMinuteOfHour, FieldMinuteOfDay, FieldHourOfAmPm, FieldClockHourOfAmPm, FieldHourOfDay, FieldClockHourOfDay, FieldAmPmOfDay:
		return true
	default:
		return false
	}
}

func GetAllFields() []Field {
	return []Field{
		FieldNanoOfSecond,
		FieldNanoOfDay,
		FieldMicroOfSecond,
		FieldMicroOfDay,
		FieldMilliOfSecond,
		FieldMilliOfDay,
		FieldSecondOfMinute,
		FieldSecondOfDay,
		FieldMinuteOfHour,
		FieldMinuteOfDay,
		FieldHourOfAmPm,
		FieldClockHourOfAmPm,
		FieldHourOfDay,
		FieldClockHourOfDay,
		FieldAmPmOfDay,
		FieldDayOfWeek,
		FieldAlignedDayOfWeekInMonth,
		FieldAlignedDayOfWeekInYear,
		FieldDayOfMonth,
		FieldDayOfYear,
		FieldEpochDay,
		FieldAlignedWeekOfMonth,
		FieldAlignedWeekOfYear,
		FieldMonthOfYear,
		FieldProlepticMonth,
		FieldYearOfEra,
		FieldYear,
		FieldEra,
		FieldInstantSeconds,
		FieldOffsetSeconds,
	}
}

func (f Field) JavaName() string {
	var j string
	switch f {
	case FieldNanoOfSecond:
		j = "NANO_OF_SECOND"
	case FieldNanoOfDay:
		j = "NANO_OF_DAY"
	case FieldMicroOfSecond:
		j = "MICRO_OF_SECOND"
	case FieldMicroOfDay:
		j = "MICRO_OF_DAY"
	case FieldMilliOfSecond:
		j = "MILLI_OF_SECOND"
	case FieldMilliOfDay:
		j = "MILLI_OF_DAY"
	case FieldSecondOfMinute:
		j = "SECOND_OF_MINUTE"
	case FieldSecondOfDay:
		j = "SECOND_OF_DAY"
	case FieldMinuteOfHour:
		j = "MINUTE_OF_HOUR"
	case FieldMinuteOfDay:
		j = "MINUTE_OF_DAY"
	case FieldHourOfAmPm:
		j = "HOUR_OF_AMPM"
	case FieldClockHourOfAmPm:
		j = "CLOCK_HOUR_OF_AMPM"
	case FieldHourOfDay:
		j = "HOUR_OF_DAY"
	case FieldClockHourOfDay:
		j = "CLOCK_HOUR_OF_DAY"
	case FieldAmPmOfDay:
		j = "AMPM_OF_DAY"
	case FieldDayOfWeek:
		j = "DAY_OF_WEEK"
	case FieldAlignedDayOfWeekInMonth:
		j = "ALIGNED_DAY_OF_WEEK_IN_MONTH"
	case FieldAlignedDayOfWeekInYear:
		j = "ALIGNED_DAY_OF_WEEK_IN_YEAR"
	case FieldDayOfMonth:
		j = "DAY_OF_MONTH"
	case FieldDayOfYear:
		j = "DAY_OF_YEAR"
	case FieldEpochDay:
		j = "EPOCH_DAY"
	case FieldAlignedWeekOfMonth:
		j = "ALIGNED_WEEK_OF_MONTH"
	case FieldAlignedWeekOfYear:
		j = "ALIGNED_WEEK_OF_YEAR"
	case FieldMonthOfYear:
		j = "MONTH_OF_YEAR"
	case FieldProlepticMonth:
		j = "PROLEPTIC_MONTH"
	case FieldYearOfEra:
		j = "YEAR_OF_ERA"
	case FieldYear:
		j = "YEAR"
	case FieldEra:
		j = "ERA"
	case FieldInstantSeconds:
		j = "INSTANT_SECONDS"
	case FieldOffsetSeconds:
		j = "OFFSET_SECONDS"
	default:
		return ""
	}
	return "ChronoField." + j
}
