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
	all     []string
	visible int
}

func NewPicker(candidates []string, visible int) *Picker {
	return &Picker{
		all:     candidates,
		visible: visible,
	}
}

type lessfn func(int, int) bool

type sortableCandidates struct {
	list []string
	by   lessfn
}

func (c sortableCandidates) Len() int           { return len(c.list) }
func (c sortableCandidates) Swap(i, j int)      { c.list[i], c.list[j] = c.list[j], c.list[i] }
func (c sortableCandidates) Less(i, j int) bool { return c.by(i, j) }

func scoreByQuery(candidates []string, query string) lessfn {
	return func(i, j int) bool {
		return Score(candidates[i], query) > Score(candidates[j], query)
	}
}

func fisrtZero(c []string, query string) int {
	for i, str := range c {
		if Score(str, query) == 0.0 {
			return i
		}
	}
	return len(c)
}

func (p *Picker) Answer(query string) *View {
	candidates := &sortableCandidates{
		list: p.all,
		by:   scoreByQuery(p.all, query),
	}
	sort.Sort(candidates)

	//XXX: We are rescoring all the candidates
	spread := fisrtZero(p.all, query)
	if spread > p.visible {
		spread = p.visible
	}

	return &View{
		index:  0,
		height: p.visible,
		lines:  p.all[:spread],
		query:  query,
	}

	// score `all`
	// sort `all`
	// pick best `n` candidates
	// restore `index` to 0
}
