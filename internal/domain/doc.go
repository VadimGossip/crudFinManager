package domain

import (
	"fmt"
	"time"
)

type Doc struct {
	ID           int       `json:"id"`
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

func (u UpdateDocInput) Validate() error {
	if u.Type == nil && u.Counterparty == nil && u.Amount == nil && u.DocCurrency == nil && u.AmountUsd == nil &&
		u.Notes == nil && u.DocDate == nil {
		return fmt.Errorf("update structure has no values")
	}
	return nil
}

type GetAllDocsResponse struct {
	Docs []Doc `json:"docs"`
}

type CreateDocResponse struct {
	ID int `json:"id"`
}

type UpdateDocResponse struct {
	Status string `json:"status"`
}

type DeleteDocResponse struct {
	Status string `json:"status"`
}
