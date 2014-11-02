package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"unicode/utf8"

	"github.com/vheon/picker/terminal"
	"github.com/vheon/picker/tty"
)

const visibleRows = 20

func appendChar(s string, b rune) string {
	return s + string(b)
}

func backspace(s string) string {
	if l := len(s); l > 0 {
		_, size := utf8.DecodeLastRuneInString(s)
		return s[:l-size]
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

func handle_input(picker *Picker, view *View, key rune) *View {
	switch key {
	case terminal.Ctrl_N:
		view.Down()
	case terminal.Ctrl_P:
		view.Up()
	case terminal.Backspace:
		view = picker.Answer(backspace(view.Query))
	case terminal.LF:
		view.Done = true
	case terminal.Ctrl_U, terminal.Ctrl_W:
		view = picker.Answer("")
	default:
		view = picker.Answer(appendChar(view.Query, key))
	}
	return view
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	picker := NewPicker(readAllCandidates(os.Stdin), visibleRows)

	tty := tty.NewTTY()
	defer tty.Restore()

	terminal := terminal.NewTerminal(tty)
	terminal.ConfigTerminal()

	view := picker.Answer("")
	terminal.MakeRoom(view.Height)
	for !view.Done {
		view.DrawOnTerminal(terminal)
		view = handle_input(picker, view, tty.ReadRune())
	}
	// should clean first?
	fmt.Println(view.Selected())
}
