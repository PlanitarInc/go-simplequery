package util

import (
	"testing"

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
