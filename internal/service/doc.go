package service

import (
	"context"
	audit "github.com/VadimGossip/grpcAuditLog/pkg/domain"
	"github.com/VadimGossip/simpleCache"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"

	"github.com/VadimGossip/crudFinManager/internal/domain"
)

type DocsRepository interface {
	Create(ctx context.Context, doc *domain.Doc) error
	GetDocByID(ctx context.Context, id int) (domain.Doc, error)
	GetAllDocs(ctx context.Context) ([]domain.Doc, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, userId, id int, inp domain.UpdateDocInput) (domain.Doc, error)
}

type Docs struct {
	repo        DocsRepository
	auditClient DocsAuditClient
	cache       simpleCache.Cashier
}

type DocsAuditClient interface {
	SendLogRequest(ctx context.Context, req audit.LogItem) error
}

func NewBooks(repo DocsRepository, auditClient DocsAuditClient, cache simpleCache.Cashier) *Docs {
	return &Docs{
		repo:        repo,
		auditClient: auditClient,
		cache:       cache,
	}
}

func (d *Docs) Create(ctx context.Context, userId int, doc domain.Doc) (int, error) {
	doc.CreatedAt = time.Now()
	doc.UpdatedAt = doc.CreatedAt
	if doc.DocDate.IsZero() {
		doc.DocDate = doc.CreatedAt
	}
	doc.AuthorID = userId
	doc.UpdaterID = doc.AuthorID

	err := d.repo.Create(ctx, &doc)
	if err != nil {
		return 0, err
	}
	d.cache.Set(strconv.Itoa(doc.ID), doc, time.Second*10)

	if err := d.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_CREATE,
		Entity:    audit.ENTITY_DOC,
		EntityID:  int64(doc.ID),
		AuthorID:  int64(doc.AuthorID),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Docs.Create",
		}).Error("failed to send log request:", err)
	}

	return doc.ID, nil
}

func (d *Docs) GetDocByID(ctx context.Context, userId, id int) (domain.Doc, error) {
	doc, err := d.cache.Get(strconv.Itoa(id))
	if err != nil {
		doc, err = d.repo.GetDocByID(ctx, id)
		if err != nil {
			return doc.(domain.Doc), err
		}
		d.cache.Set(strconv.Itoa(id), doc, time.Second*10)
		return doc.(domain.Doc), nil
	}
	if err := d.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_GET,
		Entity:    audit.ENTITY_DOC,
		EntityID:  int64(doc.(domain.Doc).ID),
		AuthorID:  int64(userId),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Docs.GetDocByID",
		}).Error("failed to send log request:", err)
	}

	return doc.(domain.Doc), nil
}

func (d *Docs) GetAllDocs(ctx context.Context) ([]domain.Doc, error) {
	return d.repo.GetAllDocs(ctx)
}

func (d *Docs) Delete(ctx context.Context, userId, id int) error {
	if err := d.repo.Delete(ctx, id); err != nil {
		return err
	}
	d.cache.Delete(strconv.Itoa(id))
	if err := d.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_DELETE,
		Entity:    audit.ENTITY_DOC,
		EntityID:  int64(id),
		AuthorID:  int64(userId),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Docs.Delete",
		}).Error("failed to send log request:", err)
	}
	return nil
}

func (d *Docs) Update(ctx context.Context, userId, id int, inp domain.UpdateDocInput) error {
	if !inp.IsValid() {
		return domain.ErrInvalidUpdateStruct
	}

	doc, err := d.repo.Update(ctx, userId, id, inp)
	if err != nil {
		return err
	}
	d.cache.Set(strconv.Itoa(id), doc, time.Second*10)

	if err := d.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_UPDATE,
		Entity:    audit.ENTITY_DOC,
		EntityID:  int64(id),
		AuthorID:  int64(doc.UpdaterID),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Docs.Delete",
		}).Error("failed to send log request:", err)
	}

	return nil
}
