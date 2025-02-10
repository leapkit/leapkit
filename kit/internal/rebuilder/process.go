package rebuilder

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runApp(restart chan bool) {
	e := entry{
		ID:      0,
		Name:    "app",
		Command: "go build -o bin/app ./cmd/app",
	}

	for {
		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			defer recoverFn()

			execCommand(ctx, e)

			e.Command = "bin/app"
			execCommand(ctx, e)
		}()

		<-restart
		cancel()

		fmt.Fprintln(wrap(os.Stderr, e), "Restarting the server...")

		time.Sleep(200 * time.Millisecond)
	}
}

func runProcess(restart chan bool, e entry) {
	for {
		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			defer recoverFn()
			execCommand(ctx, e)
		}()

		<-restart
		cancel()

		fmt.Fprintln(wrap(os.Stderr, e), "Restarting process...")

		time.Sleep(200 * time.Millisecond)
	}
}

func recoverFn() {
	var err error
	r := recover()
	if r == nil {
		return
	}

	switch t := r.(type) {
	case error:
		err = t
	case string:
		err = fmt.Errorf("%s", t)
	default:
		err = fmt.Errorf("%+v", t)
	}

	fmt.Println(err)
}

func execCommand(ctx context.Context, e entry) {
	fields := strings.Fields(e.Command)

	cmd := exec.CommandContext(ctx, fields[0], fields[1:]...)
	cmd.Stdout = wrap(os.Stdout, e)
	cmd.Stderr = wrap(os.Stderr, e)

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(cmd.Stderr, "error running process: %v\n", err)
	}
}
