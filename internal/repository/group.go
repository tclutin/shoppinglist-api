package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tclutin/shoppinglist-api/internal/domain/group"
)

type GroupRepository struct {
	db *pgxpool.Pool
}

func NewGroupRepository(db *pgxpool.Pool) *GroupRepository {
	return &GroupRepository{db: db}
}

func (g *GroupRepository) Create(ctx context.Context, group group.Group) (uint64, error) {
	sql := `INSERT INTO public.groups (name, description, code, created_at)
			VALUES ($1, $2, $3, $4)
			RETURNING group_id`

	row := g.db.QueryRow(ctx, sql, group.Name, group.Description, group.Code, group.CreatedAt)

	var groupID uint64
	if err := row.Scan(&groupID); err != nil {
		return 0, err
	}

	return groupID, nil
}

func (g *GroupRepository) Delete(ctx context.Context, groupID uint64) error {
	sql := `DELETE FROM public.groups WHERE group_id = $1`

	_, err := g.db.Exec(ctx, sql, groupID)

	return err
}

func (g *GroupRepository) GetByCode(ctx context.Context, code string) (group.Group, error) {
	sql := `SELECT * FROM public.groups WHERE code = $1`

	row := g.db.QueryRow(ctx, sql, code)

	var group group.Group
	err := row.Scan(
		&group.GroupID,
		&group.Name,
		&group.Description,
		&group.Code,
		&group.CreatedAt)

	if err != nil {
		return group, err
	}

	return group, nil
}

func (g *GroupRepository) GetById(ctx context.Context, groupID uint64) (group.Group, error) {
	sql := `SELECT * FROM public.groups WHERE group_id = $1`

	row := g.db.QueryRow(ctx, sql, groupID)

	var group group.Group
	err := row.Scan(
		&group.GroupID,
		&group.Name,
		&group.Description,
		&group.Code,
		&group.CreatedAt)

	if err != nil {
		return group, err
	}

	return group, nil
}
