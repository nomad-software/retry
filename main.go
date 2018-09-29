package main

import (
	"time"

	"github.com/nomad-software/retry/cli"
)

func main() {
	options := cli.ParseOptions()

	if options.Help {
		options.PrintUsage()

	} else if options.Valid() {
		var tries uint = 1
		for !cli.RunCommand(options) {
			if tries == options.Tries {
				return
			}
			tries++
			time.Sleep(time.Duration(options.Pause) * time.Second)
		}
	}
}
