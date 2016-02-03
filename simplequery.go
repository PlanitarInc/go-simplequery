package simplequery

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PlanitarInc/go-simplequery/util"
)

var (
	UnspecifiedValueErr = errors.New("The parameter value was no specified")
)

type Q map[string]ValueSet

func NewValues() Q {
	return make(map[string]ValueSet)
}

func FromQuery(q url.Values) Q {
	res := NewValues()
	for k, vs := range q {
		res[k] = ValueSetFrom(vs)
	}
	return res
}

func (q Q) Has(key string) bool {
	_, ok := q[key]
	return ok
}

func (q Q) Get(key string) StringValue {
	if vs, ok := q[key]; !ok {
		return StringValue("")
	} else {
		return vs[0]
	}
}

func (q Q) GetAll(key string) ValueSet {
	if vs, ok := q[key]; !ok {
		return ValueSet{}
	} else {
		return vs
	}
}

func (q Q) FilterByKey(fn func(string) bool) Q {
	res := NewValues()
	for k, vs := range q {
		if fn(k) {
			res[k] = vs
		}
	}
	return res
}

func (q Q) ForEach(fn func(string, ValueSet)) {
	for k, vs := range q {
		fn(k, vs)
	}
}

type ValueSet []StringValue

func ValueSetFrom(arr []string) ValueSet {
	res := make([]StringValue, len(arr))
	for i := range arr {
		res[i] = StringValue(arr[i])
	}
	return res
}

func (s ValueSet) Empty() bool {
	return len(s) == 0
}

func (s ValueSet) First() StringValue {
	return s.Index(0)
}

func (s ValueSet) Index(idx int) StringValue {
	if idx < 0 || len(s) <= idx {
		return StringValue("")
	}

	return StringValue(s[idx])
}

func (s ValueSet) ForEach(fn func(StringValue)) {
	for i := range s {
		fn(s[i])
	}
}

func (s ValueSet) Strings() []string {
	res := make([]string, len(s))
	for i := range s {
		res[i] = s[i].String()
	}
	return res
}

func (s ValueSet) Bools() []bool {
	res := make([]bool, len(s))
	for i := range s {
		res[i] = s[i].Bool()
	}
	return res
}

func (s ValueSet) Int64s() []int64 {
	res := make([]int64, len(s))
	for i := range s {
		res[i] = s[i].Int64()
	}
	return res
}

func (s ValueSet) Uint64s() []uint64 {
	res := make([]uint64, len(s))
	for i := range s {
		res[i] = s[i].Uint64()
	}
	return res
}

func (s ValueSet) Float64s() []float64 {
	res := make([]float64, len(s))
	for i := range s {
		res[i] = s[i].Float64()
	}
	return res
}

func (s ValueSet) Times() []time.Time {
	res := make([]time.Time, len(s))
	for i := range s {
		res[i] = s[i].Time()
	}
	return res
}

type StringValue string

func (s StringValue) String() string {
	return string(s)
}

func (s StringValue) ParseBool() (bool, error) {
	return util.ParseBool(string(s))
}

func (s StringValue) Bool(def ...bool) bool {
	defVal := true
	if len(def) > 0 {
		defVal = def[0]
	}

	if val, err := s.ParseBool(); err != nil {
		return defVal
	} else {
		return val
	}
}

func (s StringValue) ParseInt64() (int64, error) {
	return strconv.ParseInt(string(s), 0, 64)
}

func (s StringValue) Int64(def ...int64) int64 {
	defVal := int64(0)
	if len(def) > 0 {
		defVal = def[0]
	}

	if val, err := s.ParseInt64(); err != nil {
		return defVal
	} else {
		return val
	}
}

func (s StringValue) ParseUint64() (uint64, error) {
	return strconv.ParseUint(string(s), 0, 64)
}

func (s StringValue) Uint64(def ...uint64) uint64 {
	defVal := uint64(0)
	if len(def) > 0 {
		defVal = def[0]
	}

	if val, err := s.ParseUint64(); err != nil {
		return defVal
	} else {
		return val
	}
}

func (s StringValue) ParseFloat64() (float64, error) {
	return strconv.ParseFloat(string(s), 64)
}

func (s StringValue) Float64(def ...float64) float64 {
	defVal := float64(0)
	if len(def) > 0 {
		defVal = def[0]
	}

	if val, err := s.ParseFloat64(); err != nil {
		return defVal
	} else {
		return val
	}
}

func (s StringValue) ParseTime() (time.Time, error) {
	return util.ParseTime(string(s))
}

func (s StringValue) Time(def ...time.Time) time.Time {
	defVal := time.Time{}
	if len(def) > 0 {
		defVal = def[0]
	}

	if val, err := s.ParseTime(); err != nil {
		return defVal
	} else {
		return val
	}
}

func (s StringValue) List(sep ...string) []StringValue {
	separator := ","
	if len(sep) > 0 {
		separator = sep[0]
	}

	parts := strings.Split(string(s), separator)
	res := make([]StringValue, len(parts))
	for i := range parts {
		res[i] = StringValue(parts[i])
	}
	return res
}
