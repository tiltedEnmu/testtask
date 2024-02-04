package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tiltedEnmu/test-task/requester/internal/domain/models"
)

func (s *Service) Balance(ctx context.Context, walletId int64) (*models.Wallet, error) {
	const op = "service.Balance"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("walletId", walletId),
	)

	log.Info("attempt to get balance")

	wallet, err := s.bal.GetBalanceByWalletId(ctx, walletId)
	if err != nil {
		log.Error(
			"storage error", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			},
		)

		return &models.Wallet{}, fmt.Errorf("%s: %w", op, err)
	}

	return wallet, nil
}
