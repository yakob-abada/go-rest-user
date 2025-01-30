package bootstrap

import (
	"github.com/yakob-abada/go-rest-user/pkg/errorhandler"
	"github.com/yakob-abada/go-rest-user/pkg/handler"
	"github.com/yakob-abada/go-rest-user/pkg/logging"
	"github.com/yakob-abada/go-rest-user/pkg/publisher"
	"github.com/yakob-abada/go-rest-user/pkg/repository"
	"github.com/yakob-abada/go-rest-user/pkg/security"
	"github.com/yakob-abada/go-rest-user/pkg/validator"
	"gorm.io/gorm"
)

func NewUserHandler(db *gorm.DB, ampq *publisher.AMQPPublisher, logger logging.Logger) *handler.UserHandler {
	return handler.NewUserHandler(
		repository.NewGormUserRepository(db),
		logging.NewZeroLogger(),
		validator.NewDefaultUserValidator(),
		ampq,
		errorhandler.NewErrorHandler(logger),
		security.NewBcryptPasswordHasher(),
	)
}
