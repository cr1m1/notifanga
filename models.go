package main

type User struct {
	ID             int
	TelegramUserID int64
	Mangas         []int
}

type Manga struct {
	ID             int
	Name           string
	Url            string
	LastChapter    string
	LastChapterUrl string
}
