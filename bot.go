package main

import (
	"os"
	"reflect"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	startMsg = "Привет!\nС помощью этого бота ты не пропустишь новые главы твоего любимого аниме!\nЕсли выйдет новая глава, то бот сам тебе об этом напишет. Для этого он использует сайт mangalib.me.\nКоманды:\n/list - вся манга, которую бот отслеживает для тебя.\n/remove (id) - удалить их списка. id можно найти с помощью команды /list.\nДля добавления манги в коллекцию, просто отправь ссылку на нее. Пример: https://mangalib.me/onepunchman?section=info"
)

func TelegramBot(service *NotifangaService) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}

	update := tgbotapi.NewUpdate(0)
	update.Timeout = 60

	updates, _ := bot.GetUpdates(update)

	for _, u := range updates {
		if u.Message == nil {
			continue
		}

		if reflect.TypeOf(u.Message.Text).Kind() == reflect.String && u.Message.Text != "" {
			switch u.Message.Text {
			case "/start":
				user := &User{
					TelegramUserID: u.Message.Chat.UserName,
				}
				_, err = service.CreateUser(user)
				if err != nil {
					msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Ошибка в базе.Попробуйте снова.")
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(u.Message.Chat.ID, startMsg)
					bot.Send(msg)
				}
			case "/list":
				user := &User{
					TelegramUserID: u.Message.Chat.UserName,
				}
				user, err = service.CreateUser(user)
				if err != nil {
					msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Ошибка в базе.Попробуйте снова.")
					bot.Send(msg)
				}
				list, err := service.ListUserMangas(*user)
				if err != nil {
					msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Ошибка в базе.Попробуйте снова.")
					bot.Send(msg)
				}
				str := ""
				for i, m := range list {
					str += m.Name + " - " + strconv.Itoa(i) + "\n"
				}
				msg := tgbotapi.NewMessage(u.Message.Chat.ID, str)
				bot.Send(msg)
			case "/remove":
				user := &User{
					TelegramUserID: u.Message.Chat.UserName,
				}
				user, err = service.CreateUser(user)
				if err != nil {
					msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Ошибка в базе.Попробуйте снова.")
					bot.Send(msg)
				}
				list, err := service.ListUserMangas(*user)
				if err != nil {
					msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Ошибка в базе.Попробуйте снова.")
					bot.Send(msg)
				}
				arg := u.Message.Text[8:]
				num, err := strconv.Atoi(arg)
				if err != nil {
					msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Неправильный ввод.")
					bot.Send(msg)
				}
				if num < len(list) {
					m := list[num]
					if err := service.RemoveMangaFromUser(m, user); err != nil {
						msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Не удалось удалить. Попробуйте снова.")
						bot.Send(msg)
					} else {
						msg := tgbotapi.NewMessage(u.Message.Chat.ID, m.Name+" удален из вашего списка.")
						bot.Send(msg)
					}
				} else {
					msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Нет такого id в вашем списке.")
					bot.Send(msg)
				}
			default:
				if strings.Contains(u.Message.Text, "mangalib.me") {
					link := u.Message.Text
					user := &User{
						TelegramUserID: u.Message.Chat.UserName,
					}
					user, err = service.CreateUser(user)
					if err != nil {
						msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Ошибка в базе.Попробуйте снова.")
						bot.Send(msg)
					}
					manga := &Manga{
						Name:           CrawlName(link),
						Url:            link,
						LastChapter:    "",
						LastChapterUrl: "",
					}
					manga, err = service.CreateManga(manga)
					if err != nil {
						msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Манга не была найдена.")
						bot.Send(msg)
					}
					if err := service.AddMangaToUser(manga, user); err != nil {
						msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Манга не была найдена.")
						bot.Send(msg)
					}
				} else {
					msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Неправильная ссылка.")
					bot.Send(msg)
				}
			}
		}
	}
}
