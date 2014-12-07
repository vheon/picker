package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"

	"golang.org/x/crypto/ssh/terminal"
)

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

func move(step steps) []byte {
	count := step.Down + step.Up + step.Left + step.Right
	movement := make([]byte, 3*(count))
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

	return movement
}

var (
	SaveCursorPosition    = []byte{keyEscape, '[', 's'}
	RestoreCursorPosition = []byte{keyEscape, '[', 'u'}
	EraseDisplay          = []byte{keyEscape, '[', 'J'}
	ReverseColor          = []byte{keyEscape, '[', '7', 'm'}
	ResetColor            = []byte{keyEscape, '[', '0', 'm'}
	ShowCursor            = []byte{keyEscape, '[', '?', '2', '5', 'h'}
)

func OpenTTY() (*os.File, error) {
	return os.OpenFile("/dev/tty", os.O_RDWR, 0)
}

func TTYReverse(str string) string {
	return string(ReverseColor) + str + string(ResetColor)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var (
	VisibleCandidates = flag.Int("v", 20, "Number of visible candidates")
	vim               = flag.Bool("vim", false, "Print at bottom of screen")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	tty, err := OpenTTY()
	if err != nil {
		panic(err)
	}
	originalState, err := terminal.MakeRaw(int(tty.Fd()))
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(int(tty.Fd()), originalState)

	width, height, err := terminal.GetSize(int(tty.Fd()))
	if err != nil {
		panic(err)
	}

	visible := min(*VisibleCandidates, height)
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

	if *vim {
		// start from the bottom of the screen
		tty.Write(move(steps{
			Up:    0,
			Left:  width,
			Right: 0,
			Down:  height,
		}))
		tty.Write(ShowCursor)
	}

	// write the first view
	tty.WriteString(picker.View())

	// go to the start of the first line
	lastLineIndex := min(*VisibleCandidates, len(picker.all))
	tty.Write(move(steps{
		Up:    *VisibleCandidates,
		Down:  0,
		Left:  len(picker.all[lastLineIndex-1].value),
		Right: 0,
	}))

	// save the pos
	tty.Write(SaveCursorPosition)

	// focus on the right spot in the prompt
	tty.Write(move(steps{
		Up:    0,
		Down:  0,
		Left:  0,
		Right: len(picker.prompt) + len(picker.query),
	}))

	for {
		select {
		case r := <-in:
			picker.AppendToQuery(r)
			picker.Sort()
		case <-back:
			picker.Backspace()
			picker.Sort()
		case <-clear:
			picker.Clear()
		case <-quit:
			tty.Write(RestoreCursorPosition)
			os.Exit(1)
		case <-selection:
			tty.Write(RestoreCursorPosition)
			fmt.Println(picker.Selected())
			return
		case <-down:
			picker.Down()

		case <-up:
			picker.Up()
		}

		// go to the stored position
		tty.Write(RestoreCursorPosition)
		// clear the screen
		tty.Write(EraseDisplay)
		// write what we should see
		tty.WriteString(picker.View())
		// move the cursor to the right prompt position
		tty.Write(RestoreCursorPosition)
		tty.Write(move(steps{
			Up:    0,
			Down:  0,
			Left:  0,
			Right: len(picker.prompt) + len(picker.query),
		}))
	}
}
