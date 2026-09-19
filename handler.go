package main

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/telebot.v3"
)

const (
	startMsg   = "Привет!\nС помощью этого бота ты не пропустишь новые главы твоего любимого аниме!\nЕсли выйдет новая глава, то бот сам тебе об этом напишет. Для этого он использует сайт mangalib.me.\nКоманды:\n/list - вся манга, которую бот отслеживает для тебя. Знак • означает непрочитанные главы.\n/unread - только манга с непрочитанными главами.\n/read (id) - отметить прочитанной. Без id отмечает всё. id можно найти с помощью команды /list.\n/remove (id) - удалить их списка. id можно найти с помощью команды /list.\nДля добавления манги в коллекцию, просто отправь ссылку на нее. Пример: https://mangalib.me/onepunchman?section=info"
	dbErrMsg   = "Ошибка в базе.Попробуйте снова."
	unreadMark = "• "
)

func (b *Bot) StartBot() telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		fmt.Println("start")
		user := &User{
			TelegramUserID: ctx.Chat().ID,
		}
		_, err := b.service.CreateUser(user)
		if err != nil {
			return ctx.Send(dbErrMsg)
		} else {
			return ctx.Send(startMsg)
		}
	}
}

// resolves the chat's user and returns their subscriptions in stable /list order
func (b *Bot) userAndList(ctx telebot.Context) (*User, []*Subscription, error) {
	user, err := b.service.CreateUser(&User{
		TelegramUserID: ctx.Chat().ID,
	})
	if err != nil {
		return nil, nil, err
	}
	list, err := b.service.ListUserSubscriptions(*user)
	if err != nil {
		return nil, nil, err
	}
	return user, list, nil
}

// picks subscriptions by the indices given as command arguments, skipping invalid ones
func pickByIndices(list []*Subscription, args []string) []*Subscription {
	var picked []*Subscription
	for _, a := range args {
		n, err := strconv.Atoi(a)
		if err != nil {
			continue
		}
		if n >= 0 && n < len(list) {
			picked = append(picked, list[n])
		}
	}
	return picked
}

func joinNames(list []*Subscription) string {
	names := make([]string, 0, len(list))
	for _, s := range list {
		names = append(names, s.Manga.Name)
	}
	return strings.Join(names, ", ")
}

func (b *Bot) List() telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		fmt.Println("list")
		_, list, err := b.userAndList(ctx)
		if err != nil {
			return ctx.Send(dbErrMsg)
		}
		if len(list) == 0 {
			return ctx.Send("Список пуст. Отправь ссылку на mangalib.me, чтобы добавить мангу.")
		}
		str := ""
		for i, s := range list {
			if s.Unread() {
				str += unreadMark
			}
			str += s.Manga.Name + " - " + strconv.Itoa(i) + "\n"
		}
		return ctx.Send(str)
	}
}

func (b *Bot) Unread() telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		fmt.Println("unread")
		_, list, err := b.userAndList(ctx)
		if err != nil {
			return ctx.Send(dbErrMsg)
		}
		str := ""
		for i, s := range list {
			if !s.Unread() {
				continue
			}
			lastRead := s.LastReadChapter
			if lastRead == "" {
				lastRead = "—"
			}
			str += s.Manga.Name + " - " + strconv.Itoa(i) +
				"\n  прочитано: " + lastRead + ", последняя: " + s.Manga.LastChapter +
				"\n  " + s.Manga.LastChapterUrl + "\n"
		}
		if str == "" {
			return ctx.Send("Непрочитанных глав нет 🎉")
		}
		return ctx.Send(str + "\nПрочитал? Отметь: /read (id)")
	}
}

func (b *Bot) Read() telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		fmt.Println("read")
		user, list, err := b.userAndList(ctx)
		if err != nil {
			return ctx.Send(dbErrMsg)
		}
		if len(ctx.Args()) == 0 {
			if err := b.service.MarkAllRead(user); err != nil {
				return ctx.Send("Не удалось отметить. Попробуйте снова.")
			}
			return ctx.Send("Всё отмечено как прочитанное.")
		}
		readList := pickByIndices(list, ctx.Args())
		if len(readList) == 0 {
			return ctx.Send("Укажите номера манг из /list, которые хотите отметить прочитанными.")
		}
		for _, s := range readList {
			m := s.Manga
			if err := b.service.MarkRead(&m, user); err != nil {
				return ctx.Send("Не удалось отметить. Попробуйте снова.")
			}
		}
		if len(readList) == 1 {
			return ctx.Send(joinNames(readList) + " отмечена как прочитанная.")
		}
		return ctx.Send(joinNames(readList) + " отмечены как прочитанные.")
	}
}

func (b *Bot) Remove() telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		fmt.Println("remove")
		if len(ctx.Args()) == 0 {
			return ctx.Send("Укажите номера манг, которые хотите убрать из списка.")
		}
		user, list, err := b.userAndList(ctx)
		if err != nil {
			return ctx.Send(dbErrMsg)
		}
		removeList := pickByIndices(list, ctx.Args())
		for _, s := range removeList {
			m := s.Manga
			if err := b.service.RemoveMangaFromUser(&m, user); err != nil {
				return ctx.Send("Не удалось удалить. Попробуйте снова.")
			}
		}
		if len(removeList) == 1 {
			return ctx.Send(joinNames(removeList) + " был удален из списка.")
		}
		return ctx.Send(joinNames(removeList) + " были удалены из списка.")
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
			return ctx.Send(dbErrMsg)
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
