package executor

import (
	g "github.com/33cn/plugin/plugin/dapp/traceplatform/types"
	"github.com/33cn/chain33/types"
)

func (g *Traceplatform) Exec_TraceplatformAddGood(payload *g.TraceplatformAddGood, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.AddGood(payload,tx.Signature.Pubkey)
}

func (g *Traceplatform) Exec_TraceplatformAddGoods(payload *g.TraceplatformAddGoods, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.AddGoods(payload,tx.Signature.Pubkey)
}
