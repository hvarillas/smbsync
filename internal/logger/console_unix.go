//go:build !windows

package logger

func enableVirtualTerminal() bool {
	return true
}
