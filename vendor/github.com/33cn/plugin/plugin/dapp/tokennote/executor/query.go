// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/33cn/chain33/types"
	tokenty "github.com/33cn/plugin/plugin/dapp/tokennote/types"
	"github.com/33cn/chain33/common/db/table"
	"fmt"
	"reflect"
)



// Query_GetTokens 获取token
func (t *tokennote) Query_GetTokennotes(in *tokenty.ReqTokennotes) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	return t.getTokennotes(in)
}

// Query_GetTokenInfo 获取token信息
func (t *tokennote) Query_GetTokenInfo(in *types.ReqString) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	return t.getTokennoteInfo(in.GetData())
}

// Query_GetTotalAmount 获取token总量
func (t *tokennote) Query_GetTotalAmount(in *types.ReqString) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	ret, err := t.getTokennoteInfo(in.GetData())
	if err != nil {
		return nil, err
	}
	tokenInfo, ok := ret.(*tokenty.LocalTokennote)
	if !ok {
		return nil, types.ErrTypeAsset
	}
	return &types.TotalAmount{
		Total: tokenInfo.Total,
	}, nil
}

// Query_GetAddrReceiverforTokens 获取token接受人数据
func (t *tokennote) Query_GetAddrReceiverforTokens(in *tokenty.ReqAddrTokennotes) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	return t.getAddrReceiverforTokennotes(in)
}

// Query_GetAccountTokenAssets 获取账户的token资产
func (t *tokennote) Query_GetAccountTokenAssets(in *tokenty.ReqAccountTokennoteAssets) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	return t.getAccountTokennoteAssets(in)
}

// Query_GetTxByToken 获取token相关交易
func (t *tokennote) Query_GetTxByTokennote(in *tokenty.ReqTokennoteTx) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	if !cfg.SaveTokennoteTxList {
		return nil, types.ErrActionNotSupport
	}
	return t.getTxByTokennote(in)
}

//获取白条下的借款记录   当前不支持拆分借
func (t *tokennote) Query_GetAgreeTxsByTokennote(in *types.ReplyString) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	return t.getAgreeTxByTokennote(in)
}

//获取地址下的还款记录   当前不支持拆分借
func (t *tokennote) Query_GetTxsCashedByTokennote(in *tokenty.ReqTokennoteTx) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	return t.getCashedTxByTokennote(in)
}

func (t *tokennote) Query_GetTableTxs(in *tokenty.ReqTokennoteTx) (types.Message, error) {
	if in == nil && in.Symbol == "" {
		return nil, types.ErrInvalidParam
	}

	tab,err := table.NewTable(NewTransactionRow(),t.GetLocalDB(),&opt)
	query := tab.GetQuery(t.GetLocalDB())
	rows,err := query.ListIndex("Symbol",[]byte(in.Symbol),nil,0,0)
	if err != nil {
		tokennotelog.Error("GetTableTxs","rows err",err)
		return nil,err
	}
	fmt.Println("rows:",rows)
	var txs types.Transactions
	for _, v := range rows {
		tokennotelog.Error("GetTableTxs ","row ",v)
		tokennotelog.Error("GetTableTxs","type",reflect.TypeOf(v.Data))
		txs.Txs = append(txs.Txs,v.Data.(*types.Transaction))
	}
	return &txs,nil
}

func (t *tokennote) Query_GetTokennoteAccountCredit(in *types.ReplyString) (types.Message,error) {
	if in == nil || in.Data == "" {
		return nil, types.ErrInvalidParam
	}
	count := table.NewCount(tokennoteLocalPre,"createdNotes",t.GetLocalDB())
	var credit tokenty.TokennoteAccountCredit
	num,err := count.Get()
	if err != nil {
		return nil,err
	}
	credit.Addr = in.Data
	credit.CreatedNotes = num
	return &credit,nil
}

func (t *tokennote) Query_GetTokennoteMarketList(in *tokenty.ReqTokennoteTx) (types.Message,error) {
	tab,err := table.NewTable(NewTokennoteMarketRow(),t.GetLocalDB(),&marketopt)
	if err != nil {
		return nil,err
	}
	var market tokenty.TokennoteMarket
	query := tab.GetQuery(t.GetLocalDB())
	if in != nil && in.Addr != "" {
		//对应地址下的行情
		rows,err := query.ListIndex(in.Symbol,[]byte(in.Addr),nil,in.Count,in.Direction)
		if err != nil {
			return nil,err
		}
		for _,v := range rows {
			local := v.Data.(*tokenty.Tokennote)
			value,err := t.GetStateDB().Get(calcTokennoteAddrNewKeyS(local.Currency,local.Issuer))
			if err != nil {
				continue
			}
			var new tokenty.Tokennote
			err = types.Decode(value,&new)
			if err != nil {
				tokennotelog.Error("Query_GetTokennoteMarketList","decode tokennote err ",err)
				return nil,err
			}
			market.MarketList = append(market.MarketList,&new)
		}
	} else {
		//所有市场行情
		rows ,err := query.List("primary",nil,nil,in.Count,in.Direction)
		if err != nil {
			return nil,err
		}
		for _,v := range rows {
			local := v.Data.(*tokenty.Tokennote)
			value,err := t.GetStateDB().Get(calcTokennoteAddrNewKeyS(local.Currency,local.Issuer))
			if err != nil {
				continue
			}
			var new tokenty.Tokennote
			err = types.Decode(value,&new)
			if err != nil {
				tokennotelog.Error("Query_GetTokennoteMarketList","decode tokennote err ",err)
				return nil,err
			}
			market.MarketList = append(market.MarketList,&new)
		}
	}
	return &market,nil
}

func (t *tokennote) Query_GetTokennoteContract(in *types.ReplyString) (types.Message,error) {
	if in == nil || in.Data == "" {
		return nil, types.ErrInvalidParam
	}
	contract,err := getContract(t.GetStateDB(),in.Data)
	if err != nil {
		return nil,err
	}
	return contract,nil
}