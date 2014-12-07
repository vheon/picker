package main

import (
	"bufio"
	"io"
	"sort"
	"strings"
	"sync"
	"unicode/utf8"
)

type KeepBottomStack []int

func NewKeepBottomStack(d int) KeepBottomStack {
	return []int{d}
}

func (s KeepBottomStack) Peek() int   { return s[len(s)-1] }
func (s *KeepBottomStack) Push(i int) { (*s) = append((*s), i) }
func (s *KeepBottomStack) Clear()     { (*s) = (*s)[:1] }
func (s *KeepBottomStack) Drop() {
	if len(*s) > 1 {
		(*s) = (*s)[:len(*s)-1]
	}
}

type Picker struct {
	all   []Candidate
	valid KeepBottomStack

	prompt string
	query  string

	index  int
	height int
	width  int

	blank []Candidate

	view []string
}

func NewPicker(prompt string, height, width int, r io.Reader) *Picker {
	candidates := readAllCandidates(r)

	blank := make([]Candidate, height)
	copy(blank, candidates[:min(height, len(candidates))])

	return &Picker{
		all: candidates,
		// create the stack with the first value in
		valid: NewKeepBottomStack(len(candidates)),

		prompt: prompt,
		query:  "",

		index:  0,
		height: height,
		width:  width,

		blank: blank,

		view: make([]string, height),
	}
}

func cutAt(str string, width int) string {
	if len(str) > width {
		return str[:width]
	}
	return str
}

func (p *Picker) View() string {
	firstLine := p.prompt + p.query + "\n"
	candidates := p.all
	if p.query == "" {
		candidates = p.blank
	}

	for i := range p.view {
		if i < len(candidates) && candidates[i].score > 0.0 {
			p.view[i] = cutAt(candidates[i].value, p.width)
		} else {
			p.view[i] = ""
		}
	}
	p.view[p.index] = TTYReverse(p.view[p.index])

	return firstLine + strings.Join(p.view, "\n")
}

func (p *Picker) Sort() {
	if p.query == "" {
		return
	}

	// peek the top from the stack
	candidates := p.all[:p.valid.Peek()]

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

	// push the value on the stack
	p.valid.Push(sort.Search(len(candidates), func(i int) bool {
		return candidates[i].score == 0.0
	}))
}

func (p *Picker) Selected() string {
	return p.all[p.index].value
}

func (p *Picker) Up() {
	if p.index > 0 {
		p.index -= 1
	}
}

func (p *Picker) Down() {
	if p.index < p.height-1 && p.index < p.valid.Peek()-1 {
		p.index += 1
	}
}

func (p *Picker) AppendToQuery(r rune) {
	p.query += string(r)
	p.index = 0
}

func (p *Picker) Backspace() {
	_, size := utf8.DecodeLastRuneInString(p.query)
	p.query = p.query[:len(p.query)-size]

	// reset the index
	p.index = 0

	// Drop the value of valid candidates on the stack
	p.valid.Drop()

	p.Sort()
}

func (p *Picker) Clear() {
	p.query = ""
	// Clear the stack except the bottom value
	p.valid.Clear()
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
