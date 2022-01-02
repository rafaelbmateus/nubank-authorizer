package entity

import "errors"

var (
	ErrInvalidJson                  = errors.New("invalid-json")
	ErrAccountNotFound              = errors.New("account-not-found")
	ErrAccountDisabledCard          = errors.New("disabled-card")
	ErrAccountAlreadyInitialized    = errors.New("account-already-initialized")
	ErrAccountInvalidLimit          = errors.New("invalid-available-limit")
	ErrTransactionInvalidAmount     = errors.New("invalid-amount-value")
	ErrTransactionInsufficientLimit = errors.New("insufficient-limit")
	ErrTransactionHighFrequency     = errors.New("high-frequency-small-interval")
	ErrTransactionDoubled           = errors.New("doubled-transaction")
)
