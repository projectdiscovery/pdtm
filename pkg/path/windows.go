//go:build windows
// +build windows

// from https://github.com/therootcompany/pathman with some minor changes
package path

// Needs to
//   * use the registry editor directly to avoid possible PATH truncation
//     ( https://stackoverflow.com/questions/9546324/adding-directory-to-path-environment-variable-in-windows )
//     ( https://superuser.com/questions/387619/overcoming-the-1024-character-limit-with-setx )
//   * explicitly send WM_SETTINGCHANGE
//     ( https://github.com/golang/go/issues/18680#issuecomment-275582179 )

import (
	"fmt"
	"os"
	"strings"

	"github.com/projectdiscovery/gologger"
	"golang.org/x/sys/windows/registry"
)

func add(p string) (bool, error) {
	cur, err := paths()
	if nil != err {
		return false, err
	}

	index, err := IndexOf(cur, p)
	if err != nil {
		return false, err
	}
	// skip silently, successfully
	if index >= 0 {
		return false, nil
	}

	cur = append([]string{p}, cur...)
	err = write(cur)
	if nil != err {
		return false, err
	}
	return true, nil
}

func write(cur []string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("Can't open HKCU Environment for writes: %s", err)
	}
	defer k.Close()

	err = k.SetStringValue(`Path`, strings.Join(cur, string(os.PathListSeparator)))
	if nil != err {
		return fmt.Errorf("Can't set HKCU Environment[Path]: %s", err)
	}
	err = k.Close()
	if nil != err {
		return err
	}
	if nil != sendmsg {
		sendmsg()
	} else {
		gologger.Info().Label("WRN").Msg("added PATH, but you must reboot for changes to take effect\n")
	}
	return nil
}

func paths() ([]string, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE)
	if err != nil {
		return nil, fmt.Errorf("Can't open HKCU Environment for reads: %s", err)
	}
	defer k.Close()
	s, _, err := k.GetStringValue("Path")
	if err != nil {
		if strings.Contains(err.Error(), "cannot find the file") {
			return []string{}, nil
		}
		return nil, fmt.Errorf("Can't query HKCU Environment[Path]: %s", err)
	}
	return strings.Split(s, string(os.PathListSeparator)), nil
}
