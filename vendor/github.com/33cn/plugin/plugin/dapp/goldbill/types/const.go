package types

//goldbill action
const (
	GoldbillActionType_InitPlatform     = 1
	GoldbillActionType_RegisterUser     = 2
	GoldbillActionType_SetFee           = 3
	GoldbillActionType_SetStamps        = 4
	GoldbillActionType_SetBail          = 5
	GoldbillActionType_Withdraw         = 6
	GoldbillActionType_Deposit          = 7
	GoldbillActionType_DepositKcoin     = 8
	GoldbillActionType_WithdrawKcoin    = 9
	GoldbillActionType_PayKCoin         = 10
	GoldbillActionType_Invoice          = 11
	GoldbillActionType_AddBill          = 12
	GoldbillActionType_BatchSell        = 13
	GoldbillActionType_BatchCancelSell  = 14
	GoldbillActionType_BatchBuy         = 15
	GoldbillActionType_BatchPay         = 16
	GoldbillActionType_BatchTakeOut     = 17
	GoldbillActionType_SliceBill        = 18
	GoldbillActionType_BatchRegUser     = 19
)


//goldbill log

const (
	TyLogGoldbillInitPlatform = 1
	TyLogGoldbillRegisterUser = 2
	TyLogGoldbillSetFee = 3
	TyLogGoldbillSetStamps = 4
	TyLogGoldbillSetBail = 5
	TyLogGoldbillWithdraw = 6
	TyLogGoldbillDeposit = 7
	TyLogGoldbillWithdrawKcoin = 8
	TyLogGoldbillDepositKcoin = 9
	TyLogGoldbillPayKCoin = 10
	TyLogGoldbillInvoice = 11
	TyLogGoldbillAddBill = 12
	TyLogGoldbillBatchSell = 13
	TyLogGoldbillBatchCancelSell = 14
	TyLogGoldbillBatchBuy = 15
	TyLogGoldbillBatchPay = 16
	TyLogGoldbillBatchTakeOut = 17
	TyLogGoldbillSliceBill = 18
	TyLogGoldbillBatchRegUser = 19
)