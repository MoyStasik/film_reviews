package main

import (
	"fmt"
	"main/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)
}
