// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

/*
tokennote执行器支持tokennote的创建，

主要提供操作有以下几种：
1）预创建tokennote；
2）完成创建tokennote
3）撤销预创建
*/

import (
	"github.com/33cn/chain33/account"
	"github.com/33cn/chain33/common/address"
	log "github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/chain33/system/dapp"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	tokennotety "github.com/33cn/plugin/plugin/dapp/tokennote/types"
	"github.com/pkg/errors"
)

var tokennotelog = log.New("module", "execs.tokennote")

const (
	finisherKey       = "tokennote-finisher"
	tokennoteAssetsPrefix = "LODB-tokennote-assets:"
	blacklist         = "tokennote-blacklist"
)

var driverName = "tokennote"
var conf = types.ConfSub(driverName)

func init() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&tokennote{}))
}

type subConfig struct {
	SaveTokennoteTxList bool `json:"saveTokennoteTxList"`
}

var cfg subConfig

// Init 重命名执行器名称
func Init(name string, sub []byte) {
	if sub != nil {
		types.MustDecode(sub, &cfg)
	}
	drivers.Register(GetName(), newTokennote, types.GetDappFork(driverName, "Enable"))
}

// GetName 获取执行器别名
func GetName() string {
	return newTokennote().GetName()
}

type tokennote struct {
	drivers.DriverBase
}

func newTokennote() drivers.Driver {
	t := &tokennote{}
	t.SetChild(t)
	t.SetExecutorType(types.LoadExecutorType(driverName))
	return t
}

// GetDriverName 获取执行器名字
func (t *tokennote) GetDriverName() string {
	return driverName
}

// CheckTx ...
func (t *tokennote) CheckTx(tx *types.Transaction, index int) error {
	return nil
}

func (t *tokennote) queryTokennoteAssetsTimeKey(addr string) (*tokennotety.ReplyAccountTokennoteList, error) {
	key := calcTokennoteAssetsTimeKey(addr)
	value, err := t.GetLocalDB().Get(key)
	if value == nil || err != nil {
		tokennotelog.Error("tokennotedb", "GetTokennoteAssetstimeKey", types.ErrNotFound)
		return nil, types.ErrNotFound
	}
	var assets tokennotety.ReplyAccountTokennoteList
	err = types.Decode(value, &assets)
	if err != nil {
		tokennotelog.Error("tokennotedb", "GetTokennoteAssetstimeKey", err)
		return nil, err
	}
	tokennotelog.Error("queryTokennoteAssetsTimeKey:","key list",assets)
	return &assets, nil
}

//资产
func (t *tokennote) getAccountTokennoteAssets(req *tokennotety.ReqAccountTokennoteAssets) (types.Message, error) {
	var reply = &tokennotety.ReplyAccountTokennoteAssets{}
	assets, err := getTokennoteAssetsKey(req.Address,t.GetLocalDB())
	if err != nil {
		return nil, err
	}
	for _, asset := range assets.Datas {
		//判断是否可以转账
		tokennotedb,err := t.getTokennoteDB(asset)
		if err != nil {
			return nil,err
		}
		if tokennotedb.Status == tokennotety.TokennoteStatusPreCreated {
			//tokennotelog.Error("getTokennoteStatus","status",tokennotedb.Status)
			continue
		}
		acc, err := account.NewAccountDB(t.GetName(),asset, t.GetStateDB())
		if err != nil {
			return nil, err
		}
		var acc1 *types.Account
		if req.Execer == "trade" {
			execaddress := address.ExecAddress(req.Execer)
			acc1 = acc.LoadExecAccount(req.Address, execaddress)
		} else if req.Execer == t.GetName() {
			acc1 = acc.LoadAccount(req.Address)
			//tokennotelog.Error("query:","acc1:",acc1)
		}
		if acc1 == nil {
			continue
		}
		var repay, amount,loanTime int64
		var status int32
		if tokennotedb.Issuer == req.Address {
			repay = tokennotedb.Repayamount
			amount = tokennotedb.Balance
			status = tokennotedb.Status
		} else {
			if asset == tokennotety.TokenCCNY {//CNYY直接显示余额 暂时不显示
				continue
			} else {
				hold,err:= t.getTokenHoldDetail(asset,req.Address)
				if err != nil && err == types.ErrNotFound {
					continue
				} else if err != nil && err != types.ErrNotFound {
					return nil,err
				}
				repay = hold.Repayamount
				amount = hold.Amount
				status = hold.Status
				loanTime = hold.LoanTime
			}
		}

		tokennoteAsset := &tokennotety.TokennoteAsset{
			Symbol: asset,
			Creator:tokennotedb.Issuer,
			IssuerName:tokennotedb.IssuerName,
			SingleAmount:amount,
			Repay:repay,
			Account: acc1,
			AcceptanceDate:tokennotedb.AcceptanceDate,
			IssuerPhone:tokennotedb.IssuerPhone,
			LoanTime:loanTime,
			Status:status,
		}
		reply.TokennoteAssets = append(reply.TokennoteAssets, tokennoteAsset)
	}
	return reply, nil
}

func (t *tokennote) getTokenHoldDetail(currency ,addr  string) (*tokennotety.TokennoteHold,error) {
	var hold tokennotety.TokennoteHold
	value,err := t.GetStateDB().Get(calcTokennoteHoldKeyNew(currency,addr))
	//tokennotelog.Error("getTokenHoldDetail","HOLDKEY",string(calcTokennoteHoldKeyNew(currency,addr)))
	if err != nil {
		tokennotelog.Error("getTokenHoldDetail","hold data not exists",err)
		return nil,err
	}
	err = types.Decode(value,&hold)
	if err != nil {
		tokennotelog.Error("getTokenHoldDetail","hold data decode err",err)
		return nil,err
	}
	//tokennotelog.Error("getTokenHoldDetail","HOLDValue",hold)
	return &hold,nil
}

func (t *tokennote) getAddrReceiverforTokennotes(addrTokennotes *tokennotety.ReqAddrTokennotes) (types.Message, error) {
	var reply = &tokennotety.ReplyAddrRecvForTokennotes{}
	db := t.GetLocalDB()
	reciver := types.Int64{}
	for _, tokennote := range addrTokennotes.Token {
		addrRecv, err := db.Get(calcAddrKey(tokennote, addrTokennotes.Addr))
		if addrRecv == nil || err != nil {
			continue
		}
		err = types.Decode(addrRecv, &reciver)
		if err != nil {
			continue
		}

		recv := &tokennotety.TokennoteRecv{Tokennote: tokennote, Recv: reciver.Data}
		reply.TokennoteRecvs = append(reply.TokennoteRecvs, recv)
	}

	return reply, nil
}

func (t *tokennote) getTokennoteInfo(symbol string) (types.Message, error) {
	if symbol == "" {
		return nil, types.ErrInvalidParam
	}
	key := calcTokennoteStatusKeyLocal( symbol,tokennotety.TokennoteStatusCreated)
	values, err := t.GetLocalDB().Get(key)
	if err != nil {
		return nil, err
	}

	var tokennoteInfo tokennotety.LocalTokennote
	err = types.Decode(values, &tokennoteInfo)
	if err != nil {
		return &tokennoteInfo, err
	}
	return &tokennoteInfo, nil
}

func (t *tokennote) getTokennotes(reqTokennotes *tokennotety.ReqTokennotes) (types.Message, error) {
	replyTokennotes := &tokennotety.ReplyTokennotes{}
	tokennotes, err := t.listTokennoteKeys(reqTokennotes)
	if err != nil {
		return nil, err
	}
	tokennotelog.Error("tokennote Query GetTokennotes", "get count", len(tokennotes))
	if reqTokennotes.SymbolOnly {
		for _, t1 := range tokennotes {
			if len(t1) == 0 {
				continue
			}

			var tokennoteValue tokennotety.LocalTokennote
			err = types.Decode(t1, &tokennoteValue)
			if err == nil {
				tokennote := tokennotety.LocalTokennote{Currency: tokennoteValue.Currency}
				replyTokennotes.Tokens = append(replyTokennotes.Tokens, &tokennote)
			}
		}
		return replyTokennotes, nil
	}

	for _, t1 := range tokennotes {
		// delete impl by set nil
		if len(t1) == 0 {
			continue
		}

		var tokennote tokennotety.LocalTokennote
		err = types.Decode(t1, &tokennote)
		if err == nil {
			replyTokennotes.Tokens = append(replyTokennotes.Tokens, &tokennote)
		}
	}

	//tokennotelog.Info("tokennote Query", "replyTokennotes", replyTokennotes)
	return replyTokennotes, nil
}

func (t *tokennote) listTokennoteKeys(reqTokennotes *tokennotety.ReqTokennotes) ([][]byte, error) {
	querydb := t.GetLocalDB()
	if reqTokennotes.QueryAll {
		keys, err := querydb.List(calcTokennoteStatusKeyPrefixLocal(reqTokennotes.Status), nil, 0, 0)
		if err != nil && err != types.ErrNotFound {
			return nil, err
		}
		if len(keys) == 0 {
			return nil, types.ErrNotFound
		}
		tokennotelog.Debug("tokennote Query GetTokennotes", "get count", len(keys))
		return keys, nil
	}
	var keys [][]byte
	for _, tokennote := range reqTokennotes.Tokennotes {
		keys1, err := querydb.List(calcTokennoteStatusTokennoteKeyPrefixLocal(reqTokennotes.Status, tokennote), nil, 0, 0)
		if err != nil && err != types.ErrNotFound {
			return nil, err
		}
		keys = append(keys, keys1...)

		tokennotelog.Debug("tokennote Query GetTokennotes", "get count", len(keys))
	}
	if len(keys) == 0 {
		return nil, types.ErrNotFound
	}
	return keys, nil
}

// value 对应 statedb 的key
func (t *tokennote) saveLogs(receipt *tokennotety.ReceiptTokennote) []*types.KeyValue {
	var kv []*types.KeyValue

	key := calcTokennoteStatusKeyLocal(receipt.Symbol,  receipt.Status)
	var value []byte

	value = calcTokennoteAddrNewKeyS(receipt.Symbol, receipt.Owner)

	kv = append(kv, &types.KeyValue{Key: key, Value: value})
	//如果当前需要被更新的状态不是Status_PreCreated，则认为之前的状态是precreate，且其对应的key需要被删除
	if receipt.Status != tokennotety.TokennoteStatusPreCreated {
		key = calcTokennoteStatusKeyLocal(receipt.Symbol,  tokennotety.TokennoteStatusPreCreated)
		kv = append(kv, &types.KeyValue{Key: key, Value: nil})
	}
	return kv
}

func (t *tokennote) deleteLogs(receipt *tokennotety.ReceiptTokennote) []*types.KeyValue {
	var kv []*types.KeyValue

	key := calcTokennoteStatusKeyLocal(receipt.Symbol,  receipt.Status)
	kv = append(kv, &types.KeyValue{Key: key, Value: nil})
	//如果当前需要被更新的状态不是Status_PreCreated，则认为之前的状态是precreate，且其对应的key需要被恢复
	if receipt.Status != tokennotety.TokennoteStatusPreCreated {
		key = calcTokennoteStatusKeyLocal(receipt.Symbol,  tokennotety.TokennoteStatusPreCreated)
		var value []byte
		value = calcTokennoteAddrNewKeyS(receipt.Symbol, receipt.Owner)
		kv = append(kv, &types.KeyValue{Key: key, Value: value})
	}
	return kv
}

func (t *tokennote) makeTokennoteTxKvs(tx *types.Transaction, action *tokennotety.TokennoteAction, receipt *types.ReceiptData, index int, isDel bool) ([]*types.KeyValue, error) {
	var kvs []*types.KeyValue
	var symbol string
	if action.Ty == tokennotety.ActionTransfer {
		symbol = action.GetTransfer().Cointoken
	} else if action.Ty == tokennotety.ActionWithdraw {
		symbol = action.GetWithdraw().Cointoken
	} else if action.Ty == tokennotety.TokennoteActionTransferToExec {
		symbol = action.GetTransferToExec().Cointoken
	} else {
		return kvs, nil
	}

	kvs, err := tokenTxKvs(tx, symbol, t.GetHeight(), int64(index), isDel)
	return kvs, err
}

func findTokennoteTxListUtil(req *tokennotety.ReqTokennoteTx) ([]byte, []byte) {
	var key, prefix []byte
	if len(req.Addr) > 0 {
		if req.Flag == 0 {
			prefix = calcTokennoteAddrTxKey(req.Symbol, req.Addr, -1, 0)
			key = calcTokennoteAddrTxKey(req.Symbol, req.Addr, req.Height, req.Index)
		} else {
			prefix = calcTokennoteAddrTxDirKey(req.Symbol, req.Addr, req.Flag, -1, 0)
			key = calcTokennoteAddrTxDirKey(req.Symbol, req.Addr, req.Flag, req.Height, req.Index)
		}
	} else {
		prefix = calcTokennoteTxKey(req.Symbol, -1, 0)
		key = calcTokennoteTxKey(req.Symbol, req.Height, req.Index)
	}
	if req.Height == -1 {
		key = nil
	}
	return key, prefix
}

func findTokennoteCashedTxListUtil(req *tokennotety.ReqTokennoteTx) ([]byte, []byte) {
	var key, prefix []byte
	if len(req.Addr) > 0 {
		if len(req.Symbol) > 0 {
			if req.Flag == 0 {
				prefix = calcCashedTokennoteAddrTxKey(req.Symbol, req.Addr, -1, 0)
				key = calcCashedTokennoteAddrTxKey(req.Symbol, req.Addr, req.Height, req.Index)
			} else {
				prefix = calcCashedTokennoteAddrTxDirKey(req.Symbol, req.Addr, req.Flag, -1, 0)
				key = calcCashedTokennoteAddrTxDirKey(req.Symbol, req.Addr, req.Flag, req.Height, req.Index)
			}
		} else {
			prefix = calcCashedTokennoteTxKey(req.Addr, -1, 0)
			key = calcCashedTokennoteTxKey(req.Addr, req.Height, req.Index)
		}

	} else {
		prefix = calcCashedTokennoteTxKey(req.Symbol, -1, 0)
		key = calcCashedTokennoteTxKey(req.Symbol, req.Height, req.Index)
	}
	if req.Height == -1 {
		key = nil
	}
	return key, prefix
}

func (t *tokennote) getTxByTokennote(req *tokennotety.ReqTokennoteTx) (types.Message, error) {
	if req.Flag != 0 && req.Flag != dapp.TxIndexFrom && req.Flag != dapp.TxIndexTo {
		err := types.ErrInvalidParam
		return nil, errors.Wrap(err, "flag unknown")
	}
	key, prefix := findTokennoteTxListUtil(req)
	tokennotelog.Debug("GetTxByTokennote", "key", string(key), "prefix", string(prefix))

	db := t.GetLocalDB()
	txinfos, err := db.List(prefix, key, req.Count, req.Direction)
	if err != nil {
		return nil, errors.Wrap(err, "db.List to find tokennote tx list")
	}
	if len(txinfos) == 0 {
		return nil, errors.Wrapf(types.ErrNotFound, "key=%s, prefix=%s", string(key), string(prefix))
	}

	var replyTxInfos types.ReplyTxInfos
	replyTxInfos.TxInfos = make([]*types.ReplyTxInfo, len(txinfos))
	for index, txinfobyte := range txinfos {
		var replyTxInfo types.ReplyTxInfo
		err := types.Decode(txinfobyte, &replyTxInfo)
		if err != nil {
			return nil, err
		}
		replyTxInfos.TxInfos[index] = &replyTxInfo
	}
	return &replyTxInfos, nil
}

//获取地址的还款记录
func (t *tokennote) getCashedTxByTokennote(req *tokennotety.ReqTokennoteTx) (types.Message, error) {
	if req.Flag != 0 && req.Flag != dapp.TxIndexFrom && req.Flag != dapp.TxIndexTo {
		err := types.ErrInvalidParam
		return nil, errors.Wrap(err, "flag unknown")
	}
	key, prefix := findTokennoteCashedTxListUtil(req)
	tokennotelog.Debug("getCashedTxByTokennote", "key", string(key), "prefix", string(prefix))

	db := t.GetLocalDB()
	txinfos, err := db.List(prefix, key, req.Count, req.Direction)
	if err != nil {
		return nil, errors.Wrap(err, "db.List to find tokennote cased tx list")
	}
	if len(txinfos) == 0 {
		return nil, errors.Wrapf(types.ErrNotFound, "key=%s, prefix=%s", string(key), string(prefix))
	}

	var replyTxInfos tokennotety.ReceiptTokennoteCashed
	replyTxInfos.Cashlist = make([]*tokennotety.TokennoteCashDetail, len(txinfos))
	for index, txinfobyte := range txinfos {
		var replyTxInfo tokennotety.TokennoteCashDetail
		err := types.Decode(txinfobyte, &replyTxInfo)
		if err != nil {
			return nil, err
		}
		replyTxInfos.Cashlist[index] = &replyTxInfo
	}
	return &replyTxInfos, nil
}

func (t *tokennote) getAgreeTxByTokennote (req *types.ReplyString) (types.Message, error) {
	tokennote ,err := t.getTokennoteDB(req.Data)
	if err != nil {
		return nil,err
	}
	return tokennote,nil
}

// CheckReceiptExecOk return true to check if receipt ty is ok
func (t *tokennote) CheckReceiptExecOk() bool {
	return true
}
