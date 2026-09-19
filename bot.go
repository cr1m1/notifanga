package main

import (
	"log"
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Bot struct {
	bot     *tele.Bot
	service *NotifangaService
	token   string
}

type recipient struct {
	to string
}

func NewRecipient(id string) *recipient {
	return &recipient{
		to: id,
	}
}

func (r *recipient) Recipient() string {
	return r.to
}

func NewBot(t string, s *NotifangaService) (*Bot, error) {
	pref := tele.Settings{
		Token: t,
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}
	bot := &Bot{
		bot:     b,
		service: s,
		token:   t,
	}
	bot.SetCommands()
	return bot, nil
}

func (b *Bot) Start() {
	b.bot.Start()
}

func (b *Bot) Stop() {
	b.bot.Stop()
}

func (b *Bot) SetCommands() {
	b.bot.Handle(tele.OnText, b.Add())
	b.bot.Handle("/start", b.StartBot())
	b.bot.Handle("/list", b.List())
	b.bot.Handle("/remove", b.Remove())
	b.bot.Handle("/unread", b.Unread())
	b.bot.Handle("/read", b.Read())
	if err := b.bot.SetCommands([]tele.Command{
		{Text: "list", Description: "Вся отслеживаемая манга"},
		{Text: "unread", Description: "Манга с непрочитанными главами"},
		{Text: "read", Description: "Отметить прочитанной: /read (id)"},
		{Text: "remove", Description: "Удалить из списка: /remove (id)"},
	}); err != nil {
		log.Println("cannot set bot commands:", err)
	}
}

// NewChapterMessage builds the text sent to a subscriber when a new chapter is detected.
func NewChapterMessage(m Manga) string {
	return "Вышла новая " + m.LastChapter + " глава манги " + m.Name + "!\nЧитать тут - " + m.LastChapterUrl +
		"\nПрочитал? Отметь: /read (id), номер смотри в /list"
}

func (b *Bot) CrawlerBot() {
	for {
		marr, err := b.service.GetAllMangas()
		if err != nil {
			log.Println("cannot list mangas:", err)
		}

		for _, m := range marr {
			uarr, m := Crawl(*m, b.service)
			for _, u := range uarr {
				r := NewRecipient(strconv.FormatInt(u.TelegramUserID, 10))
				if _, err := b.bot.Send(r, NewChapterMessage(m)); err != nil {
					log.Println("cannot notify", r.Recipient(), err)
				}
			}
		}
		time.Sleep(time.Minute * 1)
	}
}
