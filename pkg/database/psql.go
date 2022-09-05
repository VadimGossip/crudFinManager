package database

import (
	"database/sql"
	"fmt"
)

func NewPostgresConnection(host, username, password, name, sslmode string, port int) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		host, port, username, name, sslmode, password))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
