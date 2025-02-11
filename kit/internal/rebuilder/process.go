package rebuilder

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type process struct {
	entry
	Stdout io.Writer
	Stderr io.Writer
}

func (p *process) Run(reloadSignal chan bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(p.Stderr, "error running process: %v\n", r)
		}
	}()

	var mainCmd string
	var args []string

	fields := strings.Fields(p.Command)
	if len(fields) > 0 {
		mainCmd = fields[0]
		args = fields[1:]
	}

	for {
		ctx, cancel := context.WithCancel(context.Background())
		cmd := exec.CommandContext(ctx, mainCmd, args...)
		cmd.Stdout = p.Stdout
		cmd.Stderr = p.Stderr

		errChan := make(chan error, 1)

		go func() {
			errChan <- cmd.Run()
		}()

		select {
		case <-reloadSignal:
			fmt.Fprintln(p.Stdout, "Restarting process...")

			if cmd.Process != nil {
				cmd.Process.Kill()
			}

			cancel()
			<-errChan

			time.Sleep(200 * time.Millisecond)

		case err := <-errChan:
			if err != nil {
				fmt.Fprintf(p.Stderr, "process exited with error: %v\n", err)
			}

			cancel()
			return
		}
	}
}

func newProcess(e entry) *process {
	return &process{
		entry:  e,
		Stdout: wrap(os.Stdout, e),
		Stderr: wrap(os.Stderr, e),
	}
}
