package simplequery

import (
	"net/url"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestValuesGetParamValues(t *testing.T) {
	RegisterTestingT(t)

	q, err := url.ParseQuery("k1=v1_1&k1=v1_2&k2=v2&k3")
	Ω(err).Should(BeNil())
	v := FromQuery(q)

	var vs *ParamValues

	vs = v.GetParamValues("k1")
	Ω(vs.Key).Should(Equal("k1"))
	Ω(vs.Values).Should(Equal([]string{"v1_1", "v1_2"}))

	vs = v.GetParamValues("k2")
	Ω(vs.Key).Should(Equal("k2"))
	Ω(vs.Values).Should(Equal([]string{"v2"}))

	vs = v.GetParamValues("k3")
	Ω(vs.Key).Should(Equal("k3"))
	Ω(vs.Values).Should(Equal([]string{""}))

	vs = v.GetParamValues("k9")
	Ω(vs.Key).Should(Equal("k9"))
	Ω(vs.Values).Should(BeEmpty())
}

func TestValuesFilterByKey(t *testing.T) {
	RegisterTestingT(t)

	q, err := url.ParseQuery("k1=v1_1&k1=v1_2&k2=v2&k3")
	Ω(err).Should(BeNil())
	v := FromQuery(q)

	var filtered *Values

	filtered = v.FilterByKey(func(k string) bool {
		return true
	})
	Ω(filtered).Should(Equal(v))

	filtered = v.FilterByKey(func(k string) bool {
		return k != "k2"
	})
	Ω(filtered.Values).Should(HaveKeyWithValue("k1", []string{"v1_1", "v1_2"}))
	Ω(filtered.Values).Should(HaveKeyWithValue("k3", []string{""}))
	Ω(filtered.Values).Should(HaveLen(2))
}

func TestValuesForEach(t *testing.T) {
	RegisterTestingT(t)

	q, err := url.ParseQuery("k1=v1_1&k1=v1_2&k2=v2&k3")
	Ω(err).Should(BeNil())
	v := FromQuery(q)

	res := map[string][]string{}
	v.ForEach(func(p *ParamValues) {
		res[p.Key] = p.Values
	})

	Ω(res).Should(HaveKeyWithValue("k1", []string{"v1_1", "v1_2"}))
	Ω(res).Should(HaveKeyWithValue("k2", []string{"v2"}))
	Ω(res).Should(HaveKeyWithValue("k3", []string{""}))
	Ω(res).Should(HaveLen(3))
}

func TestParamValuesEmpty(t *testing.T) {
	RegisterTestingT(t)

	Ω(ParamValues{}.Empty()).Should(BeTrue())
	Ω(ParamValues{Values: []string{}}.Empty()).Should(BeTrue())
	Ω(ParamValues{Values: []string{""}}.Empty()).Should(BeFalse())
}

func TestParamValuesGetIndex(t *testing.T) {
	RegisterTestingT(t)

	var vs ParamValues
	var p *ParamPair

	vs = ParamValues{Key: "k"}

	p = vs.First()
	Ω(p.Key).Should(Equal("k"))
	Ω(p.Value).Should(BeNil())
	p = vs.GetIndex(0)
	Ω(p.Key).Should(Equal("k"))
	Ω(p.Value).Should(BeNil())
	p = vs.GetIndex(2)
	Ω(p.Key).Should(Equal("k"))
	Ω(p.Value).Should(BeNil())

	vs = ParamValues{Key: "k", Values: []string{}}

	p = vs.First()
	Ω(p.Key).Should(Equal("k"))
	Ω(p.Value).Should(BeNil())
	p = vs.GetIndex(0)
	Ω(p.Key).Should(Equal("k"))
	Ω(p.Value).Should(BeNil())
	p = vs.GetIndex(2)
	Ω(p.Key).Should(Equal("k"))
	Ω(p.Value).Should(BeNil())

	vs = ParamValues{Key: "k", Values: []string{"v0", "v1"}}

	p = vs.First()
	Ω(p.Key).Should(Equal("k"))
	Ω(p.Value).Should(Equal(pStr("v0")))
	p = vs.GetIndex(0)
	Ω(p.Key).Should(Equal("k"))
	Ω(p.Value).Should(Equal(pStr("v0")))
	p = vs.GetIndex(2)
	Ω(p.Key).Should(Equal("k"))
	Ω(p.Value).Should(BeNil())
}

func TestParamValuesCombine(t *testing.T) {
	RegisterTestingT(t)

	var vs ParamValues
	var p *ParamPair

	concatAll := func(vs []string) *string {
		return pStr(strings.Join(vs, "-"))
	}
	concatLong := func(vs []string) *string {
		if len(vs) < 1 {
			return nil
		}
		return concatAll(vs[1:])
	}

	vs = ParamValues{Key: "k"}

	p = vs.Combine(concatAll)
	Ω(p.Key).Should(Equal("k"))
	Ω(p.Value).Should(Equal(pStr("")))
	p = vs.Combine(concatLong)
	Ω(p.Key).Should(Equal("k"))
	Ω(p.Value).Should(BeNil())

	vs = ParamValues{Key: "k", Values: []string{}}

	p = vs.Combine(concatAll)
	Ω(p.Key).Should(Equal("k"))
	Ω(p.Value).Should(Equal(pStr("")))
	p = vs.Combine(concatLong)
	Ω(p.Key).Should(Equal("k"))
	Ω(p.Value).Should(BeNil())

	vs = ParamValues{Key: "k", Values: []string{"v0", "v1"}}

	p = vs.Combine(concatAll)
	Ω(p.Key).Should(Equal("k"))
	Ω(p.Value).Should(Equal(pStr("v0-v1")))
	p = vs.Combine(concatLong)
	Ω(p.Key).Should(Equal("k"))
	Ω(p.Value).Should(Equal(pStr("v1")))
}

func TestParamPairString(t *testing.T) {
	RegisterTestingT(t)

	var p ParamPair
	var res string
	var err error

	p = ParamPair{Key: "k", Value: pStr("value")}

	res, err = p.String()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal("value"))

	res = p.StringDefault("default")
	Ω(res).Should(Equal("value"))
}

func TestParamPairStringEmpty(t *testing.T) {
	RegisterTestingT(t)

	var p ParamPair
	var res string
	var err error

	p = ParamPair{Key: "k"}

	_, err = p.String()
	Ω(err).Should(Equal(UnspecifiedValueErr))

	res = p.StringDefault("default")
	Ω(res).Should(Equal("default"))
}

func TestParamPairBool(t *testing.T) {
	RegisterTestingT(t)

	var p ParamPair
	var res bool
	var err error

	p = ParamPair{Key: "k", Value: pStr("1")}

	res, err = p.Bool()
	Ω(err).Should(BeNil())
	Ω(res).Should(BeTrue())

	res = p.BoolDefault(false)
	Ω(res).Should(BeTrue())

	p = ParamPair{Key: "k", Value: pStr("0")}

	res, err = p.Bool()
	Ω(err).Should(BeNil())
	Ω(res).Should(BeFalse())

	res = p.BoolDefault(true)
	Ω(res).Should(BeFalse())
}

func TestParamPairBool_Empty(t *testing.T) {
	RegisterTestingT(t)

	var p ParamPair
	var res bool
	var err error

	p = ParamPair{Key: "k"}
	_, err = p.Bool()
	Ω(err).Should(Equal(UnspecifiedValueErr))

	res = p.BoolDefault(true)
	Ω(res).Should(BeTrue())
}

func TestParamPairBool_Invalid(t *testing.T) {
	RegisterTestingT(t)

	var p ParamPair
	var res bool
	var err error

	p = ParamPair{Key: "k", Value: pStr("")}
	res, err = p.Bool()
	Ω(err).Should(BeNil())
	Ω(res).Should(BeTrue())

	res = p.BoolDefault(false)
	Ω(res).Should(BeTrue())
}

func TestParamPairInt64(t *testing.T) {
	RegisterTestingT(t)

	var p ParamPair
	var res int64
	var err error

	p = ParamPair{Key: "k", Value: pStr("12")}

	res, err = p.Int64()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(int64(12)))

	res = p.Int64Default(-99)
	Ω(res).Should(Equal(int64(12)))

	p = ParamPair{Key: "k", Value: pStr("-12")}

	res, err = p.Int64()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(int64(-12)))

	res = p.Int64Default(-99)
	Ω(res).Should(Equal(int64(-12)))
}

func TestParamPairInt64_Empty(t *testing.T) {
	RegisterTestingT(t)

	var p ParamPair
	var res int64
	var err error

	p = ParamPair{Key: "k"}
	_, err = p.Int64()
	Ω(err).Should(Equal(UnspecifiedValueErr))

	res = p.Int64Default(-99)
	Ω(res).Should(Equal(int64(-99)))
}

func TestParamPairInt64_Invalid(t *testing.T) {
	RegisterTestingT(t)

	var p ParamPair
	var res int64
	var err error

	p = ParamPair{Key: "k", Value: pStr("")}
	_, err = p.Int64()
	Ω(err).ShouldNot(BeNil())
	Ω(err.Error()).Should(Equal(`strconv.ParseInt: parsing "": invalid syntax`))

	res = p.Int64Default(-99)
	Ω(res).Should(Equal(int64(-99)))
}

func TestParamPairUint64(t *testing.T) {
	RegisterTestingT(t)

	var p ParamPair
	var res uint64
	var err error

	p = ParamPair{Key: "k", Value: pStr("12")}

	res, err = p.Uint64()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(uint64(12)))

	res = p.Uint64Default(99)
	Ω(res).Should(Equal(uint64(12)))
}

func TestParamPairUint64_Empty(t *testing.T) {
	RegisterTestingT(t)

	var p ParamPair
	var res uint64
	var err error

	p = ParamPair{Key: "k"}
	_, err = p.Uint64()
	Ω(err).Should(Equal(UnspecifiedValueErr))

	res = p.Uint64Default(99)
	Ω(res).Should(Equal(uint64(99)))
}

func TestParamPairUint64_Invalid(t *testing.T) {
	RegisterTestingT(t)

	var p ParamPair
	var res uint64
	var err error

	p = ParamPair{Key: "k", Value: pStr("")}
	_, err = p.Uint64()
	Ω(err).ShouldNot(BeNil())
	Ω(err.Error()).Should(Equal(`strconv.ParseUint: parsing "": invalid syntax`))

	res = p.Uint64Default(99)
	Ω(res).Should(Equal(uint64(99)))
}

func TestParamPairFloat64(t *testing.T) {
	RegisterTestingT(t)

	var p ParamPair
	var res float64
	var err error

	p = ParamPair{Key: "k", Value: pStr("12.3")}

	res, err = p.Float64()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(float64(12.3)))

	res = p.Float64Default(-99.9)
	Ω(res).Should(Equal(float64(12.3)))

	p = ParamPair{Key: "k", Value: pStr("-12.3")}

	res, err = p.Float64()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(float64(-12.3)))

	res = p.Float64Default(-99.9)
	Ω(res).Should(Equal(float64(-12.3)))
}

func TestParamPairFloat64_Empty(t *testing.T) {
	RegisterTestingT(t)

	var p ParamPair
	var res float64
	var err error

	p = ParamPair{Key: "k"}
	_, err = p.Float64()
	Ω(err).Should(Equal(UnspecifiedValueErr))

	res = p.Float64Default(-99.9)
	Ω(res).Should(Equal(float64(-99.9)))
}

func TestParamPairFloat64_Invalid(t *testing.T) {
	RegisterTestingT(t)

	var p ParamPair
	var res float64
	var err error

	p = ParamPair{Key: "k", Value: pStr("")}
	_, err = p.Float64()
	Ω(err).ShouldNot(BeNil())
	Ω(err.Error()).Should(Equal(`strconv.ParseFloat: parsing "": invalid syntax`))

	res = p.Float64Default(-99.9)
	Ω(res).Should(Equal(float64(-99.9)))
}

func TestParamPairTime(t *testing.T) {
	RegisterTestingT(t)

	var p ParamPair
	var res, expTime time.Time
	var err error

	def := time.Date(1860, 7, 2, 12, 0, 0, 0, time.UTC)

	p = ParamPair{Key: "k", Value: pStr("123")}
	expTime = time.Unix(123, 0).UTC()

	res, err = p.Time()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(expTime))

	res = p.TimeDefault(def)
	Ω(res).Should(Equal(expTime))

	p = ParamPair{Key: "k", Value: pStr("252460800000")}
	expTime = time.Date(1978, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err = p.Time()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(expTime))

	res = p.TimeDefault(def)
	Ω(res).Should(Equal(expTime))

	p = ParamPair{Key: "k", Value: pStr("2016-02-03T15:04:05Z")}
	expTime = time.Date(2016, 2, 3, 15, 4, 5, 0, time.UTC)

	res, err = p.Time()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(expTime))

	res = p.TimeDefault(def)
	Ω(res).Should(Equal(expTime))
}

func TestParamPairTime_Empty(t *testing.T) {
	RegisterTestingT(t)

	var p ParamPair
	var res time.Time
	var err error

	def := time.Date(1860, 7, 2, 12, 0, 0, 0, time.UTC)

	p = ParamPair{Key: "k"}
	_, err = p.Time()
	Ω(err).Should(Equal(UnspecifiedValueErr))

	res = p.TimeDefault(def)
	Ω(res).Should(Equal(def))
}

func TestParamPairTime_Invalid(t *testing.T) {
	RegisterTestingT(t)

	var p ParamPair
	var res time.Time
	var err error

	def := time.Date(1860, 7, 2, 12, 0, 0, 0, time.UTC)

	p = ParamPair{Key: "k", Value: pStr("wrong time value")}
	_, err = p.Time()
	Ω(err).ShouldNot(BeNil())
	Ω(err.Error()).Should(HavePrefix("parsing time "))

	res = p.TimeDefault(def)
	Ω(res).Should(Equal(def))
}

func pStr(str string) *string {
	return &str
}
