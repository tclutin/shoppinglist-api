package auth

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	SessionID    uint64
	UserID       uint64
	RefreshToken uuid.UUID
	ExpiresAt    time.Time
	CreatedAt    time.Time
}
