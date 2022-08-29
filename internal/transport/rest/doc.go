package rest

import (
	"github.com/VadimGossip/crudFinManager/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary Create Financial document
// @Tags docs
// @Description create financial document
// @ID create-doc
// @Accept  json
// @Produce json
// @Param   input body     domain.Doc               true "document info"
// @Success 201   {object} domain.CreateDocResponse
// @Failure 400   {object} domain.ErrorResponse
// @Failure 500   {object} domain.ErrorResponse
// @Router /docs [post]
func (h *Handler) createDoc(c *gin.Context) {
	var doc domain.Doc
	if err := c.ShouldBind(&doc); err != nil {
		logError("createDoc", "unmarshal request error", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "invalid doc param"})
		return
	}
	id, err := h.docsService.Create(c.Request.Context(), doc)
	if err != nil {
		logError("createDoc", "service error", err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "can't create doc"})
		return
	}

	c.JSON(http.StatusCreated, domain.CreateDocResponse{ID: id})
}

// @Summary Get financial document info by id
// @Tags docs
// @Description get financial document
// @ID get-doc-by-id
// @Accept  json
// @Produce json
// @Param   id   path    int                  true  "Doc ID"
// @Success 200 {object} domain.Doc
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /docs/{id} [get]
func (h *Handler) getDocByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logError("getDocByID", "request param error", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "invalid id param"})
		return
	}

	doc, err := h.docsService.GetDocByID(c.Request.Context(), id)
	if err != nil {
		logError("getDocByID", "service error", err)
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
// @Router /docs [get]
func (h *Handler) getAllDocs(c *gin.Context) {
	docs, err := h.docsService.GetAllDocs(c.Request.Context())
	if err != nil {
		logError("getAllDocs", "service error", err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "doc search list"})
		return
	}

	c.JSON(http.StatusOK, domain.GetAllDocsResponse{Docs: docs})
}

// @Summary Delete financial doc by id
// @Tags docs
// @Description delete financial document by id
// @ID delete-doc-by-id
// @Accept  json
// @Produce json
// @Param   id   path    int                   true  "Doc ID"
// @Success 200 {object} domain.StatusResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /docs/{id} [delete]
func (h *Handler) deleteDocByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logError("deleteDocByID", "request param error", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "invalid id param"})
		return
	}

	if err := h.docsService.Delete(c.Request.Context(), id); err != nil {
		logError("deleteDocByID", "service error", err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "delete doc"})
		return
	}

	c.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}

// @Summary Update financial document info
// @Tags docs
// @Description update financial document info
// @ID update-doc-by-id
// @Accept  json
// @Produce json
// @Param   id    path     int                   true  "Doc ID"
// @Param   input body     domain.UpdateDocInput true "document update info"
// @Success 200   {object} domain.StatusResponse
// @Failure 400   {object} domain.ErrorResponse
// @Failure 500   {object} domain.ErrorResponse
// @Router /docs/{id} [put]
func (h *Handler) updateDocByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logError("updateDocByID", "request param error", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "invalid id param"})
		return
	}

	var input domain.UpdateDocInput
	if err := c.BindJSON(&input); err != nil {
		logError("updateDocByID", "unmarshal request error", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "invalid update doc input param"})
		return
	}

	if err := h.docsService.Update(c.Request.Context(), id, input); err != nil {
		logError("updateDocByID", "service error", err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "update doc"})
		return
	}

	c.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}
