package types

import (
	"errors"
)

var (
	//supplychain
	ErrRNE             = errors.New("RMB not enough")
	ErrDateConflict    = errors.New("ErrDateConflict")
	ErrAmountLow       = errors.New("ErrAmountLow")
	ErrPurchaseTooMuch = errors.New("ErrPurchaseTooMuch")

)
