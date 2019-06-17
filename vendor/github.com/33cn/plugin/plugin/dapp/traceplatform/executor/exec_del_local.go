package executor

import (
	gty "github.com/33cn/plugin/plugin/dapp/traceplatform/types"
	"github.com/33cn/chain33/types"
	"fmt"
)

func (g *Traceplatform) ExecDelLocal_TraceplatformAddGood(payload *gty.TraceplatformAddGood,tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}

	txindex := fmt.Sprintf("%018d", g.GetHeight()*types.MaxTxsPerBlock+int64(index))
	nfcKey := traceplatformKeyNFCCode(payload.GetGoodinfo().NfcCode,txindex)

	set.KV = append(set.KV,&types.KeyValue{Key:nfcKey,Value:nil})
	kv,err := g.updateNFCCodeHistoryCount(payload.Goodinfo.NfcCode,1,false)
	if err == nil && kv != nil {
		set.KV = append(set.KV,kv)
	}
	fmt.Println("ExecDelLocal_TraceplatformAddGood",set.KV)
	return set,nil
}

func (g *Traceplatform) ExecDelLocal_TraceplatformAddGoods(payload *gty.TraceplatformAddGoods,tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}

	txindex := fmt.Sprintf("%018d", g.GetHeight()*types.MaxTxsPerBlock+int64(index))
	for _,v := range payload.GetNfcCode() {
		nfcKey := traceplatformKeyNFCCode(v,txindex)
		set.KV = append(set.KV,&types.KeyValue{Key:nfcKey,Value:nil})
		kv,err := g.updateNFCCodeHistoryCount(v,1,false)
		if err == nil && kv != nil {
			set.KV = append(set.KV,kv)
		}
	}
	fmt.Println("ExecDelLocal_TraceplatformAddGoods",set.KV)
	return set ,nil
}