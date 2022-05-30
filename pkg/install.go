package pkg

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/projectdiscovery/gologger"
)

// Install installs given tool at path
func Install(tool Tool, path string) error {
	if path, _ := exec.LookPath(tool.Name); path != "" {
		return ErrIsInstalled
	}
	gologger.Info().Msgf("installing %s...", tool.Name)
	version, err := install(tool, path)
	if err != nil {
		return err
	}
	gologger.Info().Msgf("installed %s %s(latest)", tool.Name, version)
	return nil
}

func install(tool Tool, path string) (string, error) {
	builder := &strings.Builder{}
	builder.WriteString(tool.Name)
	builder.WriteString("_")
	builder.WriteString(strings.TrimPrefix(tool.Version, "v"))
	builder.WriteString("_")
	if runtime.GOOS == "darwin" {
		builder.WriteString("macOS")
	} else {
		builder.WriteString(runtime.GOOS)
	}
	builder.WriteString("_")
	builder.WriteString(runtime.GOARCH)
	builder.WriteString(".zip")
	var id int
	for asset, assetID := range tool.Assets {
		if asset == builder.String() {
			id, _ = strconv.Atoi(assetID)
			break
		}
	}
	builder.Reset()

	_, rdurl, err := GithubClient().Repositories.DownloadReleaseAsset(context.Background(), Organization, tool.Repo, int64(id))
	if err != nil {
		return "", err
	}

	resp, err := http.Get(rdurl)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return "", err
	}

	for _, f := range zipReader.File {
		if f.Name != tool.Name {
			continue
		}
		filePath := filepath.Join(path, f.Name)
		if !strings.HasPrefix(filePath, filepath.Clean(path)+string(os.PathSeparator)) {
			return "", err
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return "", err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return "", err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return "", err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return "", err
		}
		err = os.Chmod(dstFile.Name(), 0755)
		if err != nil {
			return "", err
		}

		dstFile.Close()
		fileInArchive.Close()
	}
	return tool.Version, nil
}
