package main

import (
	"fmt"
	"main/internal/config"
)

//go run cmd/auth/main.go --config=./config_auth/auth.yaml --> Запуск

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
}
