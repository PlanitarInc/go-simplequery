package simplequery

import (
	"net/url"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestValuesHas(t *testing.T) {
	RegisterTestingT(t)

	urlQ, err := url.ParseQuery("k1=v1_1&k1=v1_2&k2=v2&k3")
	Ω(err).Should(BeNil())
	q := FromQuery(urlQ)

	Ω(q.Has("k1")).Should(BeTrue())
	Ω(q.Has("k2")).Should(BeTrue())
	Ω(q.Has("k3")).Should(BeTrue())
	Ω(q.Has("k9")).Should(BeFalse())
}

func TestValuesGet(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue

	urlQ, err := url.ParseQuery("k1=v1_1&k1=v1_2&k2=v2&k3")
	Ω(err).Should(BeNil())
	q := FromQuery(urlQ)

	v = q.Get("k1")
	Ω(v.Value).Should(Equal(pStr("v1_1")))

	v = q.Get("k2")
	Ω(v.Value).Should(Equal(pStr("v2")))

	v = q.Get("k3")
	Ω(v.Value).Should(Equal(pStr("")))

	v = q.Get("k9")
	Ω(v.Value).Should(BeNil())
}

func TestValuesGetIndex(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue

	urlQ, err := url.ParseQuery("k1=v1_1&k1=v1_2&k2=v2&k3")
	Ω(err).Should(BeNil())
	q := FromQuery(urlQ)

	v = q.GetIndex("k1", -1)
	Ω(v.Value).Should(BeNil())
	v = q.GetIndex("k1", 0)
	Ω(v.Value).Should(Equal(pStr("v1_1")))
	v = q.GetIndex("k1", 1)
	Ω(v.Value).Should(Equal(pStr("v1_2")))
	v = q.GetIndex("k1", 2)
	Ω(v.Value).Should(BeNil())

	v = q.GetIndex("k2", -1)
	Ω(v.Value).Should(BeNil())
	v = q.GetIndex("k2", 0)
	Ω(v.Value).Should(Equal(pStr("v2")))
	v = q.GetIndex("k2", 2)
	Ω(v.Value).Should(BeNil())

	v = q.GetIndex("k3", 0)
	Ω(v.Value).Should(Equal(pStr("")))

	v = q.GetIndex("k9", 0)
	Ω(v.Value).Should(BeNil())
}

func TestValuesGetCombine(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue

	urlQ, err := url.ParseQuery("k1=v1_1&k1=v1_2&k2=v2&k3")
	Ω(err).Should(BeNil())
	q := FromQuery(urlQ)

	reject := func(vs []string) *string {
		return nil
	}
	concatAll := func(vs []string) *string {
		return pStr(strings.Join(vs, "-"))
	}
	concatLong := func(vs []string) *string {
		if len(vs) < 2 {
			return nil
		}
		return concatAll(vs[1:])
	}

	v = q.GetCombine("k1", reject)
	Ω(v.Value).Should(BeNil())
	v = q.GetCombine("k2", reject)
	Ω(v.Value).Should(BeNil())
	v = q.GetCombine("k3", reject)
	Ω(v.Value).Should(BeNil())
	v = q.GetCombine("k9", reject)
	Ω(v.Value).Should(BeNil())

	v = q.GetCombine("k1", concatAll)
	Ω(v.Value).Should(Equal(pStr("v1_1-v1_2")))
	v = q.GetCombine("k2", concatAll)
	Ω(v.Value).Should(Equal(pStr("v2")))
	v = q.GetCombine("k3", concatAll)
	Ω(v.Value).Should(Equal(pStr("")))
	v = q.GetCombine("k9", concatAll)
	Ω(v.Value).Should(BeNil())

	v = q.GetCombine("k1", concatLong)
	Ω(v.Value).Should(Equal(pStr("v1_2")))
	v = q.GetCombine("k2", concatLong)
	Ω(v.Value).Should(BeNil())
	v = q.GetCombine("k3", concatLong)
	Ω(v.Value).Should(BeNil())
	v = q.GetCombine("k9", concatLong)
	Ω(v.Value).Should(BeNil())
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
	v.ForEach(func(key string, vs []string) {
		res[key] = vs
	})

	Ω(res).Should(HaveKeyWithValue("k1", []string{"v1_1", "v1_2"}))
	Ω(res).Should(HaveKeyWithValue("k2", []string{"v2"}))
	Ω(res).Should(HaveKeyWithValue("k3", []string{""}))
	Ω(res).Should(HaveLen(3))
}

func TestStringValueString(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res string
	var err error

	v = StringValue{Value: pStr("value")}

	res, err = v.String()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal("value"))

	res = v.StringDefault("default")
	Ω(res).Should(Equal("value"))
}

func TestStringValueStringEmpty(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res string
	var err error

	v = StringValue{}

	_, err = v.String()
	Ω(err).Should(Equal(UnspecifiedValueErr))

	res = v.StringDefault("default")
	Ω(res).Should(Equal("default"))
}

func TestStringValueBool(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res bool
	var err error

	v = StringValue{Value: pStr("1")}

	res, err = v.Bool()
	Ω(err).Should(BeNil())
	Ω(res).Should(BeTrue())

	res = v.BoolDefault(false)
	Ω(res).Should(BeTrue())

	v = StringValue{Value: pStr("0")}

	res, err = v.Bool()
	Ω(err).Should(BeNil())
	Ω(res).Should(BeFalse())

	res = v.BoolDefault(true)
	Ω(res).Should(BeFalse())
}

func TestStringValueBool_Empty(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res bool
	var err error

	v = StringValue{}
	_, err = v.Bool()
	Ω(err).Should(Equal(UnspecifiedValueErr))

	res = v.BoolDefault(true)
	Ω(res).Should(BeTrue())
}

func TestStringValueBool_Invalid(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res bool
	var err error

	v = StringValue{Value: pStr("")}
	res, err = v.Bool()
	Ω(err).Should(BeNil())
	Ω(res).Should(BeTrue())

	res = v.BoolDefault(false)
	Ω(res).Should(BeTrue())
}

func TestStringValueInt64(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res int64
	var err error

	v = StringValue{Value: pStr("12")}

	res, err = v.Int64()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(int64(12)))

	res = v.Int64Default(-99)
	Ω(res).Should(Equal(int64(12)))

	v = StringValue{Value: pStr("-12")}

	res, err = v.Int64()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(int64(-12)))

	res = v.Int64Default(-99)
	Ω(res).Should(Equal(int64(-12)))
}

func TestStringValueInt64_Empty(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res int64
	var err error

	v = StringValue{}
	_, err = v.Int64()
	Ω(err).Should(Equal(UnspecifiedValueErr))

	res = v.Int64Default(-99)
	Ω(res).Should(Equal(int64(-99)))
}

func TestStringValueInt64_Invalid(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res int64
	var err error

	v = StringValue{Value: pStr("")}
	_, err = v.Int64()
	Ω(err).ShouldNot(BeNil())
	Ω(err.Error()).Should(Equal(`strconv.ParseInt: parsing "": invalid syntax`))

	res = v.Int64Default(-99)
	Ω(res).Should(Equal(int64(-99)))
}

func TestStringValueUint64(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res uint64
	var err error

	v = StringValue{Value: pStr("12")}

	res, err = v.Uint64()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(uint64(12)))

	res = v.Uint64Default(99)
	Ω(res).Should(Equal(uint64(12)))
}

func TestStringValueUint64_Empty(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res uint64
	var err error

	v = StringValue{}
	_, err = v.Uint64()
	Ω(err).Should(Equal(UnspecifiedValueErr))

	res = v.Uint64Default(99)
	Ω(res).Should(Equal(uint64(99)))
}

func TestStringValueUint64_Invalid(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res uint64
	var err error

	v = StringValue{Value: pStr("")}
	_, err = v.Uint64()
	Ω(err).ShouldNot(BeNil())
	Ω(err.Error()).Should(Equal(`strconv.ParseUint: parsing "": invalid syntax`))

	res = v.Uint64Default(99)
	Ω(res).Should(Equal(uint64(99)))
}

func TestStringValueFloat64(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res float64
	var err error

	v = StringValue{Value: pStr("12.3")}

	res, err = v.Float64()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(float64(12.3)))

	res = v.Float64Default(-99.9)
	Ω(res).Should(Equal(float64(12.3)))

	v = StringValue{Value: pStr("-12.3")}

	res, err = v.Float64()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(float64(-12.3)))

	res = v.Float64Default(-99.9)
	Ω(res).Should(Equal(float64(-12.3)))
}

func TestStringValueFloat64_Empty(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res float64
	var err error

	v = StringValue{}
	_, err = v.Float64()
	Ω(err).Should(Equal(UnspecifiedValueErr))

	res = v.Float64Default(-99.9)
	Ω(res).Should(Equal(float64(-99.9)))
}

func TestStringValueFloat64_Invalid(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res float64
	var err error

	v = StringValue{Value: pStr("")}
	_, err = v.Float64()
	Ω(err).ShouldNot(BeNil())
	Ω(err.Error()).Should(Equal(`strconv.ParseFloat: parsing "": invalid syntax`))

	res = v.Float64Default(-99.9)
	Ω(res).Should(Equal(float64(-99.9)))
}

func TestStringValueTime(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res, expTime time.Time
	var err error

	def := time.Date(1860, 7, 2, 12, 0, 0, 0, time.UTC)

	v = StringValue{Value: pStr("123")}
	expTime = time.Unix(123, 0).UTC()

	res, err = v.Time()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(expTime))

	res = v.TimeDefault(def)
	Ω(res).Should(Equal(expTime))

	v = StringValue{Value: pStr("252460800000")}
	expTime = time.Date(1978, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err = v.Time()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(expTime))

	res = v.TimeDefault(def)
	Ω(res).Should(Equal(expTime))

	v = StringValue{Value: pStr("2016-02-03T15:04:05Z")}
	expTime = time.Date(2016, 2, 3, 15, 4, 5, 0, time.UTC)

	res, err = v.Time()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(expTime))

	res = v.TimeDefault(def)
	Ω(res).Should(Equal(expTime))
}

func TestStringValueTime_Empty(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res time.Time
	var err error

	def := time.Date(1860, 7, 2, 12, 0, 0, 0, time.UTC)

	v = StringValue{}
	_, err = v.Time()
	Ω(err).Should(Equal(UnspecifiedValueErr))

	res = v.TimeDefault(def)
	Ω(res).Should(Equal(def))
}

func TestStringValueTime_Invalid(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res time.Time
	var err error

	def := time.Date(1860, 7, 2, 12, 0, 0, 0, time.UTC)

	v = StringValue{Value: pStr("wrong time value")}
	_, err = v.Time()
	Ω(err).ShouldNot(BeNil())
	Ω(err.Error()).Should(HavePrefix("parsing time "))

	res = v.TimeDefault(def)
	Ω(res).Should(Equal(def))
}

func pStr(str string) *string {
	return &str
}
