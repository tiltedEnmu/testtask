package main

import (
    "log/slog"
    "os"
    "os/signal"
    "syscall"

    "github.com/tiltedEnmu/test-task/consumer/internal/app"
    "github.com/tiltedEnmu/test-task/consumer/internal/config"
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
        cfg.Postgres.Host, cfg.Postgres.Port,
        cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName,
    )

    go application.KafkaApp.MustRun()

    // Graceful shutdown
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

    sign := <-stop

    log.Info("stopping application", slog.String("signal", sign.String()))

    application.KafkaApp.Stop()

    log.Info("gracefully stopped")
}
