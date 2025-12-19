package goda

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

var locLoadMap sync.Map

func loadLocation(id string) (l *time.Location, e error) {
	v, ok := locLoadMap.Load(id)
	if ok {
		return v.(*time.Location), nil
	}
	l, e = time.LoadLocation(id)
	if e != nil {
		return nil, e
	}
	locLoadMap.Store(id, l)
	return
}

var builtinZoneShortIdMap = map[string]string{
	"ACT": "Australia/Darwin",
	"AET": "Australia/Sydney",
	"AGT": "America/Argentina/Buenos_Aires",
	"ART": "Africa/Cairo",
	"AST": "America/Anchorage",
	"BET": "America/Sao_Paulo",
	"BST": "Asia/Dhaka",
	"CAT": "Africa/Harare",
	"CNT": "America/St_Johns",
	"CST": "America/Chicago",
	"CTT": "Asia/Shanghai",
	"EAT": "Africa/Addis_Ababa",
	"ECT": "Europe/Paris",
	"IET": "America/Indiana/Indianapolis",
	"IST": "Asia/Kolkata",
	"JST": "Asia/Tokyo",
	"MIT": "Pacific/Apia",
	"NET": "Asia/Yerevan",
	"NST": "Pacific/Auckland",
	"PLT": "Asia/Karachi",
	"PNT": "America/Phoenix",
	"PRT": "America/Puerto_Rico",
	"PST": "America/Los_Angeles",
	"SST": "Pacific/Guadalcanal",
	"VST": "Asia/Ho_Chi_Minh",
	"EST": "America/Panama",
	"MST": "America/Phoenix",
	"HST": "Pacific/Honolulu",
}

// ZoneId represents a time zone identifier such as "America/New_York" or "Asia/Tokyo".
// It wraps a time.Location and provides serialization support for JSON, text, and SQL.
//
// ZoneId is comparable and can be used as a map key.
// The zero value represents an unset zone and IsZero returns true for it.
//
// ZoneId implements sql.Scanner and driver.Valuer for database operations,
// encoding.TextMarshaler and encoding.TextUnmarshaler for text serialization,
// encoding.TextAppender for efficient text appending,
// and json.Marshaler and json.Unmarshaler for JSON serialization.
//
// Format: Uses IANA Time Zone Database names (e.g., "America/New_York", "Europe/London", "UTC").
// The string representation matches Go's time.Location.String() format.
type ZoneId struct {
	loc   *time.Location
	zo    ZoneOffset
	valid bool
}

// ZoneIdOf creates a ZoneId from a time zone identifier string.
func ZoneIdOf(id string) (r ZoneId, e error) {
	defer func() { r.valid = e == nil }()
	if id == "Z" || id == "UT" || id == "UTC" || id == "GMT" {
		return ZoneIdUTC(), nil
	}
	{
		var id = id
		if strings.HasPrefix(id, "UTC") || strings.HasPrefix(id, "GMT") {
			id = id[3:]
		} else if strings.HasPrefix(id, "UT") {
			id = id[2:]
		}
		if len(id) > 0 && id[0] == '+' || id[0] == '-' {
			r.zo, e = ZoneOffsetParse(id)
			if e == nil {
				return
			}
		}
	}
	r.loc, e = loadLocation(id)
	if e != nil {
		key := builtinZoneShortIdMap[id]
		if key != "" && key != id && strings.Contains(e.Error(), "unknown time zone") {
			return ZoneIdOf(key)
		}
		e = &Error{reason: errReasonInvalidZoneId}
	}
	return
}

// MustZoneIdOf creates a ZoneId from a time zone identifier string.
// Panics if the zone ID is invalid. Use ZoneIdOf for error handling.
func MustZoneIdOf(id string) ZoneId {
	return mustValue(ZoneIdOf(id))
}

// ZoneIdUTC returns a ZoneId representing UTC (Coordinated Universal Time).
func ZoneIdUTC() ZoneId {
	return ZoneId{zo: ZoneOffsetUTC(), valid: true}
}

// ZoneIdOfGoLocation creates a ZoneId from a Go time.Location.
// This allows interoperability with Go's standard time package.
func ZoneIdOfGoLocation(l *time.Location) ZoneId {
	return ZoneId{loc: l, valid: true}
}

// ZoneIdDefault returns the system's default time zone.
// This corresponds to time.Local in Go's standard library.
// If time.Local is nil, returns ZoneIdUTC.
func ZoneIdDefault() ZoneId {
	if time.Local == nil {
		return ZoneIdUTC()
	}
	return ZoneId{loc: time.Local, valid: true}
}

// IsZero returns true if this is the zero value of ZoneId (no location set).
func (z ZoneId) IsZero() bool {
	return !z.valid
}

// Compile-time interface checks
var (
	_ encoding.TextAppender    = (*ZoneId)(nil)
	_ fmt.Stringer             = (*ZoneId)(nil)
	_ encoding.TextMarshaler   = (*ZoneId)(nil)
	_ encoding.TextUnmarshaler = (*ZoneId)(nil)
	_ json.Marshaler           = (*ZoneId)(nil)
	_ json.Unmarshaler         = (*ZoneId)(nil)
	_ driver.Valuer            = (*ZoneId)(nil)
	_ sql.Scanner              = (*ZoneId)(nil)
)

// Compile-time check that ZoneId is comparable
func _assertZoneIdIsComparable[T comparable](t T) {}

var _ = _assertZoneIdIsComparable[ZoneId]
