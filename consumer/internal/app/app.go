package app

import (
    "log/slog"

    kafkaApp "github.com/tiltedEnmu/test-task/consumer/internal/app/kafka"
    "github.com/tiltedEnmu/test-task/consumer/internal/service"
    "github.com/tiltedEnmu/test-task/consumer/internal/storage/postgres"
)

type App struct {
    KafkaApp *kafkaApp.App
}

func New(
    log *slog.Logger,
    kafkaAddr, topic string,
    pgHost, pgPort, pgUser, pgPassword, pgDBName string,
) *App {
    pgStorage, err := postgres.New(pgHost, pgPort, pgUser, pgPassword, pgDBName)
    if err != nil {
        panic("pg database module is important")
    }

    txService := service.New(log, pgStorage)

    kafka := kafkaApp.New(log, kafkaAddr, topic, txService)

    return &App{
        KafkaApp: kafka,
    }
}
