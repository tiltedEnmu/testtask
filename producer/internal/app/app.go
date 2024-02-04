package app

import (
	"log/slog"

	restapp "github.com/tiltedEnmu/test-task/requester/internal/app/rest"
	"github.com/tiltedEnmu/test-task/requester/internal/service"
	"github.com/tiltedEnmu/test-task/requester/internal/storage/kafka"
	"github.com/tiltedEnmu/test-task/requester/internal/storage/postgres"
)

type App struct {
	HttpApp *restapp.App
}

func New(
	log *slog.Logger,
	kafkaAddr, topic string,
	httpHost, httpPort string,
	pgHost, pgPort, pgUser, pgPassword, pgDbName string,
) *App {
	kafkaStorage := kafka.New(kafkaAddr, topic)
	pgStorage, err := postgres.New(pgHost, pgPort, pgUser, pgPassword, pgDbName)
	if err != nil {
		panic("pg database module is important")
	}

	services := service.New(log, kafkaStorage, pgStorage)
	httpApp := restapp.New(log, services, services, httpHost, httpPort)

	return &App{
		HttpApp: httpApp,
	}
}
