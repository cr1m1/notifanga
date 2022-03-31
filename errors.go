package main

import "errors"

var (
	ErrNotValidUrl     = errors.New("invalid url")
	ErrNotValidChapter = errors.New("invalid chapter value")
	ErrUserExists      = errors.New("user already exists")
)
