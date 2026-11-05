package utils

import (
	"testing"

	"github.com/projectdiscovery/pdtm/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestValidToolsDropsInvalidNames(t *testing.T) {
	tools := ValidTools([]types.Tool{{Name: "nuclei"}, {Name: "../../../other-dir/important.txt"}, {Name: "httpx"}})
	require.Len(t, tools, 2)
	require.Equal(t, "nuclei", tools[0].Name)
	require.Equal(t, "httpx", tools[1].Name)
}
