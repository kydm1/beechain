// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc_test

import (
	"strings"
	"testing"

	commonlog "github.com/33cn/chain33/common/log"
	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	"github.com/33cn/chain33/types"
	"github.com/33cn/chain33/util/testnode"
	pty "github.com/33cn/plugin/plugin/dapp/trade/types"
	"github.com/stretchr/testify/assert"

	_ "github.com/33cn/chain33/system"
	_ "github.com/33cn/plugin/plugin"
)

func init() {
	commonlog.SetLogLevel("error")
}

func TestJRPCChannel(t *testing.T) {
	// 启动RPCmocker
	mocker := testnode.New("--notset--", nil)
	defer func() {
		mocker.Close()
	}()
	mocker.Listen()

	jrpcClient := mocker.GetJSONC()

	testCases := []struct {
		fn func(*testing.T, *jsonclient.JSONClient) error
	}{
		{fn: testCreateRawTradeSellTxCmd},
		{fn: testCreateRawTradeBuyTxCmd},
		{fn: testCreateRawTradeRevokeTxCmd},
		{fn: testShowOnesSellOrdersCmd},
		{fn: testShowOnesSellOrdersStatusCmd},
		{fn: testShowTokenSellOrdersStatusCmd},
		{fn: testShowOnesBuyOrderCmd},
		{fn: testShowOnesBuyOrdersStatusCmd},
		{fn: testShowTokenBuyOrdersStatusCmd},
		{fn: testShowOnesOrdersStatusCmd},
	}
	for index, testCase := range testCases {
		err := testCase.fn(t, jrpcClient)
		if err == nil {
			continue
		}
		assert.NotEqualf(t, err, types.ErrActionNotSupport, "test index %d", index)
		if strings.Contains(err.Error(), "rpc: can't find") {
			assert.FailNowf(t, err.Error(), "test index %d", index)
		}
	}
}

func testCreateRawTradeSellTxCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	params := &pty.TradeSellTx{}
	return jrpc.Call("trade.CreateRawTradeSellTx", params, nil)
}

func testCreateRawTradeBuyTxCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	params := &pty.TradeBuyTx{}
	return jrpc.Call("trade.CreateRawTradeBuyTx", params, nil)
}

func testCreateRawTradeRevokeTxCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	params := &pty.TradeRevokeTx{}
	return jrpc.Call("trade.CreateRawTradeRevokeTx", params, nil)
}

func testShowOnesSellOrdersCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	params := rpctypes.Query4Jrpc{
		Execer:   "trade",
		FuncName: "GetOnesSellOrder",
		Payload:  types.MustPBToJSON(&pty.ReqAddrAssets{}),
	}
	var res pty.ReplySellOrders
	return jrpc.Call("Chain33.Query", params, &res)
}

func testShowOnesSellOrdersStatusCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	var rep interface{}
	var params rpctypes.Query4Jrpc
	req := &pty.ReqAddrAssets{}
	params.Execer = "trade"
	params.FuncName = "GetOnesSellOrderWithStatus"
	params.Payload = types.MustPBToJSON(req)
	rep = &pty.ReplySellOrders{}
	return jrpc.Call("Chain33.Query", params, rep)
}

func testShowTokenSellOrdersStatusCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	var rep interface{}
	var params rpctypes.Query4Jrpc
	req := &pty.ReqTokenSellOrder{}
	params.Execer = "trade"
	params.FuncName = "GetTokenSellOrderByStatus"
	params.Payload = types.MustPBToJSON(req)
	rep = &pty.ReplySellOrders{}

	return jrpc.Call("Chain33.Query", params, rep)
}

func testShowOnesBuyOrderCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	var rep interface{}
	var params rpctypes.Query4Jrpc
	req := &pty.ReqAddrAssets{}
	params.Execer = "trade"
	params.FuncName = "GetOnesBuyOrder"
	params.Payload = types.MustPBToJSON(req)
	rep = &pty.ReplyBuyOrders{}

	return jrpc.Call("Chain33.Query", params, rep)
}

func testShowOnesBuyOrdersStatusCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	var rep interface{}
	var params rpctypes.Query4Jrpc
	req := &pty.ReqAddrAssets{}
	params.Execer = "trade"
	params.FuncName = "GetOnesBuyOrderWithStatus"
	params.Payload = types.MustPBToJSON(req)
	rep = &pty.ReplyBuyOrders{}

	return jrpc.Call("Chain33.Query", params, rep)
}

func testShowTokenBuyOrdersStatusCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	var rep interface{}
	var params rpctypes.Query4Jrpc
	req := &pty.ReqTokenBuyOrder{}
	params.Execer = "trade"
	params.FuncName = "GetTokenBuyOrderByStatus"
	params.Payload = types.MustPBToJSON(req)
	rep = &pty.ReplyBuyOrders{}

	return jrpc.Call("Chain33.Query", params, rep)
}

func testShowOnesOrdersStatusCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	var rep interface{}
	var params rpctypes.Query4Jrpc
	req := &pty.ReqAddrAssets{}
	params.Execer = "trade"
	params.FuncName = "GetOnesOrderWithStatus"
	params.Payload = types.MustPBToJSON(req)
	rep = &pty.ReplyTradeOrders{}

	return jrpc.Call("Chain33.Query", params, rep)
}
