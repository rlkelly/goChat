package main

import (
  "fmt"
  "os"
  "runtime"
  "syscall"
  "unsafe"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getWinsize() (*winsize, error) {
	ws := new(winsize)

	var _TIOCGWINSZ int64

	switch runtime.GOOS {
	case "linux":
		_TIOCGWINSZ = 0x5413
	case "darwin":
		_TIOCGWINSZ = 1074295912
	}

	r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(_TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)),
	)

	if int(r1) == -1 {
		fmt.Println("Error:", os.NewSyscallError("GetWinsize", errno))
		return nil, os.NewSyscallError("GetWinsize", errno)
	}
	return ws, nil
}

func Height() int {
	ws, err := getWinsize()
	if err != nil {
		return -1
	}
	return int(ws.Row)
}

func Width() int {
	ws, err := getWinsize()

	if err != nil {
		return -1
	}

	return int(ws.Col)
}
