package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if len(os.Args) > 1 {
		if err := godotenv.Load(os.Args[1]); err != nil {
			log.Println(err)
		}
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

	http.ListenAndServe(":8080", nil)

	go TelegramBotReplier(service)
	TelegramBotCrawler(service)

}
