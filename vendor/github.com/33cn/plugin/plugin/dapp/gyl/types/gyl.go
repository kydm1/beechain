package types

import (
	"github.com/33cn/chain33/types"

)

var GylX = "gyl"

func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(GylX))
	types.RegistorExecutor(GylX, NewType())
	types.RegisterDappFork(GylX, "Enable", 0)
}

type GylType struct {
	types.ExecTypeBase
}

func NewType() *GylType {
	c := &GylType{}
	c.SetChild(c)
	return c
}

func (g *GylType) GetPayload() types.Message {
	return &ZsgjAction{}
}

func (g *GylType) GetTypeMap() map[string]int32 {
	return map[string]int32{
		"ZsgjSaveReceipt":				ZsgjSaveReceiptAction,
		"ZsgjSaveProduct": 				ZsgjSaveProductAction,
		"AssetRegister":       			ZsgjAssetRegisterAction,
		"ZsgjApplyRecharge":    		ZsgjApplyRechargeAction,
		"ZsgjApplyWithdraw":      		ZsgjApplyWithdrawAction,
		"ZsgjCash":     				ZsgjCashAction,
		"ZsgjClear":      				ZsgjClearAction,
		"ZsgjDelist": 					ZsgjDelistAction,
		"ZsgjCompanyCertification":     ZsgjCompanyCertificationAction,
		"ZsgjPersonCertification":     	ZsgjPersonCertificationAction,
		"DlReceiptPay":   				DlReceiptPayAction,
		"DlReceiptDelist":				DlReceiptDelistAction,
		"DlBlankNote":					DlBlankNoteAction,
		"DlReceiptList" : 				DlReceiptListAction,
		"DlBlankNotePay":				DlBlankNotePayAction,
		"DlReceiptAndNoteCash" :		DlReceiptAndNoteCashAction,
		"DlInvestCash":					DlInvestCashAction,
		"DlCredit" : 					DlCreditAction,
		"AdAssetRegister" : 			AdAssetRegisterAction,
		"AdApplyBlankNote" : 			AdApplyBlankNoteAction,
		"AdAndNoteCash":				AdAndNoteCashAction,
		"GylFinanceInfo":         		GylFinanceInfoAction,
		"GylDirectFinance":       		GylDirectFinanceAction,
	}
}

func (g *GylType) GetLogMap() map[int64]*types.LogInfo {
	return map[int64]*types.LogInfo{}
}
