package main

import (
	"fmt"
	"os"
	"os/exec"
)

func newmodule(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: new [name]")
		return
	}

	cmd := exec.Command(
		"go", "run", "rsc.io/tmp/gonew@latest",
		"github.com/leapkit/leapkit/template@latest",
		args[1],
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
