package service

import (
	"context"
	"log/slog"

	"github.com/tiltedEnmu/test-task/requester/internal/domain/models"
)

type TxProvider interface {
	Deposit(ctx context.Context, walletId int64, currCode string, amount float64) error
	Withdraw(ctx context.Context, walletId int64, currCode string, amount float64) error
}

type WalletProvider interface {
	GetBalanceByWalletId(ctx context.Context, walletId int64) (*models.Wallet, error)
}

type Service struct {
	log *slog.Logger
	tx  TxProvider
	bal WalletProvider
}

func New(log *slog.Logger, tx TxProvider, bal WalletProvider) *Service {
	return &Service{log: log, tx: tx, bal: bal}
}
