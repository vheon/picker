package main

// XXX: For now it works do not work on linux nor windows

import (
	"errors"
	"syscall"
	"unsafe"
)

func (tty *TTY) MakeRaw() (*syscall.Termios, error) {
	var old syscall.Termios
	if err := tty.syscall(syscall.TIOCGETA, unsafe.Pointer(&old)); err != nil {
		return nil, err
	}

	newState := old
	newState.Lflag &^= syscall.ECHO | syscall.ICANON

	if err := tty.syscall(syscall.TIOCSETA, unsafe.Pointer(&newState)); err != nil {
		return nil, err
	}
	return &old, nil
}

func (tty *TTY) Restore(state *syscall.Termios) error {
	return tty.syscall(syscall.TIOCSETA, unsafe.Pointer(state))
}

func (tty *TTY) GetSize() (int, int, error) {
	var dimensions [4]uint16
	if err := tty.syscall(syscall.TIOCGWINSZ, unsafe.Pointer(&dimensions)); err != nil {
		return -1, -1, err
	}
	return int(dimensions[1]), int(dimensions[0]), nil
}

func (tty *TTY) syscall(cmd uintptr, ptr unsafe.Pointer) error {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, tty.Fd(), cmd, uintptr(ptr))
	if err != 0 {
		return errors.New("Syscall SYS_IOCTL error")
	}
	return nil
}
