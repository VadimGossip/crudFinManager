package domain

import (
	"fmt"
	"time"
)

type Doc struct {
	ID           int       `json:"id" example:"1"`
	Type         string    `json:"type" binding:"required" example:"invoice"`
	Counterparty string    `json:"counterparty" binding:"required" example:"Some Company"`
	Amount       *float64  `json:"amount" binding:"required,min=0" example:"1.23554"`
	DocCurrency  string    `json:"doc_currency" binding:"required" example:"USD"`
	AmountUsd    *float64  `json:"amount_usd" binding:"required,min=0" example:"1.23554"`
	DocDate      time.Time `json:"doc_date" example:"2022-08-22T19:12:02.239488Z"`
	Notes        string    `json:"notes" example:"some notes"`
	AuthorID     int       `json:"author_id" example:"142"`
	CreatedAt    time.Time `json:"created_at" example:"2022-08-22T19:12:02.239488Z"`
	UpdaterID    int       `json:"updater_id" example:"253"`
	UpdatedAt    time.Time `json:"updated_at" example:"2022-08-22T19:12:02.239488Z"`
}

type UpdateDocInput struct {
	Type         *string    `json:"type" example:"invoice"`
	Counterparty *string    `json:"counterparty" example:"Some Company"`
	Amount       *float64   `json:"amount" example:"1.23554"`
	DocCurrency  *string    `json:"doc_currency" example:"USD"`
	AmountUsd    *float64   `json:"amount_usd" example:"1.23554"`
	DocDate      *time.Time `json:"doc_date" example:"2022-08-22T19:12:02.239488Z"`
	Notes        *string    `json:"notes" example:"some notes"`
}

func (u UpdateDocInput) IsValid() error {
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
	ID int `json:"id" example:"1"`
}
