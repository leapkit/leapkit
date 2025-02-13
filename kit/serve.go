package main

import (
	"context"
	"fmt"

	"github.com/leapkit/leapkit/kit/internal/rebuilder"
)

func serve(_ []string) error {
	err := rebuilder.Serve(context.Background())
	if err != nil {
		fmt.Println("[error] starting the server:", err)
	}

	return err
}
