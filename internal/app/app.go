package app

import (
	"github.com/VadimGossip/crudFinManager/internal/config"
	"github.com/VadimGossip/crudFinManager/pkg/database"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func Run(configDir string) {
	cfg, err := config.Init(configDir)
	if err != nil {
		logrus.Errorf("Config initialization error %s", err)
	}
	db, err := database.NewPostgresConnection(cfg.Postgres)
	if err != nil {
		logrus.Errorf("Postgres connection error %s", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := db.Close(); err != nil {
		logrus.Errorf("Error occured on postgres connection close: %s", err.Error())
	}
}
