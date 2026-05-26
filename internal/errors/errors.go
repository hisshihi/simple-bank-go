package errors

import "errors"

var (
	ErrInvalidAccountID       = errors.New("invalid account id")
	ErrBalanceIsLessThanZero  = errors.New("balance is less than zero")
	ErrCurrencyIsNotSupported = errors.New("currency is not supported")
)

var (
	ErrInvalidEntryID = errors.New("invalid entry id")
)

var (
	ErrInvalidTransferID    = errors.New("invalid transfer id")
	ErrInvalidFromAccountID = errors.New("invalid from account id")
	ErrInvalidToAccountID   = errors.New("invalid to account id")
)
