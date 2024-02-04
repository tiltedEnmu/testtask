package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type kafkaMessage struct {
	WalletId int64
	CurrCode string
	Amount   float64
}

func (s *Storage) Deposit(ctx context.Context,
	walletId int64,
	currCode string,
	amount float64,
) error {
	const op = "storage.kafka.Deposit"

	kafkaMsg := kafkaMessage{
		WalletId: walletId,
		CurrCode: currCode,
		Amount:   amount,
	}

	// deposit kafka message
	msg, err := json.Marshal(kafkaMsg)
	fmt.Println(string(msg))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.writer.WriteMessages(
		ctx,
		kafka.Message{
			Headers: []kafka.Header{
				{
					Key:   "op",
					Value: []byte("deposit"),
				},
			},
			Value: msg,
		},
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Withdraw(ctx context.Context,
	walletId int64,
	currCode string,
	amount float64,
) error {
	const op = "storage.kafka.Withdraw"

	kafkaMsg := kafkaMessage{
		WalletId: walletId,
		CurrCode: currCode,
		Amount:   amount,
	}

	// deposit kafka message
	msg, err := json.Marshal(kafkaMsg)
	fmt.Println(string(msg))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.writer.WriteMessages(
		ctx,
		kafka.Message{
			Headers: []kafka.Header{
				{
					Key:   "op",
					Value: []byte("withdraw"),
				},
			},
			Value: msg,
		},
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
