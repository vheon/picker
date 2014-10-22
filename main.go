package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
)

const (
	Ctrl_N    byte = 14
	Ctrl_P         = 16
	Ctrl_U         = 21
	Ctrl_W         = 23
	Backspace      = 127
	LF             = 10
)

const visibleRows = 20

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	scanner := bufio.NewScanner(os.Stdin)
	var lines []Candidate
	for scanner.Scan() {
		lines = append(lines, NewCandidate(scanner.Text()))
	}

	picker := NewPicker(lines, visibleRows)

	tty := NewTTY()
	defer tty.Restore()

	screen := NewScreen(tty)
	screen.ConfigScreen()

	// XXX: can we grab this from the view?
	query := ""
	view := picker.Answer(query)
	screen.MakeRoom(view.Height)
	for {
		screen.Draw(view)

		// XXX: check this
		key := tty.ReadByte()
		switch key {
		case Ctrl_N:
			view.Down()
		case Ctrl_P:
			view.Up()
		case Backspace:
			if len(query) > 0 {
				query = query[:len(query)-1]
			}
			view = picker.Answer(query)

		// XXX: check this! Especially how to read from tty
		case LF:
			fmt.Println(view.Selected())
			return
		case Ctrl_U, Ctrl_W:
			query = ""
			view = picker.Answer(query)
		default:
			query = query + string(key)
			view = picker.Answer(query)
		}
	}
}
