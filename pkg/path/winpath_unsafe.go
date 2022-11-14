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
	sendmsg = func() error {
		utf16PtrENV, err := syscall.UTF16PtrFromString("ENVIRONMENT")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return err
		}
		_, _, err = syscall.
			NewLazyDLL("user32.dll").
			NewProc("SendMessageW").
			Call(HWND_BROADCAST, WM_SETTINGCHANGE, 0, uintptr(unsafe.Pointer(utf16PtrENV)))
		if nil != err {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		return nil
	}
}
