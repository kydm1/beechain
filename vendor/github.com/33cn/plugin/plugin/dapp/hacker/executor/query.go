package executor

import (
	gty "github.com/33cn/plugin/plugin/dapp/hacker/types"
	"github.com/33cn/chain33/types"
	//"fmt"
	"errors"
	"fmt"
)

func (g *Hacker) Query_GetGoodById(in *types.ReplyString) (types.Message,error) {
	var goodInfo gty.HackerAddBill
	db := g.GetStateDB()
	value, err := db.Get(hackerKeyGood(in.Data))
	if err != nil {
		return nil, err
	}
	types.Decode(value, &goodInfo)
	return &goodInfo, nil
}

func (g *Hacker) Query_GetNfcTxsCountById(in *types.ReplyString) (types.Message,error) {
	var goodInfo types.Int64
	db := g.GetLocalDB()
	value, err := db.Get(hackerKeyNFCCodeHistory(in.Data))
	if err != nil {
		return nil, err
	}
	types.Decode(value, &goodInfo)
	return &goodInfo, nil
}

func (g *Hacker) Query_GetTxByNfcCode(addr *types.ReqAddr) (types.Message,error) {
	db := g.GetLocalDB()
	var prefix []byte
	var key []byte
	var txinfos [][]byte
	var err error
	//取最新的交易hash列表

	prefix = hackerKeyNFCCode(addr.GetAddr(), "")

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
	var replyTxInfos gty.HackerNfcCodeIndexList
	replyTxInfos.List = make([]*gty.HackerBillIndex, 0)
	for _, txinfobyte := range txinfos {
		var replyTxInfo gty.HackerBillIndex
		err := types.Decode(txinfobyte, &replyTxInfo)
		if err != nil {
			return nil, err
		}
		replyTxInfos.List = append(replyTxInfos.List,&replyTxInfo)
	}
	return &replyTxInfos, nil
}