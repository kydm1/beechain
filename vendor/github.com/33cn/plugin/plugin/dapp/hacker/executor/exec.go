package executor

import (
	g "github.com/33cn/plugin/plugin/dapp/hacker/types"
	"github.com/33cn/chain33/types"
)

func (g *Hacker) Exec_AddBill(payload *g.HackerAddBill, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.AddBill(payload,tx.Signature.Pubkey)
}
