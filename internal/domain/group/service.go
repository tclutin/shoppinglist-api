package group

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	domainErr "github.com/tclutin/shoppinglist-api/internal/domain/errors"
	"github.com/tclutin/shoppinglist-api/internal/domain/member"
	"github.com/tclutin/shoppinglist-api/internal/domain/product"
	"github.com/tclutin/shoppinglist-api/pkg/hash"
	"log/slog"
	"time"
)

type ProductService interface {
	Create(ctx context.Context, product product.Product) (uint64, error)
	Update(ctx context.Context, product product.Product) error
	GetByProductNameId(ctx context.Context, productNameID uint64) (product.ProductName, error)
	RemoveProduct(ctx context.Context, productID uint64) error
	GetById(ctx context.Context, productID uint64) (product.Product, error)
}

type MemberRepository interface {
	Create(ctx context.Context, member member.Member) (uint64, error)
	Delete(ctx context.Context, memberID uint64) error
	GetByUserId(ctx context.Context, userID uint64) (member.Member, error)
	GetByUserAndGroupId(ctx context.Context, userID uint64, groupID uint64) (member.Member, error)
	GetByMemberAndGroupId(ctx context.Context, memberID uint64, groupID uint64) (member.Member, error)
	GetMembersByGroupId(ctx context.Context, groupId uint64) ([]member.MemberDTO, error)
}

type Repository interface {
	Create(ctx context.Context, group Group) (uint64, error)
	Delete(ctx context.Context, groupID uint64) error
	GetById(ctx context.Context, groupID uint64) (Group, error)
	GetByCode(ctx context.Context, code string) (Group, error)
}

type Service struct {
	logger         *slog.Logger
	productService ProductService
	repo           Repository
	memberRepo     MemberRepository
}

func NewService(repo Repository, memberRepo MemberRepository, productService ProductService, logger *slog.Logger) *Service {
	return &Service{
		logger:         logger.With("service", "group_service"),
		productService: productService,
		repo:           repo,
		memberRepo:     memberRepo,
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

func (s *Service) GetGroupMembers(ctx context.Context, dto GroupUserDTO) ([]member.MemberDTO, error) {
	group, err := s.repo.GetById(ctx, dto.GroupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.ErrGroupNotFound
		}
	}

	membr, err := s.memberRepo.GetByUserAndGroupId(ctx, dto.UserID, group.GroupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.ErrMemberNotFound
		}
	}

	members, err := s.memberRepo.GetMembersByGroupId(ctx, membr.GroupID)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (s *Service) KickMember(ctx context.Context, dto KickMemberDTO) error {
	group, err := s.repo.GetById(ctx, dto.GroupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainErr.ErrGroupNotFound
		}
	}

	owner, err := s.memberRepo.GetByUserAndGroupId(ctx, dto.UserID, group.GroupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainErr.ErrMemberNotFound
		}
	}

	membr, err := s.memberRepo.GetByMemberAndGroupId(ctx, dto.MemberID, group.GroupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainErr.ErrMemberNotFound
		}
	}

	if owner.Role != "owner" {
		return domainErr.ErrAreNotOwner
	}

	if owner.UserID == membr.UserID {
		return domainErr.ErrCannotKickYourself
	}

	return s.memberRepo.Delete(ctx, membr.MemberID)
}

// TODO: needs to refactor
func (s *Service) AddProduct(ctx context.Context, dto CreateProductDTO) (uint64, error) {
	group, err := s.repo.GetById(ctx, dto.GroupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, domainErr.ErrGroupNotFound
		}
	}

	membr, err := s.memberRepo.GetByUserAndGroupId(ctx, dto.UserID, group.GroupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, domainErr.ErrMemberNotFound
		}
	}

	productName, err := s.productService.GetByProductNameId(ctx, dto.ProductNameID)
	if err != nil {
		return 0, err
	}

	product := product.Product{
		GroupID:       group.GroupID,
		ProductNameID: productName.ProductNameID,
		Price:         nil,
		Status:        "open",
		Quantity:      dto.Quantity,
		AddedBy:       membr.UserID,
		BoughtBy:      nil,
		CreatedAt:     time.Now().UTC(),
	}

	return s.productService.Create(ctx, product)
}

func (s *Service) RemoveProduct(ctx context.Context, dto RemoveProductDTO) error {
	group, err := s.repo.GetById(ctx, dto.GroupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainErr.ErrGroupNotFound
		}
	}

	_, err = s.memberRepo.GetByUserAndGroupId(ctx, dto.UserID, group.GroupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainErr.ErrMemberNotFound
		}
	}

	product, err := s.productService.GetById(ctx, dto.ProductID)
	if err != nil {
		return err
	}

	return s.productService.RemoveProduct(ctx, product.ProductID)
}

func (s *Service) UpdateProduct(ctx context.Context, dto UpdateProductDTO) error {
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

	product, err := s.productService.GetById(ctx, dto.ProductID)
	if err != nil {
		return err
	}

	product.Status = dto.Status
	product.Quantity = dto.Quantity
	product.Price = dto.Price
	product.BoughtBy = &membr.UserID

	return s.productService.Update(ctx, product)
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
