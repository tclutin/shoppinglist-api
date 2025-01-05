package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	User    *UserRepository
	Session *SessionRepository
}

func NewRepositories(pool *pgxpool.Pool) *Repository {
	return &Repository{
		User:    NewUserRepository(pool),
		Session: NewSessionRepository(pool),
	}
}
