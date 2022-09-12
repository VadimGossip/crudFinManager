package rest

import (
	"errors"
	"github.com/VadimGossip/crudFinManager/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce json
// @Param   input body     domain.SignUpInput          true "account info"
// @Success 201   {object} domain.StatusResponse
// @Failure 400   {object} domain.ErrorResponse
// @Failure 500   {object} domain.ErrorResponse
// @Router /auth/sign-up [post]
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

	c.JSON(http.StatusCreated, domain.StatusResponse{Status: "ok"})
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce json
// @Param   input body     domain.SignInInput          true "account credentials"
// @Success 201   {object} domain.StatusResponse
// @Failure 400   {object} domain.ErrorResponse
// @Failure 500   {object} domain.ErrorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input domain.SignInInput
	if err := c.BindJSON(&input); err != nil {
		logError("signIn", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "invalid user input param"})
		return
	}

	accessToken, refreshToken, err := h.usersService.SignIn(c.Request.Context(), input)
	if err != nil {
		logError("signIn", err)
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: domain.ErrUserNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "search user error"})
		return
	}
	refreshTokenTTL := h.usersService.GetRefreshTokenTTL().Seconds()
	c.SetCookie("refresh-token", refreshToken, int(refreshTokenTTL), "/", "localhost", false, true)
	c.JSON(http.StatusOK, domain.TokenResponse{Token: accessToken})
}

// @Summary     Refresh
// @Description Refresh tokens
// @Tags        auth
// @Produce     json
// @Success     200     {object} domain.TokenResponse
// @Failure     400,500 {object} domain.ErrorResponse
// @Router      /auth/refresh [get]
func (h *Handler) refresh(c *gin.Context) {
	cookie, err := c.Cookie("refresh-token")
	if err != nil {
		logError("refresh", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "parse refresh token error"})
		return
	}

	accessToken, refreshToken, err := h.usersService.RefreshTokens(c.Request.Context(), cookie)
	if err != nil {
		logError("refresh", err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "refresh tokens error"})
		return
	}
	refreshTokenTTL := h.usersService.GetRefreshTokenTTL().Seconds()
	c.SetCookie("refresh-token", refreshToken, int(refreshTokenTTL), "/", "localhost", false, true)
	c.JSON(http.StatusOK, domain.TokenResponse{Token: accessToken})
}
