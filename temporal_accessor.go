package goda

type TemporalAccessor interface {
	IsZero() bool
	IsSupportedField(field Field) bool
	GetField(field Field) TemporalValue
}
