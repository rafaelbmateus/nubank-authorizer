package usecase_test

import (
	"context"
	"nubank/authorizer/database/memory"
	"nubank/authorizer/entity"
	"nubank/authorizer/usecase"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var accountService = usecase.NewAccountService(
	context.Background(), &memory.Memory{}, zerolog.DefaultContextLogger)

func TestCreateAccount(t *testing.T) {
	type test struct {
		account *entity.Account
		err     error
	}

	tests := []test{
		{
			account: &entity.Account{
				ActiveCard:     true,
				AvailableLimit: 0,
			},
			err: entity.ErrAccountInvalidLimit,
		},
		{
			account: &entity.Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			},
			err: nil,
		},
		{
			account: &entity.Account{
				ActiveCard:     true,
				AvailableLimit: 200,
			},
			err: entity.ErrAccountAlreadyInitialized,
		},
	}

	for _, tc := range tests {
		_, err := accountService.CreateAccount(tc.account.ActiveCard, tc.account.AvailableLimit)
		assert.Equal(t, tc.err, err)
	}
}
