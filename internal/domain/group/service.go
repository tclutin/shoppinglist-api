package group

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	domainErr "github.com/tclutin/shoppinglist-api/internal/domain/errors"
	"github.com/tclutin/shoppinglist-api/internal/domain/member"
	"github.com/tclutin/shoppinglist-api/pkg/hash"
	"log/slog"
	"time"
)

type MemberRepository interface {
	Create(ctx context.Context, member member.Member) (uint64, error)
	Delete(ctx context.Context, memberID uint64) error
	GetByUserId(ctx context.Context, userID uint64) (member.Member, error)
	GetByUserAndGroupId(ctx context.Context, userID uint64, groupID uint64) (member.Member, error)
	GetMembersByGroupId(ctx context.Context, groupId uint64) ([]member.MemberDTO, error)
}

type Repository interface {
	Create(ctx context.Context, group Group) (uint64, error)
	Delete(ctx context.Context, groupID uint64) error

	GetById(ctx context.Context, groupID uint64) (Group, error)
	GetByCode(ctx context.Context, code string) (Group, error)
}

type Service struct {
	logger     *slog.Logger
	repo       Repository
	memberRepo MemberRepository
}

func NewService(repo Repository, memberRepo MemberRepository, logger *slog.Logger) *Service {
	return &Service{
		logger:     logger.With("service", "group_service"),
		repo:       repo,
		memberRepo: memberRepo,
	}
}

func (s *Service) CreateGroup(ctx context.Context, dto CreateGroupDTO) (uint64, error) {
	code, err := s.GenCode(5)
	if err != nil {
		return 0, err
	}

	group := Group{
		Name:        dto.Name,
		Description: dto.Description,
		Code:        code,
		CreatedAt:   time.Now().UTC(),
	}

	groupID, err := s.repo.Create(ctx, group)
	if err != nil {
		return 0, err
	}

	member := member.Member{
		UserID:   dto.OwnerID,
		GroupID:  groupID,
		Role:     "owner",
		JoinedAt: time.Now().UTC(),
	}

	_, err = s.memberRepo.Create(ctx, member)
	if err != nil {
		return 0, err
	}

	return groupID, nil
}

func (s *Service) DeleteGroup(ctx context.Context, dto GroupUserDTO) error {
	group, err := s.repo.GetById(ctx, dto.GroupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainErr.ErrGroupNotFound
		}
	}

	membr, err := s.memberRepo.GetByUserAndGroupId(ctx, dto.UserID, group.GroupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainErr.ErrMemberNotFound
		}
	}

	if membr.Role != "owner" {
		return domainErr.ErrAreNotOwner
	}

	return s.repo.Delete(ctx, group.GroupID)
}

func (s *Service) JoinToGroup(ctx context.Context, dto JoinToGroupDTO) error {
	group, err := s.repo.GetByCode(ctx, dto.Code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainErr.ErrInvalidCode
		}
	}

	_, err = s.memberRepo.GetByUserAndGroupId(ctx, dto.UserID, group.GroupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			member := member.Member{
				UserID:   dto.UserID,
				GroupID:  group.GroupID,
				Role:     "member",
				JoinedAt: time.Now().UTC(),
			}

			_, err = s.memberRepo.Create(ctx, member)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return domainErr.ErrAlreadyMember
}

func (s *Service) LeaveFromGroup(ctx context.Context, dto GroupUserDTO) error {
	group, err := s.repo.GetById(ctx, dto.GroupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainErr.ErrGroupNotFound
		}
	}

	membr, err := s.memberRepo.GetByUserAndGroupId(ctx, dto.UserID, group.GroupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainErr.ErrMemberNotFound
		}
	}

	if membr.Role == "owner" {
		return domainErr.ErrOwnerCannotLeave
	}

	return s.memberRepo.Delete(ctx, membr.MemberID)
}

func (s *Service) GenCode(size int64) (string, error) {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	alias := make([]rune, size)
	for i := range alias {
		rnd, err := hash.NewCryptoRand(int64(len(chars)))
		if err != nil {
			return "", err
		}
		alias[i] = chars[rnd]
	}
	return string(alias), nil
}
