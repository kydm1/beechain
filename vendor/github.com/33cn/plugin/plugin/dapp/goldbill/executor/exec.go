package executor

import (
	g "github.com/33cn/plugin/plugin/dapp/goldbill/types"
	"github.com/33cn/chain33/types"
)

func (g *Goldbill) Exec_InitPlatform(payload *g.GoldbillInitPlatform, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.InitPlatform(payload)
}

func (g *Goldbill) Exec_RegisterUser(payload *g.GoldbillRegisterUser, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.RegisterUser(payload)
}

func (g *Goldbill) Exec_SetFee(payload *g.GoldbillSetFee, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.SetFee(payload)
}

func (g *Goldbill) Exec_SetStamps(payload *g.GoldbillSetStamps, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.SetStamps(payload)
}

func (g *Goldbill) Exec_SetBail(payload *g.GoldbillSetBail, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.SetBail(payload)
}

func (g *Goldbill) Exec_Withdraw(payload *g.GoldbillWithdraw, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.Withdraw(payload)
}

func (g *Goldbill) Exec_Deposit(payload *g.GoldbillDeposit, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.Deposit(payload)
}

func (g *Goldbill) Exec_DepositKcoin(payload *g.GoldbillDepositKcoin, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.DepositKcoin(payload)
}

func (g *Goldbill) Exec_WithdrawKcoin(payload *g.GoldbillWithdrawKcoin, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.WithdrawKcoin(payload)
}

func (g *Goldbill) Exec_PayKCoin(payload *g.GoldbillPayKCoin, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.PayKcoin(payload)
}

func (g *Goldbill) Exec_Invoice(payload *g.GoldbillInvoice, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.Invoice(payload)
}

func (g *Goldbill) Exec_AddBill(payload *g.GoldbillAddBill, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.AddBill(payload)
}

func (g *Goldbill) Exec_BatchSell(payload *g.GoldbillBatchSell, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.BatchSell(payload)
}

func (g *Goldbill) Exec_BatchCancelSell(payload *g.GoldbillBatchCancelSell, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.BatchCancelSell(payload)
}

func (g *Goldbill) Exec_BatchBuy(payload *g.GoldbillBatchBuy, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.BatchBuy(payload)
}

func (g *Goldbill) Exec_BatchPay(payload *g.GoldbillBatchPay, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.BatchPay(payload)
}

func (g *Goldbill) Exec_BatchTakeOut(payload *g.GoldbillBatchTakeOut, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.BatchTakeOut(payload)
}

func (g *Goldbill) Exec_SliceBill(payload *g.GoldbillSliceBill, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.SliceBill(payload)
}

func (g *Goldbill) Exec_BatchRegUser(payload *g.GoldbillBatchRegUser, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.BatchRegUser(payload)
}

