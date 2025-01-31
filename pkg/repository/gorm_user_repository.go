package repository

import (
	"context"
	"errors"
	"github.com/yakob-abada/go-rest-user/pkg/model"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// GormUserRepository implements UserRepository using GORM
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository creates a new instance
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

// GetUsers retrieves paginated users with filtering and sorting
func (repo *GormUserRepository) GetUsers(ctx context.Context, page, limit int, filters map[string]string, sortBy, order string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := repo.db.WithContext(ctx).Model(&model.User{})

	// Apply filters
	if val, exists := filters["first_name"]; exists {
		query = query.Where("first_name ILIKE ?", "%"+val+"%")
	}
	if val, exists := filters["last_name"]; exists {
		query = query.Where("last_name ILIKE ?", "%"+val+"%")
	}
	if val, exists := filters["email"]; exists {
		query = query.Where("email ILIKE ?", "%"+val+"%")
	}
	if val, exists := filters["country"]; exists {
		query = query.Where("country ILIKE ?", "%"+val+"%")
	}

	// Count total users
	if err := query.Count(&total).Error; err != nil {
		log.Error().Err(err).Msg("Database query failed - Count users")
		return nil, 0, err
	}

	// Sorting
	if sortBy != "" {
		if order == "desc" {
			query = query.Order(sortBy + " DESC")
		} else {
			query = query.Order(sortBy + " ASC")
		}
	} else {
		query = query.Order("id ASC")
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		log.Error().Err(err).Msg("Failed to fetch users from db")
		return nil, 0, err
	}

	log.Info().Int("page", page).Int("limit", limit).Int64("total_users", total).Msg("Users retrieved successfully")

	return users, total, nil
}

// GetUserByID fetches a single user by UUID
func (repo *GormUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := repo.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Str("user_id", id.String()).Msg("User not found")
			return nil, nil
		}
		log.Error().Err(err).Str("user_id", id.String()).Msg("Database error while fetching user")
		return nil, err
	}

	log.Info().Str("user_id", id.String()).Str("email", user.Email).Msg("User retrieved successfully")

	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (repo *GormUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := repo.db.WithContext(ctx).Where("email = ?", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // No user found
	}

	return &user, err
}

func (repo *GormUserRepository) UpdateUser(ctx context.Context, id string, updatedData map[string]interface{}) (*model.User, error) {
	var user model.User

	// Check if user exists
	if err := repo.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Update user fields
	if err := repo.db.WithContext(ctx).Model(&user).Updates(updatedData).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// SaveUser creates or updates a user
func (repo *GormUserRepository) SaveUser(ctx context.Context, user *model.User) error {
	if err := repo.db.WithContext(ctx).Save(user).Error; err != nil {
		log.Error().Err(err).Str("email", user.Email).Msg("Failed to save user")
		return err
	}

	log.Info().Str("user_id", user.ID.String()).Str("email", user.Email).Msg("User saved successfully")
	return nil
}

// DeleteUser deletes a user by UUID
func (repo *GormUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := repo.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error; err != nil {
		log.Error().Err(err).Str("user_id", id.String()).Msg("Failed to delete user")
		return err
	}

	log.Info().Str("user_id", id.String()).Msg("User deleted successfully")
	return nil
}
