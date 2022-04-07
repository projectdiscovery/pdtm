package pkg

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/go-github/github"
)

var ErrIsInstalled = errors.New("already present")

// Install installs given tool at path
func Install(tool Tool, path string) (string, error) {
	if path, _ := exec.LookPath(tool.Name); path != "" {
		return "", ErrIsInstalled
	}
	githubClient := GithubClient()

	releases, _, err := githubClient.Repositories.ListReleases(context.Background(), Organization, tool.Repo, &github.ListOptions{
		PerPage: 1,
	})
	if err != nil {
		return "", err
	}
	if len(releases) == 0 {
		return "", errors.New("could not get latest release")

	}

	assets, _, err := githubClient.Repositories.ListReleaseAssets(context.Background(), Organization, tool.Repo, *releases[0].ID, nil)
	if err != nil {
		return "", err
	}

	builder := &strings.Builder{}
	builder.WriteString(tool.Name)
	builder.WriteString("_")
	builder.WriteString(strings.TrimPrefix(releases[0].GetTagName(), "v"))
	builder.WriteString("_")
	if runtime.GOOS == "darwin" {
		builder.WriteString("macOS")
	} else {
		builder.WriteString(runtime.GOOS)
	}
	builder.WriteString("_")
	builder.WriteString(runtime.GOARCH)
	builder.WriteString(".zip")
	var id int64
	for _, asset := range assets {
		if *asset.Name == builder.String() {
			id = *asset.ID
			break
		}
	}
	builder.Reset()

	_, rdurl, err := githubClient.Repositories.DownloadReleaseAsset(context.Background(), Organization, tool.Repo, id)
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
	return releases[0].GetTagName(), nil
}
