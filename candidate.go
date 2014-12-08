package main

import (
	"bufio"
	"io"
)

type Candidate struct {
	value string
	score float32
}

func NewCandidate(s string) Candidate {
	return Candidate{
		value: s,
		score: 1.0,
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
