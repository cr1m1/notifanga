package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if len(os.Args) > 1 {
		if err := godotenv.Load(os.Args[1]); err != nil {
			log.Println(err)
		}
	}
	dbconn, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/notifanga?sslmode=disable")
	if err != nil {
		log.Fatal("cannot get connection with db", err)
	}
	defer dbconn.Close()
	repo, err := NewRepository(dbconn)
	if err != nil {
		log.Fatal("cannot create repository", err)
	}
	service := NewNotifangaService(repo)
	b, err := NewBot("5395533657:AAHt8UeoVmtpSE5yb0MdL32WewlXyn67Vv4", service)
	if err != nil {
		log.Fatal("cannot create bot")
	}
	// go log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
	go b.Start()
	b.CrawlerBot()
}
