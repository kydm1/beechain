package executor

import (
	"github.com/33cn/chain33/types"
	dbm "github.com/33cn/chain33/common/db"
	g "github.com/33cn/plugin/plugin/dapp/traceplatform/types"
)


type Action struct {
	db     dbm.KV
	txhash []byte
	height int64
	index  int
}

func NewAction(t *Traceplatform, tx *types.Transaction,index int) *Action {
	hash := tx.Hash()
	return &Action{t.GetStateDB(), hash, t.GetHeight(),index}
}

//添加商品
func (a *Action) AddGood(payload *g.TraceplatformAddGood, pubkey []byte) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	//value,err := a.db.Get(traceplatformKeyGood(payload.Goodinfo.NfcCode))
	//if err != nil && err != types.ErrNotFound {
	//	return nil,err
	//}
	var goodinfo g.TraceplatformGoodInfo
	goodinfo = *payload.Goodinfo
	kv = append(kv,&types.KeyValue{Key:traceplatformKeyGood(payload.Goodinfo.NfcCode),Value:types.Encode(&goodinfo)})
	a.saveStateDB(kv)

	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//批量添加商品
func (a *Action) AddGoods(payload *g.TraceplatformAddGoods, pubkey []byte) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	for _, v := range payload.NfcCode {
		//_,err := a.db.Get(traceplatformKeyGood(v))
		//if err != nil && err != types.ErrNotFound {
		//	return nil,err
		//}
		var goodinfo g.TraceplatformGoodInfo
		goodinfo.Goodinfo = payload.Goodinfo
		goodinfo.NfcCode = v
		goodinfo.Name = payload.Name
		goodinfo.Amount = 1
		kv = append(kv,&types.KeyValue{Key:traceplatformKeyGood(v),Value:types.Encode(&goodinfo)})
	}
	a.saveStateDB(kv)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}


func (a *Action) saveStateDB(kv []*types.KeyValue) {
	for i:=0; i<len(kv) ; i++  {
		a.db.Set(kv[i].Key,kv[i].Value)
	}
}