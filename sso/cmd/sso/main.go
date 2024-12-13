package main

import (
	"log/slog"
	"os"
	"sso/internal/app"
	"sso/internal/config"
)

const (
	envLocal = "local"
)

// go run cmd/sso/main.go --config=./config/config.yaml --> Запуск
func main() {

	cfg := config.MustLoad() // иницилизация конфига сервиса

	log := setupLogger(cfg.Env) // логгер сервиса
	log.Info("starting app", slog.Any("config", cfg))

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
	application.GRPCSrv.MustRun() // запуск приложения

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}
	return log
}
