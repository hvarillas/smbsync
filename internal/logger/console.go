//go:build windows

package logger

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

func enableVirtualTerminal() bool {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")

	handle := windows.Handle(uintptr(1)) // stdout
	var mode uint32
	ret, _, _ := getConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&mode)))
	if ret == 0 {
		return false
	}

	mode |= 0x0004 // ENABLE_VIRTUAL_TERMINAL_PROCESSING
	ret, _, _ = setConsoleMode.Call(uintptr(handle), uintptr(mode))
	return ret != 0
}
