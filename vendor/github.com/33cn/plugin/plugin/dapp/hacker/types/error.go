package types

import (
	"errors"
)

var (
	//supplychain
	ErrDocodeErr            = errors.New("data decode error")
	ErrWrongActionType		= errors.New("wrong action type")
	ErrEmptyValue			= errors.New("empty value")
	ErrWrongPubkey			= errors.New("wrong pubkey")
)
