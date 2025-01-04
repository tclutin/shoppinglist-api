package user

import "time"

type Sex string

const (
	Male      Sex = "Male"
	Female    Sex = "Female"
	NonBinary Sex = "Non-Binary"
)

type User struct {
	UserID    uint64
	Username  string
	Password  string
	Sex       string
	Birthday  time.Time
	CreatedAt time.Time
}
