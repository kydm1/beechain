package executor

import (
	g "github.com/33cn/plugin/plugin/dapp/gyl/types"
	"github.com/33cn/chain33/types"
)

func (g *Gyl) Exec_ZsgjSaveReceipt(payload *g.ZsgjSaveReceipt, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.SaveReceipt(payload,tx.Signature.Pubkey)
}

func (g *Gyl) Exec_ZsgjSaveProduct(payload *g.ZsgjSaveProduct, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.SaveProduct(payload,tx.Signature.Pubkey)
}

func (g *Gyl) Exec_AssetRegister(payload *g.AssetRegisterAction, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.AssetRegister(payload,tx.Signature.Pubkey)
}

func (g *Gyl) Exec_ZsgjApplyRecharge(payload *g.ZsgjApplyRecharge, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.ApplyRecharge(payload,tx.Signature.Pubkey)
}

func (g *Gyl) Exec_ZsgjApplyWithdraw(payload *g.ZsgjApplyWithdraw, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.ApplyWithdraw(payload,tx.Signature.Pubkey)
}

func (g *Gyl) Exec_ZsgjCash(payload *g.ZsgjCash, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.Cash(payload,tx.Signature.Pubkey)
}

func (g *Gyl) Exec_ZsgjClear(payload *g.ZsgjClear, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.Clear(payload,tx.Signature.Pubkey)
}

func (g *Gyl) Exec_ZsgjDelist(payload *g.ZsgjDelist, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.Delist(payload,tx.Signature.Pubkey)
}

func (g *Gyl) Exec_ZsgjCompanyCertification(payload *g.ZsgjCompanyCertification, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.CompanyCertification(payload,tx.Signature.Pubkey)
}

func (g *Gyl) Exec_ZsgjPersonCertification(payload *g.ZsgjPersonCertification, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.PersonCertification(payload)
}

func (g *Gyl) Exec_DlReceiptPay(payload *g.DlReceiptPay, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.ReceiptPay(payload,tx.Signature.Pubkey)
}

func (g *Gyl) Exec_DlReceiptDelist(payload *g.DlReceiptDelist, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.ReceiptDelist(payload,tx.Signature.Pubkey)
}

func (g *Gyl) Exec_DlBlankNote(payload *g.DlBlankNote, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.BlankNote(payload,tx.Signature.Pubkey)
}

func (g *Gyl) Exec_DlReceiptList(payload *g.DlReceiptList, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.ReceiptList(payload,tx.Signature.Pubkey)
}

func (g *Gyl) Exec_DlBlankNotePay(payload *g.DlBlankNotePay, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.BlankNotePay(payload,tx.Signature.Pubkey)
}

func (g *Gyl) Exec_DlReceiptAndNoteCash(payload *g.DlReceiptAndNoteCash, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.ReceiptAndNoteCash(payload,tx.Signature.Pubkey)
}

func (g *Gyl) Exec_DlInvestCash(payload *g.DlInvestCash, tx *types.Transaction, index int) (*types.Receipt, error) {
	//action := NewAction(g, tx, index)
	//return action.InvestCash(payload,tx.Signature.Pubkey)
	return nil,nil
}

func (g *Gyl) Exec_DlCredit(payload *g.DlCredit, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.SetCredit(payload,tx.Signature.Pubkey)
}

func (g *Gyl) Exec_AdAssetRegister(payload *g.AdAssetRegister, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.AdAssetRegister(payload,tx.Signature.Pubkey)
}

func (g *Gyl) Exec_AdApplyBlankNote(payload *g.AdApplyBlankNote, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.AdApplyBlankNote(payload,tx.Signature.Pubkey)
}

func (g *Gyl) Exec_AdAndNoteCash(payload *g.AdAndNoteCash, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.AdAndNoteCash(payload,tx.Signature.Pubkey)
}

func (g *Gyl) Exec_GylFinanceInfo(payload *g.GylFinanceInfo, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.GylFinanceInfo(payload,tx.Signature.Pubkey)
}

func (g *Gyl) Exec_GylDirectFinance(payload *g.GylDirectFinance, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.GylDirectFinance(payload,tx.Signature.Pubkey)
}

