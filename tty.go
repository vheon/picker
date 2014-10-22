package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"os/exec"
)

type TTY struct {
	originalState string
	Stdin         *os.File
	Stdout        *os.File
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

func (tty *TTY) ReadByte() byte {
	buf := bufio.NewReader(tty.Stdin)
	b, err := buf.ReadByte()
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func (tty *TTY) Write(s string) {
	_, err := tty.Stdout.WriteString(s)
	if err != nil {
		log.Fatal(err)
	}
}

func (tty *TTY) Puts(s string) {
	tty.Write(s + "\n")
}
