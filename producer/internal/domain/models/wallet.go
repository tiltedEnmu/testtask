package models

type Wallet struct {
    Id       int64
    Accounts []CurrencyAccount
}

type CurrencyAccount struct {
    CurrencyCode string
    Amount       float64
}
