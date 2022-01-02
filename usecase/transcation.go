package usecase

import (
	"context"
	"nubank/authorizer/database"
	"nubank/authorizer/entity"
	"time"

	"github.com/rs/zerolog"
)

// TransactionService represents the transaction service.
type TransactionService struct {
	Context  context.Context
	Database database.Repository
	Logger   *zerolog.Logger
}

// NewTransactionService create a new TransactionService instance.
func NewTransactionService(ctx context.Context, db database.Repository, log *zerolog.Logger) *TransactionService {
	return &TransactionService{
		Context:  ctx,
		Database: db,
		Logger:   log,
	}
}

// NewTransaction add new transcation use case.
func (ts *TransactionService) NewTransaction(merchant string, amount int, time time.Time) (*entity.Transaction, error) {
	acc, err := ts.Database.GetAccount()
	if err != nil {
		return nil, err
	}

	tran, err := entity.NewTransaction(acc, merchant, amount, time)
	if err != nil {
		return nil, err
	}

	if !acc.ActiveCard {
		return tran, entity.ErrAccountDisabledCard
	}

	trans := ts.Database.GetTransactions()

	if err := ts.checkDoubledTransaction(trans, tran); err != nil {
		return tran, err
	}

	if err := ts.checkHighFrequency(trans, tran); err != nil {
		return tran, err
	}

	if err := ts.checkInsufficientLimit(acc, amount); err != nil {
		return tran, err
	}

	return ts.Database.NewTransaction(tran)
}

// check rule: doubled-transaction
func (ts *TransactionService) checkDoubledTransaction(trans []*entity.Transaction, tran *entity.Transaction) error {
	for _, t := range trans {
		if t.Merchant != tran.Merchant {
			continue
		}

		if t.Amount != tran.Amount {
			continue
		}

		tString := t.Time.String()
		tranString := tran.Time.String()
		print(tString, tranString)

		if tran.Time.After(t.Time.Add(time.Duration(2)*time.Minute)) &&
			tran.Time.After(t.Time) {
			continue
		}

		return entity.ErrTransactionDoubled
	}

	return nil
}

// check rule: high-frequency-small-interval
func (ts *TransactionService) checkHighFrequency(trans []*entity.Transaction, tran *entity.Transaction) error {
	count := 0
	for _, t := range trans {
		t1 := t.Time.String()
		t2 := tran.Time.String()
		print(t1, t2)
		if t.Time.After(tran.Time.Add(time.Duration(-2) * time.Minute)) {
			count += 1
		}

		if count > 2 {
			return entity.ErrTransactionHighFrequency
		}
	}

	return nil
}

// check rule: insufficient-limit
func (ts *TransactionService) checkInsufficientLimit(acc *entity.Account, amount int) error {
	if acc.AvailableLimit < amount {
		return entity.ErrTransactionInsufficientLimit
	}

	return nil
}
