package main

import (
	"bufio"
	"io"
	"sort"
	"sync"
)

type Picker struct {
	all    []Candidate
	valid  int
	prompt string
	query  string
	view   View
	blank  View
}

func NewPicker(prompt string, height, width int, r io.Reader) *Picker {
	candidates := readAllCandidates(r)

	view := NewView(width, height)
	view.Fill(candidates)

	blank := NewView(width, height)
	CopyView(blank, view)

	return &Picker{
		all:    candidates,
		valid:  len(candidates),
		prompt: prompt,
		query:  "",
		view:   view,
		blank:  blank,
	}
}

func (p *Picker) String() string {
	return p.prompt + p.query + "\n" + p.view.String()
}

func (cs CandidateSlice) ToChan() <-chan *Candidate {
	ch := make(chan *Candidate)
	go func() {
		for i := range cs {
			ch <- &cs[i]
		}
		close(ch)
	}()
	return ch
}

func (p *Picker) Sort() {
	if p.query == "" {
		CopyView(p.view, p.blank)
		return
	}

	candidates := p.all[:p.valid]

	ch := make(chan *Candidate)
	go func() {
		for i := range candidates {
			ch <- &candidates[i]
		}
		close(ch)
	}()

	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			for c := range ch {
				c.score = Score(c.value, p.query)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	sort.Sort(CandidateSlice(candidates))

	p.valid = sort.Search(len(candidates), func(i int) bool {
		return candidates[i].score == 0.0
	})

	p.view.Fill(p.all[:p.valid])
}

// type View struct {
// 	Height int
// 	Query  string
// 	Rows   []string
// 	Done   bool

// 	index  int
// 	prompt string
// }

// func (v *View) Index() int {
// 	return v.index
// }

// func (v *View) Selected() string {
// 	if len(v.Rows)-1 < v.index {
// 		return ""
// 	}
// 	return v.Rows[v.index]
// }

// func (v *View) Down() {
// 	if v.index < len(v.Rows)-1 {
// 		v.index++
// 	}
// }

// func (v *View) Up() {
// 	if v.index > 0 {
// 		v.index--
// 	}
// }

// func (v *View) ClearPrompt() {
// 	v.Query = ""
// 	v.prompt = ""
// }

// func (v *View) DrawOnTerminal(t *terminal.Terminal) {
// 	t.HideCursor()
// 	defer t.ShowCursor()

// 	start_row := t.Height - v.Height - 1

// 	for i, row := range v.toAnsiForm(t) {
// 		t.WriteToLine(start_row+i, row)
// 	}

// 	t.MoveTo(start_row, len(v.Query)+len(v.prompt))
// }

// func (v *View) toAnsiForm(t *terminal.Terminal) []string {
// 	rows := make([]string, v.Height+1)
// 	rows[0] = v.prompt + v.Query
// 	for i, row := range v.Rows {
// 		rows[i+1] = cutAt(row, t.Width)
// 	}
// 	if len(v.Rows) > 0 {
// 		rows[v.Index()+1] = terminal.AnsiInverted(v.Selected())
// 	}
// 	return rows
// }

// func cutAt(s string, width int) string {
// 	if len(s) > width {
// 		return s[:width]
// 	}
// 	return s
// }

// type Picker struct {
// 	all     []Candidate
// 	visible int
// 	blank   *View
// }

// func NewPicker(candidates []Candidate, visible int) *Picker {
// 	picker := &Picker{
// 		all:     candidates,
// 		visible: visible,
// 	}
// 	picker.blank = picker.doAnswer("")

// 	return picker
// }

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

// func (p *Picker) Answer(query string) *View {
// 	if query == "" {
// 		return p.blank
// 	}
// 	return p.doAnswer(query)
// }

// func (p *Picker) doAnswer(query string) *View {
// 	all := make(chan *Candidate)
// 	go func() {
// 		for i := range p.all {
// 			all <- &p.all[i]
// 		}
// 		close(all)
// 	}()

// 	var wg sync.WaitGroup
// 	for i := 0; i < 64; i++ {
// 		wg.Add(1)
// 		go func() {
// 			for c := range all {
// 				c.score = Score(c.value, query)
// 			}
// 			wg.Done()
// 		}()
// 	}
// 	wg.Wait()

// 	sort.Sort(CandidateSlice(p.all))

// 	lines := []string{}
// 	for i, c := range p.all {
// 		if c.score == 0.0 || i >= p.visible {
// 			break
// 		}
// 		lines = append(lines, c.value)
// 	}

// 	return &View{
// 		Height: p.visible,
// 		Rows:   lines,
// 		Query:  query,
// 		Done:   false,
// 		index:  0,
// 		prompt: "> ",
// 	}
// }
