package model

import "time"

type User struct {
	Username  string
	Password  string
	Isadmin   bool
	CreatedAt time.Time
	UserId    string
}

type Token struct {
	Token     string
	CreatedAt time.Time
	UserId    string
}
