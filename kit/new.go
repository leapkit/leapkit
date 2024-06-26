package main

import (
	"fmt"
	"os"
	"os/exec"
)

// newmodule creates a new module using the gonew tool
// and the leapkit/template template. It passes the name
// of the module as an argument to the gonew tool.
//
// Info: This function is not named  "new" because thats a reserved
// keyword in Go.
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
