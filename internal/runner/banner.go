package runner

import (
	"github.com/projectdiscovery/gologger"
	updateutils "github.com/projectdiscovery/utils/update"
)

// version is the release version, injected at build time via
// -X github.com/projectdiscovery/pdtm/internal/runner.version=<tag>.
// The literal below is only a fallback for plain `go build`.
var version = "v0.1.5"

var banner = `
                ____          
     ____  ____/ / /_____ ___ 
    / __ \/ __  / __/ __ __  \
   / /_/ / /_/ / /_/ / / / / /
  / .___/\__,_/\__/_/ /_/ /_/ 
 /_/                         
`

// showBanner is used to show the banner to the user
func showBanner() {
	gologger.Print().Msgf("%s\n", banner)
	gologger.Print().Msgf("\t\tprojectdiscovery.io\n\n")
}

// GetUpdateCallback returns a callback function that updates pdtm
func GetUpdateCallback() func() {
	return func() {
		showBanner()
		updateutils.GetUpdateToolCallback("pdtm", version)()
	}
}
