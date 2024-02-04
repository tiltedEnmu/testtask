package postgres

import (
	"context"
	"fmt"
)

func (s *Storage) Deposit(ctx context.Context, walletId int64, currCode string, amount float64) error {
	const op = "storage.postgres.Deposit"

	fmt.Println(op)

	_, err := s.db.ExecContext(
		ctx,
		`update currency_accounts set amount = amount + $1 where wallet_id = $2 and currency_code = $3`,
		amount, walletId, currCode,
	)
	fmt.Println(err)
	if err != nil {
		_ = s.ErrorTx(ctx, walletId, currCode, amount)
		return fmt.Errorf("%s: %w", op, err)
	}

	_ = s.SuccessTx(ctx, walletId, currCode, amount)
	return nil
}

func (s *Storage) Withdraw(ctx context.Context, walletId int64, currCode string, amount float64) error {
	const op = "storage.postgres.Withdraw"

	fmt.Println(op)

	_, err := s.db.ExecContext(
		ctx,
		`update currency_accounts set amount = amount - $1 where wallet_id = $2 and currency_code = $3`,
		amount, walletId, currCode,
	)
	fmt.Println(err)
	if err != nil {
		_ = s.ErrorTx(ctx, walletId, currCode, amount*(-1))
		return fmt.Errorf("%s: %w", op, err)
	}

	_ = s.SuccessTx(ctx, walletId, currCode, amount*(-1))
	return nil
}

// ErrorTx saves transactions that were not applied due to an error
func (s *Storage) ErrorTx(ctx context.Context, walletId int64, currCode string, amount float64) error {
	const op = "storage.postgres.ErrorTx"

	q := `insert into transactions (wallet_id, currency_code, amount, status) values (?, ?, ?, ?)`
	stmt, err := s.db.Prepare(q)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, walletId, currCode, amount, "Error")
	if err != nil {
		_ = s.ErrorTx(ctx, walletId, currCode, amount)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// SuccessTx saves applied transactions
func (s *Storage) SuccessTx(ctx context.Context, walletId int64, currCode string, amount float64) error {
	const op = "storage.postgres.SuccessTx"

	q := `insert into transactions (wallet_id, currency_code, amount, status) values ($1, $2, $3, $4)`
	stmt, err := s.db.Prepare(q)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, walletId, currCode, amount, "Success")
	if err != nil {
		_ = s.ErrorTx(ctx, walletId, currCode, amount)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
