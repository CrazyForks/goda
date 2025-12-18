package goda

const (
	fnMinusDays = iota + 1
	fnMinusHours
	fnMinusMinutes
	fnMinusMonths
	fnMinusNanos
	fnMinusSeconds
	fnMinusWeeks
	fnMinusYears
	fnPlusDays
	fnPlusHours
	fnPlusMinutes
	fnPlusMonths
	fnPlusNanos
	fnPlusSeconds
	fnPlusWeeks
	fnPlusYears
	fnWithDayOfMonth
	fnWithDayOfYear
	fnWithField
	fnWithHour
	fnWithMinute
	fnWithMonth
	fnWithNano
	fnWithSecond
	fnWithYear
)

var fnNames = []string{
	fnMinusDays:      "MinusDays",
	fnMinusHours:     "MinusHours",
	fnMinusMinutes:   "MinusMinutes",
	fnMinusMonths:    "MinusMonths",
	fnMinusNanos:     "MinusNanos",
	fnMinusSeconds:   "MinusSeconds",
	fnMinusWeeks:     "MinusWeeks",
	fnMinusYears:     "MinusYears",
	fnPlusDays:       "PlusDays",
	fnPlusHours:      "PlusHours",
	fnPlusMinutes:    "PlusMinutes",
	fnPlusMonths:     "PlusMonths",
	fnPlusNanos:      "PlusNanos",
	fnPlusSeconds:    "PlusSeconds",
	fnPlusWeeks:      "PlusWeeks",
	fnPlusYears:      "PlusYears",
	fnWithDayOfMonth: "WithDayOfMonth",
	fnWithDayOfYear:  "WithDayOfYear",
	fnWithField:      "WithField",
	fnWithHour:       "WithHour",
	fnWithMinute:     "WithMinute",
	fnWithMonth:      "WithMonth",
	fnWithNano:       "WithNano",
	fnWithSecond:     "WithSecond",
	fnWithYear:       "WithYear",
}

const (
	tyLocalDate = iota + 1
	tyLocalDateTime
	tyLocalTime
	tyOffsetDateTime
	tyYearMonth
)

var tyNames = []string{
	tyLocalDate:      "LocalDate",
	tyLocalDateTime:  "LocalDateTime",
	tyLocalTime:      "LocalTime",
	tyOffsetDateTime: "OffsetDateTime",
	tyYearMonth:      "YearMonth",
}
