package executor

import (
	gty "github.com/33cn/plugin/plugin/dapp/traceplatform/types"
	"github.com/33cn/chain33/types"
	//"fmt"
	"errors"
	"fmt"
)

func (g *Traceplatform) Query_GetGoodById(in *types.ReplyString) (types.Message,error) {
	var goodInfo gty.TraceplatformGoodInfo
	db := g.GetStateDB()
	value, err := db.Get(traceplatformKeyGood(in.Data))
	if err != nil {
		return nil, err
	}
	types.Decode(value, &goodInfo)
	return &goodInfo, nil
}

func (g *Traceplatform) Query_GetNfcTxsCountById(in *types.ReplyString) (types.Message,error) {
	var goodInfo types.Int64
	db := g.GetLocalDB()
	value, err := db.Get(traceplatformKeyNFCCodeHistory(in.Data))
	if err != nil {
		return nil, err
	}
	types.Decode(value, &goodInfo)
	return &goodInfo, nil
}

func (g *Traceplatform) Query_GetTxByNfcCode(addr *types.ReqAddr) (types.Message,error) {
	db := g.GetLocalDB()
	var prefix []byte
	var key []byte
	var txinfos [][]byte
	var err error
	//取最新的交易hash列表

	prefix = traceplatformKeyNFCCode(addr.GetAddr(), "")

	if addr.GetHeight() == -1 {
		txinfos, err = db.List(prefix, nil, addr.Count, addr.GetDirection())
		if err != nil {
			return nil, err
		}
		if len(txinfos) == 0 {
			return nil, errors.New("tx does not exist")
		}
	} else { //翻页查找指定的txhash列表
		v := addr.GetHeight()*types.MaxTxsPerBlock + addr.GetIndex()
		heightstr := fmt.Sprintf("%018d", v)
		if addr.Flag == 0 {
			key = types.CalcTxAddrHashKey(addr.GetAddr(), heightstr)
		} else if addr.Flag > 0 { //from的交易hash列表
			key = types.CalcTxAddrDirHashKey(addr.GetAddr(), addr.Flag, heightstr)
		} else {
			return nil, errors.New("flag unknown")
		}
		txinfos, err = db.List(prefix, key, addr.Count, addr.Direction)
		if err != nil {
			return nil, err
		}
		if len(txinfos) == 0 {
			return nil, errors.New("tx does not exist")
		}
	}
	var replyTxInfos gty.TraceplatformNfcCodeIndexList
	replyTxInfos.List = make([]*gty.TraceplatformNfcCodeIndex, 0)
	for _, txinfobyte := range txinfos {
		var replyTxInfo gty.TraceplatformNfcCodeIndex
		err := types.Decode(txinfobyte, &replyTxInfo)
		if err != nil {
			return nil, err
		}
		replyTxInfos.List = append(replyTxInfos.List,&replyTxInfo)
	}
	return &replyTxInfos, nil
}
