package executor

import (
	gty "github.com/33cn/plugin/plugin/dapp/hacker/types"
	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/types"
	"fmt"
)

func (g *Hacker) ExecLocal_HackerAddGood(payload *gty.HackerAddBill, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {

	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}

	//txindex := fmt.Sprintf("%018d", g.GetHeight()*types.MaxTxsPerBlock+int64(index))
	//nfcKey := hackerKeyNFCCode(payload.GetGoodinfo().NfcCode,txindex)
	//nfcIndex := &gty.HackerNfcCodeIndex{Nfc:payload.Goodinfo.NfcCode,Addr:address.PubKeyToAddress(tx.Signature.Pubkey).String(),Height:g.GetHeight(),Ty:gty.HackerAddGoodAction}
	//txByte := types.Encode(nfcIndex)
	//set.KV = append(set.KV,&types.KeyValue{Key:nfcKey,Value:txByte})
	//kv,err := g.updateNFCCodeHistoryCount(payload.Goodinfo.NfcCode,1,true)
	//if err == nil && kv != nil {
	//	set.KV = append(set.KV,kv)
	//}
	fmt.Println("ExecLocal_AddBill",set.KV)
	return set,nil
}


func (g *Hacker) updateNFCCodeHistoryCount(nfc string,amount int64,isadd bool) (*types.KeyValue,error) {
	txcount,err := g.getNfcTxsCount(g.GetLocalDB(),nfc)
	if err != nil {
		return nil,err
	}
	if isadd {
		txcount += amount
	} else {
		txcount -= amount
	}
	var c types.Int64
	c.Data = txcount
	//set cache
	g.GetLocalDB().Set(hackerKeyNFCCodeHistory(nfc),types.Encode(&c))
	return &types.KeyValue{Key:hackerKeyNFCCodeHistory(nfc),Value:types.Encode(&c)},nil
}

func (g *Hacker) getNfcTxsCount(db dbm.KVDB, addr string) (int64, error) {
	count := types.Int64{}
	TxsCount, err := db.Get(hackerKeyNFCCodeHistory(addr))
	if err != nil && err != types.ErrNotFound {
		return 0, err
	}
	if len(TxsCount) == 0 {
		return 0, nil
	}
	err = types.Decode(TxsCount, &count)
	if err != nil {
		return 0, err
	}
	return count.Data, nil
}