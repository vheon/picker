package main

import (
	"errors"
	"sort"
	"strings"
)

func indexRuneStarting(str string, r rune, start int) int {
	idx := strings.IndexRune(str[start:], r)
	if idx == -1 {
		return -1
	}
	return idx + start
}

func allIndexRune(str string, r rune) []int {
	var allIndexRune []int
	for i, c := range str {
		if c == r {
			allIndexRune = append(allIndexRune, i)
		}
	}
	return allIndexRune
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
		start += 1
	}
	return runePositions, nil
}

type Match []int

func (m Match) Length() int {
	return m[len(m)-1] - m[0] + 1
}

type matchSlice []Match

func (ms matchSlice) Len() int           { return len(ms) }
func (ms matchSlice) Swap(i, j int)      { ms[i], ms[j] = ms[j], ms[i] }
func (ms matchSlice) Less(i, j int) bool { return ms[i].Length() < ms[j].Length() }

func bestMatch(ms []Match) Match {
	sort.Sort(matchSlice(ms))
	return ms[0]
}

func Score(candidate, query string) float64 {
	if len(query) == 0 {
		return 1.0
	}
	if len(candidate) < len(query) {
		return 0.0
	}

	firstQueryCharPositions := allIndexRune(candidate, rune(query[0]))

	var matches []Match
	for _, start := range firstQueryCharPositions {
		match, err := findMatch(candidate[start:], query)
		if err != nil {
			continue
		}
		matches = append(matches, match)
	}

	if len(matches) == 0 {
		return 0.0
	}

	best := bestMatch(matches)

	score := float64(best.Length())
	score = float64(len(query)) / score
	score = score / float64(len(candidate))

	return score
}
