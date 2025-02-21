package main

import (
	"context"
	"fmt"

	"github.com/leapkit/leapkit/tools/dev/internal/rebuilder"
)

func main() {
	err := rebuilder.Serve(context.Background())
	if err != nil {
		fmt.Println("[error] starting the server:", err)
	}
}
