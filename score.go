package main

import (
	"errors"
	"strings"
	"unicode/utf8"
)

func indexRuneStarting(str string, r rune, start int) int {
	idx := strings.IndexRune(str[start:], r)
	if idx == -1 {
		return -1
	}
	return idx + start
}

func indexesRune(str string, r rune) []int {
	var indexes []int
	for i, c := range str {
		if c == r {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func findMatch(candidate string, query string) (Match, error) {
	var runePositions Match
	start := 0
	for _, r := range query {
		start = indexRuneStarting(candidate, r, start)
		if start == -1 {
			return nil, errors.New("No Match Found")
		}
		runePositions = append(runePositions, start)
		start += utf8.RuneLen(r)
	}
	return runePositions, nil
}

type Match []int

func (m Match) Length() int {
	return m[len(m)-1] - m[0] + 1
}

func bestMatch(ms []Match) Match {
	if len(ms) == 0 {
		return nil
	}
	best := ms[0]
	for _, m := range ms {
		if m.Length() < best.Length() {
			best = m
		}
	}
	return best
}

func Score(candidate, query string) float32 {
	if len(query) == 0 {
		return 1.0
	}
	if len(candidate) < len(query) {
		return 0.0
	}

	candidate = strings.ToLower(candidate)
	query = strings.ToLower(query)

	first, _ := utf8.DecodeRuneInString(query)
	firstQueryRunePositions := indexesRune(candidate, first)

	var matches []Match
	for _, start := range firstQueryRunePositions {
		match, err := findMatch(candidate[start:], query)
		if err != nil {
			continue
		}
		matches = append(matches, match)
	}

	if len(matches) == 0 {
		return 0.0
	}

	var score float32
	score = float32(bestMatch(matches).Length())
	score = float32(len(query)) / score
	score = score / float32(len(candidate))

	return score
}
