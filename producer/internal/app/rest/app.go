package rest

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/tiltedEnmu/test-task/requester/internal/transport/rest"
)

type App struct {
	log        *slog.Logger
	srv        *http.Server
	host, port string
}

func New(
	log *slog.Logger,
	txService rest.TransactionService,
	walService rest.WalletService,
	host, port string,
) *App {
	mux := rest.NewHandler(txService, walService)

	return &App{
		log: log,
		srv: &http.Server{
			Addr:           host + ":" + port,
			Handler:        mux,
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   5 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		host: host,
		port: port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "restapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.String("port", a.port),
	)

	log.Info("trying to start http server")

	err := a.srv.ListenAndServe()
	if err != nil {
		return err
	}

	log.Info("http server is running")

	return nil
}

func (a *App) Stop() {
	const op = "restapp.Stop"

	log := a.log.With(slog.String("op", op))

	log.Info("stopping http server")

	err := a.srv.Shutdown(context.Background())
	if err != nil {
		a.log.Error("graceful stop is failed", slog.String("error", err.Error()))
	}
}
