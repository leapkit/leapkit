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

	restart <-chan bool
}

func (p *process) Run() {
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
		case <-p.restart:
			fmt.Fprintln(p.Stdout, "Restarting process...")

			if cmd.Process != nil {
				cmd.Process.Kill()
			}

			cancel()
			<-errChan

			for len(p.restart) > 0 {
				<-p.restart
			}

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

func newProcess(e entry, restart chan bool) *process {
	return &process{
		entry:   e,
		Stdout:  wrap(os.Stdout, e),
		Stderr:  wrap(os.Stderr, e),
		restart: restart,
	}
}
