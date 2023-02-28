package main

import (
	"fmt"

	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/pdtm/pkg/utils"
)

type options struct {
	DisableUpdateCheck bool
}

func main() {
	options := &options{}
	flagSet := goflags.NewFlagSet()
	toolName := "nuclei"

	flagSet.CreateGroup("update", "Update",
		flagSet.CallbackVarP(utils.GetUpdaterCallback(toolName), "update", "up", fmt.Sprintf("update %v to the latest released version", toolName)),
		flagSet.BoolVarP(&options.DisableUpdateCheck, "disable-update-check", "duc", false, "disable automatic update check"),
	)

	if err := flagSet.Parse(); err != nil {
		panic(err)
	}

	if !options.DisableUpdateCheck {
		msg := utils.GetVersionCheckCallback(toolName)()
		fmt.Println(msg)
	}
}
