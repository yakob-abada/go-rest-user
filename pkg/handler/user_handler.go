package handler

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yakob-abada/go-rest-user/pkg/common"
	"github.com/yakob-abada/go-rest-user/pkg/errorhandler"
	"github.com/yakob-abada/go-rest-user/pkg/logging"
	"github.com/yakob-abada/go-rest-user/pkg/model"
	"github.com/yakob-abada/go-rest-user/pkg/publisher"
	"github.com/yakob-abada/go-rest-user/pkg/repository"
	"github.com/yakob-abada/go-rest-user/pkg/security"
	"github.com/yakob-abada/go-rest-user/pkg/validator"
)

// UserHandler handles user-related requests
type UserHandler struct {
	repo         repository.UserRepository
	logger       logging.Logger
	validator    validator.UserValidator
	publisher    publisher.Publisher
	errorHandler errorhandler.ErrorHandler
	hasher       security.PasswordHasher
}

// NewUserHandler creates a new instance of UserHandler with dependencies
func NewUserHandler(
	repo repository.UserRepository, logger logging.Logger, validator validator.UserValidator,
	publisher publisher.Publisher, errorHandler errorhandler.ErrorHandler, Hasher security.PasswordHasher,
) *UserHandler {
	return &UserHandler{
		repo:         repo,
		logger:       logger,
		validator:    validator,
		publisher:    publisher,
		errorHandler: errorHandler,
		hasher:       Hasher,
	}
}

// GetUsers retrieves users with pagination and filtering
// @Summary Get users
// @Description Retrieve all users with optional filters (first_name, last_name, country)
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Results per page (default: 10)"
// @Param first_name query string false "Filter by first name"
// @Param last_name query string false "Filter by last name"
// @Param country query string false "Filter by country"
// @Success 200 {array} model.User
// @Failure 500 {object} map[string]string
// @Router /users [get]
func (h *UserHandler) GetUsers(c echo.Context) error {
	ctx := c.Request().Context()
	correlationID := common.GetCorrelationID(ctx)

	// Parse query parameters
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	sortBy := c.QueryParam("sort_by")
	order := c.QueryParam("order")

	// Collect filters
	filters := map[string]string{}
	if firstName := c.QueryParam("first_name"); firstName != "" {
		filters["first_name"] = firstName
	}
	if lastName := c.QueryParam("last_name"); lastName != "" {
		filters["last_name"] = lastName
	}
	if email := c.QueryParam("email"); email != "" {
		filters["email"] = email
	}
	if country := c.QueryParam("country"); country != "" {
		filters["country"] = country
	}

	// Fetch users from repository
	users, total, err := h.repo.GetUsers(ctx, page, limit, filters, sortBy, order)
	if err != nil {
		return h.errorHandler.HandleInternalServerError(ctx, c, "Failed to retrieve users", map[string]interface{}{
			"error":          err.Error(),
			"correlation_id": correlationID,
		})
	}

	// Remove passwords before sending response
	safeUsers := []map[string]interface{}{}
	for _, user := range users {
		safeUsers = append(safeUsers, map[string]interface{}{
			"id":         user.ID,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"nickname":   user.Nickname,
			"email":      user.Email,
			"country":    user.Country,
		})
	}

	// Log request
	h.logger.Info(ctx, "Fetched users with pagination and filters", map[string]interface{}{
		"page":           page,
		"limit":          limit,
		"total_users":    total,
		"filters":        filters,
		"sort_by":        sortBy,
		"order":          order,
		"correlation_id": correlationID,
	})

	// Return response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"page":        page,
		"limit":       limit,
		"total_users": total,
		"users":       safeUsers,
	})
}

// SaveUser creates a new user
// @Summary Create a new user
// @Description Registers a new user with hashed password
// @Tags Users
// @Accept json
// @Produce json
// @Param user body model.User true "User data"
// @Success 201 {object} model.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) SaveUser(c echo.Context) error {
	ctx := c.Request().Context()
	correlationID := common.GetCorrelationID(ctx)

	var user model.User
	if err := c.Bind(&user); err != nil {
		return h.errorHandler.HandleBadRequest(ctx, c, "Invalid request payload", map[string]interface{}{
			"endpoint":       "SaveUser",
			"error":          err.Error(),
			"correlation_id": correlationID,
		})
	}

	// Validate user data
	if err := h.validator.ValidateUser(&user); err != nil {
		return h.errorHandler.HandleBadRequest(ctx, c, "Validation failed", map[string]interface{}{
			"user_email":     user.Email,
			"error":          err.Error(),
			"correlation_id": correlationID,
		})
	}

	existingUser, err := h.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return h.errorHandler.HandleInternalServerError(ctx, c, "Database error", map[string]interface{}{
			"user_email":     user.Email,
			"error":          err.Error(),
			"correlation_id": correlationID,
		})
	}
	if existingUser != nil {
		h.logger.Warn(ctx, "User already exists", map[string]interface{}{
			"user_email":     user.Email,
			"correlation_id": correlationID,
		})
		return h.errorHandler.HandleBadRequest(ctx, c, "User with this email already exists", map[string]interface{}{
			"user_email":     user.Email,
			"correlation_id": correlationID,
		})
	}

	// Hash password
	hashedPassword, err := h.hasher.HashPassword(user.Password)
	if err != nil {
		return h.errorHandler.HandleInternalServerError(ctx, c, "Failed to hash password", map[string]interface{}{
			"user_email":     user.Email,
			"error":          err.Error(),
			"correlation_id": correlationID,
		})
	}
	user.Password = hashedPassword

	// Save the user
	if err := h.repo.SaveUser(ctx, &user); err != nil {
		return h.errorHandler.HandleInternalServerError(ctx, c, "Failed to save user", map[string]interface{}{
			"user_email":     user.Email,
			"error":          err.Error(),
			"correlation_id": correlationID,
		})
	}

	// Publish event
	event := "user.created"
	if err := h.publisher.Publish(event, map[string]interface{}{
		"user_id":        user.ID.String(),
		"user_email":     user.Email,
		"correlation_id": correlationID,
	}); err != nil {
		h.logger.Error(ctx, "Failed to publish event", map[string]interface{}{
			"user_id":        user.ID.String(),
			"event":          event,
			"error":          err.Error(),
			"correlation_id": correlationID,
		})
	}

	h.logger.Info(ctx, "User saved successfully", map[string]interface{}{
		"user_id":        user.ID.String(),
		"user_email":     user.Email,
		"event":          event,
		"correlation_id": correlationID,
	})

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"country":    user.Country,
	})
}

// UpdateUser updates user details
// @Summary Update a user
// @Description Update an existing user's details
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Param user body map[string]interface{} true "Updated user fields"
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	correlationID := common.GetCorrelationID(ctx)
	var updateRequest model.UpdateUserRequest

	if err := c.Bind(&updateRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Convert struct to map
	updateData := map[string]interface{}{}
	if updateRequest.FirstName != "" {
		updateData["first_name"] = updateRequest.FirstName
	}
	if updateRequest.LastName != "" {
		updateData["last_name"] = updateRequest.LastName
	}
	if updateRequest.Country != "" {
		updateData["country"] = updateRequest.Country
	}

	user, err := h.repo.UpdateUser(ctx, id, updateData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Publish event
	event := "user.updated"
	if err := h.publisher.Publish(event, map[string]interface{}{
		"user_id":        user.ID.String(),
		"user_email":     user.Email,
		"correlation_id": correlationID,
	}); err != nil {
		h.logger.Error(ctx, "Failed to publish event", map[string]interface{}{
			"user_id":        user.ID.String(),
			"event":          event,
			"error":          err.Error(),
			"correlation_id": correlationID,
		})
	}

	return c.JSON(http.StatusOK, user)
}

// DeleteUser removes a user by ID
// @Summary Delete a user
// @Description Deletes a user by ID
// @Tags Users
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()
	correlationID := common.GetCorrelationID(ctx)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return h.errorHandler.HandleBadRequest(ctx, c, "Invalid UUID format", map[string]interface{}{
			"endpoint":       "DeleteUser",
			"user_id":        c.Param("id"),
			"error":          err.Error(),
			"correlation_id": correlationID,
		})
	}

	// Delete user from db
	if err := h.repo.DeleteUser(ctx, id); err != nil {
		return h.errorHandler.HandleInternalServerError(ctx, c, "Failed to delete user", map[string]interface{}{
			"user_id":        id.String(),
			"error":          err.Error(),
			"correlation_id": correlationID,
		})
	}

	// Publish event
	event := "user.deleted"
	if err := h.publisher.Publish(event, map[string]interface{}{
		"user_id":        id.String(),
		"correlation_id": correlationID,
	}); err != nil {
		h.logger.Error(ctx, "Failed to publish event", map[string]interface{}{
			"user_id":        id.String(),
			"event":          event,
			"error":          err.Error(),
			"correlation_id": correlationID,
		})
	}

	h.logger.Info(ctx, "User deleted successfully", map[string]interface{}{
		"user_id":        id.String(),
		"event":          event,
		"correlation_id": correlationID,
	})

	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
