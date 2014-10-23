package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	EraseLine  string = "\033[2K"
	Move              = "\033[%d;%dH"
	Invert            = "\033[7m"
	Reset             = "\033[0m"
	HideCursor        = "\033[?25l"
	ShowCursor        = "\033[?25h"
)

type Screen struct {
	Height int
	Width  int
	tty    *TTY
}

func NewScreen(tty *TTY) *Screen {
	height, width := parseSize(tty.Stty("size"))

	return &Screen{
		tty:    tty,
		Height: height,
		Width:  width,
	}
}

func parseSize(size string) (int, int) {
	ssize := strings.Fields(size)
	height, _ := strconv.Atoi(ssize[0])
	width, _ := strconv.Atoi(ssize[1])
	return height, width
}

func (s *Screen) ConfigScreen() {
	s.tty.Stty("-echo", "-icanon")
}

func (s *Screen) MakeRoom(rows int) {
	for i := 0; i < rows; i++ {
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
	promptRow := s.Height - view.Height
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
