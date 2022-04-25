package main

import (
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
	b.SetCommands()
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
}

func (b *Bot) CrawlerBot() {
	for {
		marr, _ := b.service.GetAllMangas()

		for _, m := range marr {
			uarr, m := Crawl(*m, b.service)
			for _, u := range uarr {
				r := NewRecipient(strconv.Itoa(int(u.TelegramUserID)))
				b.bot.Send(
					r,
					"Вышла новая ",
					m.LastChapter,
					" глава манги ",
					m.Name,
					"!\nЧитать тут - ",
					m.LastChapterUrl,
				)
			}
		}
		time.Sleep(time.Minute * 5)
	}
}
