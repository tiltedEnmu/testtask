package postgres

import (
	"context"
	"fmt"

	"github.com/tiltedEnmu/test-task/requester/internal/domain/models"
)

func (s *Storage) GetBalanceByWalletId(ctx context.Context, walletId int64) (*models.Wallet, error) {
	const op = "storage.postgres.GetBalanceByWalletId"

	wallet := models.Wallet{
		Id:       walletId,
		Accounts: make([]models.CurrencyAccount, 0),
	}

	q := `select currency_code, amount from currency_accounts where wallet_id = $1`
	stmt, err := s.db.Prepare(q)
	if err != nil {
		return &models.Wallet{}, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.QueryContext(ctx, walletId)
	defer rows.Close()
	if err != nil {
		return &models.Wallet{}, fmt.Errorf("%s: %w", op, err)
	}

	var currCode string
	var amount float64
	for rows.Next() {
		if err := rows.Scan(&currCode, &amount); err != nil {
			return &models.Wallet{}, fmt.Errorf("%s: %w", op, err)
		}
		wallet.Accounts = append(
			wallet.Accounts, models.CurrencyAccount{
				CurrencyCode: currCode,
				Amount:       amount,
			},
		)
	}

	if err := rows.Err(); err != nil {
		return &models.Wallet{}, fmt.Errorf("%s: %w", op, err)
	}

	return &wallet, nil
}
