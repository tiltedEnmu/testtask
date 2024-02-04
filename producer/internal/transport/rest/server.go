package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tiltedEnmu/test-task/requester/internal/domain/models"
)

type TransactionService interface {
	Invoice(ctx context.Context, currencyCode string, amount float64, walletId int64) error
	Withdraw(ctx context.Context, currencyCode string, amount float64, walletId int64) error
}

type WalletService interface {
	Balance(ctx context.Context, walletId int64) (*models.Wallet, error)
}

type txReq struct {
	WalletId     int64   `json:"walletId"`
	Amount       float64 `json:"amount"`
	CurrencyCode string  `json:"currencyCode"`
}

type Api struct {
	tx  TransactionService
	wal WalletService
}

func NewHandler(txService TransactionService, walService WalletService) *http.ServeMux {
	api := &Api{tx: txService, wal: walService}
	mux := http.NewServeMux()
	mux.HandleFunc("/invoice", api.invoice)
	mux.HandleFunc("/withdraw", api.withdraw)
	mux.HandleFunc("/balance", api.balance)
	return mux
}

func (a *Api) invoice(w http.ResponseWriter, r *http.Request) {
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad request"))
		return
	}

	var req txReq
	err = json.Unmarshal(buf, &req)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad request"))
		return
	}
	err = validateInvoice(req)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	err = a.tx.Invoice(context.TODO(), req.CurrencyCode, req.Amount, req.WalletId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (a *Api) withdraw(w http.ResponseWriter, r *http.Request) {
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad request"))
		return
	}

	var req txReq
	err = json.Unmarshal(buf, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad request"))
		return
	}
	err = validateWithdraw(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	err = a.tx.Withdraw(context.TODO(), req.CurrencyCode, req.Amount, req.WalletId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (a *Api) balance(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recovered from panic: %v\n", r)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("internal server error"))
		}
	}()

	buf, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad request"))
		return
	}

	var req txReq
	err = json.Unmarshal(buf, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad request"))
		return
	}

	fmt.Println(req)

	fmt.Println(1)
	wallet, err := a.wal.Balance(context.Background(), req.WalletId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal"))
		return
	}

	marshalled, err := json.Marshal(wallet.Accounts)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(marshalled)
}

func validateInvoice(req txReq) error {
	if req.WalletId == 0 {
		return fmt.Errorf("wallet_id is empty")
	}
	if req.Amount < 0.0 {
		fmt.Println(req.Amount)
		return fmt.Errorf("amount must be bigger than 0")
	}
	if req.CurrencyCode == "" {
		return fmt.Errorf("currency_code is empty")
	}
	return nil
}

func validateWithdraw(req txReq) error {
	if req.WalletId == 0 {
		return fmt.Errorf("wallet_id is empty")
	}
	if req.Amount < 0.0 {
		return fmt.Errorf("amount must be bigger than 0")
	}
	if req.CurrencyCode == "" {
		return fmt.Errorf("currency_code is empty")
	}
	return nil
}
