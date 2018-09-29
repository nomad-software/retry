package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Output defines the output channels.
type Output struct {
	Stdout chan string
	Stderr chan string
	Closed chan bool
}

// NewOutput returns a new console.
func NewOutput() Output {
	return Output{
		Stdout: make(chan string),
		Stderr: make(chan string),
		Closed: make(chan bool),
	}
}

// Start the output channels.
func (o *Output) Start() {
	go o.processOutput()
}

func (o *Output) processOutput() {
	for {
		select {
		case s, ok := <-o.Stdout:
			if !ok {
				o.Stdout = nil
			} else {
				fmt.Fprintln(os.Stdout, s)
			}
		case s, ok := <-o.Stderr:
			if !ok {
				o.Stderr = nil
			} else {
				fmt.Fprintln(os.Stderr, color.RedString(s))
			}
		}

		if o.Stdout == nil && o.Stderr == nil {
			break
		}
	}
	o.Closed <- true
}
