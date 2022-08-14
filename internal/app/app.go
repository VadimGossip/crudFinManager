package app

import (
	"github.com/VadimGossip/crudFinManager/internal/config"
	"github.com/VadimGossip/crudFinManager/pkg/database"
	"github.com/sirupsen/logrus"
)

func Run(configDir string) {
	cfg, err := config.Init(configDir)
	if err != nil {
		logrus.Errorf("Config initialization error %s", err)
	}
	_, err = database.NewPostgresConnection(cfg.Postgres)
	if err != nil {
		logrus.Errorf("Postgres connection error %s", err)
	}
}
