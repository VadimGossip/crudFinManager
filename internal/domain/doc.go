package domain

import "time"

type Doc struct {
	ID           int64     `json:"id"`
	Type         string    `json:"type"`
	Counterparty string    `json:"counterparty"`
	Amount       float64   `json:"amount"`
	DocCurrency  string    `json:"doc_currency"`
	AmountUsd    float64   `json:"amount_usd"`
	DocDate      time.Time `json:"doc_date"`
	Notes        string    `json:"notes"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
}

type UpdateDocInput struct {
	Type         *string    `json:"type"`
	Counterparty *string    `json:"counterparty"`
	Amount       *float64   `json:"amount"`
	DocCurrency  *string    `json:"doc_currency"`
	AmountUsd    *float64   `json:"amount_usd"`
	DocDate      *time.Time `json:"doc_date"`
	Notes        *string    `json:"notes"`
}
