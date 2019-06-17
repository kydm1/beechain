package types

import "errors"

var (
	//goldbill
	ErrPlatformExists 		= errors.New("ErrPlatformExists")
	ErrPlatformNotExists 	= errors.New("ErrPlatformNotExists")
	ErrUserExists        	= errors.New("ErrUserExists")
	ErrUserNotExists     	= errors.New("ErrUserNotExists")
	ErrAdminExists       	= errors.New("ErrAdminExists")
	ErrAdminNotExists    	= errors.New("ErrAdminNotExists")
	ErrPubkeyExists      	= errors.New("ErrPubkeyExists")
	ErrBillExists        	= errors.New("ErrBillExists")
	ErrBillNotExists     	= errors.New("ErrBillNotExists")
	ErrRMBNotEnough      	= errors.New("ErrRMBNotEnough")
	ErrCoinNotEnough     	= errors.New("ErrCoinNotEnough")
	ErrDupBillId         	= errors.New("ErrDupBillId")
	ErrWrongState           = errors.New("ErrWrongState")
	ErrNoPrivilege          = errors.New("ErrNoPrivilege")
)
