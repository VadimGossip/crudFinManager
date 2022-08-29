package rest

import (
	"context"
	_ "github.com/VadimGossip/crudFinManager/docs"
	"github.com/VadimGossip/crudFinManager/internal/domain"
	"github.com/gin-gonic/gin"
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

type Users interface {
	SignUp(ctx context.Context, inp domain.SignUpInput) error
	SignIn(ctx context.Context, inp domain.SignInInput) (string, error)
}

type Handler struct {
	usersService Users
	docsService  Docs
}

func NewHandler(users Users, docs Docs) *Handler {
	return &Handler{
		usersService: users,
		docsService:  docs,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	usersApi := router.Group("/auth")
	{
		usersApi.POST("/sign-up", h.signUp)
		usersApi.GET("/sign-in", h.signIn)
	}

	docsApi := router.Group("/docs")
	{
		docsApi.POST("", h.createDoc)
		docsApi.GET("/:id", h.getDocByID)
		docsApi.GET("", h.getAllDocs)
		docsApi.DELETE("/:id", h.deleteDocByID)
		docsApi.PUT("/:id", h.updateDocByID)
	}

	return router
}
