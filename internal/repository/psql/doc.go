package psql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/VadimGossip/crudFinManager/internal/domain"
)

type Docs struct {
	db *sql.DB
}

func NewDocs(db *sql.DB) *Docs {
	return &Docs{db: db}
}

func (d *Docs) Create(ctx context.Context, doc *domain.Doc) error {
	createStmt := "INSERT INTO docs(type, counterparty, amount, doc_currency, amount_usd, doc_date, notes, author_id, created_at, updater_id, updated_at)" +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id"
	err := d.db.QueryRowContext(ctx, createStmt,
		doc.Type, doc.Counterparty, doc.Amount, doc.DocCurrency, doc.AmountUsd, doc.DocDate, doc.Notes, doc.AuthorID, doc.CreatedAt, doc.UpdaterID, doc.UpdatedAt).
		Scan(&doc.ID)

	return err
}

func (d *Docs) GetDocByID(ctx context.Context, id int) (domain.Doc, error) {
	var doc domain.Doc
	selectStmt := "SELECT id, type, counterparty, amount, doc_currency, amount_usd, doc_date, notes, author_id, created_at, updater_id, updated_at" +
		" FROM docs where id=$1"
	err := d.db.QueryRowContext(ctx, selectStmt, id).
		Scan(&doc.ID, &doc.Type, &doc.Counterparty, &doc.Amount, &doc.DocCurrency,
			&doc.AmountUsd, &doc.DocDate, &doc.Notes, &doc.AuthorID, &doc.CreatedAt, &doc.UpdaterID, &doc.UpdatedAt)
	if err == sql.ErrNoRows {
		return doc, fmt.Errorf("document with id = %d not found", id)
	}
	return doc, err
}

func (d *Docs) GetAllDocs(ctx context.Context) ([]domain.Doc, error) {
	selectStmt := "SELECT id, type, counterparty, amount, doc_currency, amount_usd, doc_date, notes, author_id, created_at, updater_id, updated_at" +
		" FROM docs"
	rows, err := d.db.QueryContext(ctx, selectStmt)
	if err != nil {
		return nil, err
	}

	docs := make([]domain.Doc, 0)
	for rows.Next() {
		var doc domain.Doc
		err := rows.Scan(&doc.ID, &doc.Type, &doc.Counterparty, &doc.Amount, &doc.DocCurrency,
			&doc.AmountUsd, &doc.DocDate, &doc.Notes, &doc.AuthorID, &doc.CreatedAt, &doc.UpdaterID, &doc.UpdatedAt)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, rows.Err()
}

func (d *Docs) Delete(ctx context.Context, id int) error {
	deleteStmt := "DELETE FROM docs WHERE id=$1"
	_, err := d.db.ExecContext(ctx, deleteStmt, id)
	return err
}

func (d *Docs) Update(ctx context.Context, userId, id int, inp domain.UpdateDocInput) (domain.Doc, error) {
	var doc domain.Doc
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if inp.Type != nil {
		setValues = append(setValues, fmt.Sprintf("type=$%d", argId))
		args = append(args, *inp.Type)
		argId++
	}

	if inp.Counterparty != nil {
		setValues = append(setValues, fmt.Sprintf("counterparty=$%d", argId))
		args = append(args, *inp.Counterparty)
		argId++
	}

	if inp.Amount != nil {
		setValues = append(setValues, fmt.Sprintf("amount=$%d", argId))
		args = append(args, *inp.Amount)
		argId++
	}

	if inp.DocCurrency != nil {
		setValues = append(setValues, fmt.Sprintf("doc_currency=$%d", argId))
		args = append(args, *inp.DocCurrency)
		argId++
	}

	if inp.AmountUsd != nil {
		setValues = append(setValues, fmt.Sprintf("amount_usd=$%d", argId))
		args = append(args, *inp.AmountUsd)
		argId++
	}

	if inp.DocDate != nil {
		setValues = append(setValues, fmt.Sprintf("doc_date=$%d", argId))
		args = append(args, *inp.DocDate)
		argId++
	}

	if inp.Notes != nil {
		setValues = append(setValues, fmt.Sprintf("notes=$%d", argId))
		args = append(args, *inp.Notes)
		argId++
	}

	setValues = append(setValues, fmt.Sprintf("updater_id=$%d", argId))
	args = append(args, userId)
	argId++

	setQuery := strings.Join(setValues, ", ")
	updStmt := "UPDATE docs SET %s where id=$%d RETURNING" +
		" id, type, counterparty, amount, doc_currency, amount_usd, doc_date, notes, author_id, created_at, updater_id, updated_at"

	query := fmt.Sprintf(updStmt, setQuery, argId)
	args = append(args, id)

	err := d.db.QueryRowContext(ctx, query, args...).
		Scan(&doc.ID, &doc.Type, &doc.Counterparty, &doc.Amount, &doc.DocCurrency,
			&doc.AmountUsd, &doc.DocDate, &doc.Notes, &doc.AuthorID, &doc.CreatedAt, &doc.UpdaterID, &doc.UpdatedAt)

	return doc, err
}
