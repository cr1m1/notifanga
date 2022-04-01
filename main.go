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

	u := &User{
		TelegramUserID: "telegram1",
	}
	u, err = service.CreateUser(u)
	if err != nil {
		log.Println("cant create a user", err)
	}
	log.Println("user was created", u)

	m := &Manga{
		Name:           CrawlName("https://mangalib.me/one-piece?section=info"),
		Url:            "https://mangalib.me/one-piece?section=info",
		LastChapter:    "",
		LastChapterUrl: "",
	}
	log.Println("manga name", m.Name)
	m, err = service.CreateManga(m)
	if err != nil {
		log.Println("cant create manga", err)
	}
	log.Println("manga was created", m)

	if err := service.AddMangaToUser(m, u); err != nil {
		log.Println("cant add manga to user", err)
	}

	m = &Manga{
		Name:           CrawlName("https://mangalib.me/toukyou-revengers?section=info"),
		Url:            "https://mangalib.me/toukyou-revengers?section=info",
		LastChapter:    "",
		LastChapterUrl: "",
	}
	log.Println("manga name", m.Name)
	m, err = service.CreateManga(m)
	if err != nil {
		log.Println("cant create manga", err)
	}
	log.Println("manga was created", m)

	if err := service.AddMangaToUser(m, u); err != nil {
		log.Println("cant add manga to user", err)
	}

	marr, err := service.GetAllMangas()
	if err != nil {
		log.Println("cannot get all mangas", err)
	}

	for _, m := range marr {
		log.Println("loop started")
		log.Println("manga list", m)
		uarr := Crawl(*m, service)
		for _, u := range uarr {
			fmt.Println("new chapter!", u.TelegramUserID)
		}
	}
}
