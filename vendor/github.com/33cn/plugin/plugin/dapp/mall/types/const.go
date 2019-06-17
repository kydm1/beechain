package types

//mall action
const (
	MallActionType = iota
	MallUserRegisterAction
	MallPlatformDepositAction
	MallUserWithdrawAction
	MallUserGiveAction
	MallAddGoodAction
	MallPayAction
	MallDeliveryAction

)

const (
	//log for mall

	TyLogMallUserRegister  = 900
	TyLogMallPlatformDeposit = 901
	TyLogMallUserWithdraw = 902
	TyLogMallUserGive = 903
	TyLogMallAddGood = 904
	TyLogMallPay = 905
	TyLogMallDelivery = 906
)

