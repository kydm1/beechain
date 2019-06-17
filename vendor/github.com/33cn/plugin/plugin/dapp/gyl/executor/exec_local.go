package executor

import (
	gty "github.com/33cn/plugin/plugin/dapp/gyl/types"
	"github.com/33cn/chain33/types"
)

func (g *Gyl) ExecLocal_ZsgjSaveReceipt(payload *gty.ZsgjSaveReceipt, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {

	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	//dba := g.GetLocalDB()
	//if payload.ZsgjReceipt.State == 3 { //单据签收
		//上游企业单据资产增加
		//pubKey, _ := dba.Get([]byte("mavl-gyl-" + payload.ZsgjReceipt.ReceiveCompany.Name))
		//roleKey := append(pubKey, []byte(payload.ZsgjReceipt.ReceiveCompany.RoleId)...)
		//value, _ := dba.Get(roleKey)
		//zsgjRoleAccount := &gty.ZsgjRoleAccount{}
		//types.Decode(value, zsgjRoleAccount)
		//zsgjRoleAccount.Enterprise = payload.ZsgjReceipt.ReceiveCompany
		//if zsgjRoleAccount.ReceiptAssets == nil {
		//	receiptAssets := &gty.ReceiptAssets{}
		//	receiptAssets.Available += payload.ZsgjReceipt.SumAmount
		//	receiptAssets.Receipt += payload.ZsgjReceipt.SumAmount
		//	zsgjRoleAccount.ReceiptAssets = receiptAssets
		//} else {
		//	zsgjRoleAccount.ReceiptAssets.Available += payload.ZsgjReceipt.SumAmount
		//	zsgjRoleAccount.ReceiptAssets.Receipt += payload.ZsgjReceipt.SumAmount
		//}
		//
		//value = types.Encode(zsgjRoleAccount)
		//kv := &types.KeyValue{Key: roleKey, Value: value}
		//dba.Set(kv.Key, kv.Value)
		//set.KV = append(set.KV, kv)
		////核心企业未兑付增加
		//pubKey, _ = dba.Get([]byte("mavl-gyl-" + payload.ZsgjReceipt.CoreCompany.Name))
		//value, _ = dba.Get(pubKey)
		//zsgjAccount := &gty.ZsgjAccount{}
		//types.Decode(value, zsgjAccount)
		//if zsgjAccount.Rmb == nil {
		//	rmbassets := &gty.RmbAssets{}
		//	rmbassets.WaitCash += payload.ZsgjReceipt.SumAmount
		//	zsgjAccount.Rmb = rmbassets
		//} else {
		//	zsgjAccount.Rmb.WaitCash += payload.ZsgjReceipt.SumAmount
		//}
		//
		//value = types.Encode(zsgjAccount)
		//kv = &types.KeyValue{Key: pubKey, Value: value}
		//dba.Set(kv.Key, kv.Value)
		//set.KV = append(set.KV, kv)
	//} else if payload.ZsgjReceipt.State == 8 && payload.ZsgjReceipt.Opty == 11 { //单据转让发行机构
		////上游企业单据资产减少
		//pubKey, _ := dba.Get([]byte("mavl-gyl-" + payload.ZsgjReceipt.ReceiveCompany.Name))
		//roleKey := append(pubKey, []byte(payload.ZsgjReceipt.ReceiveCompany.RoleId)...)
		//value, _ := dba.Get(roleKey)
		//zsgjRoleAccount := &gty.ZsgjRoleAccount{}
		//types.Decode(value, zsgjRoleAccount)
		//zsgjRoleAccount.Enterprise = payload.ZsgjReceipt.ReceiveCompany
		//if zsgjRoleAccount.ReceiptAssets == nil {
		//	receiptAssets := &gty.ReceiptAssets{}
		//	receiptAssets.Available -= payload.ZsgjReceipt.SumAmount
		//	receiptAssets.Receipt -= payload.ZsgjReceipt.SumAmount
		//	zsgjRoleAccount.ReceiptAssets = receiptAssets
		//} else {
		//	zsgjRoleAccount.ReceiptAssets.Available -= payload.ZsgjReceipt.SumAmount
		//	zsgjRoleAccount.ReceiptAssets.Receipt -= payload.ZsgjReceipt.SumAmount
		//}
		//value = types.Encode(zsgjRoleAccount)
		//kv := &types.KeyValue{Key: roleKey, Value: value}
		//dba.Set(kv.Key, kv.Value)
		//set.KV = append(set.KV, kv)
		////发行机构单据资产增加
		//pubKey, _ = dba.Get([]byte("mavl-gyl-" + payload.ZsgjReceipt.IssuedAgency.Name))
		//roleKey = append(pubKey, []byte(payload.ZsgjReceipt.IssuedAgency.RoleId)...)
		//value, _ = dba.Get(roleKey)
		//zsgjRoleAccount = &gty.ZsgjRoleAccount{}
		//types.Decode(value, zsgjRoleAccount)
		//zsgjRoleAccount.Enterprise = payload.ZsgjReceipt.IssuedAgency
		//if zsgjRoleAccount.ReceiptAssets == nil {
		//	receiptAssets := &gty.ReceiptAssets{}
		//	receiptAssets.Available += payload.ZsgjReceipt.SumAmount
		//	receiptAssets.Receipt += payload.ZsgjReceipt.SumAmount
		//	zsgjRoleAccount.ReceiptAssets = receiptAssets
		//} else {
		//	zsgjRoleAccount.ReceiptAssets.Available += payload.ZsgjReceipt.SumAmount
		//	zsgjRoleAccount.ReceiptAssets.Receipt += payload.ZsgjReceipt.SumAmount
		//}
		//value = types.Encode(zsgjRoleAccount)
		//kv = &types.KeyValue{Key: roleKey, Value: value}
		//dba.Set(kv.Key, kv.Value)
		//set.KV = append(set.KV, kv)
	//} else if payload.ZsgjReceipt.State == 10 { //发行机构发行产品
		////发行机构单据资产减少、产品资产增加
		//pubKey, _ := dba.Get([]byte("mavl-gyl-" + payload.ZsgjReceipt.IssuedAgency.Name))
		//roleKey := append(pubKey, []byte(payload.ZsgjReceipt.IssuedAgency.RoleId)...)
		//value, _ := dba.Get(roleKey)
		//zsgjRoleAccount := &gty.ZsgjRoleAccount{}
		//types.Decode(value, zsgjRoleAccount)
		//zsgjRoleAccount.Enterprise = payload.ZsgjReceipt.IssuedAgency
		//if zsgjRoleAccount.ReceiptAssets == nil {
		//	receiptAssets := &gty.ReceiptAssets{}
		//	receiptAssets.Available -= payload.ZsgjReceipt.SumAmount
		//	receiptAssets.Receipt -= payload.ZsgjReceipt.SumAmount
		//	zsgjRoleAccount.ReceiptAssets = receiptAssets
		//} else {
		//	zsgjRoleAccount.ReceiptAssets.Available -= payload.ZsgjReceipt.SumAmount
		//	zsgjRoleAccount.ReceiptAssets.Receipt -= payload.ZsgjReceipt.SumAmount
		//}
		//if zsgjRoleAccount.ProductAssets == nil {
		//	productAssets := &gty.ProductAssets{}
		//	productAssets.Available += payload.ZsgjReceipt.SumAmount
		//	productAssets.Product += payload.ZsgjReceipt.SumAmount
		//	zsgjRoleAccount.ProductAssets = productAssets
		//} else {
		//	zsgjRoleAccount.ProductAssets.Available += payload.ZsgjReceipt.SumAmount
		//	zsgjRoleAccount.ProductAssets.Product += payload.ZsgjReceipt.SumAmount
		//}
		//value = types.Encode(zsgjRoleAccount)
		//kv := &types.KeyValue{Key: roleKey, Value: value}
		//dba.Set(kv.Key, kv.Value)
		//set.KV = append(set.KV, kv)
		////核心机构未兑付减少
		//pubKey, _ = dba.Get([]byte("mavl-gyl-" + payload.ZsgjReceipt.CoreCompany.Name))
		//value, _ = dba.Get(pubKey)
		//zsgjAccount := &gty.ZsgjAccount{}
		//types.Decode(value, zsgjAccount)
		//if zsgjAccount.Rmb == nil {
		//	rmbassets := &gty.RmbAssets{}
		//	rmbassets.WaitCash -= payload.ZsgjReceipt.SumAmount
		//	zsgjAccount.Rmb = rmbassets
		//} else {
		//	zsgjAccount.Rmb.WaitCash -= payload.ZsgjReceipt.SumAmount
		//}
		//value = types.Encode(zsgjAccount)
		//kv = &types.KeyValue{Key: pubKey, Value: value}
		//dba.Set(kv.Key, kv.Value)
		//set.KV = append(set.KV, kv)
	//} else if payload.ZsgjReceipt.State == 8 && payload.ZsgjReceipt.Opty == 21 { //发行机构撤销发行产品，变为单据
		//TODO在产品撤销操作中改变产品资产//融资中产品资产减少
		//pubKey, _ := dba.Get([]byte("mavl-gyl-" + payload.ZsgjReceipt.IssuedAgency.Name))
		//roleKey := append(pubKey, []byte(payload.ZsgjReceipt.ReceiveCompany.RoleId)...)
		//value, _ := dba.Get(roleKey)
		//zsgjRoleAccount := &gty.ZsgjRoleAccount{}
		//types.Decode(value, zsgjRoleAccount)
		//zsgjRoleAccount.Enterprise = payload.ZsgjReceipt.ReceiveCompany
		//if zsgjRoleAccount.ReceiptAssets == nil {
		//	receiptAssets := &gty.ReceiptAssets{}
		//	receiptAssets.Available += payload.ZsgjReceipt.SumAmount
		//	receiptAssets.Receipt += payload.ZsgjReceipt.SumAmount
		//	zsgjRoleAccount.ReceiptAssets = receiptAssets
		//} else {
		//	zsgjRoleAccount.ReceiptAssets.Available += payload.ZsgjReceipt.SumAmount
		//	zsgjRoleAccount.ReceiptAssets.Receipt += payload.ZsgjReceipt.SumAmount
		//}
		//
		//value = types.Encode(zsgjRoleAccount)
		//kv := &types.KeyValue{Key: roleKey, Value: value}
		//dba.Set(kv.Key, kv.Value)
		//set.KV = append(set.KV, kv)
	//}

	return set,nil
}

func (g *Gyl) ExecLocal_ZsgjSaveProduct(payload *gty.ZsgjSaveProduct, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	//dba := g.GetLocalDB()
	//if payload.ZsgjProduct.State == 20 { //挂牌产品
	//	pubKey, _ := dba.Get([]byte("mavl-gyl-" + payload.ZsgjProduct.OpCompany.Name))
	//	roleKey := append(pubKey, []byte(payload.ZsgjProduct.OpCompany.RoleId)...)
	//	value, _ := dba.Get(roleKey)
	//	zsgjRoleAccount := &gty.ZsgjRoleAccount{}
	//	types.Decode(value, zsgjRoleAccount)
	//	zsgjRoleAccount.Enterprise = payload.ZsgjProduct.OpCompany
	//	if zsgjRoleAccount.ProductAssets == nil {
	//		productAssets := &gty.ProductAssets{}
	//		productAssets.Available -= payload.ZsgjProduct.IssueScale
	//		productAssets.Financing += payload.ZsgjProduct.IssueScale
	//		zsgjRoleAccount.ProductAssets = productAssets
	//	} else {
	//		zsgjRoleAccount.ProductAssets.Available -= payload.ZsgjProduct.IssueScale
	//		zsgjRoleAccount.ProductAssets.Financing += payload.ZsgjProduct.IssueScale
	//	}
	//
	//	value = types.Encode(zsgjRoleAccount)
	//	kv := &types.KeyValue{Key: roleKey, Value: value}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV, kv)
	//} else if payload.ZsgjProduct.State == 12 { //产品撤牌
	//	pubKey, _ := dba.Get([]byte("mavl-gyl-" + payload.ZsgjProduct.OpCompany.Name))
	//	roleKey := append(pubKey, []byte(payload.ZsgjProduct.OpCompany.RoleId)...)
	//	value, _ := dba.Get(roleKey)
	//	zsgjRoleAccount := &gty.ZsgjRoleAccount{}
	//	types.Decode(value, zsgjRoleAccount)
	//	zsgjRoleAccount.Enterprise = payload.ZsgjProduct.OpCompany
	//	if zsgjRoleAccount.ProductAssets == nil {
	//		productAssets := &gty.ProductAssets{}
	//		productAssets.Available += payload.ZsgjProduct.IssueScale
	//		productAssets.Financing -= payload.ZsgjProduct.IssueScale
	//		zsgjRoleAccount.ProductAssets = productAssets
	//	} else {
	//		zsgjRoleAccount.ProductAssets.Available += payload.ZsgjProduct.IssueScale
	//		zsgjRoleAccount.ProductAssets.Financing -= payload.ZsgjProduct.IssueScale
	//	}
	//	value = types.Encode(zsgjRoleAccount)
	//	kv := &types.KeyValue{Key: roleKey, Value: value}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV, kv)
	//}

	return set ,nil
}

func (g *Gyl) ExecLocal_AssetRegister(payload *gty.AssetRegisterAction, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	//dba := g.GetLocalDB()
	//info := gty.AssetRegisterAction{RegisterDate: payload.RegisterDate, PayCompany: payload.PayCompany, ReceiveCompany: payload.ReceiveCompany, SumAmount: payload.SumAmount,
	//	StartDate: payload.StartDate, EndDate: payload.EndDate, ContractName: payload.ContractName, ContractNo: payload.ContractNo, State: payload.State}
	//value := types.Encode(&info)
	//kv1 := &types.KeyValue{Key: []byte("mavl-gyl-rid-" + payload.ReceiptId), Value: value}
	//dba.Set(kv1.Key, kv1.Value)
	//set.KV = append(set.KV, kv1)
	//
	//kv, _ := saveInfo(&info, g.GetLocalDB())
	//if kv != nil {
	//
	//}
	//set.KV = append(set.KV, kv)

	return set,nil
}

func (g *Gyl) ExecLocal_ZsgjApplyRecharge(payload *gty.ZsgjApplyRecharge, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	//dba := g.GetLocalDB()
	//pubKey, err := dba.Get(Key(payload.CompanyName))
	//if err != nil {
	//	return nil, err
	//}
	//value, err := dba.Get(pubKey)
	//if err != nil {
	//	return nil, err
	//}
	//zsgjAccount := &gty.ZsgjAccount{}
	//types.Decode(value, zsgjAccount)
	//if zsgjAccount.Rmb == nil {
	//	rmbassets := &gty.RmbAssets{}
	//	rmbassets.Rmb += rmbassets.Rmb + payload.Amount
	//	zsgjAccount.Rmb = rmbassets
	//} else {
	//	zsgjAccount.Rmb.Rmb += payload.Amount
	//}
	//value = types.Encode(zsgjAccount)
	//kv := &types.KeyValue{Key: pubKey, Value: value}
	//dba.Set(kv.Key, kv.Value)

	return set,nil
}

func (g *Gyl) ExecLocal_ZsgjApplyWithdraw(payload *gty.ZsgjApplyWithdraw, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	//dba := g.GetLocalDB()
	//if payload.State == 27 {
	//	pubKey, _ := dba.Get(Key(payload.CompanyName)) // using company key
	//	value, _ := dba.Get(pubKey)
	//	zsgjAccount := &gty.ZsgjAccount{}
	//	types.Decode(value, zsgjAccount)
	//	if zsgjAccount.Rmb == nil {
	//		rmbassets := &gty.RmbAssets{}
	//		rmbassets.Rmb -= payload.Amount
	//		zsgjAccount.Rmb = rmbassets
	//	} else {
	//		zsgjAccount.Rmb.Rmb -= payload.Amount
	//	}
	//	value = types.Encode(zsgjAccount)
	//	kv := &types.KeyValue{Key: pubKey, Value: value}
	//	dba.Set(kv.Key, kv.Value)
	//	if kv != nil {
	//		set.KV = append(set.KV, kv)
	//	}
	//}

	return set,nil
}

func (g *Gyl) ExecLocal_ZsgjCash(payload *gty.ZsgjCash, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	//dba := g.GetLocalDB()
	//pubKey, _ := dba.Get(Key(payload.InvestAgency.Name))
	//value, _ := dba.Get(pubKey)
	//roleKey := pubKey
	//roleKey = append(roleKey, []byte(payload.InvestAgency.RoleId)...)
	//roleValue, err := dba.Get(roleKey) //企业角色对应账户
	//if err != nil {
	//	dba.Set(roleKey, nil)
	//}
	//zsgjAccount := &gty.ZsgjAccount{}
	//types.Decode(value, zsgjAccount)
	//if zsgjAccount.Rmb == nil {
	//	rmbassets := &gty.RmbAssets{}
	//	rmbassets.Rmb += payload.WaitSumAmount                     //wait cash amount
	//	rmbassets.Income += (payload.WaitSumAmount - payload.WaitAmount) //累积收益增加
	//	zsgjAccount.Rmb = rmbassets
	//} else {
	//	zsgjAccount.Rmb.Rmb += payload.WaitSumAmount                     //wait cash amount
	//	zsgjAccount.Rmb.Income += (payload.WaitSumAmount - payload.WaitAmount) //累积收益增加
	//}
	//
	//zsgjRoleAccount := &gty.ZsgjRoleAccount{}
	//types.Decode(roleValue, zsgjRoleAccount)
	//zsgjRoleAccount.Enterprise = payload.InvestAgency
	//if zsgjRoleAccount.ProductAssets == nil {
	//	productAssets := &gty.ProductAssets{}
	//	productAssets.Product -= payload.WaitAmount //产品资产减少待付本金
	//	zsgjRoleAccount.ProductAssets = productAssets
	//} else {
	//	zsgjRoleAccount.ProductAssets.Product -= payload.WaitAmount
	//}
	//roleValue = types.Encode(zsgjRoleAccount)
	//
	//value = types.Encode(zsgjAccount)
	//kv := &types.KeyValue{Key: pubKey, Value: value}
	//dba.Set(kv.Key, kv.Value)
	//set.KV = append(set.KV, kv)
	//
	//kv = &types.KeyValue{Key: roleKey, Value: roleValue}
	//dba.Set(kv.Key, kv.Value)
	//set.KV = append(set.KV, kv)
	////发行机构资产变化
	//value1, _ := dba.Get(tx.Signature.Pubkey)
	//zsgjAccount1 := &gty.ZsgjAccount{}
	//types.Decode(value1, zsgjAccount1)
	//if zsgjAccount.Rmb == nil {
	//	rmbassets := &gty.RmbAssets{}
	//	rmbassets.Rmb -= payload.WaitSumAmount
	//	rmbassets.Cashed += payload.WaitSumAmount
	//	rmbassets.WaitCash -= payload.WaitSumAmount
	//	zsgjAccount.Rmb = rmbassets
	//} else {
	//	zsgjAccount1.Rmb.Rmb -= payload.WaitSumAmount
	//	zsgjAccount1.Rmb.Cashed += payload.WaitSumAmount
	//	zsgjAccount1.Rmb.WaitCash -= payload.WaitSumAmount
	//}
	//
	//value1 = types.Encode(zsgjAccount1)
	//kv1 := &types.KeyValue{Key: tx.Signature.Pubkey, Value: value1}
	//dba.Set(kv1.Key, kv1.Value)
	//if kv1 != nil {
	//	set.KV = append(set.KV, kv1)
	//}

	return set,nil
}

func (g *Gyl) ExecLocal_ZsgjClear(payload *gty.ZsgjClear, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	//dba := g.GetLocalDB()
	//if payload.State == 23 {
	//	pubKey, _ := dba.Get([]byte("mavl-gyl-" + payload.IssuedName.Name))
	//	value, _ := dba.Get(pubKey)
	//	zsgjAccount := &gty.ZsgjAccount{}
	//	types.Decode(value, zsgjAccount)
	//	if zsgjAccount.Rmb == nil {
	//		rmbassets := &gty.RmbAssets{}
	//		rmbassets.Rmb += payload.Amount
	//		zsgjAccount.Rmb = rmbassets
	//	} else {
	//		zsgjAccount.Rmb.Rmb += payload.Amount
	//	}
	//	value = types.Encode(zsgjAccount)
	//	kv := &types.KeyValue{Key: pubKey, Value: value}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV, kv)
	//} else if payload.State == 24 { //清算失败数据回滚
	//	delistValue, err := dba.Get([]byte("mavl-gyl-delist" + payload.ProductId))
	//	if err != nil {
	//		return set, nil
	//	}
	//	var delist gty.Delist
	//	types.Decode(delistValue, &delist)
	//	for _, v := range delist.Info {
	//		pubKey, _ := dba.Get([]byte("mavl-gyl-" + v.Name)) //获得企业对应公钥
	//		value, _ := dba.Get(pubKey)
	//		zsgjAccount := &gty.ZsgjAccount{}
	//		types.Decode(value, zsgjAccount)
	//		// zsgjAccount.ProductAssets.Product -= v.Amount
	//		if zsgjAccount.Rmb == nil {
	//			rmbassets := &gty.RmbAssets{}
	//			rmbassets.Rmb += v.Amount
	//			rmbassets.Invest -= v.Amount //累积投资是否需要回滚，存疑
	//			zsgjAccount.Rmb = rmbassets
	//		} else {
	//			zsgjAccount.Rmb.Rmb += v.Amount
	//			zsgjAccount.Rmb.Invest -= v.Amount //累积投资是否需要回滚，存疑
	//		}
	//		value = types.Encode(zsgjAccount)
	//		kv := &types.KeyValue{Key: pubKey, Value: value}
	//		dba.Set(kv.Key, kv.Value)
	//		set.KV = append(set.KV, kv)
	//
	//		roleKey := append(pubKey, []byte("6")...)
	//		roleValue, _ := dba.Get(roleKey)
	//		zsgjRoleAccount := &gty.ZsgjRoleAccount{}
	//		types.Decode(roleValue, zsgjRoleAccount)
	//		enterprise := gty.Enterprise{}
	//		enterprise.Name = v.Name
	//		enterprise.RoleId = "6"
	//		zsgjRoleAccount.Enterprise = &enterprise
	//		if zsgjRoleAccount.ProductAssets == nil {
	//			productAssets := &gty.ProductAssets{}
	//			productAssets.Product -= v.Amount
	//			zsgjRoleAccount.ProductAssets = productAssets
	//		} else {
	//			zsgjRoleAccount.ProductAssets.Product -= v.Amount
	//		}
	//
	//		roleValue = types.Encode(zsgjRoleAccount)
	//		kv = &types.KeyValue{Key: roleKey, Value: roleValue}
	//		dba.Set(kv.Key, kv.Value)
	//		set.KV = append(set.KV, kv)
	//	}
	//
	//}

	return set,nil
}

func (g *Gyl) ExecLocal_ZsgjDelist(payload *gty.ZsgjDelist, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	//dba := g.GetLocalDB()
	//db := g.GetStateDB()
	//productInfo, _ := db.Get([]byte("mavl-gyl-pid-" + payload.ProductId))
	//product := &gty.ZsgjProductInfo{}
	//err := types.Decode(productInfo, product)
	//if err != nil {
	//	clog.Error("execLocal","ZsgjDelistAction",err)
	//	return set,nil
	//}
	//IssuedAgency := product.IssuedAgency
	//if IssuedAgency == nil {
	//	clog.Error("execLocal","ZsgjDelistAction","nil err")
	//	return set,nil
	//}
	////发行机构变化
	//pubkey, _ := dba.Get(Key(IssuedAgency.Name))
	//pubKey_Issued := append(pubkey, []byte(IssuedAgency.RoleId)...)
	//value_Issued, _ := dba.Get(pubKey_Issued)
	//zsgjRoleAccount_Issued := &gty.ZsgjRoleAccount{}
	//types.Decode(value_Issued, zsgjRoleAccount_Issued)
	//zsgjRoleAccount_Issued.Enterprise = product.IssuedAgency
	//if zsgjRoleAccount_Issued.ProductAssets == nil {
	//	productAssets := &gty.ProductAssets{}
	//	productAssets.Product -= payload.Purchase
	//	productAssets.Financing -= payload.Purchase
	//	zsgjRoleAccount_Issued.ProductAssets = productAssets
	//} else {
	//	zsgjRoleAccount_Issued.ProductAssets.Product -= payload.Purchase
	//	zsgjRoleAccount_Issued.ProductAssets.Financing -= payload.Purchase
	//}
	//value_Issued = types.Encode(zsgjRoleAccount_Issued)
	//kv := &types.KeyValue{Key: pubKey_Issued, Value: value_Issued}
	//dba.Set(kv.Key, kv.Value)
	//set.KV = append(set.KV, kv)
	////投资机构
	//value, _ := dba.Get(tx.Signature.Pubkey)
	//zsgjAccount := &gty.ZsgjAccount{}
	//types.Decode(value, zsgjAccount)
	//if zsgjAccount.Rmb == nil {
	//	rmbassets := &gty.RmbAssets{}
	//	rmbassets.Rmb -= payload.Purchase
	//	rmbassets.Invest += payload.Purchase
	//	zsgjAccount.Rmb = rmbassets
	//} else {
	//	zsgjAccount.Rmb.Rmb -= payload.Purchase
	//	zsgjAccount.Rmb.Invest += payload.Purchase
	//}
	//
	//value = types.Encode(zsgjAccount)
	//kv = &types.KeyValue{Key: tx.Signature.Pubkey, Value: value}
	//dba.Set(kv.Key, kv.Value)
	//if kv != nil {
	//	set.KV = append(set.KV, kv)
	//}
	//
	//roleKey := append(tx.Signature.Pubkey, []byte("6")...)
	//roleValue, err := dba.Get(roleKey) //企业角色对应账户
	//if err != nil {
	//	dba.Set(roleKey, nil)
	//}
	//zsgjRoleAccount := &gty.ZsgjRoleAccount{}
	//types.Decode(roleValue, zsgjRoleAccount)
	//zsgjRoleAccount.Enterprise = payload.OpCompany
	//if zsgjRoleAccount.ProductAssets == nil {
	//	productAssets := &gty.ProductAssets{}
	//	productAssets.Product += payload.Purchase
	//	zsgjRoleAccount.ProductAssets = productAssets
	//} else {
	//	zsgjRoleAccount.ProductAssets.Product += payload.Purchase
	//}
	//
	//roleValue = types.Encode(zsgjRoleAccount)
	//kv = &types.KeyValue{Key: roleKey, Value: roleValue}
	//dba.Set(kv.Key, kv.Value)
	//if kv != nil {
	//	set.KV = append(set.KV, kv)
	//}
	//
	//var delist gty.Delist
	//delistInfo := &gty.DelistInfo{Amount: payload.Purchase, Name: payload.OpCompany.Name}
	//delistValue, err := dba.Get([]byte("mavl-gyl-delist" + payload.ProductId))
	//if err == nil { //数据内容不为空
	//	types.Decode(delistValue, &delist)
	//}
	//delist.Info = append(delist.Info, delistInfo)
	//delistValue = types.Encode(&delist)
	//d_kv := &types.KeyValue{Key: []byte("mavl-gyl-delist" + payload.ProductId), Value: delistValue}
	//set.KV = append(set.KV, d_kv)

	return set,nil
}

func (g *Gyl) ExecLocal_ZsgjCompanyCertification(payload *gty.ZsgjCompanyCertification, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	//dba := g.GetLocalDB()
	////kv1 := &types.KeyValue{Key: Key(payload.CompanyName), Value: payload.PubKey}
	////dba.Set(kv1.Key, kv1.Value)
	////set.KV = append(set.KV, kv1)
	////zsgjAccount := &gty.ZsgjAccount{}
	////zsgjAccount.CompanyName = payload.CompanyName
	////rmbAssets := &gty.RmbAssets{0, 0, 0, 0, 0}
	////zsgjAccount.Rmb = rmbAssets
	////val := types.Encode(zsgjAccount)
	////kv2 := &types.KeyValue{Key: payload.PubKey, Value: val}
	////dba.Set(kv2.Key, kv2.Value)
	////set.KV = append(set.KV, kv2)
	//info := gty.ZsgjCompanyCertification{Name: payload.Name, IdCard: payload.IdCard, PhoneNumber: payload.PhoneNumber, CertificateDate: payload.CertificateDate,
	//	Info: payload.Info, CompanyAddress: payload.CompanyAddress, CompanyName: payload.CompanyName, LicenseNumber: payload.LicenseNumber,
	//	LegalPersonName: payload.LegalPersonName}
	//kv, _ := saveComCerInfo(&info, g.GetLocalDB())
	//set.KV = append(set.KV, kv)

	return set,nil
}

func (g *Gyl) ExecLocal_ZsgjPersonCertification(payload *gty.ZsgjPersonCertification, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	////dba := g.GetLocalDB()
	//info := gty.ZsgjPersonCertification{Name: payload.Name, IdCard: payload.IdCard, PhoneNumber: payload.PhoneNumber, CertificateDate: payload.CertificateDate,
	//	State: payload.State, Info: payload.Info}
	//kv, _ := savePerCerInfo(&info, g.GetLocalDB())
	//if kv != nil {
	//	set.KV = append(set.KV, kv)
	//}

	return set,nil
}

func (g *Gyl) ExecLocal_DlReceiptPay(payload *gty.DlReceiptPay, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	//dba := g.GetLocalDB()
	//if payload.State == 51 { //支付冻结
	//	pubKey, _ := dba.Get([]byte("mavl-gyl-" + payload.OpCompany.Name))
	//	roleKey := append(pubKey, []byte(payload.OpCompany.RoleId)...)
	//	value, _ := dba.Get(roleKey)
	//	zsgjRoleAccount := &gty.ZsgjRoleAccount{}
	//	types.Decode(value, zsgjRoleAccount)
	//	zsgjRoleAccount.Enterprise = payload.OpCompany
	//	if zsgjRoleAccount.ReceiptAssets == nil {
	//		receiptAssets := &gty.ReceiptAssets{}
	//		receiptAssets.Available -= payload.Amount
	//		receiptAssets.Payments += payload.Amount
	//		zsgjRoleAccount.ReceiptAssets = receiptAssets
	//	} else {
	//		zsgjRoleAccount.ReceiptAssets.Available -= payload.Amount
	//		zsgjRoleAccount.ReceiptAssets.Payments += payload.Amount
	//	}
	//	value = types.Encode(zsgjRoleAccount)
	//	kv := &types.KeyValue{Key: roleKey, Value: value}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV, kv)
	//} else if payload.State == 52 { //确认
	//	//上游企业A资产变化
	//	pubKey, _ := dba.Get([]byte("mavl-gyl-" + payload.OpCompany.Name))
	//	roleKey := append(pubKey, []byte(payload.OpCompany.RoleId)...)
	//	value, _ := dba.Get(roleKey)
	//	zsgjRoleAccount := &gty.ZsgjRoleAccount{}
	//	types.Decode(value, zsgjRoleAccount)
	//	zsgjRoleAccount.Enterprise = payload.OpCompany
	//	if zsgjRoleAccount.ReceiptAssets == nil {
	//		receiptAssets := &gty.ReceiptAssets{}
	//		receiptAssets.Receipt -= payload.Amount
	//		receiptAssets.Payments -= payload.Amount
	//		zsgjRoleAccount.ReceiptAssets = receiptAssets
	//	} else {
	//		zsgjRoleAccount.ReceiptAssets.Receipt -= payload.Amount
	//		zsgjRoleAccount.ReceiptAssets.Payments -= payload.Amount
	//	}
	//	value = types.Encode(zsgjRoleAccount)
	//	kv := &types.KeyValue{Key: roleKey, Value: value}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV, kv)
	//	//上游企业B资产变化
	//	pubKey, _ = dba.Get([]byte("mavl-gyl-" + payload.ReceiveCompany.Name))
	//	roleKey = append(pubKey, []byte(payload.ReceiveCompany.RoleId)...)
	//	value, _ = dba.Get(pubKey)
	//	zsgjRoleAccount = &gty.ZsgjRoleAccount{}
	//	types.Decode(value, zsgjRoleAccount)
	//	zsgjRoleAccount.Enterprise = payload.ReceiveCompany
	//	if zsgjRoleAccount.ReceiptAssets == nil {
	//		receiptAssets := &gty.ReceiptAssets{}
	//		receiptAssets.Receipt += payload.Amount
	//		receiptAssets.Available += payload.Amount
	//		zsgjRoleAccount.ReceiptAssets = receiptAssets
	//	} else {
	//		zsgjRoleAccount.ReceiptAssets.Receipt += payload.Amount
	//		zsgjRoleAccount.ReceiptAssets.Available += payload.Amount
	//	}
	//	value = types.Encode(zsgjRoleAccount)
	//	kv = &types.KeyValue{Key: roleKey, Value: value}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV, kv)
	//} else if payload.State == 53 && payload.Opty == 38 { //撤销中 && 同意支付撤销
	//	pubKey, _ := dba.Get([]byte("mavl-gyl-" + payload.OpCompany.Name))
	//	roleKey := append(pubKey, []byte(payload.OpCompany.RoleId)...)
	//	value, _ := dba.Get(roleKey)
	//	zsgjRoleAccount := &gty.ZsgjRoleAccount{}
	//	types.Decode(value, zsgjRoleAccount)
	//	zsgjRoleAccount.Enterprise = payload.OpCompany
	//	if zsgjRoleAccount.ReceiptAssets == nil {
	//		receiptAssets := &gty.ReceiptAssets{}
	//		receiptAssets.Available += payload.Amount
	//		receiptAssets.Payments -= payload.Amount
	//		zsgjRoleAccount.ReceiptAssets = receiptAssets
	//	} else {
	//		zsgjRoleAccount.ReceiptAssets.Available += payload.Amount
	//		zsgjRoleAccount.ReceiptAssets.Payments -= payload.Amount
	//	}
	//
	//	value = types.Encode(zsgjRoleAccount)
	//	kv := &types.KeyValue{Key: roleKey, Value: value}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV, kv)
	//}

	return set,nil
}

func (g *Gyl) ExecLocal_DlReceiptDelist(payload *gty.DlReceiptDelist, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	//dba := g.GetLocalDB()
	//pubKey, _ := dba.Get([]byte("mavl-gyl-" + payload.UpstreamFirm.Name))
	//value, _ := dba.Get(pubKey)
	//zsgjAccount := &gty.ZsgjAccount{}
	//types.Decode(value, zsgjAccount)
	//if zsgjAccount.Rmb == nil {
	//	rmbassets := &gty.RmbAssets{}
	//	rmbassets.Rmb += payload.PayAmount
	//	zsgjAccount.Rmb = rmbassets
	//} else {
	//	zsgjAccount.Rmb.Rmb += payload.PayAmount
	//}
	//value = types.Encode(zsgjAccount)
	//kv := &types.KeyValue{Key: pubKey, Value: value}
	//dba.Set(kv.Key, kv.Value)
	//set.KV = append(set.KV, kv)
	//
	//roleKey := append(pubKey, []byte(payload.UpstreamFirm.RoleId)...)
	//roleValue, _ := dba.Get(roleKey)
	//zsgjRoleAccount := &gty.ZsgjRoleAccount{}
	//types.Decode(value, zsgjRoleAccount)
	//zsgjRoleAccount.Enterprise = payload.UpstreamFirm
	//if zsgjRoleAccount.ReceiptAssets == nil {
	//	receiptAssets := &gty.ReceiptAssets{}
	//	receiptAssets.Receipt -= payload.DelistAmount
	//	receiptAssets.Financing -= payload.DelistAmount
	//	zsgjRoleAccount.ReceiptAssets = receiptAssets
	//} else {
	//	zsgjRoleAccount.ReceiptAssets.Receipt -= payload.DelistAmount
	//	zsgjRoleAccount.ReceiptAssets.Financing -= payload.DelistAmount
	//}
	//roleValue = types.Encode(zsgjRoleAccount)
	//kv = &types.KeyValue{Key: roleKey, Value: roleValue}
	//dba.Set(kv.Key, kv.Value)
	//set.KV = append(set.KV, kv)
	////TODO 投资机构变化
	//pubKey, _ = dba.Get([]byte("mavl-gyl-" + payload.OpCompany.Name))
	//value, _ = dba.Get(pubKey)
	//zsgjAccount = &gty.ZsgjAccount{}
	//types.Decode(value, zsgjAccount)
	//if zsgjAccount.Rmb == nil {
	//	rmbassets := &gty.RmbAssets{}
	//	rmbassets.Rmb -= payload.DelistAmount
	//	rmbassets.Invest += payload.DelistAmount
	//	zsgjAccount.Rmb = rmbassets
	//} else {
	//	zsgjAccount.Rmb.Rmb -= payload.DelistAmount
	//	zsgjAccount.Rmb.Invest += payload.DelistAmount
	//}
	//
	//value = types.Encode(zsgjAccount)
	//kv = &types.KeyValue{Key: pubKey, Value: value}
	//dba.Set(kv.Key, kv.Value)
	//set.KV = append(set.KV, kv)
	//
	//roleKey = append(pubKey, []byte(payload.OpCompany.RoleId)...)
	//roleValue, _ = dba.Get(roleKey)
	//zsgjRoleAccount = &gty.ZsgjRoleAccount{}
	//types.Decode(value, zsgjRoleAccount)
	//zsgjRoleAccount.Enterprise = payload.OpCompany
	//if zsgjRoleAccount.ReceiptAssets == nil {
	//	receiptAssets := &gty.ReceiptAssets{}
	//	receiptAssets.Receipt += payload.DelistAmount
	//	//receiptAssets.Available += payload.DelistAmount
	//	zsgjRoleAccount.ReceiptAssets = receiptAssets
	//} else {
	//	zsgjRoleAccount.ReceiptAssets.Receipt += payload.DelistAmount
	//}
	//
	//roleValue = types.Encode(zsgjRoleAccount)
	//kv = &types.KeyValue{Key: roleKey, Value: roleValue}
	//dba.Set(kv.Key, kv.Value)
	//set.KV = append(set.KV, kv)

	return set,nil
}

func (g *Gyl) ExecLocal_DlBlankNote(payload *gty.DlBlankNote, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	//dba := g.GetLocalDB()
	//pubKey, _ := dba.Get([]byte("mavl-gyl-" + payload.OpCompany.Name))
	//roleKey := append(pubKey, []byte(payload.OpCompany.RoleId)...)
	//value, _ := dba.Get(roleKey)
	//zsgjRoleAccount := &gty.ZsgjRoleAccount{}
	//types.Decode(value, zsgjRoleAccount)
	//zsgjRoleAccount.Enterprise = payload.OpCompany
	//if zsgjRoleAccount.BlankNoteAssets == nil {
	//	blankNoteAssets := &gty.BlankNoteAssets{}
	//	blankNoteAssets.Available += payload.Amount
	//	blankNoteAssets.BlankNote += payload.Amount
	//	zsgjRoleAccount.BlankNoteAssets = blankNoteAssets
	//} else {
	//	zsgjRoleAccount.BlankNoteAssets.Available += payload.Amount
	//	zsgjRoleAccount.BlankNoteAssets.BlankNote += payload.Amount
	//}
	//value = types.Encode(zsgjRoleAccount)
	//kv := &types.KeyValue{Key: roleKey, Value: value}
	//dba.Set(kv.Key, kv.Value)
	//set.KV = append(set.KV, kv)
	////投资机构资产变化
	//pubKey, _ = dba.Get([]byte("mavl-gyl-" + payload.CashAgency.Name)) //承兑机构、基金方
	//roleKey = append(pubKey, []byte(payload.CashAgency.RoleId)...)
	//value, _ = dba.Get(roleKey)
	//zsgjRoleAccount = &gty.ZsgjRoleAccount{}
	//types.Decode(value, zsgjRoleAccount)
	//zsgjRoleAccount.Enterprise = payload.CashAgency
	//if zsgjRoleAccount.ReceiptAssets == nil {
	//	receiptAssets := &gty.ReceiptAssets{}
	//	receiptAssets.Receipt -= payload.Amount
	//	zsgjRoleAccount.ReceiptAssets = receiptAssets
	//} else {
	//	zsgjRoleAccount.ReceiptAssets.Receipt -= payload.Amount
	//}
	//
	//value = types.Encode(zsgjRoleAccount)
	//kv = &types.KeyValue{Key: roleKey, Value: value}
	//dba.Set(kv.Key, kv.Value)
	//set.KV = append(set.KV, kv)

	return set,nil
}

func (g *Gyl) ExecLocal_DlReceiptList(payload *gty.DlReceiptList, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	//dba := g.GetLocalDB()
	//if payload.Opty == 45 { //单据挂牌 两个挂牌完成状态
	//	pubKey, _ := dba.Get([]byte("mavl-gyl-" + payload.OpCompany.Name))
	//	roleKey := append(pubKey, []byte(payload.OpCompany.RoleId)...)
	//	value, _ := dba.Get(roleKey)
	//	zsgjRoleAccount := &gty.ZsgjRoleAccount{}
	//	types.Decode(value, zsgjRoleAccount)
	//	zsgjRoleAccount.Enterprise = payload.OpCompany
	//	if zsgjRoleAccount.ReceiptAssets == nil {
	//		receiptAssets := &gty.ReceiptAssets{}
	//		receiptAssets.Available -= payload.TransferAmount
	//		receiptAssets.Receipt -= payload.TransferAmount
	//		receiptAssets.Financing += payload.Finance
	//		zsgjRoleAccount.ReceiptAssets = receiptAssets
	//	} else {
	//		zsgjRoleAccount.ReceiptAssets.Available -= payload.TransferAmount
	//		zsgjRoleAccount.ReceiptAssets.Receipt -= payload.TransferAmount
	//		zsgjRoleAccount.ReceiptAssets.Financing += payload.Finance
	//	}
	//
	//	value = types.Encode(zsgjRoleAccount)
	//	kv := &types.KeyValue{Key: roleKey, Value: value}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV, kv)
	//} else if payload.Opty == 46 { //单据撤牌
	//	pubKey, _ := dba.Get([]byte("mavl-gyl-" + payload.OpCompany.Name))
	//	roleKey := append(pubKey, []byte(payload.OpCompany.RoleId)...)
	//	value, _ := dba.Get(roleKey)
	//	zsgjRoleAccount := &gty.ZsgjRoleAccount{}
	//	types.Decode(value, zsgjRoleAccount)
	//	zsgjRoleAccount.Enterprise = payload.OpCompany
	//	if zsgjRoleAccount.ReceiptAssets == nil {
	//		receiptAssets := &gty.ReceiptAssets{}
	//		receiptAssets.Available += payload.TransferAmount
	//		receiptAssets.Receipt += payload.TransferAmount
	//		receiptAssets.Financing -= payload.Finance
	//		zsgjRoleAccount.ReceiptAssets = receiptAssets
	//	} else {
	//		zsgjRoleAccount.ReceiptAssets.Available += payload.TransferAmount
	//		zsgjRoleAccount.ReceiptAssets.Receipt += payload.TransferAmount
	//		zsgjRoleAccount.ReceiptAssets.Financing -= payload.Finance
	//	}
	//	value = types.Encode(zsgjRoleAccount)
	//	kv := &types.KeyValue{Key: roleKey, Value: value}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV, kv)
	//}

	return set,nil
}

func (g *Gyl) ExecLocal_DlBlankNotePay(payload *gty.DlBlankNotePay, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	//dba := g.GetLocalDB()
	//if payload.State == 63 { //支付
	//	pubKey, _ := dba.Get([]byte("mavl-gyl-" + payload.OpCompany.Name))
	//	roleKey := append(pubKey, []byte(payload.OpCompany.RoleId)...)
	//	value, _ := dba.Get(roleKey)
	//	zsgjRoleAccount := &gty.ZsgjRoleAccount{}
	//	types.Decode(value, zsgjRoleAccount)
	//	zsgjRoleAccount.Enterprise = payload.OpCompany
	//	if zsgjRoleAccount.BlankNoteAssets == nil {
	//		blankNoteAssets := &gty.BlankNoteAssets{}
	//		blankNoteAssets.BlankNote -= payload.Amount
	//		blankNoteAssets.Available -= payload.Amount
	//		blankNoteAssets.Payments += payload.Amount
	//		zsgjRoleAccount.BlankNoteAssets = blankNoteAssets
	//	} else {
	//		zsgjRoleAccount.BlankNoteAssets.BlankNote -= payload.Amount
	//		zsgjRoleAccount.BlankNoteAssets.Available -= payload.Amount
	//		zsgjRoleAccount.BlankNoteAssets.Payments += payload.Amount
	//	}
	//	value = types.Encode(zsgjRoleAccount)
	//	kv := &types.KeyValue{Key: roleKey, Value: value}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV, kv)
	//} else if payload.State == 67 { //支付确认
	//	//上游企业A资产变化
	//	pubKey, _ := dba.Get([]byte("mavl-gyl-" + payload.OpCompany.Name))
	//	roleKey := append(pubKey, []byte(payload.OpCompany.RoleId)...)
	//	value, _ := dba.Get(roleKey)
	//	zsgjRoleAccount := &gty.ZsgjRoleAccount{}
	//	types.Decode(value, zsgjRoleAccount)
	//	zsgjRoleAccount.Enterprise = payload.OpCompany
	//	if zsgjRoleAccount.BlankNoteAssets == nil {
	//		blankNoteAssets := &gty.BlankNoteAssets{}
	//		blankNoteAssets.Payments -= payload.Amount
	//		zsgjRoleAccount.BlankNoteAssets = blankNoteAssets
	//	} else {
	//		zsgjRoleAccount.BlankNoteAssets.Payments -= payload.Amount
	//	}
	//	value = types.Encode(zsgjRoleAccount)
	//	kv := &types.KeyValue{Key: roleKey, Value: value}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV, kv)
	//	//上游企业B资产变化
	//	pubKey, _ = dba.Get([]byte("mavl-gyl-" + payload.ReceiveCompany.Name))
	//	roleKey = append(pubKey, []byte(payload.ReceiveCompany.RoleId)...)
	//	value, _ = dba.Get(roleKey)
	//	zsgjRoleAccount = &gty.ZsgjRoleAccount{}
	//	types.Decode(value, zsgjRoleAccount)
	//	zsgjRoleAccount.Enterprise = payload.ReceiveCompany
	//	if zsgjRoleAccount.BlankNoteAssets == nil {
	//		blankNoteAssets := &gty.BlankNoteAssets{}
	//		blankNoteAssets.BlankNote += payload.Amount
	//		blankNoteAssets.Available += payload.Amount
	//		zsgjRoleAccount.BlankNoteAssets = blankNoteAssets
	//	} else {
	//		zsgjRoleAccount.BlankNoteAssets.BlankNote += payload.Amount
	//		zsgjRoleAccount.BlankNoteAssets.Available += payload.Amount
	//	}
	//	value = types.Encode(zsgjRoleAccount)
	//	kv = &types.KeyValue{Key: roleKey, Value: value}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV, kv)
	//} else if payload.State == 64 && payload.Opty == 38 { //白条支付撤销 && 同意撤销
	//	pubKey, _ := dba.Get([]byte("mavl-gyl-" + payload.OpCompany.Name))
	//	roleKey := append(pubKey, []byte(payload.OpCompany.RoleId)...)
	//	value, _ := dba.Get(roleKey)
	//	zsgjRoleAccount := &gty.ZsgjRoleAccount{}
	//	types.Decode(value, zsgjRoleAccount)
	//	zsgjRoleAccount.Enterprise = payload.OpCompany
	//	if zsgjRoleAccount.BlankNoteAssets == nil {
	//		blankNoteAssets := &gty.BlankNoteAssets{}
	//		blankNoteAssets.BlankNote += payload.Amount
	//		blankNoteAssets.Available += payload.Amount
	//		blankNoteAssets.Payments -= payload.Amount
	//		zsgjRoleAccount.BlankNoteAssets = blankNoteAssets
	//	} else {
	//		zsgjRoleAccount.BlankNoteAssets.BlankNote += payload.Amount
	//		zsgjRoleAccount.BlankNoteAssets.Available += payload.Amount
	//		zsgjRoleAccount.BlankNoteAssets.Payments -= payload.Amount
	//	}
	//	value = types.Encode(zsgjRoleAccount)
	//	kv := &types.KeyValue{Key: roleKey, Value: value}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV, kv)
	//}


	return set,nil
}

func (g *Gyl) ExecLocal_DlReceiptAndNoteCash(payload *gty.DlReceiptAndNoteCash, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	//dba := g.GetLocalDB()
	//if payload.CashType == "单据" {
	//	//核心企业资产变化
	//	pubKey, _ := dba.Get([]byte("mavl-gyl-" + payload.OpCompany.Name))
	//	value, _ := dba.Get(pubKey)
	//	zsgjAccount := &gty.ZsgjAccount{}
	//	types.Decode(value, zsgjAccount)
	//	if zsgjAccount.Rmb == nil {
	//		rmbAssets := &gty.RmbAssets{}
	//		rmbAssets.Rmb -= payload.CashAmount
	//		rmbAssets.Cashed += payload.CashAmount
	//		rmbAssets.WaitCash -= payload.CashAmount
	//		zsgjAccount.Rmb = rmbAssets
	//
	//	} else {
	//		zsgjAccount.Rmb.Rmb -= payload.CashAmount
	//		zsgjAccount.Rmb.Cashed += payload.CashAmount
	//		zsgjAccount.Rmb.WaitCash -= payload.CashAmount
	//	}
	//	value = types.Encode(zsgjAccount)
	//	kv := &types.KeyValue{Key: pubKey, Value: value}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV, kv)
	//	//被兑付企业资产变化
	//	pubKey, _ = dba.Get([]byte("mavl-gyl-" + payload.OwnerEnterprise.Name))
	//	value, _ = dba.Get(pubKey)
	//	zsgjAccount = &gty.ZsgjAccount{}
	//	types.Decode(value, zsgjAccount)
	//	if zsgjAccount.Rmb == nil {
	//		rmbAssets := &gty.RmbAssets{}
	//		rmbAssets.Rmb += payload.CashAmount
	//		zsgjAccount.Rmb = rmbAssets
	//	} else {
	//		zsgjAccount.Rmb.Rmb += payload.CashAmount
	//	}
	//
	//	value = types.Encode(zsgjAccount)
	//	kv = &types.KeyValue{Key: pubKey, Value: value}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV, kv)
	//
	//	roleKey := append(pubKey, []byte(payload.OwnerEnterprise.RoleId)...)
	//	roleValue, _ := dba.Get(roleKey)
	//	zsgjRoleAccount := &gty.ZsgjRoleAccount{}
	//	types.Decode(roleValue, zsgjRoleAccount)
	//	zsgjRoleAccount.Enterprise = payload.OwnerEnterprise
	//	if zsgjRoleAccount.ReceiptAssets == nil {
	//		receiptAssets := &gty.ReceiptAssets{}
	//		receiptAssets.Receipt -= payload.Amount
	//		if payload.OwnerEnterprise.RoleId == "2" {
	//			receiptAssets.Available -= payload.Amount
	//		}
	//		zsgjRoleAccount.ReceiptAssets = receiptAssets
	//	} else {
	//		zsgjRoleAccount.ReceiptAssets.Receipt -= payload.Amount
	//		if payload.OwnerEnterprise.RoleId == "2" {
	//			zsgjRoleAccount.ReceiptAssets.Available -= payload.Amount
	//		}
	//	}
	//	roleValue = types.Encode(zsgjRoleAccount)
	//	kv = &types.KeyValue{Key: roleKey, Value: value}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV, kv)
	//} else if payload.CashType == "白条" {
	//	//投资机构变化
	//	pubKey, _ := dba.Get([]byte("mavl-gyl-" + payload.OpCompany.Name))
	//	value, _ := dba.Get(pubKey)
	//	zsgjAccount := &gty.ZsgjAccount{}
	//	types.Decode(value, zsgjAccount)
	//	if zsgjAccount.Rmb == nil {
	//		rmbAssets := &gty.RmbAssets{}
	//		rmbAssets.Rmb -= payload.CashAmount
	//		rmbAssets.Cashed += payload.CashAmount
	//		rmbAssets.WaitCash -= payload.CashAmount
	//		zsgjAccount.Rmb = rmbAssets
	//	} else {
	//		zsgjAccount.Rmb.Rmb -= payload.CashAmount
	//		zsgjAccount.Rmb.Cashed += payload.CashAmount
	//		zsgjAccount.Rmb.WaitCash -= payload.CashAmount
	//	}
	//	value = types.Encode(zsgjAccount)
	//	kv := &types.KeyValue{Key: pubKey, Value: value}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV, kv)
	//	//上游企业变化
	//	pubKey, _ = dba.Get([]byte("mavl-gyl-" + payload.OwnerEnterprise.Name))
	//	value, _ = dba.Get(pubKey)
	//	zsgjAccount = &gty.ZsgjAccount{}
	//	types.Decode(value, zsgjAccount)
	//	if zsgjAccount.Rmb == nil {
	//		rmbAssets := &gty.RmbAssets{}
	//		rmbAssets.Rmb += payload.CashAmount
	//		zsgjAccount.Rmb = rmbAssets
	//	} else {
	//		zsgjAccount.Rmb.Rmb += payload.CashAmount
	//	}
	//
	//	value = types.Encode(zsgjAccount)
	//	kv = &types.KeyValue{Key: pubKey, Value: value}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV, kv)
	//
	//	roleKey := append(pubKey, []byte(payload.OwnerEnterprise.RoleId)...)
	//	roleValue, _ := dba.Get(roleKey)
	//	zsgjRoleAccount := &gty.ZsgjRoleAccount{}
	//	types.Decode(roleValue, zsgjRoleAccount)
	//	zsgjRoleAccount.Enterprise = payload.OwnerEnterprise
	//	if zsgjRoleAccount.BlankNoteAssets == nil {
	//		blankNoteAssets := &gty.BlankNoteAssets{}
	//		blankNoteAssets.Available -= payload.Amount
	//		blankNoteAssets.BlankNote -= payload.Amount
	//		zsgjRoleAccount.BlankNoteAssets = blankNoteAssets
	//	} else {
	//		zsgjRoleAccount.BlankNoteAssets.Available -= payload.Amount
	//		zsgjRoleAccount.BlankNoteAssets.BlankNote -= payload.Amount
	//	}
	//	roleValue = types.Encode(zsgjRoleAccount)
	//	kv = &types.KeyValue{Key: roleKey, Value: value}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV, kv)
	//}
	return set,nil
}

func (g *Gyl) ExecLocal_DlInvestCash(payload *gty.DlInvestCash, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {

	return nil,nil
}

func (g *Gyl) ExecLocal_DlCredit(payload *gty.DlCredit, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	//set := &types.LocalDBSet{}
	//if receipt.GetTy() != types.ExecOk {
	//	return set, nil
	//}
	//dba := g.GetLocalDB()


	return nil,nil
}

func (g *Gyl) ExecLocal_AdAssetRegister(payload *gty.AdAssetRegister, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	//dba := g.GetLocalDB()
	//if payload.State == gty.ZsgjState_WAIT_CORE_CONFIRM_PAYMENT {//预付款持有中，等待核心确认
	//	info := &gty.AdanceInfo{
	//		 payload.AdanceId,
	//		 payload.CoreCompany,
	//		 payload.DownCompany,
	//		 payload.OrderAmount,
	//		 payload.SignDate,
	//		 payload.EndDate,
	//		 payload.GoodsNums,
	//		 payload.GoodsDate,
	//		 payload.GoodsName,
	//		 payload.State,
	//		 payload.Opty,
	//		 payload.Rate,
	//		 payload.DownCompany,
	//	}
	//	kv := &types.KeyValue{[]byte("mavl-gyl-aid-"+payload.AdanceId),types.Encode(info)}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV,kv)
	//} else {//其他
	//	value ,err := dba.Get([]byte("mavl-gyl-aid-"+payload.AdanceId))
	//	if err != nil {
	//		panic(err)
	//	}
	//	var info gty.AdanceInfo
	//	err = types.Decode(value,&info)
	//	if err != nil {
	//		panic(err)
	//	}
	//	info.State = payload.State
	//	kv := &types.KeyValue{[]byte("mavl-gyl-aid-"+payload.AdanceId),types.Encode(&info)}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV,kv)
	//}

	return set,nil
}

func (g *Gyl) ExecLocal_AdApplyBlankNote(payload *gty.AdApplyBlankNote, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	//dba := g.GetLocalDB()
	//if payload.State == gty.ZsgjState_BLANK_NOTE_IN_APPLYING {//预付款申请白条
	//	note := gty.BlankNoteInfo{
	//		 payload.Amount,
	//		 payload.FinanceCompany,
	//		 payload.EndDate,
	//		 payload.Rate,
	//		 payload.Rate,
	//		 payload.State,
	//		 payload.Opty,
	//		 payload.BlankId,
	//		 payload.AdanceId,
	//		 payload.OpCompany,
	//		0,
	//	}
	//	kv := &types.KeyValue{[]byte("mavl-gyl-noteid-"+payload.BlankId),types.Encode(&note)}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV,kv)
	//} else {//其他白条操作
	//	value,err := dba.Get([]byte("mavl-gyl-noteid-"+payload.BlankId))
	//	if err != nil {
	//		panic(err)
	//	}
	//	var info gty.BlankNoteInfo
	//	err = types.Decode(value,&info)
	//	if err != nil {
	//		panic(err)
	//	}
	//	info.State = payload.State
	//	kv := &types.KeyValue{[]byte("mavl-gyl-noteid-"+payload.BlankId),types.Encode(&info)}
	//	dba.Set(kv.Key, kv.Value)
	//	set.KV = append(set.KV,kv)
	//}

		return set,nil
}

func (g *Gyl) ExecLocal_AdAndNoteCash(payload *gty.AdAndNoteCash, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	//dba := g.GetLocalDB()
	//if payload.Opty == gty.OperationType_OP_ADANCE_INVEST_CASH {//预付款兑付
	//
	//	avalue ,err := dba.Get([]byte("mavl-gyl-aid-"+payload.AdanceId))
	//	if err != nil {
	//		panic(err)
	//	}
	//	var adance gty.AdanceInfo
	//	err = types.Decode(avalue,&adance)
	//	if err != nil {
	//		panic(err)
	//	}
	//	adance.State = payload.State
	//	kv := &types.KeyValue{[]byte("mavl-gyl-aid-"+payload.AdanceId),types.Encode(&adance)}
	//	dba.Set(kv.Key,kv.Value)
	//	set.KV = append(set.KV,kv)
	//
	//
	//} else if payload.Opty == gty.OperationType_OP_ADANCE_DOWN_PAY {//预付款同意兑付
	//	avalue ,err := dba.Get([]byte("mavl-gyl-aid-"+payload.AdanceId))
	//	if err != nil {
	//		panic(err)
	//	}
	//	var adance gty.AdanceInfo
	//	err = types.Decode(avalue,&adance)
	//	if err != nil {
	//		panic(err)
	//	}
	//	adance.State = payload.State
	//	kv := &types.KeyValue{[]byte("mavl-gyl-aid-"+payload.AdanceId),types.Encode(&adance)}
	//	dba.Set(kv.Key,kv.Value)
	//	set.KV = append(set.KV,kv)
	//
	//	pubkey,err := dba.Get([]byte("mavl-gyl-"+payload.OwnerEnterprise.Name))
	//	if err != nil {
	//		return nil,err
	//	}
	//	value,err := dba.Get(pubkey)
	//	if err != nil {
	//		return nil,err
	//	}
	//	var invest gty.ZsgjAccount
	//	types.Decode(value,&invest)
	//
	//	downpubkey,err := dba.Get([]byte("mavl-gyl-"+payload.CashEnterprise.Name))
	//	if err != nil {
	//		return nil,err
	//	}
	//	downvalue,err := dba.Get(downpubkey)
	//	if err != nil {
	//		return nil,err
	//	}
	//	var down gty.ZsgjAccount
	//	types.Decode(downvalue,&down)
	//	if down.Rmb.Rmb < payload.CashAmount {
	//		return nil,errors.New("rmb not enough")
	//	}
	//	down.Rmb.Rmb -= payload.CashAmount
	//	invest.Rmb.Rmb += payload.CashAmount
	//	kv1 := &types.KeyValue{pubkey,types.Encode(&invest)}
	//	kv2 := &types.KeyValue{downpubkey,types.Encode(&down)}
	//	dba.Set(kv1.Key,kv1.Value)
	//	dba.Set(kv2.Key,kv2.Value)
	//	set.KV = append(set.KV,kv1,kv2)
	//} else if payload.Opty == gty.OperationType_OP_ADANCE_CORE_CASH_IOU {//预付款申请兑付白条
	//	value ,err := dba.Get([]byte("mavl-gyl-noteid-"+payload.BlankNoteId))
	//	if err != nil {
	//		panic(err)
	//	}
	//	var note gty.BlankNoteInfo
	//	err = types.Decode(value,&note)
	//	if err != nil {
	//		panic(err)
	//	}
	//	note.State = payload.State
	//	kv := &types.KeyValue{[]byte("mavl-gyl-noteid-"+payload.BlankNoteId),types.Encode(&note)}
	//	dba.Set(kv.Key,kv.Value)
	//	set.KV = append(set.KV,kv)
	//} else if payload.Opty == gty.OperationType_OP_ADANCE_INVEST_PAY_IOU {//预付款同意兑付白条
	//	value ,err := dba.Get([]byte("mavl-gyl-noteid-"+payload.BlankNoteId))
	//	if err != nil {
	//		panic(err)
	//	}
	//	var note gty.BlankNoteInfo
	//	err = types.Decode(value,&note)
	//	if err != nil {
	//		panic(err)
	//	}
	//	note.State = payload.State
	//	kv := &types.KeyValue{[]byte("mavl-gyl-noteid-"+payload.BlankNoteId),types.Encode(&note)}
	//	dba.Set(kv.Key,kv.Value)
	//	set.KV = append(set.KV,kv)
	//
	//	pubkey,err := dba.Get([]byte("mavl-gyl-"+payload.CashEnterprise.Name))
	//	if err != nil {
	//		return nil,err
	//	}
	//	investvalue,err := dba.Get(pubkey)
	//	if err != nil {
	//		return nil,err
	//	}
	//	var invest gty.ZsgjAccount
	//	types.Decode(investvalue,&invest)
	//
	//	corepubkey,err := dba.Get([]byte("mavl-gyl-"+payload.OwnerEnterprise.Name))
	//	if err != nil {
	//		return nil,err
	//	}
	//	corevalue,err := dba.Get(corepubkey)
	//	if err != nil {
	//		return nil,err
	//	}
	//	var core gty.ZsgjAccount
	//	types.Decode(corevalue,&core)
	//	if invest.Rmb.Rmb < payload.CashAmount {
	//		return nil,errors.New("rmb not enough")
	//	}
	//	core.Rmb.Rmb += payload.CashAmount
	//	invest.Rmb.Rmb -= payload.CashAmount
	//	kv1 := &types.KeyValue{pubkey,types.Encode(&invest)}
	//	kv2 := &types.KeyValue{corepubkey,types.Encode(&core)}
	//	dba.Set(kv1.Key,kv1.Value)
	//	dba.Set(kv2.Key,kv2.Value)
	//	set.KV = append(set.KV,kv1,kv2)
	//}

	return set,nil
}

func (g *Gyl) ExecLocal_GylFinanceInfo(payload *gty.GylFinanceInfo, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	return set,nil
}


func (g *Gyl) ExecLocal_GylDirectFinance(payload *gty.GylDirectFinance, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}
	return set,nil
}
