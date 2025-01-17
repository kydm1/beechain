// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"fmt"

	"strings"

	"github.com/33cn/chain33/account"
	"github.com/33cn/chain33/common/address"
	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	pty "github.com/33cn/plugin/plugin/dapp/token/types"
)

type tokenDB struct {
	token pty.Token
}

func newTokenDB(preCreate *pty.TokenPreCreate, creator string, height int64) *tokenDB {
	t := &tokenDB{}
	t.token.Name = preCreate.GetName()
	t.token.Symbol = preCreate.GetSymbol()
	t.token.Introduction = preCreate.GetIntroduction()
	t.token.Total = preCreate.GetTotal()
	t.token.Price = preCreate.GetPrice()
	//token可以由自己进行创建，也可以通过委托给其他地址进行创建
	t.token.Owner = preCreate.GetOwner()
	t.token.Creator = creator
	t.token.Status = pty.TokenStatusPreCreated
	if types.IsDappFork(height, pty.TokenX, pty.ForkTokenSymbolWithNumberX) {
		t.token.Category = preCreate.Category
	}
	return t
}

func (t *tokenDB) save(db dbm.KV, key []byte) {
	set := t.getKVSet(key)
	for i := 0; i < len(set); i++ {
		err := db.Set(set[i].GetKey(), set[i].Value)
		if err != nil {
			panic(err)
		}
	}
}

func (t *tokenDB) getLogs(ty int32, status int32) []*types.ReceiptLog {
	var log []*types.ReceiptLog
	value := types.Encode(&pty.ReceiptToken{Symbol: t.token.Symbol, Owner: t.token.Owner, Status: t.token.Status})
	log = append(log, &types.ReceiptLog{Ty: ty, Log: value})

	return log
}

//key:mavl-create-token-addr-xxx or mavl-token-xxx <-----> value:token
func (t *tokenDB) getKVSet(key []byte) (kvset []*types.KeyValue) {
	value := types.Encode(&t.token)
	kvset = append(kvset, &types.KeyValue{Key: key, Value: value})
	return kvset
}

func loadTokenDB(db dbm.KV, symbol string) (*tokenDB, error) {
	token, err := db.Get(calcTokenKey(symbol))
	if err != nil {
		tokenlog.Error("tokendb load ", "Can't get token form db for token", symbol)
		return nil, pty.ErrTokenNotExist
	}
	var t pty.Token
	err = types.Decode(token, &t)
	if err != nil {
		tokenlog.Error("tokendb load", "Can't decode token info", symbol)
		return nil, err
	}
	return &tokenDB{t}, nil
}

func (t *tokenDB) mint(db dbm.KV, addr string, amount int64) ([]*types.KeyValue, []*types.ReceiptLog, error) {
	if t.token.Owner != addr {
		return nil, nil, types.ErrNotAllow
	}
	if t.token.Total+amount > types.MaxTokenBalance {
		return nil, nil, types.ErrAmount
	}
	prevToken := t.token
	t.token.Total += amount

	kvs := append(t.getKVSet(calcTokenKey(t.token.Symbol)), t.getKVSet(calcTokenAddrNewKeyS(t.token.Symbol, t.token.Owner))...)
	logs := []*types.ReceiptLog{{Ty: pty.TyLogTokenMint, Log: types.Encode(&pty.ReceiptTokenAmount{Prev: &prevToken, Current: &t.token})}}
	return kvs, logs, nil
}

func (t *tokenDB) burn(db dbm.KV, amount int64) ([]*types.KeyValue, []*types.ReceiptLog, error) {
	if t.token.Total < amount {
		return nil, nil, types.ErrNoBalance
	}
	prevToken := t.token
	t.token.Total -= amount

	kvs := append(t.getKVSet(calcTokenKey(t.token.Symbol)), t.getKVSet(calcTokenAddrNewKeyS(t.token.Symbol, t.token.Owner))...)
	logs := []*types.ReceiptLog{{Ty: pty.TyLogTokenBurn, Log: types.Encode(&pty.ReceiptTokenAmount{Prev: &prevToken, Current: &t.token})}}
	return kvs, logs, nil
}

func getTokenFromDB(db dbm.KV, symbol string, owner string) (*pty.Token, error) {
	key := calcTokenAddrKeyS(symbol, owner)
	value, err := db.Get(key)
	if err != nil {
		// not found old key
		key = calcTokenAddrNewKeyS(symbol, owner)
		value, err = db.Get(key)
		if err != nil {
			return nil, err
		}
	}

	var token pty.Token
	if err = types.Decode(value, &token); err != nil {
		tokenlog.Error("getTokenFromDB", "Fail to decode types.token for key", string(key), "err info is", err)
		return nil, err
	}
	return &token, nil
}

type tokenAction struct {
	coinsAccount *account.DB
	db           dbm.KV
	txhash       []byte
	fromaddr     string
	toaddr       string
	blocktime    int64
	height       int64
	execaddr     string
}

func newTokenAction(t *token, toaddr string, tx *types.Transaction) *tokenAction {
	hash := tx.Hash()
	fromaddr := tx.From()
	return &tokenAction{t.GetCoinsAccount(), t.GetStateDB(), hash, fromaddr, toaddr,
		t.GetBlockTime(), t.GetHeight(), dapp.ExecAddress(string(tx.Execer))}
}

func (action *tokenAction) preCreate(token *pty.TokenPreCreate) (*types.Receipt, error) {
	tokenlog.Debug("preCreate")
	if token == nil {
		return nil, types.ErrInvalidParam
	}
	if len(token.GetName()) > pty.TokenNameLenLimit {
		return nil, pty.ErrTokenNameLen
	} else if len(token.GetIntroduction()) > pty.TokenIntroLenLimit {
		return nil, pty.ErrTokenIntroLen
	} else if len(token.GetSymbol()) > pty.TokenSymbolLenLimit {
		return nil, pty.ErrTokenSymbolLen
	} else if token.GetTotal() > types.MaxTokenBalance || token.GetTotal() <= 0 {
		return nil, pty.ErrTokenTotalOverflow
	}
	if types.IsDappFork(action.height, pty.TokenX, pty.ForkTokenCheckX) {
		if err := address.CheckAddress(token.Owner); err != nil {
			return nil, err
		}
	}
	if !types.IsDappFork(action.height, pty.TokenX, pty.ForkTokenSymbolWithNumberX) {
		if token.Category != 0 {
			return nil, types.ErrNotSupport
		}
	}

	if !validSymbolWithHeight([]byte(token.GetSymbol()), action.height) {
		tokenlog.Error("token precreate ", "symbol need be upper", token.GetSymbol())
		return nil, pty.ErrTokenSymbolUpper
	}

	if checkTokenExist(token.GetSymbol(), action.db) {
		return nil, pty.ErrTokenExist
	}

	if checkTokenHasPrecreateWithHeight(token.GetSymbol(), token.GetOwner(), action.db, action.height) {
		return nil, pty.ErrTokenHavePrecreated
	}

	if types.IsDappFork(action.height, pty.TokenX, pty.ForkTokenBlackListX) {
		found, err := inBlacklist(token.GetSymbol(), blacklist, action.db)
		if err != nil {
			return nil, err
		}
		if found {
			return nil, pty.ErrTokenBlacklist
		}
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	if types.IsDappFork(action.height, pty.TokenX, pty.ForkTokenPriceX) && token.GetPrice() == 0 {
		// pay for create token offline
	} else {
		receipt, err := action.coinsAccount.ExecFrozen(action.fromaddr, action.execaddr, token.GetPrice())
		if err != nil {
			tokenlog.Error("token precreate ", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", token.GetTotal())
			return nil, err
		}
		logs = append(logs, receipt.Logs...)
		kv = append(kv, receipt.KV...)
	}

	tokendb := newTokenDB(token, action.fromaddr, action.height)
	var statuskey []byte
	var key []byte
	if types.IsFork(action.height, "ForkExecKey") {
		statuskey = calcTokenStatusNewKeyS(tokendb.token.Symbol, tokendb.token.Owner, pty.TokenStatusPreCreated)
		key = calcTokenAddrNewKeyS(tokendb.token.Symbol, tokendb.token.Owner)
	} else {
		statuskey = calcTokenStatusKey(tokendb.token.Symbol, tokendb.token.Owner, pty.TokenStatusPreCreated)
		key = calcTokenAddrKeyS(tokendb.token.Symbol, tokendb.token.Owner)
	}

	tokendb.save(action.db, statuskey)
	tokendb.save(action.db, key)

	logs = append(logs, tokendb.getLogs(pty.TyLogPreCreateToken, pty.TokenStatusPreCreated)...)
	kv = append(kv, tokendb.getKVSet(key)...)
	kv = append(kv, tokendb.getKVSet(statuskey)...)
	//tokenlog.Info("func token preCreate", "token:", tokendb.token.Symbol, "owner:", tokendb.token.Owner,
	//	"key:", key, "key string", string(key), "value:", tokendb.getKVSet(key)[0].Value)

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (action *tokenAction) finishCreate(tokenFinish *pty.TokenFinishCreate) (*types.Receipt, error) {
	tokenlog.Debug("finishCreate")
	if tokenFinish == nil {
		return nil, types.ErrInvalidParam
	}
	token, err := getTokenFromDB(action.db, tokenFinish.GetSymbol(), tokenFinish.GetOwner())
	if err != nil || token.Status != pty.TokenStatusPreCreated {
		return nil, pty.ErrTokenNotPrecreated
	}

	approverValid := false

	for _, approver := range conf.GStrList("tokenApprs") {
		if approver == action.fromaddr {
			approverValid = true
			break
		}
	}

	hasPriv, ok := validFinisher(action.fromaddr, action.db)
	if (ok != nil || !hasPriv) && !approverValid {
		return nil, pty.ErrTokenCreatedApprover
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	if types.IsDappFork(action.height, pty.TokenX, "ForkTokenPrice") && token.GetPrice() == 0 {
		// pay for create token offline
	} else {
		//将之前冻结的资金转账到fund合约账户中
		receiptForCoin, err := action.coinsAccount.ExecTransferFrozen(token.Creator, action.toaddr, action.execaddr, token.Price)
		if err != nil {
			tokenlog.Error("token finishcreate ", "addr", action.fromaddr, "execaddr", action.execaddr, "token", token.Symbol)
			return nil, err
		}
		logs = append(logs, receiptForCoin.Logs...)
		kv = append(kv, receiptForCoin.KV...)
	}

	//创建token类型的账户，同时需要创建的额度存入

	tokenAccount, err := account.NewAccountDB("token", tokenFinish.GetSymbol(), action.db)
	if err != nil {
		return nil, err
	}
	tokenlog.Debug("finishCreate", "token.Owner", token.Owner, "token.GetTotal()", token.GetTotal())
	receiptForToken, err := tokenAccount.GenesisInit(token.Owner, token.GetTotal())
	if err != nil {
		return nil, err
	}
	//更新token的状态为已经创建
	token.Status = pty.TokenStatusCreated
	tokendb := &tokenDB{*token}
	var key []byte
	if types.IsFork(action.height, "ForkExecKey") {
		key = calcTokenAddrNewKeyS(tokendb.token.Symbol, tokendb.token.Owner)
	} else {
		key = calcTokenAddrKeyS(tokendb.token.Symbol, tokendb.token.Owner)
	}
	tokendb.save(action.db, key)

	logs = append(logs, receiptForToken.Logs...)
	logs = append(logs, tokendb.getLogs(pty.TyLogFinishCreateToken, pty.TokenStatusCreated)...)
	kv = append(kv, receiptForToken.KV...)
	kv = append(kv, tokendb.getKVSet(key)...)

	key = calcTokenKey(tokendb.token.Symbol)
	//因为该token已经被创建，需要保存一个全局的token，防止其他用户再次创建
	tokendb.save(action.db, key)
	kv = append(kv, tokendb.getKVSet(key)...)
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (action *tokenAction) revokeCreate(tokenRevoke *pty.TokenRevokeCreate) (*types.Receipt, error) {
	if tokenRevoke == nil {
		return nil, types.ErrInvalidParam
	}
	token, err := getTokenFromDB(action.db, tokenRevoke.GetSymbol(), tokenRevoke.GetOwner())
	if err != nil {
		tokenlog.Error("token revokeCreate ", "Can't get token form db for token", tokenRevoke.GetSymbol())
		return nil, pty.ErrTokenNotPrecreated
	}

	if token.Status != pty.TokenStatusPreCreated {
		tokenlog.Error("token revokeCreate ", "token's status should be precreated to be revoked for token", tokenRevoke.GetSymbol())
		return nil, pty.ErrTokenCanotRevoked
	}

	//确认交易发起者的身份，token的发起人可以撤销该项token的创建
	//token的owner允许撤销交易
	if action.fromaddr != token.Owner && action.fromaddr != token.Creator {
		tokenlog.Error("tprocTokenRevokeCreate, different creator/owner vs actor of this revoke",
			"action.fromaddr", action.fromaddr, "creator", token.Creator, "owner", token.Owner)
		return nil, pty.ErrTokenRevoker
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	if types.IsDappFork(action.height, pty.TokenX, pty.ForkTokenPriceX) && token.GetPrice() == 0 {
		// pay for create token offline
	} else {
		//解锁之前冻结的资金
		receipt, err := action.coinsAccount.ExecActive(token.Creator, action.execaddr, token.Price)
		if err != nil {
			tokenlog.Error("token revokeCreate error ", "error info", err, "creator addr", token.Creator, "execaddr", action.execaddr, "token", token.Symbol)
			return nil, err
		}
		logs = append(logs, receipt.Logs...)
		kv = append(kv, receipt.KV...)
	}

	token.Status = pty.TokenStatusCreateRevoked
	tokendb := &tokenDB{*token}
	var key []byte
	if types.IsFork(action.height, "ForkExecKey") {
		key = calcTokenAddrNewKeyS(tokendb.token.Symbol, tokendb.token.Owner)
	} else {
		key = calcTokenAddrKeyS(tokendb.token.Symbol, tokendb.token.Owner)
	}
	tokendb.save(action.db, key)

	logs = append(logs, tokendb.getLogs(pty.TyLogRevokeCreateToken, pty.TokenStatusCreateRevoked)...)
	kv = append(kv, tokendb.getKVSet(key)...)

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func checkTokenExist(token string, db dbm.KV) bool {
	_, err := db.Get(calcTokenKey(token))
	return err == nil
}

// bug: prepare again after revoke, need to check status, fixed in fork ForkTokenCheckX
func checkTokenHasPrecreate(token, owner string, status int32, db dbm.KV) bool {
	_, err := db.Get(calcTokenAddrKeyS(token, owner))
	if err == nil {
		return true
	}
	_, err = db.Get(calcTokenAddrNewKeyS(token, owner))
	return err == nil
}

func checkTokenHasPrecreateWithHeight(token, owner string, db dbm.KV, height int64) bool {
	if !types.IsDappFork(height, pty.TokenX, pty.ForkTokenCheckX) {
		return checkTokenHasPrecreate(token, owner, pty.TokenStatusPreCreated, db)
	}

	tokenStatus, err := db.Get(calcTokenAddrNewKeyS(token, owner))
	if err != nil {
		tokenStatus, err = db.Get(calcTokenAddrKeyS(token, owner))
		if err != nil {
			return false
		}
	}

	var t pty.Token
	if err = types.Decode(tokenStatus, &t); err != nil {
		tokenlog.Error("checkTokenHasPrecreateWithHeight", "Fail to decode types.token for key err info is", err)
		panic("data err: checkTokenHasPrecreateWithHeight Fail to decode types.token for key err info is" + err.Error())
	}
	return t.Status == pty.TokenStatusPreCreated
}

func validFinisher(addr string, db dbm.KV) (bool, error) {
	return validOperator(addr, finisherKey, db)
}

func getManageKey(key string, db dbm.KV) ([]byte, error) {
	manageKey := types.ManageKey(key)
	value, err := db.Get([]byte(manageKey))
	if err != nil {
		tokenlog.Info("tokendb", "get db key", "not found manageKey", "key", manageKey)
		return getConfigKey(key, db)
	}
	return value, nil
}

func getConfigKey(key string, db dbm.KV) ([]byte, error) {
	configKey := types.ConfigKey(key)
	value, err := db.Get([]byte(configKey))
	if err != nil {
		tokenlog.Info("tokendb", "get db key", "not found configKey", "key", configKey)
		return nil, err
	}
	return value, nil
}

func validOperator(addr, key string, db dbm.KV) (bool, error) {
	value, err := getManageKey(key, db)
	if err != nil {
		tokenlog.Info("tokendb", "get db key", "not found", "key", key)
		return false, err
	}
	if value == nil {
		tokenlog.Info("tokendb", "get db key", "  found nil value", "key", key)
		return false, nil
	}

	var item types.ConfigItem
	err = types.Decode(value, &item)
	if err != nil {
		tokenlog.Error("tokendb", "get db key", err)
		return false, err // types.ErrBadConfigValue
	}

	for _, op := range item.GetArr().Value {
		if op == addr {
			return true, nil
		}
	}

	return false, nil
}

func calcTokenAssetsKey(addr string) []byte {
	return []byte(fmt.Sprintf(tokenAssetsPrefix+"%s", addr))
}

func getTokenAssetsKey(addr string, db dbm.KVDB) (*types.ReplyStrings, error) {
	key := calcTokenAssetsKey(addr)
	value, err := db.Get(key)
	if err != nil && err != types.ErrNotFound {
		tokenlog.Error("tokendb", "GetTokenAssetsKey", err)
		return nil, err
	}
	var assets types.ReplyStrings
	if err == types.ErrNotFound {
		return &assets, nil
	}
	err = types.Decode(value, &assets)
	if err != nil {
		tokenlog.Error("tokendb", "GetTokenAssetsKey", err)
		return nil, err
	}
	return &assets, nil
}

// AddTokenToAssets 添加个人资产列表
func AddTokenToAssets(addr string, db dbm.KVDB, symbol string) []*types.KeyValue {
	tokenAssets, err := getTokenAssetsKey(addr, db)
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
	kv = append(kv, &types.KeyValue{Key: calcTokenAssetsKey(addr), Value: types.Encode(tokenAssets)})
	return kv
}

func inBlacklist(symbol, key string, db dbm.KV) (bool, error) {
	found, err := validOperator(symbol, key, db)
	return found, err
}

func isUpperChar(a byte) bool {
	res := (a <= 'Z' && a >= 'A')
	return res
}

func isDigit(a byte) bool {
	return (a <= '9' && a >= '0')
}

func validSymbolForkBadTokenSymbol(cs []byte) bool {
	for _, c := range cs {
		if !isUpperChar(c) {
			return false
		}
	}
	return true
}

func validSymbolForkTokenSymbolWithNumber(cs []byte) bool {
	for _, c := range cs {
		if !isUpperChar(c) && !isDigit(c) {
			return false
		}
	}
	return true
}

func validSymbolOriginal(cs []byte) bool {
	symbol := string(cs)
	upSymbol := strings.ToUpper(symbol)
	return upSymbol == symbol
}

func validSymbolWithHeight(cs []byte, height int64) bool {
	if types.IsDappFork(height, pty.TokenX, pty.ForkTokenSymbolWithNumberX) {
		return validSymbolForkTokenSymbolWithNumber(cs)
	} else if types.IsDappFork(height, pty.TokenX, pty.ForkBadTokenSymbolX) {
		return validSymbolForkBadTokenSymbol(cs)
	}
	return validSymbolOriginal(cs)
}

// 铸币不可控， 也是麻烦。 2选1
// 1. 谁可以发起
// 2. 是否需要审核  这个会增加管理的成本
// 现在实现选择 1
func (action *tokenAction) mint(mint *pty.TokenMint) (*types.Receipt, error) {
	if mint == nil {
		return nil, types.ErrInvalidParam
	}
	if mint.GetAmount() < 0 || mint.GetAmount() > types.MaxTokenBalance || mint.GetSymbol() == "" {
		return nil, types.ErrInvalidParam
	}

	tokendb, err := loadTokenDB(action.db, mint.GetSymbol())
	if err != nil {
		return nil, err
	}

	if tokendb.token.Category&pty.CategoryMintBurnSupport == 0 {
		tokenlog.Error("Can't mint category", "category", tokendb.token.Category, "support", pty.CategoryMintBurnSupport)
		return nil, types.ErrNotSupport
	}

	kvs, logs, err := tokendb.mint(action.db, action.fromaddr, mint.Amount)
	if err != nil {
		tokenlog.Error("token mint ", "symbol", mint.GetSymbol(), "error", err, "from", action.fromaddr, "owner", tokendb.token.Owner)
		return nil, err
	}

	tokenAccount, err := account.NewAccountDB("token", mint.GetSymbol(), action.db)
	if err != nil {
		return nil, err
	}
	tokenlog.Debug("mint", "token.Owner", mint.Symbol, "token.GetTotal()", mint.Amount)
	receipt, err := tokenAccount.Mint(action.fromaddr, mint.Amount)
	if err != nil {
		return nil, err
	}

	logs = append(logs, receipt.Logs...)
	kvs = append(kvs, receipt.KV...)

	return &types.Receipt{Ty: types.ExecOk, KV: kvs, Logs: logs}, nil
}

func (action *tokenAction) burn(burn *pty.TokenBurn) (*types.Receipt, error) {
	if burn == nil {
		return nil, types.ErrInvalidParam
	}
	if burn.GetAmount() < 0 || burn.GetAmount() > types.MaxTokenBalance || burn.GetSymbol() == "" {
		return nil, types.ErrInvalidParam
	}

	tokendb, err := loadTokenDB(action.db, burn.GetSymbol())
	if err != nil {
		return nil, err
	}

	if tokendb.token.Category&pty.CategoryMintBurnSupport == 0 {
		tokenlog.Error("Can't burn category", "category", tokendb.token.Category, "support", pty.CategoryMintBurnSupport)
		return nil, types.ErrNotSupport
	}

	kvs, logs, err := tokendb.burn(action.db, burn.Amount)
	if err != nil {
		tokenlog.Error("token burn ", "symbol", burn.GetSymbol(), "error", err, "from", action.fromaddr, "owner", tokendb.token.Owner)
		return nil, err
	}

	tokenAccount, err := account.NewAccountDB("token", burn.GetSymbol(), action.db)
	if err != nil {
		return nil, err
	}
	tokenlog.Debug("burn", "token.Owner", burn.Symbol, "token.GetTotal()", burn.Amount)
	receipt, err := tokenAccount.Burn(action.fromaddr, burn.Amount)
	if err != nil {
		return nil, err
	}

	logs = append(logs, receipt.Logs...)
	kvs = append(kvs, receipt.KV...)

	return &types.Receipt{Ty: types.ExecOk, KV: kvs, Logs: logs}, nil
}
