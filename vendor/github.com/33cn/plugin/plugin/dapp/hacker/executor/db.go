package executor

import (
	"github.com/33cn/chain33/types"
	dbm "github.com/33cn/chain33/common/db"
	g "github.com/33cn/plugin/plugin/dapp/hacker/types"
)


type Action struct {
	db     dbm.KV
	txhash []byte
	height int64
	index  int
}

func NewAction(t *Hacker, tx *types.Transaction,index int) *Action {
	hash := tx.Hash()
	return &Action{t.GetStateDB(), hash, t.GetHeight(),index}
}

//添加商品
func (a *Action) AddBill(payload *g.HackerAddBill, pubkey []byte) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	kv = append(kv,&types.KeyValue{Key:hackerKeyGood(payload.StockNumber),Value:types.Encode(payload)})
	a.saveStateDB(kv)
	logs = append(logs,&types.ReceiptLog{Ty:g.TyLogHackerAddBill,Log:types.Encode(payload)})
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

func (a *Action) saveStateDB(kv []*types.KeyValue) {
	for i:=0; i<len(kv) ; i++  {
		a.db.Set(kv[i].Key,kv[i].Value)
	}
}