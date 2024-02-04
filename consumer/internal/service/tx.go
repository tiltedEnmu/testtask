package service

import (
	"context"
	"fmt"
	"log/slog"
)

func (s *Service) Invoice(ctx context.Context, currencyCode string, amount float64, walletId int64) error {
	const op = "service.Invoice"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("walletId", walletId),
	)

	log.Info("attempt to perform an invoice operation")

	err := s.tx.Deposit(ctx, walletId, currencyCode, amount)
	if err != nil {
		log.Error(
			"storage error", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			},
		)

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) Withdraw(ctx context.Context, currencyCode string, amount float64, walletId int64) error {
	const op = "service.Withdraw"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("walletId", walletId),
	)

	log.Info("attempt to perform a withdraw operation")

	err := s.tx.Withdraw(ctx, walletId, currencyCode, amount)
	if err != nil {
		log.Error(
			"storage error", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			},
		)

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
