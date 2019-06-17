// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/33cn/chain33/account"
	"github.com/33cn/chain33/types"
	tokenty "github.com/33cn/plugin/plugin/dapp/tokennote/types"
)

func (t *tokennote) Exec_Transfer(payload *types.AssetsTransfer, tx *types.Transaction, index int) (*types.Receipt, error) {
	token := payload.GetCointoken()
	db, err := account.NewAccountDB(t.GetName(), token, t.GetStateDB())
	if err != nil {
		return nil, err
	}
	tokenAction := tokenty.TokennoteAction{
		Ty: tokenty.ActionTransfer,
		Value: &tokenty.TokennoteAction_Transfer{
			Transfer: payload,
		},
	}
	err = t.CheckTokennoteStatus(token,payload.To,payload.Amount)
	if err != nil {
		tokennotelog.Error("Transfer","tokennote status err ",err)
		return nil,err
	}
	return t.ExecTransWithdraw(db, tx, &tokenAction, index)
}

func (t *tokennote) Exec_Withdraw(payload *types.AssetsWithdraw, tx *types.Transaction, index int) (*types.Receipt, error) {
	token := payload.GetCointoken()
	db, err := account.NewAccountDB(t.GetName(), token, t.GetStateDB())
	if err != nil {
		return nil, err
	}
	tokenAction := tokenty.TokennoteAction{
		Ty: tokenty.ActionWithdraw,
		Value: &tokenty.TokennoteAction_Withdraw{
			Withdraw: payload,
		},
	}
	err = t.CheckTokennoteStatus(token,payload.To,payload.Amount)
	if err != nil {
		tokennotelog.Error("Withdraw","tokennote status err ",err)
		return nil,err
	}
	return t.ExecTransWithdraw(db, tx, &tokenAction, index)
}


func (t *tokennote) Exec_TransferToExec(payload *types.AssetsTransferToExec, tx *types.Transaction, index int) (*types.Receipt, error) {
	token := payload.GetCointoken()
	db, err := account.NewAccountDB(t.GetName(), token, t.GetStateDB())
	if err != nil {
		return nil, err
	}
	tokenAction := tokenty.TokennoteAction{
		Ty: tokenty.TokennoteActionTransferToExec,
		Value: &tokenty.TokennoteAction_TransferToExec{
			TransferToExec: payload,
		},
	}
	err = t.CheckTokennoteStatus(token,payload.To,payload.Amount)
	if err != nil {
		tokennotelog.Error("TransferToExec","tokennote status err ",err)
		return nil,err
	}
	return t.ExecTransWithdraw(db, tx, &tokenAction, index)
}


//创建白条
func (t *tokennote) Exec_TokennoteCreate(payload *tokenty.TokennoteCreate, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newTokennoteAction(t, "", tx)
	return action.create(payload)
}


//各方拒绝
func (t *tokennote) Exec_TokennoteLoanedReject(payload *tokenty.TokennoteLoanedReject, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newTokennoteAction(t, "", tx)
	return action.loanedReject(payload)
}

//借款人借款
func (t *tokennote) Exec_TokennoteLoan(payload *tokenty.TokennoteLoan, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newTokennoteAction(t, "", tx)
	return action.loan(payload)
}

//出借人放款确认
func (t *tokennote) Exec_TokennoteLoanedAgree(payload *tokenty.TokennoteLoanedAgree, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newTokennoteAction(t, "", tx)
	return action.loanedAgree(payload)
}

//出借人放款确认
func (t *tokennote) Exec_TokennoteCashed(payload *tokenty.TokennoteCashed, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newTokennoteAction(t, "", tx)
	return action.cashed(payload)
}


//白条增发
func (t *tokennote) Exec_TokennoteMint(payload *tokenty.TokennoteMint, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newTokennoteAction(t, "", tx)
	return action.mint(payload)
}

//白条销毁
func (t *tokennote) Exec_TokennoteBurn(payload *tokenty.TokennoteBurn, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newTokennoteAction(t, "", tx)
	return action.burn(payload)
}