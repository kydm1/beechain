// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

const (
	// ActionTransfer for transfer
	ActionTransfer = 4
	// ActionGenesis for genesis
	ActionGenesis = 5
	// ActionWithdraw for Withdraw
	ActionWithdraw = 6
	// TokennoteActionPreCreate for token pre create
	TokennoteActionCreate = 7
	// TokennoteActionFinishCreate for token finish create
	TokennoteActionLoan = 8
	// TokennoteActionRevokeCreate for token revoke create
	TokennoteActionLoanedAgree = 9
	TokennoteActionCashed = 10
	// TokennoteActionTransferToExec for token transfer to exec
	TokennoteActionTransferToExec = 11
	TokennoteActionLoanedReject = 12
	TokennoteActionMint = 13
	TokennoteActionBurn = 14
)

// token status
const (
	// TokennoteStatusPreCreated token pre create status
	TokennoteStatusPreCreated = iota
	// TokennoteStatusCreated token create status
	TokennoteStatusCreated
	//tokennote cashed
	TokennoteStatusCashed
	//
	TokennoteStatusReadyLoan
	//
	TokennoteStatusAgreeLoan
	//
	TokennoteStatusRejectLoan
	//
	TokennoteStatusFrozenInExecer
)

var (
	// TokennoteX token name
	TokennoteX = "tokennote"
	TokenX     = "token"
	TokenCCNY  = "CCNY"
	TokennoteRateMax = 2000
	TokennoteRateUnit = 10000
	TokennoteExpireTime = 3600*24
	TokennoteOverdueRate = 1000 //0.1 * 1e4   年利率
	TokennoteOverdueDayRate = 28 //0.00028 * 1e5  日利率
)

const (
	// TyLogPreCreateTokennote log for pre create token
	TyLogTokennoteCreate = 900
	// TyLogFinishCreateTokennote log for finish create token
	TyLogTokennoteLoan = 901
	// TyLogRevokeCreateTokennote log for revoke create token
	TyLogTokennoteLoanedAgree = 902
	TyLogTokennoteCashed = 903
	TyLogTokennoteLoanedReject = 904
	TyLogTokennoteMint = 905
	TyLogTokennoteBurn = 906
	TyLogTokennoteMarket = 907
	// TyLogTokennoteTransfer log for token tranfer
	TyLogTokennoteTransfer = 313
	// TyLogTokennoteGenesis log for token genesis
	TyLogTokennoteGenesis = 314
	// TyLogTokennoteDeposit log for token deposit
	TyLogTokennoteDeposit = 315
	// TyLogTokennoteExecTransfer log for token exec transfer
	TyLogTokennoteExecTransfer = 316
	// TyLogTokennoteExecWithdraw log for token exec withdraw
	TyLogTokennoteExecWithdraw = 317
	// TyLogTokennoteExecDeposit log for token exec deposit
	TyLogTokennoteExecDeposit = 318
	// TyLogTokennoteExecFrozen log for token exec frozen
	TyLogTokennoteExecFrozen = 319
	// TyLogTokennoteExecActive log for token exec active
	TyLogTokennoteExecActive = 320
	// TyLogTokennoteGenesisTransfer log for token genesis rransfer
	TyLogTokennoteGenesisTransfer = 321
	// TyLogTokennoteGenesisDeposit log for token genesis deposit
	TyLogTokennoteGenesisDeposit = 322
)

const (
	// TokennoteNameLenLimit token name length limit
	TokennoteNameLenLimit = 128
	// TokennoteSymbolLenLimit token symbol length limit
	TokennoteSymbolLenLimit = 32
	// TokennoteIntroLenLimit token introduction length limit
	TokennoteIntroLenLimit = 1024
	// less than const should be signed by admin
	TokennoteSymbolAdminLimit = 6
)
