package main

import "fmt"

const (
	ClearLine  string = "\033[2K"
	Move              = "\033[%d;%dH"
	Invert            = "\033[7m"
	Reset             = "\033[0m"
	HideCursor        = "\033[?25l"
	ShowCursor        = "\033[?25h"
)

type Screen struct {
	nLines int
	tty    *TTY
}

func NewScreen(tty *TTY, visibleLines int) *Screen {
	return &Screen{
		nLines: visibleLines,
		tty:    tty,
	}
}

func (s *Screen) ConfigScreen() {
	s.tty.Stty("-echo", "-icanon")
}

func (s *Screen) MakeRoom() {
	for i := 0; i < s.nLines; i++ {
		s.tty.Puts("")
	}
}

func (s *Screen) HideCursor() {
	s.tty.Write(HideCursor)
}

func (s *Screen) ShowCursor() {
	s.tty.Write(ShowCursor)
}

func (s *Screen) MoveTo(x, y int) {
	s.tty.Write(fmt.Sprintf(Move, x, y))
}

func (s *Screen) MoveToRow(x int) {
	s.tty.Write(fmt.Sprintf(Move, x, 1))
}

func (s *Screen) Draw(view *View) {
	s.HideCursor()
	promptRow := s.tty.Height - s.nLines
	for i, row := range ttyView(view) {
		s.MoveToRow(promptRow + i)
		s.tty.Write(ClearLine)
		s.tty.Write(row)
	}
	// XXX: 3 magic number
	s.MoveTo(promptRow, len(view.Query)+3)
	s.ShowCursor()
}

func ansiInverted(s string) string {
	return Invert + s + Reset
}

func ttyView(view *View) []string {
	rows := make([]string, view.Height+1)
	rows[0] = "> " + view.Query
	copy(rows[1:], view.Rows)
	if len(view.Rows) > 0 {
		rows[view.Index()+1] = ansiInverted(view.Selected())
	}
	return rows
}
