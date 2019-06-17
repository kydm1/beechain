package executor

import (
	gty "github.com/33cn/plugin/plugin/dapp/mall/types"
	"github.com/33cn/chain33/types"
)

func (g *Mall) ExecLocal_MallUserRegister(payload *gty.MallUserRegister, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {

	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}

	return set,nil
}

func (g *Mall) ExecLocal_MallPlatformDeposit(payload *gty.MallPlatformDeposit, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {

	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}

	return set,nil
}

func (g *Mall) ExecLocal_MallUserWithdraw(payload *gty.MallUserWithdraw, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {

	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}

	return set,nil
}

func (g *Mall) ExecLocal_MallUserGive(payload *gty.MallUserGive, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {

	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}

	return set,nil
}

func (g *Mall) ExecLocal_MallAddGood(payload *gty.MallAddGood, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {

	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}

	return set,nil
}

func (g *Mall) ExecLocal_MallPay(payload *gty.MallPay, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {

	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}

	return set,nil
}


func (g *Mall) ExecLocal_MallDelivery(payload *gty.MallDelivery, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {

	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}

	return set,nil
}