package executor

import (
	gty "github.com/33cn/plugin/plugin/dapp/gyl/types"
	"github.com/33cn/chain33/types"
	//"fmt"
	"strconv"
	"errors"
)

func (g *Gyl) Query_GetAssetsByName(in *gty.ZsgjAccount) (types.Message,error) {
	var account gty.ZsgjAccount
	var roleAccount gty.ZsgjRoleAccount
	var allAccount gty.AllAccount
	db := g.GetStateDB()
	pubKey, err := db.Get(gylKeyUser(in.CompanyName))
	if err != nil {
		if err == types.ErrNotFound {
			return nil, errors.New("company not found")
		}
		return nil, err
	}
	value, err := db.Get(gylKeyUserInfo(pubKey))
	if err != nil {
		return nil, errors.New("pubkey not found")
	}
	err = types.Decode(value, &account)
	if err != nil {
		return nil, errors.New("decode error")
	}
	allAccount.Account = &account
	for i := 1; i <= 6; i++ {
		roleKey := append(gylKeyUserInfo(pubKey), []byte(strconv.Itoa(i))...)
		roleValue, err := db.Get(roleKey)
		if err != nil {
			continue
		}
		err = types.Decode(roleValue, &roleAccount)
		if err != nil {
			continue
		}
		allAccount.RoleAccount = append(allAccount.RoleAccount, &roleAccount)
	}
	return &allAccount, nil

}

func (g *Gyl) Query_GetProductById(in *gty.ZsgjProduct) (types.Message,error) {
	var productInfo gty.ZsgjProductInfo
	db := g.GetStateDB()
	value, err := db.Get(gylKeyProduct(in.ProductId))
	if err != nil {
		return nil, err
	}
	types.Decode(value, &productInfo)
	return &productInfo, nil

}

func (g *Gyl) Query_GetReceiptById(in *gty.ZsgjReceipt) (types.Message,error) {

	var receiptInfo gty.ZsgjReceiptInfo
	db := g.GetStateDB()
	value, err := db.Get(gylKeyReceipt(in.ReceiptId))
	if err != nil {
		return nil, err
	}
	types.Decode(value, &receiptInfo)
	return &receiptInfo, nil

}

func (g *Gyl) Query_GetBlankNoteInfoById(in *gty.DlBlankNote) (types.Message,error) {
	var blankNoteInfo gty.ZsgjBlankNoteInfo
	db := g.GetStateDB()
	value, err := db.Get(gylKeyBlankNote(in.BlankNoteId))
	if err != nil {
		return nil, err
	}
	types.Decode(value, &blankNoteInfo)
	return &blankNoteInfo, nil

}

func (g *Gyl) Query_GetProductDelistInfo(in *types.ReplyString) (types.Message,error) {
	var info gty.Delist
	db := g.GetStateDB()
	value, err := db.Get(gylKeyDelist(in.Data))
	if err != nil {
		return nil, err
	}
	types.Decode(value, &info)
	return &info, nil

}