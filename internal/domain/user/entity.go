package user

import "time"

type User struct {
	UserID    uint64
	Username  string
	Password  string
	Gender    string
	CreatedAt time.Time
}
