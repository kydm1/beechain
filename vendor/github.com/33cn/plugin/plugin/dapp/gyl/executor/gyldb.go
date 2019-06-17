package executor

import (
	"github.com/33cn/chain33/types"
	"bytes"
	"encoding/hex"
	"fmt"
	"errors"
	dbm "github.com/33cn/chain33/common/db"
	g "github.com/33cn/plugin/plugin/dapp/gyl/types"
)


type Action struct {
	db     dbm.KV
	txhash []byte
	height int64
	index  int
}

func NewAction(t *Gyl, tx *types.Transaction,index int) *Action {
	hash := tx.Hash()
	return &Action{t.GetStateDB(), hash, t.GetHeight(),index}
}

//单据保存操作，单据的转让、签收都用这个方法
//下面所有使用到的State和Opty都在proto文件中定义
func (a *Action) SaveReceipt(payload *g.ZsgjSaveReceipt, pubkey []byte) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	clog.Info("GYL", "ENTERPRISE_NAME", payload.ZsgjReceipt.OpCompany)
	clog.Info("GYL", "ENTERPRISE_NAME", payload.ZsgjReceipt.ReceiveCompany)
	pub, err := a.db.Get(gylKeyUser(payload.ZsgjReceipt.OpCompany.Name))
	if err != nil {
		return nil, errors.New("data not found")
	}
	clog.Error("GYL", "LOCAL_KEY", hex.EncodeToString(pub))
	clog.Error("GYL", "INPUT_KEY", hex.EncodeToString(pubkey))
	if bytes.Compare(pub, pubkey) != 0 {
		return nil, errors.New("different pubkey")
	}
	if payload.ZsgjReceipt.StartDate > payload.ZsgjReceipt.EndDate {
		return nil, g.ErrDateConflict
	}
	if payload.ZsgjReceipt.SumAmount < 0 {
		return nil, g.ErrAmountLow
	}

	if payload.ZsgjReceipt.ReceiptViceId != "" && payload.ZsgjReceipt.ReceiptViceId != payload.ZsgjReceipt.ReceiptId {
		oldRid := gylKeyReceipt(payload.ZsgjReceipt.ReceiptViceId)
		receiveCompanyKey := gylKeyUser(payload.ZsgjReceipt.ReceiveCompany.Name)
		coreCompanyKey := gylKeyUser(payload.ZsgjReceipt.CoreCompany.Name)
		issueComKey := gylKeyUser(payload.ZsgjReceipt.IssuedAgency.Name)
		receivePubkey, err := a.db.Get(receiveCompanyKey)
		if err != nil {
			return nil,err
		}
		corePubkey, err := a.db.Get(coreCompanyKey)
		if err != nil {
			return nil,err
		}
		issuePubkey, err := a.db.Get(issueComKey)
		if err != nil {
			return nil,err
		}
		issueCompanyRoleKey := append(gylKeyUserInfo(issuePubkey), []byte(payload.ZsgjReceipt.IssuedAgency.RoleId)...)
		receiveCompanyRoleKey := append(gylKeyUserInfo(receivePubkey), []byte(payload.ZsgjReceipt.ReceiveCompany.RoleId)...)
		zsgjReceiveRoleAccount := &g.ZsgjRoleAccount{}
		roleval ,err := a.db.Get(receiveCompanyRoleKey)
		if err != nil {
			if err != types.ErrNotFound {
				return nil,err
			}
			zsgjReceiveRoleAccount.Enterprise = payload.ZsgjReceipt.ReceiveCompany
		} else if err == nil && roleval != nil {
			types.Decode(roleval,zsgjReceiveRoleAccount)
		}
		zsgjIssueRoleAccount := &g.ZsgjRoleAccount{}
		issueval ,err := a.db.Get(issueCompanyRoleKey)
		if err != nil {
			if err != types.ErrNotFound {
				return nil,err
			}
			zsgjIssueRoleAccount.Enterprise = payload.ZsgjReceipt.IssuedAgency
		} else if err == nil && issueval != nil {
			err = types.Decode(issueval,zsgjReceiveRoleAccount)
			if err != nil {
				return nil,err
			}
		}

		oldReceiptValue, err := a.db.Get(oldRid)
		if err != nil {
			return nil, err
		}
		var val g.ZsgjReceiptInfo
		err = types.Decode(oldReceiptValue, &val)
		if err != nil {
			return nil, errors.New("receipt decode err")
		}

		switch payload.ZsgjReceipt.State {
		case g.ZsgjState_BILL_HOLDING:// //单据签收
			//上游企业单据资产增加
			if zsgjReceiveRoleAccount.ReceiptAssets == nil {
				receiptAssets := &g.ReceiptAssets{}
				receiptAssets.Available += payload.ZsgjReceipt.SumAmount
				receiptAssets.Receipt += payload.ZsgjReceipt.SumAmount
				zsgjReceiveRoleAccount.ReceiptAssets = receiptAssets
			} else {
				zsgjReceiveRoleAccount.ReceiptAssets.Available += payload.ZsgjReceipt.SumAmount
				zsgjReceiveRoleAccount.ReceiptAssets.Receipt += payload.ZsgjReceipt.SumAmount
			}

			kv1 := &types.KeyValue{Key: receiveCompanyRoleKey, Value: types.Encode(zsgjReceiveRoleAccount)}
			kv = append(kv, kv1)
			//核心企业未兑付增加
			coreComVal, err := a.db.Get(gylKeyUserInfo(corePubkey))
			if err != nil {
				return nil,err
			}
			zsgjCoreAccount := &g.ZsgjAccount{}
			types.Decode(coreComVal, zsgjCoreAccount)
			if zsgjCoreAccount.Rmb == nil {
				rmbassets := &g.RmbAssets{}
				rmbassets.WaitCash += payload.ZsgjReceipt.SumAmount
				zsgjCoreAccount.Rmb = rmbassets
			} else {
				zsgjCoreAccount.Rmb.WaitCash += payload.ZsgjReceipt.SumAmount
			}

			kv2 := &types.KeyValue{Key: gylKeyUserInfo(corePubkey), Value: types.Encode(zsgjCoreAccount)}
			kv = append(kv, kv2)
		case g.ZsgjState_WAIT_CORE_CHECKED: //待核心企业确认
			val.State = 57 //单据转让中
		case g.ZsgjState_BILL_HOLDING_ISSUE: //发行机构持有中
			val.State = 58 //转让完成
			//上游企业单据资产减少
			if payload.ZsgjReceipt.Opty == g.OperationType_AGREE_TRANSFER_SIGN {
				if zsgjReceiveRoleAccount.ReceiptAssets == nil {
					receiptAssets := &g.ReceiptAssets{}
					receiptAssets.Available -= payload.ZsgjReceipt.SumAmount
					receiptAssets.Receipt -= payload.ZsgjReceipt.SumAmount
					zsgjReceiveRoleAccount.ReceiptAssets = receiptAssets
				} else {
					zsgjReceiveRoleAccount.ReceiptAssets.Available -= payload.ZsgjReceipt.SumAmount
					zsgjReceiveRoleAccount.ReceiptAssets.Receipt -= payload.ZsgjReceipt.SumAmount
				}
				kv1 := &types.KeyValue{Key: receiveCompanyRoleKey, Value: types.Encode(zsgjReceiveRoleAccount)}
				kv = append(kv, kv1)
				//发行机构单据资产增加

				if zsgjIssueRoleAccount.ReceiptAssets == nil {
					receiptAssets := &g.ReceiptAssets{}
					receiptAssets.Available += payload.ZsgjReceipt.SumAmount
					receiptAssets.Receipt += payload.ZsgjReceipt.SumAmount
					zsgjIssueRoleAccount.ReceiptAssets = receiptAssets
				} else {
					zsgjIssueRoleAccount.ReceiptAssets.Available += payload.ZsgjReceipt.SumAmount
					zsgjIssueRoleAccount.ReceiptAssets.Receipt += payload.ZsgjReceipt.SumAmount
				}
				kv2 := &types.KeyValue{Key: issueCompanyRoleKey, Value: types.Encode(zsgjIssueRoleAccount)}
				kv = append(kv, kv2)
			} else if payload.ZsgjReceipt.Opty == g.OperationType_ISSUER_REVOKE {
				//TODO在产品撤销操作中改变产品资产//融资中产品资产减少

				if zsgjIssueRoleAccount.ReceiptAssets == nil {
					receiptAssets := &g.ReceiptAssets{}
					receiptAssets.Available += payload.ZsgjReceipt.SumAmount
					receiptAssets.Receipt += payload.ZsgjReceipt.SumAmount
					zsgjIssueRoleAccount.ReceiptAssets = receiptAssets
				} else {
					zsgjIssueRoleAccount.ReceiptAssets.Available += payload.ZsgjReceipt.SumAmount
					zsgjIssueRoleAccount.ReceiptAssets.Receipt += payload.ZsgjReceipt.SumAmount
				}
				kv3 := &types.KeyValue{Key: issueCompanyRoleKey, Value: types.Encode(zsgjIssueRoleAccount)}
				kv = append(kv, kv3)
			}
		case g.ZsgjState_ISSUE_BE_PRODUCT://发行产品
			//发行机构单据资产减少、产品资产增加

			if zsgjIssueRoleAccount.ReceiptAssets == nil {
				receiptAssets := &g.ReceiptAssets{}
				receiptAssets.Available -= payload.ZsgjReceipt.SumAmount
				receiptAssets.Receipt -= payload.ZsgjReceipt.SumAmount
				zsgjIssueRoleAccount.ReceiptAssets = receiptAssets
			} else {
				zsgjIssueRoleAccount.ReceiptAssets.Available -= payload.ZsgjReceipt.SumAmount
				zsgjIssueRoleAccount.ReceiptAssets.Receipt -= payload.ZsgjReceipt.SumAmount
			}
			if zsgjIssueRoleAccount.ProductAssets == nil {
				productAssets := &g.ProductAssets{}
				productAssets.Available += payload.ZsgjReceipt.SumAmount
				productAssets.Product += payload.ZsgjReceipt.SumAmount
				zsgjIssueRoleAccount.ProductAssets = productAssets
			} else {
				zsgjIssueRoleAccount.ProductAssets.Available += payload.ZsgjReceipt.SumAmount
				zsgjIssueRoleAccount.ProductAssets.Product += payload.ZsgjReceipt.SumAmount
			}
			kv1 := &types.KeyValue{Key: issueCompanyRoleKey, Value: types.Encode(zsgjIssueRoleAccount)}

			kv = append(kv, kv1)
			//核心机构未兑付减少
			coreComVal, err := a.db.Get(gylKeyUserInfo(corePubkey))
			if err != nil {
				return nil,err
			}
			zsgjCoreAccount := &g.ZsgjAccount{}
			types.Decode(coreComVal, zsgjCoreAccount)
			if zsgjCoreAccount.Rmb == nil {
				rmbassets := &g.RmbAssets{}
				rmbassets.WaitCash -= payload.ZsgjReceipt.SumAmount
				zsgjCoreAccount.Rmb = rmbassets
			} else {
				zsgjCoreAccount.Rmb.WaitCash -= payload.ZsgjReceipt.SumAmount
			}
			kv2 := &types.KeyValue{Key: gylKeyUserInfo(corePubkey), Value: types.Encode(zsgjCoreAccount)}
			kv = append(kv, kv2)
		default:
		}
		if payload.ZsgjReceipt.Opty == 15 { //上游企业撤销转让
			val.State = 3 //单据被上游企业持有中
		}
		newValue := types.Encode(&val)
		kv = append(kv, &types.KeyValue{Key: oldRid, Value: newValue}) //更新母单状态
	}

	receiptInfo := NewReceiptInfo(payload)
	id := gylKeyReceipt(payload.ZsgjReceipt.ReceiptId)
	value := types.Encode(&receiptInfo.zri)
	kv = append(kv, &types.KeyValue{Key: id, Value: value})


	a.saveStateDB(kv)

	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//产品保存的操作，产品的发行、确认等
func (a *Action) SaveProduct(save *g.ZsgjSaveProduct, pubkey []byte) (*types.Receipt, error) {
	clog.Info("GYL", "ENTERPRISE_NAME", save.ZsgjProduct.OpCompany)
	opKey := gylKeyUser(save.ZsgjProduct.OpCompany.Name)
	pub, err := a.db.Get(opKey)
	if err != nil {
		return nil, errors.New("data not found")
	}
	clog.Error("GYL", "LOCAL_KEY", hex.EncodeToString(pub))
	clog.Error("GYL", "INPUT_KEY", hex.EncodeToString(pubkey))
	if bytes.Compare(pub, pubkey) != 0 {
		return nil, errors.New("different pubkey")
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	oldRid := gylKeyReceipt(save.ZsgjProduct.ReceiptId)
	oldReceiptValue, err := a.db.Get(oldRid)
	if err != nil {
		return nil, err
	}
	var val g.ZsgjReceiptInfo
	var productValue g.ZsgjProductInfo
	err = types.Decode(oldReceiptValue, &val)
	if err != nil {
		return nil, errors.New("receipt decode err")
	}

	//操作时产品状态判断
	pid := gylKeyProduct(save.ZsgjProduct.ProductId)
	rid :=gylKeyReceipt(save.ZsgjProduct.ReceiptId)
	if save.ZsgjProduct.State == g.ZsgjState_WAIT_GUARANTEE && save.ZsgjProduct.Opty == g.OperationType_TRANSFER_TO_PRODUCT {//待担保  操作转化产品
		_, err := a.db.Get(pid)
		if err == nil {
			return nil, errors.New("product id exist")
		}
		_, err = a.db.Get(rid)
		if err != nil {
			return nil, errors.New("receipt id not exist")
		}
		product := NewProductInfo(save)
		value := types.Encode(&product.zpi)
		kv = append(kv, &types.KeyValue{Key: pid, Value: value})
		a.saveStateDB(kv)
		return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
	} else if save.ZsgjProduct.State == g.ZsgjState_WAIT_GUARANTEE && save.ZsgjProduct.Opty == g.OperationType_UPDATE_PRODUCT {//待担保  编辑产品
		_, err := a.db.Get(rid)
		if err != nil {
			return nil, errors.New("receipt id not exist")
		}
		product := NewProductInfo(save)
		value := types.Encode(&product.zpi)
		kv = append(kv, &types.KeyValue{Key: pid, Value: value})
		a.saveStateDB(kv)
		return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
	} else if save.ZsgjProduct.State == g.ZsgjState_ON_LISTING {//挂牌产品
		opKey, err := a.db.Get(opKey)
		if err != nil {
			return nil,err
		}
		zsgjRoleAccount := &g.ZsgjRoleAccount{}
		roleKey := append(gylKeyUserInfo(opKey), []byte(save.ZsgjProduct.OpCompany.RoleId)...)
		roleAccountValue, err := a.db.Get(roleKey)

		if err != nil {
			if err != types.ErrNotFound {
				return nil,err
			}
			zsgjRoleAccount.Enterprise = save.ZsgjProduct.OpCompany
		}
		if roleAccountValue != nil {
			err := types.Decode(roleAccountValue, zsgjRoleAccount)
			if err != nil {
				return nil,err
			}
		}


		if zsgjRoleAccount.ProductAssets == nil {
			productAssets := &g.ProductAssets{}
			productAssets.Available -= save.ZsgjProduct.IssueScale
			productAssets.Financing += save.ZsgjProduct.IssueScale
			zsgjRoleAccount.ProductAssets = productAssets
		} else {
			zsgjRoleAccount.ProductAssets.Available -= save.ZsgjProduct.IssueScale
			zsgjRoleAccount.ProductAssets.Financing += save.ZsgjProduct.IssueScale
		}

		kv1 := &types.KeyValue{Key: roleKey, Value: types.Encode(zsgjRoleAccount)}
		kv = append(kv, kv1)
	} else if save.ZsgjProduct.State == g.ZsgjState_PRODUCT_CANCE { //产品撤牌
		opKey, err := a.db.Get(opKey)
		if err != nil {
			return nil,err
		}
		zsgjRoleAccount := &g.ZsgjRoleAccount{}
		roleKey := append(gylKeyUserInfo(opKey), []byte(save.ZsgjProduct.OpCompany.RoleId)...)
		roleAccountValue, err := a.db.Get(roleKey)

		if err != nil {
			if err != types.ErrNotFound {
				return nil,err
			}
			zsgjRoleAccount.Enterprise = save.ZsgjProduct.OpCompany
		}
		if roleAccountValue != nil {
			err := types.Decode(roleAccountValue, zsgjRoleAccount)
			if err != nil {
				return nil,errors.New("zsgjRoleAccount deocde err")
			}
		}

		if zsgjRoleAccount.ProductAssets == nil {
			productAssets := &g.ProductAssets{}
			productAssets.Available += save.ZsgjProduct.IssueScale
			productAssets.Financing -= save.ZsgjProduct.IssueScale
			zsgjRoleAccount.ProductAssets = productAssets
		} else {
			zsgjRoleAccount.ProductAssets.Available += save.ZsgjProduct.IssueScale
			zsgjRoleAccount.ProductAssets.Financing -= save.ZsgjProduct.IssueScale
		}
		kv1 := &types.KeyValue{Key: roleKey, Value: types.Encode(zsgjRoleAccount)}

		kv = append(kv, kv1)
	}


	oldProductValue, err := a.db.Get(pid)
	if err != nil {
		return nil, err
	}
	err = types.Decode(oldProductValue, &productValue)
	if err != nil {
		return nil, errors.New("product decode err")
	}

	product := NewProductInfo(save)
	value := types.Encode(&product.zpi)
	kv = append(kv, &types.KeyValue{Key: pid, Value: value})
	a.saveStateDB(kv)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil


}

//开单操作
func (a *Action) AssetRegister(asset *g.AssetRegisterAction, pubkey []byte) (*types.Receipt, error) {
	if asset.SumAmount < 0 {
		return nil, g.ErrAmountLow
	}
	if asset.StartDate > asset.EndDate {
		return nil, g.ErrDateConflict
	}
	//bo, _ := isNumberExist(asset.ContractNo, db)
	//if bo {
	//	return nil, errors.New("number is exist")
	//}
	clog.Info("GYL", "ASSET_DETAILS", asset)
	pub, err := a.db.Get(gylKeyUser(asset.PayCompany.Name))
	if err != nil {
		return nil, errors.New("data not found")
	}
	clog.Error("GYL", "LOCAL_KEY", hex.EncodeToString(pub))
	clog.Error("GYL", "INPUT_KEY", hex.EncodeToString(pubkey))
	if bytes.Compare(pub, pubkey) != 0 {
		return nil, errors.New("different pubkey")
	}
	if asset.State != 1 {
		return nil, errors.New("state not match")
	}
	rid := gylKeyReceipt(asset.ReceiptId)
	_, err = a.db.Get(rid)
	if err == nil {
		return nil, errors.New("receiptId exist")
	}


	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	//TODO完整对象的赋值操作
	receipt := &g.ZsgjReceiptInfo{}
	receipt.ReceiptId = asset.ReceiptId
	receipt.CoreCompany = asset.PayCompany
	receipt.ReceiveCompany = asset.ReceiveCompany
	receipt.InvestAgency = asset.InvestAgency
	receipt.SumAmount = asset.SumAmount
	receipt.Opty = asset.Opty
	receipt.State = asset.State
	//初始单据标示
	receipt.ReceiptViceId = asset.ReceiptId
	value := types.Encode(receipt)
	kv = append(kv, &types.KeyValue{Key: rid, Value: value})
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//充值请求
func (a *Action) ApplyRecharge(payload *g.ZsgjApplyRecharge, pubkey []byte) (*types.Receipt, error) {
	if payload.Amount < 0 {
		return nil, g.ErrAmountLow
	}
	clog.Info("GYL", "ENTERPRISE_NAME", payload.CompanyName)
	pub, err := a.db.Get(gylKeyUser(payload.CompanyName))
	if err != nil {
		return nil, err
	}
	clog.Error("GYL", "LOCAL_KEY", hex.EncodeToString(pub))
	clog.Error("GYL", "INPUT_KEY", hex.EncodeToString(pubkey))
	if bytes.Compare(pub, pubkey) != 0 {
		return nil, errors.New("different pubkey")
	}
	keyid := gylKeyDeposit(payload.Id)
	_, err = a.db.Get(keyid)
	if err == nil {
		return nil, errors.New("this id is exist")
	}
	opAccountKey := gylKeyUserInfo(pubkey)
	opAccountVal ,err := a.db.Get(opAccountKey)
	if err != nil {
		clog.Error("GYL","APPLYRECHARGE",err)
		return nil,err
	}
	opAccount := &g.ZsgjAccount{}
	err = types.Decode(opAccountVal,opAccount)
	if err != nil {
		clog.Error("GYL","opAccount DECODE ERR",err)
		return nil,errors.New("opAccount DECODE ERR")
	}
	if opAccount.Rmb == nil {
		rmbassets := &g.RmbAssets{}
		rmbassets.Rmb += rmbassets.Rmb + payload.Amount
		opAccount.Rmb = rmbassets
	} else {
		opAccount.Rmb.Rmb += payload.Amount
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	kv = append(kv, &types.KeyValue{Key: keyid, Value: types.Encode(payload)},&types.KeyValue{Key:opAccountKey,Value:types.Encode(opAccount)})
	a.saveStateDB(kv)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//提款请求
func (a *Action) ApplyWithdraw(payload *g.ZsgjApplyWithdraw, pubkey []byte) (*types.Receipt, error) {
	if payload.State != 27 {
		return nil, errors.New("非出账成功状态")
	}
	if payload.Amount < 0 {
		return nil, g.ErrAmountLow
	}
	clog.Info("GYL", "ENTERPRISE_NAME", payload.CompanyName)
	pub, err := a.db.Get(gylKeyUser(payload.CompanyName))
	if err != nil {
		return nil, err
	}
	clog.Error("GYL", "LOCAL_KEY", hex.EncodeToString(pub))
	clog.Error("GYL", "INPUT_KEY", hex.EncodeToString(pubkey))
	if bytes.Compare(pub, pubkey) != 0 {
		return nil, errors.New("different pubkey")
	}
	keyid := gylKeyWithdraw(payload.Id)
	_, err = a.db.Get(keyid)
	if err == nil {
		return nil, errors.New("this id is exist")
	}
	val, err := a.db.Get(gylKeyUserInfo(pubkey))
	if err != nil {
		return nil, err
	}
	zsgjAccount := &g.ZsgjAccount{}
	err = types.Decode(val, zsgjAccount)
	if err != nil {
		return nil, errors.New("account decode err")
	}
	fmt.Println(zsgjAccount)
	if zsgjAccount.Rmb.Rmb < payload.Amount {
		return nil, g.ErrRNE
	}

	if zsgjAccount.Rmb == nil {
		rmbassets := &g.RmbAssets{}
		rmbassets.Rmb -= payload.Amount
		zsgjAccount.Rmb = rmbassets
	} else {
		zsgjAccount.Rmb.Rmb -= payload.Amount
	}
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	kv = append(kv, &types.KeyValue{Key: keyid, Value: types.Encode(payload)},&types.KeyValue{Key:gylKeyUserInfo(pubkey),Value:types.Encode(zsgjAccount)})
	a.saveStateDB(kv)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//产品的兑付操作
func (a *Action) Cash(payload *g.ZsgjCash, pubkey []byte) (*types.Receipt, error) {
	clog.Info("GYL", "ENTERPRISE_NAME", payload.OpCompany.Name)

	pub, err := a.db.Get(gylKeyUser(payload.OpCompany.Name))
	if err != nil {
		return nil, err
	}
	clog.Error("GYL", "LOCAL_KEY", hex.EncodeToString(pub))
	clog.Error("GYL", "INPUT_KEY", hex.EncodeToString(pubkey))
	if bytes.Compare(pub, pubkey) != 0 {
		return nil, errors.New("different pubkey")
	}
	investKey := gylKeyUser(payload.InvestAgency.Name)
	investRoleKey := append(investKey,[]byte(payload.InvestAgency.RoleId)...)
	investPubkey ,err := a.db.Get(investKey)
	if err != nil {
		return nil,err
	}

	zsgjInvestAccount := &g.ZsgjAccount{}
	zsgjInvestRoleAccount := &g.ZsgjRoleAccount{}
	investVal ,err := a.db.Get(gylKeyUserInfo(investPubkey))
	if err != nil {
		return nil,err
	}
	investRoleVal ,err := a.db.Get(investRoleKey)
	if err != nil {
		if err != types.ErrNotFound {
			return nil,err
		}
		zsgjInvestRoleAccount.Enterprise = payload.InvestAgency
	}
	if investRoleVal != nil {
		err = types.Decode(investRoleVal,zsgjInvestRoleAccount)
		if err != nil {
			return nil,errors.New("zsgjInvestRoleAccount decode err")
		}
	}

	err = types.Decode(investVal,zsgjInvestAccount)
	if err != nil {
		return nil,errors.New("zsgjInvestAccount decode err")
	}
	if zsgjInvestAccount.Rmb == nil {
		rmbassets := &g.RmbAssets{}
		rmbassets.Rmb += payload.WaitSumAmount                     //wait cash amount
		rmbassets.Income += (payload.WaitSumAmount - payload.WaitAmount) //累积收益增加
		zsgjInvestAccount.Rmb = rmbassets
	} else {
		zsgjInvestAccount.Rmb.Rmb += payload.WaitSumAmount                     //wait cash amount
		zsgjInvestAccount.Rmb.Income += (payload.WaitSumAmount - payload.WaitAmount) //累积收益增加
	}
	if zsgjInvestRoleAccount.ProductAssets == nil {
		productAssets := &g.ProductAssets{}
		productAssets.Product -= payload.WaitAmount //产品资产减少待付本金
		zsgjInvestRoleAccount.ProductAssets = productAssets
	} else {
		zsgjInvestRoleAccount.ProductAssets.Product -= payload.WaitAmount
	}
	pid := gylKeyProduct(payload.Id)
	_, err = a.db.Get(pid)
	if err != nil {
		return nil, errors.New("产品不存在")
	}
	_, err = a.db.Get(gylKeyUser(payload.InvestAgency.Name))
	if err != nil {
		return nil, errors.New("invest aengcy not found")
	}
	//发行机构资产变化
	issueval, err := a.db.Get(gylKeyUserInfo(pubkey))
	if err != nil {
		return nil, err
	}
	zsgjIssueAccount := &g.ZsgjAccount{}
	err = types.Decode(issueval, zsgjIssueAccount)
	if err != nil {
		return nil, errors.New("account decode err")
	}

	if zsgjIssueAccount.Rmb == nil {
		rmbassets := &g.RmbAssets{}
		rmbassets.Rmb -= payload.WaitSumAmount
		rmbassets.Cashed += payload.WaitSumAmount
		rmbassets.WaitCash -= payload.WaitSumAmount
		zsgjIssueAccount.Rmb = rmbassets
	} else {
		if zsgjIssueAccount.Rmb.Rmb < payload.WaitSumAmount {
			return nil, g.ErrRNE
		}
		zsgjIssueAccount.Rmb.Rmb -= payload.WaitSumAmount
		zsgjIssueAccount.Rmb.Cashed += payload.WaitSumAmount
		zsgjIssueAccount.Rmb.WaitCash -= payload.WaitSumAmount
	}
	var productValue g.ZsgjProductInfo
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	oldId := gylKeyProduct(payload.Id)
	oldProductValue, err := a.db.Get(oldId)
	if err != nil {
		return nil, err
	}
	err = types.Decode(oldProductValue, &productValue)
	if err != nil {
		return nil, errors.New("product decode err")
	}
	if productValue.State != 23 {
		return nil, errors.New("状态错误")
	}
	productValue.Cashed += payload.WaitAmount
	//已经兑付的金额等于发行规模
	if productValue.Cashed == productValue.IssueScale {
		productValue.State = 25
	}

	kv = append(kv, &types.KeyValue{Key: oldId, Value: types.Encode(&productValue)},
	&types.KeyValue{Key:gylKeyUserInfo(pubkey),Value:types.Encode(zsgjIssueAccount)},
	&types.KeyValue{Key:investRoleKey,Value:types.Encode(zsgjInvestRoleAccount)},
	&types.KeyValue{Key:gylKeyUserInfo(investPubkey),Value:types.Encode(zsgjInvestAccount)})
	a.saveStateDB(kv)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//产品清算操作
func (a *Action) Clear(payload *g.ZsgjClear, pubkey []byte) (*types.Receipt, error) {
	pub, err := a.db.Get(gylKeyUser(payload.IssuedName.Name))
	if err != nil {
		return nil, err
	}
	clog.Error("GYL", "LOCAL_KEY", hex.EncodeToString(pub))
	clog.Error("GYL", "INPUT_KEY", hex.EncodeToString(pubkey))
	if bytes.Compare(pub, pubkey) != 0 {
		return nil, errors.New("different pubkey")
	}
	oldId := gylKeyProduct(payload.ProductId)
	_, err = a.db.Get(oldId)
	if err != nil {
		return nil, errors.New("产品不存在")
	}

	var productValue g.ZsgjProductInfo
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	oldProductValue, err := a.db.Get(oldId)
	if err != nil {
		return nil, err
	}
	err = types.Decode(oldProductValue, &productValue)
	if err != nil {
		return nil, errors.New("product decode err")
	}
	if payload.State != 23 && payload.State != 24 {
		return nil, errors.New("状态错误")
	} else if payload.State == g.ZsgjState_SUCCESS_RELEASE {//发行成功待兑付

		issuevalue, err := a.db.Get(gylKeyUserInfo(pub))
		if err != nil {
			return nil,err
		}
		zsgjAccount := &g.ZsgjAccount{}
		types.Decode(issuevalue, zsgjAccount)
		if zsgjAccount.Rmb == nil {
			rmbassets := &g.RmbAssets{}
			rmbassets.Rmb += payload.Amount
			zsgjAccount.Rmb = rmbassets
		} else {
			zsgjAccount.Rmb.Rmb += payload.Amount
		}

		kv1 := &types.KeyValue{Key: gylKeyUserInfo(pub), Value: types.Encode(zsgjAccount)}
		kv = append(kv, kv1)
	} else if payload.State == g.ZsgjState_FAIL_RELEASE {//发行失败
		delistValue, err := a.db.Get(gylKeyDelist(payload.ProductId))
		if err != nil {
			return nil, err
		}
		var delist g.Delist
		types.Decode(delistValue, &delist)
		for _, v := range delist.Info {
			dpubKey, err := a.db.Get(gylKeyUser(v.Name)) //获得企业对应公钥
			if err != nil {
				return nil,err
			}
			dvalue, err  := a.db.Get(gylKeyUserInfo(dpubKey))
			if err != nil {
				return nil,err
			}
			dzsgjAccount := &g.ZsgjAccount{}
			types.Decode(dvalue, dzsgjAccount)
			// zsgjAccount.ProductAssets.Product -= v.Amount
			if dzsgjAccount.Rmb == nil {
				rmbassets := &g.RmbAssets{}
				rmbassets.Rmb += v.Amount
				rmbassets.Invest -= v.Amount //累积投资是否需要回滚，存疑
				dzsgjAccount.Rmb = rmbassets
			} else {
				dzsgjAccount.Rmb.Rmb += v.Amount
				dzsgjAccount.Rmb.Invest -= v.Amount //累积投资是否需要回滚，存疑
			}
			kv1 := &types.KeyValue{Key: gylKeyUserInfo(dpubKey), Value: types.Encode(dzsgjAccount)}

			kv = append(kv, kv1)

			dzsgjRoleAccount := &g.ZsgjRoleAccount{}
			roleKey := append(gylKeyUserInfo(dpubKey), []byte("6")...)
			roleValue, err := a.db.Get(roleKey)
			if err != nil {
				if err != types.ErrNotFound {
					return nil,err
				}
				enterprise := g.Enterprise{}
				enterprise.Name = v.Name
				enterprise.RoleId = "6"
				dzsgjRoleAccount.Enterprise = &enterprise
			}

			if roleValue != nil {
				err := types.Decode(roleValue, dzsgjRoleAccount)
				if err != nil {
					return nil,errors.New("dzsgjRoleAccount decode err")
				}
			}

			if dzsgjRoleAccount.ProductAssets == nil {
				productAssets := &g.ProductAssets{}
				productAssets.Product -= v.Amount
				dzsgjRoleAccount.ProductAssets = productAssets
			} else {
				dzsgjRoleAccount.ProductAssets.Product -= v.Amount
			}

			roleValue = types.Encode(dzsgjRoleAccount)
			kv2 := &types.KeyValue{Key: roleKey, Value: roleValue}
			kv = append(kv, kv2)
		}
	}
	productValue.State = payload.State
	value := types.Encode(&productValue)
	kv = append(kv, &types.KeyValue{Key: oldId, Value: value})
	a.saveStateDB(kv)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//产品摘牌操作
func (a *Action) Delist(delist *g.ZsgjDelist, pubkey []byte) (*types.Receipt, error) {
	var productValue g.ZsgjProductInfo
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	pidKey := gylKeyProduct(delist.ProductId)
	pvalue, err := a.db.Get(pidKey)
	if err != nil {
		return nil, errors.New("产品不存在")
	}
	err = types.Decode(pvalue, &productValue)
	if err != nil {
		return nil, errors.New("product decode err")
	}
	if delist.Purchase > delist.Purchaseable {
		return nil, g.ErrPurchaseTooMuch
	}
	clog.Info("GYL", "ENTERPRISE_NAME", delist.OpCompany)
	pub, err := a.db.Get(gylKeyUser(delist.OpCompany.Name))
	if err != nil {
		return nil, err
	}
	clog.Error("GYL", "LOCAL_KEY", hex.EncodeToString(pub))
	clog.Error("GYL", "INPUT_KEY", hex.EncodeToString(pubkey))
	if bytes.Compare(pub, pubkey) != 0 {
		return nil, errors.New("different pubkey")
	}
	val, err := a.db.Get(gylKeyUserInfo(pubkey))
	if err != nil {
		return nil, err
	}
	//投资机构变化
	zsgjAccount := &g.ZsgjAccount{}
	err = types.Decode(val, zsgjAccount)
	if err != nil {
		return nil, errors.New("account decode err")
	}
	clog.Info("zsgj", "delist", zsgjAccount)
	if zsgjAccount.Rmb.Rmb < delist.Purchase {
		return nil, g.ErrRNE
	}
	if zsgjAccount.Rmb == nil {
		rmbassets := &g.RmbAssets{}
		rmbassets.Rmb -= delist.Purchase
		rmbassets.Invest += delist.Purchase
		zsgjAccount.Rmb = rmbassets
	} else {
		zsgjAccount.Rmb.Rmb -= delist.Purchase
		zsgjAccount.Rmb.Invest += delist.Purchase
	}
	kv1 := &types.KeyValue{Key:gylKeyUserInfo(pubkey),Value:types.Encode(zsgjAccount)}
	zsgjInvestAccount := &g.ZsgjRoleAccount{}
	investRoleKey := append(gylKeyUserInfo(pubkey),[]byte("6")...)
	investRoleVale,err := a.db.Get(investRoleKey)
	if err != nil {
		if err != types.ErrNotFound {
			return nil,err
		}
		zsgjInvestAccount.Enterprise = delist.OpCompany
	}
	if investRoleVale != nil {
		err = types.Decode(investRoleVale,zsgjInvestAccount)
		if err != nil {
			return nil,errors.New("zsgjInvestAccount decode err")
		}
	}

	if zsgjInvestAccount.ProductAssets == nil {
		productAssets := &g.ProductAssets{}
		productAssets.Product += delist.Purchase
		zsgjInvestAccount.ProductAssets = productAssets
	} else {
		zsgjInvestAccount.ProductAssets.Product += delist.Purchase
	}
	kv2 := &types.KeyValue{Key:investRoleKey,Value:types.Encode(zsgjInvestAccount)}

	//发行机构变化
	if productValue.IssuedAgency == nil {
		clog.Error("execLocal","ZsgjDelistAction","nil err")
		return nil,errors.New("product issue nil")
	}
	zsgjIssueAccount := &g.ZsgjRoleAccount{}
	issueKey := gylKeyUser(productValue.IssuedAgency.Name)
	issuePubkey ,err := a.db.Get(issueKey)
	issueRoleKey := append(gylKeyUserInfo(issuePubkey),[]byte(productValue.IssuedAgency.RoleId)...)
	issueRoleVale,err := a.db.Get(issueRoleKey)
	if err != nil {
		if err != types.ErrNotFound {
			return nil,err
		}
		zsgjIssueAccount.Enterprise = productValue.IssuedAgency
	}
	if issueRoleVale != nil {
		err = types.Decode(issueRoleVale,zsgjIssueAccount)
		if err != nil {
			return nil,errors.New("zsgjIssueAccount decode err")
		}
	}

	if zsgjIssueAccount.ProductAssets == nil {
		productAssets := &g.ProductAssets{}
		productAssets.Product -= delist.Purchase
		productAssets.Financing -= delist.Purchase
		zsgjIssueAccount.ProductAssets = productAssets
	} else {
		zsgjIssueAccount.ProductAssets.Product -= delist.Purchase
		zsgjIssueAccount.ProductAssets.Financing -= delist.Purchase
	}
	kv3 := &types.KeyValue{Key:issueRoleKey,Value:types.Encode(zsgjIssueAccount)}

	//产品变化
	productValue.Delisted += delist.Purchase
	if productValue.Delisted == productValue.IssueScale {
		productValue.State = 22
	}
	kv4 := &types.KeyValue{Key: pidKey, Value: types.Encode(&productValue)}

	//摘牌信息保存
	var delistinfo g.Delist
	delistInfo := &g.DelistInfo{Amount: delist.Purchase, Name: delist.OpCompany.Name}
	delistValue, err := a.db.Get(gylKeyDelist(delist.ProductId))
	if err == nil { //数据内容不为空
		err := types.Decode(delistValue, &delistinfo)
		if err != nil {
			return nil, errors.New("delist decode err")
		}
	} else {
		if err != types.ErrNotFound {
			return nil,err
		}
	}

	delistinfo.Info = append(delistinfo.Info, delistInfo)

	kv5 := &types.KeyValue{Key: gylKeyDelist(delist.ProductId), Value: types.Encode(&delistinfo)}

	kv = append(kv,kv1,kv2,kv3,kv4,kv5)
	a.saveStateDB(kv)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//个人认证（废弃）
func (a *Action) PersonCertification(perCer *g.ZsgjPersonCertification) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//企业认证
func (a *Action) CompanyCertification(payload *g.ZsgjCompanyCertification, pubkey []byte) (*types.Receipt, error) {
	clog.Info("GYL", "ENTERPRISE_NAME", payload.CompanyName)
	clog.Error("GYL", "ENTERPRISE_KEY", hex.EncodeToString(payload.PubKey))
	clog.Error("GYL", "CERT_INPUT_KEY", hex.EncodeToString(pubkey))
	if bytes.Compare(payload.PubKey, pubkey) != 0 {
		return nil, errors.New("different pubkey")
	}
	value,err := a.db.Get(gylKeyUser(payload.CompanyName))
	if err == nil && value != nil {
		return nil,errors.New("user already register ") //判断公钥是否已经绑定
	}
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	zsgjAccount := &g.ZsgjAccount{}
	zsgjAccount.CompanyName = payload.CompanyName
	rmbAssets := &g.RmbAssets{Rmb:0, Cashed:0, WaitCash:0, Invest:0, Income:0}
	zsgjAccount.Rmb = rmbAssets
	val := types.Encode(zsgjAccount)
	kv2 := &types.KeyValue{Key: gylKeyUserInfo(pubkey), Value: val}
	info := &g.ZsgjCompanyCertification{Name: payload.Name, IdCard: payload.IdCard, PhoneNumber: payload.PhoneNumber, CertificateDate: payload.CertificateDate,
		Info: payload.Info, CompanyAddress: payload.CompanyAddress, CompanyName: payload.CompanyName, LicenseNumber: payload.LicenseNumber,
		LegalPersonName: payload.LegalPersonName}
	val3 := types.Encode(info)
	kv3 := &types.KeyValue{Key: calComKey(info.PhoneNumber), Value: val3}
	kv = append(kv,&types.KeyValue{Key:gylKeyUser(payload.CompanyName),Value:pubkey},kv2,kv3)
	a.saveStateDB(kv)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//单据挂牌
func (a *Action) ReceiptList(recList *g.DlReceiptList, pubkey []byte) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	clog.Info("GYL", "ENTERPRISE_NAME", recList.OpCompany)
	pub, err := a.db.Get(gylKeyUser(recList.OpCompany.Name))
	if err != nil {
		return nil, err
	}
	clog.Error("GYL", "LOCAL_KEY", hex.EncodeToString(pub))
	clog.Error("GYL", "INPUT_KEY", hex.EncodeToString(pubkey))
	if bytes.Compare(pub, pubkey) != 0 {
		return nil, errors.New("different pubkey")
	}
	//判断操作状态，在有父单据的情况下
	if recList.ReceiptParentsId != "" {
		oldId := gylKeyReceipt(recList.ReceiptParentsId)
		oldReceiptValue, err := a.db.Get(oldId)
		if err != nil {
			return nil, err
		}
		var val g.ZsgjReceiptInfo
		err = types.Decode(oldReceiptValue, &val)
		if err != nil {
			return nil, errors.New("receipt decode err")
		}
		//if val.State != 3 {
		//	clog.Error("GYL","receiptinfo",val.State)
		//	return nil, errors.New("状态错误")
		//}
		opRoleAccount := &g.ZsgjRoleAccount{}
		opRoleKey := append(gylKeyUserInfo(pubkey),[]byte(recList.OpCompany.RoleId)...)
		opRoleVal,err := a.db.Get(opRoleKey)
		if err != nil {
			if err != types.ErrNotFound {
				return nil,err
			}
			opRoleAccount.Enterprise = recList.OpCompany
		}
		if opRoleVal != nil {
			err = types.Decode(opRoleVal,opRoleAccount)
			if err != nil {
				return nil,errors.New("opRoleAccount decode err")
			}
		}


		if recList.Opty == 45 && recList.State == 41 { //单据挂牌操作
			//单据变化
			val.State = recList.State
			val.Opty = recList.Opty
			//角色资产变化
			if opRoleAccount.ReceiptAssets == nil {
				receiptAssets := &g.ReceiptAssets{}
				receiptAssets.Available -= recList.TransferAmount
				receiptAssets.Receipt -= recList.TransferAmount
				receiptAssets.Financing += recList.Finance
				opRoleAccount.ReceiptAssets = receiptAssets
			} else {
				opRoleAccount.ReceiptAssets.Available -= recList.TransferAmount
				opRoleAccount.ReceiptAssets.Receipt -= recList.TransferAmount
				opRoleAccount.ReceiptAssets.Financing += recList.Finance
			}
		}
		if recList.Opty == 46 { //单据撤牌操作
			//单据变化
			val.State = 3 //单据被发行机构持有中
			val.Opty = recList.Opty
			//角色资产变化
			if opRoleAccount.ReceiptAssets == nil {
				receiptAssets := &g.ReceiptAssets{}
				receiptAssets.Available += recList.TransferAmount
				receiptAssets.Receipt += recList.TransferAmount
				receiptAssets.Financing -= recList.Finance
				opRoleAccount.ReceiptAssets = receiptAssets
			} else {
				opRoleAccount.ReceiptAssets.Available += recList.TransferAmount
				opRoleAccount.ReceiptAssets.Receipt += recList.TransferAmount
				opRoleAccount.ReceiptAssets.Financing -= recList.Finance
			}
		}
		kv = append(kv,&types.KeyValue{Key:opRoleKey,Value:types.Encode(opRoleAccount)})

		newValue := types.Encode(&val)
		kv = append(kv, &types.KeyValue{Key: oldId, Value: newValue})
		receiptInfo := val
		receiptInfo.ReceiptId = recList.ReceiptId
		newid := gylKeyReceipt(recList.ReceiptId)
		value := types.Encode(&receiptInfo)

		kv = append(kv, &types.KeyValue{Key: newid, Value: value})
	} else {
		return nil, errors.New("没有父单据ID")
	}
	a.saveStateDB(kv)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//单据摘牌
func (a *Action) ReceiptDelist(payload *g.DlReceiptDelist, pubkey []byte) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	pub, err := a.db.Get(gylKeyUser(payload.OpCompany.Name))
	if err != nil {
		return nil, err
	}
	clog.Error("GYL", "LOCAL_KEY", hex.EncodeToString(pub))
	clog.Error("GYL", "INPUT_KEY", hex.EncodeToString(pubkey))
	if bytes.Compare(pub, pubkey) != 0 {
		return nil, errors.New("different pubkey")
	}
	if payload.UpstreamFirm == nil {
		return nil,errors.New("empty upstream data")
	}
	val, err := a.db.Get(gylKeyUserInfo(pubkey))
	if err != nil {
		return nil, err
	}
	zsgjacc := &g.ZsgjAccount{}
	err = types.Decode(val, zsgjacc)
	if err != nil {
		return nil, errors.New("account decode err")
	}
	clog.Info("zsgj", "delist", zsgjacc)
	if zsgjacc.Rmb == nil {
		return nil,errors.New("rmb account is nil")
	}
	if zsgjacc.Rmb.Rmb < payload.PayAmount {
		return nil, g.ErrRNE
	}
	zsgjacc.Rmb.Rmb -=  payload.PayAmount
	zsgjacc.Rmb.Invest += payload.PayAmount

	kv = append(kv,&types.KeyValue{Key:gylKeyUserInfo(pubkey),Value:types.Encode(zsgjacc)})

	//摘牌操作状态判断
	if payload.ReceiptParentsId == "" {
		return nil, errors.New("没有父单据ID")
	}
	oldId := gylKeyReceipt(payload.ReceiptParentsId)
	oldReceiptValue, err := a.db.Get(oldId)
	if err != nil {
		return nil, err
	}
	var oldval g.ZsgjReceiptInfo
	err = types.Decode(oldReceiptValue, &oldval)
	if err != nil {
		return nil, errors.New("receipt decode err")
	}
	if payload.Opty == 30 { //单据摘牌牌操作
		oldval.State = 48 //母单状态改为认购中
		oldval.FinanceAmount += payload.DelistAmount
		if oldval.FinanceAmount == oldval.TransferAmount {
			oldval.State = 49 //母单状态改为认购完成
		}
	}
	newValue := types.Encode(&oldval)

	srcReceiptId := gylKeyReceipt(oldval.ReceiptViceId)
	srcReceiptValue, err := a.db.Get(srcReceiptId)
	if err != nil {
		return nil, err
	}
	var srcval g.ZsgjReceiptInfo
	err = types.Decode(srcReceiptValue, &srcval)
	if err != nil {
		return nil, errors.New("receipt decode err")
	}
	srcval.FinanceAmount += payload.DelistAmount
	if srcval.FinanceAmount == srcval.SumAmount {
		srcval.State = 47 //源单据变为挂牌完成
	}
	srcValue := types.Encode(&srcval)

	kv = append(kv, &types.KeyValue{Key: oldId, Value: newValue})
	kv = append(kv, &types.KeyValue{Key: srcReceiptId, Value: srcValue})

	//上游企业资金变化
	upstreamkey := gylKeyUser(payload.UpstreamFirm.Name)
	upstreamPubkey ,err := a.db.Get(upstreamkey)
	if err != nil {
		return nil,err
	}
	upval ,err := a.db.Get(gylKeyUserInfo(upstreamPubkey))
	if err != nil {
		return nil,err
	}
	zsgjUpstreamAccount := &g.ZsgjAccount{}
	err = types.Decode(upval,zsgjUpstreamAccount)
	if err != nil {
		return nil,errors.New("zsgjUpstreamAccount decode err")
	}
	if zsgjUpstreamAccount.Rmb == nil {
		rmbassets := &g.RmbAssets{}
		rmbassets.Rmb += payload.PayAmount
		zsgjUpstreamAccount.Rmb = rmbassets
	} else {
		zsgjUpstreamAccount.Rmb.Rmb += payload.PayAmount
	}
	kv = append(kv,&types.KeyValue{Key:gylKeyUserInfo(upstreamPubkey),Value:types.Encode(zsgjUpstreamAccount)})
	//上游企业角色资金变化
	upRoleAccount := &g.ZsgjRoleAccount{}
	upRoleKey := append(gylKeyUserInfo(upstreamPubkey),[]byte(payload.UpstreamFirm.RoleId)...)
	upRoleVal ,err := a.db.Get(upRoleKey)
	if err != nil {
		if err != types.ErrNotFound {
			return nil, err
		}
		upRoleAccount.Enterprise = payload.UpstreamFirm
	}
	if upRoleVal != nil {
		err = types.Decode(upRoleVal,upRoleAccount)
		if err != nil {
			return nil,errors.New("upRoleAccount decode err")
		}
	}


	if upRoleAccount.ReceiptAssets == nil {
		receiptAssets := &g.ReceiptAssets{}
		receiptAssets.Receipt -= payload.DelistAmount
		receiptAssets.Financing -= payload.DelistAmount
		upRoleAccount.ReceiptAssets = receiptAssets
	} else {
		upRoleAccount.ReceiptAssets.Receipt -= payload.DelistAmount
		upRoleAccount.ReceiptAssets.Financing -= payload.DelistAmount
	}
	kv = append(kv,&types.KeyValue{Key:upRoleKey,Value:types.Encode(upRoleAccount)})

	//操作企业角色资产变化
	opRoleAccount := &g.ZsgjRoleAccount{}
	opRoleKey := append(gylKeyUserInfo(pubkey),[]byte(payload.OpCompany.RoleId)...)
	opRoleVal ,err := a.db.Get(opRoleKey)
	if err != nil {
		if err != types.ErrNotFound {
			return nil, err
		}
		opRoleAccount.Enterprise = payload.OpCompany
	}
	if opRoleVal != nil {
		err = types.Decode(opRoleVal,opRoleAccount)
		if err != nil {
			return nil,errors.New("opRoleAccount decode err")
		}
	}

	if opRoleAccount.ReceiptAssets == nil {
		receiptAssets := &g.ReceiptAssets{}
		receiptAssets.Receipt += payload.DelistAmount
		//receiptAssets.Available += payload.DelistAmount
		opRoleAccount.ReceiptAssets = receiptAssets
	} else {
		opRoleAccount.ReceiptAssets.Receipt += payload.DelistAmount
	}
	kv = append(kv,&types.KeyValue{Key:opRoleKey,Value:types.Encode(opRoleAccount)})

	a.saveStateDB(kv)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//单据支付
func (a *Action) ReceiptPay(payload *g.DlReceiptPay, pubkey []byte) (*types.Receipt, error) {
	pub, err := a.db.Get(gylKeyUser(payload.OpCompany.Name))
	if err != nil {
		return nil, err
	}
	clog.Error("GYL", "LOCAL_KEY", hex.EncodeToString(pub))
	clog.Error("GYL", "INPUT_KEY", hex.EncodeToString(pubkey))
	if bytes.Compare(pub, pubkey) != 0 {
		return nil, errors.New("different pubkey")
	}
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	//状态判断，单据父单据ID不为空时
	if payload.ReceiptParentsId != "" {
		//父单据信息
		oldId := gylKeyReceipt(payload.ReceiptParentsId)
		oldReceiptValue, err := a.db.Get(oldId)
		if err != nil {
			return nil, err
		}
		var oldval g.ZsgjReceiptInfo
		err = types.Decode(oldReceiptValue, &oldval)
		if err != nil {
			return nil, errors.New("receipt decode err")
		}
		//当前单据已有信息
		id := gylKeyReceipt(payload.ReceiptId)

		var val g.ZsgjReceiptInfo

		oldReceiptValue = types.Encode(&oldval)
		kv = append(kv, &types.KeyValue{Key: oldId, Value: oldReceiptValue})
		val = oldval
		val.SumAmount = payload.Amount
		val.State = payload.State
		val.Opty = payload.Opty
		val.Payed = 0
		val.Owner = payload.ReceiveCompany
		val.FinanceAmount = 0
		val.ReceiptId = payload.ReceiptId
		val.ReceiptViceId = payload.ReceiptParentsId
		value := types.Encode(&val)
		kv = append(kv, &types.KeyValue{Key: id, Value: value})

		//操作企业角色资产
		opRoleAccount := &g.ZsgjRoleAccount{}
		opRoleKey := append(gylKeyUserInfo(pubkey),[]byte(payload.OpCompany.RoleId)...)
		opRoleVal ,err := a.db.Get(opRoleKey)
		if err != nil {
			if err != types.ErrNotFound {
				return nil, err
			}
			opRoleAccount.Enterprise = payload.OpCompany
		}
		if opRoleVal != nil {
			err := types.Decode(opRoleVal,opRoleAccount)
			if err != nil {
				return nil,errors.New("opRoleAccount decode err")
			}
		}

		if payload.State == g.ZsgjState_BILL_PAY_FROZEN {//支付冻结
			if opRoleAccount.ReceiptAssets == nil {
				receiptAssets := &g.ReceiptAssets{}
				receiptAssets.Available -= payload.Amount
				receiptAssets.Payments += payload.Amount
				opRoleAccount.ReceiptAssets = receiptAssets
			} else {
				opRoleAccount.ReceiptAssets.Available -= payload.Amount
				opRoleAccount.ReceiptAssets.Payments += payload.Amount
			}
			kv = append(kv,&types.KeyValue{Key:opRoleKey,Value:types.Encode(opRoleAccount)})
		} else if payload.State == g.ZsgjState_BILL_PAY_FINISH {//确认
			//上游企业A资产变化
			if opRoleAccount.ReceiptAssets == nil {
				receiptAssets := &g.ReceiptAssets{}
				receiptAssets.Receipt -= payload.Amount
				receiptAssets.Payments -= payload.Amount
				opRoleAccount.ReceiptAssets = receiptAssets
			} else {
				opRoleAccount.ReceiptAssets.Receipt -= payload.Amount
				opRoleAccount.ReceiptAssets.Payments -= payload.Amount
			}
			kv = append(kv,&types.KeyValue{Key:opRoleKey,Value:types.Encode(opRoleAccount)})
			//上游企业B资产变化
			bpubkey, err := a.db.Get(gylKeyUser(payload.ReceiveCompany.Name))
			if err != nil {
				return nil,err
			}
			broleAccount := &g.ZsgjRoleAccount{}
			broleKey := append(gylKeyUserInfo(bpubkey), []byte(payload.ReceiveCompany.RoleId)...)
			bvalue, err := a.db.Get(broleKey)
			if err != nil {
				if err != types.ErrNotFound {
					return nil,err
				}
				broleAccount.Enterprise = payload.ReceiveCompany
			}
			if bvalue != nil {
				err = types.Decode(bvalue, broleAccount)
				if err != nil {
					return nil,err
				}
			}
			if broleAccount.ReceiptAssets == nil {
				receiptAssets := &g.ReceiptAssets{}
				receiptAssets.Receipt += payload.Amount
				receiptAssets.Available += payload.Amount
				broleAccount.ReceiptAssets = receiptAssets
			} else {
				broleAccount.ReceiptAssets.Receipt += payload.Amount
				broleAccount.ReceiptAssets.Available += payload.Amount
			}
			kv = append(kv,&types.KeyValue{Key:broleKey,Value:types.Encode(broleAccount)})
		} else if payload.State == g.ZsgjState_BILL_PAY_CANCEL && payload.Opty == g.OperationType_IOU_AGREE_PAY_REVOKE{//撤销中 && 同意支付撤销
			if opRoleAccount.ReceiptAssets == nil {
				receiptAssets := &g.ReceiptAssets{}
				receiptAssets.Available += payload.Amount
				receiptAssets.Payments -= payload.Amount
				opRoleAccount.ReceiptAssets = receiptAssets
			} else {
				opRoleAccount.ReceiptAssets.Available += payload.Amount
				opRoleAccount.ReceiptAssets.Payments -= payload.Amount
			}
			kv = append(kv,&types.KeyValue{Key:opRoleKey,Value:types.Encode(opRoleAccount)})
		}
	} else {
		return nil, errors.New("没有父单据ID")
	}
	a.saveStateDB(kv)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//保存白条
func (a *Action) BlankNote(blankNote *g.DlBlankNote, pubkey []byte) (*types.Receipt, error) {
	pub, err := a.db.Get(gylKeyUser(blankNote.OpCompany.Name))
	if err != nil {
		return nil, err
	}
	clog.Error("GYL", "LOCAL_KEY", hex.EncodeToString(pub))
	clog.Error("GYL", "INPUT_KEY", hex.EncodeToString(pubkey))
	if bytes.Compare(pub, pubkey) != 0 {
		return nil, errors.New("different pubkey")
	}
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	if blankNote.ReceiptId == "" {
		return nil, errors.New("单据不存在")
	}
	//得到单据信息
	recId := gylKeyReceipt(blankNote.ReceiptId)
	recValue, err := a.db.Get(recId)
	if err != nil {
		return nil, err
	}
	var val g.ZsgjReceiptInfo
	err = types.Decode(recValue, &val)
	if err != nil {
		return nil, errors.New("receipt decode err")
	}
	switch blankNote.Opty {
	case 32, 33://申请白条
		val.State = blankNote.State
	case 41://同意白条申请
		if val.State != 44 && blankNote.State != 61 {
			return nil, errors.New("状态错误")
		} else {
			val.State = 45 //单据状态改为白条申请完成
		}
		//操作企业角色资产变化
		roleAccount := &g.ZsgjRoleAccount{}
		upubkey,err := a.db.Get(gylKeyUser(blankNote.UpstreamFirm.Name))
		if err != nil {
			return nil,err
		}
		roleKey := append(gylKeyUserInfo(upubkey), []byte(blankNote.UpstreamFirm.RoleId)...)
		rolevalue, err := a.db.Get(roleKey)
		if err != nil {
			if err != types.ErrNotFound {
				return nil,err
			}
			roleAccount.Enterprise = blankNote.UpstreamFirm
		}
		if rolevalue != nil {
			err = types.Decode(rolevalue,roleAccount)
			if err != nil {
				return nil,err
			}
		}

		if roleAccount.BlankNoteAssets == nil {
			blankNoteAssets := &g.BlankNoteAssets{}
			blankNoteAssets.Available += blankNote.Amount
			blankNoteAssets.BlankNote += blankNote.Amount
			roleAccount.BlankNoteAssets = blankNoteAssets
		} else {
			roleAccount.BlankNoteAssets.Available += blankNote.Amount
			roleAccount.BlankNoteAssets.BlankNote += blankNote.Amount
		}
		kv = append(kv,&types.KeyValue{Key:roleKey,Value:types.Encode(roleAccount)})

		upAccount,err := a.getUserAccount(blankNote.UpstreamFirm.Name)
		if err != nil {
			return nil,err
		}
		if upAccount.Rmb == nil {
			return nil,errors.New("upaccount rmb is nil")
		}
		upAccount.Rmb.Rmb -= blankNote.BlankNoteAmount
		upAccount.Rmb.WaitCash += blankNote.ApplyBlankNote
		kv = append(kv,&types.KeyValue{Key:gylKeyUserInfo(upubkey),Value:types.Encode(upAccount)})

		//投资机构资产变化
		croleAccount := &g.ZsgjRoleAccount{}
		if blankNote.CashAgency == nil {
			return nil,errors.New("empty cashagency")
		}
		cpubkey, err := a.db.Get(gylKeyUser(blankNote.CashAgency.Name)) //承兑机构、基金方
		if err != nil {
			return nil,err
		}
		croleKey := append(gylKeyUserInfo(cpubkey), []byte(blankNote.CashAgency.RoleId)...)
		cvalue, err := a.db.Get(croleKey)
		if err != nil {
			if err != nil {
				return nil,err
			}
			croleAccount.Enterprise = blankNote.CashAgency
		}
		if cvalue != nil {
			err = types.Decode(cvalue,croleAccount)
			if err != nil {
				return nil,err
			}
		}

		if croleAccount.ReceiptAssets == nil {
			receiptAssets := &g.ReceiptAssets{}
			receiptAssets.Receipt -= blankNote.Amount
			croleAccount.ReceiptAssets = receiptAssets
		} else {
			croleAccount.ReceiptAssets.Receipt -= blankNote.Amount
		}
		kv = append(kv,&types.KeyValue{Key:croleKey,Value:types.Encode(croleAccount)})

		cAccount,err := a.getUserAccount(blankNote.CashAgency.Name)
		if err != nil {
			return nil,err
		}
		if cAccount.Rmb == nil {
			return nil,errors.New("CashAgency rmb is nil")
		}
		cAccount.Rmb.Rmb += blankNote.BlankNoteAmount
		cAccount.Rmb.WaitCash += blankNote.ApplyBlankNote
		kv = append(kv,&types.KeyValue{Key:gylKeyUserInfo(cpubkey),Value:types.Encode(cAccount)})
	case 42://拒绝白条申请
		val.State = blankNote.State
		//申请撤销
	case 34:
		val.State = blankNote.State
	default:
		return nil, errors.New("操作类型不符合")
	}
	recValue = types.Encode(&val)

	kv = append(kv, &types.KeyValue{Key: recId, Value: recValue})
	blankNoteInfo := NewBlankNoteInfo(blankNote)
	id := gylKeyBlankNote(blankNote.BlankNoteId)
	value := types.Encode(&blankNoteInfo.bzi)

	kv = append(kv, &types.KeyValue{Key: id, Value: value})



	a.saveStateDB(kv)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//白条支付
func (a *Action) BlankNotePay(blankNote *g.DlBlankNotePay, pubkey []byte) (*types.Receipt, error) {
	pub, err := a.db.Get(gylKeyUser(blankNote.OpCompany.Name))
	if err != nil {
		return nil, err
	}
	clog.Error("GYL", "LOCAL_KEY", hex.EncodeToString(pub))
	clog.Error("GYL", "INPUT_KEY", hex.EncodeToString(pubkey))
	if bytes.Compare(pub, pubkey) != 0 {
		return nil, errors.New("different pubkey")
	}
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	if blankNote.BlankNoteParentsId != "" {
		//父单据信息
		oldId := gylKeyBlankNote(blankNote.BlankNoteParentsId)
		oldBlankNoteValue, err := a.db.Get(oldId)
		if err != nil {
			return nil, err
		}
		var oldval g.ZsgjBlankNoteInfo
		err = types.Decode(oldBlankNoteValue, &oldval)
		if err != nil {
			return nil, errors.New("note decode err")
		}
		//当前单据已有信息
		id := gylKeyBlankNote(blankNote.BlankNoteId)


		//if oldval.State != 61 {
		//	return nil, errors.New("父白条状态不允许操作")
		//}

		opRoleAccount := &g.ZsgjRoleAccount{}
		opRoleKey := append(gylKeyUserInfo(pubkey),[]byte(blankNote.OpCompany.RoleId)...)
		opRoleVal ,err := a.db.Get(opRoleKey)
		if err != nil {
			if err != types.ErrNotFound {
				return nil,err
			}
			opRoleAccount.Enterprise = blankNote.OpCompany
		}
		if opRoleVal != nil {
			err = types.Decode(opRoleVal,opRoleAccount)
			if err != nil {
				return nil,err
			}
		}
		switch blankNote.Opty {
		//白条支付操作
		case g.OperationType_IOU_PAY:
			if blankNote.State != g.ZsgjState_IOU_FROZEN {
				return nil, errors.New("状态错误")
			}
			//保存新的白条id
			value, err := a.db.Get(id)
			if err == nil && value != nil {
				return nil, errors.New("BlankNoteId exists")
			}
			var val g.ZsgjBlankNoteInfo

			//TODO剩余资金可能还有问题
			val = oldval
			val.ApplyBlankNote = blankNote.PayAmount
			val.State = blankNote.State
			val.Opty = blankNote.Opty
			val.Payed = 0
			val.Owner = blankNote.ReceiveCompany
			val.BlankNoteId = blankNote.BlankNoteId
			val.BlankNoteParentsId = blankNote.BlankNoteParentsId
			value = types.Encode(&val)

			kv = append(kv, &types.KeyValue{Key: id, Value: value})


			if opRoleAccount.BlankNoteAssets == nil {
				blankNoteAssets := &g.BlankNoteAssets{}
				blankNoteAssets.BlankNote -= blankNote.Amount
				blankNoteAssets.Available -= blankNote.Amount
				blankNoteAssets.Payments += blankNote.Amount
				opRoleAccount.BlankNoteAssets = blankNoteAssets
			} else {
				opRoleAccount.BlankNoteAssets.BlankNote -= blankNote.Amount
				opRoleAccount.BlankNoteAssets.Available -= blankNote.Amount
				opRoleAccount.BlankNoteAssets.Payments += blankNote.Amount
			}
			kv = append(kv,&types.KeyValue{Key:opRoleKey,Value:types.Encode(opRoleAccount)})
		//白条确认支付操作
		case g.OperationType_IOU_PAY_CONFIRM:
			if blankNote.State != g.ZsgjState_APPLY_OK {
				return nil, errors.New("状态错误")
			}
			oldval.Payed += blankNote.PayAmount
			if oldval.Payed == oldval.ApplyBlankNote {
				oldval.State = g.ZsgjState_PAY_SUCCESS
			}
			//更新白条id状态
			value, err := a.db.Get(id)
			if err != nil {
				return nil, err
			}
			var val g.ZsgjBlankNoteInfo
			err = types.Decode(value,&val)
			if err != nil {
				return nil,err
			}
			val.State = blankNote.State
			value = types.Encode(&val)

			kv = append(kv, &types.KeyValue{Key: id, Value: value})
			//上游企业A资产变化

			if opRoleAccount.BlankNoteAssets == nil {
				blankNoteAssets := &g.BlankNoteAssets{}
				blankNoteAssets.Payments -= blankNote.Amount
				opRoleAccount.BlankNoteAssets = blankNoteAssets
			} else {
				opRoleAccount.BlankNoteAssets.Payments -= blankNote.Amount
			}
			kv = append(kv,&types.KeyValue{Key:opRoleKey,Value:types.Encode(opRoleAccount)})
			//上游企业B资产变化
			broleAccount := &g.ZsgjRoleAccount{}
			bpubkey, err := a.db.Get(gylKeyUser(blankNote.ReceiveCompany.Name))
			if err != nil {
				return nil,err
			}
			broleKey := append(gylKeyUserInfo(bpubkey), []byte(blankNote.ReceiveCompany.RoleId)...)
			bvalue, err := a.db.Get(broleKey)
			if err != nil {
				if err != types.ErrNotFound {
					return nil,err
				}
				broleAccount.Enterprise = blankNote.ReceiveCompany
			}
			if bvalue != nil {
				err = types.Decode(bvalue,broleAccount)
				if err != nil {
					return nil,err
				}
			}
			if broleAccount.BlankNoteAssets == nil {
				blankNoteAssets := &g.BlankNoteAssets{}
				blankNoteAssets.BlankNote += blankNote.Amount
				blankNoteAssets.Available += blankNote.Amount
				broleAccount.BlankNoteAssets = blankNoteAssets
			} else {
				broleAccount.BlankNoteAssets.BlankNote += blankNote.Amount
				broleAccount.BlankNoteAssets.Available += blankNote.Amount
			}
		//白条撤销支付操作
		case g.OperationType_IOU_PAY_REVOKE:
			if  blankNote.State != g.ZsgjState_IOU_PAYMENT_REVOCATION {
				return nil, errors.New("状态错误")
			}
			//更新白条id状态
			value, err := a.db.Get(id)
			if err != nil {
				return nil, err
			}
			var val g.ZsgjBlankNoteInfo
			err = types.Decode(value,&val)
			if err != nil {
				return nil,err
			}
			val.State = blankNote.State
			value = types.Encode(&val)

			kv = append(kv, &types.KeyValue{Key: id, Value: value})
		//白条同意撤销支付操作
		case g.OperationType_IOU_AGREE_PAY_REVOKE:
			if  blankNote.State != g.ZsgjState_IOU_PAYMENT_HAVE_REVOCATION {
				return nil, errors.New("状态错误")
			}
			if opRoleAccount.BlankNoteAssets == nil {
				blankNoteAssets := &g.BlankNoteAssets{}
				blankNoteAssets.BlankNote += blankNote.Amount
				blankNoteAssets.Available += blankNote.Amount
				blankNoteAssets.Payments -= blankNote.Amount
				opRoleAccount.BlankNoteAssets = blankNoteAssets
			} else {
				opRoleAccount.BlankNoteAssets.BlankNote += blankNote.Amount
				opRoleAccount.BlankNoteAssets.Available += blankNote.Amount
				opRoleAccount.BlankNoteAssets.Payments -= blankNote.Amount
			}
			//更新白条id状态
			value, err := a.db.Get(id)
			if err != nil {
				return nil, err
			}
			var val g.ZsgjBlankNoteInfo
			err = types.Decode(value,&val)
			if err != nil {
				return nil,err
			}
			val.State = blankNote.State
			value = types.Encode(&val)

			kv = append(kv, &types.KeyValue{Key: id, Value: value})
		//白条拒绝撤销支付操作
		case g.OperationType_IOU__REJECT_PAY_REVOKE:
			if  blankNote.State != g.ZsgjState_IOU_FROZEN {
				return nil, errors.New("状态错误")
			}
			//更新白条id状态
			value, err := a.db.Get(id)
			if err != nil {
				return nil, err
			}
			var val g.ZsgjBlankNoteInfo
			err = types.Decode(value,&val)
			if err != nil {
				return nil,err
			}
			val.State = blankNote.State
			value = types.Encode(&val)

			kv = append(kv, &types.KeyValue{Key: id, Value: value})
		default:
			return nil, errors.New("操作类型不符合")
		}
		oldBlankNoteValue = types.Encode(&oldval)

		kv = append(kv, &types.KeyValue{Key: oldId, Value: oldBlankNoteValue})

	} else {
		return nil, errors.New("没有父单据ID")
	}
	a.saveStateDB(kv)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//单据和白条兑付
func (a *Action) ReceiptAndNoteCash(payload *g.DlReceiptAndNoteCash, pubkey []byte) (*types.Receipt, error) {
	pub, err := a.db.Get(gylKeyUser(payload.OpCompany.Name))
	if err != nil {
		return nil, err
	}
	clog.Error("GYL", "LOCAL_KEY", hex.EncodeToString(pub))
	clog.Error("GYL", "INPUT_KEY", hex.EncodeToString(pubkey))
	if bytes.Compare(pub, pubkey) != 0 {
		return nil, errors.New("different pubkey")
	}
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	//TODO兑付状态的变化（单据、白条）

	opAccount := &g.ZsgjAccount{}
	opKey := gylKeyUserInfo(pubkey)
	opVal,err := a.db.Get(opKey)
	if err != nil {
		return nil,err
	}
	err = types.Decode(opVal,opAccount)
	if err != nil {
		return nil,err
	}

	ownerAccount := &g.ZsgjAccount{}
	ownerpubkey,err := a.db.Get(gylKeyUser(payload.OwnerEnterprise.Name))
	if err != nil {
		return nil,err
	}
	ownerKey := gylKeyUserInfo(ownerpubkey)
	ownerVal,err := a.db.Get(ownerKey)
	if err != nil {
		return nil,err
	}
	err = types.Decode(ownerVal,ownerAccount)
	if err != nil {
		return nil,err
	}

	ownerRoleAccount := &g.ZsgjRoleAccount{}
	ownerRoleKey := append(gylKeyUserInfo(ownerpubkey),[]byte(payload.OwnerEnterprise.RoleId)...)
	ownerRoleVal,err := a.db.Get(ownerRoleKey)
	if err != nil {
		if err != types.ErrNotFound {
			return nil,err
		}
		ownerRoleAccount.Enterprise = payload.OwnerEnterprise
	}
	if ownerRoleVal != nil {
		err = types.Decode(ownerRoleVal,ownerRoleAccount)
		if err != nil {
			return nil,err
		}
	}

	if payload.CashType == "单据" {
		if opAccount.Rmb == nil {
			rmbAssets := &g.RmbAssets{}
			rmbAssets.Rmb -= payload.CashAmount
			rmbAssets.Cashed += payload.CashAmount
			rmbAssets.WaitCash -= payload.CashAmount
			opAccount.Rmb = rmbAssets

		} else {
			opAccount.Rmb.Rmb -= payload.CashAmount
			opAccount.Rmb.Cashed += payload.CashAmount
			opAccount.Rmb.WaitCash -= payload.CashAmount
		}
		kv = append(kv, &types.KeyValue{Key: opKey, Value: types.Encode(opAccount)})

		if ownerAccount.Rmb == nil {
			rmbAssets := &g.RmbAssets{}
			rmbAssets.Rmb += payload.CashAmount
			ownerAccount.Rmb = rmbAssets
		} else {
			ownerAccount.Rmb.Rmb += payload.CashAmount
		}
		kv = append(kv, &types.KeyValue{Key: ownerKey, Value: types.Encode(ownerAccount)})

		if ownerRoleAccount.ReceiptAssets == nil {
			receiptAssets := &g.ReceiptAssets{}
			receiptAssets.Receipt -= payload.Amount
			if payload.OwnerEnterprise.RoleId == "2" {
				receiptAssets.Available -= payload.Amount
			}
			ownerRoleAccount.ReceiptAssets = receiptAssets
		} else {
			ownerRoleAccount.ReceiptAssets.Receipt -= payload.Amount
			if payload.OwnerEnterprise.RoleId == "2" {
				ownerRoleAccount.ReceiptAssets.Available -= payload.Amount
			}
		}
		kv = append(kv, &types.KeyValue{Key: ownerRoleKey, Value: types.Encode(ownerRoleAccount)})


	} else if payload.CashType == "白条" {
		if opAccount.Rmb == nil {
			rmbAssets := &g.RmbAssets{}
			rmbAssets.Rmb -= payload.CashAmount
			rmbAssets.Cashed += payload.CashAmount
			rmbAssets.WaitCash -= payload.CashAmount
			opAccount.Rmb = rmbAssets
		} else {
			opAccount.Rmb.Rmb -= payload.CashAmount
			opAccount.Rmb.Cashed += payload.CashAmount
			opAccount.Rmb.WaitCash -= payload.CashAmount
		}
		kv = append(kv, &types.KeyValue{Key: opKey, Value: types.Encode(opAccount)})

		if ownerAccount.Rmb == nil {
			rmbAssets := &g.RmbAssets{}
			rmbAssets.Rmb += payload.CashAmount
			ownerAccount.Rmb = rmbAssets
		} else {
			ownerAccount.Rmb.Rmb += payload.CashAmount
		}
		kv = append(kv, &types.KeyValue{Key: ownerKey, Value: types.Encode(ownerAccount)})

		if ownerRoleAccount.BlankNoteAssets == nil {
			blankNoteAssets := &g.BlankNoteAssets{}
			blankNoteAssets.Available -= payload.Amount
			blankNoteAssets.BlankNote -= payload.Amount
			ownerRoleAccount.BlankNoteAssets = blankNoteAssets
		} else {
			ownerRoleAccount.BlankNoteAssets.Available -= payload.Amount
			ownerRoleAccount.BlankNoteAssets.BlankNote -= payload.Amount
		}
	}
	a.saveStateDB(kv)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//设置信用额度
func (a *Action) SetCredit(credit *g.DlCredit, pubkey []byte) (*types.Receipt, error) {
	pub, err := a.db.Get(gylKeyUser(credit.OpCompany.Name))
	if err != nil {
		return nil, err
	}
	clog.Error("GYL", "LOCAL_KEY", hex.EncodeToString(pub))
	clog.Error("GYL", "INPUT_KEY", hex.EncodeToString(pubkey))
	if bytes.Compare(pub, pubkey) != 0 {
		return nil, errors.New("different pubkey")
	}
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}


//预付款登记
func (a *Action) AdAssetRegister(payload *g.AdAssetRegister,pubkey []byte) (*types.Receipt,error) {
	pub,err := a.db.Get(gylKeyUser(payload.OpCompany.Name))
	if err != nil {
		return nil,err
	}
	if !bytes.Equal(pub,pubkey) {
		return nil, errors.New("different pubkey")
	}
	if payload.State != g.ZsgjState_WAIT_CORE_CONFIRM_PAYMENT {
		_,err := a.db.Get(gylKeyAdance(payload.AdanceId))
		if err != nil {
			return nil,err
		}
	}
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	if payload.State == g.ZsgjState_WAIT_CORE_CONFIRM_PAYMENT {//预付款持有中，等待核心确认
		info := &g.AdanceInfo{
			AdanceId:payload.AdanceId,
			CoreCompany:payload.CoreCompany,
			DownCompany:payload.DownCompany,
			OrderAmount:payload.OrderAmount,
			SignDate:payload.SignDate,
			EndDate:payload.EndDate,
			GoodsNums:payload.GoodsNums,
			GoodsDate:payload.GoodsDate,
			GoodsName:payload.GoodsName,
			State:payload.State,
			Opty:payload.Opty,
			Rate:payload.Rate,
			OwnerCompany:payload.DownCompany,
		}
		kv1 := &types.KeyValue{Key:gylKeyAdance(payload.AdanceId),Value:types.Encode(info)}
		kv = append(kv,kv1)
	} else {//其他
		value ,err := a.db.Get(gylKeyAdance(payload.AdanceId))
		if err != nil {
			panic(err)
		}
		var info g.AdanceInfo
		err = types.Decode(value,&info)
		if err != nil {
			panic(err)
		}
		info.State = payload.State
		kv1 := &types.KeyValue{Key:gylKeyAdance(payload.AdanceId),Value:types.Encode(&info)}

		kv = append(kv,kv1)
	}
	a.saveStateDB(kv)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//预付款申请白条
func (a *Action) AdApplyBlankNote(payload *g.AdApplyBlankNote,pubkey []byte) (*types.Receipt,error) {
	pub,err := a.db.Get(gylKeyUser(payload.OpCompany.Name))
	if err != nil {
		return nil,err
	}
	if !bytes.Equal(pub,pubkey) {
		return nil, errors.New("different pubkey")
	}
	if payload.State != g.ZsgjState_BLANK_NOTE_IN_APPLYING {
		_,err := a.db.Get(gylKeyAdance(payload.AdanceId))
		if err != nil {
			return nil,err
		}
		_,err = a.db.Get(gylKeyNote(payload.BlankId))
		if err != nil {
			return nil,err
		}
	}
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	if payload.State == g.ZsgjState_BLANK_NOTE_IN_APPLYING {//预付款申请白条
		note := g.BlankNoteInfo{
			Amount:payload.Amount,
			CashAgency:payload.FinanceCompany,
			CashDate:payload.EndDate,
			Rate:payload.Rate,
			OverdueRate:payload.Rate,
			State:payload.State,
			Opty:payload.Opty,
			BlankNoteId:payload.BlankId,
			ReceiptId:payload.AdanceId,
			Owner:payload.OpCompany,
			Payed:0,
		}
		kv1 := &types.KeyValue{Key:gylKeyNote(payload.BlankId),Value:types.Encode(&note)}

		kv = append(kv,kv1)
	} else {//其他白条操作
		value,err := a.db.Get(gylKeyNote(payload.BlankId))
		if err != nil {
			panic(err)
		}
		var info g.BlankNoteInfo
		err = types.Decode(value,&info)
		if err != nil {
			panic(err)
		}
		info.State = payload.State
		kv1 := &types.KeyValue{Key:gylKeyNote(payload.BlankId),Value:types.Encode(&info)}

		kv = append(kv,kv1)
	}
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//预付款兑付
func (a *Action) AdAndNoteCash(payload *g.AdAndNoteCash,pubkey []byte) (*types.Receipt,error) {
	pub,err := a.db.Get(gylKeyUser(payload.OpCompany.Name))
	if err != nil {
		return nil,err
	}
	if !bytes.Equal(pub,pubkey) {
		return nil, errors.New("different pubkey")
	}
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	if payload.Opty == g.OperationType_OP_ADANCE_INVEST_CASH {//预付款兑付

		avalue ,err := a.db.Get(gylKeyAdance(payload.AdanceId))
		if err != nil {
			panic(err)
		}
		var adance g.AdanceInfo
		err = types.Decode(avalue,&adance)
		if err != nil {
			panic(err)
		}
		adance.State = payload.State
		kv1 := &types.KeyValue{Key:gylKeyAdance(payload.AdanceId),Value:types.Encode(&adance)}

		kv = append(kv,kv1)


	} else if payload.Opty == g.OperationType_OP_ADANCE_DOWN_PAY {//预付款同意兑付
		avalue ,err := a.db.Get(gylKeyAdance(payload.AdanceId))
		if err != nil {
			panic(err)
		}
		var adance g.AdanceInfo
		err = types.Decode(avalue,&adance)
		if err != nil {
			panic(err)
		}
		adance.State = payload.State
		kv1 := &types.KeyValue{Key:gylKeyAdance(payload.AdanceId),Value:types.Encode(&adance)}
		kv = append(kv,kv1)

		pubkey,err := a.db.Get(gylKeyUser(payload.OwnerEnterprise.Name))
		if err != nil {
			return nil,err
		}
		value,err := a.db.Get(pubkey)
		if err != nil {
			return nil,err
		}
		var invest g.ZsgjAccount
		types.Decode(value,&invest)

		downpubkey,err := a.db.Get(gylKeyUser(payload.CashEnterprise.Name))
		if err != nil {
			return nil,err
		}
		downvalue,err := a.db.Get(downpubkey)
		if err != nil {
			return nil,err
		}
		var down g.ZsgjAccount
		types.Decode(downvalue,&down)
		if down.Rmb.Rmb < payload.CashAmount {
			return nil,errors.New("rmb not enough")
		}
		down.Rmb.Rmb -= payload.CashAmount
		invest.Rmb.Rmb += payload.CashAmount
		kv2 := &types.KeyValue{Key:gylKeyUserInfo(pubkey),Value:types.Encode(&invest)}
		kv3 := &types.KeyValue{Key:gylKeyUserInfo(downpubkey),Value:types.Encode(&down)}

		kv = append(kv,kv2,kv3)
	} else if payload.Opty == g.OperationType_OP_ADANCE_CORE_CASH_IOU {//预付款申请兑付白条
		value ,err := a.db.Get(gylKeyNote(payload.BlankNoteId))
		if err != nil {
			panic(err)
		}
		var note g.BlankNoteInfo
		err = types.Decode(value,&note)
		if err != nil {
			panic(err)
		}
		note.State = payload.State
		kv1 := &types.KeyValue{Key:[]byte("mavl-gyl-noteid-"+payload.BlankNoteId),Value:types.Encode(&note)}

		kv = append(kv,kv1)
	} else if payload.Opty == g.OperationType_OP_ADANCE_INVEST_PAY_IOU {//预付款同意兑付白条
		value ,err := a.db.Get(gylKeyNote(payload.BlankNoteId))
		if err != nil {
			panic(err)
		}
		var note g.BlankNoteInfo
		err = types.Decode(value,&note)
		if err != nil {
			panic(err)
		}
		note.State = payload.State
		kv1 := &types.KeyValue{Key:gylKeyNote(payload.BlankNoteId),Value:types.Encode(&note)}

		kv = append(kv,kv1)

		pubkey,err := a.db.Get(gylKeyUser(payload.CashEnterprise.Name))
		if err != nil {
			return nil,err
		}
		investvalue,err := a.db.Get(pubkey)
		if err != nil {
			return nil,err
		}
		var invest g.ZsgjAccount
		types.Decode(investvalue,&invest)

		corepubkey,err := a.db.Get(gylKeyUser(payload.OwnerEnterprise.Name))
		if err != nil {
			return nil,err
		}
		corevalue,err := a.db.Get(corepubkey)
		if err != nil {
			return nil,err
		}
		var core g.ZsgjAccount
		types.Decode(corevalue,&core)
		if invest.Rmb.Rmb < payload.CashAmount {
			return nil,errors.New("rmb not enough")
		}
		core.Rmb.Rmb += payload.CashAmount
		invest.Rmb.Rmb -= payload.CashAmount
		kv2 := &types.KeyValue{Key:gylKeyUserInfo(pubkey),Value:types.Encode(&invest)}
		kv3 := &types.KeyValue{Key:gylKeyUserInfo(corepubkey),Value:types.Encode(&core)}

		kv = append(kv,kv2,kv3)
	}
	a.saveStateDB(kv)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//融资信息
func (a *Action) GylFinanceInfo(payload *g.GylFinanceInfo,pubkey []byte) (*types.Receipt,error) {
	if payload.OpCompany == nil {
		return nil,errors.New("opcompany is nil")
	}

	if payload.OpCompany.Name == "" {
		return nil,errors.New("opcompany name is nil")
	}
	pub,err := a.db.Get(gylKeyUser(payload.OpCompany.Name))
	if err != nil {
		return nil,err
	}
	if !bytes.Equal(pub,pubkey) {
		return nil, errors.New("different pubkey")
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	kv = append(kv,&types.KeyValue{Key:gylKeyFinanceInfo(pubkey),Value:types.Encode(payload)})

	a.saveStateDB(kv)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//定向融资
func (a *Action) GylDirectFinance(payload *g.GylDirectFinance,pubkey []byte) (*types.Receipt,error) {
	if payload.OpCompany == nil {
		return nil,errors.New("opcompany is nil")
	}

	if payload.OpCompany.Name == "" {
		return nil,errors.New("opcompany name is nil")
	}
	pub,err := a.db.Get(gylKeyUser(payload.OpCompany.Name))
	if err != nil {
		return nil,err
	}
	if !bytes.Equal(pub,pubkey) {
		return nil, errors.New("different pubkey")
	}
	if payload.FinanceCompany == nil {
		return nil,errors.New("FinanceCompany is nil")
	}

	if payload.FinanceCompany.Name == "" {
		return nil,errors.New("FinanceCompany name is nil")
	}
	fpub,err := a.db.Get(gylKeyUser(payload.FinanceCompany.Name))
	if err != nil {
		return nil,err
	}
	if payload.UpstreamCompany == nil {
		return nil,errors.New("UpstreamCompany is nil")
	}

	if payload.UpstreamCompany.Name == "" {
		return nil,errors.New("UpstreamCompany name is nil")
	}
	upub,err := a.db.Get(gylKeyUser(payload.UpstreamCompany.Name))
	if err != nil {
		return nil,err
	}

	if payload.ReceiptParentsId == "" || payload.ReceiptId == "" {
		return nil,errors.New("receipt id is empty")
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue


	if payload.Opty == g.OperationType_DIRECTIONAL_BILL {//定向融资单据挂牌
		opRoleAccount ,err := a.getUserRoleAccount(payload.OpCompany.Name,payload.OpCompany.RoleId)
		if err != nil {
			if err != types.ErrNotFound {
				return nil,err
			}
			opRoleAccount.Enterprise = payload.OpCompany
		}
		if opRoleAccount.ReceiptAssets == nil {
			receiptAssets := &g.ReceiptAssets{}
			receiptAssets.Available -= payload.TransferAmount
			receiptAssets.Receipt -= payload.TransferAmount
			receiptAssets.Financing += payload.Finance
			opRoleAccount.ReceiptAssets = receiptAssets
		} else {
			opRoleAccount.ReceiptAssets.Available -= payload.TransferAmount
			opRoleAccount.ReceiptAssets.Receipt -= payload.TransferAmount
			opRoleAccount.ReceiptAssets.Financing += payload.Finance
		}
		opRoleKey := append(gylKeyUserInfo(pub),[]byte(payload.OpCompany.RoleId)...)
		kv = append(kv,&types.KeyValue{Key:opRoleKey,Value:types.Encode(opRoleAccount)})
	} else if payload.Opty == g.OperationType_DIRECTIONAL_BILL_DELISTING {//定向融资单据摘牌
		//上游企业资金变化
		upAccount ,err := a.getUserAccount(payload.UpstreamCompany.Name)
		if err != nil {
			return nil,err
		}
		if upAccount.Rmb == nil {
			rmbassets := &g.RmbAssets{}
			rmbassets.Rmb += payload.Finance
			upAccount.Rmb = rmbassets
		} else {
			upAccount.Rmb.Rmb += payload.Finance
		}
		kv = append(kv,&types.KeyValue{Key:gylKeyUserInfo(upub),Value:types.Encode(upAccount)})


		//上游企业角色资产变化
		upRoleAccount,err := a.getUserRoleAccount(payload.UpstreamCompany.Name,payload.UpstreamCompany.RoleId)
		if err != nil {
			return nil,err
		}
		if upRoleAccount.ReceiptAssets == nil {
			return nil,errors.New("uproleaccount receiptassets is nil")
		} else {
			upRoleAccount.ReceiptAssets.Receipt -= payload.Finance
			upRoleAccount.ReceiptAssets.Financing -= payload.Finance
		}
		upRoleKey := append(gylKeyUserInfo(upub),[]byte(payload.UpstreamCompany.RoleId)...)
		kv = append(kv,&types.KeyValue{Key:upRoleKey,Value:types.Encode(upRoleAccount)})
		//资金方资金变化
		fAccount ,err := a.getUserAccount(payload.FinanceCompany.Name)
		if err != nil {
			return nil,err
		}
		if fAccount.Rmb == nil {

		} else {
			if fAccount.Rmb.Rmb < payload.Finance {
				return nil,errors.New("rmb not enough")
			}
			fAccount.Rmb.Rmb -= payload.Finance
		}
		kv = append(kv,&types.KeyValue{Key:gylKeyUserInfo(fpub),Value:types.Encode(fAccount)})
		//资金方角色资产变化
		fRoleAccount,err := a.getUserRoleAccount(payload.FinanceCompany.Name,payload.FinanceCompany.RoleId)
		if err != nil {
			return nil,err
		}
		if fRoleAccount.ReceiptAssets == nil {
			receiptAssets := &g.ReceiptAssets{}
			receiptAssets.Receipt +=  payload.ReceiptAmount
			fRoleAccount.ReceiptAssets = receiptAssets
		} else {
			fRoleAccount.ReceiptAssets.Receipt += payload.ReceiptAmount
		}
		fRoleKey := append(gylKeyUserInfo(fpub),[]byte(payload.FinanceCompany.RoleId)...)
		kv = append(kv,&types.KeyValue{Key:fRoleKey,Value:types.Encode(fRoleAccount)})

	} else if payload.Opty == g.OperationType_DIRECTIONAL_BILL_CANCEL_LIST {//定向融资挂牌撤销
		opRoleAccount,err := a.getUserRoleAccount(payload.OpCompany.Name,payload.OpCompany.RoleId)
		if err != nil {
			return nil,err
		}
		if opRoleAccount.ReceiptAssets == nil {
			return nil,errors.New("oproleaccount receiptassets is nil")
		}
		opRoleAccount.ReceiptAssets.Available += payload.TransferAmount
		opRoleAccount.ReceiptAssets.Receipt += payload.TransferAmount
		opRoleAccount.ReceiptAssets.Financing -= payload.Finance
		opRoleKey := append(gylKeyUserInfo(pub),[]byte(payload.OpCompany.RoleId)...)
		kv = append(kv,&types.KeyValue{Key:opRoleKey,Value:types.Encode(opRoleAccount)})
	} else if payload.Opty == g.OperationType_DIRECTIONAL_BILL_INVEST_REFUSE {//定向融资摘牌拒绝
		upRoleAccount,err := a.getUserRoleAccount(payload.UpstreamCompany.Name,payload.UpstreamCompany.RoleId)
		if err != nil {
			return nil,err
		}
		if upRoleAccount.ReceiptAssets == nil {
			return nil,errors.New("oproleaccount receiptassets is nil")
		}
		upRoleAccount.ReceiptAssets.Available += payload.TransferAmount
		upRoleAccount.ReceiptAssets.Receipt += payload.TransferAmount
		upRoleAccount.ReceiptAssets.Financing -= payload.Finance
		upRoleKey := append(gylKeyUserInfo(upub),[]byte(payload.OpCompany.RoleId)...)
		kv = append(kv,&types.KeyValue{Key:upRoleKey,Value:types.Encode(upRoleAccount)})
	}

	a.saveStateDB(kv)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

func (a *Action) saveStateDB(kv []*types.KeyValue) {
	for i:=0; i<len(kv) ; i++  {
		a.db.Set(kv[i].Key,kv[i].Value)
	}
}

func (a *Action) getUserPubkey(name string) (pub []byte,err error) {
	pub,err = a.db.Get(gylKeyUser(name))
	if err != nil {
		return nil,err
	}
	return pub,nil
}

func (a *Action) getUserRoleAccount(name ,id string) (account *g.ZsgjRoleAccount,err error) {
	pub,err := a.db.Get(gylKeyUser(name))
	if err != nil {
		return nil,err
	}
	roleKey := append(gylKeyUserInfo(pub),[]byte(id)...)
	roleAccountVal,err := a.db.Get(roleKey)
	if err != nil {
		return nil,err
	}

	var roleAccount g.ZsgjRoleAccount
	err = types.Decode(roleAccountVal,&roleAccount)
	if err != nil {
		return nil,err
	}

	return &roleAccount ,nil
}

func (a *Action) getUserAccount(name string) (account *g.ZsgjAccount,err error) {
	pub,err := a.db.Get(gylKeyUser(name))
	if err != nil {
		return nil,err
	}
	accountVal,err := a.db.Get(gylKeyUserInfo(pub))
	if err != nil {
		return nil,err
	}

	var aaccount g.ZsgjAccount
	err = types.Decode(accountVal,&aaccount)
	if err != nil {
		return nil,err
	}

	return &aaccount ,nil
}

func (a *Action) getReceiptInfo(id string) (receipt *g.ZsgjReceiptInfo,err error) {
	rval,err := a.db.Get(gylKeyReceipt(id))
	if err != nil {
		return nil,err
	}
	err = types.Decode(rval,receipt)
	if err != nil {
		return nil,err
	}
	return
}