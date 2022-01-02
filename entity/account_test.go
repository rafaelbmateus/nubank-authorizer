package entity_test

import (
	"testing"

	"nubank/authorizer/entity"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	acc, _ := entity.CreateAccount(true, 100)
	assert.Equal(t, true, acc.ActiveCard)
	assert.Equal(t, 100, acc.AvailableLimit)
}

func TestValidadeAccount(t *testing.T) {
	type test struct {
		activeCard     bool
		availableLimit int
		err            error
	}

	tests := []test{
		{
			activeCard:     true,
			availableLimit: 0,
			err:            entity.ErrAccountInvalidLimit,
		},
		{
			activeCard:     true,
			availableLimit: 100,
			err:            nil,
		},
	}

	for _, tc := range tests {
		_, err := entity.CreateAccount(tc.activeCard, tc.availableLimit)
		assert.Equal(t, tc.err, err)
	}
}
