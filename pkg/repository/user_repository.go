package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/yakob-abada/go-rest-user/pkg/model"
)

type UserRepository interface {
	GetUsers(ctx context.Context, page, limit int, filters map[string]string, sortBy, order string) ([]model.User, int64, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	SaveUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, id string, updatedData map[string]interface{}) (*model.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
