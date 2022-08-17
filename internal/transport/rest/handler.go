package rest

import (
	"context"
	"fmt"
	"github.com/VadimGossip/crudFinManager/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	_ "github.com/VadimGossip/crudFinManager/docs"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

type Docs interface {
	Create(ctx context.Context, doc domain.Doc) (int, error)
	GetDocByID(ctx context.Context, id int) (domain.Doc, error)
	GetAllDocs(ctx context.Context) ([]domain.Doc, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, inp domain.UpdateDocInput) error
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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	docsApi := router.Group("/docs")
	{
		docsApi.POST("/create", h.createDoc)
		docsApi.GET("", h.getDocByID)
		docsApi.GET("/list", h.GetAllDocs)
		docsApi.DELETE("", h.deleteDoc)
		docsApi.PUT("", h.updateDocByID)
	}
	return router
}

// @Summary Create Financial document
// @Tags docs
// @Description create financial document
// @ID create-doc
// @Accept json
// @Produce json
// @Param input body domain.Doc true "document info"
// @Success 201 {integer} {object} createDocResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Failure default {object} string
// @Router /docs/create [post]
func (h *Handler) createDoc(c *gin.Context) {
	var doc domain.Doc
	if err := c.ShouldBind(&doc); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.docsService.Create(context.TODO(), doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Can't create doc. Error: %s", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, domain.CreateDocResponse{ID: id})
}

func (h *Handler) getDocByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
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

func (h *Handler) deleteDoc(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid id param: %s", err.Error()))
		return
	}

	if err := h.docsService.Delete(context.TODO(), id); err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Deleting doc. Error: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, domain.DeleteDocResponse{Status: "ok"})
}

func (h *Handler) updateDocByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid id param: %s", err.Error()))
		return
	}

	var input domain.UpdateDocInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := h.docsService.Update(context.TODO(), id, input); err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Updating doc. Error: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, domain.UpdateDocResponse{Status: "ok"})
}
