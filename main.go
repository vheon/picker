package main

import (
	"bufio"
	"flag"
	"os"
	"runtime"
	"unicode/utf8"

	"golang.org/x/crypto/ssh/terminal"
)

const VisibleCandidates int = 20

const (
	keyCtrlC     = 3
	keyCtrlD     = 4
	keyCtrlU     = 21
	keyCtrlW     = 23
	keyEnter     = '\r'
	keyEscape    = 27
	keyBackspace = 127
	// keyUp
	// keyDown
	// keyLeft
	// keyRight
	// keyHome
	// keyEnd
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
	flag.Parse()

	tty, err := OpenTTY()
	if err != nil {
		panic(err)
	}
	width, height, err := GetSize(tty)
	if err != nil {
		panic(err)
	}
	originalState, err := MakeRaw(tty)
	if err != nil {
		panic(err)
	}
	defer Restore(tty, originalState)

	visible := min(VisibleCandidates, height)
	prompt := "> "
	picker := NewPicker(prompt, visible, width, os.Stdin)

	in := make(chan rune)
	quit := make(chan struct{})
	end := make(chan struct{})
	back := make(chan struct{})
	clear := make(chan struct{})
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
				end <- struct{}{}
			case keyBackspace:
				back <- struct{}{}
			case keyCtrlU, keyCtrlW:
				clear <- struct{}{}
			default:
				in <- r
			}
		}
		close(quit)
		close(end)
		close(in)
	}()

	tty.WriteString(picker.String())

	// go to the start of the first line
	tty.WriteString(move(movements{
		Up:    VisibleCandidates,
		Down:  0,
		Left:  len(picker.view.lines[VisibleCandidates-1]),
		Right: 0,
	}))

	// save the pos
	tty.WriteString(SaveCursorPosition)

	// move the cursor to the right prompt position
	tty.WriteString(move(movements{
		Up:    0,
		Down:  0,
		Left:  0,
		Right: len(picker.prompt),
	}))

	for {
		select {
		case r := <-in:
			picker.query += string(r)
		case <-back:
			_, size := utf8.DecodeLastRuneInString(picker.query)
			picker.query = picker.query[:len(picker.query)-size]
			picker.valid = len(picker.all)
		case <-clear:
			picker.query = ""
			picker.valid = len(picker.all)
		case <-quit:
			os.Exit(1)
		case <-end:
			tty.WriteString(move(movements{
				Up:    0,
				Down:  0,
				Right: 0,
				Left:  len(picker.prompt) + len(picker.query),
			}))
			tty.WriteString(picker.view.Selected() + "\n")
			return
		}

		// do the reorder
		picker.Sort()

		// go to the stored position
		tty.WriteString(RestoreCursorPosition)

		// clear the screen
		tty.WriteString(EraseDisplay)

		// write what we should see
		tty.WriteString(picker.String())

		// move the cursor to the right prompt position
		tty.WriteString(RestoreCursorPosition)
		tty.WriteString(move(movements{
			Up:    0,
			Down:  0,
			Left:  0,
			Right: len(picker.prompt) + len(picker.query),
		}))
	}
}
