package rest

import (
	"context"
	"github.com/VadimGossip/crudFinManager/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"

	_ "github.com/VadimGossip/crudFinManager/docs"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

type Docs interface {
	Create(ctx context.Context, doc domain.Doc) (int, error)
	GetDocByID(ctx context.Context, id int) (interface{}, error)
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
		docsApi.GET("/list", h.getAllDocs)
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
// @Success 201 {object} domain.CreateDocResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /docs/create [post]
func (h *Handler) createDoc(c *gin.Context) {
	var doc domain.Doc
	if err := c.ShouldBind(&doc); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "createDoc",
			"problem": "unmarshal request error",
		}).Error(err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "invalid doc param"})
		return
	}
	id, err := h.docsService.Create(context.TODO(), doc)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "createDoc",
			"problem": "service err",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "can't create doc"})
		return
	}

	c.JSON(http.StatusCreated, domain.CreateDocResponse{ID: id})
}

// @Summary Get financial document info by id
// @Tags docs
// @Description get financial document
// @ID get-doc
// @Accept json
// @Produce json
// @Param id query int true "doc id"
// @Success 200 {object} domain.Doc
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /docs [get]
func (h *Handler) getDocByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "getDocByID",
			"problem": "unmarshal request error",
		}).Error(err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "invalid id param"})
		return
	}

	doc, err := h.docsService.GetDocByID(context.TODO(), id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "getDocByID",
			"problem": "service err",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "doc search"})
		return
	}

	c.JSON(http.StatusOK, doc)
}

// @Summary Get all financial documents
// @Tags docs
// @Description get all financial documents
// @ID get-list
// @Accept json
// @Produce json
// @Success 200 {object} domain.GetAllDocsResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /docs/list [get]
func (h *Handler) getAllDocs(c *gin.Context) {
	docs, err := h.docsService.GetAllDocs(context.TODO())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "getAllDocs",
			"problem": "service err",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "doc search list"})
		return
	}

	c.JSON(http.StatusOK, domain.GetAllDocsResponse{Docs: docs})
}

// @Summary Delete financial doc by id
// @Tags docs
// @Description delete financial document by id
// @ID delete-doc
// @Accept json
// @Produce json
// @Param id query int true "doc id"
// @Success 200 {object} domain.StatusResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /docs [delete]
func (h *Handler) deleteDoc(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "deleteDoc",
			"problem": "unmarshal request error",
		}).Error(err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "invalid id param"})
		return
	}

	if err := h.docsService.Delete(context.TODO(), id); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "deleteDoc",
			"problem": "service err",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "delete doc"})
		return
	}

	c.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}

// @Summary Update financial document info
// @Tags docs
// @Description update financial document info
// @ID update-doc
// @Accept json
// @Produce json
// @Param id query int true "doc id"
// @Param input body domain.UpdateDocInput true "document update info"
// @Success 200 {object} domain.StatusResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /docs [put]
func (h *Handler) updateDocByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "updateDocByID",
			"problem": "unmarshal request error",
		}).Error(err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "invalid id param"})
		return
	}

	var input domain.UpdateDocInput
	if err := c.BindJSON(&input); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "updateDocByID",
			"problem": "unmarshal request error",
		}).Error(err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "invalid update doc input param"})
		return
	}

	if err := h.docsService.Update(context.TODO(), id, input); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "updateDocByID",
			"problem": "service err",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "update doc"})
		return
	}

	c.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}
