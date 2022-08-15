package http

import (
	"context"
	"encoding/json"
	"github.com/VadimGossip/crudFinManager/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Docs interface {
	Create(ctx context.Context, doc domain.Doc) error
}

type Handler struct {
	docsService Docs
}

func NewHandler(docs Docs) *Handler {
	return &Handler{
		docsService: docs,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	return router
}

func (h *Handler) createDoc(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var doc domain.Doc
	if err = json.Unmarshal(reqBytes, &doc); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.docsService.Create(context.TODO(), doc)
	if err != nil {
		logrus.Errorf("Can't create doc %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
