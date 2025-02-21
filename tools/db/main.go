package main

import (
	"fmt"

	"github.com/leapkit/leapkit/tools/db/internal/database"
)

func main() {
	err := database.Exec()
	if err != nil {
		fmt.Printf("[error] %v\n", err)
	}
}
