package initializr

import (
	"feklistova/config"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// DbConnectionInit инициализирует подключение к базе данных PostgreSQL
func DbConnectionInit() (*sql.DB, error) {
	db, err := sql.Open("postgres", config.ConnStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("Connection to database opened successfully")
	return db, nil
}
