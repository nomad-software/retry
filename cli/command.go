package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
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

	forwardOutput(cmd)

	err := cmd.Wait()

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

// ForwardOutput forwards all output from the command to the console.
func forwardOutput(cmd *exec.Cmd) bool {
	output := NewOutput()
	output.Start()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		output.Stderr <- err.Error()
		os.Exit(1)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		output.Stderr <- err.Error()
		os.Exit(1)
	}

	err = cmd.Start()
	if err != nil {
		output.Stderr <- err.Error()
		os.Exit(1)
	}

	go forwardPipe(stdout, output.Stdout)
	go forwardPipe(stderr, output.Stderr)

	return <-output.Closed
}

// ForwardPipe forwards output from the specified pipe to the console.
func forwardPipe(pipe io.ReadCloser, channel chan string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		channel <- scanner.Text()
	}
	close(channel)
}
