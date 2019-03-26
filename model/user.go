package model

import "time"

type User struct {
	Username  string
	Password  string
	Isadmin   bool
	CreatedAt time.Time
	UserId    string
	Token     string
}
