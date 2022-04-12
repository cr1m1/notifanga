package main

import (
	"database/sql"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	b, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}
	wh, _ := tgbotapi.NewWebhook("https://notifanga-bot.herokuapp.com/" + b.Token)
	_, err = b.Request(wh)
	if err != nil {
		log.Fatal(err)
	}
	info, err := b.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	bot := &Bot{
		bot: b,
	}

	// go log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))

	go bot.TelegramBotReplier(service)
	bot.TelegramBotCrawler(service)

}
