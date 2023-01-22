package runner

import (
	"fmt"

	"github.com/projectdiscovery/gologger"
)

const Version = "v0.0.3"

var banner = fmt.Sprintf(`
                ____          
     ____  ____/ / /_____ ___ 
    / __ \/ __  / __/ __ __  \
   / /_/ / /_/ / /_/ / / / / /
  / .___/\__,_/\__/_/ /_/ /_/ 
 /_/                          %s
`, Version)

// showBanner is used to show the banner to the user
func showBanner() {
	gologger.Print().Msgf("%s\n", banner)
	gologger.Print().Msgf("\t\tprojectdiscovery.io\n\n")
}
