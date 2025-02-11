package rebuilder

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
)

func Serve() error {
	entries, err := readProcfile()
	if err != nil {
		return err
	}

	fmt.Println("[kit] Starting app")
	for _, entry := range entries {
		maxServiceNameLen = max(maxServiceNameLen, len(entry.Name))
	}

	fmt.Printf("\nName%s | Command\n", strings.Repeat(" ", maxServiceNameLen-len("Name")))
	for _, entry := range entries {
		fmt.Printf("%s%s | %s\n", entry.Name, strings.Repeat(" ", maxServiceNameLen-len(entry.Name)), entry.Command)
	}

	fmt.Println()

	reload := make(chan bool)

	go watcher().Watch(reload)
	for _, e := range entries {
		go newProcess(e).Run(reload)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan
	close(reload)

	fmt.Println("\n[kit] Shutting down...")

	return nil
}
