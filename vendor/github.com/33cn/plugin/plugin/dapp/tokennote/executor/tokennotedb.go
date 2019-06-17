// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"fmt"

	"strings"

	"github.com/33cn/chain33/account"
	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	pty "github.com/33cn/plugin/plugin/dapp/tokennote/types"
	"github.com/33cn/chain33/common/address"
	"math"
	"time"
	"math/big"
	"strconv"
	"github.com/33cn/chain33/common"
)

type tokennoteDB struct {
	token pty.Tokennote
}

func newTokennoteDB(create *pty.TokennoteCreate, creator string,createTime int64) *tokennoteDB {
	t := &tokennoteDB{}
	t.token.Issuer = creator
	t.token.IssuerName = create.IssuerName
	t.token.IssuerPhone = create.IssuerPhone
	t.token.IssuerId = create.IssuerId
	t.token.Acceptor = creator
	t.token.AcceptanceDate = create.AcceptanceDate
	t.token.Rate = create.Rate
	if create.Rate > 0 {
		t.token.OverdueRate = create.Rate
	} else {
		t.token.OverdueRate = int64(pty.TokennoteOverdueDayRate)
	}
	t.token.Currency = create.GetCurrency()
	t.token.Introduction = create.GetIntroduction()
	t.token.Balance = create.Balance
	t.token.Total = create.GetBalance()
	//t.token.Repayamount = create.Repayamount
	t.token.CreateTime =createTime
	return t
}

func (t *tokennoteDB) save(db dbm.KV, key []byte) {
	set := t.getKVSet(key)
	for i := 0; i < len(set); i++ {
		db.Set(set[i].GetKey(), set[i].Value)
	}
}


func (action *tokennoteAction) saveStateDB(kv []*types.KeyValue) {
	for i:=0; i<len(kv) ; i++  {
		action.db.Set(kv[i].Key,kv[i].Value)
	}
}

func (t *tokennoteDB) getLogs(ty int32, status int32) []*types.ReceiptLog {
	var log []*types.ReceiptLog
	value := types.Encode(&pty.ReceiptTokennote{Symbol: t.token.Currency, Owner: t.token.Issuer, Status: t.token.Status})
	log = append(log, &types.ReceiptLog{Ty: ty, Log: value})

	return log
}

//key:mavl-create-token-addr-xxx or mavl-token-xxx <-----> value:token
func (t *tokennoteDB) getKVSet(key []byte) (kvset []*types.KeyValue) {
	value := types.Encode(&t.token)
	kvset = append(kvset, &types.KeyValue{Key: key, Value: value})
	return kvset
}


func getTokennoteFromDB(db dbm.KV, symbol string, owner string) (*pty.Tokennote, error) {
	key := calcTokennoteAddrNewKeyS(symbol, owner)
	value, err := db.Get(key)
	if err != nil {
		return nil, err
	}
	var token pty.Tokennote
	if err = types.Decode(value, &token); err != nil {
		tokennotelog.Error("getTokennoteFromDB", "Fail to decode types.token for key", string(key), "err info is", err)
		return nil, err
	}
	return &token, nil
}

func getContract(db dbm.KV, symbol string) (*pty.TokennoteContract,error) {
	key := calcTokennoteContractKey(symbol)
	value, err := db.Get(key)
	if err != nil {
		return nil, err
	}
	var contract pty.TokennoteContract
	if err = types.Decode(value, &contract); err != nil {
		tokennotelog.Error("getTokennoteFromDB", "Fail to decode types.token for key", string(key), "err info is", err)
		return nil, err
	}
	return &contract, nil
}

type tokennoteAction struct {
	coinsAccount *account.DB
	db           dbm.KV
	txhash       []byte
	fromaddr     string
	toaddr       string
	blocktime    int64
	height       int64
	execaddr     string
}

func newTokennoteAction(t *tokennote, toaddr string, tx *types.Transaction) *tokennoteAction {
	hash := tx.Hash()
	fromaddr := tx.From()
	return &tokennoteAction{t.GetCoinsAccount(), t.GetStateDB(), hash, fromaddr, toaddr,
		t.GetBlockTime(), t.GetHeight(), dapp.ExecAddress(string(tx.Execer))}
}

func (action *tokennoteAction) create(token *pty.TokennoteCreate) (*types.Receipt, error) {
	tokennotelog.Debug("create note")
	if token == nil {
		return nil, types.ErrInvalidParam
	}
	if token.Issuer == "" || token.Acceptor == "" {
		return nil,types.ErrInvalidParam
	}
	//检查地址是否符合格式
	e := address.CheckAddress(token.Acceptor)
	if e != nil {
		return nil,e
	}


	if token.Rate > int64(pty.TokennoteRateMax) || token.Rate < 0 || token.AcceptanceDate < 1000000000 || token.AcceptanceDate > 10000000000 {
		return nil,pty.ErrTokennoteNumberFormat
	}

	if len(token.GetIssuerName()) > pty.TokennoteNameLenLimit {
		return nil, pty.ErrTokennoteNameLen
	} else if len(token.GetIntroduction()) > pty.TokennoteIntroLenLimit {
		return nil, pty.ErrTokennoteIntroLen
	} else if len(token.GetCurrency()) > pty.TokennoteSymbolLenLimit || len(token.GetCurrency()) <= 0 {
		return nil, pty.ErrTokennoteSymbolLen
	} else if token.GetBalance() > types.MaxTokenBalance || token.GetBalance() <= 0 {
		return nil, pty.ErrTokennoteTotalOverflow
	}
	if action.blocktime >= token.AcceptanceDate {
		tokennotelog.Error("create","token AcceptanceDate should >= action.blocktime",token.AcceptanceDate)
		return nil,pty.ErrTokennoteExcept
	}
	//
	if len(token.GetCurrency()) <= pty.TokennoteSymbolAdminLimit {
		approverValid := false

		for _, approver := range conf.GStrList("tokennoteApprs") {
			if approver == action.fromaddr {
				approverValid = true
				break
			}
		}

		hasPriv, ok := validFinisher(action.fromaddr, action.db)
		if (ok != nil || !hasPriv) && !approverValid {
			return nil, pty.ErrTokennoteCreatedNotAllowed
		}
	} else {//特殊tokennote  需要管理员注册
		if token.Issuer != action.fromaddr {
			return nil,pty.ErrTokennoteOwner
		}
	}




	if !validSymbol([]byte(token.GetCurrency())) {
		tokennotelog.Error("token precreate ", "symbol need be upper", token.GetCurrency())
		return nil, pty.ErrTokennoteSymbolUpper
	}

	if checkTokenExist(token.GetCurrency(), action.db) {
		return nil, pty.ErrTokennoteExist
	}


	found, err := inBlacklist(token.GetCurrency(), blacklist, action.db)
	if err != nil {
		return nil, err
	}
	if found {
		return nil, pty.ErrTokennoteBlacklist
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue


	tokendb := newTokennoteDB(token, token.Issuer,action.blocktime)
	day := getSubDays(action.blocktime,tokendb.token.AcceptanceDate)
	repay,err := getCalcAmount(token.Balance,tokendb.token.Rate,day,1)
	if err != nil {
		tokennotelog.Error("loan","calc err",err)
		return nil,err
	}
	if repay < token.Balance {
		tokennotelog.Error("create ","loan repay err ",repay)
	}
	tokennotelog.Error("create ","loan repay",repay)
	tokendb.token.Repayamount = repay

	if token.Currency == pty.TokenCCNY {
		tokendb.token.Repayamount = token.Balance
		tokendb.token.Rate = 0
	}
	var statuskey []byte
	var key []byte

	statuskey = calcTokennoteStatusNewKeyS(tokendb.token.Currency, tokendb.token.Issuer, pty.TokennoteStatusCreated)
	key = calcTokennoteAddrNewKeyS(tokendb.token.Currency, tokendb.token.Issuer)
	tokendb.token.Status = pty.TokennoteStatusCreated

	//非定向，上市场

	marketLogs := &types.ReceiptLog{Ty:pty.TyLogTokennoteMarket,Log:types.Encode(&tokendb.token)}
	logs = append(logs,marketLogs)

	tokendb.save(action.db, statuskey)
	tokendb.save(action.db, key)

	logs = append(logs, tokendb.getLogs(pty.TyLogTokennoteCreate, pty.TokennoteStatusCreated)...)
	kv = append(kv, tokendb.getKVSet(key)...)
	kv = append(kv, tokendb.getKVSet(statuskey)...)

	key = calcTokennoteKey(tokendb.token.Currency)
	//因为该token已经被创建，需要保存一个全局的token，防止其他用户再次创建
	tokendb.save(action.db, key)
	kv = append(kv, tokendb.getKVSet(key)...)

	//创建token类型的账户，同时需要创建的额度存入

	tokennoteAccount, err := account.NewAccountDB(pty.TokennoteX, token.GetCurrency(), action.db)
	if err != nil {
		return nil, err
	}
	tokennotelog.Debug("create", "tokennote.Owner", token.Issuer, "token.GetTotal()", token.GetBalance())
	receiptForToken, err := tokennoteAccount.GenesisInit(token.Issuer, token.GetBalance())
	if err != nil {
		return nil, err
	}

	logs = append(logs, receiptForToken.Logs...)
	kv = append(kv, receiptForToken.KV...)

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	for _,v := range kv {
		tokennotelog.Error("kv:","key ",string(v.Key))
	}
	return receipt, nil
}

func (action *tokennoteAction) loan(loan *pty.TokennoteLoan) (*types.Receipt, error) {
	tokennotelog.Error("loan")
	if loan == nil {
		return nil, types.ErrInvalidParam
	}

	e := address.CheckAddress(loan.GetTo())
	if e != nil {
		return nil,e
	}

	tokennote, err := getTokennoteFromDB(action.db, loan.GetSymbol(), action.fromaddr)
	if err != nil || tokennote.Status == pty.TokennoteStatusPreCreated {
		return nil, pty.ErrTokennoteNotCreated
	}
	if tokennote.Status == pty.TokennoteStatusCashed {
		return nil,pty.ErrTokennoteCashed
	}
	if tokennote.Issuer != action.fromaddr {
		return nil,pty.ErrTokennoteOwner
	}
	if tokennote.Issuer == loan.To {
		return nil,pty.ErrTokennoteNotLoanToSelf
	}

	if (action.blocktime - tokennote.CreateTime) > int64(pty.TokennoteExpireTime) {
		return nil,pty.ErrTokennoteOverdue
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue


	tokendb := &tokennoteDB{*tokennote}
	var key []byte

	key = calcTokennoteAddrNewKeyS(tokendb.token.Currency, tokendb.token.Issuer)

	tokendb.token.Balance -= loan.Amount
	if tokendb.token.Balance < 0 {
		return nil,pty.ErrTokennoteAmountLow
	}

	tokenoteAccount ,err := account.NewAccountDB("tokennote",loan.Symbol,action.db)
	if err != nil {
		return nil,err
	}

	//借款人白条转移到合约子账户
	receiptBorrow,err := tokenoteAccount.TransferToExec(action.fromaddr,action.execaddr,loan.Amount)
	if err != nil {
		tokennotelog.Error("loan ","借款人白条资产转移至合约子账户出现错误",err)
		return nil,err
	}
	//执行冻结
	receiptBorrowExec,err := tokenoteAccount.ExecFrozen(action.fromaddr,action.execaddr,loan.Amount)
	if err != nil {
		tokennotelog.Error("loan ","冻结子账户 白条出现异常",err)
		return nil,err
	}
	hold := &pty.TokennoteHold{Addr:loan.To,Amount:loan.Amount,LoanTime:action.blocktime,Currency:loan.GetSymbol(),Status:pty.TokennoteStatusReadyLoan}

	//借款to key
	//key1 := calcTokennoteHoldKey(loan.Symbol,loan.To,action.blocktime)
	keynew := calcTokennoteHoldKeyNew(loan.Symbol,loan.To)
	//同一个白条只能借一次
	kvalue,err := action.db.Get(keynew)
	if err != nil && err != types.ErrNotFound {
		tokennotelog.Error("loan ","checkhold key exists",err)
		return nil,err
	}
	if kvalue != nil {
		return nil,pty.ErrTokennoteNotAllowedLoanTwice
	}
	bigamount := big.NewInt(loan.Amount)
	bigtotal := big.NewInt(tokendb.token.Total)
	bigrepay := big.NewInt(tokendb.token.Repayamount)
	base := big.NewInt(int64(1))

	repay :=base.Mul(bigamount,bigrepay).Div(base,bigtotal).Int64()
	tokennotelog.Error("loan","repay amount ",repay)
	if repay < loan.Amount {
		return nil,pty.ErrTokennoteExcept
	}

	hold.Repayamount = repay
	hold.Creator = action.fromaddr
	tokendb.token.Holds = append(tokendb.token.Holds,hold)
	value1 := types.Encode(hold)

	tokendb.save(action.db,key)
	kv = append(kv,&types.KeyValue{Key:keynew,Value:value1})

	logs = append(logs, receiptBorrow.Logs...)
	logs = append(logs, receiptBorrowExec.Logs...)
	logs = append(logs, tokendb.getLogs(pty.TyLogTokennoteLoan, pty.TokennoteStatusCreated)...)
	kv = append(kv, receiptBorrow.KV...)
	kv = append(kv, receiptBorrowExec.KV...)
	kv = append(kv, tokendb.getKVSet(key)...)

	//保存合同信息
	var contract pty.TokennoteContract
	*contract.Tokennote = tokendb.token
	contract.BorrowHash = common.ToHex(action.txhash)
	contractKey := calcTokennoteContractKey(loan.Symbol)

	kv = append(kv,&types.KeyValue{Key:contractKey,Value:types.Encode(&contract)})

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	for _,v := range kv {
		tokennotelog.Error("kv:","key ",string(v.Key))
	}
	return receipt, nil
}

func (action *tokennoteAction) loanedAgree(loan *pty.TokennoteLoanedAgree) (*types.Receipt, error) {
	tokennotelog.Error("loanAgree")
	if loan == nil {
		return nil, types.ErrInvalidParam
	}
	tokennote, err := getTokennoteFromDB(action.db, loan.GetSymbol(), loan.GetOwner())
	if err != nil {
		tokennotelog.Error("token loanedConfirm  ", "Can't get token form db for token", loan.GetSymbol())
		return nil, pty.ErrTokennoteNotCreated
	}

	if tokennote.Status != pty.TokennoteStatusCreated {
		tokennotelog.Error("token loanedConfirm ", "token's status should be created to be confirm", loan.GetSymbol())
		return nil, pty.ErrTokennoteNotCreated
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	tokennoteAccount ,err := account.NewAccountDB(pty.TokennoteX,loan.Symbol,action.db)
	if err != nil {
		return nil,err
	}
	tokenCNYYAccount ,err := account.NewAccountDB(pty.TokennoteX,pty.TokenCCNY,action.db)
	if err != nil {
		return nil,err
	}

	tokendb := &tokennoteDB{*tokennote}
	var key []byte
	key = calcTokennoteAddrNewKeyS(tokendb.token.Currency,tokendb.token.Issuer)

	//确认放款key
	key2 := calcTokennoteHoldKeyNew(loan.Symbol,action.fromaddr)
	v,err := action.db.Get(key2)
	if err != nil {
		tokennotelog.Error("loan agree"," note hold key  not exists",err)
		return nil,err
	}
	var hold pty.TokennoteHold
	err = types.Decode(v,&hold)
	if err != nil {
		tokennotelog.Error("loanedConfirm","decode TokennoteHold err",key2)
		return nil,err
	}
	if hold.Status != pty.TokennoteStatusReadyLoan {
		return nil,pty.ErrTokennoteNotReadyLoan
	}

	//day := getSubDays(action.blocktime,tokendb.token.AcceptanceDate)
	//repay,err := getCalcAmount(hold.Amount,tokendb.token.Rate,day)
	//if err != nil {
	//	tokennotelog.Error("loan","calc err",err)
	//	return nil,err
	//}
	//tokennotelog.Error("loan agree ","repay amount ",repay)
	//hold.Repayamount = repay
	hold.Status = pty.TokennoteStatusAgreeLoan
	hold.LoanTime = action.blocktime
	value1 := types.Encode(&hold)

	kv = append(kv,&types.KeyValue{Key:key2,Value:value1})

	c := false
	tokennotelog.Error("loanedAgree","holds:",tokendb.token.Holds)
	for _,v := range tokendb.token.Holds {
		if  v.Addr == action.fromaddr && v.Status == pty.TokennoteStatusReadyLoan {
			v.LoanTime = action.blocktime
			//v.Repayamount = repay
			v.Status = pty.TokennoteStatusAgreeLoan
			c = true
			break
		}
	}
	if !c {
		return nil,pty.ErrTokennoteExcept
	}
	tokendb.save(action.db,key)
	kv = append(kv,tokendb.getKVSet(key)...)
	accfrom := tokenCNYYAccount.LoadAccount(action.fromaddr)
	tokennotelog.Error("loanedAgree","accfrom",accfrom)
	//出借人CNYY资产转移给借款人
	receiptLend,err := tokenCNYYAccount.Transfer(action.fromaddr,loan.Owner,hold.Amount)
	if err != nil {
		tokennotelog.Error("loanedAgree","出借人CNYY资产转移异常 ",err)
		return nil,err
	}
	//出借人获取白条
	receiptLendNote,err := tokennoteAccount.ExecTransferFrozen(loan.Owner,action.fromaddr,action.execaddr,hold.Amount)
	if err != nil {
		tokennotelog.Error("loanedAgree","出借人获取白条异常 ",err)
		return nil,err
	}

	logs = append(logs, receiptLend.Logs...)
	logs = append(logs, receiptLendNote.Logs...)
	logs = append(logs, tokendb.getLogs(pty.TyLogTokennoteLoanedAgree, pty.TokennoteStatusCreated)...)
	kv = append(kv, receiptLend.KV...)
	kv = append(kv, receiptLendNote.KV...)

	//更新合同相关信息
	contract,error := getContract(action.db,loan.Symbol)
	if error != nil {
		tokennotelog.Error("loanedAgree","get contract err",error)
		return nil,pty.ErrTokennoteExcept
	}
	*contract.Tokennote = tokendb.token
	contract.LoanHash = common.ToHex(action.txhash)

	kv = append(kv,&types.KeyValue{Key:calcTokennoteContractKey(loan.Symbol),Value:types.Encode(contract)})

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	for _,v := range kv {
		tokennotelog.Error("kv:","key ",string(v.Key))
	}
	return receipt, nil
}

func (action *tokennoteAction) loanedReject(loan *pty.TokennoteLoanedReject) (*types.Receipt, error) {
	tokennotelog.Error("loanReject")
	if loan == nil {
		return nil,types.ErrInvalidParam
	}
	tokennote, err := getTokennoteFromDB(action.db, loan.GetSymbol(), loan.GetOwner())
	if err != nil {
		tokennotelog.Error("loanedReject","getTokennoteFromDB err  ",err)
		return nil,err
	}
	tokennotedb := &tokennoteDB{token:*tokennote}

	tokennoteAccount ,err := account.NewAccountDB(pty.TokennoteX,loan.Symbol,action.db)
	if err != nil {
		return nil,err
	}
	//key1 := calcTokennoteHoldKey(loan.GetSymbol(),action.fromaddr,loan.LoanTime)
	key_new := calcTokennoteHoldKeyNew(loan.Symbol,action.fromaddr)
	v,err := action.db.Get(key_new)
	if err != nil {
		tokennotelog.Error("loanedReject","get user hold key not exists",err)
		return nil,err
	}
	var hold pty.TokennoteHold

	err = types.Decode(v,&hold)
	if err != nil {
		tokennotelog.Error("loanedReject ","deocde hold err" ,err)
		return nil,err
	}
	hold.Status = pty.TokennoteStatusRejectLoan

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	//value1 := types.Encode(&hold)
	//拒绝的话  删除holdkey
	kv = append(kv,&types.KeyValue{Key:key_new,Value:nil})

	//激活冻结资金
	receiptIssuerActive,err := tokennoteAccount.ExecActive(loan.Owner,action.execaddr,hold.Amount)
	if err != nil {
		tokennotelog.Error("loanedReject","激活冻结资金异常",err)
		return nil,err
	}
	//合约归坏借款人白条
	receiptIssuer,err := tokennoteAccount.TransferWithdraw(loan.Owner,action.execaddr,hold.Amount)
	if err != nil {
		tokennotelog.Error("loanedReject","issert transfer note back err ",err)
		return nil,err
	}
	c:= false
	for k,v := range tokennotedb.token.Holds {
		if v.Addr == action.fromaddr  && v.Status != pty.TokennoteStatusRejectLoan {
			tokennotedb.token.Holds = append(tokennotedb.token.Holds[:k],tokennotedb.token.Holds[k+1:]...)
			c = true
			break
		}
	}
	if !c {
		tokennotelog.Error("loanedreject","tokennotedb dont include loan",c)
		return nil,pty.ErrTokennoteExcept
	}
	tokennotedb.token.Balance += hold.Amount
	key := calcTokennoteAddrNewKeyS(tokennotedb.token.Currency, tokennotedb.token.Issuer)
	tokennotedb.save(action.db,key)

	kv = append(kv,tokennotedb.getKVSet(key)...)
	logs = append(logs, receiptIssuerActive.Logs...)
	logs = append(logs, receiptIssuer.Logs...)
	kv = append(kv, receiptIssuerActive.KV...)
	kv = append(kv, receiptIssuer.KV...)
	logs = append(logs, tokennotedb.getLogs(pty.TyLogTokennoteLoanedReject,pty.TokennoteStatusCreated)...)
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	for _,v := range kv {
		tokennotelog.Error("kv:","key ",string(v.Key))
	}
	return receipt, nil
}


func (action *tokennoteAction) cashed(loan *pty.TokennoteCashed) (*types.Receipt, error) {
	tokennotelog.Error("loanCashed")
	if loan == nil {
		return nil, types.ErrInvalidParam
	}
	tokennote, err := getTokennoteFromDB(action.db, loan.GetSymbol(), action.fromaddr)
	if err != nil {
		tokennotelog.Error("cashed","getTokennoteFromDB err ",err)
		return nil,err
	}
	tokendb := &tokennoteDB{*tokennote}
	out := tokendb.token.Total-tokendb.token.Balance
	if out < 0 {
		return nil,pty.ErrTokennoteExcept
	}
	var out1,repay int64
	for _, v := range tokendb.token.Holds {
		if v.Status == pty.TokennoteStatusAgreeLoan {
			out1 += v.Amount
			repay += v.Repayamount
		}
	}
	if out1 != out {
		tokennotelog.Error("cashed","out != out1 err",out,out1)
		return nil,pty.ErrTokennoteExcept
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	tokenCNYYAccount ,err := account.NewAccountDB(pty.TokennoteX,pty.TokenCCNY,action.db)
	if err != nil {
		return nil,err
	}
	var reply pty.ReceiptTokennoteCashed
	tokennotelog.Error("loancashed","token db ",tokendb)
	for _,v := range tokendb.token.Holds {
		//借款人还款 出借人获取CNYY
		if v.Status == pty.TokennoteStatusAgreeLoan {

			k1 := calcTokennoteHoldKeyNew(v.Currency,v.Addr)
			kvalue ,err := action.db.Get(k1)
			if err != nil {
				tokennotelog.Error("cashed","get hold key err ",err)
				return nil,pty.ErrTokennoteExcept
			}
			var h pty.TokennoteHold
			err = types.Decode(kvalue,&h)
			if err != nil {
				tokennotelog.Error("cashed","decode hold key err ",err)
				return nil,err
			}
			if h.Status != pty.TokennoteStatusAgreeLoan {
				return nil,pty.ErrTokennoteWrongStateKey
			}
			h.Status = pty.TokennoteStatusCashed

			kv = append(kv,&types.KeyValue{Key:k1,Value:types.Encode(&h)})

			receiptTranToLend,err := tokenCNYYAccount.Transfer(action.fromaddr,v.Addr,v.Repayamount)
			if err != nil {
				tokennotelog.Error("cashed","借款人 转账错误 ",err)
				return nil,err
			}
			v.Status = pty.TokennoteStatusCashed
			logs = append(logs, receiptTranToLend.Logs...)
			kv = append(kv, receiptTranToLend.KV...)
			reply.Cashlist = append(reply.Cashlist,&pty.TokennoteCashDetail{
				Amount:v.Repayamount,
				Addr:v.Addr,
				Currency:v.Currency,
				Height:action.height,
				Time:action.blocktime,
			})
		}

	}
	tokennotelog.Error("loanCashed","replay ",reply)
	key := calcTokennoteAddrNewKeyS(loan.Symbol,action.fromaddr)
	tokendb.token.Status = pty.TokennoteStatusCashed
	tokendb.save(action.db, key)
	kv = append(kv, tokendb.getKVSet(key)...)
	//logs = append(logs, tokendb.getLogs(pty.TyLogTokennoteCashed,pty.TokennoteStatusCreated)...)
	logs = append(logs,&types.ReceiptLog{Ty:pty.TyLogTokennoteCashed,Log:types.Encode(&reply)})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	for _,v := range kv {
		tokennotelog.Error("kv:","key ",string(v.Key))
	}
	return receipt, nil
}

func (action *tokennoteAction) mint(loan *pty.TokennoteMint) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	logs = append(logs,&types.ReceiptLog{Ty:pty.TyLogTokennoteMint,Log:nil})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (action *tokennoteAction) burn(loan *pty.TokennoteBurn) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	logs = append(logs,&types.ReceiptLog{Ty:pty.TyLogTokennoteBurn,Log:nil})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func getSubDays(starttime,endtime int64) int64 {
	s := time.Unix(starttime,0)
	e := time.Unix(endtime,0)
	d := int64(math.Ceil(e.Sub(s).Hours()/24))
	return d
}

func checkTokenExist(token string, db dbm.KV) bool {
	_, err := db.Get(calcTokennoteKey(token))
	return err == nil
}


func getManageKey(key string, db dbm.KV) ([]byte, error) {
	manageKey := types.ManageKey(key)
	value, err := db.Get([]byte(manageKey))
	if err != nil {
		tokennotelog.Info("tokendb", "get db key", "not found manageKey", "key", manageKey)
		return getConfigKey(key, db)
	}
	return value, nil
}

func getConfigKey(key string, db dbm.KV) ([]byte, error) {
	configKey := types.ConfigKey(key)
	value, err := db.Get([]byte(configKey))
	if err != nil {
		tokennotelog.Info("tokendb", "get db key", "not found configKey", "key", configKey)
		return nil, err
	}
	return value, nil
}

func validOperator(addr, key string, db dbm.KV) (bool, error) {
	value, err := getManageKey(key, db)
	if err != nil {
		tokennotelog.Info("tokendb", "get db key", "not found", "key", key)
		return false, err
	}
	if value == nil {
		tokennotelog.Info("tokendb", "get db key", "  found nil value", "key", key)
		return false, nil
	}

	var item types.ConfigItem
	err = types.Decode(value, &item)
	if err != nil {
		tokennotelog.Error("tokendb", "get db key", err)
		return false, err // types.ErrBadConfigValue
	}

	for _, op := range item.GetArr().Value {
		if op == addr {
			return true, nil
		}
	}

	return false, nil
}

func calcTokennoteAssetsKey(addr string) []byte {
	return []byte(fmt.Sprintf(tokennoteAssetsPrefix+"%s", addr))
}

func calcTokennoteAssetsTimeKey(addr string) []byte {
	return []byte(fmt.Sprintf(tokennoteAssetsPrefix+"-time-"+"%s", addr))
}

func getTokennoteAssetsKey(addr string, db dbm.KVDB) (*types.ReplyStrings, error) {

	key := calcTokennoteAssetsKey(addr)
	value, err := db.Get(key)
	if err != nil && err != types.ErrNotFound {
		tokennotelog.Error("tokendb", "GetTokenAssetsKey", err)
		return nil, err
	}
	var assets types.ReplyStrings
	if err == types.ErrNotFound {
		return &assets, nil
	}
	err = types.Decode(value, &assets)
	if err != nil {
		tokennotelog.Error("tokendb", "GetTokenAssetsKey", err)
		return nil, err
	}
	return &assets, nil
}

func getTokennoteAssetsTimeKey(addr string, db dbm.KVDB) (*pty.ReplyAccountTokennoteList, error) {
	key := calcTokennoteAssetsTimeKey(addr)
	value, err := db.Get(key)
	if err != nil && err != types.ErrNotFound {
		tokennotelog.Error("tokendb", "GetTokenAssetstimeKey", err)
		return nil, err
	}
	var assets pty.ReplyAccountTokennoteList
	if err == types.ErrNotFound {
		return &assets, nil
	}
	err = types.Decode(value, &assets)
	if err != nil {
		tokennotelog.Error("tokendb", "GetTokenAssetstimeKey", err)
		return nil, err
	}
	return &assets, nil
}

// AddTokenToAssets 添加个人资产列表
func AddTokennoteToAssets(addr string, db dbm.KVDB, symbol string,time int64) []*types.KeyValue {
	tokenAssets, err := getTokennoteAssetsKey(addr, db)
	if err != nil {
		return nil
	}
	if tokenAssets == nil {
		tokenAssets = &types.ReplyStrings{}
	}

	var found = false
	for _, sym := range tokenAssets.Datas {
		if sym == symbol {
			found = true
			break
		}
	}
	if !found {
		tokenAssets.Datas = append(tokenAssets.Datas, symbol)
	}
	var kv []*types.KeyValue
	kv = append(kv, &types.KeyValue{Key: calcTokennoteAssetsKey(addr), Value: types.Encode(tokenAssets)})
	//fmt.Println("1:",kv)
	//if time > 0 {
	//	tokenAssets1, err := getTokennoteAssetsTimeKey(addr, db)
	//	if err != nil {
	//		return nil
	//	}
	//	if tokenAssets1 == nil {
	//		tokenAssets1 = &pty.ReplyAccountTokennoteList{}
	//	}
	//
	//	var found = false
	//	for _, sym := range tokenAssets1.List {
	//		if sym.Currency == symbol && sym.Time == time {
	//			found = true
	//			break
	//		}
	//	}
	//	if !found {
	//		tokenAssets1.List = append(tokenAssets1.List, &pty.TokennoteAddrTime{Currency:symbol,Time:time})
	//	}
	//	kv = append(kv, &types.KeyValue{Key: calcTokennoteAssetsTimeKey(addr), Value: types.Encode(tokenAssets1)})
	//}
	//tokennotelog.Error("AddTokennoteToAssets:","localassets:",kv)
	//fmt.Println("2:",kv)
	return kv
}

func DeleteTokennoteInAssets(addr string, db dbm.KVDB, symbol string,time int64) []*types.KeyValue {
	var kv []*types.KeyValue
	//if time > 0 {
		tokenAssets1, err := getTokennoteAssetsTimeKey(addr, db)
		if err != nil {
			return nil
		}
		if tokenAssets1 == nil {
			tokenAssets1 = &pty.ReplyAccountTokennoteList{}
		}


		for k, sym := range tokenAssets1.List {
			if sym.Currency == symbol  {
				temp := tokenAssets1.List
				tokenAssets1.List = temp[:k]
				tokenAssets1.List = append(tokenAssets1.List,temp[k+1:]...)
				break
			}
		}

		kv = append(kv, &types.KeyValue{Key: calcTokennoteAssetsTimeKey(addr), Value: types.Encode(tokenAssets1)})
	//}
	tokennotelog.Error("AddTokennoteToAssets:","localassets:",kv)
	return kv
}

func inBlacklist(symbol, key string, db dbm.KV) (bool, error) {
	found, err := validOperator(symbol, key, db)
	return found, err
}


func validFinisher(addr string, db dbm.KV) (bool, error) {
	return validOperator(addr, finisherKey, db)
}

func isUpperChar(a byte) bool {
	res := (a <= 'Z' && a >= 'A')
	return res
}

func isNumberStr(a byte) bool {
	res := (a <= '9' && a >= '0')
	return res
}

//第一个必须是字母
func validSymbol(cs []byte) bool {
	for k, c := range cs {
		if k == 0 {
			if !isUpperChar(c) {
				return false
			}
		} else {
			if isUpperChar(c) || isNumberStr(c) {
				return true
			}
			return false
		}

	}
	return true
}

func validSymbolWithHeight(cs []byte, height int64) bool {
	if !types.IsDappFork(height, pty.TokennoteX, "ForkBadTokenSymbol") {
		symbol := string(cs)
		upSymbol := strings.ToUpper(symbol)
		return upSymbol == symbol
	}
	return validSymbol(cs)
}

func getCalcAmount(amount,rate,day int64,ty int64) (int64,error) {
	amountbig := big.NewInt(amount)
	ratebig := big.NewInt(rate)
	daybig := big.NewInt(day)
	s := big.NewInt(1)
	if ty == 0 {//年利率

		day360 := big.NewInt(int64(360))
		unit := big.NewInt(int64(10000))
		sm := big.NewInt(1)
		repayson := s.Mul(amountbig,ratebig).Mul(s,daybig).String()
		repaymoth := sm.Mul(sm,day360).Mul(sm,unit).String()
		repaysonf ,err := strconv.ParseFloat(repayson,64)
		if err != nil {
			return 0,err
		}
		repayMathf,err := strconv.ParseFloat(repaymoth,64)
		if err != nil {
			return 0,err
		}
		repayfloat64 := repaysonf/repayMathf

		repaycalc := fmt.Sprintf("%0.f",repayfloat64)

		repayint64, err := strconv.ParseInt(repaycalc, 10, 64)
		if err != nil {
			return 0,err
		}
		return repayint64+amount,nil
	} else {//日利率
		unit := big.NewInt(int64(100000))
		repayson := s.Mul(amountbig,ratebig).Mul(s,daybig).String()
		repaymoth := unit.String()
		repaysonf ,err := strconv.ParseFloat(repayson,64)
		if err != nil {
			return 0,err
		}
		repayMathf,err := strconv.ParseFloat(repaymoth,64)
		if err != nil {
			return 0,err
		}
		repayfloat64 := repaysonf/repayMathf
		repaycalc := fmt.Sprintf("%0.f",repayfloat64)

		repayint64, err := strconv.ParseInt(repaycalc, 10, 64)
		if err != nil {
			return 0,err
		}
		return repayint64+amount,nil
	}

}