package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/tiltedEnmu/test-task/requester/internal/app"
	"github.com/tiltedEnmu/test-task/requester/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	log.Info("starting application", slog.Any("config", cfg))

	application := app.New(
		log,
		cfg.Kafka.Addr, cfg.Kafka.Topic,
		cfg.HTTP.Host, cfg.HTTP.Port,
		cfg.Postgres.Host, cfg.Postgres.Port,
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName,
	)

	go application.HttpApp.MustRun()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping application", slog.String("signal", sign.String()))

	application.HttpApp.Stop()

	log.Info("gracefully stopped")
}
