package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/VadimGossip/crudFinManager/internal/domain"
)

type Docs interface {
	Create(ctx context.Context, doc domain.Doc) (int64, error)
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
	id, err := h.docsService.Create(context.TODO(), doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Can't create doc. Error: ": err.Error()})
	}

	c.JSON(http.StatusCreated, id)
}
