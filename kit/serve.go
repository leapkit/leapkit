package main

import (
	"fmt"

	"github.com/leapkit/leapkit/kit/internal/rebuilder"
)

func serve(_ []string) error {
	err := rebuilder.Serve()
	if err != nil {
		fmt.Println("[error] starting the server:", err)
	}

	return err
}
