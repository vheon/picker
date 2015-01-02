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

var (
	SaveCursorPosition     = []byte{keyEscape, '[', 's'}
	RestoreCursorPosition  = []byte{keyEscape, '[', 'u'}
	EraseDisplayFromCursor = []byte{keyEscape, '[', 'J'}
	ReverseColor           = []byte{keyEscape, '[', '7', 'm'}
	ResetColor             = []byte{keyEscape, '[', '0', 'm'}
	ShowCursor             = []byte{keyEscape, '[', '?', '2', '5', 'h'}
)

type TTY struct {
	*os.File
}

func OpenTTY() (*TTY, error) {
	file, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	tty := TTY{file}
	return &tty, err
}

func (tty *TTY) Fd() int {
	return int(tty.File.Fd())
}

type steps struct {
	Up    int
	Down  int
	Left  int
	Right int
}

func (tty *TTY) moveCursor(step steps) {
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

	tty.Write(movement)
}

func (tty *TTY) showCursor() {
	tty.Write(ShowCursor)
}

func (tty *TTY) saveCursorPosition() {
	tty.Write(SaveCursorPosition)
}

func (tty *TTY) restoreCursorPosition() {
	tty.Write(RestoreCursorPosition)
}

func (tty *TTY) eraseDisplayFromCursor() {
	tty.Write(EraseDisplayFromCursor)
}

func (tty *TTY) focusWritingPoint(picker *Picker) {
	tty.restoreCursorPosition()
	tty.moveCursor(steps{
		Right: len(picker.prompt) + len(picker.query),
	})
}

type Renderer struct {
	tty     *TTY
	width   int
	height  int
	visible int
}

func (r *Renderer) PrepareForTerminalVim() {
	// start from the bottom of the screen
	r.tty.moveCursor(steps{
		Left: r.width,
		Down: r.height,
	})
	r.tty.showCursor()
}

func (r *Renderer) renderFirstFrame(picker *Picker) {
	// write the first view
	r.tty.WriteString(picker.View())

	// going width time to the left is more than necessary but it works in all
	// situations and is simpler
	r.tty.moveCursor(steps{
		Up:   r.visible,
		Left: r.width,
	})

	// save the pos
	r.tty.saveCursorPosition()

	r.tty.focusWritingPoint(picker)
}

func (r *Renderer) Start(channel chan *Picker) {

	r.renderFirstFrame(<-channel)

	for picker := range channel {
		r.tty.restoreCursorPosition()
		r.tty.eraseDisplayFromCursor()

		// write what we should see
		r.tty.WriteString(picker.View())

		r.tty.focusWritingPoint(picker)
	}
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
	originalState, err := terminal.MakeRaw(tty.Fd())
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(tty.Fd(), originalState)

	width, height, err := terminal.GetSize(tty.Fd())
	if err != nil {
		panic(err)
	}

	if height < visible {
		visible = height
	}

	renderChan := make(chan *Picker)
	renderer := &Renderer{
		tty:     tty,
		height:  height,
		width:   width,
		visible: visible,
	}

	if vim {
		renderer.PrepareForTerminalVim()
	}
	go renderer.Start(renderChan)

	picker := NewPicker("> ", visible, width, os.Stdin, renderChan)

	input := make(chan rune)
	go func() {
		reader := bufio.NewReader(tty)
		for {
			r, _, err := reader.ReadRune()
			if err != nil || r == keyEscape || r == keyCtrlC {
				tty.restoreCursorPosition()
				os.Exit(1)
			}
			input <- r
		}
	}()

	for r := range input {
		switch r {
		case keyEnter:
			tty.restoreCursorPosition()
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
	}
}
