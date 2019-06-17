// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc

import (
	"encoding/hex"

	"github.com/33cn/chain33/account"
	"github.com/33cn/chain33/common/address"
	rpctypes "github.com/33cn/chain33/rpc/types"
	"github.com/33cn/chain33/types"
	tokennotety "github.com/33cn/plugin/plugin/dapp/tokennote/types"
	"golang.org/x/net/context"
)

//TODO:和GetBalance进行泛化处理，同时LoadAccounts和LoadExecAccountQueue也需要进行泛化处理, added by hzj
func (c *channelClient) getTokennoteBalance(in *tokennotety.ReqTokennoteBalance) ([]*types.Account, error) {
	accountTokennotedb, err := account.NewAccountDB(tokennotety.TokennoteX, in.GetTokennoteSymbol(), nil)
	if err != nil {
		return nil, err
	}
	switch in.GetExecer() {
	case types.ExecName(tokennotety.TokennoteX):
		addrs := in.GetAddresses()
		var queryAddrs []string
		for _, addr := range addrs {
			if err := address.CheckAddress(addr); err != nil {
				addr = string(accountTokennotedb.AccountKey(addr))
			}
			queryAddrs = append(queryAddrs, addr)
		}

		accounts, err := accountTokennotedb.LoadAccounts(c.QueueProtocolAPI, queryAddrs)
		if err != nil {
			log.Error("GetTokennoteBalance", "err", err.Error(), "tokennote symbol", in.GetTokennoteSymbol(), "address", queryAddrs)
			return nil, err
		}
		return accounts, nil

	default: //trade
		execaddress := address.ExecAddress(in.GetExecer())
		addrs := in.GetAddresses()
		var accounts []*types.Account
		for _, addr := range addrs {
			acc, err := accountTokennotedb.LoadExecAccountQueue(c.QueueProtocolAPI, addr, execaddress)
			if err != nil {
				log.Error("GetTokennoteBalance for exector", "err", err.Error(), "tokennote symbol", in.GetTokennoteSymbol(),
					"address", addr)
				continue
			}
			accounts = append(accounts, acc)
		}

		return accounts, nil
	}
}

// GetTokennoteBalance 获取tokennote金额（channelClient）
func (c *channelClient) GetTokennoteBalance(ctx context.Context, in *tokennotety.ReqTokennoteBalance) (*types.Accounts, error) {
	reply, err := c.getTokennoteBalance(in)
	if err != nil {
		return nil, err
	}
	return &types.Accounts{Acc: reply}, nil
}

// GetTokennoteBalance 获取tokennote金额 (Jrpc)
func (c *Jrpc) GetTokennoteBalance(in tokennotety.ReqTokennoteBalance, result *interface{}) error {
	balances, err := c.cli.GetTokennoteBalance(context.Background(), &in)
	if err != nil {
		return err
	}
	var accounts []*rpctypes.Account
	for _, balance := range balances.Acc {
		accounts = append(accounts, &rpctypes.Account{Addr: balance.GetAddr(),
			Balance:  balance.GetBalance(),
			Currency: balance.GetCurrency(),
			Frozen:   balance.GetFrozen()})
	}
	*result = accounts
	return nil
}

// CreateRawTokennotePreCreateTx 创建未签名的创建Tokennote交易
func (c *Jrpc) CreateRawTokennotePreCreateTx(param *tokennotety.TokennoteCreate, result *interface{}) error {
	if param == nil || param.Currency == "" {
		return types.ErrInvalidParam
	}
	data, err := types.CallCreateTx(types.ExecName(tokennotety.TokennoteX), "TokennoteCreate", param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(data)
	return nil
}

