package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println(err)
	}
	dbconn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println("cannot get connection with db", err)
	}
	defer dbconn.Close()
	repo, err := NewRepository(dbconn)
	if err != nil {
		log.Println("cannot create repository", err)
	}
	service := NewNotifangaService(repo)
	marr, err := service.GetAllMangas()
	if err != nil {
		log.Println("cannot get all mangas", err)
	}

	for _, m := range marr {
		uarr := Crawl(*m, service)
		for _, u := range uarr {
			fmt.Println("new chapter!", u.TelegramUserID)
		}
	}
}
