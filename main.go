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

	// user1
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
		Name:           CrawlName("https://mangalib.me/onepunchman?section=info"),
		Url:            "https://mangalib.me/onepunchman?section=info",
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
	// if err := service.RemoveMangaFromUser(m, u); err != nil {
	// }

	// // user2
	// u = &User{
	// 	TelegramUserID: "telegram2",
	// }
	// u, err = service.CreateUser(u)
	// if err != nil {
	// 	log.Println("cant create a user", err)
	// }
	// log.Println("user was created", u)

	// m = &Manga{
	// 	Name:           CrawlName("https://mangalib.me/toukyou-revengers?section=info"),
	// 	Url:            "https://mangalib.me/toukyou-revengers?section=info",
	// 	LastChapter:    "",
	// 	LastChapterUrl: "",
	// }
	// log.Println("manga name", m.Name)
	// m, err = service.CreateManga(m)
	// if err != nil {
	// 	log.Println("cant create manga", err)
	// }
	// log.Println("manga was created", m)

	// if err := service.AddMangaToUser(m, u); err != nil {
	// 	log.Println("cant add manga to user", err)
	// }

	// for {
	marr, err := service.GetAllMangas()
	for _, manga := range marr {
		log.Println(manga)
	}
	if err != nil {
		log.Println("cannot get all mangas", err)
	}

	for _, m := range marr {
		log.Println("loop started")
		log.Println("manga list", m)
		uarr := Crawl(*m, service)
		for _, user := range uarr {
			log.Println("user list", user)
		}
		for _, u := range uarr {
			fmt.Println("new chapter!", u.TelegramUserID)
		}
	}
	// time.Sleep(time.Minute * 1)
	// }

	users, err := service.GetUsers()
	if err != nil {
		log.Println("cant get list of all users", err)
	}
	log.Println("users:")
	for _, u := range users {
		log.Println(u)
		mangas, err := service.ListUserMangas(*u)
		if err != nil {
		}
		for _, m := range mangas {
			log.Println(m)
		}
	}
}
