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
	b, err := NewBot(os.Getenv("TOKEN"), service)
	if err != nil {
		log.Println("cannot create bot")
	}
	go log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
	b.Start()
	b.CrawlerBot()
}
