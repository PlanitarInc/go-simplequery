package util

import (
	"errors"
	"strconv"
	"time"
)

var (
	unknownValueError = errors.New("unknown value")
)

// Similar to strconv.ParseBool, in addition handles "on" and "off" values.
// It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False, on, On,
// ON, off, Off, OFF.
// Any other value returns an error.
func ParseBool(val string) (bool, error) {
	switch val {
	case "1", "t", "T", "true", "True", "TRUE", "on", "On", "ON":
		return true, nil
	case "0", "f", "F", "false", "False", "FALSE", "off", "Off", "OFF":
		return false, nil
	}
	return false, unknownValueError
}

const (
	// Epoch timestamp in seconds for "2100/01/01 00:00:00 UTC"
	timeSecEpochTreshold = 4102444800
)

// ParseTime converts given string to Time object.
//
// It tries to parse the string as int, if it succeeds it tries to interpret the
// number as Epoch time in seconds. If the resulting time is greater than
// 2100/01/01 00:00:00 UTC, it interprets the given number as Epoch in
// milliseconds.
//
// If the given string does not represent a decimal integer number,
// ParseTime tries to parse it as a quoted string in RFC 3339 format, with
// sub-second precision added if present.
func ParseTime(str string) (time.Time, error) {
	if val, err := strconv.ParseInt(str, 10, 64); err == nil {
		sec := val
		nsec := int64(0)
		// Heuristic: if the val is too high, treat it as milliseconds rather
		// than seconds.
		if val > timeSecEpochTreshold {
			sec = val / 1000
			nsec = (val - sec*1000) * 1000 * 1000
		}
		return time.Unix(sec, nsec).UTC(), nil
	}

	t := time.Time{}
	err := t.UnmarshalJSON([]byte(`"` + str + `"`))
	return t, err
}
