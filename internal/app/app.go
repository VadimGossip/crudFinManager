package app

import (
	"context"
	"github.com/VadimGossip/crudFinManager/internal/config"
	"github.com/VadimGossip/crudFinManager/internal/repository/psql"
	"github.com/VadimGossip/crudFinManager/internal/server/http"
	"github.com/VadimGossip/crudFinManager/internal/service"
	"github.com/VadimGossip/crudFinManager/internal/transport/grps"
	"github.com/VadimGossip/crudFinManager/internal/transport/rest"
	"github.com/VadimGossip/crudFinManager/pkg/database"
	"github.com/VadimGossip/crudFinManager/pkg/hash"
	"github.com/VadimGossip/simpleCache"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
)

// @title Fin Manager App API
// @version 1.0
// @description API Server for Fin Manager Application

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func Run(configDir string) {
	cfg, err := config.Init(configDir)
	if err != nil {
		logrus.Fatalf("Config initialization error %s", err)
	}
	db, err := database.NewPostgresConnection(cfg.Postgres.Host, cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.Name, cfg.Postgres.SSLMode, cfg.Postgres.Port)
	if err != nil {
		logrus.Fatalf("Postgres connection error %s", err)
	}

	usersRepo := psql.NewUsers(db)
	tokensRepo := psql.NewTokens(db)
	hasher := hash.NewSHA1Hasher(cfg.Auth.Salt)
	auditClient, err := grps.NewClient(cfg.AuditServer.Host, cfg.AuditServer.Port)
	if err != nil {
		logrus.Fatalf("error occured while creating client for audit server %s", err)
	}

	usersService := service.NewUsers(usersRepo, tokensRepo, auditClient, hasher, []byte(cfg.Auth.Secret), cfg.Auth.AccessTokenTTL, cfg.Auth.RefreshTokenTTL)

	docsRepo := psql.NewDocs(db)
	cache := simpleCache.NewCache()
	docsService := service.NewBooks(docsRepo, cache)

	handler := rest.NewHandler(usersService, docsService)
	server := http.NewServer()

	go func() {
		if err := server.Run(cfg.Server, handler.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running rest server: %s", err.Error())
		}
	}()

	logrus.Info("Http Server for fin manager service started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Info("Http Server for fin manager service stopped")

	if err := db.Close(); err != nil {
		logrus.Errorf("Error occured on postgres connection close: %s", err.Error())
	}

	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Error occured on http server for fin manager service shutting down: %s", err.Error())
	}
}
