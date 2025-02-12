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

	reloadCh := make([]chan bool, len(entries))
	exitCh := make(chan error, len(entries))

	go new(watcher).Watch(reloadCh)
	for i, e := range entries {
		reloadCh[i] = make(chan bool)
		go func() {
			exitCh <- newProcess(e).Run(ctx, reloadCh[i])
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
