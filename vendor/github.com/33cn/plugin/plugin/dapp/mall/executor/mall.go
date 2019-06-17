package executor

import (
	g "github.com/33cn/plugin/plugin/dapp/mall/types"
	log "github.com/33cn/chain33/common/log/log15"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	"bytes"
	"encoding/hex"
	"github.com/33cn/chain33/common/address"
)

var clog = log.New("module", "execs.mall")
var driverName = "mall"
var adminkey = "0307b368e8c2e9a7f9ed5f4b85a845ee6cae9fac7b9743e1dddfa823161284a13c"
var currencyBty = "BTY"

func init() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&Mall{}))
}

func Init(name string, sub []byte) {
	clog.Debug("register mall execer")
	drivers.Register(GetName(), newMall, types.GetDappFork(driverName, "Enable"))
}

func GetName() string {
	return newMall().GetName()
}

type Mall struct {
	drivers.DriverBase
}

func newMall() drivers.Driver {
	n := &Mall{}
	n.SetChild(n)
	n.SetIsFree(true)
	n.SetExecutorType(types.LoadExecutorType(driverName))
	return n
}

func (n *Mall) GetDriverName() string {
	return driverName
}

//简单对交易数据和签名用户做下校验
func (n *Mall) CheckTx(tx *types.Transaction, index int) error {
	var payload g.MallAction
	err := types.Decode(tx.Payload,&payload)
	if err != nil {
		return g.ErrTxErr
	}
	if payload.Value == nil {
		return g.ErrEmptyValue
	}
	switch payload.Ty {
	case g.MallUserRegisterAction:
	case g.MallPlatformDepositAction:
		//return n.checkPlatformDeposit(payload.GetMallPlatformDeposit(),tx.Signature.Pubkey)
	case g.MallUserWithdrawAction:
		//return n.checkUserWithdraw(payload.GetMallUserWithdraw(),tx.Signature.Pubkey)
	case g.MallUserGiveAction:
		//return n.checkUserGive(payload.GetMallUserGive(),tx.Signature.Pubkey)
	case g.MallAddGoodAction:
		//return n.checkAddGood(payload.GetMallAddGood(),tx.Signature.Pubkey)
	case g.MallPayAction:
		//return n.checkPay(payload.GetMallPay(),tx.Signature.Pubkey)
	case g.MallDeliveryAction:
		//return n.checkUser(tx.Signature.Pubkey)
	default :
		return g.ErrWrongActionType
	}
	return nil
}

func (n *Mall) checkPlatformDeposit(payload *g.MallPlatformDeposit,pubkey []byte) error {
	clog.Info("MALL","checkPlatformDeposit","start")
	err := n.checkAdmin(pubkey)
	if err != nil {
		return err
	}
	for _,v := range payload.GetDeposit() {
		addr := address.PubKeyToAddress(v.Pubkey).String()
		_,err := n.GetStateDB().Get(mallKeyUser(addr))
		if err != nil {
			return err
		}
	}
	clog.Info("MALL","checkPlatformDeposit","end")
	return nil
}

func (n *Mall) checkUserWithdraw(payload *g.MallUserWithdraw,pubkey []byte) error {

	err := n.checkUser(pubkey)
	if err != nil {
		return err
	}
	if payload.Name == "" || payload.Amount <= 0 {
		return g.ErrWrongData
	}
	return nil
}

func (n *Mall) checkUserGive(payload *g.MallUserGive,pubkey []byte) error {

	err := n.checkUser(pubkey)
	if err != nil {
		return err
	}
	err = n.checkUser(payload.To)
	if err != nil {
		return err
	}
	if payload.Name == "" || payload.Amount <= 0 {
		return g.ErrWrongData
	}
	return nil
}

func (n *Mall) checkAddGood(payload *g.MallAddGood,pubkey []byte) error {

	err := n.checkUser(pubkey)
	if err != nil {
		return err
	}
	if payload.BaseInfo == nil || payload.OtherInfo == nil || payload.SpecInfo == nil {
		return g.ErrEmptyValue
	}
	if payload.BaseInfo.GoodId == "" || payload.BaseInfo.GoodName == "" {
		return g.ErrEmptyValue
	}
	return nil
}

func (n *Mall) checkPay(payload *g.MallPay,pubkey []byte) error {

	err := n.checkUser(pubkey)
	if err != nil {
		return err
	}
	for _,v := range payload.SinglePay {
		if v.PayOrder == nil || v.PayInfo == nil {
			return g.ErrEmptyValue
		}
		if bytes.Equal(v.PayOrder.BuyerPubkey,v.PayOrder.SellerPubkey) {
			return g.ErrWrongData
		}
		err = n.checkUser(v.PayOrder.BuyerPubkey)
		if err != nil {
			return err
		}
		err = n.checkUser(v.PayOrder.SellerPubkey)
		if err != nil {
			return err
		}
	}


	return nil
}

func (n *Mall) checkAdmin(pubkey []byte) error {

	platformKey,err := hex.DecodeString(adminkey)
	if err != nil {
		return err
	}
	if !bytes.Equal(pubkey,platformKey) {
		return g.ErrWrongPlatformKey
	}

	return nil
}

func (n *Mall) checkUser(pubkey []byte) error {
	addr := address.PubKeyToAddress(pubkey).String()
	_,err := n.GetStateDB().Get(mallKeyUser(addr))
	if err != nil {
		clog.Error("mall","checkUser","sign user not exists")
		return err
	}
	return nil
}