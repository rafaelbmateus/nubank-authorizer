package entity_test

import (
	"testing"
	"time"

	"nubank/authorizer/entity"

	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	tran, _ := entity.NewTransaction(nil, "abc", 10, time.Now())
	assert.NotEmpty(t, tran.ID)
	assert.Equal(t, "abc", tran.Merchant)
	assert.Equal(t, 10, tran.Amount)
}

func TestValidadeTransaction(t *testing.T) {
	type test struct {
		merchant string
		amount   int
		time     time.Time
		err      error
	}

	tests := []test{
		{
			merchant: "Burger King",
			amount:   5,
			time:     time.Now(),
			err:      nil,
		},
		{
			merchant: "ifood",
			amount:   15,
			time:     time.Now(),
			err:      nil,
		},
		{
			merchant: "uber",
			amount:   0,
			time:     time.Now(),
			err:      entity.ErrTransactionInvalidAmount,
		},
	}

	for _, tc := range tests {
		_, err := entity.NewTransaction(nil, tc.merchant, tc.amount, tc.time)
		assert.Equal(t, tc.err, err)
	}
}
