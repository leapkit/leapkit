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
	"time"
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

func (p *process) Run(ctx context.Context, reloadSignal chan bool) error {
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

		errCh := make(chan error, 1)
		go func() {
			if err := cmd.Wait(); err != nil {
				errCh <- err
			}
		}()

		select {
		case <-reloadSignal:
			fmt.Fprintln(p.Stdout, "Reloading...")
			syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
			cancel()
		case <-ctx.Done():
			fmt.Fprintln(p.Stdout, "Stopping...")
			syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
			cancel()
			return nil
		case err := <-errCh:
			if err != nil {
				fmt.Fprintf(p.Stderr, "process exited with error: %v\n", err)
				cancel()
				return err
			}
		}

		time.Sleep(200 * time.Millisecond)
		fmt.Fprintln(p.Stdout, "Started...")
	}
}
