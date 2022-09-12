package service

import (
	"context"
	"os"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/ashah360/nyte-auth/internal/api/cache"
	"github.com/ashah360/nyte-auth/internal/api/cerror"
	"github.com/ashah360/nyte-auth/internal/api/model"
	"github.com/ashah360/nyte-auth/internal/api/permissions"
	"github.com/ashah360/nyte-auth/internal/api/repository"
	tok "github.com/ashah360/nyte-auth/internal/api/token"
)

type AuthService interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	AuthenticateUser(ctx context.Context, email, password string) (string, error)
	RefreshToken(ctx context.Context, userId, token string) (string, error)
	RefreshUserSnapshots(ctx context.Context, userId string) error

	GetUserSnapshot(ctx context.Context, userId, token string) (*model.UserSnapshot, error)
	VerifyJWT(ctx context.Context, token string) (*model.UserSnapshot, error)

	DeleteTokensByUserID(ctx context.Context, userId string) error
}

type authService struct {
	userRepo   repository.UserRepository
	tokenStore cache.TokenStore
}

// GetUserByID is for testing only.
func (s *authService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

// GetUserByEmail takes in an email and returns the associated user.
func (s *authService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.userRepo.GetUserByEmail(ctx, email)
}

func (s *authService) VerifyJWT(ctx context.Context, token string) (*model.UserSnapshot, error) {
	p, err := tok.ValidateJWT(token)
	if err != nil {
		return nil, err
	}

	return s.tokenStore.Get(ctx, p.ID, p.Token)
}

func (s *authService) AuthenticateUser(ctx context.Context, email, password string) (string, error) {
	// Fetch the User account via email
	u, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	// Compare the password
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", cerror.ErrInvalidAccountDetails
	}

	return s.getAccessToken(ctx, u)
}

func (s *authService) RefreshToken(ctx context.Context, userId, token string) (string, error) {
	// Check if token is valid (it's valid if it exists in the cache)
	if _, err := s.tokenStore.Get(ctx, userId, token); err != nil {
		return "", err
	}

	u, err := s.userRepo.GetUserByID(ctx, userId)
	if err != nil {
		return "", err
	}

	go s.tokenStore.Delete(ctx, userId, token)

	return s.getAccessToken(ctx, u)
}

func (s *authService) RefreshUserSnapshots(ctx context.Context, userId string) error {
	ttl, ok := ctx.Value("TOKEN_TTL_SECONDS").(time.Duration)
	if !ok || ttl == 0 {
		ttl = time.Hour * 24 * 7
	}

	u, err := s.userRepo.GetUserByID(ctx, userId)
	if err != nil {
		return err
	}

	snapshots, err := s.tokenStore.GetByUserID(ctx, userId)

	for _, us := range snapshots {
		us.Role = u.Role
		us.Upgraded = u.Upgraded
		us.AccessExpiresAt = u.AccessExpiresAt

		go s.tokenStore.Update(context.Background(), u.ID, us.Token, us)
	}

	return nil
}

func (s *authService) GetUserSnapshot(ctx context.Context, userId, token string) (*model.UserSnapshot, error) {
	return s.tokenStore.Get(ctx, userId, token)
}

func (s *authService) DeleteTokensByUserID(ctx context.Context, userId string) error {
	return s.tokenStore.DeleteByUserID(ctx, userId)
}

func (s *authService) getAccessToken(ctx context.Context, u *model.User) (string, error) {
	opaque, _ := uuid.NewRandom()

	us := &model.UserSnapshot{
		ID:              u.ID,
		Token:           opaque.String(),
		Role:            u.Role,
		Upgraded:        u.Upgraded,
		Grants:          []string{},
		AccessExpiresAt: u.AccessExpiresAt,
	}

	if us.AccessExpiresAt != nil && time.Until(*us.AccessExpiresAt) > 0 {
		us.Grants = append(us.Grants, permissions.AccessClient)
	}

	token := jwt.New(jwt.SigningMethodHS256)

	ttl, ok := ctx.Value("TOKEN_TTL_SECONDS").(time.Duration)
	if !ok || ttl == 0 {
		ttl = time.Hour * 24 * 7
	}

	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = u.ID
	claims["token"] = opaque.String()
	claims["exp"] = time.Now().Add(ttl).Unix()
	claims["iat"] = time.Now().Unix()

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	if err = s.tokenStore.Create(ctx, u.ID, opaque.String(), us, ttl); err != nil {
		return "", err
	}

	return t, nil
}

func NewAuthService(repo repository.UserRepository, tokenStore cache.TokenStore) AuthService {
	return &authService{
		repo,
		tokenStore,
	}
}
