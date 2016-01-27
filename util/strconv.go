package util

import "errors"

var (
	unknownValueError = errors.New("unknown value")
)

// Similar to strconv.ParseBool, in addition handles "on" and "off" values
func ParseBool(val string) (bool, error) {
	switch val {
	case "1", "t", "T", "true", "True", "TRUE", "on", "On", "ON":
		return true, nil
	case "0", "f", "F", "false", "False", "FALSE", "off", "Off", "OFF":
		return false, nil
	}
	return false, unknownValueError
}
