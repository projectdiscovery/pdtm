package pkg

import (
	"os"
	"testing"

	ospath "github.com/projectdiscovery/pdtm/pkg/path"
	"github.com/projectdiscovery/pdtm/pkg/types"
	"github.com/stretchr/testify/require"
)

func GetToolStruct() types.Tool {
	tool := types.Tool{
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

func TestInstall(t *testing.T) {
	tool := GetToolStruct()

	pathBin, err := os.MkdirTemp("", "test-dir")
	require.Nil(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(pathBin))
	}()

	// install first time
	err = Install(pathBin, tool)
	require.Nil(t, err)

	// installing again should trigger an error
	err = Install(pathBin, tool)
	require.NotNil(t, err)
}

func TestRemove(t *testing.T) {
	tool := GetToolStruct()

	pathBin, err := os.MkdirTemp("", "test-dir")
	require.Nil(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(pathBin))
	}()

	// install the tool
	err = Install(pathBin, tool)
	require.Nil(t, err)

	// remove it from path
	err = Remove(pathBin, tool)
	require.Nil(t, err)

	// removing non existing tool triggers an error
	err = Remove(pathBin, tool)
	require.NotNil(t, err)
}

func TestUpdateSameVersion(t *testing.T) {
	tool := GetToolStruct()

	pathBin, err := os.MkdirTemp("", "test-dir")
	require.Nil(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(pathBin))
	}()

	// install the tool
	err = Install(pathBin, tool)
	require.Nil(t, err)

	// updating a tool to the same version should trigger an error
	err = Update(pathBin, tool, true)
	require.Equal(t, "already up to date", err.Error())
}

func TestUpdateNonExistingTool(t *testing.T) {
	tool := GetToolStruct()

	pathBin, err := os.MkdirTemp("", "test-dir")
	require.Nil(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(pathBin))
	}()

	// updating non existing tool should error
	err = Update(pathBin, tool, true)
	require.NotNil(t, err)
}

func TestUpdateToolWithoutAssets(t *testing.T) {
	tool := GetToolStruct()

	pathBin, err := os.MkdirTemp("", "test-dir")
	require.Nil(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(pathBin))
	}()

	// install the tool
	err = Install(pathBin, tool)
	require.Nil(t, err)

	// remove assets
	tool.Assets = nil

	// updating a tool without assets should trigger an error
	err = Update(pathBin, tool, true)
	require.NotNil(t, err)
	// and leave the original binary in place
	_, exists := ospath.GetExecutablePath(pathBin, tool.Name)
	require.True(t, exists)
}
