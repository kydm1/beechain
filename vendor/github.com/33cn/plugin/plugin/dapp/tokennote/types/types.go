// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"encoding/json"
	"reflect"

	log "github.com/33cn/chain33/common/log/log15"

	"github.com/33cn/chain33/types"
)

var tokenlog = log.New("module", "execs.tokennote.types")

func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(TokennoteX))
	types.RegistorExecutor(TokennoteX, NewType())
	types.RegisterDappFork(TokennoteX, "Enable", 0)
}

// TokennoteType 执行器基类结构体
type TokennoteType struct {
	types.ExecTypeBase
}

// NewType 创建执行器类型
func NewType() *TokennoteType {
	c := &TokennoteType{}
	c.SetChild(c)
	return c
}

// GetName 获取执行器名称
func (t *TokennoteType) GetName() string {
	return TokennoteX
}

// GetPayload 获取token action
func (t *TokennoteType) GetPayload() types.Message {
	return &TokennoteAction{}
}

// GetTypeMap 根据action的name获取type
func (t *TokennoteType) GetTypeMap() map[string]int32 {
	return map[string]int32{
		"Transfer":          ActionTransfer,
		"Genesis":           ActionGenesis,
		"Withdraw":          ActionWithdraw,
		"TokennoteCreate":    TokennoteActionCreate,
		"TokennoteLoan": TokennoteActionLoan,
		"TokennoteLoanedAgree": TokennoteActionLoanedAgree,
		"TokennoteCashed":TokennoteActionCashed,
		"TokennoteLoanedReject":TokennoteActionLoanedReject,
		"TransferToExec":    TokennoteActionTransferToExec,
		"TokennoteMint":	TokennoteActionMint,
		"TokennoteBurn":	TokennoteActionBurn,
	}
}

// GetLogMap 获取log的映射对应关系
func (t *TokennoteType) GetLogMap() map[int64]*types.LogInfo {
	return map[int64]*types.LogInfo{
		TyLogTokennoteTransfer:        {Ty: reflect.TypeOf(types.ReceiptAccountTransfer{}), Name: "LogTokenTransfer"},
		TyLogTokennoteDeposit:         {Ty: reflect.TypeOf(types.ReceiptAccountTransfer{}), Name: "LogTokenDeposit"},
		TyLogTokennoteExecTransfer:    {Ty: reflect.TypeOf(types.ReceiptExecAccountTransfer{}), Name: "LogTokenExecTransfer"},
		TyLogTokennoteExecWithdraw:    {Ty: reflect.TypeOf(types.ReceiptExecAccountTransfer{}), Name: "LogTokenExecWithdraw"},
		TyLogTokennoteExecDeposit:     {Ty: reflect.TypeOf(types.ReceiptExecAccountTransfer{}), Name: "LogTokenExecDeposit"},
		TyLogTokennoteExecFrozen:      {Ty: reflect.TypeOf(types.ReceiptExecAccountTransfer{}), Name: "LogTokenExecFrozen"},
		TyLogTokennoteExecActive:      {Ty: reflect.TypeOf(types.ReceiptExecAccountTransfer{}), Name: "LogTokenExecActive"},
		TyLogTokennoteGenesisTransfer: {Ty: reflect.TypeOf(types.ReceiptAccountTransfer{}), Name: "LogTokenGenesisTransfer"},
		TyLogTokennoteGenesisDeposit:  {Ty: reflect.TypeOf(types.ReceiptExecAccountTransfer{}), Name: "LogTokenGenesisDeposit"},
		TyLogTokennoteCreate:       {Ty: reflect.TypeOf(ReceiptTokennote{}), Name: "LogCreateToken"},
		TyLogTokennoteLoan:       {Ty: reflect.TypeOf(ReceiptTokennote{}), Name: "TyLogTokennoteLoan"},
		TyLogTokennoteLoanedAgree:       {Ty: reflect.TypeOf(ReceiptTokennote{}), Name: "TyLogTokennoteLoanedAgree"},
		TyLogTokennoteLoanedReject:       {Ty: reflect.TypeOf(ReceiptTokennote{}), Name: "TyLogTokennoteLoanedReject"},
		TyLogTokennoteCashed:       {Ty: reflect.TypeOf(ReceiptTokennoteCashed{}), Name: "TyLogTokennoteCashed"},
		TyLogTokennoteMint:       {Ty: reflect.TypeOf(ReceiptTokennote{}), Name: "TyLogTokennoteMint"},
		TyLogTokennoteBurn:       {Ty: reflect.TypeOf(ReceiptTokennote{}), Name: "TyLogTokennoteBurn"},
		TyLogTokennoteMarket:       {Ty: reflect.TypeOf(Tokennote{}), Name: "TyLogTokennoteMarket"},
	}
}

// RPC_Default_Process rpc 默认处理
func (t *TokennoteType) RPC_Default_Process(action string, msg interface{}) (*types.Transaction, error) {
	var create *types.CreateTx
	if _, ok := msg.(*types.CreateTx); !ok {
		return nil, types.ErrInvalidParam
	}
	create = msg.(*types.CreateTx)
	if !create.IsToken {
		return nil, types.ErrNotSupport
	}
	tx, err := t.AssertCreate(create)
	if err != nil {
		return nil, err
	}
	//to地址的问题,如果是主链交易，to地址就是直接是设置to
	if !types.IsPara() {
		tx.To = create.To
	}
	return tx, err
}

// CreateTx token 创建合约
func (t *TokennoteType) CreateTx(action string, msg json.RawMessage) (*types.Transaction, error) {
	tx, err := t.ExecTypeBase.CreateTx(action, msg)
	if err != nil {
		tokenlog.Error("token CreateTx failed", "err", err, "action", action, "msg", string(msg))
		return nil, err
	}
	if !types.IsPara() {
		var transfer TokennoteAction
		err = types.Decode(tx.Payload, &transfer)
		if err != nil {
			tokenlog.Error("token CreateTx failed", "decode payload err", err, "action", action, "msg", string(msg))
			return nil, err
		}
		if action == "Transfer" {
			tx.To = transfer.GetTransfer().To
		} else if action == "Withdraw" {
			tx.To = transfer.GetWithdraw().To
		} else if action == "TransferToExec" {
			tx.To = transfer.GetTransferToExec().To
		}
	}
	return tx, nil
}
