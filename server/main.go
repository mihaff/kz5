package main

import (
	_ "feklistova/docs"
	"feklistova/filestorage"
	"feklistova/python"
	"feklistova/repository"
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("secret"))
var repo repository.Repository
var fileRepo filestorage.FileStorage
var pyModel python.PyModel

//	@title			Social Network API
//	@version		1.0
//	@description	API-документация для реализации простейшей соцсети

//	@contact.name	Ellina Aleshina
//	@contact.email	esalesina_1@edu.hse.ru

// @host		localhost:8080
// @BasePath	/home
func main() {
	defer log.Println("Shutting down completed")
	log.Println("Starting")

	log.Println("Setting up file storage")
	if err := fileRepo.NewFileStorage("/root/uploads", "/root/downloads"); err != nil {
		panic(err)
	}

	log.Println("Opening database connection")
	repo.NewRepository()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
		log.Println("Connection to database closed successfully")
	}(repo.Db)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigChannel := make(chan os.Signal, 1)
		signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)
		<-sigChannel
		close(sigChannel)
		cancel()
	}()
	Serve(ctx)
}
