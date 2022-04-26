package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"

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
		log.Fatal("cannot get connection with db", err)
	}
	defer dbconn.Close()
	repo, err := NewRepository(dbconn)
	if err != nil {
		log.Fatal("cannot create repository", err)
	}
	service := NewNotifangaService(repo)
	b, err := NewBot(os.Getenv("TOKEN"), service)
	if err != nil {
		log.Fatal("cannot create bot")
	}
	go func() {
		b.Start()
		defer b.Stop()
	}()
	go b.CrawlerBot()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("Gracefully closed")
}
