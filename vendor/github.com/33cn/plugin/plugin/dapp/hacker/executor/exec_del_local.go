package executor

import (
	gty "github.com/33cn/plugin/plugin/dapp/hacker/types"
	"github.com/33cn/chain33/types"
	"fmt"
)

func (g *Hacker) ExecDelLocal_AddBill(payload *gty.HackerAddBill,tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}

	//txindex := fmt.Sprintf("%018d", g.GetHeight()*types.MaxTxsPerBlock+int64(index))
	//nfcKey := traceplatformKeyNFCCode(payload.GetGoodinfo().NfcCode,txindex)
	//
	//set.KV = append(set.KV,&types.KeyValue{Key:nfcKey,Value:nil})
	//kv,err := g.updateNFCCodeHistoryCount(payload.Goodinfo.NfcCode,1,false)
	//if err == nil && kv != nil {
	//	set.KV = append(set.KV,kv)
	//}
	fmt.Println("AddBill",set.KV)
	return set,nil
}
