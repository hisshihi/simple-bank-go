package errors

import "errors"

var (
	ErrInvalidAccountID       = errors.New("invalid account id")
	ErrBalanceIsLessThanZero  = errors.New("balance is less than zero")
	ErrCurrencyIsNotSupported = errors.New("currency is not supported")
	ErrOwnerRequired          = errors.New("owner is required")
	ErrCurrencyMismatch       = errors.New("currency mismatch")
	ErrInsufficientFunds      = errors.New("insufficient funds")
)

var (
	ErrInvalidEntryID = errors.New("invalid entry id")
)

var (
	ErrInvalidTransferID    = errors.New("invalid transfer id")
	ErrInvalidFromAccountID = errors.New("invalid from account id")
	ErrInvalidToAccountID   = errors.New("invalid to account id")
)
