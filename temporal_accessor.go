package goda

// TemporalAccessor provides read-only access to date and time fields.
// This interface is implemented by types that can provide field-based access
// to their components, such as LocalDate, LocalTime, LocalDateTime, and OffsetDateTime.
//
// The accessor pattern allows querying date/time components in a uniform way
// across different temporal types. For example, both LocalDate and LocalDateTime
// support querying the year field, but only LocalDateTime supports querying time fields.
//
// This is similar to Java's TemporalAccessor interface.
type TemporalAccessor interface {
	// IsZero returns true if this temporal object represents a zero/unset value.
	IsZero() bool

	// IsSupportedField checks if the specified field is supported by this temporal type.
	// For example, LocalDate supports date fields but not time fields.
	IsSupportedField(field Field) bool

	// GetField returns the value of the specified field.
	// If the field is not supported or the temporal is zero, returns an unsupported TemporalValue.
	GetField(field Field) TemporalValue
}
