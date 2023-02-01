package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetToolStruct() Tool {
	tool := Tool{
		Name:    "dnsx",
		Repo:    "dnsx",
		Version: "1.1.1",
		Assets: map[string]string{
			"dnsx_1.1.1_checksums.txt":     "79344865",
			"dnsx_1.1.1_linux_386.zip":     "79344862",
			"dnsx_1.1.1_linux_amd64.zip":   "79344859",
			"dnsx_1.1.1_linux_arm64.zip":   "79344852",
			"dnsx_1.1.1_linux_armv6.zip":   "79344864",
			"dnsx_1.1.1_macOS_amd64.zip":   "79344851",
			"dnsx_1.1.1_macOS_arm64.zip":   "79344856",
			"dnsx_1.1.1_windows_386.zip":   "79344855",
			"dnsx_1.1.1_windows_amd64.zip": "79344857",
		},
	}

	return tool
}

func TestInstallTool(t *testing.T) {
	tool := GetToolStruct()

	pathBin, err := os.MkdirTemp("", "test-dir")
	assert.Nil(t, err, err)
	defer os.RemoveAll(pathBin)

	// create directory
	err = os.MkdirAll(pathBin, 0777)
	assert.Nil(t, err)

	// install first time
	err = Install(pathBin, tool)
	assert.Nil(t, err, err)

	// check if its installed in path
	// need to throw exeption
	err = Install(pathBin, tool)
	assert.NotNil(t, err, err)
}

func TestRemoveTool(t *testing.T) {
	tool := GetToolStruct()

	pathBin, err := os.MkdirTemp("", "test-dir")
	assert.Nil(t, err)
	defer os.RemoveAll(pathBin)

	// install first time
	err = Install(pathBin, tool)
	assert.Nil(t, err)

	// remove when tool exist in path
	err = Remove(pathBin, tool)
	assert.Nil(t, err)

	// throws error
	err = Remove(pathBin, tool)
	assert.Equal(t, err, fmt.Errorf(ErrToolNotFound, tool.Name, filepath.Join(pathBin, tool.Name)))
}

func TestUpdateToolUpToDate(t *testing.T) {
	tool := GetToolStruct()

	pathBin, err := os.MkdirTemp("", "test-dir")
	assert.Nil(t, err)
	defer os.RemoveAll(pathBin)

	// install first time
	err = Install(pathBin, tool)
	assert.Nil(t, err)

	// remove when tool exist in path
	err = Update(pathBin, tool)
	assert.Equal(t, "already up to date", err.Error())
}

func TestUpdateToolDufferentVersion(t *testing.T) {
	tool := GetToolStruct()

	pathBin, err := os.MkdirTemp("", "test-dir")
	assert.Nil(t, err)
	defer os.RemoveAll(pathBin)

	// install first time
	err = Install(pathBin, tool)
	assert.Nil(t, err)

	// remove when tool exist in path
	err = Update(pathBin, tool)
	assert.Equal(t, "already up to date", err.Error())

	// remove tool
	err = Remove(pathBin, tool)
	assert.Nil(t, err)

	// update tool removed
	// will install new one
	err = Update(pathBin, tool)
	assert.Nil(t, err)

	// check if tool is installed post update
	_, err = os.Stat(filepath.Join(pathBin, tool.Name))
	assert.Nil(t, err, err)
}
