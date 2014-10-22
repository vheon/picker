package main

import "sort"

type View struct {
	height int
	query  string
	lines  []string
	index  int
}

func (v *View) Index() int {
	return v.index
}

func (v *View) Selected() string {
	return v.lines[v.index]
}

func (v *View) Query() string {
	return v.query
}

func (v *View) Down() *View {
	if v.index < len(v.lines)-1 {
		v.index++
	}
	return v
}

func (v *View) Up() {
	if v.index > 0 {
		v.index--
	}
}

type Picker struct {
	all     []Candidate
	visible int
}

func NewPicker(candidates []string, visible int) *Picker {

	all := make([]Candidate, len(candidates))
	for i, c := range candidates {
		all[i] = NewCandidate(c)
	}

	return &Picker{
		all:     all,
		visible: visible,
	}
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
	for i, c := range p.all {
		p.all[i].score = Score(c.value, query)
	}
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
		height: p.visible,
		lines:  lines,
		query:  query,
	}
}
