package executor

import (
	gty "github.com/33cn/plugin/plugin/dapp/mall/types"
	"github.com/33cn/chain33/types"
	//"fmt"
	//"strconv"
	//"errors"
)

func (g *Mall) Query_GetUserInfoByAddr(in *types.ReplyString) (types.Message,error) {

	var userInfo gty.MallUserInfo
	db := g.GetStateDB()
	value, err := db.Get(mallKeyUser(in.Data))
	if err != nil {
		return nil, err
	}
	types.Decode(value, &userInfo)
	return &userInfo, nil

}

func (g *Mall) Query_GetGoodInfoByGoodId(in *types.ReplyString) (types.Message,error) {

	var goodInfo gty.MallAddGood
	db := g.GetStateDB()
	value, err := db.Get(mallKeyGoodInfo(in.Data))
	if err != nil {
		return nil, err
	}
	types.Decode(value, &goodInfo)
	return &goodInfo, nil
}