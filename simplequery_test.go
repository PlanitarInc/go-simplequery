package simplequery

import (
	"net/url"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestQHas(t *testing.T) {
	RegisterTestingT(t)

	urlQ, err := url.ParseQuery("k1=v1_1&k1=v1_2&k2=v2&k3")
	Ω(err).Should(BeNil())
	q := FromQuery(urlQ)

	Ω(q.Has("k1")).Should(BeTrue())
	Ω(q.Has("k2")).Should(BeTrue())
	Ω(q.Has("k3")).Should(BeTrue())
	Ω(q.Has("k9")).Should(BeFalse())
}

func TestQGet(t *testing.T) {
	RegisterTestingT(t)

	var v *StringValue

	urlQ, err := url.ParseQuery("k1=v1_1&k1=v1_2&k2=v2&k3")
	Ω(err).Should(BeNil())
	q := FromQuery(urlQ)

	v = q.Get("k1")
	Ω(v).Should(Equal(psv("v1_1")))

	v = q.Get("k2")
	Ω(v).Should(Equal(psv("v2")))

	v = q.Get("k3")
	Ω(v).Should(Equal(psv("")))

	v = q.Get("k9")
	Ω(v).Should(BeNil())
}

func TestQGetAll(t *testing.T) {
	RegisterTestingT(t)

	var v ValueSet

	urlQ, err := url.ParseQuery("k1=v1_1&k1=v1_2&k2=v2&k3")
	Ω(err).Should(BeNil())
	q := FromQuery(urlQ)

	v = q.GetAll("k1")
	Ω(v).Should(Equal(ValueSetFrom([]string{"v1_1", "v1_2"})))

	v = q.GetAll("k2")
	Ω(v).Should(Equal(ValueSetFrom([]string{"v2"})))

	v = q.GetAll("k3")
	Ω(v).Should(Equal(ValueSetFrom([]string{""})))

	v = q.GetAll("k9")
	Ω(v).Should(Equal(ValueSetFrom([]string{})))
}

func TestQFilterByKey(t *testing.T) {
	RegisterTestingT(t)

	urlQ, err := url.ParseQuery("k1=v1_1&k1=v1_2&k2=v2&k3")
	Ω(err).Should(BeNil())
	q := FromQuery(urlQ)

	var filtered Q

	filtered = q.FilterByKey(func(k string) bool {
		return true
	})
	Ω(filtered).Should(Equal(q))

	filtered = q.FilterByKey(func(k string) bool {
		return k != "k2"
	})
	Ω(filtered).Should(HaveKeyWithValue("k1",
		ValueSetFrom([]string{"v1_1", "v1_2"})))
	Ω(filtered).Should(HaveKeyWithValue("k3",
		ValueSetFrom([]string{""})))
	Ω(filtered).Should(HaveLen(2))
}

func TestQForEach(t *testing.T) {
	RegisterTestingT(t)

	urlQ, err := url.ParseQuery("k1=v1_1&k1=v1_2&k2=v2&k3")
	Ω(err).Should(BeNil())
	q := FromQuery(urlQ)

	res := map[string]ValueSet{}
	q.ForEach(func(key string, vs ValueSet) {
		res[key] = vs
	})

	Ω(res).Should(HaveKeyWithValue("k1",
		ValueSetFrom([]string{"v1_1", "v1_2"})))
	Ω(res).Should(HaveKeyWithValue("k2",
		ValueSetFrom([]string{"v2"})))
	Ω(res).Should(HaveKeyWithValue("k3",
		ValueSetFrom([]string{""})))
	Ω(res).Should(HaveLen(3))
}

func TestValueSetFrom(t *testing.T) {
	RegisterTestingT(t)

	var vs ValueSet

	vs = ValueSetFrom(nil)
	Ω(vs).Should(BeEmpty())

	vs = ValueSetFrom([]string{})
	Ω(vs).Should(BeEmpty())

	vs = ValueSetFrom([]string{"v2"})
	Ω(vs).Should(Equal(ValueSet{StringValue("v2")}))

	vs = ValueSetFrom([]string{"v1_1", "v1_2"})
	Ω(vs).Should(Equal(ValueSet{StringValue("v1_1"), StringValue("v1_2")}))
}

func TestValueSetEmpty(t *testing.T) {
	RegisterTestingT(t)

	var vs ValueSet

	vs = ValueSet{}
	Ω(vs.Empty()).Should(BeTrue())

	vs = ValueSet{StringValue("")}
	Ω(vs.Empty()).Should(BeFalse())

	vs = ValueSetFrom([]string{"v1_1", "v1_3"})
	Ω(vs.Empty()).Should(BeFalse())
}

func TestValueSetFirst(t *testing.T) {
	RegisterTestingT(t)

	var vs ValueSet

	vs = ValueSetFrom([]string{})
	Ω(vs.First()).Should(BeNil())

	vs = ValueSetFrom([]string{"v2"})
	Ω(vs.First()).Should(Equal(psv("v2")))

	vs = ValueSetFrom([]string{"v1_1", "v1_2"})
	Ω(vs.First()).Should(Equal(psv("v1_1")))
}

func TestValueSetIndex(t *testing.T) {
	RegisterTestingT(t)

	vs := ValueSetFrom([]string{"v1_1", "v1_2"})

	Ω(vs.Index(-1)).Should(BeNil())
	Ω(vs.Index(0)).Should(Equal(psv("v1_1")))
	Ω(vs.Index(1)).Should(Equal(psv("v1_2")))
	Ω(vs.Index(2)).Should(BeNil())
}

func TestValueSetStrings(t *testing.T) {
	RegisterTestingT(t)

	vs := ValueSetFrom([]string{"11", "f", "-1", "FALSE", "qwe"})
	Ω(vs.Strings()).Should(Equal([]string{"11", "f", "-1", "FALSE", "qwe"}))
}

func TestValueSetBools(t *testing.T) {
	RegisterTestingT(t)

	vs := ValueSetFrom([]string{"11", "1", "f", "off", "FALSE", "TRUE", "qwe", ""})
	Ω(vs.Bools()).Should(Equal([]bool{true, true, false, false, false, true, true, true}))
}

func TestValueSetInt64s(t *testing.T) {
	RegisterTestingT(t)

	vs := ValueSetFrom([]string{"11", "a", "-1"})
	Ω(vs.Int64s()).Should(Equal([]int64{11, 0, -1}))
}

func TestValueSetUint64s(t *testing.T) {
	RegisterTestingT(t)

	vs := ValueSetFrom([]string{"11", "a", "-1"})
	Ω(vs.Uint64s()).Should(Equal([]uint64{11, 0, 0}))
}

func TestValueSetFloat64s(t *testing.T) {
	RegisterTestingT(t)

	vs := ValueSetFrom([]string{"11", "a", "-1.2"})
	Ω(vs.Float64s()).Should(Equal([]float64{11, 0, -1.2}))
}

func TestValueSetTimes(t *testing.T) {
	RegisterTestingT(t)

	vs := ValueSetFrom([]string{"11", "2016-02-03T15:04:05Z", "-1.2"})
	Ω(vs.Times()).Should(Equal([]time.Time{
		time.Unix(11, 0).UTC(),
		time.Date(2016, 2, 3, 15, 4, 5, 0, time.UTC),
		time.Time{},
	}))
}

func TestStringValueString(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res string
	var err error

	v = StringValue("")

	Ω(v.String()).Should(Equal(""))

	res, err = v.ParseString()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(""))

	v = StringValue("value")

	Ω(v.String()).Should(Equal("value"))

	res, err = v.ParseString()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal("value"))
}

func TestStringValueString_Empty(t *testing.T) {
	RegisterTestingT(t)

	var v *StringValue
	var err error

	v = nil

	Ω(v.String()).Should(Equal(""))

	_, err = v.ParseString()
	Ω(err).Should(Equal(UnspecifiedValueErr))
}

func TestStringValueBool(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res bool
	var err error

	v = StringValue("1")

	res, err = v.ParseBool()
	Ω(err).Should(BeNil())
	Ω(res).Should(BeTrue())

	res = v.Bool()
	Ω(res).Should(BeTrue())

	res = v.Bool(true)
	Ω(res).Should(BeTrue())

	res = v.Bool(false)
	Ω(res).Should(BeTrue())

	v = StringValue("0")

	res, err = v.ParseBool()
	Ω(err).Should(BeNil())
	Ω(res).Should(BeFalse())

	res = v.Bool()
	Ω(res).Should(BeFalse())

	res = v.Bool(true)
	Ω(res).Should(BeFalse())

	res = v.Bool(false)
	Ω(res).Should(BeFalse())
}

func TestStringValueBool_Empty(t *testing.T) {
	RegisterTestingT(t)

	var v *StringValue
	var res bool
	var err error

	v = nil

	_, err = v.ParseBool()
	Ω(err).Should(Equal(UnspecifiedValueErr))

	// Flag is missing hence should return false
	res = v.Bool()
	Ω(res).Should(BeFalse())

	res = v.Bool(true)
	Ω(res).Should(BeTrue())

	res = v.Bool(false)
	Ω(res).Should(BeFalse())

	v = psv("")

	_, err = v.ParseBool()
	Ω(err).ShouldNot(BeNil())
	Ω(err.Error()).Should(Equal("unknown value"))

	// Flag is not missing but has empty value, hence should return true
	res = v.Bool()
	Ω(res).Should(BeTrue())

	res = v.Bool(true)
	Ω(res).Should(BeTrue())

	res = v.Bool(false)
	Ω(res).Should(BeFalse())
}

func TestStringValueBool_Invalid(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res bool
	var err error

	v = StringValue("qwe")

	res, err = v.ParseBool()
	Ω(err).ShouldNot(BeNil())
	Ω(err.Error()).Should(Equal("unknown value"))

	res = v.Bool()
	Ω(res).Should(BeTrue())

	res = v.Bool(true)
	Ω(res).Should(BeTrue())

	res = v.Bool(false)
	Ω(res).Should(BeFalse())
}

func TestStringValueInt64(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res int64
	var err error

	v = StringValue("12")

	res, err = v.ParseInt64()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(int64(12)))

	res = v.Int64()
	Ω(res).Should(Equal(int64(12)))

	res = v.Int64(-99)
	Ω(res).Should(Equal(int64(12)))

	v = StringValue("-12")

	res, err = v.ParseInt64()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(int64(-12)))

	res = v.Int64()
	Ω(res).Should(Equal(int64(-12)))

	res = v.Int64(-99)
	Ω(res).Should(Equal(int64(-12)))
}

func TestStringValueInt64_Empty(t *testing.T) {
	RegisterTestingT(t)

	var v *StringValue
	var res int64
	var err error

	v = nil

	_, err = v.ParseInt64()
	Ω(err).Should(Equal(UnspecifiedValueErr))

	res = v.Int64()
	Ω(res).Should(Equal(int64(0)))

	res = v.Int64(-99)
	Ω(res).Should(Equal(int64(-99)))

	v = psv("")

	_, err = v.ParseInt64()
	Ω(err).ShouldNot(BeNil())
	Ω(err.Error()).Should(HavePrefix("strconv.ParseInt: parsing "))

	res = v.Int64()
	Ω(res).Should(Equal(int64(0)))

	res = v.Int64(-99)
	Ω(res).Should(Equal(int64(-99)))
}

func TestStringValueInt64_Invalid(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res int64
	var err error

	v = StringValue("qwe")

	_, err = v.ParseInt64()
	Ω(err).ShouldNot(BeNil())
	Ω(err.Error()).Should(HavePrefix("strconv.ParseInt: parsing "))

	res = v.Int64()
	Ω(res).Should(Equal(int64(0)))

	res = v.Int64(-99)
	Ω(res).Should(Equal(int64(-99)))
}

func TestStringValueUint64(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res uint64
	var err error

	v = StringValue("12")

	res, err = v.ParseUint64()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(uint64(12)))

	res = v.Uint64()
	Ω(res).Should(Equal(uint64(12)))

	res = v.Uint64(99)
	Ω(res).Should(Equal(uint64(12)))
}

func TestStringValueUint64_Empty(t *testing.T) {
	RegisterTestingT(t)

	var v *StringValue
	var res uint64
	var err error

	v = nil

	_, err = v.ParseUint64()
	Ω(err).Should(Equal(UnspecifiedValueErr))

	res = v.Uint64()
	Ω(res).Should(Equal(uint64(0)))

	res = v.Uint64(99)
	Ω(res).Should(Equal(uint64(99)))

	v = psv("")

	_, err = v.ParseUint64()
	Ω(err).ShouldNot(BeNil())
	Ω(err.Error()).Should(HavePrefix("strconv.ParseUint: parsing "))

	res = v.Uint64()
	Ω(res).Should(Equal(uint64(0)))

	res = v.Uint64(99)
	Ω(res).Should(Equal(uint64(99)))
}

func TestStringValueUint64_Invalid(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res uint64
	var err error

	v = StringValue("qwe")

	_, err = v.ParseUint64()
	Ω(err).ShouldNot(BeNil())
	Ω(err.Error()).Should(HavePrefix("strconv.ParseUint: parsing "))

	res = v.Uint64()
	Ω(res).Should(Equal(uint64(0)))

	res = v.Uint64(99)
	Ω(res).Should(Equal(uint64(99)))
}

func TestStringValueFloat64(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res float64
	var err error

	v = StringValue("12.3")

	res, err = v.ParseFloat64()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(float64(12.3)))

	res = v.Float64()
	Ω(res).Should(Equal(float64(12.3)))

	res = v.Float64(-99.9)
	Ω(res).Should(Equal(float64(12.3)))

	v = StringValue("-12.3")

	res, err = v.ParseFloat64()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(float64(-12.3)))

	res = v.Float64()
	Ω(res).Should(Equal(float64(-12.3)))

	res = v.Float64(-99.9)
	Ω(res).Should(Equal(float64(-12.3)))
}

func TestStringValueFloat64_Empty(t *testing.T) {
	RegisterTestingT(t)

	var v *StringValue
	var res float64
	var err error

	v = nil

	_, err = v.ParseFloat64()
	Ω(err).Should(Equal(UnspecifiedValueErr))

	res = v.Float64()
	Ω(res).Should(Equal(float64(0)))

	res = v.Float64(-99.9)
	Ω(res).Should(Equal(float64(-99.9)))

	v = psv("")

	_, err = v.ParseFloat64()
	Ω(err).ShouldNot(BeNil())
	Ω(err.Error()).Should(HavePrefix("strconv.ParseFloat: parsing "))

	res = v.Float64()
	Ω(res).Should(Equal(float64(0)))

	res = v.Float64(-99.9)
	Ω(res).Should(Equal(float64(-99.9)))
}

func TestStringValueFloat64_Invalid(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res float64
	var err error

	v = StringValue("qwe")

	_, err = v.ParseFloat64()
	Ω(err).ShouldNot(BeNil())
	Ω(err.Error()).Should(HavePrefix("strconv.ParseFloat: parsing "))

	res = v.Float64()
	Ω(res).Should(Equal(float64(0)))

	res = v.Float64(-99.9)
	Ω(res).Should(Equal(float64(-99.9)))
}

func TestStringValueTime(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res, expTime time.Time
	var err error

	def := time.Date(1860, 7, 2, 12, 0, 0, 0, time.UTC)

	v = StringValue("123")
	expTime = time.Unix(123, 0).UTC()

	res, err = v.ParseTime()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(expTime))

	res = v.Time()
	Ω(res).Should(Equal(expTime))

	res = v.Time(def)
	Ω(res).Should(Equal(expTime))

	v = StringValue("252460800000")
	expTime = time.Date(1978, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err = v.ParseTime()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(expTime))

	res = v.Time()
	Ω(res).Should(Equal(expTime))

	res = v.Time(def)
	Ω(res).Should(Equal(expTime))

	v = StringValue("2016-02-03T15:04:05Z")
	expTime = time.Date(2016, 2, 3, 15, 4, 5, 0, time.UTC)

	res, err = v.ParseTime()
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(expTime))

	res = v.Time()
	Ω(res).Should(Equal(expTime))

	res = v.Time(def)
	Ω(res).Should(Equal(expTime))
}

func TestStringValueTime_Empty(t *testing.T) {
	RegisterTestingT(t)

	var v *StringValue
	var res time.Time
	var err error

	def := time.Date(1860, 7, 2, 12, 0, 0, 0, time.UTC)

	v = nil

	_, err = v.ParseTime()
	Ω(err).Should(Equal(UnspecifiedValueErr))

	res = v.Time()
	Ω(res).Should(Equal(time.Time{}))

	res = v.Time(def)
	Ω(res).Should(Equal(def))

	v = psv("")

	_, err = v.ParseTime()
	Ω(err).ShouldNot(BeNil())
	Ω(err.Error()).Should(HavePrefix("parsing time"))

	res = v.Time()
	Ω(res).Should(Equal(time.Time{}))

	res = v.Time(def)
	Ω(res).Should(Equal(def))
}

func TestStringValueTime_Invalid(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue
	var res time.Time
	var err error

	def := time.Date(1860, 7, 2, 12, 0, 0, 0, time.UTC)

	v = StringValue("wrong time value")
	_, err = v.ParseTime()
	Ω(err).ShouldNot(BeNil())
	Ω(err.Error()).Should(HavePrefix("parsing time "))

	res = v.Time()
	Ω(res).Should(Equal(time.Time{}))

	res = v.Time(def)
	Ω(res).Should(Equal(def))
}

func TestStringValueList(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue

	v = StringValue("123")
	Ω(v.List()).Should(Equal(ValueSet{
		StringValue("123"),
	}))

	v = StringValue("1,2,3")
	Ω(v.List()).Should(Equal(ValueSet{
		StringValue("1"),
		StringValue("2"),
		StringValue("3"),
	}))
}

func TestStringValueList_Empty(t *testing.T) {
	RegisterTestingT(t)

	var v *StringValue

	v = nil
	Ω(v.List()).Should(BeEmpty())

	v = psv("")
	Ω(v.List()).Should(Equal(ValueSet{
		StringValue(""),
	}))
}

func TestStringValueList_CustomSeparator(t *testing.T) {
	RegisterTestingT(t)

	var v StringValue

	v = StringValue("1|2|3")

	Ω(v.List("|")).Should(Equal(ValueSet{
		StringValue("1"),
		StringValue("2"),
		StringValue("3"),
	}))

	Ω(v.List("2")).Should(Equal(ValueSet{
		StringValue("1|"),
		StringValue("|3"),
	}))
}

func psv(str string) *StringValue {
	v := StringValue(str)
	return &v
}
