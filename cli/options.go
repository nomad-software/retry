package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Options define the options passed to this program.
type Options struct {
	Command string
	Help    bool
	Pause   uint
	Timeout uint
	Tries   uint
}

// ParseOptions parses all command line options.
func ParseOptions() Options {
	var opt Options

	flag.StringVar(&opt.Command, "cmd", "", "The command to run.")
	flag.BoolVar(&opt.Help, "help", false, "Show help.")
	flag.UintVar(&opt.Pause, "p", 1, "A pause in seconds between each retry.")
	flag.UintVar(&opt.Timeout, "t", 0, "A timeout in seconds which if met, will terminate the command. The default of 0 means no timeout.")
	flag.UintVar(&opt.Tries, "r", 0, "The amount of retries to perform. The default of 0 means try forever.")
	flag.Parse()

	return opt
}

// Valid checks if the passed options are valid.
func (o *Options) Valid() bool {
	if o.Command == "" {
		fmt.Fprintln(os.Stderr, color.RedString("Command cannot be empty."))
		return false
	}

	return true
}

// PrintUsage prints the usage of this program.
func (o *Options) PrintUsage() {
	var banner = ` ____      _
|  _ \ ___| |_ _ __ _   _
| |_) / _ \ __| '__| | | |
|  _ <  __/ |_| |  | |_| |
|_| \_\___|\__|_|   \__, |
                    |___/
`
	color.Cyan(banner)
	flag.Usage()
}
