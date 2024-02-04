package models

type Transaction struct {
    WalletId     int64
    Amount       float64
    CurrencyCode string
    Status       string // Created | Error | Success
}
