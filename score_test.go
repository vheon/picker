package main

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestScore_WithEmptyQuery_ReturnOne(t *testing.T) {
	RegisterTestingT(t)

	Expect(Score("a", "")).To(Equal(1.0))
}

func TestScore_WithQueryLongerThanCandidate_ReturnZero(t *testing.T) {
	RegisterTestingT(t)

	Expect(Score("a", "aa")).To(Equal(0.0))
}

func TestScore_AllTheLettersInQueryMustBePresent(t *testing.T) {
	RegisterTestingT(t)

	Expect(Score("abcx", "abcd")).To(Equal(0.0))
	Expect(Score("abcx", "abc")).To(BeNumerically(">", 0.0))
	Expect(Score("abcx", "abc")).To(BeNumerically(">", 0.0))
	Expect(Score("axbcx", "abc")).To(BeNumerically(">", 0.0))
	Expect(Score("axbcx", "aa")).To(Equal(0.0))
}

func TestScore_PreferCompactMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(Score("yxxxabxc", "abc")).To(BeNumerically(">", Score("axxxybxc", "abc")))
	Expect(Score("axxxabxc", "abc")).To(Equal(Score("yxxxabxc", "abc")))

	Expect(Score("abcxxxabxxxc", "abc")).To(Equal(Score("abcxxxyyxxxy", "abc")))
}

func TestScore_PreferShorterCandidate(t *testing.T) {
	RegisterTestingT(t)

	Expect(Score("long string abc", "abc")).To(BeNumerically("<", Score("xabc", "abc")))
}
