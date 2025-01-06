package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	User    *UserRepository
	Session *SessionRepository
	Group   *GroupRepository
	Member  *MemberRepository
}

func NewRepositories(pool *pgxpool.Pool) *Repository {
	return &Repository{
		User:    NewUserRepository(pool),
		Session: NewSessionRepository(pool),
		Group:   NewGroupRepository(pool),
		Member:  NewMemberRepository(pool),
	}
}
