package executor

import (
	"github.com/33cn/chain33/types"
	g "github.com/33cn/plugin/plugin/dapp/gyl/types"
	dbm "github.com/33cn/chain33/common/db"
	"fmt"
)


var (
	gylPrefix = "mavl-gyl-"
)

func gylKeyUser(id string) (key []byte) {
	return []byte(fmt.Sprintf(gylPrefix+"user-"+"%s", id))
}

func gylKeyUserInfo(id []byte) (key []byte) {
	key = append(key,[]byte(gylPrefix+"account-")...)
	key = append(key,id...)
	return key
}

func gylKeyFinanceInfo(id []byte) (key []byte) {
	key = append(key,[]byte(gylPrefix+"finfo-")...)
	key = append(key,id...)
	return key
}

func gylKeyProduct(id string) (key []byte) {
	return []byte(fmt.Sprintf(gylPrefix+"pid-"+"%s", id))
}

func gylKeyDelist(id string) (key []byte) {
	return []byte(fmt.Sprintf(gylPrefix+"delist-"+"%s", id))
}

func gylKeyReceipt(id string) (key []byte) {
	return []byte(fmt.Sprintf(gylPrefix+"rid-"+"%s", id))
}

func gylKeyNote(id string) (key []byte) {
	return []byte(fmt.Sprintf(gylPrefix+"noteid-"+"%s", id))
}

func gylKeyBlankNote(id string) (key []byte) {
	return []byte(fmt.Sprintf(gylPrefix+"bid-"+"%s", id))
}

func gylKeyAdance(id string) (key []byte) {
	return []byte(fmt.Sprintf(gylPrefix+"aid-"+"%s", id))
}

func gylKeyWithdraw(id string) (key []byte) {
	return []byte(fmt.Sprintf(gylPrefix+"out-"+"%s", id))
}

func gylKeyDeposit(id string) (key []byte) {
	return []byte(fmt.Sprintf(gylPrefix+"in-"+"%s", id))
}

type receiptInfo struct {
	zri g.ZsgjReceiptInfo
}

type productInfo struct {
	zpi g.ZsgjProductInfo
}

type blankNoteInfo struct {
	bzi g.ZsgjBlankNoteInfo
}

func Key(id string) (key []byte) {
	key = append(key, []byte("mavl-gyl-")...)
	key = append(key, []byte(id)...)
	return key
}

//用ZsgjSaveReceipt结构给ReceiptInfo结构赋值
func NewReceiptInfo(receipt *g.ZsgjSaveReceipt) *receiptInfo {
	t := &receiptInfo{}
	dl := &g.ZsgjReceiptInfo{}
	dl.State = receipt.ZsgjReceipt.State
	dl.ReceiptId = receipt.ZsgjReceipt.ReceiptId
	dl.ReceiptViceId = receipt.ZsgjReceipt.ReceiptViceId
	dl.CoreCompany = receipt.ZsgjReceipt.CoreCompany
	dl.ReceiveCompany = receipt.ZsgjReceipt.ReceiveCompany
	dl.SumAmount = receipt.ZsgjReceipt.SumAmount
	dl.IssuedAgency = receipt.ZsgjReceipt.IssuedAgency
	dl.FinanceAmount = receipt.ZsgjReceipt.FinanceAmount
	dl.Opty = receipt.ZsgjReceipt.Opty
	dl.InvestAgency = receipt.ZsgjReceipt.InvestAgency
	t.zri = *dl
	return t
}

//用ZsgjSaveProduct结构给ProductInfo结构赋值
func NewProductInfo(product *g.ZsgjSaveProduct) *productInfo {
	t := &productInfo{}
	dl := &g.ZsgjProductInfo{}
	dl.ProductId = product.ZsgjProduct.GetProductId()
	dl.ProductName = product.ZsgjProduct.GetProductName()
	dl.AnnualizedRate = product.ZsgjProduct.GetAnnualizedRate()
	dl.IssuedAgency = product.ZsgjProduct.GetIssuedAgency()
	dl.Frequency = product.ZsgjProduct.GetFrequency()
	dl.PayInterestMethod = product.ZsgjProduct.GetPayInterestMethod()
	dl.CalInterestMethod = product.ZsgjProduct.GetCalInterestMethod()
	dl.CoinType = product.ZsgjProduct.GetCoinType()
	dl.GuaranteeAgency = product.ZsgjProduct.GetGuaranteeAgency()
	dl.HostingAgency = product.ZsgjProduct.GetHostingAgency()
	dl.IssueScale = product.ZsgjProduct.GetIssueScale()
	dl.SaleTarget = product.ZsgjProduct.GetSaleTarget()
	dl.RiskLevel = product.ZsgjProduct.GetRiskLevel()
	dl.RaiseStartDate = product.ZsgjProduct.GetRaiseStartDate()
	dl.RaiseEndDate = product.ZsgjProduct.GetRaiseEndDate()
	dl.CalInterestStart = product.ZsgjProduct.GetCalInterestStart()
	dl.CalInterestEnd = product.ZsgjProduct.GetCalInterestEnd()
	dl.StartAmount = product.ZsgjProduct.GetStartAmount()
	dl.IncreamentingAmount = product.ZsgjProduct.GetIncreamentingAmount()
	dl.SubscriptionPerson = product.ZsgjProduct.GetSubscriptionPerson()
	dl.IsTranferable = product.ZsgjProduct.GetIsTranferable()
	dl.Info = product.ZsgjProduct.GetInfo()
	dl.TransferStart = product.ZsgjProduct.GetTransferStart()
	dl.TransferEnd = product.ZsgjProduct.GetTransferEnd()
	dl.MinHoldingDuration = product.ZsgjProduct.GetMinHoldingDuration()
	dl.IstransferRateFloat = product.ZsgjProduct.GetIstransferRateFloat()
	dl.MinTransferAmount = product.ZsgjProduct.GetMinTransferAmount()
	dl.RateFloatInterregional = product.ZsgjProduct.GetRateFloatInterregional()
	dl.TransferIncreamentingAmount = product.ZsgjProduct.GetTransferIncreamentingAmount()
	dl.TransferFeeTarget = product.ZsgjProduct.GetTransferFeeTarget()
	dl.Opty = product.ZsgjProduct.GetOpty()
	dl.OpCompany = product.ZsgjProduct.GetOpCompany()
	dl.ReceiptId = product.ZsgjProduct.GetReceiptId()
	t.zpi = *dl
	return t
}

//用BlankNote结构给BlankNoteInfo结构赋值
func NewBlankNoteInfo(blank *g.DlBlankNote) *blankNoteInfo {
	t := &blankNoteInfo{}
	dl := &g.ZsgjBlankNoteInfo{}
	dl.CoreCompany = blank.GetCoreCompany()
	dl.UpstreamFirm = blank.GetUpstreamFirm()
	dl.Amount = blank.GetAmount()
	dl.CashAgency = blank.GetCashAgency()
	dl.CompleteDate = blank.GetCompleteDate()
	dl.CashDate = blank.GetCashDate()
	dl.ApplyBlankNote = blank.GetApplyBlankNote()
	dl.BlankNoteAmount = blank.GetBlankNoteAmount()
	dl.OverdueRate = blank.GetOverdueRate()
	dl.LeftAmount = blank.GetLeftAmount()
	dl.OpCompany = blank.GetOpCompany()
	dl.State = blank.GetState()
	dl.Opty = blank.GetOpty()
	dl.ReceiptId = blank.GetReceiptId()
	dl.BlankNoteId = blank.GetBlankNoteId()
	dl.BlankNoteParentsId = blank.GetBlankNoteParentsId()
	dl.ReceiptId = blank.GetReceiptId()
	t.bzi = *dl
	return t
}


func saveInfo(info *g.AssetRegisterAction, db dbm.KVDB) (*types.KeyValue, error) {
	value := types.Encode(info)
	kv := &types.KeyValue{Key: calCtrKey(info.ContractNo), Value: value}
	db.Set(kv.Key, kv.Value)
	return kv, nil
}
func saveComCerInfo(info *g.ZsgjCompanyCertification, db dbm.KVDB) (*types.KeyValue, error) {
	value := types.Encode(info)
	kv := &types.KeyValue{Key: calComKey(info.PhoneNumber), Value: value}
	db.Set(kv.Key, kv.Value)
	return kv, nil
}
func savePerCerInfo(info *g.ZsgjPersonCertification, db dbm.KVDB) (*types.KeyValue, error) {
	value := types.Encode(info)
	kv := &types.KeyValue{Key: calPerKey(info.IdCard), Value: value}
	db.Set(kv.Key, kv.Value)
	return kv, nil
}
func calPerKey(str string) []byte {
	key := fmt.Sprintf("mavl-gyl-personCertification:%s", str)
	return []byte(key)
}
func calCtrKey(str string) []byte {
	key := fmt.Sprintf("mavl-gyl-contract:%s", str)
	return []byte(key)
}
func calComKey(str string) []byte {
	key := fmt.Sprintf("mavl-gyl-companyCertification:%s", str)
	return []byte(key)
}
func isNumberExist(no string, db dbm.KVDB) (bool, error) {
	if db == nil {
		clog.Info("ZSGJ", "contract", "contract db is empty, contractNo available")
		return false, nil
	}
	//合同编号数据
	_, err := db.Get(calCtrKey(no))
	if err != nil {
		clog.Info("ZSGJ", "contract", "contractNo available")
		return false, err
	}
	return true, nil
}
