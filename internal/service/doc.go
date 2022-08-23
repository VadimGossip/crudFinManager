package service

import (
	"context"
	"github.com/VadimGossip/simpleCache"
	"strconv"
	"time"

	"github.com/VadimGossip/crudFinManager/internal/domain"
)

type DocsRepository interface {
	Create(ctx context.Context, doc domain.Doc) (domain.Doc, error)
	GetDocByID(ctx context.Context, id int) (domain.Doc, error)
	GetAllDocs(ctx context.Context) ([]domain.Doc, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, inp domain.UpdateDocInput) (domain.Doc, error)
}

type Docs struct {
	repo  DocsRepository
	cache simpleCache.Cashier
}

func NewBooks(repo DocsRepository, cache simpleCache.Cashier) *Docs {
	return &Docs{
		repo:  repo,
		cache: cache,
	}
}

func (d *Docs) Create(ctx context.Context, doc domain.Doc) (int, error) {
	doc.Created = time.Now()
	doc.Updated = doc.Created
	if doc.DocDate.IsZero() {
		doc.DocDate = doc.Created
	}
	doc, err := d.repo.Create(ctx, doc)
	if err != nil {
		return 0, err
	}
	d.cache.Set(strconv.Itoa(doc.ID), doc, time.Second*10)

	return doc.ID, nil
}

func (d *Docs) GetDocByID(ctx context.Context, id int) (domain.Doc, error) {
	doc, err := d.cache.Get(strconv.Itoa(id))
	if err != nil {
		doc, err = d.repo.GetDocByID(ctx, id)
		if err != nil {
			return doc.(domain.Doc), err
		}
		d.cache.Set(strconv.Itoa(id), doc, time.Second*10)
		return doc.(domain.Doc), nil
	}
	return doc.(domain.Doc), nil
}

func (d *Docs) GetAllDocs(ctx context.Context) ([]domain.Doc, error) {
	return d.repo.GetAllDocs(ctx)
}

func (d *Docs) Delete(ctx context.Context, id int) error {
	if err := d.repo.Delete(ctx, id); err != nil {
		return err
	}
	d.cache.Delete(strconv.Itoa(id))
	return nil
}

func (d *Docs) Update(ctx context.Context, id int, inp domain.UpdateDocInput) error {
	if err := inp.Validate(); err != nil {
		return err
	}
	doc, err := d.repo.Update(ctx, id, inp)
	if err != nil {
		return err
	}
	d.cache.Set(strconv.Itoa(id), doc, time.Second*10)

	return nil
}
