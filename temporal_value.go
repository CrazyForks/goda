package goda

type TemporalValue struct {
	v           int64
	overflow    bool
	unsupported bool
}

func (t TemporalValue) Int64() int64 {
	return t.v
}

func (t TemporalValue) Int() int {
	return int(t.v)
}

func (t TemporalValue) Overflow() bool {
	return t.overflow
}

func (t TemporalValue) Unsupported() bool {
	return t.unsupported
}

func (t TemporalValue) Valid() bool {
	return !t.unsupported && !t.overflow
}
