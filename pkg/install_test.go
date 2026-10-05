package pkg

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/projectdiscovery/pdtm/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestGoInstallUsesRepository(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("fake-go shell scripts are not executable on windows")
	}

	goDir := t.TempDir()
	script := "#!/bin/sh\nprintf '%s\\n' \"$@\" > \"$GOBIN/go-install-args\"\n"
	require.NoError(t, os.WriteFile(filepath.Join(goDir, "go"), []byte(script), 0o755))
	t.Setenv("PATH", goDir)

	tests := []struct {
		name string
		tool types.Tool
		want string
	}{
		{
			name: "interactsh server",
			tool: types.Tool{Name: "interactsh-server", Repo: "interactsh", GoInstallPath: "cmd/interactsh-server@latest"},
			want: "github.com/projectdiscovery/interactsh/cmd/interactsh-server@latest",
		},
		{
			name: "interactsh client",
			tool: types.Tool{Name: "interactsh-client", Repo: "interactsh", GoInstallPath: "cmd/interactsh-client@latest"},
			want: "github.com/projectdiscovery/interactsh/cmd/interactsh-client@latest",
		},
		{
			name: "matching tool and repository names",
			tool: types.Tool{Name: "dnsx", Repo: "dnsx", GoInstallPath: "cmd/dnsx@latest"},
			want: "github.com/projectdiscovery/dnsx/cmd/dnsx@latest",
		},
		{
			name: "versioned module path",
			tool: types.Tool{Name: "nuclei", Repo: "nuclei", GoInstallPath: "v3/cmd/nuclei@latest"},
			want: "github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest",
		},
		{
			name: "pinned version",
			tool: types.Tool{Name: "interactsh-client", Repo: "interactsh", GoInstallPath: "cmd/interactsh-client@v1.2.4"},
			want: "github.com/projectdiscovery/interactsh/cmd/interactsh-client@v1.2.4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			binDir := t.TempDir()
			require.NoError(t, GoInstall(binDir, tt.tool))

			args, err := os.ReadFile(filepath.Join(binDir, "go-install-args"))
			require.NoError(t, err)
			require.Equal(t, "install\n-v\n"+tt.want+"\n", string(args))
		})
	}
}
