// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"encoding/hex"

	"github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	tokenty "github.com/33cn/plugin/plugin/dapp/token/types"
)

func (t *token) ExecLocal_Transfer(payload *types.AssetsTransfer, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set, err := t.ExecLocalTransWithdraw(tx, receiptData, index)
	if err != nil {
		return nil, err
	}
	// 添加个人资产列表
	//tokenlog.Info("ExecLocalTransWithdraw", "addr", tx.GetRealToAddr(), "asset", transfer.Cointoken)
	kv := AddTokenToAssets(tx.GetRealToAddr(), t.GetLocalDB(), payload.Cointoken)
	if kv != nil {
		set.KV = append(set.KV, kv...)
	}
	if cfg.SaveTokenTxList {
		tokenAction := tokenty.TokenAction{
			Ty: tokenty.ActionTransfer,
			Value: &tokenty.TokenAction_Transfer{
				Transfer: payload,
			},
		}
		kvs, err := t.makeTokenTxKvs(tx, &tokenAction, receiptData, index, false)
		if err != nil {
			return nil, err
		}
		set.KV = append(set.KV, kvs...)
	}
	return set, nil
}

func (t *token) ExecLocal_Withdraw(payload *types.AssetsWithdraw, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set, err := t.ExecLocalTransWithdraw(tx, receiptData, index)
	if err != nil {
		return nil, err
	}
	// 添加个人资产列表
	kv := AddTokenToAssets(tx.From(), t.GetLocalDB(), payload.Cointoken)
	if kv != nil {
		set.KV = append(set.KV, kv...)
	}
	if cfg.SaveTokenTxList {
		tokenAction := tokenty.TokenAction{
			Ty: tokenty.ActionWithdraw,
			Value: &tokenty.TokenAction_Withdraw{
				Withdraw: payload,
			},
		}
		kvs, err := t.makeTokenTxKvs(tx, &tokenAction, receiptData, index, false)
		if err != nil {
			return nil, err
		}
		set.KV = append(set.KV, kvs...)
	}
	return set, nil
}

func (t *token) ExecLocal_TransferToExec(payload *types.AssetsTransferToExec, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set, err := t.ExecLocalTransWithdraw(tx, receiptData, index)
	if err != nil {
		return nil, err
	}
	if cfg.SaveTokenTxList {
		tokenAction := tokenty.TokenAction{
			Ty: tokenty.TokenActionTransferToExec,
			Value: &tokenty.TokenAction_TransferToExec{
				TransferToExec: payload,
			},
		}
		kvs, err := t.makeTokenTxKvs(tx, &tokenAction, receiptData, index, false)
		if err != nil {
			return nil, err
		}
		set.KV = append(set.KV, kvs...)
	}
	return set, nil
}

func (t *token) ExecLocal_TokenPreCreate(payload *tokenty.TokenPreCreate, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	localToken := newLocalToken(payload)
	localToken = setPrepare(localToken, tx.From(), t.GetHeight(), t.GetBlockTime())
	key := calcTokenStatusKeyLocal(payload.Symbol, payload.Owner, tokenty.TokenStatusPreCreated)

	var set []*types.KeyValue
	set = append(set, &types.KeyValue{Key: key, Value: types.Encode(localToken)})
	return &types.LocalDBSet{KV: set}, nil
}

func (t *token) ExecLocal_TokenFinishCreate(payload *tokenty.TokenFinishCreate, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	prepareKey := calcTokenStatusKeyLocal(payload.Symbol, payload.Owner, tokenty.TokenStatusPreCreated)
	localToken, err := loadLocalToken(payload.Symbol, payload.Owner, tokenty.TokenStatusPreCreated, t.GetLocalDB())
	if err != nil {
		return nil, err
	}
	localToken = setCreated(localToken, t.GetHeight(), t.GetBlockTime())
	key := calcTokenStatusKeyLocal(payload.Symbol, payload.Owner, tokenty.TokenStatusCreated)
	var set []*types.KeyValue
	set = append(set, &types.KeyValue{Key: prepareKey, Value: nil})
	set = append(set, &types.KeyValue{Key: key, Value: types.Encode(localToken)})
	kv := AddTokenToAssets(payload.Owner, t.GetLocalDB(), payload.Symbol)
	set = append(set, kv...)

	table := NewLogsTable(t.GetLocalDB())
	txIndex := dapp.HeightIndexStr(t.GetHeight(), int64(index))
	err = table.Add(&tokenty.LocalLogs{Symbol: payload.Symbol, TxIndex: txIndex, ActionType: tokenty.TokenActionFinishCreate, TxHash: "0x" + hex.EncodeToString(tx.Hash())})
	if err != nil {
		return nil, err
	}
	kv, err = table.Save()
	if err != nil {
		return nil, err
	}
	set = append(set, kv...)

	return &types.LocalDBSet{KV: set}, nil
}

func (t *token) ExecLocal_TokenRevokeCreate(payload *tokenty.TokenRevokeCreate, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	prepareKey := calcTokenStatusKeyLocal(payload.Symbol, payload.Owner, tokenty.TokenStatusPreCreated)
	localToken, err := loadLocalToken(payload.Symbol, payload.Owner, tokenty.TokenStatusPreCreated, t.GetLocalDB())
	if err != nil {
		return nil, err
	}
	localToken = setRevoked(localToken, t.GetHeight(), t.GetBlockTime())
	key := calcTokenStatusKeyLocal(payload.Symbol, payload.Owner, tokenty.TokenStatusCreateRevoked)
	var set []*types.KeyValue
	set = append(set, &types.KeyValue{Key: prepareKey, Value: nil})
	set = append(set, &types.KeyValue{Key: key, Value: types.Encode(localToken)})
	return &types.LocalDBSet{KV: set}, nil
}

func newLocalToken(payload *tokenty.TokenPreCreate) *tokenty.LocalToken {
	localToken := tokenty.LocalToken{
		Name:                payload.Name,
		Symbol:              payload.Symbol,
		Introduction:        payload.Introduction,
		Total:               payload.Total,
		Price:               payload.Price,
		Owner:               payload.Owner,
		Creator:             "",
		Status:              tokenty.TokenStatusPreCreated,
		CreatedHeight:       0,
		CreatedTime:         0,
		PrepareCreateHeight: 0,
		PrepareCreateTime:   0,
		Precision:           8,
		TotalTransferTimes:  0,
		RevokedHeight:       0,
		RevokedTime:         0,
		Category:            payload.Category,
	}
	return &localToken
}

func setPrepare(t *tokenty.LocalToken, creator string, height, time int64) *tokenty.LocalToken {
	t.Creator = creator
	t.PrepareCreateHeight = height
	t.PrepareCreateTime = time
	return t
}

func loadLocalToken(symbol, owner string, status int32, db db.KVDB) (*tokenty.LocalToken, error) {
	key := calcTokenStatusKeyLocal(symbol, owner, status)
	v, err := db.Get(key)
	if err != nil {
		return nil, err
	}
	var localToken tokenty.LocalToken
	err = types.Decode(v, &localToken)
	if err != nil {
		return nil, err
	}
	return &localToken, nil
}

func setCreated(t *tokenty.LocalToken, height, time int64) *tokenty.LocalToken {
	t.CreatedTime = time
	t.CreatedHeight = height
	t.Status = tokenty.TokenStatusCreated
	return t
}

func setRevoked(t *tokenty.LocalToken, height, time int64) *tokenty.LocalToken {
	t.RevokedTime = time
	t.RevokedHeight = height
	t.Status = tokenty.TokenStatusCreateRevoked
	return t
}

func setMint(t *tokenty.LocalToken, height, time, amount int64) *tokenty.LocalToken {
	t.Total = t.Total + amount
	return t
}

func setBurn(t *tokenty.LocalToken, height, time, amount int64) *tokenty.LocalToken {
	t.Total = t.Total - amount
	return t
}

func resetCreated(t *tokenty.LocalToken) *tokenty.LocalToken {
	t.CreatedTime = 0
	t.CreatedHeight = 0
	t.Status = tokenty.TokenStatusPreCreated
	return t
}

func resetRevoked(t *tokenty.LocalToken) *tokenty.LocalToken {
	t.RevokedTime = 0
	t.RevokedHeight = 0
	t.Status = tokenty.TokenStatusPreCreated
	return t
}

func resetMint(t *tokenty.LocalToken, height, time, amount int64) *tokenty.LocalToken {
	t.Total = t.Total - amount
	return t
}

func resetBurn(t *tokenty.LocalToken, height, time, amount int64) *tokenty.LocalToken {
	t.Total = t.Total + amount
	return t
}

func (t *token) ExecLocal_TokenMint(payload *tokenty.TokenMint, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	localToken, err := loadLocalToken(payload.Symbol, tx.From(), tokenty.TokenStatusCreated, t.GetLocalDB())
	if err != nil {
		return nil, err
	}
	localToken = setMint(localToken, t.GetHeight(), t.GetBlockTime(), payload.Amount)
	var set []*types.KeyValue
	key := calcTokenStatusKeyLocal(payload.Symbol, tx.From(), tokenty.TokenStatusCreated)
	set = append(set, &types.KeyValue{Key: key, Value: types.Encode(localToken)})

	table := NewLogsTable(t.GetLocalDB())
	txIndex := dapp.HeightIndexStr(t.GetHeight(), int64(index))
	err = table.Add(&tokenty.LocalLogs{Symbol: payload.Symbol, TxIndex: txIndex, ActionType: tokenty.TokenActionMint, TxHash: "0x" + hex.EncodeToString(tx.Hash())})
	if err != nil {
		return nil, err
	}
	kv, err := table.Save()
	if err != nil {
		return nil, err
	}
	set = append(set, kv...)

	return &types.LocalDBSet{KV: set}, nil
}

func (t *token) ExecLocal_TokenBurn(payload *tokenty.TokenBurn, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	localToken, err := loadLocalToken(payload.Symbol, tx.From(), tokenty.TokenStatusCreated, t.GetLocalDB())
	if err != nil {
		return nil, err
	}
	localToken = setBurn(localToken, t.GetHeight(), t.GetBlockTime(), payload.Amount)
	var set []*types.KeyValue
	key := calcTokenStatusKeyLocal(payload.Symbol, tx.From(), tokenty.TokenStatusCreated)
	set = append(set, &types.KeyValue{Key: key, Value: types.Encode(localToken)})

	table := NewLogsTable(t.GetLocalDB())
	txIndex := dapp.HeightIndexStr(t.GetHeight(), int64(index))
	err = table.Add(&tokenty.LocalLogs{Symbol: payload.Symbol, TxIndex: txIndex, ActionType: tokenty.TokenActionBurn, TxHash: "0x" + hex.EncodeToString(tx.Hash())})
	if err != nil {
		return nil, err
	}
	kv, err := table.Save()
	if err != nil {
		return nil, err
	}
	set = append(set, kv...)

	return &types.LocalDBSet{KV: set}, nil
}
