//go:build windows

// path from https://github.com/therootcompany/pathman with some minor changes
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
	sliceutil "github.com/projectdiscovery/utils/slice"
	"golang.org/x/sys/windows/registry"
)

func add(p string) (bool, error) {
	cur, err := getPathsFromRegistry()
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

	cur = append(cur, p)
	err = write(p, cur)
	if nil != err {
		return false, err
	}
	return true, nil
}

func remove(p string) (bool, error) {
	cur, err := getPathsFromRegistry()
	if nil != err {
		return false, err
	}

	index, err := IndexOf(cur, p)
	if err != nil {
		return false, err
	}
	// skip silently, successfully
	if index == -1 {
		return false, nil
	}

	cur = sliceutil.PruneEqual(cur, p)

	err = write(p, cur)
	if nil != err {
		return false, err
	}
	return true, nil
}

func write(path string, cur []string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("can't open HKCU Environment for writes: %s", err)
	}
	defer k.Close()

	err = k.SetStringValue(`Path`, strings.Join(cur, string(os.PathListSeparator)))
	if nil != err {
		return fmt.Errorf("can't set HKCU Environment[Path]: %s", err)
	}
	err = k.Close()
	if nil != err {
		return err
	}
	if nil != sendmsg {
		err := sendmsg()
		if err != nil {
			gologger.Info().Label("WRN").Msgf("Please reboot to load newly added $PATH (%s)", path)
		}
	} else {
		gologger.Info().Label("WRN").Msgf("Please reboot to load newly added $PATH (%s)", path)
	}
	gologger.Info().Label("WRN").Msgf("Please reload terminal to load newly added $PATH (%s)", path)
	return nil
}

func paths() []string {
	configuredPaths, _ := getPathsFromRegistry()
	return configuredPaths
}

func getPathsFromRegistry() ([]string, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE)
	if err != nil {
		return nil, fmt.Errorf("can't open HKCU Environment for reads: %s", err)
	}
	defer k.Close()
	s, _, err := k.GetStringValue("Path")
	if err != nil {
		if strings.Contains(err.Error(), "cannot find the file") {
			return []string{}, nil
		}
		return nil, fmt.Errorf("can't query HKCU Environment[Path]: %s", err)
	}
	return strings.Split(s, string(os.PathListSeparator)), nil
}

func isSet(path string) (bool, error) {
	cur, err := getPathsFromRegistry()
	if nil != err {
		return false, err
	}

	index, err := IndexOf(cur, path)
	if err != nil {
		return false, err
	}
	return index >= 0, nil
}
