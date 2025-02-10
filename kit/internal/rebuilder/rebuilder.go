package rebuilder

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
)

func Serve() error {
	entries, err := procfile()
	if err != nil {
		return err
	}

	reload := make(chan bool)
	go watcher().Watch(reload)

	printHeader(entries)

	go runApp(reload)
	for _, e := range entries {
		if e.Name == "app" {
			continue
		}

		go runProcess(reload, e)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan
	close(reload)

	fmt.Println()
	fmt.Println("[kit] Shutting down...")

	return nil
}

func printHeader(entries []entry) {
	fmt.Println("[kit] Starting app")

	for _, entry := range entries {
		maxServiceNameLen = max(maxServiceNameLen, len(entry.Name))
	}

	for _, entry := range entries {
		fmt.Printf("%s%s | %s\n", entry.Name, strings.Repeat(" ", maxServiceNameLen-len(entry.Name)), entry.Command)
	}

	fmt.Println()
}
