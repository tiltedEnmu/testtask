package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type TransactionService interface {
	Invoice(ctx context.Context, currencyCode string, amount float64, walletId int64) error
	Withdraw(ctx context.Context, currencyCode string, amount float64, walletId int64) error
}

type Api struct {
	tx TransactionService
}

type ParsedMsg struct {
	WalletId int64   `json:"WalletId"`
	CurrCode string  `json:"CurrCode"`
	Amount   float64 `json:"Amount"`
}

func NewApi(tx TransactionService) *Api {
	return &Api{tx: tx}
}

func (a *Api) Invoice(msg kafka.Message) {
	const op = "kafka.Invoice"

	input := ParsedMsg{}

	fmt.Println(string(msg.Value))

	err := json.Unmarshal(msg.Value, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = validateMsg(input)
	if err != nil {
		fmt.Println(err)
		return
	}

	_ = a.tx.Invoice(context.TODO(), input.CurrCode, input.Amount, input.WalletId)
}

func (a *Api) Withdraw(msg kafka.Message) {
	const op = "kafka.Withdraw"

	input := ParsedMsg{}

	fmt.Println(string(msg.Value))

	err := json.Unmarshal(msg.Value, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = validateMsg(input)
	if err != nil {
		fmt.Println(err)
		return
	}

	_ = a.tx.Withdraw(context.TODO(), input.CurrCode, input.Amount, input.WalletId)
}

func validateMsg(msg ParsedMsg) error {
	if msg.WalletId == 0 {
		fmt.Println(msg.WalletId)
		return fmt.Errorf("walletId is zero value")
	}
	if msg.Amount == 0 {
		return fmt.Errorf("amount is zero value")
	}
	if msg.CurrCode == "" {
		return fmt.Errorf("currCode is zero value")
	}

	return nil
}
