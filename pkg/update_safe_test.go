package pkg

import (
	"archive/zip"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/projectdiscovery/pdtm/pkg/types"
	"github.com/stretchr/testify/require"
)

// rewriteRT sends every request (GitHub API, asset host) to a local server.
type rewriteRT struct{ target *url.URL }

func (r rewriteRT) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := req.Clone(req.Context())
	req2.URL.Scheme = r.target.Scheme
	req2.URL.Host = r.target.Host
	req2.Host = r.target.Host
	return http.DefaultTransport.RoundTrip(req2)
}

// fakeRelease serves asset 42 of projectdiscovery/faketool with the given
// status and body, and routes the default client to it.
func fakeRelease(t *testing.T, status int, body []byte) {
	t.Helper()
	mux := http.NewServeMux()
	var srv *httptest.Server
	mux.HandleFunc("/repos/projectdiscovery/faketool/releases/assets/42", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, srv.URL+"/dl/faketool.zip", http.StatusFound)
	})
	mux.HandleFunc("/dl/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write(body)
	})
	srv = httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	u, err := url.Parse(srv.URL)
	require.NoError(t, err)
	http.DefaultClient.Transport = rewriteRT{u}
	t.Cleanup(func() { http.DefaultClient.Transport = nil })
}

func zipWith(t *testing.T, name string, content []byte) []byte {
	t.Helper()
	var b bytes.Buffer
	zw := zip.NewWriter(&b)
	w, err := zw.Create(name)
	require.NoError(t, err)
	_, err = w.Write(content)
	require.NoError(t, err)
	require.NoError(t, zw.Close())
	return b.Bytes()
}

func fakeTool(assets map[string]string) types.Tool {
	return types.Tool{Name: "faketool", Repo: "faketool", Version: "2.0.0", Assets: assets}
}

func platformAsset() string {
	osn := runtime.GOOS
	if osn == "darwin" {
		osn = "macOS"
	}
	return fmt.Sprintf("faketool_2.0.0_%s_%s.zip", osn, runtime.GOARCH)
}

const oldBinary = "#!/bin/sh\necho faketool v1.0.0\n"

func installedTool(t *testing.T) (string, string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("fake tool is a shell script")
	}
	dir := t.TempDir()
	bin := filepath.Join(dir, "faketool")
	require.NoError(t, os.WriteFile(bin, []byte(oldBinary), 0o755))
	return dir, bin
}

func requireUnchanged(t *testing.T, dir, bin string) {
	t.Helper()
	data, err := os.ReadFile(bin)
	require.NoError(t, err, "installed binary was removed")
	require.Equal(t, oldBinary, string(data))
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	require.Len(t, entries, 1, "temporary files left in the install path")
}

func TestUpdateKeepsBinaryOnFailedDownload(t *testing.T) {
	dir, bin := installedTool(t)
	fakeRelease(t, http.StatusServiceUnavailable, []byte("upstream busy"))
	require.Error(t, Update(dir, fakeTool(map[string]string{platformAsset(): "42"}), true))
	requireUnchanged(t, dir, bin)
}

func TestUpdateKeepsBinaryWithoutPlatformAsset(t *testing.T) {
	dir, bin := installedTool(t)
	require.Error(t, Update(dir, fakeTool(map[string]string{"faketool_2.0.0_plan9_mips.zip": "42"}), true))
	requireUnchanged(t, dir, bin)
}

func TestUpdateKeepsBinaryWhenArchiveHasNoBinary(t *testing.T) {
	dir, bin := installedTool(t)
	fakeRelease(t, http.StatusOK, zipWith(t, "README.md", []byte("no binary")))
	require.Error(t, Update(dir, fakeTool(map[string]string{platformAsset(): "42"}), true))
	requireUnchanged(t, dir, bin)
}

func TestUpdateReplacesBinary(t *testing.T) {
	dir, bin := installedTool(t)
	newBinary := "#!/bin/sh\necho faketool v2.0.0\n"
	fakeRelease(t, http.StatusOK, zipWith(t, "faketool", []byte(newBinary)))
	require.NoError(t, Update(dir, fakeTool(map[string]string{platformAsset(): "42"}), true))
	data, err := os.ReadFile(bin)
	require.NoError(t, err)
	require.Equal(t, newBinary, string(data))
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	require.Len(t, entries, 1, "temporary files left in the install path")
}
