package types

import (
	"github.com/33cn/chain33/types"

	"reflect"
)

var GoldbillX = "goldbill"

func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(GoldbillX))
	types.RegistorExecutor(GoldbillX, NewType())
	types.RegisterDappFork(GoldbillX, "Enable", 0)
}

type GoldbillType struct {
	types.ExecTypeBase
}

func NewType() *GoldbillType {
	c := &GoldbillType{}
	c.SetChild(c)
	return c
}

func (goldbill *GoldbillType) GetName() string {
	return GoldbillX
}

func (goldbill *GoldbillType) GetPayload() types.Message {
	return &GoldbillAction{}
}

func (goldbill *GoldbillType) GetTypeMap() map[string]int32 {
	return map[string]int32{
		"InitPlatform":		GoldbillActionType_InitPlatform,
		"RegisterUser": 	GoldbillActionType_RegisterUser,
		"SetFee":       	GoldbillActionType_SetFee,
		"SetStamps":    	GoldbillActionType_SetStamps,
		"SetBail":      	GoldbillActionType_SetBail,
		"Withdraw":     	GoldbillActionType_Withdraw,
		"Deposit":      	GoldbillActionType_Deposit,
		"DepositKcoin": 	GoldbillActionType_DepositKcoin,
		"WithdrawKcoin":    GoldbillActionType_WithdrawKcoin,
		"PayKCoin":     	GoldbillActionType_PayKCoin,
		"Invoice":   		GoldbillActionType_Invoice,
		"AddBill":			GoldbillActionType_AddBill,
		"BatchSell":		GoldbillActionType_BatchSell,
		"BatchCancelSell" : GoldbillActionType_BatchCancelSell,
		"BatchBuy":			GoldbillActionType_BatchBuy,
		"BatchPay" :		GoldbillActionType_BatchPay,
		"BatchTakeOut" : 	GoldbillActionType_BatchTakeOut,
		"SliceBill" : 		GoldbillActionType_SliceBill,
		"BatchRegUser" : 	GoldbillActionType_BatchRegUser,
	}
}

func (goldbill *GoldbillType) GetLogMap() map[int64]*types.LogInfo {
	return map[int64]*types.LogInfo{
		TyLogGoldbillInitPlatform: {reflect.TypeOf(GoldbillResponseInitPlatform{}), "TyLogGoldbillInitPlatform"},
		TyLogGoldbillRegisterUser: {reflect.TypeOf(GoldbillResponseRegisterUser{}), "TyLogGoldbillRegisterUser"},
		TyLogGoldbillSetFee: {reflect.TypeOf(GoldbillResponseSetFee{}), "TyLogGoldbillSetFee"},
		TyLogGoldbillSetStamps: {reflect.TypeOf(GoldbillResponseSetStamps{}), "TyLogGoldbillSetStamps"},
		TyLogGoldbillSetBail: {reflect.TypeOf(GoldbillResponseSetBail{}), "TyLogGoldbillSetBail"},
		TyLogGoldbillWithdraw: {reflect.TypeOf(GoldbillResponseWithdraw{}), "TyLogGoldbillWithdraw"},
		TyLogGoldbillDeposit: {reflect.TypeOf(GoldbillResponseDeposit{}), "TyLogGoldbillDeposit"},
		TyLogGoldbillDepositKcoin: {reflect.TypeOf(GoldbillResponseDepositKcoin{}), "TyLogGoldbillDepositKcoin"},
		TyLogGoldbillWithdrawKcoin: {reflect.TypeOf(GoldbillResponseWithdrawKcoin{}), "TyLogGoldbillWithdrawKcoin"},
		TyLogGoldbillPayKCoin: {reflect.TypeOf(GoldbillResponsePayKCoin{}), "TyLogGoldbillPayKcoin"},
		TyLogGoldbillInvoice: {reflect.TypeOf(GoldbillResponseInvoice{}), "TyLogGoldbillInvoice"},
		TyLogGoldbillAddBill: {reflect.TypeOf(GoldbillResponseAddBill{}), "TyLogGoldbillAddBill"},
		TyLogGoldbillBatchSell: {reflect.TypeOf(GoldbillResponseBatchSell{}), "TyLogGoldbillBatchSell"},
		TyLogGoldbillBatchCancelSell: {reflect.TypeOf(GoldbillResponseBatchCancelSell{}), "TyLogGoldbillBatchCancelSell"},
		TyLogGoldbillBatchBuy: {reflect.TypeOf(GoldbillResponseBatchBuy{}), "TyLogGoldbillBatchBuy"},
		TyLogGoldbillBatchPay: {reflect.TypeOf(GoldbillResponseBatchPay{}), "TyLogGoldbillBatchPay"},
		TyLogGoldbillBatchTakeOut: {reflect.TypeOf(GoldbillResponseBatchTakeOut{}), "TyLogGoldbillBatchTakeOut"},
		TyLogGoldbillSliceBill: {reflect.TypeOf(GoldbillResponseSliceBill{}), "TyLogGoldbillSliceBill"},
		TyLogGoldbillBatchRegUser: {reflect.TypeOf(GoldbillResponseBatchRegUser{}), "TyLogGoldbillBatchRegUser"},

	}
}
