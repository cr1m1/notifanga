package main

type User struct {
	ID             int
	TelegramUserID string
	Mangas         []int
}

type Manga struct {
	ID             int
	Name           string
	Url            string
	LastChapter    string
	LastChapterUrl string
}
