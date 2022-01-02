package usecase_test

import (
	"context"
	"nubank/authorizer/database/memory"
	"nubank/authorizer/entity"
	"nubank/authorizer/usecase"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	type test struct {
		name     string
		merchant string
		amount   int
		time     string
		err      error
	}

	accountService.Database.CreateAccount(&entity.Account{
		ActiveCard: true, AvailableLimit: 100})

	tests := []test{
		{
			name:     "basic transaction",
			merchant: "Uber",
			amount:   2,
			time:     "2020-12-13T10:00:00.000Z",
			err:      nil,
		},
		{
			name:     "double transaction same time",
			merchant: "Uber",
			amount:   2,
			time:     "2020-12-13T10:00:00.000Z",
			err:      entity.ErrTransactionDoubled,
		},
		{
			name:     "double transaction after two minutes",
			merchant: "Uber",
			amount:   2,
			time:     "2020-12-13T10:02:00.000Z",
			err:      entity.ErrTransactionDoubled,
		},
		{
			name:     "transaction successfully one second double transaction after",
			merchant: "Uber",
			amount:   2,
			time:     "2020-12-13T10:02:01.000Z",
			err:      nil,
		},
		{
			name:     "insufficient limit",
			merchant: "Taxi",
			amount:   101,
			time:     "2020-12-13T11:00:00.000Z",
			err:      entity.ErrTransactionInsufficientLimit,
		},
		{
			name:     "bk number 1",
			merchant: "Burger King",
			amount:   3,
			time:     "2020-12-13T12:00:00.000Z",
			err:      nil,
		},
		{
			name:     "bk number 2",
			merchant: "Burger King",
			amount:   2,
			time:     "2020-12-13T12:00:10.000Z",
			err:      nil,
		},
		{
			name:     "bk number 3",
			merchant: "Burger King",
			amount:   1,
			time:     "2020-12-13T12:00:25.000Z",
			err:      nil,
		},
		{
			name:     "bk high frequency in small interval",
			merchant: "Burger King",
			amount:   10,
			time:     "2020-12-13T12:00:47.000Z",
			err:      entity.ErrTransactionHighFrequency,
		},
		{
			name:     "invalid amount value",
			merchant: "Bob's",
			amount:   0,
			time:     "2020-12-13T10:00:00.000Z",
			err:      entity.ErrTransactionInvalidAmount,
		},
		{
			name:     "insufficient limit after spending a little money",
			merchant: "Bar",
			amount:   91,
			time:     "2020-12-14T11:00:00.000Z",
			err:      entity.ErrTransactionInsufficientLimit,
		},
	}

	var transactionService = usecase.NewTransactionService(context.Background(),
		&memory.Memory{Account: &entity.Account{ActiveCard: true, AvailableLimit: 100}},
		zerolog.DefaultContextLogger)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			time, _ := time.Parse("2006-01-02T15:04:05.000Z", tc.time)
			_, err := transactionService.NewTransaction(tc.merchant, tc.amount, time)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestNewTransactionAccountNotFound(t *testing.T) {
	var transactionService = usecase.NewTransactionService(context.Background(),
		&memory.Memory{}, zerolog.DefaultContextLogger)
	_, err := transactionService.NewTransaction("Amazon", 1000, time.Now())
	assert.Equal(t, entity.ErrAccountNotFound, err)
}

func TestDisabledCard(t *testing.T) {
	var transactionService = usecase.NewTransactionService(context.Background(),
		&memory.Memory{Account: &entity.Account{ActiveCard: false, AvailableLimit: 2000}},
		zerolog.DefaultContextLogger)
	_, err := transactionService.NewTransaction("Amazon", 1000, time.Now())
	assert.Equal(t, entity.ErrAccountDisabledCard, err)
}
