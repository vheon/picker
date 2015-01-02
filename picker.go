package main

import (
	"io"
	"sort"
	"strings"
	"sync"
	"unicode/utf8"
)

type Stack []int

func (s *Stack) Empty() bool       { return len(*s) == 0 }
func (s *Stack) Peek() int         { return (*s)[len(*s)-1] }
func (s *Stack) Push(i int)        { (*s) = append((*s), i) }
func (s *Stack) ClearUntilBottom() { (*s) = (*s)[:1] }
func (s *Stack) DropExceptBottom() {
	if len(*s) > 1 {
		(*s) = (*s)[:len(*s)-1]
	}
}

type PickerView struct {
	firstLine string
	lines     string
}

func (pv *PickerView) String() string {
	return pv.firstLine + "\n" + pv.lines
}

type Picker struct {
	originals []Candidate
	all       []Candidate
	validSize Stack

	prompt string
	query  string

	index  int
	height int
	width  int

	view []string

	render chan *PickerView
}

func NewPicker(prompt string, height, width int, r io.Reader, renderChan chan *PickerView) *Picker {
	candidates := readAllCandidates(r)

	blank := make([]Candidate, height)
	copy(blank, candidates[:min(height, len(candidates))])

	picker := &Picker{
		all: candidates,
		// create the stack with the first value in
		validSize: Stack([]int{len(candidates)}),

		prompt: prompt,
		query:  "",

		index:  0,
		height: height,
		width:  width,

		originals: blank,

		view: make([]string, height),

		render: renderChan,
	}

	// render the first frame
	renderChan <- picker.View()

	return picker
}

func cutAt(str string, width int) string {
	if len(str) > width {
		return str[:width]
	}
	return str
}

func TTYReverse(str string) string {
	return string(ReverseColor) + str + string(ResetColor)
}

func (p *Picker) View() *PickerView {
	firstLine := p.prompt + p.query
	candidates := p.all
	if p.query == "" {
		candidates = p.originals
	}

	for i := range p.view {
		if i < len(candidates) && candidates[i].score > 0.0 {
			p.view[i] = cutAt(candidates[i].value, p.width)
		} else {
			p.view[i] = ""
		}
	}
	p.view[p.index] = TTYReverse(p.view[p.index])

	return &PickerView{firstLine, strings.Join(p.view, "\n")}
}

func (p *Picker) Sort() {
	if p.query == "" {
		return
	}

	// peek the top from the stack
	candidates := p.all[:p.validSize.Peek()]

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
				c.score = Score(c, p.query)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	sort.Sort(CandidateSlice(candidates))
}

func (p *Picker) Selected() string {
	return p.all[p.index].value
}

func (p *Picker) Up() {
	if p.index > 0 {
		p.index -= 1
	}

	p.render <- p.View()
}

func (p *Picker) Down() {
	if p.index < p.height-1 && p.index < p.validSize.Peek()-1 {
		p.index += 1
	}

	p.render <- p.View()
}

func (p *Picker) More(r rune) {
	p.query += string(r)
	p.index = 0

	p.Sort()

	// push the value on the stack
	p.validSize.Push(sort.Search(p.validSize.Peek(), func(i int) bool {
		return p.all[i].score == 0.0
	}))

	p.render <- p.View()
}

func (p *Picker) Back() {
	_, size := utf8.DecodeLastRuneInString(p.query)
	p.query = p.query[:len(p.query)-size]

	// reset the index
	p.index = 0

	// Drop the value of valid candidates on the stack
	p.validSize.DropExceptBottom()

	p.Sort()

	p.render <- p.View()
}

func (p *Picker) Clear() {
	p.query = ""
	// Clear the stack except the bottom value
	p.validSize.ClearUntilBottom()

	p.render <- p.View()
}
