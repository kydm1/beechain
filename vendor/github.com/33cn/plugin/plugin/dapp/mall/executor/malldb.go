package executor

import (
	"github.com/33cn/chain33/types"
	dbm "github.com/33cn/chain33/common/db"
	g "github.com/33cn/plugin/plugin/dapp/mall/types"
	"github.com/33cn/chain33/common/address"
)


type Action struct {
	db     dbm.KV
	txhash []byte
	height int64
	pubkey []byte
	addr   string
	index  int
}

func NewAction(t *Mall, tx *types.Transaction,index int) *Action {
	hash := tx.Hash()
	return &Action{t.GetStateDB(), hash, t.GetHeight(),tx.Signature.Pubkey,address.PubKeyToAddress(tx.Signature.Pubkey).String(),index}
}

//用户注册
func (a *Action) UserRegister(payload *g.MallUserRegister) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	//检查用户是否已经存在\

	var userinfo g.MallUserInfo
	userinfo.Uid = payload.Uid
	userinfo.Phone = payload.Phone
	userinfo.Addr = a.addr

	value := types.Encode(&userinfo)
	kv = append(kv, &types.KeyValue{Key: mallKeyUser(a.addr), Value: value})
	logs = append(logs,&types.ReceiptLog{Ty:g.TyLogMallUserRegister,Log:types.Encode(&userinfo)})
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//平台充币
func (a *Action) PlatformDeposit(payload *g.MallPlatformDeposit) (*types.Receipt, error) {
	clog.Info("MALL","PlatformDeposit","start")
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	for _,v := range payload.Deposit {
		addr := address.PubKeyToAddress(v.Pubkey).String()
		userInfo,err := a.getUserInfo(addr)
		if err != nil {
			return nil,err
		}
		if v.IsToken && v.Amount >0 {
			c := false
			for _,vv := range userInfo.Token {
				if vv.Name == v.Name {
					vv.Amount += v.Amount
					c = true
					break
				}
			}
			if !c {
				newToken := &g.MallUserToken{Name:v.Name,Amount:v.Amount}
				userInfo.Token = append(userInfo.Token,newToken)
			}
		} else if !v.IsToken && v.Amount >0{
			c := false
			for _,vv := range userInfo.Currency {
				if vv.Name == v.Name {
					vv.Amount += v.Amount
					c = true
					break
				}
			}
			if !c {
				newCurrency := &g.MallUserCurrency{Name:v.Name,Amount:v.Amount}
				userInfo.Currency = append(userInfo.Currency,newCurrency)
			}

		}
		//及时保存在缓存中
		a.db.Set(mallKeyUser(addr),types.Encode(userInfo))
		kv = append(kv,&types.KeyValue{Key:mallKeyUser(addr),Value:types.Encode(userInfo)})
	}
	//logs = append(logs,&types.ReceiptLog{Ty:g.TyLogMallPlatformDeposit,Log:nil})
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//用户提币
func (a *Action) UserWithdraw(payload *g.MallUserWithdraw) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	//检查用户是否已经存在
	addr := a.addr
	userinfo,err := a.getUserInfo(addr)
	if err != nil {
		return nil, err
	}

	//扣除对应的币
	c := false
	if payload.IsToken {
		for _,v := range userinfo.Token {
			if v.Name == payload.Name {
				if v.Amount < payload.Amount {
					return nil,g.ErrAmountLow
				}
				v.Amount -= payload.Amount
				c = true
				break
			}
		}
	} else {
		for _,v := range userinfo.Currency {
			if v.Name == payload.Name {
				if v.Amount < payload.Amount {
					return nil,g.ErrAmountLow
				}
				v.Amount -= payload.Amount
				c = true
				break
			}
		}
	}

	//扣除手续费
	if c {
		for _,v := range userinfo.Currency {
			if v.Name == currencyBty {
				if v.Amount < payload.Fee {
					return nil,g.ErrCurrencyNotEnough
				}
				v.Amount -= payload.Fee
				break
			}
		}
		platform,err := a.getPlatformInfo()
		if err != nil {
			return nil,err
		}
		platform.FeeAmount += payload.Fee
		kv = append(kv,&types.KeyValue{Key:mallKeyPlatform(),Value:types.Encode(platform)})
	}

	value := types.Encode(userinfo)
	kv = append(kv, &types.KeyValue{Key: mallKeyUser(addr), Value: value})
	a.saveStateDB(kv)
	//logs = append(logs,&types.ReceiptLog{Ty:g.TyLogMallUserWithdraw,Log:types.Encode(userinfo)})
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//用户赠送
func (a *Action) UserGive(payload *g.MallUserGive) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	//检查用户是否已经存在
	fromaddr := a.addr
	sendUser,err := a.getUserInfo(fromaddr)
	if err != nil {
		return nil, err
	}

	toUser,err := a.getUserInfo(address.PubKeyToAddress(payload.To).String())
	if err != nil {
		return nil, err
	}

	if payload.IsToken {
		c := 0
		for _,v := range sendUser.Token {
			if v.Name == payload.Name {
				if v.Amount < payload.Amount {
					return nil,g.ErrAmountLow
				}
				v.Amount -= payload.Amount
				c++
				break
			}
		}
		if c == 0 {
			return nil,g.ErrTokenNotExist
		}
		for _,v := range toUser.Token {
			if v.Name == payload.Name {
				v.Amount += payload.Amount
				c++
				break
			}
		}
		if c == 1 {
			toUser.Token = append(toUser.Token,&g.MallUserToken{Name:payload.Name,Amount:payload.Amount})
		}

	} else {
		c := 0
		for _,v := range sendUser.Currency {
			if v.Name == payload.Name {
				if v.Amount < payload.Amount {
					return nil,g.ErrAmountLow
				}
				v.Amount -= payload.Amount
				c++
				break
			}
		}
		if c == 0 {
			return nil,g.ErrTokenNotExist
		}
		for _,v := range toUser.Currency {
			if v.Name == payload.Name {
				v.Amount += payload.Amount
				c++
				break
			}
		}
		if c == 1 {
			toUser.Currency = append(toUser.Currency,&g.MallUserCurrency{Name:payload.Name,Amount:payload.Amount})
		}
	}

	kv = append(kv, &types.KeyValue{Key: mallKeyUser(fromaddr), Value: types.Encode(sendUser)},&types.KeyValue{Key:mallKeyUser(address.PubKeyToAddress(payload.To).String()),Value:types.Encode(toUser)})
	//logs = append(logs,&types.ReceiptLog{Ty:g.TyLogMallUserGive,Log:nil})
	a.saveStateDB(kv)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//添加商品
func (a *Action) AddGood(payload *g.MallAddGood) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	_,err := a.getGoodInfo(payload.BaseInfo.GoodId)
	if err != types.ErrNotFound {
		return nil,g.ErrGoodExists
	}
	for _, v := range payload.SpecInfo.TokenInfo {
		//全新商品 更新数据
		if payload.BaseInfo.GoodType == g.MallAllType_GOOD_NEW {
			kv1,err := a.updateUserToken(v.Num,v.Token)
			if err != nil {
				return nil,err
			}
			a.saveStateDB(kv1)
			kv = append(kv,kv1...)
		}
	}
	kv = append(kv,&types.KeyValue{Key:mallKeyGoodInfo(payload.BaseInfo.GoodId),Value:types.Encode(payload)})
	a.saveStateDB(kv)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//支付
func (a *Action) Pay(payload *g.MallPay) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	for _,v := range payload.SinglePay {

		kv1,err := a.saveUserAssets(v)//v.PayOrder.SellerPubkey,v.PayOrder.BuyerPubkey,v.PayInfo.IsToken,v.PayInfo.Name,v.PayInfo.Num,v.PayOrder.Name,v.PayOrder.Total
		if err != nil {
			return nil,err
		}
		a.saveStateDB(kv1)
		kv = append(kv,kv1...)
	}

	logs = append(logs,&types.ReceiptLog{Ty:g.TyLogMallPay,Log:nil})

	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}


//提货
func (a *Action) Delivery(payload *g.MallDelivery) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	addr := a.addr
	userinfo ,err := a.getUserInfo(addr)
	if err != nil {
		return nil,err
	}
	c := false
	for _,v := range userinfo.Token {
		if v.Name == payload.Coin {
			if v.Amount < payload.Num {
				return nil,g.ErrAmountLow
			}
			v.Amount -= payload.Num
			c = true
			break
		}
	}
	if !c {
		return nil,g.ErrTokenNotExist
	}
	kv = append(kv,&types.KeyValue{Key:mallKeyUser(addr),Value:types.Encode(userinfo)})
	logs = append(logs,&types.ReceiptLog{Ty:g.TyLogMallDelivery,Log:types.Encode(userinfo)})
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}


func (a *Action) updateUserToken(amount int64, name string) ([]*types.KeyValue,error) {
	var kv []*types.KeyValue

	info,err := a.getUserInfo(a.addr)
	if err != nil {
		return nil,err
	}
	c := false
	for _,v := range info.Token {
		if v.Name == name {
			v.Amount += amount
			c = true
			break
		}
	}
	if !c {
		info.Token = append(info.Token,&g.MallUserToken{Name:name,Amount:amount})
	}
	kv = append(kv,&types.KeyValue{Key:mallKeyUser(a.addr),Value:types.Encode(info)})

	return kv ,nil
}

func (a *Action) saveUserAssets(pay *g.MallSinglePay) ([]*types.KeyValue,error) {
	var kv []*types.KeyValue

	switch pay.PayInfo.PayType {
	case g.MallAllType_PAY_FUTURES_FULL://付全款
		kv1 ,err := a.updatePayInfo(pay)
		if err != nil {
			return nil,err
		}
		kv = append(kv,kv1...)
	case g.MallAllType_PAY_FUTURES_DEPOSIT://付定金
		kv1 ,err := a.updatePayInfo(pay)
		if err != nil {
			return nil,err
		}
		kv = append(kv,kv1...)
	case g.MallAllType_PAY_FUTURES_TAIL://付尾款
		kv1 ,err := a.updatePayInfo(pay)
		if err != nil {
			return nil,err
		}
		kv = append(kv,kv1...)
	case g.MallAllType_PAY_SPOT:
		kv1 ,err := a.updatePayInfo(pay)
		if err != nil {
			return nil,err
		}
		kv = append(kv,kv1...)
	default :
		return nil,g.ErrWrongType
	}


	return kv,nil
}

func (a *Action) updatePayInfo(pay *g.MallSinglePay) ([]*types.KeyValue,error) {
	var kv []*types.KeyValue
	selleraddr := address.PubKeyToAddress(pay.PayOrder.SellerPubkey).String()
	buyueraddr := address.PubKeyToAddress(pay.PayOrder.BuyerPubkey).String()
	payeraddr := a.addr
	if selleraddr == payeraddr {
		return nil,g.ErrSellerNotSamePayer
	}
	sellerinfo, err := a.getUserInfo(selleraddr)
	if err != nil {
		return nil, err
	}
	buyuerinfo, err := a.getUserInfo(buyueraddr)
	if err != nil {
		return nil, err
	}
	payerinfo, err := a.getUserInfo(payeraddr)
	if err != nil {
		return nil, err
	}
	n := 0
	if pay.PayInfo.IsToken {
		//付的是token类型  卖家+ 付款家-
		for _,v := range sellerinfo.Token {
			if v.Name == pay.PayInfo.PayCoin {
				v.Amount += pay.PayInfo.PayNum
				n++
				break
			}
		}
		if n == 0 {
			sellerinfo.Token = append(sellerinfo.Token,&g.MallUserToken{Name:pay.PayInfo.PayCoin,Amount:pay.PayInfo.PayNum})
			n++
		}
		for _,v := range payerinfo.Token {
			if v.Name == pay.PayInfo.PayCoin {
				v.Amount -= pay.PayInfo.PayNum
				n++
				break
			}
		}
	} else {
		for _,v := range sellerinfo.Currency {
			if v.Name == pay.PayInfo.PayCoin {
				v.Amount += pay.PayInfo.PayNum
				n++
				break
			}
		}
		if n == 0 {
			sellerinfo.Currency = append(sellerinfo.Currency,&g.MallUserCurrency{Name:pay.PayInfo.PayCoin,Amount:pay.PayInfo.PayNum})
			n++
		}
		for _,v := range payerinfo.Currency {
			if v.Name == pay.PayInfo.PayCoin {
				v.Amount -= pay.PayInfo.PayNum
				n++
				break
			}
		}
	}
	if n != 2 {
		return nil,g.ErrTokenNotExist
	}
	//付定金不交换商品token
	if pay.PayInfo.PayType != g.MallAllType_PAY_FUTURES_DEPOSIT {
		m := 0
		for _,v := range sellerinfo.Token {
			if v.Name == pay.PayInfo.BuyCoin {
				if v.Amount <  pay.PayInfo.BuyNum {
					return nil,g.ErrAmountLow
				}
				v.Amount -= pay.PayInfo.BuyNum
				m++
				break
			}
		}
		if m == 0 {
			return nil,g.ErrTokenNotExist
		}
		for  _,v := range buyuerinfo.Token {
			if v.Name == pay.PayInfo.BuyCoin {
				v.Amount += pay.PayInfo.BuyNum
				m++
				break
			}
		}
		if m == 1 {
			buyuerinfo.Token = append(buyuerinfo.Token,&g.MallUserToken{Name:pay.PayInfo.BuyCoin,Amount:pay.PayInfo.BuyNum})
		}
		kv = append(kv,&types.KeyValue{Key:mallKeyUser(buyueraddr),Value:types.Encode(buyuerinfo)})
	}

	kv = append(kv,&types.KeyValue{Key:mallKeyUser(selleraddr),Value:types.Encode(sellerinfo)})
	kv = append(kv,&types.KeyValue{Key:mallKeyUser(payeraddr),Value:types.Encode(payerinfo)})

	return kv,nil
}


func (a *Action) getUserInfo(id string) (receipt *g.MallUserInfo,err error) {
	var userinfo g.MallUserInfo
	rval,err := a.db.Get(mallKeyUser(id))
	if err != nil {
		return nil,err
	}
	err = types.Decode(rval,&userinfo)
	if err != nil {
		return nil,err
	}
	return &userinfo,nil
}

func (a *Action) getPlatformInfo() (receipt *g.MallPlatformInfo,err error) {
	var info g.MallPlatformInfo
	rval,err := a.db.Get(mallKeyPlatform())
	if err != nil {
		if err != types.ErrNotFound {
			return nil,err
		}
		return &info,nil
	}
	err = types.Decode(rval,&info)
	if err != nil {
		return nil,err
	}
	return &info,nil
}

func (a *Action) getGoodInfo(id string) (receipt *g.MallAddGood,err error) {
	var goodInfo g.MallAddGood
	value ,err := a.db.Get(mallKeyGoodInfo(id))
	if err != nil {
		return nil,err
	}
	err = types.Decode(value,&goodInfo)
	if err != nil {
		return nil,err
	}
	return &goodInfo ,err
}

func (a *Action) saveStateDB(kv []*types.KeyValue) {
	for i:=0; i<len(kv) ; i++  {
		a.db.Set(kv[i].Key,kv[i].Value)
	}
}