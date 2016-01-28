package util

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestParseBool(t *testing.T) {
	RegisterTestingT(t)

	var res bool
	var err error

	res, err = ParseBool("0")
	Ω(err).Should(BeNil())
	Ω(res).Should(BeFalse())

	res, err = ParseBool("1")
	Ω(err).Should(BeNil())
	Ω(res).Should(BeTrue())

	res, err = ParseBool("T")
	Ω(err).Should(BeNil())
	Ω(res).Should(BeTrue())

	res, err = ParseBool("F")
	Ω(err).Should(BeNil())
	Ω(res).Should(BeFalse())

	res, err = ParseBool("true")
	Ω(err).Should(BeNil())
	Ω(res).Should(BeTrue())

	res, err = ParseBool("false")
	Ω(err).Should(BeNil())
	Ω(res).Should(BeFalse())

	res, err = ParseBool("on")
	Ω(err).Should(BeNil())
	Ω(res).Should(BeTrue())

	res, err = ParseBool("off")
	Ω(err).Should(BeNil())
	Ω(res).Should(BeFalse())

	res, err = ParseBool("On")
	Ω(err).Should(BeNil())
	Ω(res).Should(BeTrue())

	res, err = ParseBool("Off")
	Ω(err).Should(BeNil())
	Ω(res).Should(BeFalse())
}

func TestParseBool_Invalid(t *testing.T) {
	RegisterTestingT(t)

	var err error

	_, err = ParseBool("")
	Ω(err).Should(Equal(unknownValueError))

	_, err = ParseBool("00")
	Ω(err).Should(Equal(unknownValueError))

	_, err = ParseBool("10")
	Ω(err).Should(Equal(unknownValueError))

	_, err = ParseBool("a")
	Ω(err).Should(Equal(unknownValueError))
}

func TestParseTime(t *testing.T) {
	RegisterTestingT(t)

	var res time.Time
	var err error

	res, err = ParseTime("0")
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)))

	res, err = ParseTime("1445486493")
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(time.Date(2015, 10, 22, 4, 1, 33, 0, time.UTC)))

	res, err = ParseTime("2524608000")
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(time.Date(2050, 1, 1, 0, 0, 0, 0, time.UTC)))

	res, err = ParseTime("252460800000")
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(time.Date(1978, 1, 1, 0, 0, 0, 0, time.UTC)))

	res, err = ParseTime("2016-02-03T15:04:05Z")
	Ω(err).Should(BeNil())
	Ω(res).Should(Equal(time.Date(2016, 2, 3, 15, 4, 5, 0, time.UTC)))
}

func TestParseTime_Invalid(t *testing.T) {
	RegisterTestingT(t)

	var err error

	_, err = ParseTime("1445486493A")
	Ω(err).ShouldNot(BeNil())

	_, err = ParseTime("")
	Ω(err).ShouldNot(BeNil())

	_, err = ParseTime("2016.02.03")
	Ω(err).ShouldNot(BeNil())
}
