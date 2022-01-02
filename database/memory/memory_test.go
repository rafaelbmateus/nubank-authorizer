package memory_test

import (
	"nubank/authorizer/database/memory"
	"nubank/authorizer/entity"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	tests := []struct {
		account *entity.Account
	}{
		{
			account: &entity.Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			},
		},
		{
			account: &entity.Account{
				ActiveCard:     false,
				AvailableLimit: 1,
			},
		},
	}

	var database = memory.Memory{}
	for _, tc := range tests {
		database.CreateAccount(tc.account)
	}
}

func TestGetAccount(t *testing.T) {
	var database = memory.Memory{}
	database.CreateAccount(&entity.Account{
		ActiveCard:     true,
		AvailableLimit: 100,
	})
	acc, _ := database.GetAccount()
	assert.Equal(t, true, acc.ActiveCard)
	assert.Equal(t, 100, acc.AvailableLimit)
}

func TestAccountNotFound(t *testing.T) {
	var database = memory.Memory{}
	_, err := database.GetAccount()
	assert.Equal(t, entity.ErrAccountNotFound, err)
}

func TestNewTransaction(t *testing.T) {
	tests := []struct {
		tran *entity.Transaction
		err  error
	}{
		{
			tran: &entity.Transaction{
				ID:       uuid.New(),
				Merchant: "Bob's",
				Amount:   15,
				Time:     time.Now(),
			},
			err: nil,
		},
		{
			tran: &entity.Transaction{
				ID:       uuid.New(),
				Merchant: "ifood",
				Amount:   20,
				Time:     time.Now(),
			},
			err: nil,
		},
	}

	var database = memory.Memory{Account: &entity.Account{ActiveCard: true, AvailableLimit: 100}}
	for _, tc := range tests {
		tran, err := database.NewTransaction(tc.tran)
		assert.Equal(t, tc.err, err)
		assert.Equal(t, tc.tran.Merchant, tran.Merchant)
		assert.Equal(t, tc.tran.Amount, tran.Amount)
		assert.Equal(t, tc.tran.Time, tran.Time)
	}
}

func TestGetTransactions(t *testing.T) {
	var database = memory.Memory{Account: &entity.Account{ActiveCard: true, AvailableLimit: 100}}
	database.NewTransaction(&entity.Transaction{})
	trans := database.GetTransactions()
	assert.Equal(t, 1, len(trans))
}
