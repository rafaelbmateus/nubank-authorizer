package memory

import (
	"nubank/authorizer/entity"
)

// Memory repository.
type Memory struct {
	Account      *entity.Account
	Transactions []*entity.Transaction
}

// CreateAccount create a new account in the memory.
func (m *Memory) CreateAccount(acc *entity.Account) (*entity.Account, error) {
	m.Account = acc
	return m.Account, nil
}

// GetAccount get account.
func (m *Memory) GetAccount() (*entity.Account, error) {
	if m.Account == nil {
		return nil, entity.ErrAccountNotFound
	}
	return m.Account, nil
}

// NewTransaction save a new transaction in the memory.
func (m *Memory) NewTransaction(tran *entity.Transaction) (*entity.Transaction, error) {
	m.updateAvailableLimit(tran.Amount)
	transactions := append(m.Transactions, tran)
	m.Transactions = transactions
	m.Account, _ = m.GetAccount()
	return tran, nil
}

// GetTransactions get all transactions in the memory.
func (m *Memory) GetTransactions() []*entity.Transaction {
	return m.Transactions
}

// updateAvailableLimit update account available limit.
func (m *Memory) updateAvailableLimit(amount int) {
	m.Account.AvailableLimit -= amount
}
