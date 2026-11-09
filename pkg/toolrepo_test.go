package pkg

import (
	"io"
	"net/http"
	"runtime"
	"strings"
	"testing"

	"github.com/projectdiscovery/pdtm/pkg/types"
	"github.com/stretchr/testify/require"
)

type recordingTransport struct{ urls []string }

func (r *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r.urls = append(r.urls, req.URL.String())
	return &http.Response{StatusCode: http.StatusNotFound, Header: http.Header{}, Body: io.NopCloser(strings.NewReader("{}")), Request: req}, nil
}

// A repo from the API must not move the GitHub API path to another owner.
func TestToolRepoWithPathIsRejected(t *testing.T) {
	transport := &recordingTransport{}
	previous := http.DefaultTransport
	http.DefaultTransport = transport
	defer func() { http.DefaultTransport = previous }()

	osName := runtime.GOOS
	if osName == "darwin" {
		osName = "macOS"
	}
	tool := types.Tool{
		Name:    "faketool",
		Repo:    "../attacker/evil",
		Version: "1.0.0",
		Assets:  map[string]string{"faketool_1.0.0_" + osName + "_" + runtime.GOARCH + ".zip": "42"},
	}

	require.ErrorIs(t, Install(t.TempDir(), tool), types.ErrInvalidRepoName)
	require.ErrorIs(t, GoInstall(t.TempDir(), tool), types.ErrInvalidRepoName)
	require.Empty(t, transport.urls)
}
