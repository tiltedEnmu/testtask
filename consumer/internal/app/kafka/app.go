package kafka

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/segmentio/kafka-go"
	tsKafka "github.com/tiltedEnmu/test-task/consumer/internal/transport/kafka" // internal transport kafka module
)

type App struct {
	log    *slog.Logger
	reader *kafka.Reader
	addr   string
	api    *tsKafka.Api
}

func New(
	log *slog.Logger,
	addr string,
	topic string,
	txService tsKafka.TransactionService,
) *App {

	return &App{
		log:  log,
		addr: addr,
		reader: kafka.NewReader(
			kafka.ReaderConfig{
				Brokers:   []string{addr},
				Topic:     topic,
				Partition: 0,
				MaxBytes:  1 << 20,
			},
		),
		api: tsKafka.NewApi(txService),
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "kafkaApp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.String("addr", a.addr),
	)

	log.Info("trying to start kafka reader")

	go func() {
		for {
			m, err := a.reader.ReadMessage(context.Background())
			log.Info(
				"message received", slog.Attr{
					Key:   "offset",
					Value: slog.StringValue(strconv.FormatInt(m.Offset, 10)),
				},
			)

			if err != nil {
				log.Error(
					"kafka message error: ", slog.Attr{
						Key:   "error",
						Value: slog.StringValue(err.Error()),
					},
				)
				break
			}
			for _, k := range m.Headers {
				if k.Key == "op" {
					if "deposit" == string(k.Value) {
						log.Info("invoice")
						a.api.Invoice(m)
					}
					if "withdraw" == string(k.Value) {
						log.Info("withdraw")
						a.api.Withdraw(m)
					}
				}
			}
		}
	}()

	log.Info("kafka reader is running")

	return nil
}

func (a *App) Stop() {
	const op = "kafkaApp.Stop"

	log := a.log.With(
		slog.String("op", op),
		slog.String("addr", a.addr),
	)

	err := a.reader.Close()
	if err != nil {
		log.Error(
			"kafka stop failed", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			},
		)
	}
}
