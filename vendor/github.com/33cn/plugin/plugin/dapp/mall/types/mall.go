package types

import (
	"github.com/33cn/chain33/types"
	"reflect"
)

var MallX = "mall"

func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(MallX))
	types.RegistorExecutor(MallX, NewType())
	types.RegisterDappFork(MallX, "Enable", 0)
}

type MallType struct {
	types.ExecTypeBase
}

func NewType() *MallType {
	c := &MallType{}
	c.SetChild(c)
	return c
}

func (g *MallType) GetPayload() types.Message {
	return &MallAction{}
}

func (g *MallType) GetTypeMap() map[string]int32 {
	return map[string]int32{
		"MallUserRegister":		MallUserRegisterAction,
		"MallPlatformDeposit":		MallPlatformDepositAction,
		"MallUserWithdraw":		MallUserWithdrawAction,
		"MallUserGive":		MallUserGiveAction,
		"MallAddGood":		MallAddGoodAction,
		"MallPay":			MallPayAction,
		"MallDelivery":					MallDeliveryAction,
	}
}

func (g *MallType) GetLogMap() map[int64]*types.LogInfo {
	return map[int64]*types.LogInfo{
		TyLogMallUserRegister: {reflect.TypeOf(MallUserInfo{}), "TyLogMallUserRegister"},
		TyLogMallPlatformDeposit: {reflect.TypeOf(MallUserInfo{}), "TyLogMallPlatformDeposit"},
		TyLogMallUserWithdraw: {reflect.TypeOf(MallUserInfo{}), "TyLogMallUserWithdraw"},
		TyLogMallUserGive: {reflect.TypeOf(MallUserInfo{}), "TyLogMallUserGive"},
		TyLogMallAddGood: {reflect.TypeOf(MallUserInfo{}), "TyLogMallAddGood"},
		TyLogMallPay: {reflect.TypeOf(MallUserInfo{}), "TyLogMallPay"},
		TyLogMallDelivery: {reflect.TypeOf(MallUserInfo{}), "TyLogMallDelivery"},
	}
}
