package rebuilder

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Serve() error {
	entries, err := readProcfile()
	if err != nil {
		return err
	}

	fmt.Println("[kit] Starting app")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	reload := make(chan bool)
	exitCh := make(chan error, len(entries))

	go watcher().Watch(reload)
	for _, e := range entries {
		go func() {
			exitCh <- newProcess(e).Run(ctx, reload)
		}()
	}

	<-ctx.Done()
	fmt.Println()

	for i := 0; i < len(entries); i++ {
		<-exitCh
	}

	fmt.Println("[kit] Shutting down...")

	return nil
}
