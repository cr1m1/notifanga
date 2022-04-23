package main

import (
	tele "gopkg.in/telebot.v3"
)

type Bot struct {
	bot     *tele.Bot
	service *NotifangaService
	token   string
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
	return bot, nil
}

func (b *Bot) Start() {
	b.bot.Start()
}

func (b *Bot) Stop() {
	b.bot.Stop()
}

// func (b *Bot) TelegramBotReplier(service *NotifangaService) {
// 	// update := tgbotapi.NewUpdate(0)
// 	// update.Timeout = 60

// 	updates := b.bot.ListenForWebhook("/" + b.bot.Token)
// 	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)
// 	for u := range updates {
// 		fmt.Println(u.Message.Text)
// 		if u.Message == nil {
// 			continue
// 		}

// 		if reflect.TypeOf(u.Message.Text).Kind() == reflect.String && u.Message.Text != "" {
// 			message := strings.Fields(u.Message.Text)
// 			switch message[0] {
// 			case "/start":
// 				user := &User{
// 					TelegramUserID: u.Message.Chat.ID,
// 				}
// 				_, err := service.CreateUser(user)
// 				if err != nil {
// 					msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Ошибка в базе.Попробуйте снова.")
// 					b.bot.Send(msg)
// 				} else {
// 					msg := tgbotapi.NewMessage(u.Message.Chat.ID, startMsg)
// 					b.bot.Send(msg)
// 				}
// 			case "/list":
// 				user := &User{
// 					TelegramUserID: u.Message.Chat.ID,
// 				}
// 				user, err := service.CreateUser(user)
// 				if err != nil {
// 					msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Ошибка в базе.Попробуйте снова.")
// 					b.bot.Send(msg)
// 				}
// 				list, err := service.ListUserMangas(*user)
// 				if err != nil {
// 					msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Ошибка в базе.Попробуйте снова.")
// 					b.bot.Send(msg)
// 				}
// 				str := ""
// 				for i, m := range list {
// 					str += m.Name + " - " + strconv.Itoa(i) + "\n"
// 				}
// 				msg := tgbotapi.NewMessage(u.Message.Chat.ID, str)
// 				b.bot.Send(msg)
// 			case "/remove":
// 				if len(message) > 1 {
// 					user := &User{
// 						TelegramUserID: u.Message.Chat.ID,
// 					}
// 					user, err := service.CreateUser(user)
// 					if err != nil {
// 						msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Ошибка в базе.Попробуйте снова.")
// 						b.bot.Send(msg)
// 					}
// 					list, err := service.ListUserMangas(*user)
// 					if err != nil {
// 						msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Ошибка в базе.Попробуйте снова.")
// 						b.bot.Send(msg)
// 					}
// 					arg := message[1]
// 					num, err := strconv.Atoi(arg)
// 					if err != nil {
// 						msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Неправильный ввод.")
// 						b.bot.Send(msg)
// 					}
// 					if num < len(list) {
// 						m := list[num]
// 						if err := service.RemoveMangaFromUser(m, user); err != nil {
// 							msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Не удалось удалить. Попробуйте снова.")
// 							b.bot.Send(msg)
// 						} else {
// 							msg := tgbotapi.NewMessage(u.Message.Chat.ID, m.Name+" удален из вашего списка.")
// 							b.bot.Send(msg)
// 						}
// 					} else {
// 						msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Нет такого id в вашем списке.")
// 						b.bot.Send(msg)
// 					}
// 				}
// 			default:
// 				if strings.Contains(u.Message.Text, "mangalib.me") {
// 					link := u.Message.Text
// 					user := &User{
// 						TelegramUserID: u.Message.Chat.ID,
// 					}
// 					user, err := service.CreateUser(user)
// 					if err != nil {
// 						msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Ошибка в базе.Попробуйте снова.")
// 						b.bot.Send(msg)
// 					}
// 					manga := &Manga{
// 						Name:           CrawlName(link),
// 						Url:            link,
// 						LastChapter:    "",
// 						LastChapterUrl: "",
// 					}
// 					manga, err = service.CreateManga(manga)
// 					if err != nil {
// 						msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Манга не была найдена.")
// 						b.bot.Send(msg)
// 					}
// 					if err := service.AddMangaToUser(manga, user); err != nil {
// 						msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Манга не была найдена.")
// 						b.bot.Send(msg)
// 					} else {
// 						msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Манга добавлена в ваш список.")
// 						b.bot.Send(msg)
// 					}
// 				} else {
// 					msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Неправильная ссылка.")
// 					b.bot.Send(msg)
// 				}
// 			}
// 		}
// 	}
// }

// func (b *Bot) TelegramBotCrawler(service *NotifangaService) {
// 	for {
// 		marr, _ := service.GetAllMangas()

// 		for _, m := range marr {
// 			uarr, m := Crawl(*m, service)
// 			for _, u := range uarr {
// 				msg := tgbotapi.NewMessage(
// 					u.TelegramUserID,
// 					"Вышла новая "+m.LastChapter+" глава манги "+m.Name+"!\nЧитать тут - "+m.LastChapterUrl,
// 				)
// 				b.bot.Send(msg)
// 			}
// 		}
// 		time.Sleep(time.Minute * 5)
// 	}
// }
