package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tclutin/shoppinglist-api/internal/domain/auth"
)

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{db: db}
}

func (s *SessionRepository) CreateSession(ctx context.Context, session auth.Session) (uint64, error) {
	sql := `INSERT INTO public.sessions (user_id, refresh_token, expires_at, created_at) VALUES ($1, $2, $3, $4) RETURNING session_id`

	row := s.db.QueryRow(
		ctx,
		sql,
		session.UserID,
		session.RefreshToken,
		session.ExpiresAt,
		session.CreatedAt)

	var sessionID uint64
	if err := row.Scan(&sessionID); err != nil {
		return 0, err
	}

	return sessionID, nil
}

func (s *SessionRepository) DeleteSession(ctx context.Context, sessionID uint64) error {
	sql := `DELETE FROM public.sessions WHERE session_id = $1`

	_, err := s.db.Exec(ctx, sql, sessionID)

	return err
}

func (s *SessionRepository) GetSessionByRefreshToken(ctx context.Context, token uuid.UUID) (auth.Session, error) {
	sql := `SELECT * FROM public.sessions WHERE refresh_token=$1;`

	row := s.db.QueryRow(ctx, sql, token)

	var session auth.Session
	err := row.Scan(
		&session.SessionID,
		&session.UserID,
		&session.RefreshToken,
		&session.ExpiresAt,
		&session.CreatedAt)

	if err != nil {
		return session, err
	}

	return session, nil
}
