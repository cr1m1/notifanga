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

// Subscription is a user's link to a manga together with reading progress.
type Subscription struct {
	Manga           Manga
	LastReadChapter string
}

// Unread reports whether the manga has a chapter the user has not marked as read yet.
func (s Subscription) Unread() bool {
	return s.Manga.LastChapter != "" && s.Manga.LastChapter != s.LastReadChapter
}
