package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type notifangaRepository struct {
	conn *sql.DB
}

func NewRepository(db *sql.DB) (*notifangaRepository, error) {
	return &notifangaRepository{
		conn: db,
	}, nil
}

// checks does user exists or not
// if exists returns user's id
// if not exists creates user and returns user's id
func (r *notifangaRepository) UserCreate(u *User) (*User, error) {
	err := r.conn.QueryRow(`
		WITH s AS (
    		SELECT id, telegram_user_id
    		FROM users
    		WHERE telegram_user_id = $1
		), i AS (
    		INSERT INTO users (telegram_user_id)
    		SELECT $1
    		WHERE NOT EXISTS (SELECT 1 FROM s)
    		RETURNING id
		)
		SELECT id
		FROM i
		UNION ALL
		SELECT id
		FROM s
	`, u.TelegramUserID).Scan(&u.ID)
	return u, err
}

func (r *notifangaRepository) UserList() ([]*User, error) {
	rows, err := r.conn.Query(`
		SELECT id, telegram_user_id
		FROM users;
	`)
	if err != nil {
		log.Fatal(err)
	}
	var (
		id       int
		tgUserId string
		users    []*User
	)

	for rows.Next() {
		if err := rows.Scan(
			&id,
			&tgUserId,
		); err != nil {
			return nil, err
		}
		users = append(users, &User{
			ID:             id,
			TelegramUserID: tgUserId,
		})
	}
	return users, err
}

// checks does manga exists or not
// if exists returns manga's id
// if not exists creates manga and returns manga's id
func (r *notifangaRepository) MangaCreate(m *Manga) (*Manga, error) {
	err := r.conn.QueryRow(`
		WITH s AS (
    		SELECT id, name
    		FROM mangas
    		WHERE name = $1
		), i AS (
    		INSERT INTO mangas (name, link, last_chapter, last_chapter_url)
    		SELECT $1, $2, $3, $4
    		WHERE NOT EXISTS (SELECT 1 FROM s)
    		RETURNING id
		)
		SELECT id
		FROM i
		UNION ALL
		SELECT id
		FROM s
	`, m.Name, m.Url, m.LastChapter, m.LastChapterUrl).Scan(&m.ID)
	return m, err
}

func (r *notifangaRepository) MangaList(u User) (map[int]*Manga, error) {
	rows, err := r.conn.Query(`
		SELECT id, name, link, last_chapter, last_chapter_url
		FROM users_mangas
		INNER JOIN mangas
			ON manga_id = id
		WHERE user_id = $1;
	`, u.ID)
	if err != nil {
		log.Fatal(err)
	}
	var (
		id                                      int
		name, link, lastChapter, lastChapterUrl string
	)
	mangas := make(map[int]*Manga)

	for rows.Next() {
		if err := rows.Scan(
			&id,
			&name,
			&link,
			&lastChapter,
			&lastChapterUrl,
		); err != nil {
			return nil, err
		}
		mangas[id] = &Manga{
			Name:           name,
			Url:            link,
			LastChapter:    lastChapter,
			LastChapterUrl: lastChapterUrl,
		}
	}
	return mangas, err
}

func (r *notifangaRepository) UserListByManga(m Manga) ([]*User, error) {
	rows, err := r.conn.Query(`
		SELECT id, telegram_user_id
		FROM users_mangas
		INNER JOIN users
			ON user_id = id
		WHERE manga_id = $1;
	`, m.ID)
	if err != nil {
		log.Fatal(err)
	}
	var (
		id       int
		tgUserId string
		uarr     []*User
	)

	for rows.Next() {
		if err := rows.Scan(
			&id,
			&tgUserId,
		); err != nil {
			return nil, err
		}
		uarr = append(uarr, &User{
			ID:             id,
			TelegramUserID: tgUserId,
		})
	}
	return uarr, err
}

func (r *notifangaRepository) UpdateManga(m Manga) error {
	err := r.conn.QueryRow(`
		UPDATE mangas
		SET last_chapter = $1,
			last_chapter_url = $2
		WHERE id = $3
	`, m.LastChapter, m.LastChapterUrl, m.ID).Err()
	return err
}

func (r *notifangaRepository) AddMangaToUser(m *Manga, u *User) error {
	err := r.conn.QueryRow(`
		INSERT INTO users_mangas (user_id, manga_id)
		VALUES ($1, $2);
	`, m.ID, u.ID).Err()
	return err
}

func (r *notifangaRepository) DeleteMangaFromUser(m Manga, u User) error {
	err := r.conn.QueryRow(`
		DELETE 
		FROM users_mangas
		WHERE manga_id=$1
		AND
		user_id=$2;
	`, m.ID, u.ID).Err()
	return err
}
