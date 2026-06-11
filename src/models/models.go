package models

import "time"

type User struct {
	Id           int
	UserId       string
	Name         string
	Band         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

type Session struct {
	Id        int
	UsersId   string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type Message struct {
	ID int
}
