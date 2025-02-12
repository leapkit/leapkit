package rebuilder

import (
	"context"
	"errors"
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

func (p *process) Run(ctx context.Context, reload chan bool) error {
	fields := strings.Fields(p.Command)
	if len(fields) == 0 {
		fmt.Fprintln(p.Stderr, "error: command is empty")
		return errors.New("command is empty")
	}

	mainCmd, args := fields[0], fields[1:]

	for {
		pCtx, cancel := context.WithCancel(context.Background())
		cmd := exec.CommandContext(pCtx, mainCmd, args...)

		cmd.Stdout = p.Stdout
		cmd.Stderr = p.Stderr
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

		if err := cmd.Start(); err != nil {
			fmt.Fprintf(p.Stderr, "failed to start process: %v\n", err)
			cancel()
			return err
		}

		fmt.Fprintln(p.Stdout, "Started... (pid:", cmd.Process.Pid, ")")

		errCh := make(chan error, 1)
		go func() {
			if err := cmd.Wait(); err != nil {
				errCh <- err
			}
		}()

		select {
		case <-reload:
			fmt.Fprintln(p.Stdout, "Reloading...")
			if err := Signal(cmd, syscall.SIGTERM, cancel); err != nil {
				fmt.Fprintf(p.Stderr, "error sending SIGTERM: %v\n", err)
			}
		case <-ctx.Done():
			fmt.Fprintln(p.Stdout, "Stopping...")
			if err := Signal(cmd, syscall.SIGTERM, cancel); err != nil {
				fmt.Fprintf(p.Stderr, "error sending SIGTERM: %v\n", err)
			}

			return nil
		case err := <-errCh:
			fmt.Fprintf(p.Stderr, "process exited with error: %v\n", err)
			cancel()
			return err
		}

		fmt.Fprintln(p.Stdout, "Restarted...")
	}
}

func Signal(cmd *exec.Cmd, s syscall.Signal, cancel context.CancelFunc) error {
	group, err := os.FindProcess(-cmd.Process.Pid)
	if err != nil {
		return err
	}

	if err := group.Signal(s); err != nil {
		return err
	}

	cancel()

	return nil
}
