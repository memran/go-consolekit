//go:build windows

package console

import "syscall"

func processExists(pid int) bool {
	const wantAccess = 0x0400
	handle, err := syscall.OpenProcess(wantAccess, false, uint32(pid))
	if err != nil {
		return false
	}
	syscall.CloseHandle(handle)
	return true
}
