// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/33cn/chain33/types"
	tokenty "github.com/33cn/plugin/plugin/dapp/tokennote/types"
	"io/ioutil"
	"github.com/33cn/chain33/common/db/table"
	"github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/util"
)

func (t *tokennote) execDelLocal(receiptData *types.ReceiptData) ([]*types.KeyValue, error) {
	var set []*types.KeyValue
	for i := 0; i < len(receiptData.Logs); i++ {
		item := receiptData.Logs[i]
		if item.Ty == tokenty.TyLogTokennoteCreate || item.Ty == tokenty.TyLogTokennoteLoan || item.Ty == tokenty.TyLogTokennoteLoanedAgree || item.Ty == tokenty.TyLogTokennoteLoanedReject || item.Ty == tokenty.TyLogTokennoteCashed {
			var receipt tokenty.ReceiptTokennote
			err := types.Decode(item.Log, &receipt)
			if err != nil {
				tokennotelog.Error("Failed to decode ReceiptTokennote in ExecDelLocal")
				continue
			}
			set = append(set, t.deleteLogs(&receipt)...)
		}
	}
	return set, nil
}

func (t *tokennote) ExecDelLocal_Transfer(payload *types.AssetsTransfer, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set, err := t.ExecDelLocalLocalTransWithdraw(tx, receiptData, index)
	if err != nil {
		return nil, err
	}
	if cfg.SaveTokennoteTxList {
		tokenAction := tokenty.TokennoteAction{
			Ty: tokenty.ActionTransfer,
			Value: &tokenty.TokennoteAction_Transfer{
				Transfer: payload,
			},
		}
		kvs, err := t.makeTokennoteTxKvs(tx, &tokenAction, receiptData, index, true)
		if err != nil {
			return nil, err
		}
		set.KV = append(set.KV, kvs...)
	}
	return set, nil
}

func (t *tokennote) ExecDelLocal_Withdraw(payload *types.AssetsWithdraw, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set, err := t.ExecDelLocalLocalTransWithdraw(tx, receiptData, index)
	if err != nil {
		return nil, err
	}
	if cfg.SaveTokennoteTxList {
		tokenAction := tokenty.TokennoteAction{
			Ty: tokenty.ActionWithdraw,
			Value: &tokenty.TokennoteAction_Withdraw{
				Withdraw: payload,
			},
		}
		kvs, err := t.makeTokennoteTxKvs(tx, &tokenAction, receiptData, index, true)
		if err != nil {
			return nil, err
		}
		set.KV = append(set.KV, kvs...)
	}
	return set, nil
}

func (t *tokennote) ExecDelLocal_TransferToExec(payload *types.AssetsTransferToExec, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set, err := t.ExecDelLocalLocalTransWithdraw(tx, receiptData, index)
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
		kvs, err := t.makeTokennoteTxKvs(tx, &tokenAction, receiptData, index, true)
		if err != nil {
			return nil, err
		}
		set.KV = append(set.KV, kvs...)
	}
	return set, nil
}


func (t *tokennote) ExecDelLocal_TokennoteCreate(payload *tokenty.TokennoteCreate, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {

	localTokennote, err := loadLocalTokennote(payload.Currency, tokenty.TokennoteStatusCreated, t.GetLocalDB())
	if err != nil {
		return nil, err
	}
	localTokennote = resetCreated(localTokennote)
	key := calcTokennoteStatusKeyLocal(payload.Currency, tokenty.TokennoteStatusCreated)
	var set []*types.KeyValue

	set = append(set, &types.KeyValue{Key: key, Value: nil})

	tab ,err := table.NewTable(NewTokennoteMarketRow(),t.GetLocalDB(),&marketopt)
	if err != nil {
		tokennotelog.Error("ExecLocal_TokennoteCreate","new tokennotemarket init err ",err)
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
			tab.DelRow(TokennoteMarketRow{&t})
			ccount.Dec()
		}
	}
	tabkv,err := tab.Save()
	if err != nil {
		return nil,err
	}
	ccountkv ,err := ccount.Save()
	if err != nil {
		return nil,err
	}
	set = append(set,tabkv...)
	set = append(set,ccountkv...)
	return &types.LocalDBSet{KV: set}, nil
}

func (t *tokennote) ExecDelLocal_TokennoteLoan(payload *tokenty.TokennoteLoan, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {

	var set []*types.KeyValue

	//set = append(set, &types.KeyValue{Key: key, Value: nil})
	return &types.LocalDBSet{KV: set}, nil
}

func (t *tokennote) ExecDelLocal_TokennoteLoanedAgree(payload *tokenty.TokennoteLoanedAgree, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {

	var set []*types.KeyValue

	//set = append(set, &types.KeyValue{Key: key, Value: nil})
	return &types.LocalDBSet{KV: set}, nil
}

func (t *tokennote) ExecDelLocal_TokennoteLoanedReject(payload *tokenty.TokennoteLoanedReject, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {

	var set []*types.KeyValue

	//set = append(set, &types.KeyValue{Key: key, Value: nil})
	return &types.LocalDBSet{KV: set}, nil
}

func (t *tokennote) ExecDelLocal_TokennoteCashed(payload *tokenty.TokennoteCashed, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {

	var set []*types.KeyValue
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
					set = append(set,&types.KeyValue{Key:vvv,Value:nil})
				}
			}
		}
	}
	return &types.LocalDBSet{KV: set}, nil
}

func (t *tokennote) ExecDelLocal_TokennoteMint(payload *tokenty.TokennoteMint, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {

	dir, err := ioutil.TempDir("", "tokennotedb")
	if err != nil {
		panic(err)
	}
	ldb,err := db.NewGoLevelDB("tokennotedb",dir,128)
	if err != nil {
		return nil,err
	}
	var set []*types.KeyValue

	tab,err := table.NewTable(NewTransactionRow(),t.GetLocalDB(),&opt)
	if err != nil {
		return nil,err
	}
	tab.DelRow(tx)
	kvs,err := tab.Save()
	if err != nil {
		return nil,err
	}
	util.SaveKVList(ldb,kvs)

	//set = append(set, &types.KeyValue{Key: key, Value: nil})
	return &types.LocalDBSet{KV: set}, nil
}

func (t *tokennote) ExecDelLocal_TokennoteBurn(payload *tokenty.TokennoteBurn, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {

	dir, err := ioutil.TempDir("", "tokennotedb")
	if err != nil {
		panic(err)
	}
	ldb,err := db.NewGoLevelDB("tokennotedb",dir,128)
	if err != nil {
		return nil,err
	}
	var set []*types.KeyValue

	tab,err := table.NewTable(NewTransactionRow(),t.GetLocalDB(),&opt)
	if err != nil {
		return nil,err
	}
	tab.DelRow(tx)
	kvs,err := tab.Save()
	if err != nil {
		return nil,err
	}
	util.SaveKVList(ldb,kvs)

	//set = append(set, &types.KeyValue{Key: key, Value: nil})
	return &types.LocalDBSet{KV: set}, nil
}