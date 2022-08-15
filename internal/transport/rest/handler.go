package rest

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"

	"github.com/VadimGossip/crudFinManager/internal/domain"
)

type Docs interface {
	Create(ctx context.Context, doc *domain.Doc) error
	GetDocByID(ctx context.Context, id int64) (domain.Doc, error)
	GetAllDocs(ctx context.Context) ([]domain.Doc, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, inp domain.UpdateDocInput) error
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
		docsApi.GET("/list")
		docsApi.GET("/", h.getDocByID)
		docsApi.DELETE("/{id:[0-9]+}")
		docsApi.PUT("/{id:[0-9]+}")
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
		return
	}

	c.JSON(http.StatusCreated, fmt.Sprintf("Financial doc id = %d created %s", doc.ID, doc.Created.Format(time.RFC3339)))
}

func (h *Handler) getDocByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid id param: ": err.Error()})
		return
	}

	doc, err := h.docsService.GetDocByID(context.TODO(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Searching for doc. Error: ": err.Error()})
		return
	}

	c.JSON(http.StatusOK, doc)
}
