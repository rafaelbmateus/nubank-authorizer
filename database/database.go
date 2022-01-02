package database

import "nubank/authorizer/entity"

type Repository interface {
	CreateAccount(*entity.Account) (*entity.Account, error)
	GetAccount() (*entity.Account, error)
	NewTransaction(*entity.Transaction) (*entity.Transaction, error)
	GetTransactions() []*entity.Transaction
}
