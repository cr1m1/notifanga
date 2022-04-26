package main

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/telebot.v3"
)

const (
	startMsg = "Привет!\nС помощью этого бота ты не пропустишь новые главы твоего любимого аниме!\nЕсли выйдет новая глава, то бот сам тебе об этом напишет. Для этого он использует сайт mangalib.me.\nКоманды:\n/list - вся манга, которую бот отслеживает для тебя.\n/remove (id) - удалить их списка. id можно найти с помощью команды /list.\nДля добавления манги в коллекцию, просто отправь ссылку на нее. Пример: https://mangalib.me/onepunchman?section=info"
)

func (b *Bot) StartBot() telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		fmt.Println("start")
		user := &User{
			TelegramUserID: ctx.Chat().ID,
		}
		_, err := b.service.CreateUser(user)
		if err != nil {
			return ctx.Send("Ошибка в базе.Попробуйте снова.")
		} else {
			return ctx.Send(startMsg)
		}
	}
}

func (b *Bot) List() telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		fmt.Println("list")
		user := &User{
			TelegramUserID: ctx.Chat().ID,
		}
		user, err := b.service.CreateUser(user)
		if err != nil {
			return ctx.Send("Ошибка в базе.Попробуйте снова.")
		}
		list, err := b.service.ListUserMangas(*user)
		if err != nil {
			return ctx.Send("Ошибка в базе.Попробуйте снова.")
		}
		str := ""
		for i, m := range list {
			str += m.Name + " - " + strconv.Itoa(i) + "\n"
		}
		return ctx.Send(str)
	}
}

func (b *Bot) Remove() telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		fmt.Println("remove")
		if len(ctx.Args()) == 0 {
			return ctx.Send("Укажите номера манг, которые хотите убрать из списка.")
		}
		user := &User{
			TelegramUserID: ctx.Chat().ID,
		}
		user, err := b.service.CreateUser(user)
		if err != nil {
			return ctx.Send("Ошибка в базе.Попробуйте снова.")
		}
		list, err := b.service.ListUserMangas(*user)
		if err != nil {
			return ctx.Send("Ошибка в базе.Попробуйте снова.")
		}
		args := ctx.Args()
		var removeList []*Manga
		for _, a := range args {
			n, err := strconv.Atoi(a)
			if err != nil {
				continue
			}
			if n < len(list) {
				removeList = append(removeList, list[n])
			}
		}
		msg := ""
		for i, m := range removeList {
			if err := b.service.RemoveMangaFromUser(m, user); err != nil {
				return ctx.Send("Не удалось удалить. Попробуйте снова.")
			} else {
				if i < len(removeList)-1 {
					msg += m.Name + ", "
				} else {
					msg += m.Name + " "
				}
			}
		}
		if len(removeList) == 1 {
			return ctx.Send(msg + "был удален из списка.")
		}
		return ctx.Send(msg + "были удалены из списка.")
	}
}

func (b *Bot) Add() telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		fmt.Println("add")
		link := ctx.Text()
		if !strings.Contains(link, "mangalib.me") {
			return ctx.Send("Неправильная ссылка.")
		}
		user := &User{
			TelegramUserID: ctx.Chat().ID,
		}
		user, err := b.service.CreateUser(user)
		if err != nil {
			return ctx.Send("Ошибка в базе.Попробуйте снова.")
		}
		manga := &Manga{
			Name:           CrawlName(link),
			Url:            link,
			LastChapter:    "",
			LastChapterUrl: "",
		}
		manga, err = b.service.CreateManga(manga)
		if err != nil {
			return ctx.Send("Манга не была найдена.")
		}
		if err = b.service.AddMangaToUser(manga, user); err != nil {
			return ctx.Send("Манга не была найдена.")
		}
		return ctx.Send("Манга была добавлена в список.")
	}
}
