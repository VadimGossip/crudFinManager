package rest

import (
	"errors"
	"net/http"

	"github.com/VadimGossip/crudFinManager/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var input domain.SignUpInput
	if err := c.BindJSON(&input); err != nil {
		logError("signUp", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "invalid user input param"})
		return
	}

	err := h.usersService.SignUp(c.Request.Context(), input)
	if err != nil {
		logError("signUp", err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "can't create user"})
		return
	}

	c.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}

func (h *Handler) signIn(c *gin.Context) {
	var input domain.SignInInput
	if err := c.BindJSON(&input); err != nil {
		logError("signIn", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "invalid user input param"})
		return
	}

	token, err := h.usersService.SignIn(c.Request.Context(), input)
	if err != nil {
		logError("signIn", err)
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: domain.ErrUserNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "search user error"})
		return
	}

	c.JSON(http.StatusOK, domain.TokenResponse{AccessToken: token})
}
