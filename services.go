package main

import (
	"strings"
)

type NotifangaRepository interface {
	UserCreate(u *User) (*User, error)
	UserList() ([]*User, error)
	MangaCreate(m *Manga) (*Manga, error)
	MangaList(u User) ([]*Manga, error)
	SubscriptionList(u User) ([]*Subscription, error)
	MarkRead(m *Manga, u *User, chapter string) error
	UserListByManga(m Manga) ([]*User, error)
	UpdateManga(m Manga) error
	AddMangaToUser(m *Manga, u *User) error
	DeleteMangaFromUser(m *Manga, u *User) error
}

type NotifangaService struct {
	repo NotifangaRepository
}

func NewNotifangaService(r NotifangaRepository) *NotifangaService {
	return &NotifangaService{
		repo: r,
	}
}

func (s *NotifangaService) CreateUser(u *User) (*User, error) {
	u, err := s.repo.UserCreate(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *NotifangaService) GetUsers() ([]*User, error) {
	uarr, err := s.repo.UserList()
	if err != nil {
		return nil, err
	}
	return uarr, err
}

func (s *NotifangaService) CreateManga(m *Manga) (*Manga, error) {
	if !strings.Contains(m.Url, "mangalib.me/") {
		return nil, ErrNotValidUrl
	}
	m, err := s.repo.MangaCreate(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// gets list of mangas of a user
func (s *NotifangaService) ListUserMangas(u User) ([]*Manga, error) {
	m, err := s.repo.MangaList(u)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// gets list of a user's subscriptions with reading progress
func (s *NotifangaService) ListUserSubscriptions(u User) ([]*Subscription, error) {
	subs, err := s.repo.SubscriptionList(u)
	if err != nil {
		return nil, err
	}
	return subs, nil
}

// gets only the subscriptions that have unread chapters
func (s *NotifangaService) ListUnread(u User) ([]*Subscription, error) {
	subs, err := s.repo.SubscriptionList(u)
	if err != nil {
		return nil, err
	}
	var unread []*Subscription
	for _, sub := range subs {
		if sub.Unread() {
			unread = append(unread, sub)
		}
	}
	return unread, nil
}

// marks manga as read up to its current last chapter
func (s *NotifangaService) MarkRead(m *Manga, u *User) error {
	return s.repo.MarkRead(m, u, m.LastChapter)
}

// marks every unread subscription of the user as read
func (s *NotifangaService) MarkAllRead(u *User) error {
	subs, err := s.repo.SubscriptionList(*u)
	if err != nil {
		return err
	}
	for _, sub := range subs {
		if !sub.Unread() {
			continue
		}
		m := sub.Manga
		if err := s.repo.MarkRead(&m, u, m.LastChapter); err != nil {
			return err
		}
	}
	return nil
}

// gets a list of users of a manga
func (s *NotifangaService) ListMangaUsers(m Manga) ([]*User, error) {
	u, err := s.repo.UserListByManga(m)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *NotifangaService) UpdateManga(m Manga) error {
	if err := s.repo.UpdateManga(m); err != nil {
		return err
	}
	return nil
}

func (s *NotifangaService) AddMangaToUser(m *Manga, u *User) error {
	if !strings.Contains(m.Url, "mangalib.me/") {
		return ErrNotValidUrl
	}
	if err := s.repo.AddMangaToUser(m, u); err != nil {
		return err
	}
	return nil
}

func (s *NotifangaService) RemoveMangaFromUser(m *Manga, u *User) error {
	if err := s.repo.DeleteMangaFromUser(m, u); err != nil {
		return err
	}
	return nil
}

func (s *NotifangaService) GetAllMangas() ([]*Manga, error) {
	uarr, err := s.repo.UserList()
	if err != nil {
		return nil, err
	}
	var marr []*Manga
	allKeys := make(map[int]bool)
	for _, u := range uarr {
		mList, err := s.repo.MangaList(*u)
		if err != nil {
			return nil, err
		}
		for _, m := range mList {
			if _, val := allKeys[m.ID]; !val {
				allKeys[m.ID] = true
				marr = append(marr, m)
			}
		}
	}
	return marr, nil
}
