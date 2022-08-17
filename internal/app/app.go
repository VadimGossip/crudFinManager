package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/VadimGossip/crudFinManager/internal/config"
	"github.com/VadimGossip/crudFinManager/internal/repository/psql"
	"github.com/VadimGossip/crudFinManager/internal/server/http"
	"github.com/VadimGossip/crudFinManager/internal/service"
	"github.com/VadimGossip/crudFinManager/internal/transport/rest"
	"github.com/VadimGossip/crudFinManager/pkg/database"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

// @title Fin Manager App API
// @version 1.0
// @description API Server for Fin Manager Application

// @host localhost:8080
// @BasePath /

func Run(configDir string) {
	cfg, err := config.Init(configDir)
	if err != nil {
		logrus.Fatalf("Config initialization error %s", err)
	}
	db, err := database.NewPostgresConnection(cfg.Postgres)
	if err != nil {
		logrus.Fatalf("Postgres connection error %s", err)
	}

	repo := psql.NewDocs(db)
	docsService := service.NewBooks(repo)
	handler := rest.NewHandler(docsService)
	server := http.NewServer()

	go func() {
		if err := server.Run(cfg.Server, handler.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running rest server: %s", err.Error())
		}
	}()

	logrus.Print("Http Server for fin manager service started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Http Server for fin manager service stopped")

	if err := db.Close(); err != nil {
		logrus.Errorf("Error occured on postgres connection close: %s", err.Error())
	}

	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Error occured on http server for fin manager service shutting down: %s", err.Error())
	}
}
