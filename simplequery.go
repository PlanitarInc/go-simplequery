package simplequery

import (
	"errors"
	"net/url"
	"strconv"
	"time"

	"github.com/PlanitarInc/go-simplequery/util"
)

var (
	UnspecifiedValueErr = errors.New("The parameter value was no specified")
)

type Values struct {
	url.Values
}

func NewValues() *Values {
	return &Values{
		Values: make(url.Values),
	}
}

func FromQuery(q url.Values) *Values {
	return &Values{
		Values: q,
	}
}

func (s Values) Has(key string) bool {
	_, ok := s.Values[key]
	return ok
}

func (s Values) Get(key string) StringValue {
	return s.GetIndex(key, 0)
}

func (s Values) GetIndex(key string, idx int) StringValue {
	if vs, ok := s.Values[key]; !ok || idx < 0 || len(vs) <= idx {
		return StringValue{}
	} else {
		return StringValue{&vs[idx]}
	}
}

func (s Values) GetCombine(key string, fn func([]string) *string) StringValue {
	if vs, ok := s.Values[key]; !ok {
		// if no values, then there's nothing to combine
		return StringValue{}
	} else {
		return StringValue{fn(vs)}
	}
}

func (s Values) FilterByKey(fn func(string) bool) *Values {
	res := NewValues()
	for k, vs := range s.Values {
		if fn(k) {
			res.Values[k] = vs
		}
	}
	return res
}

func (s *Values) ForEach(fn func(string, []string)) *Values {
	for k, vs := range s.Values {
		fn(k, vs)
	}
	return s
}

type StringValue struct {
	Value *string
}

func (s StringValue) String() (string, error) {
	if s.Value == nil {
		return "", UnspecifiedValueErr
	}

	return *s.Value, nil
}

func (s StringValue) StringDefault(def string) string {
	if val, err := s.String(); err != nil {
		return def
	} else {
		return val
	}
}

func (s StringValue) Bool() (bool, error) {
	if s.Value == nil {
		return false, UnspecifiedValueErr
	}

	if val, err := util.ParseBool(*s.Value); err != nil {
		// Any value other than defined true/false values would fall here.
		// Assume true
		return true, nil
	} else {
		return val, nil
	}
}

func (s StringValue) BoolDefault(def bool) bool {
	if val, err := s.Bool(); err != nil {
		return def
	} else {
		return val
	}
}

func (s StringValue) Int64() (int64, error) {
	if s.Value == nil {
		return 0, UnspecifiedValueErr
	}

	return strconv.ParseInt(*s.Value, 0, 64)
}

func (s StringValue) Int64Default(def int64) int64 {
	if val, err := s.Int64(); err != nil {
		return def
	} else {
		return val
	}
}

func (s StringValue) Uint64() (uint64, error) {
	if s.Value == nil {
		return 0, UnspecifiedValueErr
	}

	return strconv.ParseUint(*s.Value, 0, 64)
}

func (s StringValue) Uint64Default(def uint64) uint64 {
	if val, err := s.Uint64(); err != nil {
		return def
	} else {
		return val
	}
}

func (s StringValue) Float64() (float64, error) {
	if s.Value == nil {
		return 0, UnspecifiedValueErr
	}

	return strconv.ParseFloat(*s.Value, 64)
}

func (s StringValue) Float64Default(def float64) float64 {
	if val, err := s.Float64(); err != nil {
		return def
	} else {
		return val
	}
}

func (s StringValue) Time() (time.Time, error) {
	if s.Value == nil {
		return time.Time{}, UnspecifiedValueErr
	}

	return util.ParseTime(*s.Value)
}

func (s StringValue) TimeDefault(def time.Time) time.Time {
	if val, err := s.Time(); err != nil {
		return def
	} else {
		return val
	}
}
