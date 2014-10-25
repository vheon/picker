package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
)

const visibleRows = 20

func appendChar(s string, b byte) string {
	return s + string(b)
}

func backspace(s string) string {
	if l := len(s); l > 0 {
		return s[:l-1]
	}
	return s
}

func readAllCandidates(r io.Reader) []Candidate {
	scanner := bufio.NewScanner(r)
	var lines []Candidate
	for scanner.Scan() {
		lines = append(lines, NewCandidate(scanner.Text()))
	}
	return lines
}

func handle_input(picker *Picker, view *View, key byte) *View {
	switch key {
	case Ctrl_N:
		view.Down()
	case Ctrl_P:
		view.Up()
	case Backspace:
		view = picker.Answer(backspace(view.Query))
	case LF:
		view.Done = true
	case Ctrl_U, Ctrl_W:
		view = picker.Answer("")
	default:
		view = picker.Answer(appendChar(view.Query, key))
	}
	return view
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	picker := NewPicker(readAllCandidates(os.Stdin), visibleRows)

	tty := NewTTY()
	defer tty.Restore()

	terminal := NewTerminal(tty)
	terminal.ConfigTerminal()

	view := picker.Answer("")
	terminal.MakeRoom(view.Height)
	for !view.Done {
		terminal.Draw(view)
		view = handle_input(picker, view, tty.ReadByte())
	}

	terminal.MoveBottom()
	fmt.Println(view.Selected())
}
