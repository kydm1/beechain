// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "errors"

var (
	// ErrTokennoteNameLen error token name length
	ErrTokennoteNameLen = errors.New("ErrTokennoteNameLength")
	// ErrTokennoteSymbolLen error token symbol length
	ErrTokennoteSymbolLen = errors.New("ErrTokennoteSymbolLength")
	// ErrTokennoteTotalOverflow error token total Overflow
	ErrTokennoteTotalOverflow = errors.New("ErrTokennoteTotalOverflow")
	// ErrTokennoteSymbolUpper error token total Overflow
	ErrTokennoteSymbolUpper = errors.New("ErrTokennoteSymbolUpper")
	// ErrTokennoteIntroLen error token introduction length
	ErrTokennoteIntroLen = errors.New("ErrTokennoteIntroductionLen")
	// ErrTokennoteExist error token symbol exist already
	ErrTokennoteExist = errors.New("ErrTokennoteSymbolExistAlready")
	// ErrTokennoteNotPrecreated error token not pre created
	ErrTokennoteNotPrecreated = errors.New("ErrTokennoteNotPrecreated")
	// ErrTokennoteCreatedApprover error token created approver
	ErrTokennoteCreatedApprover = errors.New("ErrTokennoteCreatedApprover")
	// ErrTokennoteRevoker error token revoker
	ErrTokennoteRevoker = errors.New("ErrTokennoteRevoker")
	// ErrTokennoteCanotRevoked error token canot revoked with wrong status
	ErrTokennoteCanotRevoked = errors.New("ErrTokennoteCanotRevokedWithWrongStatus")
	// ErrTokennoteOwner error token symbol owner not match
	ErrTokennoteOwner = errors.New("ErrTokennoteSymbolOwnerNotMatch")
	// ErrTokennoteHavePrecreated error owner have token pre create yet
	ErrTokennoteHavePrecreated = errors.New("ErrOwnerHaveTokennotePrecreateYet")
	// ErrTokennoteBlacklist error token blacklist
	ErrTokennoteBlacklist = errors.New("ErrTokennoteBlacklist")
	// ErrTokennoteNotExist error token symbol not exist
	ErrTokennoteNotExist = errors.New("ErrTokennoteSymbolNotExist")
	ErrTokennoteNotCreated = errors.New("ErrTokennoteNotCreated")
	ErrTokennoteAmountLow  = errors.New("ErrTokennoteAmountLow")
	ErrTokennoteCashed     = errors.New("ErrTokennoteCashed")
	ErrTokennoteExcept     = errors.New("ErrTokennoteExcept")
	ErrTokennoteNumberFormat = errors.New("ErrTokennoteNumberFormat")
	ErrTokennoteNotReadyLoan = errors.New("ErrTokennoteNotReadyLoan")
	ErrTokenoteNotAvaible    = errors.New("ErrTokenoteNotAvaible")
	ErrTokennoteNotLoanToSelf = errors.New("ErrTokennoteNotLoanToSelf")
	ErrTokennoteWrongStateKey = errors.New("ErrTokennoteWrongStateKey")
	ErrTokennoteCreatedNotAllowed = errors.New("ErrTokennoteCreatedNotAllowed")
	ErrTokennoteNotAllowedLoanTwice = errors.New("ErrTokennoteNotAllowedLoanTwice")
	ErrTokennoteOverdue           = errors.New("ErrTokennoteOverdue")
)
