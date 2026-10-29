package pkg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/projectdiscovery/pdtm/pkg/types"
	"github.com/projectdiscovery/pdtm/pkg/version"
	"github.com/stretchr/testify/require"
)

// A tool name from the API must not reach files outside the install path.
func TestToolNameWithPathIsRejected(t *testing.T) {
	home := t.TempDir()
	binDir := filepath.Join(home, ".pdtm", "go", "bin")
	require.NoError(t, os.MkdirAll(binDir, 0o755))
	other := filepath.Join(home, "other-dir", "important.txt")
	require.NoError(t, os.MkdirAll(filepath.Dir(other), 0o755))
	require.NoError(t, os.WriteFile(other, []byte("do not delete"), 0o644))

	tool := types.Tool{Name: "../../../other-dir/important.txt", Repo: "x", Version: "9.9.9", Assets: map[string]string{"a.zip": "1"}}

	require.ErrorIs(t, Remove(binDir, tool), types.ErrInvalidToolName)
	require.ErrorIs(t, Update(binDir, tool, true), types.ErrInvalidToolName)
	require.ErrorIs(t, Install(binDir, tool), types.ErrInvalidToolName)
	require.ErrorIs(t, GoInstall(binDir, tool), types.ErrInvalidToolName)
	_, err := version.ExtractInstalledVersion(tool, binDir)
	require.ErrorIs(t, err, types.ErrInvalidToolName)

	_, err = os.Stat(other)
	require.NoError(t, err, "file outside the install path was removed")
}

func TestIsValidToolName(t *testing.T) {
	for _, name := range []string{"nuclei", "interactsh-client", "chaos-client", "tool_1", "tlsx"} {
		require.True(t, types.IsValidToolName(name), name)
	}
	for _, name := range []string{"", "..", ".", "-x", "a/b", `a\b`, "../x", "a b", "a.exe", "x\x00"} {
		require.False(t, types.IsValidToolName(name), "%q", name)
	}
}
