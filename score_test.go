package main

import (
	"testing"
)

func TestScore(t *testing.T) {

	var absTests = []struct {
		candidate Candidate
		query     string
		wanted    float32
	}{
		{candidate: NewCandidate("a"), query: "", wanted: 1.0},
		{candidate: NewCandidate("a"), query: "aa", wanted: 0.0},
		{candidate: NewCandidate("abcx"), query: "abcd", wanted: 0.0},
		{candidate: NewCandidate("axbxcx"), query: "aac", wanted: 0.0},
	}

	for _, test := range absTests {
		s := Score(&test.candidate, test.query)
		if s != test.wanted {
			t.Errorf("Score(%q, %q) = %f, wanted %f",
				test.candidate,
				test.query,
				s,
				test.wanted)
		}
	}

	var greaterTests = []struct {
		candidate Candidate
		query     string
		lower     float32
	}{
		{candidate: NewCandidate("abcx"), query: "abc", lower: 0.0},
		{candidate: NewCandidate("abcx"), query: "abc", lower: 0.0},
		{candidate: NewCandidate("aaa/bbb/File"), query: "abf", lower: 0.0},
		{candidate: NewCandidate("aaa/bbb/file"), query: "abF", lower: 0.0},
	}

	for _, test := range greaterTests {
		s := Score(&test.candidate, test.query)
		if s <= test.lower {
			t.Errorf("Score(%q, %q) = %f, wanted > %f",
				test.candidate,
				test.query,
				s,
				test.lower)
		}
	}

	var comparingTests = []struct {
		candidate1 Candidate
		query1     string

		candidate2 Candidate
		query2     string
	}{
		{NewCandidate("yxxxabxc"), "abc", NewCandidate("axxxybxc"), "abc"},
		{NewCandidate("xabc"), "abc", NewCandidate("long string abc"), "abc"},
	}

	for _, test := range comparingTests {
		s1 := Score(&test.candidate1, test.query1)
		s2 := Score(&test.candidate2, test.query2)
		if s1 <= s2 {
			t.Errorf("Score(%q, %q) <= Score(%q, %q), wanted >",
				test.candidate1,
				test.query1,
				test.candidate2,
				test.query2)
		}

	}

	var equalTests = []struct {
		candidate1 Candidate
		query1     string

		candidate2 Candidate
		query2     string
	}{
		{NewCandidate("abcxxxabxxxc"), "abc", NewCandidate("xabcxxxyyxxy"), "abc"},
		{NewCandidate("axxxabxc"), "abc", NewCandidate("yxxxabxc"), "abc"},
	}

	for _, test := range equalTests {
		s1 := Score(&test.candidate1, test.query1)
		s2 := Score(&test.candidate2, test.query2)
		if s1 != s2 {
			t.Errorf("Score(%q, %q) != Score(%q, %q), wanted =",
				test.candidate1,
				test.query1,
				test.candidate2,
				test.query2)
		}

	}

	var lessTests = []struct {
		candidate1 Candidate
		query1     string

		candidate2 Candidate
		query2     string
	}{
		{NewCandidate("axbxcxd"), "abcd", NewCandidate("AxBxCxD"), "abcd"},
		{NewCandidate("axbxcxd"), "abcd", NewCandidate("a/b/c/d"), "abcd"},
	}

	for _, test := range lessTests {
		s1 := Score(&test.candidate1, test.query1)
		s2 := Score(&test.candidate2, test.query2)
		if s1 >= s2 {
			t.Errorf("Score(%q, %q)[%q] >= Score(%q, %q)[%q], wanted <",
				test.candidate1,
				test.query1,
				s1,
				test.candidate2,
				test.query2,
				s2)

		}

	}

}
