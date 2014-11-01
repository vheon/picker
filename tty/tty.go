package tty

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"os/exec"
)

type TTY struct {
	Stdin  *os.File
	Stdout *os.File

	originalState string
	reader        *bufio.Reader
}

func NewTTY() *TTY {
	stdin, err := os.Open("/dev/tty")
	if err != nil {
		log.Fatal(err)
	}
	stdout, err := os.OpenFile("/dev/tty", os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	tty := &TTY{
		Stdin:  stdin,
		Stdout: stdout,
	}

	tty.originalState = tty.Stty("-g")
	tty.reader = bufio.NewReader(stdin)

	return tty
}

func (tty *TTY) Restore() {
	tty.Stty(tty.originalState)
}

func (tty *TTY) Stty(params ...string) string {
	return tty.command("stty", params...)
}

func (tty *TTY) command(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	cmd.Stdin = tty.Stdin
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return bytes.NewBuffer(out).String()
}

func (tty *TTY) ReadRune() rune {
	r, _, err := tty.reader.ReadRune()
	if err != nil {
		log.Fatal(err)
	}
	return r
}

func (tty *TTY) Write(s string) {
	_, err := tty.Stdout.WriteString(s)
	if err != nil {
		log.Fatal(err)
	}
}
