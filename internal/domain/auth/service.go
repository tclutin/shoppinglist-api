package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tclutin/shoppinglist-api/internal/config"
	domainErr "github.com/tclutin/shoppinglist-api/internal/domain/errors"
	"github.com/tclutin/shoppinglist-api/internal/domain/user"
	"github.com/tclutin/shoppinglist-api/pkg/hash"
	"github.com/tclutin/shoppinglist-api/pkg/jwt/manager"
	"time"
)

type UserService interface {
	GetById(ctx context.Context, userID uint64) (user.User, error)
	Create(ctx context.Context, user user.User) (uint64, error)
	GetByUsername(ctx context.Context, username string) (user.User, error)
}

type Repository interface {
	CreateSession(ctx context.Context, session Session) (uint64, error)
	GetSessionByRefreshToken(ctx context.Context, token uuid.UUID) (Session, error)
	DeleteSession(ctx context.Context, sessionID uint64) error
}

type Service struct {
	cfg          *config.Config
	tokenManager manager.Manager
	userService  UserService
	repo         Repository
}

func NewService(cfg *config.Config, userService UserService, tokenManager manager.Manager, repo Repository) *Service {
	return &Service{
		tokenManager: tokenManager,
		cfg:          cfg,
		userService:  userService,
		repo:         repo,
	}
}

func (s *Service) LogIn(ctx context.Context, dto LogInDTO) (TokenDTO, error) {
	usr, err := s.userService.GetByUsername(ctx, dto.Username)
	if err != nil {
		return TokenDTO{}, err
	}

	if !hash.CompareBcryptHash(usr.Password, dto.Password) {
		return TokenDTO{}, domainErr.ErrUserNotValid
	}

	accessToken, err := s.tokenManager.NewAccessToken(usr.UserID, s.cfg.JWT.AccessExpire)
	if err != nil {
		return TokenDTO{}, err
	}

	refreshToken := s.tokenManager.NewRefreshToken()

	session := Session{
		UserID:       usr.UserID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().UTC().Add(s.cfg.JWT.RefreshExpire),
		CreatedAt:    time.Now().UTC(),
	}

	_, err = s.repo.CreateSession(ctx, session)
	if err != nil {
		return TokenDTO{}, err
	}

	return TokenDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.String(),
	}, nil
}

func (s *Service) SignUp(ctx context.Context, dto SignUpDTO) (TokenDTO, error) {
	_, err := s.userService.GetByUsername(ctx, dto.Username)
	if err == nil {
		return TokenDTO{}, domainErr.ErrUserAlreadyExists
	}

	bcryptHash, err := hash.NewBcryptHash(dto.Password)
	if err != nil {
		return TokenDTO{}, fmt.Errorf("failed to get crypthash of password: %w", err)
	}

	entity := user.User{
		Username:  dto.Username,
		Password:  bcryptHash,
		Gender:    dto.Gender,
		CreatedAt: time.Now().UTC(),
	}

	userID, err := s.userService.Create(ctx, entity)
	if err != nil {
		return TokenDTO{}, err
	}

	accessToken, err := s.tokenManager.NewAccessToken(userID, s.cfg.JWT.AccessExpire)
	if err != nil {
		return TokenDTO{}, err
	}

	refreshToken := s.tokenManager.NewRefreshToken()

	session := Session{
		UserID:       userID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().UTC().Add(s.cfg.JWT.RefreshExpire),
		CreatedAt:    time.Now().UTC(),
	}

	_, err = s.repo.CreateSession(ctx, session)
	if err != nil {
		return TokenDTO{}, err
	}

	return TokenDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.String(),
	}, nil
}

func (s *Service) Refresh(ctx context.Context, dto RefreshTokenDTO) (TokenDTO, error) {
	session, err := s.repo.GetSessionByRefreshToken(ctx, dto.RefreshToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return TokenDTO{}, domainErr.ErrSessionNotFound
		}
	}

	if time.Now().UTC().After(session.ExpiresAt) {
		return TokenDTO{}, domainErr.ErrRefreshTokenExpired
	}

	accessToken, err := s.tokenManager.NewAccessToken(session.UserID, s.cfg.JWT.AccessExpire)
	if err != nil {
		return TokenDTO{}, err
	}

	refreshToken := s.tokenManager.NewRefreshToken()

	session.RefreshToken = refreshToken
	session.ExpiresAt = time.Now().UTC().Add(s.cfg.JWT.RefreshExpire)
	session.CreatedAt = time.Now().UTC()

	if err = s.repo.DeleteSession(ctx, session.SessionID); err != nil {
		return TokenDTO{}, err
	}

	_, err = s.repo.CreateSession(ctx, session)
	if err != nil {
		return TokenDTO{}, err
	}

	return TokenDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.String(),
	}, nil

}

func (s *Service) VerifyCredentials(accessToken string) (uint64, error) {
	return s.tokenManager.ParseToken(accessToken)
}

func (s *Service) Who(ctx context.Context, userID uint64) (user.User, error) {
	usr, err := s.userService.GetById(ctx, userID)
	if err != nil {
		return usr, err
	}

	return usr, nil
}
