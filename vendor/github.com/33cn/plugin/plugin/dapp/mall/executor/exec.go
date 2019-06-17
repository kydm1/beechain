package executor

import (
	g "github.com/33cn/plugin/plugin/dapp/mall/types"
	"github.com/33cn/chain33/types"
)

func (g *Mall) Exec_MallUserRegister(payload *g.MallUserRegister, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.UserRegister(payload)
}

func (g *Mall) Exec_MallPlatformDeposit(payload *g.MallPlatformDeposit, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.PlatformDeposit(payload)
}

func (g *Mall) Exec_MallUserWithdraw(payload *g.MallUserWithdraw, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.UserWithdraw(payload)
}

func (g *Mall) Exec_MallUserGive(payload *g.MallUserGive, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.UserGive(payload)
}

func (g *Mall) Exec_MallAddGood(payload *g.MallAddGood, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.AddGood(payload)
}

func (g *Mall) Exec_MallPay(payload *g.MallPay, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.Pay(payload)
}

func (g *Mall) Exec_MallDelivery(payload *g.MallDelivery, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.Delivery(payload)
}