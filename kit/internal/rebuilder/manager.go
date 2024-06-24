package rebuilder

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// Starts the application and reacts to changes by listening
// to the changed channel.
func runManager(changed chan bool) {
	ctx, cancel := context.WithCancel(context.Background())
	run := func(ctx context.Context, path string) {
		// Build the application
		cmd := exec.CommandContext(ctx, "go", "build")
		cmd.Args = append(cmd.Args, "-o", "bin/app", "-v", path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Println("[error]", err)
			return
		}

		// Run the application.
		cmd = exec.CommandContext(ctx, "bin/app")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			fmt.Println("[error] running:", err)
		}
	}

	go run(ctx, config.path)
	for {
		select {
		case <-changed:
			fmt.Println("[info] Restarting the server")
			cancel()

			ctx, cancel = context.WithCancel(context.Background())
			go run(ctx, config.path)
		default:
			time.Sleep(200 * time.Millisecond)
		}
	}
}
