package main

import (
	"sort"
	"sync"
)

type View struct {
	Height int
	Query  string
	Rows   []string

	index int
}

func (v *View) Index() int {
	return v.index
}

func (v *View) Selected() string {
	return v.Rows[v.index]
}

func (v *View) Down() {
	if v.index < len(v.Rows)-1 {
		v.index++
	}
}

func (v *View) Up() {
	if v.index > 0 {
		v.index--
	}
}

type Picker struct {
	all     []Candidate
	visible int
	blank   *View
}

func NewPicker(candidates []Candidate, visible int) *Picker {
	picker := &Picker{
		all:     candidates,
		visible: visible,
	}
	blank := picker.doAnswer("")
	picker.blank = blank

	return picker
}

type Candidate struct {
	value string
	score float64
}

func NewCandidate(s string) Candidate {
	return Candidate{
		value: s,
		score: 1.0,
	}
}

type CandidateSlice []Candidate

func (cs CandidateSlice) Len() int           { return len(cs) }
func (cs CandidateSlice) Swap(i, j int)      { cs[i], cs[j] = cs[j], cs[i] }
func (cs CandidateSlice) Less(i, j int) bool { return cs[i].score > cs[j].score }

func (p *Picker) Answer(query string) *View {
	if query == "" {
		return p.blank
	}
	return p.doAnswer(query)
}

func (p *Picker) doAnswer(query string) *View {
	var wg sync.WaitGroup
	wg.Add(len(p.all))
	for i := range p.all {
		candidate := &p.all[i]
		go func(c *Candidate) {
			c.score = Score(c.value, query)
			wg.Done()
		}(candidate)
	}
	wg.Wait()

	sort.Sort(CandidateSlice(p.all))

	lines := []string{}
	for i := 0; i < p.visible; i++ {
		if p.all[i].score == 0.0 {
			break
		}
		lines = append(lines, p.all[i].value)
	}

	return &View{
		index:  0,
		Height: p.visible,
		Rows:   lines,
		Query:  query,
	}
}
