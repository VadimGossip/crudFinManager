package psql

import "database/sql"

type Docs struct {
	db *sql.DB
}

func NewDocs(db *sql.DB) *Docs {
	return &Docs{db: db}
}
