package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID       uuid.UUID `json:"id"`
	Account  *Account  `json:"account"`
	Merchant string    `json:"merchant"`
	Amount   int       `json:"amount"`
	Time     time.Time `json:"time"`
}

func NewTransaction(account *Account, merchant string, amount int, time time.Time) (*Transaction, error) {
	t := &Transaction{
		ID:       uuid.New(),
		Account:  account,
		Merchant: merchant,
		Amount:   amount,
		Time:     time,
	}

	err := t.Validade()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Transaction) Validade() error {
	if t.Amount == 0 {
		return ErrTransactionInvalidAmount
	}

	return nil
}
