package path

import (
	"os"
	"path/filepath"
	"strings"
)

// IsSubPath determines if sub is a subdirectory of parent
func IsSubPath(parent, sub string) bool {
	up := ".." + string(os.PathSeparator)
	rel, err := filepath.Rel(parent, sub)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, up) && rel != ".."
}
