package executor

import (
	gty "github.com/33cn/plugin/plugin/dapp/goldbill/types"
	"github.com/33cn/chain33/types"
)

func (g *Goldbill) ExecLocal_InitPlatform(payload *gty.GoldbillInitPlatform, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return g.execLocal(receipt)
}

func (g *Goldbill) ExecLocal_RegisterUser(payload *gty.GoldbillRegisterUser, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return g.execLocal(receipt)
}

func (g *Goldbill) ExecLocal_SetFee(payload *gty.GoldbillSetFee, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return g.execLocal(receipt)
}

func (g *Goldbill) ExecLocal_SetStamps(payload *gty.GoldbillSetStamps, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return g.execLocal(receipt)
}

func (g *Goldbill) ExecLocal_SetBail(payload *gty.GoldbillSetBail, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return g.execLocal(receipt)
}

func (g *Goldbill) ExecLocal_Withdraw(payload *gty.GoldbillWithdraw, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return g.execLocal(receipt)
}

func (g *Goldbill) ExecLocal_Deposit(payload *gty.GoldbillDeposit, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return g.execLocal(receipt)
}

func (g *Goldbill) ExecLocal_DepositKcoin(payload *gty.GoldbillDepositKcoin, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return g.execLocal(receipt)
}

func (g *Goldbill) ExecLocal_WithdrawKcoin(payload *gty.GoldbillWithdrawKcoin, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return g.execLocal(receipt)
}

func (g *Goldbill) ExecLocal_PayKCoin(payload *gty.GoldbillPayKCoin, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return g.execLocal(receipt)
}

func (g *Goldbill) ExecLocal_Invoice(payload *gty.GoldbillInvoice, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return g.execLocal(receipt)
}

func (g *Goldbill) ExecLocal_AddBill(payload *gty.GoldbillAddBill, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return g.execLocal(receipt)
}

func (g *Goldbill) ExecLocal_BatchSell(payload *gty.GoldbillBatchSell, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return g.execLocal(receipt)
}

func (g *Goldbill) ExecLocal_BatchCancelSell(payload *gty.GoldbillBatchCancelSell, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return g.execLocal(receipt)
}

func (g *Goldbill) ExecLocal_BatchBuy(payload *gty.GoldbillBatchBuy, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return g.execLocal(receipt)
}

func (g *Goldbill) ExecLocal_BatchPay(payload *gty.GoldbillBatchPay, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return g.execLocal(receipt)
}

func (g *Goldbill) ExecLocal_BatchTakeOut(payload *gty.GoldbillBatchTakeOut, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return g.execLocal(receipt)
}

func (g *Goldbill) ExecLocal_SliceBill(payload *gty.GoldbillSliceBill, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return g.execLocal(receipt)
}

func (g *Goldbill) ExecLocal_BatchRegUser(payload *gty.GoldbillBatchRegUser, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return g.execLocal(receipt)
}

func (g *Goldbill) execLocal(receipt *types.ReceiptData) (*types.LocalDBSet,error) {
	dbSet := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return dbSet, nil
	}
	for _, log := range receipt.Logs {
		switch log.Ty {
		case gty.TyLogGoldbillRegisterUser:
			var r gty.GoldbillResponseRegisterUser
			types.Decode(log.Log,&r)
			kv ,err := g.UpdateUserLocalState(r.GoldbillUser.Type,true)
			if err != nil {
				return nil,err
			}
			dbSet.KV = append(dbSet.KV, kv...)
		case gty.TyLogGoldbillBatchRegUser:
			var r gty.GoldbillResponseBatchRegUser
			types.Decode(log.Log,&r)
			for i:=0;i<len(r.Succ) ;i++ {
				kv, err := g.UpdateUserLocalState(gty.GoldbillUserType_UT_USER, true)
				if err != nil {
					return nil, err
				}
				dbSet.KV = append(dbSet.KV, kv...)
			}
		}
	}
	return dbSet, nil
}