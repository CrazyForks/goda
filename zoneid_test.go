package goda

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZoneIdOf(t *testing.T) {
	t.Run("valid zone IDs", func(t *testing.T) {
		testCases := []string{
			"UTC",
			"America/New_York",
			"Europe/London",
			"Asia/Tokyo",
			"Australia/Sydney",
		}

		for _, tc := range testCases {
			z, err := ZoneIdOf(tc)
			require.NoError(t, err, "zone: %s", tc)
			assert.False(t, z.IsZero())
			assert.Equal(t, tc, z.String())
		}
	})

	t.Run("invalid zone ID", func(t *testing.T) {
		_, err := ZoneIdOf("Invalid/Zone")
		assert.Error(t, err)
	})

	t.Run("empty zone ID", func(t *testing.T) {
		z, err := ZoneIdOf("")
		assert.Error(t, err)
		assert.True(t, z.IsZero())
	})
}

func TestMustZoneIdOf(t *testing.T) {
	t.Run("valid zone ID", func(t *testing.T) {
		assert.NotPanics(t, func() {
			z := MustZoneIdOf("America/New_York")
			assert.Equal(t, "America/New_York", z.String())
		})
	})

	t.Run("invalid zone ID panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustZoneIdOf("Invalid/Zone")
		})
	})
}

func TestZoneIdUTC(t *testing.T) {
	z := ZoneIdUTC()
	assert.False(t, z.IsZero())
	assert.Equal(t, "UTC", z.String())
}

func TestZoneIdOfGoLocation(t *testing.T) {
	t.Run("from time.UTC", func(t *testing.T) {
		z := ZoneIdOfGoLocation(time.UTC)
		assert.False(t, z.IsZero())
		assert.Equal(t, "UTC", z.String())
	})

	t.Run("from time.Local", func(t *testing.T) {
		z := ZoneIdOfGoLocation(time.Local)
		assert.False(t, z.IsZero())
		assert.NotEmpty(t, z.String())
	})

	t.Run("from loaded location", func(t *testing.T) {
		loc, err := time.LoadLocation("Asia/Shanghai")
		require.NoError(t, err)
		z := ZoneIdOfGoLocation(loc)
		assert.Equal(t, "Asia/Shanghai", z.String())
	})
}

func TestZoneIdDefault(t *testing.T) {
	z := ZoneIdDefault()
	assert.False(t, z.IsZero())
	assert.NotEmpty(t, z.String())
}

func TestZoneId_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		var z ZoneId
		assert.True(t, z.IsZero())
		assert.Empty(t, z.String())
	})

	t.Run("non-zero value", func(t *testing.T) {
		z := ZoneIdUTC()
		assert.False(t, z.IsZero())
	})
}

func TestZoneId_String(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		var z ZoneId
		assert.Empty(t, z.String())
	})

	t.Run("non-zero value", func(t *testing.T) {
		z := MustZoneIdOf("America/Los_Angeles")
		assert.Equal(t, "America/Los_Angeles", z.String())
	})
}

func TestZoneId_MarshalJSON(t *testing.T) {
	t.Run("valid zone", func(t *testing.T) {
		z := MustZoneIdOf("Asia/Tokyo")
		data, err := json.Marshal(z)
		require.NoError(t, err)
		assert.Equal(t, `"Asia/Tokyo"`, string(data))
	})

	t.Run("zero value", func(t *testing.T) {
		var z ZoneId
		data, err := json.Marshal(z)
		require.NoError(t, err)
		assert.Equal(t, `""`, string(data))
	})

	t.Run("UTC", func(t *testing.T) {
		z := ZoneIdUTC()
		data, err := json.Marshal(z)
		require.NoError(t, err)
		assert.Equal(t, `"UTC"`, string(data))
	})
}

func TestZoneId_UnmarshalJSON(t *testing.T) {
	t.Run("valid zone", func(t *testing.T) {
		var z ZoneId
		err := json.Unmarshal([]byte(`"America/New_York"`), &z)
		require.NoError(t, err)
		assert.Equal(t, "America/New_York", z.String())
	})

	t.Run("null", func(t *testing.T) {
		var z = MustZoneIdOf("UTC") // Set to non-zero first
		err := json.Unmarshal([]byte(`null`), &z)
		require.NoError(t, err)
		assert.True(t, z.IsZero())
	})

	t.Run("empty string", func(t *testing.T) {
		var z = MustZoneIdOf("UTC") // Set to non-zero first
		err := json.Unmarshal([]byte(`""`), &z)
		require.NoError(t, err)
		assert.True(t, z.IsZero())
	})

	t.Run("invalid zone", func(t *testing.T) {
		var z ZoneId
		err := json.Unmarshal([]byte(`"Invalid/Zone"`), &z)
		assert.Error(t, err)
	})
}

func TestZoneId_MarshalText(t *testing.T) {
	t.Run("valid zone", func(t *testing.T) {
		z := MustZoneIdOf("Europe/Berlin")
		data, err := z.MarshalText()
		require.NoError(t, err)
		assert.Equal(t, "Europe/Berlin", string(data))
	})

	t.Run("zero value", func(t *testing.T) {
		var z ZoneId
		data, err := z.MarshalText()
		require.NoError(t, err)
		assert.Empty(t, data)
	})
}

func TestZoneId_UnmarshalText(t *testing.T) {
	t.Run("valid zone", func(t *testing.T) {
		var z ZoneId
		err := z.UnmarshalText([]byte("Asia/Seoul"))
		require.NoError(t, err)
		assert.Equal(t, "Asia/Seoul", z.String())
	})

	t.Run("empty", func(t *testing.T) {
		var z = MustZoneIdOf("UTC") // Set to non-zero first
		err := z.UnmarshalText([]byte(""))
		require.NoError(t, err)
		assert.True(t, z.IsZero())
	})

	t.Run("invalid zone", func(t *testing.T) {
		var z ZoneId
		err := z.UnmarshalText([]byte("Not/A/Zone"))
		assert.Error(t, err)
	})
}

func TestZoneId_AppendText(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		var z ZoneId
		buf := []byte("prefix:")
		result, err := z.AppendText(buf)
		require.NoError(t, err)
		assert.Equal(t, "prefix:", string(result))
	})

	t.Run("non-zero value", func(t *testing.T) {
		z := MustZoneIdOf("Pacific/Auckland")
		buf := []byte("zone=")
		result, err := z.AppendText(buf)
		require.NoError(t, err)
		assert.Equal(t, "zone=Pacific/Auckland", string(result))
	})

	t.Run("empty buffer", func(t *testing.T) {
		z := MustZoneIdOf("UTC")
		result, err := z.AppendText(nil)
		require.NoError(t, err)
		assert.Equal(t, "UTC", string(result))
	})
}

func TestZoneId_Scan(t *testing.T) {
	t.Run("from nil", func(t *testing.T) {
		var z = MustZoneIdOf("UTC") // Set to non-zero first
		err := z.Scan(nil)
		require.NoError(t, err)
		assert.True(t, z.IsZero())
	})

	t.Run("from string", func(t *testing.T) {
		var z ZoneId
		err := z.Scan("Europe/Paris")
		require.NoError(t, err)
		assert.Equal(t, "Europe/Paris", z.String())
	})

	t.Run("from []byte", func(t *testing.T) {
		var z ZoneId
		err := z.Scan([]byte("Asia/Dubai"))
		require.NoError(t, err)
		assert.Equal(t, "Asia/Dubai", z.String())
	})

	t.Run("from invalid type", func(t *testing.T) {
		var z ZoneId
		err := z.Scan(12345)
		assert.Error(t, err)
	})

	t.Run("from invalid zone", func(t *testing.T) {
		var z ZoneId
		err := z.Scan("Bad/Zone")
		assert.Error(t, err)
	})
}

func TestZoneId_Value(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		var z ZoneId
		val, err := z.Value()
		require.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("non-zero value", func(t *testing.T) {
		z := MustZoneIdOf("America/Chicago")
		val, err := z.Value()
		require.NoError(t, err)
		assert.Equal(t, "America/Chicago", val)
	})
}

func TestZoneId_RoundTrip(t *testing.T) {
	testZones := []string{
		"UTC",
		"America/New_York",
		"Europe/London",
		"Asia/Tokyo",
		"Australia/Sydney",
		"Pacific/Honolulu",
	}

	for _, zoneName := range testZones {
		t.Run(zoneName, func(t *testing.T) {
			// Create original
			original := MustZoneIdOf(zoneName)

			// JSON round-trip
			jsonData, err := json.Marshal(original)
			require.NoError(t, err)

			var fromJSON ZoneId
			err = json.Unmarshal(jsonData, &fromJSON)
			require.NoError(t, err)
			assert.Equal(t, original.String(), fromJSON.String())

			// Text round-trip
			textData, err := original.MarshalText()
			require.NoError(t, err)

			var fromText ZoneId
			err = fromText.UnmarshalText(textData)
			require.NoError(t, err)
			assert.Equal(t, original.String(), fromText.String())

			// SQL round-trip
			val, err := original.Value()
			require.NoError(t, err)

			var fromSQL ZoneId
			err = fromSQL.Scan(val)
			require.NoError(t, err)
			assert.Equal(t, original.String(), fromSQL.String())
		})
	}
}

func TestZoneId_Comparable(t *testing.T) {
	z1 := MustZoneIdOf("America/New_York")
	z2 := MustZoneIdOf("America/New_York")
	z3 := MustZoneIdOf("Europe/London")
	var zero ZoneId

	// Test equality
	assert.True(t, z1 == z2)
	assert.False(t, z1 == z3)
	assert.True(t, zero == ZoneId{})

	// Can be used as map key
	m := make(map[ZoneId]string)
	m[z1] = "New York"
	m[z3] = "London"
	assert.Equal(t, "New York", m[z2]) // z2 equals z1
	assert.Equal(t, "London", m[z3])
}

func TestZoneId_ConcurrentLoadLocation(t *testing.T) {
	// Test that concurrent zone ID creation is safe
	const numGoroutines = 100
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			z, err := ZoneIdOf("America/New_York")
			assert.NoError(t, err)
			assert.Equal(t, "America/New_York", z.String())
			done <- true
		}()
	}

	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}

func TestZoneId_OfShort(t *testing.T) {
	var a = []struct{ name, zone string }{
		{"ACT", "Australia/Darwin"},
		{"AET", "Australia/Sydney"},
		{"AGT", "America/Argentina/Buenos_Aires"},
		{"ART", "Africa/Cairo"},
		{"AST", "America/Anchorage"},
		{"BET", "America/Sao_Paulo"},
		{"BST", "Asia/Dhaka"},
		{"CAT", "Africa/Harare"},
		{"CNT", "America/St_Johns"},
		{"CST", "America/Chicago"},
		{"CTT", "Asia/Shanghai"},
		{"EAT", "Africa/Addis_Ababa"},
		{"ECT", "Europe/Paris"},
		{"IET", "America/Indiana/Indianapolis"},
		{"IST", "Asia/Kolkata"},
		{"JST", "Asia/Tokyo"},
		{"MIT", "Pacific/Apia"},
		{"NET", "Asia/Yerevan"},
		{"NST", "Pacific/Auckland"},
		{"PLT", "Asia/Karachi"},
		{"PNT", "America/Phoenix"},
		{"PRT", "America/Puerto_Rico"},
		{"PST", "America/Los_Angeles"},
		{"SST", "Pacific/Guadalcanal"},
		{"VST", "Asia/Ho_Chi_Minh"},
	}
	for _, it := range a {
		t.Run(it.name, func(t *testing.T) {
			z, e := ZoneIdOf(it.name)
			assert.NoError(t, e)
			assert.Equal(t, it.zone, z.String())
		})
	}

}

func TestZoneId_GetOffset_LosAngeles2(t *testing.T) {
	var data = [][2]string{
		{"2025-11-02T00:00", "-07:00"},
		{"2025-11-02T00:15", "-07:00"},
		{"2025-11-02T00:30", "-07:00"},
		{"2025-11-02T00:45", "-07:00"},
		{"2025-11-02T01:00", "-07:00"},
		{"2025-11-02T01:15", "-07:00"},
		{"2025-11-02T01:30", "-07:00"},
		{"2025-11-02T01:45", "-07:00"},
		{"2025-11-02T02:00", "-08:00"},
		{"2025-11-02T02:15", "-08:00"},
		{"2025-11-02T02:30", "-08:00"},
		{"2025-11-02T02:45", "-08:00"},
		{"2025-11-02T03:00", "-08:00"},
		{"2025-11-02T03:15", "-08:00"},
		{"2025-11-02T03:30", "-08:00"},
		{"2025-11-02T03:45", "-08:00"},
		{"2025-11-02T04:00", "-08:00"},
		{"2025-11-02T04:15", "-08:00"},
		{"2025-11-02T04:30", "-08:00"},
		{"2025-11-02T04:45", "-08:00"},
		{"2025-11-02T05:00", "-08:00"},
		{"2025-11-02T05:15", "-08:00"},
		{"2025-11-02T05:30", "-08:00"},
		{"2025-11-02T05:45", "-08:00"},
		{"2025-11-02T06:00", "-08:00"},
		{"2025-11-02T06:15", "-08:00"},
		{"2025-11-02T06:30", "-08:00"},
		{"2025-11-02T06:45", "-08:00"},
		{"2025-11-02T07:00", "-08:00"},
		{"2025-11-02T07:15", "-08:00"},
		{"2025-11-02T07:30", "-08:00"},
		{"2025-11-02T07:45", "-08:00"},
		{"2025-11-02T08:00", "-08:00"},
		{"2025-11-02T08:15", "-08:00"},
		{"2025-11-02T08:30", "-08:00"},
		{"2025-11-02T08:45", "-08:00"},
		{"2025-11-02T09:00", "-08:00"},
		{"2025-11-02T09:15", "-08:00"},
		{"2025-11-02T09:30", "-08:00"},
		{"2025-11-02T09:45", "-08:00"},
		{"2025-11-02T10:00", "-08:00"},
		{"2025-11-02T10:15", "-08:00"},
		{"2025-11-02T10:30", "-08:00"},
		{"2025-11-02T10:45", "-08:00"},
		{"2025-11-02T11:00", "-08:00"},
		{"2025-11-02T11:15", "-08:00"},
		{"2025-11-02T11:30", "-08:00"},
		{"2025-11-02T11:45", "-08:00"},
		{"2025-11-02T12:00", "-08:00"},
		{"2025-11-02T12:15", "-08:00"},
		{"2025-11-02T12:30", "-08:00"},
		{"2025-11-02T12:45", "-08:00"},
		{"2025-11-02T13:00", "-08:00"},
		{"2025-11-02T13:15", "-08:00"},
		{"2025-11-02T13:30", "-08:00"},
		{"2025-11-02T13:45", "-08:00"},
		{"2025-11-02T14:00", "-08:00"},
		{"2025-11-02T14:15", "-08:00"},
		{"2025-11-02T14:30", "-08:00"},
		{"2025-11-02T14:45", "-08:00"},
		{"2025-11-02T15:00", "-08:00"},
		{"2025-11-02T15:15", "-08:00"},
		{"2025-11-02T15:30", "-08:00"},
		{"2025-11-02T15:45", "-08:00"},
		{"2025-11-02T16:00", "-08:00"},
		{"2025-11-02T16:15", "-08:00"},
		{"2025-11-02T16:30", "-08:00"},
		{"2025-11-02T16:45", "-08:00"},
		{"2025-11-02T17:00", "-08:00"},
		{"2025-11-02T17:15", "-08:00"},
		{"2025-11-02T17:30", "-08:00"},
		{"2025-11-02T17:45", "-08:00"},
		{"2025-11-02T18:00", "-08:00"},
		{"2025-11-02T18:15", "-08:00"},
		{"2025-11-02T18:30", "-08:00"},
		{"2025-11-02T18:45", "-08:00"},
		{"2025-11-02T19:00", "-08:00"},
		{"2025-11-02T19:15", "-08:00"},
		{"2025-11-02T19:30", "-08:00"},
		{"2025-11-02T19:45", "-08:00"},
		{"2025-11-02T20:00", "-08:00"},
		{"2025-11-02T20:15", "-08:00"},
		{"2025-11-02T20:30", "-08:00"},
		{"2025-11-02T20:45", "-08:00"},
		{"2025-11-02T21:00", "-08:00"},
		{"2025-11-02T21:15", "-08:00"},
		{"2025-11-02T21:30", "-08:00"},
		{"2025-11-02T21:45", "-08:00"},
		{"2025-11-02T22:00", "-08:00"},
		{"2025-11-02T22:15", "-08:00"},
		{"2025-11-02T22:30", "-08:00"},
		{"2025-11-02T22:45", "-08:00"},
		{"2025-11-02T23:00", "-08:00"},
		{"2025-11-02T23:15", "-08:00"},
		{"2025-11-02T23:30", "-08:00"},
		{"2025-11-02T23:45", "-08:00"},
		{"2025-11-03T00:00", "-08:00"},
		{"2025-11-03T00:15", "-08:00"},
		{"2025-11-03T00:30", "-08:00"},
		{"2025-11-03T00:45", "-08:00"},
	}
	var zoneId = MustZoneIdOf("America/Los_Angeles")
	for _, it := range data {
		ldt := MustLocalDateTimeParse(it[0])
		off := MustZoneOffsetParse(it[1])
		assert.Equal(t, off, zoneId.GetOffset(ldt))
	}
}

func TestZoneId_GetOffset_LosAngeles1(t *testing.T) {
	// var s = LocalDateTime.of(2025, 3, 9, 0, 0, 0);
	// var rules = ZoneId.of("America/Los_Angeles").getRules();
	// for(var i = 0;i<100;i++){
	//     var off = rules.getOffset(s);
	//     System.out.println(""+s+","+off);
	//     s = s.plusMinutes(15);
	// }

	var data = [][2]string{
		{"2025-03-09T00:00", "-08:00"},
		{"2025-03-09T00:15", "-08:00"},
		{"2025-03-09T00:30", "-08:00"},
		{"2025-03-09T00:45", "-08:00"},
		{"2025-03-09T01:00", "-08:00"},
		{"2025-03-09T01:15", "-08:00"},
		{"2025-03-09T01:30", "-08:00"},
		{"2025-03-09T01:45", "-08:00"},
		{"2025-03-09T02:00", "-08:00"},
		{"2025-03-09T02:15", "-08:00"},
		{"2025-03-09T02:30", "-08:00"},
		{"2025-03-09T02:45", "-08:00"},
		{"2025-03-09T03:00", "-07:00"},
		{"2025-03-09T03:15", "-07:00"},
		{"2025-03-09T03:30", "-07:00"},
		{"2025-03-09T03:45", "-07:00"},
		{"2025-03-09T04:00", "-07:00"},
		{"2025-03-09T04:15", "-07:00"},
		{"2025-03-09T04:30", "-07:00"},
		{"2025-03-09T04:45", "-07:00"},
		{"2025-03-09T05:00", "-07:00"},
		{"2025-03-09T05:15", "-07:00"},
		{"2025-03-09T05:30", "-07:00"},
		{"2025-03-09T05:45", "-07:00"},
		{"2025-03-09T06:00", "-07:00"},
		{"2025-03-09T06:15", "-07:00"},
		{"2025-03-09T06:30", "-07:00"},
		{"2025-03-09T06:45", "-07:00"},
		{"2025-03-09T07:00", "-07:00"},
		{"2025-03-09T07:15", "-07:00"},
		{"2025-03-09T07:30", "-07:00"},
		{"2025-03-09T07:45", "-07:00"},
		{"2025-03-09T08:00", "-07:00"},
		{"2025-03-09T08:15", "-07:00"},
		{"2025-03-09T08:30", "-07:00"},
		{"2025-03-09T08:45", "-07:00"},
		{"2025-03-09T09:00", "-07:00"},
		{"2025-03-09T09:15", "-07:00"},
		{"2025-03-09T09:30", "-07:00"},
		{"2025-03-09T09:45", "-07:00"},
		{"2025-03-09T10:00", "-07:00"},
		{"2025-03-09T10:15", "-07:00"},
		{"2025-03-09T10:30", "-07:00"},
		{"2025-03-09T10:45", "-07:00"},
		{"2025-03-09T11:00", "-07:00"},
		{"2025-03-09T11:15", "-07:00"},
		{"2025-03-09T11:30", "-07:00"},
		{"2025-03-09T11:45", "-07:00"},
		{"2025-03-09T12:00", "-07:00"},
		{"2025-03-09T12:15", "-07:00"},
		{"2025-03-09T12:30", "-07:00"},
		{"2025-03-09T12:45", "-07:00"},
		{"2025-03-09T13:00", "-07:00"},
		{"2025-03-09T13:15", "-07:00"},
		{"2025-03-09T13:30", "-07:00"},
		{"2025-03-09T13:45", "-07:00"},
		{"2025-03-09T14:00", "-07:00"},
		{"2025-03-09T14:15", "-07:00"},
		{"2025-03-09T14:30", "-07:00"},
		{"2025-03-09T14:45", "-07:00"},
		{"2025-03-09T15:00", "-07:00"},
		{"2025-03-09T15:15", "-07:00"},
		{"2025-03-09T15:30", "-07:00"},
		{"2025-03-09T15:45", "-07:00"},
		{"2025-03-09T16:00", "-07:00"},
		{"2025-03-09T16:15", "-07:00"},
		{"2025-03-09T16:30", "-07:00"},
		{"2025-03-09T16:45", "-07:00"},
		{"2025-03-09T17:00", "-07:00"},
		{"2025-03-09T17:15", "-07:00"},
		{"2025-03-09T17:30", "-07:00"},
		{"2025-03-09T17:45", "-07:00"},
		{"2025-03-09T18:00", "-07:00"},
		{"2025-03-09T18:15", "-07:00"},
		{"2025-03-09T18:30", "-07:00"},
		{"2025-03-09T18:45", "-07:00"},
		{"2025-03-09T19:00", "-07:00"},
		{"2025-03-09T19:15", "-07:00"},
		{"2025-03-09T19:30", "-07:00"},
		{"2025-03-09T19:45", "-07:00"},
		{"2025-03-09T20:00", "-07:00"},
		{"2025-03-09T20:15", "-07:00"},
		{"2025-03-09T20:30", "-07:00"},
		{"2025-03-09T20:45", "-07:00"},
		{"2025-03-09T21:00", "-07:00"},
		{"2025-03-09T21:15", "-07:00"},
		{"2025-03-09T21:30", "-07:00"},
		{"2025-03-09T21:45", "-07:00"},
		{"2025-03-09T22:00", "-07:00"},
		{"2025-03-09T22:15", "-07:00"},
		{"2025-03-09T22:30", "-07:00"},
		{"2025-03-09T22:45", "-07:00"},
		{"2025-03-09T23:00", "-07:00"},
		{"2025-03-09T23:15", "-07:00"},
		{"2025-03-09T23:30", "-07:00"},
		{"2025-03-09T23:45", "-07:00"},
		{"2025-03-10T00:00", "-07:00"},
		{"2025-03-10T00:15", "-07:00"},
		{"2025-03-10T00:30", "-07:00"},
		{"2025-03-10T00:45", "-07:00"},
	}
	var zoneId = MustZoneIdOf("America/Los_Angeles")
	for _, it := range data {
		ldt := MustLocalDateTimeParse(it[0])
		off := MustZoneOffsetParse(it[1])
		assert.Equal(t, off, zoneId.GetOffset(ldt))
	}
}
