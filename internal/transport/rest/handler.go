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
		docsApi.GET("", h.getDocByID)
		docsApi.GET("/list", h.GetAllDocs)
		docsApi.DELETE("/{id:[0-9]+}")
		docsApi.PUT("/{id:[0-9]+}")
	}
	return router
}

func (h *Handler) createDoc(c *gin.Context) {
	var doc domain.Doc
	if err := c.BindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err := h.docsService.Create(context.TODO(), &doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Can't create doc. Error: %s", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, fmt.Sprintf("Financial doc id = %d created %s", doc.ID, doc.Created.Format(time.RFC3339)))
}

func (h *Handler) getDocByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid id param: %s", err.Error()))
		return
	}

	doc, err := h.docsService.GetDocByID(context.TODO(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Searching for doc. Error: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, doc)
}

func (h *Handler) GetAllDocs(c *gin.Context) {
	docs, err := h.docsService.GetAllDocs(context.TODO())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Searching for doc list. Error: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, domain.GetAllDocsResponse{Docs: docs})
}
