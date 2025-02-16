package group

import "time"

type Group struct {
	GroupID     uint64
	Name        string
	Description string
	Code        string
	CreatedAt   time.Time
}
