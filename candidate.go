package main

import (
	"bufio"
	"io"
	"unicode"
	"unicode/utf8"
)

type Candidate struct {
	value      string
	boundaries []rune
	score      float32
}

func FindBoundaries(s string) []rune {
	var res []rune

	first, size := utf8.DecodeRuneInString(s)
	if !unicode.IsPunct(first) {
		res = append(res, first)
	}

	prev := first
	for _, r := range s[size:] {
		is_good_uppercase := unicode.IsUpper(r) && !unicode.IsUpper(prev)
		is_alpha_after_punctuation := unicode.IsPunct(prev) && unicode.In(r, unicode.Number, unicode.Letter)

		if is_good_uppercase || is_alpha_after_punctuation {
			res = append(res, r)
		}
		prev = r
	}
	return res
}

func NewCandidate(s string) Candidate {
	return Candidate{
		value:      s,
		boundaries: FindBoundaries(s),
		score:      1.0,
	}
}

func readAllCandidates(r io.Reader) []Candidate {
	scanner := bufio.NewScanner(r)
	var candidates []Candidate
	for scanner.Scan() {
		candidates = append(candidates, NewCandidate(scanner.Text()))
	}
	return candidates
}

type CandidateSlice []Candidate

func (cs CandidateSlice) Len() int           { return len(cs) }
func (cs CandidateSlice) Swap(i, j int)      { cs[i], cs[j] = cs[j], cs[i] }
func (cs CandidateSlice) Less(i, j int) bool { return cs[i].score > cs[j].score }
