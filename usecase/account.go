package usecase

import (
	"context"
	"nubank/authorizer/database"
	"nubank/authorizer/entity"

	"github.com/rs/zerolog"
)

type AccountService struct {
	Context  context.Context
	Database database.Repository
	Logger   *zerolog.Logger
}

func NewAccountService(ctx context.Context, db database.Repository, log *zerolog.Logger) *AccountService {
	return &AccountService{
		Context:  ctx,
		Database: db,
		Logger:   log,
	}
}

func (as AccountService) CreateAccount(activeCard bool, availableLimit int) (*entity.Account, error) {
	if acc, _ := as.Database.GetAccount(); acc != nil {
		return acc, entity.ErrAccountAlreadyInitialized
	}

	acc, err := entity.CreateAccount(activeCard, availableLimit)
	if err != nil {
		return nil, err
	}

	return as.Database.CreateAccount(acc)
}
