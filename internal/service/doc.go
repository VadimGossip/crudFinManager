package service

import (
	"context"
	"time"

	"github.com/VadimGossip/crudFinManager/internal/domain"
)

type DocsRepository interface {
	Create(ctx context.Context, doc *domain.Doc) error
	GetDocByID(ctx context.Context, id int64) (domain.Doc, error)
	GetAllDocs(ctx context.Context) ([]domain.Doc, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, inp domain.UpdateDocInput) error
}

type Docs struct {
	repo DocsRepository
}

func NewBooks(repo DocsRepository) *Docs {
	return &Docs{
		repo: repo,
	}
}

func (d *Docs) Create(ctx context.Context, doc *domain.Doc) error {
	doc.Created = time.Now()
	doc.Updated = doc.Created
	if doc.DocDate.IsZero() {
		doc.DocDate = doc.Created
	}
	return d.repo.Create(ctx, doc)
}

func (d *Docs) GetDocByID(ctx context.Context, id int64) (domain.Doc, error) {
	return d.repo.GetDocByID(ctx, id)
}

func (d *Docs) GetAllDocs(ctx context.Context) ([]domain.Doc, error) {
	return d.repo.GetAllDocs(ctx)
}

func (d *Docs) Delete(ctx context.Context, id int64) error {
	return d.repo.Delete(ctx, id)
}

func (d *Docs) Update(ctx context.Context, id int64, inp domain.UpdateDocInput) error {
	return d.repo.Update(ctx, id, inp)
}
