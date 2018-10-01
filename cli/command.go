package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
)

// RunCommand runs the command specified in the options.
func RunCommand(options Options) bool {
	args := strings.Fields(options.Command)
	var cmd *exec.Cmd

	if options.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(options.Timeout)*time.Second)
		defer cancel()
		cmd = exec.CommandContext(ctx, args[0], args[1:]...)
	} else {
		cmd = exec.Command(args[0], args[1:]...)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, color.RedString(err.Error()))
		os.Exit(1)
	}

	err = cmd.Wait()
	if err == nil {
		return true
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, color.RedString(err.Error()))

		switch e := err.(type) {
		case *exec.ExitError:
			return e.Success()
		}
	}

	return false
}
