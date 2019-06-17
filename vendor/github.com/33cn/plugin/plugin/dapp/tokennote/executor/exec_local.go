// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/types"
	tokenty "github.com/33cn/plugin/plugin/dapp/tokennote/types"
	"time"
	"github.com/33cn/chain33/common/db/table"
)

func (t *tokennote) ExecLocal_Transfer(payload *types.AssetsTransfer, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set, err := t.ExecLocalTransWithdraw(tx, receiptData, index)
	if err != nil {
		return nil, err
	}
	// 添加个人资产列表
	//tokenlog.Info("ExecLocalTransWithdraw", "addr", tx.GetRealToAddr(), "asset", transfer.Cointoken)
	kv := AddTokennoteToAssets(tx.GetRealToAddr(), t.GetLocalDB(), payload.Cointoken,t.GetBlockTime())
	if kv != nil {
		set.KV = append(set.KV, kv...)
	}
	if cfg.SaveTokennoteTxList {
		tokenAction := tokenty.TokennoteAction{
			Ty: tokenty.ActionTransfer,
			Value: &tokenty.TokennoteAction_Transfer{
				Transfer: payload,
			},
		}
		kvs, err := t.makeTokennoteTxKvs(tx, &tokenAction, receiptData, index, false)
		if err != nil {
			return nil, err
		}
		set.KV = append(set.KV, kvs...)
	}
	return set, nil
}

func (t *tokennote) ExecLocal_Withdraw(payload *types.AssetsWithdraw, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set, err := t.ExecLocalTransWithdraw(tx, receiptData, index)
	if err != nil {
		return nil, err
	}
	// 添加个人资产列表
	kv := AddTokennoteToAssets(tx.From(), t.GetLocalDB(), payload.Cointoken,t.GetBlockTime())
	if kv != nil {
		set.KV = append(set.KV, kv...)
	}
	if cfg.SaveTokennoteTxList {
		tokenAction := tokenty.TokennoteAction{
			Ty: tokenty.ActionWithdraw,
			Value: &tokenty.TokennoteAction_Withdraw{
				Withdraw: payload,
			},
		}
		kvs, err := t.makeTokennoteTxKvs(tx, &tokenAction, receiptData, index, false)
		if err != nil {
			return nil, err
		}
		set.KV = append(set.KV, kvs...)
	}
	return set, nil
}

func (t *tokennote) ExecLocal_TransferToExec(payload *types.AssetsTransferToExec, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	if receiptData.GetTy() != types.ExecOk {
		return nil, nil
	}
	set, err := t.ExecLocalTransWithdraw(tx, receiptData, index)
	if err != nil {
		return nil, err
	}
	if cfg.SaveTokennoteTxList {
		tokenAction := tokenty.TokennoteAction{
			Ty: tokenty.TokennoteActionTransferToExec,
			Value: &tokenty.TokennoteAction_TransferToExec{
				TransferToExec: payload,
			},
		}
		kvs, err := t.makeTokennoteTxKvs(tx, &tokenAction, receiptData, index, false)
		if err != nil {
			return nil, err
		}
		set.KV = append(set.KV, kvs...)
	}
	return set, nil
}


func (t *tokennote) ExecLocal_TokennoteCreate(payload *tokenty.TokennoteCreate, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	if receiptData.GetTy() != types.ExecOk {
		return nil, nil
	}
	localToken := newLocalTokennote(payload)

	localToken = setCreated(localToken, t.GetHeight(), t.GetBlockTime())
	key := calcTokennoteStatusKeyLocal(payload.Currency, tokenty.TokennoteStatusCreated)
	var set []*types.KeyValue
	set = append(set, &types.KeyValue{Key: key, Value: types.Encode(localToken)})
	tokennotelog.Error("create note","height time",t.GetBlockTime())
	kv := AddTokennoteToAssets(payload.Issuer, t.GetLocalDB(), payload.Currency,time.Now().Unix())
	set = append(set, kv...)

	tokennotelog.Error("ExecLocal_TokennoteCreate","marketopt",marketopt)
	tab ,err := table.NewTable(NewTokennoteMarketRow(),t.GetLocalDB(),&marketopt)
	if err != nil {
		tokennotelog.Error("ExecLocal_TokennoteCreate","new tokennotemarket init err ",err)
		return nil, err
	}

	ccount  := table.NewCount(tokennoteLocalPre,"createdNotes",t.GetLocalDB())

	for _ ,v := range receiptData.Logs {
		if v.Ty == tokenty.TyLogTokennoteMarket {
			var t tokenty.Tokennote
			err := types.Decode(v.Log,&t)
			if err != nil {
				tokennotelog.Error("ExecLocal_TokennoteCreate","decode tokennote err",err)
				return nil,err
			}
			tokennotelog.Error("ExecLocal_TokennoteCreate","table ","ok")
			tab.Add(&t)
			ccount.Inc()
		}
	}
	tabkv,err := tab.Save()
	if err != nil {
		return nil,err
	}
	for _,v := range tabkv {
		tokennotelog.Error("ExeclLocal_Create","set market kv",string(v.Key))
	}
	ccountkv ,err := ccount.Save()
	if err != nil {
		return nil,err
	}
	set = append(set,tabkv...)
	set = append(set,ccountkv...)
	for _, v := range set {
		tokennotelog.Error("ExeclLocal_Create","set kv",string(v.Key))
	}
	return &types.LocalDBSet{KV: set}, nil
}

func (t *tokennote) ExecLocal_TokennoteLoan(payload *tokenty.TokennoteLoan, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	if receiptData.GetTy() != types.ExecOk {
		return nil, nil
	}
	var set []*types.KeyValue
	//白条在合约地址中
	kv := AddTokennoteToAssets(tx.To, t.GetLocalDB(), payload.Symbol,t.GetBlockTime())
	kv1 := AddTokennoteToAssets(payload.To, t.GetLocalDB(), payload.Symbol,t.GetBlockTime())
	set = append(set, kv...)
	set = append(set,kv1...)
	return &types.LocalDBSet{KV: set}, nil
}

func (t *tokennote) ExecLocal_TokennoteLoanedAgree(payload *tokenty.TokennoteLoanedAgree, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	if receiptData.GetTy() != types.ExecOk {
		return nil, nil
	}
	var set []*types.KeyValue
	kv := AddTokennoteToAssets(tx.From(), t.GetLocalDB(), payload.Symbol,t.GetBlockTime())
	//kv_del := DeleteTokennoteInAssets(tx.From(), t.GetLocalDB(), payload.Symbol,payload.Loantime)
	set = append(set, kv...)
	//set = append(set,kv_del...)
	return &types.LocalDBSet{KV: set}, nil
}

func (t *tokennote) ExecLocal_TokennoteLoanedReject(payload *tokenty.TokennoteLoanedReject, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	if receiptData.GetTy() != types.ExecOk {
		return nil, nil
	}
	var set []*types.KeyValue
	kv_del := DeleteTokennoteInAssets(tx.From(), t.GetLocalDB(), payload.Symbol,payload.LoanTime)
	set = append(set,kv_del...)
	return &types.LocalDBSet{KV: set}, nil
}

func (t *tokennote) ExecLocal_TokennoteCashed(payload *tokenty.TokennoteCashed, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {

	var set []*types.KeyValue
	if receiptData.GetTy() != types.ExecOk {
		return nil, nil
	}
	for _, v := range receiptData.Logs {
		if v.Ty == tokenty.TyLogTokennoteCashed {
			var reply tokenty.ReceiptTokennoteCashed
			err := types.Decode(v.Log,&reply)
			if err != nil {
				return nil,err
			}
			for _,vv := range reply.Cashlist {
				keys := tokenCashedTxkeys(payload.Symbol,tx.From(),vv.Addr,t.GetHeight(),int64(index))
				for _, vvv := range keys {
					set = append(set,&types.KeyValue{Key:vvv,Value:types.Encode(vv)})
					tokennotelog.Error("ExecLocal_TokennoteCashed","set key",string(vvv), " value ",vv)
				}

			}

		}
	}

	return &types.LocalDBSet{KV: set}, nil
}

func (t *tokennote) ExecLocal_TokennoteMint(payload *tokenty.TokennoteMint, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	if receiptData.GetTy() != types.ExecOk {
		return nil, nil
	}
	var set []*types.KeyValue

	tab,err := table.NewTable(NewTransactionRow(),t.GetLocalDB(),&opt)
	if err != nil {
		return nil,err
	}
	tab.Add(tx)
	kvs,err := tab.Save()
	if err != nil {
		return nil,err
	}

	for _,v := range kvs {
		tokennotelog.Error("mint","table key value ",string(v.Key))
		t.GetLocalDB().Set(v.Key,v.Value)
	}
	set = append(set,kvs...)
	return &types.LocalDBSet{KV: set}, nil
}


func (t *tokennote) ExecLocal_TokennoteBurn(payload *tokenty.TokennoteBurn, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	if receiptData.GetTy() != types.ExecOk {
		return nil, nil
	}
	var set []*types.KeyValue

	tab,err := table.NewTable(NewTransactionRow(),t.GetLocalDB(),&opt)
	if err != nil {
		return nil,err
	}
	tab.Add(tx)
	kvs,err := tab.Save()
	if err != nil {
		return nil,err
	}
	for _,v := range kvs {
		tokennotelog.Error("burn","table list key ",string(v.Key))
		t.GetLocalDB().Set(v.Key,v.Value)
	}
	set = append(set,kvs...)
	return &types.LocalDBSet{KV: set}, nil
}

func newLocalTokennote(payload *tokenty.TokennoteCreate) *tokenty.LocalTokennote {
	localTokennote := tokenty.LocalTokennote{
		Issuer:                payload.Issuer,
		IssuerPhone:           payload.IssuerPhone,
		IssuerName:            payload.IssuerName,
		IssuerId:		        payload.IssuerId,
		Acceptor:			    payload.Acceptor,
		AcceptanceDate:			    payload.AcceptanceDate,
		Rate:			    payload.Rate,
		Currency:              payload.Currency,
		Introduction:        payload.Introduction,
		Total:               payload.Balance,
		Status:              tokenty.TokennoteStatusCreated,
		CreatedHeight:       0,
		CreatedTime:         0,

	}
	return &localTokennote
}


func loadLocalTokennote(symbol string, status int32, db db.KVDB) (*tokenty.LocalTokennote, error) {
	key := calcTokennoteStatusKeyLocal(symbol,  status)
	v, err := db.Get(key)
	if err != nil {
		return nil, err
	}
	var localTokennote tokenty.LocalTokennote
	err = types.Decode(v, &localTokennote)
	if err != nil {
		return nil, err
	}
	return &localTokennote, nil
}

func setCreated(t *tokenty.LocalTokennote, height, time int64) *tokenty.LocalTokennote {
	t.CreatedTime = time
	t.CreatedHeight = height
	t.Status = tokenty.TokennoteStatusCreated
	return t
}


func resetCreated(t *tokenty.LocalTokennote) *tokenty.LocalTokennote {
	t.CreatedTime = 0
	t.CreatedHeight = 0
	t.Status = tokenty.TokennoteStatusPreCreated
	return t
}

