package pkg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// fetchReleaseBody must return the body of the release tagged with the
// supplied version, not whatever GitHub currently considers "latest".
// Regression coverage for https://github.com/projectdiscovery/pdtm/issues/435.
func TestFetchReleaseBody_PinsToVersion(t *testing.T) {
	body, err := fetchReleaseBody("dnsx", "1.1.1")
	require.NoError(t, err)
	require.NotEmpty(t, body, "release body for dnsx v1.1.1 should be non-empty")
}
