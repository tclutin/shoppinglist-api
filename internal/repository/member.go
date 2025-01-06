package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tclutin/shoppinglist-api/internal/domain/member"
)

type MemberRepository struct {
	db *pgxpool.Pool
}

func NewMemberRepository(db *pgxpool.Pool) *MemberRepository {
	return &MemberRepository{db: db}
}

func (m *MemberRepository) Create(ctx context.Context, member member.Member) (uint64, error) {
	sql := `INSERT INTO public.members (user_id, group_id, role, joined_at)
			VALUES ($1, $2, $3, $4)
			RETURNING member_id`

	var memberID uint64
	row := m.db.QueryRow(ctx, sql, member.UserID, member.GroupID, member.Role, member.JoinedAt)

	if err := row.Scan(&memberID); err != nil {
		return 0, err
	}

	return memberID, nil
}

func (m *MemberRepository) Delete(ctx context.Context, memberID uint64) error {
	sql := `DELETE FROM public.members WHERE member_id = $1`

	_, err := m.db.Exec(ctx, sql, memberID)

	return err
}

func (m *MemberRepository) GetByUserId(ctx context.Context, userID uint64) (member.Member, error) {
	sql := `SELECT * FROM public.members WHERE user_id = $1`

	row := m.db.QueryRow(ctx, sql, userID)

	var member member.Member
	err := row.Scan(
		&member.MemberID,
		&member.UserID,
		&member.GroupID,
		&member.Role,
		&member.JoinedAt)

	if err != nil {
		return member, err
	}

	return member, nil
}

func (m *MemberRepository) GetByUserAndGroupId(ctx context.Context, userID uint64, groupID uint64) (member.Member, error) {
	sql := `SELECT * FROM public.members WHERE user_id = $1 AND group_id = $2`

	row := m.db.QueryRow(ctx, sql, userID, groupID)

	var member member.Member
	err := row.Scan(
		&member.MemberID,
		&member.UserID,
		&member.GroupID,
		&member.Role,
		&member.JoinedAt)

	if err != nil {
		return member, err
	}

	return member, nil
}

func (m *MemberRepository) GetMembersByGroupId(ctx context.Context, groupId uint64) ([]member.MemberDTO, error) {
	sql := `SELECT m.member_id, u.username, u.gender, m.role FROM public.members as m
			INNER JOIN public.users as u ON u.user_id = m.user_id
			WHERE m.group_id = $1`

	rows, err := m.db.Query(ctx, sql, groupId)
	if err != nil {
		return nil, err
	}

	members, err := pgx.CollectRows(rows, pgx.RowToStructByName[member.MemberDTO])
	if err != nil {
		return nil, err
	}

	return members, nil
}
