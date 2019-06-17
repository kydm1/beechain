package executor

import (
	"github.com/33cn/chain33/account"
	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/types"
	"github.com/33cn/chain33/common/address"
	"encoding/hex"
	g "github.com/33cn/plugin/plugin/dapp/goldbill/types"
	"github.com/33cn/chain33/system/dapp"
)

type Action struct {
	coinsAccount *account.DB
	db           dbm.KV
	txhash       []byte
	fromaddr     string
	blocktime    int64
	height       int64
	execaddr     string
	index        int
}

func NewAction(t *Goldbill, tx *types.Transaction,index int) *Action {
	hash := tx.Hash()
	fromaddr := tx.From()
	return &Action{t.GetCoinsAccount(), t.GetStateDB(), hash, fromaddr,
		t.GetBlockTime(), t.GetHeight(),dapp.ExecAddress(string(tx.Execer)),index}
}

type calcbuy struct {
	amount int64
	list  []string
}

func (a *Action) InitPlatform(req *g.GoldbillInitPlatform) (*types.Receipt ,error) {
	ok := checkIsSupper(a.fromaddr)
	if  !ok {
		return nil,g.ErrNoPrivilege
	}
	v,err := a.db.Get(calcGoldbillPlatformKey())
	if err != types.ErrNotFound && err != nil{
		return nil,err
	}
	if v != nil && err == nil{
		return nil,g.ErrPlatformExists
	}
	var platform g.GoldbillPlatform
	platform.Info = req.GetInfo()
	platform.Pubkey = req.GetPlatformKey()
	value := types.Encode(&platform)
	a.db.Set(calcGoldbillPlatformKey(),value)
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var respLog g.GoldbillResponseInitPlatform
	respLog.GoldbillPlatform = &platform
	logs = append(logs,&types.ReceiptLog{Ty:g.TyLogGoldbillInitPlatform,Log:types.Encode(&respLog)})
	kv = append(kv,&types.KeyValue{Key:calcGoldbillPlatformKey(),Value:value})
	return &types.Receipt{Ty:types.ExecOk,KV:kv,Logs:logs},nil
}

func (a *Action) RegisterUser(req *g.GoldbillRegisterUser) (*types.Receipt ,error) {
	clog.Error("RegisterUser")
	_,err := a.db.Get(calcGoldbillUserKey(a.fromaddr))
	if err != nil && err != types.ErrNotFound {
		return nil,err
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var respLog g.GoldbillResponseRegisterUser
	var u g.GoldbillUser
	u.Uid = req.GetUserId()
	u.Type = req.GetUserType()
	u.Pubkey = req.GetUserPubkey()

	uvalue := types.Encode(&u)
	addr := ""

	if req.UserType == g.GoldbillUserType_UT_ADMIN {
		ok := checkIsSupper(a.fromaddr)
		if !ok {
			return nil,g.ErrNoPrivilege
		}
		addr = address.PubKeyToAddress(req.GetUserPubkey()).String()

	} else {
		addr = a.fromaddr
	}

	kv = append(kv,&types.KeyValue{Key:calcGoldbillUserKey(addr),Value:uvalue})
	a.db.Set(calcGoldbillUserKey(addr),uvalue)
	respLog.GoldbillUser = &u
	logs = append(logs,&types.ReceiptLog{Ty:g.TyLogGoldbillRegisterUser,Log:types.Encode(&respLog)})
	return &types.Receipt{Ty:types.ExecOk,KV:kv,Logs:logs},nil
}

func (a *Action) SetFee(req *g.GoldbillSetFee) (*types.Receipt ,error) {
	value ,err := a.db.Get(calcGoldbillPlatformKey())
	if err != nil {
		return nil,g.ErrNoPrivilege
	}
	_,err = a.db.Get(calcGoldbillAdminKey(a.fromaddr))
	if err != nil {
		if err == types.ErrNotFound {
			return nil,g.ErrUserNotExists
		}
		return nil,g.ErrNoPrivilege
	}
	var platform g.GoldbillPlatform
	types.Decode(value,&platform)
	platform.Fee = req.GetFee()
	pvalue := types.Encode(&platform)
	a.db.Set(calcGoldbillPlatformKey(),pvalue)
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var respLog g.GoldbillResponseSetFee
	respLog.GoldbillPlatform = &platform
	kv = append(kv,&types.KeyValue{Key:calcGoldbillPlatformKey(),Value:pvalue})
	logs = append(logs,&types.ReceiptLog{Ty:g.TyLogGoldbillSetFee,Log:types.Encode(&respLog)})
	return &types.Receipt{Ty:types.ExecOk,KV:kv,Logs:logs},nil
}

func (a *Action) SetStamps(req *g.GoldbillSetStamps) (*types.Receipt ,error) {
	value ,err := a.db.Get(calcGoldbillPlatformKey())
	if err != nil {
		return nil,g.ErrNoPrivilege
	}
	_,err = a.db.Get(calcGoldbillAdminKey(a.fromaddr))
	if err != nil {
		if err == types.ErrNotFound {
			return nil,g.ErrUserNotExists
		}
		return nil,g.ErrNoPrivilege
	}
	var platform g.GoldbillPlatform
	types.Decode(value,&platform)
	platform.Stamps = req.GetStamps()
	pvalue := types.Encode(&platform)
	a.db.Set(calcGoldbillPlatformKey(),pvalue)
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var respLog g.GoldbillResponseSetStamps
	respLog.GoldbillPlatform = &platform
	kv = append(kv,&types.KeyValue{Key:calcGoldbillPlatformKey(),Value:pvalue})
	logs = append(logs,&types.ReceiptLog{Ty:g.TyLogGoldbillSetStamps,Log:types.Encode(&respLog)})
	return &types.Receipt{Ty:types.ExecOk,KV:kv,Logs:logs},nil
}

func (a *Action) SetBail(req *g.GoldbillSetBail) (*types.Receipt ,error) {
	value ,err := a.db.Get(calcGoldbillPlatformKey())
	if err != nil {
		return nil,g.ErrNoPrivilege
	}
	_,err = a.db.Get(calcGoldbillAdminKey(a.fromaddr))
	if err != nil {
		if err == types.ErrNotFound {
			return nil,g.ErrUserNotExists
		}
		return nil,g.ErrNoPrivilege
	}
	var platform g.GoldbillPlatform
	types.Decode(value,&platform)
	platform.Bail = req.GetBail()
	pvalue := types.Encode(&platform)
	a.db.Set(calcGoldbillPlatformKey(),pvalue)
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var respLog g.GoldbillResponseSetBail
	respLog.GoldbillPlatform = &platform
	kv = append(kv,&types.KeyValue{Key:calcGoldbillPlatformKey(),Value:pvalue})
	logs = append(logs,&types.ReceiptLog{Ty:g.TyLogGoldbillSetBail,Log:types.Encode(&respLog)})
	return &types.Receipt{Ty:types.ExecOk,KV:kv,Logs:logs},nil
}

func (a *Action) Withdraw(req *g.GoldbillWithdraw) (*types.Receipt ,error) {
	value ,err := a.db.Get(calcGoldbillPlatformKey())
	if err != nil {
		return nil,g.ErrNoPrivilege
	}
	avalue,err := a.db.Get(calcGoldbillUserKey(a.fromaddr))
	if err != nil {
		if err == types.ErrNotFound {
			return nil,g.ErrUserNotExists
		}
		return nil,g.ErrNoPrivilege
	}
	var platform g.GoldbillPlatform
	types.Decode(value,&platform)
	var admin g.GoldbillUser
	types.Decode(avalue,&admin)
	if admin.Type != g.GoldbillUserType_UT_ADMIN {
		return nil,g.ErrNoPrivilege
	}
	toaddr := address.PubKeyToAddress(req.GetPubkey()).String()
	uvalue ,err := a.db.Get(calcGoldbillUserKey(toaddr))
	if err != nil {
		return nil,err
	}
	var touser g.GoldbillUser
	types.Decode(uvalue,&touser)
	if touser.Uid != req.Uid {
		return nil,g.ErrNoPrivilege
	}
	if touser.Rmb < req.GetAmount() {
		clog.Error("user rmb not enough")
		return nil,g.ErrRMBNotEnough
	}
	touser.Rmb -= req.GetAmount()
	platform.Kcoin -= req.GetAmount()
	psave := types.Encode(&platform)
	usave := types.Encode(&touser)
	a.db.Set(calcGoldbillPlatformKey(),psave)
	a.db.Set(calcGoldbillUserKey(toaddr),usave)
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var respLog g.GoldbillResponseWithdraw
	respLog.GoldbillPlatform = &platform
	respLog.GoldbillUser = &touser
	kv = append(kv,&types.KeyValue{Key:calcGoldbillPlatformKey(),Value:psave})
	kv = append(kv,&types.KeyValue{Key:calcGoldbillUserKey(toaddr),Value:usave})
	logs = append(logs,&types.ReceiptLog{Ty:g.TyLogGoldbillWithdraw,Log:types.Encode(&respLog)})
	return &types.Receipt{Ty:types.ExecOk,KV:kv,Logs:logs},nil
}

func (a *Action) Deposit(req *g.GoldbillDeposit) (*types.Receipt ,error) {
	value ,err := a.db.Get(calcGoldbillPlatformKey())
	if err != nil {
		return nil,g.ErrNoPrivilege
	}
	uvalue,err := a.db.Get(calcGoldbillUserKey(a.fromaddr))
	if err != nil {
		if err == types.ErrNotFound {
			return nil,g.ErrUserNotExists
		}
		return nil,g.ErrNoPrivilege
	}
	var platform g.GoldbillPlatform
	types.Decode(value,&platform)
	var user g.GoldbillUser
	types.Decode(uvalue,&user)
	//if user.Rmb < req.GetDeposit().GetAmount() {
	//	clog.Error("user rmb not enough")
	//	return nil,g.ErrRMBNotEnough
	//}
	user.Rmb += req.GetAmount()
	platform.Kcoin += req.GetAmount()
	psave := types.Encode(&platform)
	usave := types.Encode(&user)
	a.db.Set(calcGoldbillPlatformKey(),psave)
	a.db.Set(calcGoldbillUserKey(a.fromaddr),usave)
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var respLog g.GoldbillResponseDeposit
	respLog.GoldbillPlatform = &platform
	respLog.GoldbillUser = &user
	kv = append(kv,&types.KeyValue{Key:calcGoldbillPlatformKey(),Value:psave})
	kv = append(kv,&types.KeyValue{Key:calcGoldbillUserKey(a.fromaddr),Value:usave})
	logs = append(logs,&types.ReceiptLog{Ty:g.TyLogGoldbillDeposit,Log:types.Encode(&respLog)})

	return &types.Receipt{Ty:types.ExecOk,KV:kv,Logs:logs},nil
}

func (a *Action) WithdrawKcoin(req *g.GoldbillWithdrawKcoin) (*types.Receipt ,error) {
	value ,err := a.db.Get(calcGoldbillPlatformKey())
	if err != nil {
		return nil,g.ErrNoPrivilege
	}
	uvalue,err := a.db.Get(calcGoldbillUserKey(a.fromaddr))
	if err != nil {
		if err == types.ErrNotFound {
			return nil,g.ErrUserNotExists
		}
		return nil,g.ErrNoPrivilege
	}
	var platform g.GoldbillPlatform
	types.Decode(value,&platform)
	var user g.GoldbillUser
	types.Decode(uvalue,&user)
	if user.Kcoin < req.GetAmount() {
		clog.Error("user rmb not enough")
		return nil,g.ErrCoinNotEnough
	}
	user.Rmb += req.GetAmount()
	user.Kcoin -= req.GetAmount()
	platform.Kcoin += req.GetAmount()
	platform.Rmb -= req.GetAmount()
	psave := types.Encode(&platform)
	usave := types.Encode(&user)
	a.db.Set(calcGoldbillPlatformKey(),psave)
	a.db.Set(calcGoldbillUserKey(a.fromaddr),usave)
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var respLog g.GoldbillResponseWithdrawKcoin
	respLog.GoldbillPlatform = &platform
	respLog.GoldbillUser = &user
	kv = append(kv,&types.KeyValue{Key:calcGoldbillPlatformKey(),Value:psave})
	kv = append(kv,&types.KeyValue{Key:calcGoldbillUserKey(a.fromaddr),Value:usave})
	logs = append(logs,&types.ReceiptLog{Ty:g.TyLogGoldbillWithdrawKcoin,Log:types.Encode(&respLog)})
	return &types.Receipt{Ty:types.ExecOk,KV:kv,Logs:logs},nil
}

func (a *Action) DepositKcoin(req *g.GoldbillDepositKcoin) (*types.Receipt ,error) {
	value ,err := a.db.Get(calcGoldbillPlatformKey())
	if err != nil {
		return nil,g.ErrNoPrivilege
	}
	uvalue,err := a.db.Get(calcGoldbillUserKey(a.fromaddr))
	if err != nil {
		if err == types.ErrNotFound {
			return nil,g.ErrUserNotExists
		}
		return nil,g.ErrNoPrivilege
	}
	var platform g.GoldbillPlatform
	types.Decode(value,&platform)
	var user g.GoldbillUser
	types.Decode(uvalue,&user)
	if user.Rmb < req.GetAmount() {
		clog.Error("user rmb not enough")
		return nil,g.ErrRMBNotEnough
	}
	user.Rmb -= req.GetAmount()
	user.Kcoin += req.GetAmount()
	platform.Kcoin -= req.GetAmount()
	platform.Rmb += req.GetAmount()
	psave := types.Encode(&platform)
	usave := types.Encode(&user)
	a.db.Set(calcGoldbillPlatformKey(),psave)
	a.db.Set(calcGoldbillUserKey(a.fromaddr),usave)
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var respLog g.GoldbillResponseDepositKcoin
	respLog.GoldbillPlatform = &platform
	respLog.GoldbillUser = &user
	kv = append(kv,&types.KeyValue{Key:calcGoldbillPlatformKey(),Value:psave})
	kv = append(kv,&types.KeyValue{Key:calcGoldbillUserKey(a.fromaddr),Value:usave})
	logs = append(logs,&types.ReceiptLog{Ty:g.TyLogGoldbillDepositKcoin,Log:types.Encode(&respLog)})
	return &types.Receipt{Ty:types.ExecOk,KV:kv,Logs:logs},nil
}

func (a *Action) PayKcoin(req *g.GoldbillPayKCoin) (*types.Receipt ,error) {
	_ ,err := a.db.Get(calcGoldbillPlatformKey())
	if err != nil {
		return nil,g.ErrNoPrivilege
	}
	uvalue,err := a.db.Get(calcGoldbillUserKey(a.fromaddr))
	if err != nil {
		if err == types.ErrNotFound {
			return nil,g.ErrUserNotExists
		}
		return nil,g.ErrNoPrivilege
	}
	toaddr := address.PubKeyToAddress(req.GetTopubkey()).String()
	tovalue,err := a.db.Get(calcGoldbillUserKey(toaddr))

	var user g.GoldbillUser
	types.Decode(uvalue,&user)
	var to g.GoldbillUser
	types.Decode(tovalue,&to)
	if user.Kcoin < req.GetAmount() {
		clog.Error("user coin not enough")
		return nil,g.ErrCoinNotEnough
	}
	user.Kcoin -= req.GetAmount()
	to.Kcoin += req.GetAmount()
	tsave := types.Encode(&to)
	usave := types.Encode(&user)
	a.db.Set(calcGoldbillUserKey(toaddr),tsave)
	a.db.Set(calcGoldbillUserKey(a.fromaddr),usave)
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var respLog g.GoldbillResponsePayKCoin
	respLog.From = &user
	respLog.To = &to
	kv = append(kv,&types.KeyValue{Key:calcGoldbillUserKey(toaddr),Value:tsave})
	kv = append(kv,&types.KeyValue{Key:calcGoldbillUserKey(a.fromaddr),Value:usave})
	logs = append(logs,&types.ReceiptLog{Ty:g.TyLogGoldbillPayKCoin,Log:types.Encode(&respLog)})

	return &types.Receipt{Ty:types.ExecOk,KV:kv,Logs:logs},nil
}

func (a *Action) Invoice(req *g.GoldbillInvoice) (*types.Receipt ,error) {

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	return &types.Receipt{Ty:types.ExecOk,KV:kv,Logs:logs},nil
}

func (a *Action) AddBill(req *g.GoldbillAddBill) (*types.Receipt ,error) {
	_,err := a.db.Get(calcGoldbillPlatformKey())
	if err != nil {
		return nil,g.ErrNoPrivilege
	}
	uvalue,err := a.db.Get(calcGoldbillUserKey(a.fromaddr))
	if err != nil {
		if err == types.ErrNotFound {
			return nil,g.ErrUserNotExists
		}
		return nil,g.ErrNoPrivilege
	}
	var user g.GoldbillUser
	types.Decode(uvalue,&user)
	var kvs []*types.KeyValue
	var respLog g.GoldbillResponseAddBill

	if user.Type == g.GoldbillUserType_UT_ADMIN {
		keylist := make([]string,0)
		list := make(map[string][]string,0)
		mapList := make(map[string]*g.GoldbillDetail,0)
		for _, v := range req.GetBillList() {
			addr := address.PubKeyToAddress(v.Pubkey).String()
			_,ok := list[addr]
			if !ok {
				list[addr] = []string{}
			}
			keylist = appendKeyList(keylist,addr)
			check := false
			for _,vv := range list[addr] {
				if vv == v.BillId {
					check = true
					break
				}
			}
			if check {
				return nil,g.ErrDupBillId
			}
			mapList[v.BillId] = v
			respLog.BillList = append(respLog.BillList,v.BillId)
			list[addr] = append(list[addr],v.BillId)
		}
		clog.Info("ORDER LIST:",list)
		for _,v := range keylist {
			uuvalue,err := a.db.Get(calcGoldbillUserKey(v))
			if err != nil {
				if err == types.ErrNotFound {
					return nil,g.ErrUserNotExists
				}
				return nil,g.ErrNoPrivilege
			}
			var uuser g.GoldbillUser
			types.Decode(uuvalue,&uuser)
			for _,vv := range list[v] {
				bill := &g.GoldbillDetail{}
				bill = mapList[vv]
				bsave := types.Encode(bill)
				kvs = append(kvs,&types.KeyValue{Key:calcGoldbillDetailKey(vv),Value:bsave})
			}
			uuser.WrList = append(user.WrList,list[v]...)
			uusave := types.Encode(&uuser)
			kvs = append(kvs,&types.KeyValue{Key:calcGoldbillUserKey(v),Value:uusave})
		}

	} else if user.Type == g.GoldbillUserType_UT_USER {
		list := make(map[string]string,0)
		add := make([]string,0)

		for _, v:= range req.GetBillList() {
			_,	ok := list[v.BillId]
			if ok {
				return nil,g.ErrDupBillId
			}
			_,err := a.db.Get(calcGoldbillDetailKey(v.BillId))
			if err != nil {
				if err == types.ErrNotFound {
					return nil,g.ErrBillNotExists
				}
				return nil,g.ErrNoPrivilege
			}
			bill := &g.GoldbillDetail{}
			bill = v
			bsave := types.Encode(bill)
			list[v.BillId] = v.BillId
			add = append(add,v.BillId)
			respLog.BillList = append(respLog.BillList,v.BillId)
			kvs = append(kvs,&types.KeyValue{Key:calcGoldbillDetailKey(v.BillId),Value:bsave})
		}
		user.WrList = append(user.WrList,add...)
		usave := types.Encode(&user)
		kvs = append(kvs,&types.KeyValue{Key:calcGoldbillUserKey(a.fromaddr),Value:usave})
	}
	for i:=0 ;i < len(kvs);i++ {
		a.db.Set(kvs[i].Key,kvs[i].Value)
	}
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	kv = append(kv,kvs...)
	logs = append(logs,&types.ReceiptLog{Ty:g.TyLogGoldbillAddBill,Log:types.Encode(&respLog)})
	clog.Info("logs kv",kv)
	clog.Info("logs logs",logs)
	return &types.Receipt{Ty:types.ExecOk,KV:kv,Logs:logs},nil
}

func (a *Action) BatchSell(req *g.GoldbillBatchSell) (*types.Receipt ,error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var respLog g.GoldbillResponseBatchSell

	_,err := a.db.Get(calcGoldbillPlatformKey())
	if err != nil {
		return nil,g.ErrNoPrivilege
	}
	_,err = a.db.Get(calcGoldbillUserKey(a.fromaddr))
	if err != nil {
		if err == types.ErrNotFound {
			return nil,g.ErrUserNotExists
		}
		return nil,g.ErrNoPrivilege
	}
	list := make(map[string]string,0)
	var kvs []*types.KeyValue
	for _, v:= range req.GetBillList() {
		_,	ok := list[v.BillId]
		if ok {
			return nil,g.ErrDupBillId
		}
		bvalue,err := a.db.Get(calcGoldbillDetailKey(v.BillId))
		if err != nil {
			if err == types.ErrNotFound {
				return nil,g.ErrBillNotExists
			}
			return nil,g.ErrNoPrivilege
		}
		var bill g.GoldbillDetail
		types.Decode(bvalue,&bill)
		if bill.BillState == g.GoldbillState_BS_SELL || bill.BillState == g.GoldbillState_BS_TAKEOUT {
			return nil,g.ErrWrongState
		}
		bill.BillState = g.GoldbillState_BS_SELL
		bill.Amount = v.Amount
		bsave := types.Encode(&bill)
		respLog.SellList = append(respLog.SellList,v.BillId)
		kvs = append(kvs,&types.KeyValue{Key:calcGoldbillDetailKey(v.BillId),Value:bsave})
	}
	//var platform g.GoldbillPlatform
	//types.Decode(value,&platform)
	//var user g.GoldbillUser
	//types.Decode(uvalue,&user)

	kv = append(kv,kvs...)
	logs = append(logs,&types.ReceiptLog{Ty:g.TyLogGoldbillBatchSell,Log:types.Encode(&respLog)})
	clog.Info("logs kv",kv)
	clog.Info("logs logs",logs)
	return &types.Receipt{Ty:types.ExecOk,KV:kv,Logs:logs},nil
}

func (a *Action) BatchCancelSell(req *g.GoldbillBatchCancelSell) (*types.Receipt ,error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var respLog g.GoldbillResponseBatchCancelSell
	_,err := a.db.Get(calcGoldbillPlatformKey())
	if err != nil {
		return nil,g.ErrNoPrivilege
	}
	_,err = a.db.Get(calcGoldbillUserKey(a.fromaddr))
	if err != nil {
		if err == types.ErrNotFound {
			return nil,g.ErrUserNotExists
		}
		return nil,g.ErrNoPrivilege
	}
	list := make(map[string]string,0)
	var kvs []*types.KeyValue
	for _, v:= range req.GetBillList() {
		_,	ok := list[v]
		if ok {
			return nil,g.ErrDupBillId
		}
		bvalue,err := a.db.Get(calcGoldbillDetailKey(v))
		if err != nil {
			if err == types.ErrNotFound {
				return nil,g.ErrBillNotExists
			}
			return nil,g.ErrNoPrivilege
		}
		var bill g.GoldbillDetail
		types.Decode(bvalue,&bill)
		if bill.BillState != g.GoldbillState_BS_SELL {
			return nil,g.ErrWrongState
		}
		bill.BillState = g.GoldbillState_BS_FREE
		bill.Amount = 0
		bsave := types.Encode(&bill)
		kvs = append(kvs,&types.KeyValue{Key:calcGoldbillDetailKey(v),Value:bsave})
	}
	//var platform g.GoldbillPlatform
	//types.Decode(value,&platform)
	//var user g.GoldbillUser
	//types.Decode(uvalue,&user)
	respLog.CancelList = req.BillList
	kv = append(kv,kvs...)
	logs = append(logs,&types.ReceiptLog{Ty:g.TyLogGoldbillBatchCancelSell,Log:types.Encode(&respLog)})

	return &types.Receipt{Ty:types.ExecOk,KV:kv,Logs:logs},nil
}

func (a *Action) BatchBuy(req *g.GoldbillBatchBuy) (*types.Receipt ,error) {
	_,err := a.db.Get(calcGoldbillPlatformKey())
	if err != nil {
		return nil,g.ErrNoPrivilege
	}
	value,err := a.db.Get(calcGoldbillUserKey(a.fromaddr))
	if err != nil {
		if err == types.ErrNotFound {
			return nil,g.ErrUserNotExists
		}
		return nil,g.ErrNoPrivilege
	}
	var user g.GoldbillUser
	types.Decode(value,&user)
	var amount int64
	var kvs []*types.KeyValue
	maplist := make(map[string]*calcbuy,0)
	keylist := make([]string,0)
	for _, v := range req.GetBillList() {
		bvalue,err := a.db.Get(calcGoldbillDetailKey(v))
		if err != nil {
			if err == types.ErrNotFound {
				return nil,g.ErrBillNotExists
			}
			return nil,g.ErrNoPrivilege
		}
		var bill g.GoldbillDetail
		types.Decode(bvalue,&bill)
		addr := address.PubKeyToAddress(bill.Pubkey).String()
		keylist = appendKeyList(keylist,addr)
		_,ok := maplist[addr]
		if !ok {
			maplist[addr] = &calcbuy{amount:0,list:[]string{}}
		}
		maplist[addr].list = append(maplist[addr].list,bill.BillId)
		maplist[addr].amount += bill.Amount
		if bill.BillState != g.GoldbillState_BS_SELL {
			return nil,g.ErrWrongState
		}
		amount += bill.Amount
		bill.BillState = g.GoldbillState_BS_FREE
		bill.Amount = 0
		bill.Oid = user.Uid
		bill.Pubkey = user.Pubkey
		kvs = append(kvs,&types.KeyValue{Key:calcGoldbillDetailKey(v),Value:types.Encode(&bill)})
	}
	if user.Rmb < amount {
		return nil,g.ErrRMBNotEnough
	}
	for _, v := range keylist {
		nvalue ,err := a.db.Get(calcGoldbillUserKey(v))
		if err != nil {
			return nil,err
		}
		var seller g.GoldbillUser
		types.Decode(nvalue,&seller)
		seller.Rmb += maplist[v].amount
		for _,vv := range maplist[v].list {
			for kkk,vvv := range seller.WrList {
				if vv == vvv {
					seller.WrList = append(seller.WrList[:kkk], seller.WrList[kkk+1:]...)
					break
				}
			}
		}

		kvs = append(kvs,&types.KeyValue{Key:calcGoldbillUserKey(v),Value:types.Encode(&seller)})
	}
	user.Rmb -= amount
	user.WrList = append(user.WrList,req.GetBillList()...)
	kvs = append(kvs,&types.KeyValue{Key:calcGoldbillUserKey(a.fromaddr),Value:types.Encode(&user)})
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var respLog g.GoldbillResponseBatchBuy
	respLog.BuyList = req.BillList
	kv = append(kv,kvs...)
	logs = append(logs,&types.ReceiptLog{Ty:g.TyLogGoldbillBatchBuy,Log:types.Encode(&respLog)})
	clog.Info("logs kv",kv)
	clog.Info("logs logs",logs)
	return &types.Receipt{Ty:types.ExecOk,KV:kv,Logs:logs},nil
}

func (a *Action) BatchPay(req *g.GoldbillBatchPay) (*types.Receipt ,error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var respLog g.GoldbillResponseBatchPay

	_,err := a.db.Get(calcGoldbillPlatformKey())
	if err != nil {
		return nil,g.ErrNoPrivilege
	}
	value,err := a.db.Get(calcGoldbillUserKey(a.fromaddr))
	if err != nil {
		if err == types.ErrNotFound {
			return nil,g.ErrAdminNotExists
		}
		return nil,g.ErrNoPrivilege
	}
	var user g.GoldbillUser
	types.Decode(value,&user)
	var kvs []*types.KeyValue
	maplist := make(map[string][]string,0)
	keylist := make([]string,0)
	rest := make([]string,0)
	for _,v := range req.GetPayTo() {
		bcheck := false
		for _,vv := range user.WrList {
			if v.BillId == vv {
				bcheck = true
				break
			}
		}
		if !bcheck {
			return nil,g.ErrNoPrivilege
		}
		toaddr := address.PubKeyToAddress(v.Toaddr).String()
		clog.Info("pubkey:",hex.EncodeToString(v.Toaddr),"addr:",toaddr)
		_,ok := maplist[toaddr]
		if !ok {
			maplist[toaddr] = []string{}
		}

		bvalue ,err := a.db.Get(calcGoldbillDetailKey(v.BillId))
		if err != nil {
			return nil,types.ErrNotFound
		}
		var bill g.GoldbillDetail
		types.Decode(bvalue,&bill)
		bill.Pubkey = v.Toaddr
		bill.Oid = v.Toid
		keylist = appendKeyList(keylist,toaddr)
		maplist[toaddr] = append(maplist[toaddr],v.BillId)
		respLog.BillList = append(respLog.BillList,v.BillId)
		kvs = append(kvs,&types.KeyValue{Key:calcGoldbillDetailKey(v.BillId),Value:types.Encode(&bill)})
	}
	clog.Info("key",keylist)
	for _,v :=range user.WrList {
		check := false
		for _,vv := range req.GetPayTo() {
			if v == vv.BillId {
				check = true
				break
			}
		}
		if !check {
			rest = append(rest,v)
		}
	}
	user.WrList = rest
	kvs = append(kvs,&types.KeyValue{Key:calcGoldbillUserKey(a.fromaddr),Value:types.Encode(&user)})
	for _,v := range keylist {
		tovalue,err := a.db.Get(calcGoldbillUserKey(v))
		if err != nil {
			if err == types.ErrNotFound {
				clog.Error("hello : user not exists",v)
				return nil,g.ErrUserNotExists
			}
			return nil,g.ErrNoPrivilege
		}
		var touser g.GoldbillUser
		types.Decode(tovalue,&touser)
		touser.WrList = append(touser.WrList,maplist[v]...)
		kvs = append(kvs,&types.KeyValue{Key:calcGoldbillUserKey(v),Value:types.Encode(&touser)})
	}

	for i:=0 ;i < len(kvs);i++ {
		a.db.Set(kvs[i].Key,kvs[i].Value)
	}
	respLog.From = &user
	kv = append(kv,kvs...)
	logs = append(logs,&types.ReceiptLog{Ty:g.TyLogGoldbillBatchPay,Log:types.Encode(&respLog)})
	clog.Info("logs kv",kv)
	clog.Info("logs logs",logs)
	return &types.Receipt{Ty:types.ExecOk,KV:kv,Logs:logs},nil
}

func (a *Action) BatchTakeOut(req *g.GoldbillBatchTakeOut) (*types.Receipt ,error) {
	_,err := a.db.Get(calcGoldbillPlatformKey())
	if err != nil {
		return nil,g.ErrNoPrivilege
	}
	_,err = a.db.Get(calcGoldbillUserKey(a.fromaddr))
	if err != nil {
		if err == types.ErrNotFound {
			return nil,g.ErrUserNotExists
		}
		return nil,g.ErrNoPrivilege
	}
	var kvs []*types.KeyValue
	for _ ,v := range req.GetBillId() {
		bvalue ,err := a.db.Get(calcGoldbillDetailKey(v))
		if err != nil {
			return nil,g.ErrBillNotExists
		}
		var bill g.GoldbillDetail
		types.Decode(bvalue,&bill)
		if address.PubKeyToAddress(bill.Pubkey).String() != a.fromaddr {
			return nil,g.ErrNoPrivilege
		}
		bill.BillState = g.GoldbillState_BS_TAKEOUT
		bsave := types.Encode(&bill)
		kvs = append(kvs,&types.KeyValue{Key:calcGoldbillDetailKey(v),Value:bsave})
	}
	for i:=0 ;i < len(kvs);i++ {
		a.db.Set(kvs[i].Key,kvs[i].Value)
	}
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var respLog g.GoldbillResponseBatchTakeOut
	respLog.BillId = req.BillId
	kv = append(kv,kvs...)
	logs = append(logs,&types.ReceiptLog{Ty:g.TyLogGoldbillBatchTakeOut,Log:types.Encode(&respLog)})
	return &types.Receipt{Ty:types.ExecOk,KV:kv,Logs:logs},nil
}

func (a *Action) SliceBill(req *g.GoldbillSliceBill) (*types.Receipt ,error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var respLog g.GoldbillResponseSliceBill
	var newList []*g.GoldbillDetail

	value,err := a.db.Get(calcGoldbillUserKey(a.fromaddr))
	if err != nil {
		return nil,err
	}
	var kvs []*types.KeyValue
	var admin g.GoldbillUser
	types.Decode(value,&admin)
	if admin.Type != g.GoldbillUserType_UT_ADMIN {
		return nil,g.ErrNoPrivilege
	}
	bvalue,err := a.db.Get(calcGoldbillDetailKey(req.GetBigId()))
	if err != nil {
		return nil,err
	}
	var bigbill g.GoldbillDetail
	types.Decode(bvalue,&bigbill)
	addr := address.PubKeyToAddress(bigbill.Pubkey).String()
	uvalue ,err := a.db.Get(calcGoldbillUserKey(addr))
	if err != nil {
		return nil,err
	}
	var user g.GoldbillUser
	types.Decode(uvalue,&user)
	for k,v := range user.WrList {
		if v == req.BigId {
			tmp := user.WrList[:]
			user.WrList = tmp[:k]
			user.WrList = append(user.WrList,tmp[k+1:]...)
		}
	}
	for _,v := range req.GetSmallBill() {
		smallbill := &g.GoldbillDetail{}
		smallbill = &bigbill
		smallbill.BillNum = v.BillNum
		newList = append(newList,smallbill)
		user.WrList = append(user.WrList,v.BillId)
		kvs = append(kvs,&types.KeyValue{Key:calcGoldbillDetailKey(v.BillId),Value:types.Encode(smallbill)})
	}
	bigbill = g.GoldbillDetail{}
	kvs = append(kvs,&types.KeyValue{Key:calcGoldbillUserKey(addr),Value:types.Encode(&user)})
	kvs = append(kvs,&types.KeyValue{Key:calcGoldbillDetailKey(req.BigId),Value:types.Encode(&bigbill)})
	for i:=0 ;i < len(kvs);i++ {
		a.db.Set(kvs[i].Key,kvs[i].Value)
	}

	respLog.NewBill = newList
	kv = append(kv,kvs...)
	logs = append(logs,&types.ReceiptLog{Ty:g.TyLogGoldbillSliceBill,Log:types.Encode(&respLog)})
	return &types.Receipt{Ty:types.ExecOk,KV:kv,Logs:logs},nil
}

func (a *Action) BatchRegUser(req *g.GoldbillBatchRegUser) (*types.Receipt ,error) {
	_,err := a.db.Get(calcGoldbillAdminKey(a.fromaddr))
	if err != nil && err != types.ErrNotFound {
		return nil,err
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var kvs []*types.KeyValue
	var respLog g.GoldbillResponseBatchRegUser
	for _,v := range req.GetUserList() {
		addr := address.PubKeyToAddress(v.GetPubkey()).String()
		_,err := a.db.Get(calcGoldbillUserKey(addr))
		if err != types.ErrNotFound {
			return nil,g.ErrUserExists
		}
		var u g.GoldbillUser
		u.Uid = v.GetUid()
		u.Type = g.GoldbillUserType_UT_USER
		u.Pubkey = v.GetPubkey()
		clog.Info("uid:",v.GetUid(),"-addr:",addr,"-PUBKEY",hex.EncodeToString(v.Pubkey))
		uvalue := types.Encode(&u)
		kvs = append(kvs,&types.KeyValue{Key:calcGoldbillUserKey(addr),Value:uvalue})
		respLog.Succ = append(respLog.Succ,v.GetUid())
	}
	for i:=0 ;i < len(kvs);i++ {
		a.db.Set(kvs[i].Key, kvs[i].Value)
	}
	kv = append(kv,kvs...)
	logs = append(logs,&types.ReceiptLog{Ty:g.TyLogGoldbillBatchRegUser,Log:types.Encode(&respLog)})
	return &types.Receipt{Ty:types.ExecOk,KV:kv,Logs:logs},nil
}


func (a *Action) checkAccountExists(addr string) bool {
	acc := a.coinsAccount.AccountKey(addr)
	_ ,err := a.db.Get(acc)
	if err != nil {
		return false
	}
	return true
}

func checkIsSupper(addr string) bool {
	for _, m := range conf.GStrList("superManager") {
		if addr == m {
			return true
		}
	}
	return false
}

func appendKeyList(list []string,add string) []string {
	ch := false
	for _,vv := range list {
		if vv ==  add {
			ch = true
			break
		}
	}
	if !ch {
		list = append(list,add)
	}
	return list
}
