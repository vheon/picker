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

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	picker := NewPicker(readAllCandidates(os.Stdin), visibleRows)

	tty := NewTTY()
	defer tty.Restore()

	terminal := NewTerminal(tty)
	terminal.ConfigTerminal()

	view := picker.Answer("")
	terminal.MakeRoom(view.Height)
	for {
		terminal.Draw(view)

		// XXX: check this
		key := tty.ReadByte()
		switch key {
		case Ctrl_N:
			view.Down()
		case Ctrl_P:
			view.Up()
		case Backspace:
			view = picker.Answer(backspace(view.Query))

		// XXX: check this! Especially how to read from tty
		case LF:
			terminal.MoveBottom()
			fmt.Println(view.Selected())
			return
		case Ctrl_U, Ctrl_W:
			view = picker.Answer("")
		default:
			view = picker.Answer(appendChar(view.Query, key))
		}
	}
}
