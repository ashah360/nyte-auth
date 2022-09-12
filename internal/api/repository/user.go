package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/ashah360/nyte-auth/internal/api/cerror"
	"github.com/ashah360/nyte-auth/internal/api/model"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
}

type userRepository struct {
	db *sqlx.DB
}

// GetUserByID takes in a User ID and returns the associated User object.
func (r *userRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var u model.User

	if err := r.db.GetContext(ctx, &u, "select * from users where id=$1", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, cerror.ErrUserDoesNotExist
		}

		return nil, err
	}

	return &u, nil
}

// GetUserByID takes in a User ID and returns the associated User object.
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User

	if err := r.db.GetContext(ctx, &u, "select * from users where lower(email)=lower($1)", email); err != nil {
		if err == sql.ErrNoRows {
			return nil, cerror.ErrUserDoesNotExist
		}

		return nil, err
	}

	return &u, nil
}

// CreateUser inserts a new user into the database with a provided model.User
func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	q := fmt.Sprintf(`insert into %s (email, password, role, discord_id, upgraded, access_expires_at) 
					values (:email, :password, :role, :discord_id, :upgraded, :access_expires_at)`, "users")

	if _, err := r.db.NamedExecContext(ctx, q, user); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// NewUserRepository creates a new UserRepository with the provided sql db connection.
func NewUserRepository(conn *sqlx.DB) UserRepository {
	return &userRepository{
		db: conn,
	}
}
