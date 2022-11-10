// Package winpath is useful for managing PATH as part of the Environment
// in the Windows HKey Local User registry. It returns an error for most
// operations on non-Windows systems.
package path

import (
	"os"
	"path/filepath"
	"strings"
)

// sendmsg uses a syscall to broadcast the registry change so that
// new shells will get the new PATH immediately, without a reboot
var sendmsg func()

// NormalizePathEntry will return the given directory path relative
// from its absolute path to the %USERPROFILE% (home) directory.
func NormalizePathEntry(pathentry string) (string, string, error) {
	home, err := os.UserHomeDir()
	if nil != err {
		return "", "", err
	}

	sep := string(os.PathSeparator)
	absentry, _ := filepath.Abs(pathentry)
	home, _ = filepath.Abs(home)

	var homeentry string
	if strings.HasPrefix(strings.ToLower(absentry)+sep, strings.ToLower(home)+sep) {
		// %USERPROFILE% is allowed, but only for user PATH
		// https://superuser.com/a/442163/73857
		homeentry = `%USERPROFILE%` + absentry[len(home):]
	}

	if absentry == pathentry {
		absentry = ""
	}
	if homeentry == pathentry {
		homeentry = ""
	}
	return absentry, homeentry, nil
}

// IndexOf searches the given path list for first occurence
// of the given path entry and returns the index, or -1
func IndexOf(paths []string, p string) (int, error) {
	index := -1

	abspath, homepath, err := NormalizePathEntry(p)
	if err != nil {
		return index, nil
	}
	for i, path := range paths {
		if path == "" {
			continue
		}
		if strings.ToLower(p) == strings.ToLower(path) {
			index = i
			break
		}
		if strings.ToLower(abspath) == strings.ToLower(path) {
			index = i
			break
		}
		if strings.ToLower(homepath) == strings.ToLower(path) {
			index = i
			break
		}
	}
	return index, nil
}
