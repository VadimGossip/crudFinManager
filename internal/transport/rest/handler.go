package rest

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"github.com/VadimGossip/crudFinManager/internal/domain"
)

type Docs interface {
	Create(ctx context.Context, doc *domain.Doc) error
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
	router := gin.Default()
	docsApi := router.Group("/docs")
	{
		docsApi.POST("/create", h.createDoc)
	}
	return router
}

func (h *Handler) createDoc(c *gin.Context) {
	var doc domain.Doc
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
		return
	}
	err := h.docsService.Create(context.TODO(), &doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Can't create doc. Error: ": err.Error()})
	}

	c.JSON(http.StatusCreated, fmt.Sprintf("Financial doc id = %d created %s", doc.ID, doc.Created.Format(time.RFC3339)))
}
