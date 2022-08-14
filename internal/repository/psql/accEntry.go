package psql

import "database/sql"

type AccEntries struct {
	db *sql.DB
}

func NewAccEntries(db *sql.DB) *AccEntries {
	return &AccEntries{db: db}
}
