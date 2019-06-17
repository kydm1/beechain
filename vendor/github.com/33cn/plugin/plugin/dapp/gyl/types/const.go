package types

//goldbill action
const (
	ZsgjActionType = iota
	ZsgjSaveReceiptAction
	ZsgjSaveProductAction
	ZsgjAssetRegisterAction
	ZsgjApplyRechargeAction        //充值
	ZsgjApplyWithdrawAction        //提款
	ZsgjCashAction                 //兑付
	ZsgjClearAction                //清算
	ZsgjDelistAction               //摘牌
	ZsgjCompanyCertificationAction //企业认证
	ZsgjPersonCertificationAction  //个人认证
	DlReceiptPayAction             //单据支付
	DlReceiptDelistAction          //单据摘牌
	DlBlankNoteAction              //白条
	DlReceiptListAction            //单据挂牌
	DlBlankNotePayAction           //白条支付
	DlReceiptAndNoteCashAction     //白条、单据持有兑付
	DlInvestCashAction             //渠道企业兑付
	DlCreditAction                 //信用额度
	AdAssetRegisterAction          //预付款登记
	AdApplyBlankNoteAction		   //申请白条
	AdAndNoteCashAction		       //白条兑付   预付款兑付
	GylFinanceInfoAction           //融资信息
	GylDirectFinanceAction               //定向融资
)

const (
	//log for zsgj
	TyLogSaveReceipt   = 370
	TyLogSaveProduct   = 371
	TyLogAssetRegister = 372

)

