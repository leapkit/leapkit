package main

import (
	"fmt"
	"os/exec"
)

func update() error {
	fmt.Println("Updating leapkit/kit...")
	cmd := exec.Command("go", "install", "github.com/leapkit/leapkit/kit@latest")
	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Println("✅ Updated leapkit/kit to the latest version.")
	return nil
}
