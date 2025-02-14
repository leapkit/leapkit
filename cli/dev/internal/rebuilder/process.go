package rebuilder

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func newProcess(e entry) *process {
	return &process{
		entry:  e,
		Stdout: wrap(os.Stdout, e),
		Stderr: wrap(os.Stderr, e),
	}
}

type process struct {
	entry
	Stdout io.Writer
	Stderr io.Writer
}

func (p *process) Run(parentCtx context.Context, reload chan bool) error {
	fields := strings.Fields(p.Command)
	name, args := fields[0], fields[1:]

	var restarted bool

	for {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		cmd := exec.CommandContext(ctx, name, args...)

		cmd.Stdout = p.Stdout
		cmd.Stderr = p.Stderr
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

		if restarted {
			fmt.Fprintln(p.Stdout, "Restarted...")
		}

		if err := cmd.Start(); err != nil {
			fmt.Fprintf(p.Stderr, "failed to start process: %v\n", err)
			return err
		}

		errCh := make(chan error, 1)
		go func() {
			if err := cmd.Wait(); err != nil {
				errCh <- err
			}
		}()

		select {
		case <-reload:
			if err := Signal(cmd, syscall.SIGTERM); err != nil {
				fmt.Fprintf(p.Stdout, "error restarting process: %v\n", err)
			}
		case <-parentCtx.Done():
			fmt.Fprintln(p.Stdout, "Stopping...")
			if err := Signal(cmd, syscall.SIGTERM); err != nil {
				fmt.Fprintf(p.Stdout, "error stopping process: %v\n", err)
			}

			return nil
		case err := <-errCh:
			fmt.Fprintf(p.Stderr, "process exited with error: %v\n", err)

			select {
			case <-reload:
			case <-parentCtx.Done():
				return nil
			}
		}

		restarted = true
	}
}

func Signal(cmd *exec.Cmd, s syscall.Signal) error {
	group, err := os.FindProcess(-cmd.Process.Pid)
	if err != nil {
		return err
	}

	return group.Signal(s)
}
