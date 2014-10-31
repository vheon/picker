package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	Ctrl_N    byte = 14
	Ctrl_P         = 16
	Ctrl_U         = 21
	Ctrl_W         = 23
	Backspace      = 127
	LF             = 10
)

const (
	EraseLine  string = "\033[2K"
	Move              = "\033[%d;%dH"
	Invert            = "\033[7m"
	Reset             = "\033[0m"
	HideCursor        = "\033[?25l"
	ShowCursor        = "\033[?25h"
)

type Terminal struct {
	Height int
	Width  int
	tty    *TTY
}

func NewTerminal(tty *TTY) *Terminal {
	height, width := parseSize(tty.Stty("size"))

	return &Terminal{
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

func (t *Terminal) ConfigTerminal() {
	t.tty.Stty("-echo", "-icanon")
}

func (t *Terminal) MakeRoom(rows int) {
	for i := 0; i < rows; i++ {
		t.tty.Write("\n")
	}
}

func (t *Terminal) HideCursor() {
	t.tty.Write(HideCursor)
}

func (t *Terminal) ShowCursor() {
	t.tty.Write(ShowCursor)
}

func (t *Terminal) MoveTo(x, y int) {
	t.tty.Write(fmt.Sprintf(Move, x+1, y+1))
}

func (t *Terminal) MoveToRow(x int) {
	t.MoveTo(x, 0)
}

func (t *Terminal) MoveBottom() {
	t.MoveToRow(t.Height)
}

func (t *Terminal) Draw(view *View) {
	t.HideCursor()
	defer t.ShowCursor()

	start_row := t.Height - view.Height - 1

	for i, row := range ttyView(view, t.Width) {
		t.MoveToRow(start_row + i)
		t.tty.Write(EraseLine)
		t.tty.Write(row)
	}

	t.MoveTo(start_row, len(view.Query)+len(view.prompt))
}

func ansiInverted(s string) string {
	return Invert + s + Reset
}

func ttyView(view *View, width int) []string {
	rows := make([]string, view.Height+1)
	rows[0] = view.prompt + view.Query
	for i, row := range view.Rows {
		if len(row) > width {
			row = row[:width]
		}
		rows[i+1] = row
	}
	if len(view.Rows) > 0 {
		rows[view.Index()+1] = ansiInverted(view.Selected())
	}
	return rows
}
