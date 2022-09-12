package service

import (
	"context"
	"github.com/ashah360/nyte-auth/internal/api/cerror"
	"github.com/ashah360/nyte-auth/internal/api/model"
	"github.com/ashah360/nyte-auth/internal/api/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, email string, password string) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func (us *userService) CreateUser(ctx context.Context, email, password string) error {
	ue, _ := us.userRepo.GetUserByEmail(ctx, email)
	if ue != nil {
		return cerror.ErrUserAlreadyExists
	}

	p, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u := model.NewUserWithDefaults()

	u.Email = email
	u.Password = string(p)

	return us.userRepo.CreateUser(ctx, u)
}

func (us *userService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return us.userRepo.GetUserByID(ctx, id)
}

func (us *userService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return us.userRepo.GetUserByEmail(ctx, email)
}

func NewUserService(urepo repository.UserRepository) UserService {
	return &userService{
		urepo,
	}
}
