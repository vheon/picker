package main

import (
	"testing"
)

func TestScore(t *testing.T) {

	var absTests = []struct {
		candidate string
		query     string
		wanted    float64
	}{
		{candidate: "a", query: "", wanted: 1.0},
		{candidate: "a", query: "aa", wanted: 0.0},
		{candidate: "abcx", query: "abcd", wanted: 0.0},
		{candidate: "axbxcx", query: "aac", wanted: 0.0},
	}

	for _, test := range absTests {
		s := Score(test.candidate, test.query)
		if s != test.wanted {
			t.Errorf("Score(%q, %q) = %f, wanted %f",
				test.candidate,
				test.query,
				s,
				test.wanted)
		}
	}

	var greaterTests = []struct {
		candidate string
		query     string
		lower     float64
	}{
		{candidate: "abcx", query: "abc", lower: 0.0},
		{candidate: "abcx", query: "abc", lower: 0.0},
		{candidate: "aaa/bbb/File", query: "abf", lower: 0.0},
		{candidate: "aaa/bbb/file", query: "abF", lower: 0.0},
	}

	for _, test := range greaterTests {
		s := Score(test.candidate, test.query)
		if s <= test.lower {
			t.Errorf("Score(%q, %q) = %f, wanted > %f",
				test.candidate,
				test.query,
				s,
				test.lower)
		}
	}

	var comparingTests = []struct {
		candidate1 string
		query1     string

		candidate2 string
		query2     string
	}{
		{"yxxxabxc", "abc", "axxxybxc", "abc"},
		{"xabc", "abc", "long string abc", "abc"},
	}

	for _, test := range comparingTests {
		s1 := Score(test.candidate1, test.query1)
		s2 := Score(test.candidate2, test.query2)
		if s1 <= s2 {
			t.Errorf("Score(%q, %q) <= Score(%q, %q), wanted >",
				test.candidate1,
				test.query1,
				test.candidate2,
				test.query2)
		}

	}

	var equalTests = []struct {
		candidate1 string
		query1     string

		candidate2 string
		query2     string
	}{
		{"abcxxxabxxxc", "abc", "xabcxxxyyxxy", "abc"},
		{"axxxabxc", "abc", "yxxxabxc", "abc"},
	}

	for _, test := range equalTests {
		s1 := Score(test.candidate1, test.query1)
		s2 := Score(test.candidate2, test.query2)
		if s1 != s2 {
			t.Errorf("Score(%q, %q) != Score(%q, %q), wanted =",
				test.candidate1,
				test.query1,
				test.candidate2,
				test.query2)
		}

	}
}
