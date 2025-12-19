package goda

import "database/sql/driver"

// String returns the zone ID as a string, or empty string for zero value.
func (z ZoneId) String() string {
	if z.IsZero() {
		return ""
	}
	if z.loc != nil {
		return z.loc.String()
	}
	if z.zo == ZoneOffsetUTC() {
		return "UTC"
	}
	return z.zo.String()
}

// Scan implements the sql.Scanner interface.
// It supports scanning from nil, string, and []byte.
// Nil values are converted to the zero value of ZoneId.
func (z *ZoneId) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*z = ZoneId{}
		return nil
	case string:
		return z.UnmarshalText([]byte(v))
	case []byte:
		return z.UnmarshalText(v)
	default:
		return sqlScannerDefaultBranch(v)
	}
}

// Value implements the driver.Valuer interface.
// It returns nil for zero values, otherwise returns the zone ID as a string.
func (z ZoneId) Value() (driver.Value, error) {
	if z.IsZero() {
		return nil, nil
	}
	return z.String(), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It accepts JSON strings representing a zone ID or JSON null.
func (z *ZoneId) UnmarshalJSON(bytes []byte) error {
	if len(bytes) == 4 && string(bytes) == "null" {
		*z = ZoneId{}
		return nil
	}
	return unmarshalJsonImpl(z, bytes)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It parses zone IDs. Empty input is treated as zero value.
func (z *ZoneId) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*z = ZoneId{}
		return nil
	}
	zoneId, err := ZoneIdOf(string(text))
	if err != nil {
		return err
	}
	*z = zoneId
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// It returns the zone ID as text, or empty for zero value.
func (z ZoneId) MarshalText() (text []byte, err error) {
	return marshalTextImpl(z)
}

// MarshalJSON implements the json.Marshaler interface.
// It returns the zone ID as a JSON string, or empty string for zero value.
func (z ZoneId) MarshalJSON() ([]byte, error) {
	return marshalJsonImpl(z)
}

// AppendText implements the encoding.TextAppender interface.
// It appends the zone ID to b and returns the extended buffer.
func (z ZoneId) AppendText(b []byte) ([]byte, error) {
	if z.IsZero() {
		return b, nil
	}
	return append(b, z.String()...), nil
}
