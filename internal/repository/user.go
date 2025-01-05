package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tclutin/shoppinglist-api/internal/domain/user"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) Create(ctx context.Context, user user.User) (uint64, error) {
	sql := `INSERT INTO public.users (username, password, gender, created_at)
			VALUES ($1, $2, $3, $4)
			RETURNING user_id`

	row := u.db.QueryRow(
		ctx,
		sql,
		user.Username,
		user.Password,
		user.Gender,
		user.CreatedAt)

	var userID uint64
	if err := row.Scan(&userID); err != nil {
		return 0, err
	}

	return userID, nil
}

func (u *UserRepository) GetById(ctx context.Context, userID uint64) (user.User, error) {
	sql := `SELECT * FROM public.users WHERE user_id = $1`

	row := u.db.QueryRow(ctx, sql, userID)

	var usr user.User
	err := row.Scan(
		&usr.UserID,
		&usr.Username,
		&usr.Password,
		&usr.Gender,
		&usr.CreatedAt)

	if err != nil {
		return usr, err
	}

	return usr, nil
}

func (u *UserRepository) GetByUsername(ctx context.Context, username string) (user.User, error) {
	sql := `SELECT * FROM public.users WHERE username = $1`

	row := u.db.QueryRow(ctx, sql, username)

	var usr user.User
	err := row.Scan(
		&usr.UserID,
		&usr.Username,
		&usr.Password,
		&usr.Gender,
		&usr.CreatedAt)

	if err != nil {
		return usr, err
	}

	return usr, nil
}
