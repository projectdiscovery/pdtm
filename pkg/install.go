package pkg

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/logrusorgru/aurora/v4"
	"github.com/projectdiscovery/gologger"
	ospath "github.com/projectdiscovery/pdtm/pkg/path"
	"github.com/projectdiscovery/pdtm/pkg/types"
)

var (
	extIfFound = ".exe"
	au         = aurora.New(aurora.WithColors(true))
)

// Install installs given tool at path
func Install(path string, tool types.Tool) error {
	if _, exists := ospath.GetExecutablePath(path, tool.Name); exists {
		return types.ErrIsInstalled
	}
	gologger.Info().Msgf("installing %s...", tool.Name)
	version, err := install(tool, path)
	if err != nil {
		return err
	}
	gologger.Info().Msgf("installed %s %s (%s)", tool.Name, version, au.BrightGreen("latest").String())
	return nil
}

func install(tool types.Tool, path string) (string, error) {
	builder := &strings.Builder{}
	builder.WriteString(tool.Name)
	builder.WriteString("_")
	builder.WriteString(strings.TrimPrefix(tool.Version, "v"))
	builder.WriteString("_")
	if strings.EqualFold(runtime.GOOS, "darwin") {
		builder.WriteString("macOS")
	} else {
		builder.WriteString(runtime.GOOS)
	}
	builder.WriteString("_")
	builder.WriteString(runtime.GOARCH)
	var id int
	var isZip, isTar bool
loop:
	for asset, assetID := range tool.Assets {
		switch {
		case strings.Contains(asset, ".zip"):
			if strings.EqualFold(asset, builder.String()+".zip") {
				id, _ = strconv.Atoi(assetID)
				isZip = true
				break loop
			}
		case strings.Contains(asset, ".tar.gz"):
			if strings.EqualFold(asset, builder.String()+".tar.gz") {
				id, _ = strconv.Atoi(assetID)
				isTar = true
				break loop
			}
		}
	}
	builder.Reset()

	// handle if id is zero (no asset found)
	if id == 0 {
		return "", fmt.Errorf(types.ErrNoAssetFound, runtime.GOOS, runtime.GOARCH)
	}

	_, rdurl, err := GithubClient().Repositories.DownloadReleaseAsset(context.Background(), types.Organization, tool.Repo, int64(id))
	if err != nil {
		if arlErr, ok := err.(*github.AbuseRateLimitError); ok {
			// Provide user with more info regarding the rate limit
			gologger.Error().Msgf("error for remaining request per hour: %s, RetryAfter: %s", err.Error(), arlErr.RetryAfter)
		}
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

	switch {
	case isZip:
		err := downloadZip(resp.Body, tool.Name, path)
		if err != nil {
			return "", err
		}
	case isTar:
		err := downloadTar(resp.Body, tool.Name, path)
		if err != nil {
			return "", err
		}
	}
	return tool.Version, nil
}

func downloadTar(reader io.Reader, toolName, path string) error {
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	tarReader := tar.NewReader(gzipReader)
	// iterate through the files in the archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if !strings.EqualFold(strings.TrimSuffix(header.FileInfo().Name(), extIfFound), toolName) {
			continue
		}
		// if the file is not a directory, extract it
		if !header.FileInfo().IsDir() {
			filePath := filepath.Join(path, header.FileInfo().Name())
			if !strings.HasPrefix(filePath, filepath.Clean(path)+string(os.PathSeparator)) {
				return err
			}

			if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				return err
			}

			dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, header.FileInfo().Mode())
			if err != nil {
				return err
			}
			defer dstFile.Close()
			// copy the file data from the archive
			_, err = io.Copy(dstFile, tarReader)
			if err != nil {
				return err
			}
			// set the file permissions
			err = os.Chmod(dstFile.Name(), 0755)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func downloadZip(reader io.Reader, toolName, path string) error {
	buff := bytes.NewBuffer([]byte{})
	size, err := io.Copy(buff, reader)
	if err != nil {
		return err
	}
	zipReader, err := zip.NewReader(bytes.NewReader(buff.Bytes()), size)
	if err != nil {
		return err
	}
	for _, f := range zipReader.File {
		if !strings.EqualFold(strings.TrimSuffix(f.Name, extIfFound), toolName) {
			continue
		}
		filePath := filepath.Join(path, f.Name)
		if !strings.HasPrefix(filePath, filepath.Clean(path)+string(os.PathSeparator)) {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}
		err = os.Chmod(dstFile.Name(), 0755)
		if err != nil {
			return err
		}

		dstFile.Close()
		fileInArchive.Close()
	}
	return nil
}
