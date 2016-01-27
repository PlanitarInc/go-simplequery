package simplequery

import (
	"errors"
	"net/url"
	"strconv"

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

func (s Values) GetParamValues(key string) *ParamValues {
	return &ParamValues{
		Key:    key,
		Values: s.Values[key],
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

func (s *Values) ForEach(fn func(*ParamValues)) *Values {
	tmp := ParamValues{}
	for k, vs := range s.Values {
		tmp.Key = k
		tmp.Values = vs

		fn(&tmp)
	}

	return s
}

type ParamValues struct {
	Key    string
	Values []string
}

func (vs ParamValues) Empty() bool {
	return len(vs.Values) == 0
}

func (s ParamValues) First() *ParamPair {
	return s.GetIndex(0)
}

func (s ParamValues) GetIndex(idx int) *ParamPair {
	var val *string
	if idx >= 0 && idx < len(s.Values) {
		val = &s.Values[idx]
	}

	return &ParamPair{
		Key:   s.Key,
		Value: val,
	}
}

func (s ParamValues) Combine(fn func([]string) *string) *ParamPair {
	return &ParamPair{
		Key:   s.Key,
		Value: fn(s.Values),
	}
}

type ParamPair struct {
	Key   string
	Value *string
}

func (s ParamPair) String() (string, error) {
	if s.Value == nil {
		return "", UnspecifiedValueErr
	}

	return *s.Value, nil
}

func (s ParamPair) StringDefault(def string) string {
	if val, err := s.String(); err != nil {
		return def
	} else {
		return val
	}
}

func (s ParamPair) Bool() (bool, error) {
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

func (s ParamPair) BoolDefault(def bool) bool {
	if val, err := s.Bool(); err != nil {
		return def
	} else {
		return val
	}
}

func (s ParamPair) Int64() (int64, error) {
	if s.Value == nil {
		return 0, UnspecifiedValueErr
	}

	return strconv.ParseInt(*s.Value, 0, 64)
}

func (s ParamPair) Int64Default(def int64) int64 {
	if val, err := s.Int64(); err != nil {
		return def
	} else {
		return val
	}
}

func (s ParamPair) Uint64() (uint64, error) {
	if s.Value == nil {
		return 0, UnspecifiedValueErr
	}

	return strconv.ParseUint(*s.Value, 0, 64)
}

func (s ParamPair) Uint64Default(def uint64) uint64 {
	if val, err := s.Uint64(); err != nil {
		return def
	} else {
		return val
	}
}

func (s ParamPair) Float64() (float64, error) {
	if s.Value == nil {
		return 0, UnspecifiedValueErr
	}

	return strconv.ParseFloat(*s.Value, 64)
}

func (s ParamPair) Float64Default(def float64) float64 {
	if val, err := s.Float64(); err != nil {
		return def
	} else {
		return val
	}
}
