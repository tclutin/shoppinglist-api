package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tclutin/shoppinglist-api/internal/domain/group"
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

func (u *UserRepository) GetGroupsByUserId(ctx context.Context, userId uint64) ([]group.GroupDTO, error) {
	sql := `SELECT g.group_id, g.name, g.description, g.code FROM public.members as m
			INNER JOIN public.groups as g ON g.group_id = m.group_id
			WHERE m.user_id = $1`

	rows, err := u.db.Query(ctx, sql, userId)
	if err != nil {
		return nil, err
	}

	groups, err := pgx.CollectRows(rows, pgx.RowToStructByName[group.GroupDTO])
	if err != nil {
		return nil, err
	}

	return groups, nil
}
