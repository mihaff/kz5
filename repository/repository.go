package repository

import (
	"feklistova/initializr"
	"database/sql"
	"log"
)

type Repository struct {
	Db *sql.DB
}

func (r *Repository) NewRepository() {
	db, err := initializr.DbConnectionInit()
	if err != nil {
		log.Fatal(err)
	}
	r.Db = db
}
