package psql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/VadimGossip/crudFinManager/internal/domain"
)

type Docs struct {
	db *sql.DB
}

func NewDocs(db *sql.DB) *Docs {
	return &Docs{db: db}
}

func (d *Docs) Create(ctx context.Context, doc *domain.Doc) error {
	createStmt := "insert into docs(type, counterparty, amount, doc_currency, amount_usd, doc_date, notes, created, updated)" +
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id, doc_date, created, updated"
	err := d.db.QueryRowContext(ctx, createStmt,
		doc.Type, doc.Counterparty, doc.Amount, doc.DocCurrency, doc.AmountUsd, doc.DocDate, doc.Notes, doc.Created, doc.Updated).
		Scan(&doc.ID, &doc.DocDate, &doc.Created, &doc.Updated)
	if err != nil {
		return err
	}
	fmt.Println(doc)
	return err
}

//
//func (b *Books) GetByID(ctx context.Context, id int64) (domain.Book, error) {
//	var book domain.Book
//	err := b.db.QueryRow("SELECT id, title, author, publish_date, rating FROM books WHERE id=$1", id).
//		Scan(&book.ID, &book.Title, &book.Author, &book.PublishDate, &book.Rating)
//	if err == sql.ErrNoRows {
//		return book, domain.ErrBookNotFound
//	}
//
//	return book, err
//}
//
//func (b *Books) GetAll(ctx context.Context) ([]domain.Book, error) {
//	rows, err := b.db.Query("SELECT id, title, author, publish_date, rating FROM books")
//	if err != nil {
//		return nil, err
//	}
//
//	books := make([]domain.Book, 0)
//	for rows.Next() {
//		var book domain.Book
//		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishDate, &book.Rating); err != nil {
//			return nil, err
//		}
//
//		books = append(books, book)
//	}
//
//	return books, rows.Err()
//}
//
//func (b *Books) Delete(ctx context.Context, id int64) error {
//	_, err := b.db.Exec("DELETE FROM books WHERE id=$1", id)
//
//	return err
//}
//
//func (b *Books) Update(ctx context.Context, id int64, inp domain.UpdateBookInput) error {
//	setValues := make([]string, 0)
//	args := make([]interface{}, 0)
//	argId := 1
//
//	if inp.Title != nil {
//		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
//		args = append(args, *inp.Title)
//		argId++
//	}
//
//	if inp.Author != nil {
//		setValues = append(setValues, fmt.Sprintf("author=$%d", argId))
//		args = append(args, *inp.Author)
//		argId++
//	}
//
//	if inp.PublishDate != nil {
//		setValues = append(setValues, fmt.Sprintf("publish_date=$%d", argId))
//		args = append(args, *inp.PublishDate)
//		argId++
//	}
//
//	if inp.Rating != nil {
//		setValues = append(setValues, fmt.Sprintf("rating=$%d", argId))
//		args = append(args, *inp.Rating)
//		argId++
//	}
//
//	setQuery := strings.Join(setValues, ", ")
//
//	query := fmt.Sprintf("UPDATE books SET %s WHERE id=$%d", setQuery, argId)
//	args = append(args, id)
//
//	_, err := b.db.Exec(query, args...)
//	return err
//}
