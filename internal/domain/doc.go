package domain

import "time"

type Doc struct {
	ID           int64     `json:"id"`
	Type         string    `json:"type" binding:"required"`
	Counterparty string    `json:"counterparty" binding:"required"`
	Amount       float64   `json:"amount" binding:"required"`
	DocCurrency  string    `json:"doc_currency" binding:"required"`
	AmountUsd    float64   `json:"amount_usd" binding:"required"`
	DocDate      time.Time `json:"doc_date"`
	Notes        string    `json:"notes" binding:"required"`
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

type GetAllDocsResponse struct {
	Docs []Doc `json:"docs"`
}
