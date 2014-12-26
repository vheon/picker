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
	visible int
	vim     bool
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.IntVar(&visible, "v", 20, "Number of visible candidates")
	flag.BoolVar(&vim, "vim", false, "Print at bottom of screen")
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

	if height < visible {
		visible = height
	}

	picker := NewPicker("> ", visible, width, os.Stdin)

	if vim {
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

	// going width time to the left is more than necessary but it works in all
	// situations and is simpler
	tty.Write(move(steps{
		Up:    visible,
		Down:  0,
		Left:  width,
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

	input := make(chan rune)
	go func() {
		reader := bufio.NewReader(tty)
		for {
			r, _, err := reader.ReadRune()
			if err != nil || r == keyEscape || r == keyCtrlC {
				tty.Write(RestoreCursorPosition)
				os.Exit(1)
			}
			input <- r
		}
	}()

	for r := range input {
		switch r {
		case keyEnter:
			tty.Write(RestoreCursorPosition)
			fmt.Println(picker.Selected())
			return
		case keyCtrlU, keyCtrlW:
			picker.Clear()
		case keyCtrlN, keyDown:
			picker.Down()
		case keyCtrlP, keyUp:
			picker.Up()
		case keyBackspace:
			picker.Back()
		default:
			picker.More(r)
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
