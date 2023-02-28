package main

import (
	"fmt"
	"os"

	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/pdtm/pkg/utils"
)

type options struct {
	DisableUpdateCheck bool
}

func main() {
	options := &options{}
	flagSet := goflags.NewFlagSet()

	flagSet.CreateGroup("update", "Update",
		flagSet.CallbackVarP(utils.GetUpdaterCallback("subfinder"), "update", "up", "update pdtm to the latest released version"),
		flagSet.BoolVarP(&options.DisableUpdateCheck, "disable-update-check", "duc", false, "disable automatic update check"),
	)

	if err := flagSet.Parse(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	runWithOptions(options)
}

func runWithOptions(options *options) {
	if !options.DisableUpdateCheck {
		msg := utils.GetVersionCheckCallback("subfinder")()
		fmt.Println(msg)
	}
}
