package main

import (
	"database/sql"
)

type Store struct {
	db   *sql.DB
	repo *Repo
}

func New(url string) (*Store, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Store{
		db:   db,
		repo: NewRepository(db),
	}, nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) NotifangaRepo() *Repo {
	if s.repo != nil {
		return s.repo
	}

	s.repo = &Repo{
		conn: s.db,
	}
	return s.repo
}
