package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type TTY struct {
	originalState string
	Stdin         *os.File
	Stdout        *os.File
	Height        int
	Width         int
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

	height, width := parseSize(tty.Stty("size"))
	tty.Height = height
	tty.Width = width

	return tty
}

func parseSize(size string) (int, int) {
	ssize := strings.Fields(size)
	height, _ := strconv.Atoi(ssize[0])
	width, _ := strconv.Atoi(ssize[1])
	return height, width
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

func (tty *TTY) Readc() byte {
	b := make([]byte, 1)
	_, err := tty.Stdin.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return b[0]
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
