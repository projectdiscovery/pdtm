//go:build windows && !nounsafe
// +build windows,!nounsafe

package path

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	HWND_BROADCAST   = uintptr(0xffff)
	WM_SETTINGCHANGE = uintptr(0x001A)
)

func init() {
	// WM_SETTING_CHANGE
	// https://gist.github.com/microo8/c1b9525efab9bb462adf9d123e855c52
	sendmsg = func() {
		//x, y, err := syscall.
		_, _, err := syscall.
			NewLazyDLL("user32.dll").
			NewProc("SendMessageW").
			Call(HWND_BROADCAST, WM_SETTINGCHANGE, 0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("ENVIRONMENT"))))
		//fmt.Fprintf(os.Stderr, "%d, %d, %s\n", x, y, err)
		if nil != err {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
}
