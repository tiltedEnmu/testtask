package service

import (
	"context"
	"log/slog"
)

type TxProvider interface {
	Deposit(ctx context.Context, walletId int64, currCode string, amount float64) error
	Withdraw(ctx context.Context, walletId int64, currCode string, amount float64) error
}

type Service struct {
	log *slog.Logger
	tx  TxProvider
}

func New(log *slog.Logger, tx TxProvider) *Service {
	return &Service{log: log, tx: tx}
}
