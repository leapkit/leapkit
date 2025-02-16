package main

import (
	"fmt"

	"github.com/leapkit/leapkit/cli/migration/internal/migration"
)

func main() {
	err := migration.Migration()
	if err != nil {
		fmt.Printf("[error] %v\n", err)
	}
}
