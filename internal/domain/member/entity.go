package member

import "time"

type Member struct {
	MemberID uint64
	UserID   uint64
	GroupID  uint64
	Role     string
	JoinedAt time.Time
}
