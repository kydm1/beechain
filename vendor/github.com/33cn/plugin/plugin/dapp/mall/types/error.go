package types

import (
	"errors"
)

var (
	//supplychain
	ErrTxErr             	= errors.New("transaction format err")
	ErrWrongActionType		= errors.New("wrong tx action type")
	ErrWrongLength          = errors.New("wrong length")
	ErrWrongPlatformKey     = errors.New("wrong platform key")
	ErrEmptyValue			= errors.New("empty value")
	ErrWrongSignPubkey      = errors.New("wrong sign pubkey")
	ErrAmountLow            = errors.New("amount low")
	ErrDateConflict         = errors.New("wrong date")
	ErrWrongData            = errors.New("wrong data")
	ErrTokenNotExist        = errors.New("token not exists")
	ErrCurrencyNotEnough    = errors.New("currency not enough")
	ErrGoodExists           = errors.New("good already exists")
	ErrWrongType			= errors.New("wrong type")
	ErrSellerNotSamePayer   = errors.New("seller is same to payer")
)
