package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/pdtm/internal/runner"
)

func main() {

	options := runner.ParseOptions()

	pdtmRunner, err := runner.NewRunner(options)
	if err != nil {
		gologger.Fatal().Msgf("Could not create runner: %s\n", err)
	}

	// Setup close handler
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			fmt.Println("\r- Ctrl+C pressed in Terminal, Exiting...")
			pdtmRunner.Close()
			os.Exit(0)
		}()
	}()

	err = pdtmRunner.Run()
	if err != nil {
		gologger.Fatal().Msgf("Could not run pdtm: %s\n", err)
	}
}
