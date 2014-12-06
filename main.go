package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"

	"golang.org/x/crypto/ssh/terminal"
)

const VisibleCandidates int = 20

const (
	keyCtrlC     = 3
	keyCtrlD     = 4
	keyCtrlN     = 14
	keyCtrlP     = 16
	keyCtrlU     = 21
	keyCtrlW     = 23
	keyEnter     = '\r'
	keyEscape    = 27
	keyUp        = 38
	keyDown      = 40
	keyBackspace = 127
)

type steps struct {
	Up    int
	Down  int
	Left  int
	Right int
}

func move(step steps) string {
	count := step.Down + step.Up + step.Left + step.Right
	movement := make([]rune, 3*(count))
	m := movement
	for i := 0; i < step.Up; i++ {
		m[0] = keyEscape
		m[1] = '['
		m[2] = 'A'
		m = m[3:]
	}
	for i := 0; i < step.Down; i++ {
		m[0] = keyEscape
		m[1] = '['
		m[2] = 'B'
		m = m[3:]
	}
	for i := 0; i < step.Left; i++ {
		m[0] = keyEscape
		m[1] = '['
		m[2] = 'D'
		m = m[3:]
	}
	for i := 0; i < step.Right; i++ {
		m[0] = keyEscape
		m[1] = '['
		m[2] = 'C'
		m = m[3:]
	}

	return string(movement)
}

var (
	SaveCursorPosition    string = string([]rune{keyEscape, '[', 's'})
	RestoreCursorPosition string = string([]rune{keyEscape, '[', 'u'})
	EraseDisplay          string = string([]rune{keyEscape, '[', 'J'})

	ReverseColor string = string([]rune{keyEscape, '[', '7', 'm'})
	ResetColor   string = string([]rune{keyEscape, '[', '0', 'm'})
)

func OpenTTY() (*os.File, error) {
	return os.OpenFile("/dev/tty", os.O_RDWR, 0)
}

func MakeRaw(tty *os.File) (*terminal.State, error) {
	return terminal.MakeRaw(int(tty.Fd()))
}

func Restore(tty *os.File, state *terminal.State) error {
	return terminal.Restore(int(tty.Fd()), state)
}

func GetSize(tty *os.File) (int, int, error) {
	return terminal.GetSize(int(tty.Fd()))
}

func TTYReverse(str string) string {
	return ReverseColor + str + ResetColor
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	tty, err := OpenTTY()
	if err != nil {
		panic(err)
	}
	originalState, err := MakeRaw(tty)
	if err != nil {
		panic(err)
	}
	defer Restore(tty, originalState)

	width, height, err := GetSize(tty)
	if err != nil {
		panic(err)
	}

	visible := min(VisibleCandidates, height)
	prompt := "> "
	picker := NewPicker(prompt, visible, width, os.Stdin)

	in := make(chan rune)
	quit := make(chan struct{})
	selection := make(chan struct{})
	back := make(chan struct{})
	clear := make(chan struct{})
	down := make(chan struct{})
	up := make(chan struct{})
	go func() {
		reader := bufio.NewReader(tty)
		for {
			r, _, err := reader.ReadRune()
			if err != nil {
				break
			}
			switch r {
			case keyEscape, keyCtrlC:
				quit <- struct{}{}
			case keyEnter:
				selection <- struct{}{}
			case keyBackspace:
				back <- struct{}{}
			case keyCtrlU, keyCtrlW:
				clear <- struct{}{}
			case keyCtrlN, keyDown:
				down <- struct{}{}
			case keyCtrlP, keyUp:
				up <- struct{}{}
			default:
				in <- r
			}
		}
		close(in)
		close(quit)
		close(selection)
		close(back)
		close(clear)
		close(down)
		close(up)
	}()

	// write the first view
	tty.WriteString(picker.View())

	// go to the start of the first line
	lastLineIndex := min(VisibleCandidates, len(picker.all))
	tty.WriteString(move(steps{
		Up:    VisibleCandidates,
		Down:  0,
		Left:  len(picker.all[lastLineIndex-1].value),
		Right: 0,
	}))

	// save the pos
	tty.WriteString(SaveCursorPosition)

	// focus on the right spot in the prompt
	tty.WriteString(move(steps{
		Up:    0,
		Down:  0,
		Left:  0,
		Right: len(picker.prompt) + len(picker.query),
	}))

	for {
		select {
		case r := <-in:
			picker.query += string(r)
			picker.index = 0
			picker.Sort()
		case <-back:
			picker.Backspace()
		case <-clear:
			picker.Clear()
		case <-quit:
			tty.WriteString(RestoreCursorPosition)
			os.Exit(1)
		case <-selection:
			tty.WriteString(RestoreCursorPosition)
			fmt.Println(picker.Selected())
			return
		case <-down:
			picker.Down()

		case <-up:
			picker.Up()
		}

		// go to the stored position
		tty.WriteString(RestoreCursorPosition)
		// clear the screen
		tty.WriteString(EraseDisplay)
		// write what we should see
		tty.WriteString(picker.View())
		// move the cursor to the right prompt position
		tty.WriteString(RestoreCursorPosition)
		tty.WriteString(move(steps{
			Up:    0,
			Down:  0,
			Left:  0,
			Right: len(picker.prompt) + len(picker.query),
		}))
	}
}
