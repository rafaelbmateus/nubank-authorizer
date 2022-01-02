package entity

type Account struct {
	ActiveCard     bool `json:"active_card"`
	AvailableLimit int  `json:"available_limit"`
}

func CreateAccount(activeCard bool, availableLimit int) (*Account, error) {
	a := &Account{
		ActiveCard:     activeCard,
		AvailableLimit: availableLimit,
	}

	err := a.Validade()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Account) Validade() error {
	if a.AvailableLimit == 0 {
		return ErrAccountInvalidLimit
	}

	return nil
}
