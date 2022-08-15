package service

import (
	"context"
	"time"

	"github.com/VadimGossip/crudFinManager/internal/domain"
)

type DocsRepository interface {
	Create(ctx context.Context, doc domain.Doc) (int64, error)
}

type Docs struct {
	repo DocsRepository
}

func NewBooks(repo DocsRepository) *Docs {
	return &Docs{
		repo: repo,
	}
}

func (d *Docs) Create(ctx context.Context, doc domain.Doc) (int64, error) {
	doc.Created = time.Now()
	doc.Updated = doc.Created
	if doc.DocDate.IsZero() {
		doc.DocDate = doc.Created
	}
	return d.repo.Create(ctx, doc)
}
