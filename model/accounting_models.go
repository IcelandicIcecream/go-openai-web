package model

import "github.com/jackc/pgx/v5/pgtype"

type BankTransaction struct {
	Id              pgtype.UUID `json:"id",omitempty"`
	Description     string      `json:"description"`
	TransactionDate pgtype.Date `json:"transaction_date"`
	Reference       string      `json:"reference"`
	CurrencyCode    string      `json:"currency_code"`
	MoneyIn         float64     `json:"money_in"`
	MoneyOut        float64     `json:"money_out"`
	Balance         float64     `json:"balance"`
	Status          string      `json:"status",omitempty`
}

type BankTransactionLink struct {
	Id                pgtype.UUID `json:"id",omitempty"`
	BankTransactionId pgtype.UUID `json:"bank_transaction_id"`
	TransactionId     pgtype.UUID `json:"transaction_id"`
}

type AccountTransaction struct {
	Id              pgtype.UUID `json:"id"`
	Description     string      `json:"description"`
	TransactionDate pgtype.Date `json:"transaction_date"`
	CurrencyCode    string      `json:"currency_code"`
	Amount          float64     `json:"amount"`
	Status          string      `json:"status"`
}

type AddTransactionEntriesRequest struct {
	TransactionId string             `json:"transaction_id"`
	Entries       []TransactionEntry `json:"entries"`
}

type TransactionEntry struct {
	Id            pgtype.UUID `json:"id"`
	AccountName   string      `json:"account_name"`
	TransactionID pgtype.UUID `json:"transaction_id"`
	Amount        float64     `json:"amount"`
}

// TransactionTags struct
type TransactionTags struct {
	Id            pgtype.UUID `json:"id"`
	TransactionID pgtype.UUID `json:"transaction_id"`
	Tags          []string    `json:"tags"`
}

type Account struct {
	Id          pgtype.UUID `json:"id"`
	Name        string      `json:"name"`
	Category    string      `json:"category"`
	Description string      `json:"description"`
}

type Currency struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Category struct {
	Id                pgtype.UUID `json:"id",omitempty"`
	CategoryName      string      `json:"category_name"`
	NormalBalanceSide string      `json:"normal_balance_side"`
}
