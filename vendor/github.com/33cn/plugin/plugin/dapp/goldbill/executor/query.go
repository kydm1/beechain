package executor

import (
	g "github.com/33cn/plugin/plugin/dapp/goldbill/types"
	"github.com/33cn/chain33/types"
)

func (c *Goldbill) Query_GetGoldbillDetail(in *types.ReqString) (types.Message,error) {
	value,err := c.GetStateDB().Get(calcGoldbillDetailKey(in.Data))
	if err != nil {
		if err == types.ErrNotFound {
			return nil,types.ErrNotFound
		}
		return nil,err
	}
	var bill g.GoldbillDetail
	types.Decode(value,&bill)
	return &bill,nil
}

func (c *Goldbill) Query_GetGoldbillUserState(in *types.ReqString) (types.Message,error) {
	value,err := c.GetLocalDB().Get(calcGoldbillUserState())
	if err != nil {
		if err == types.ErrNotFound {
			return nil,types.ErrNotFound
		}
		return nil,err
	}
	var userstate g.GoldbillUserState
	types.Decode(value,&userstate)
	return &userstate,nil
}

func (c *Goldbill) Query_GetGoldbillPlatform(in *types.ReqString) (types.Message,error) {
	value,err := c.GetStateDB().Get(calcGoldbillPlatformKey())
	if err != nil {
		if err == types.ErrNotFound {
			return nil,types.ErrNotFound
		}
		return nil,err
	}
	var platform g.GoldbillPlatform
	types.Decode(value,&platform)
	return &platform,nil
}

func (c *Goldbill) Query_GetGoldbillUser(in *types.ReqString) (types.Message,error) {
	value,err := c.GetStateDB().Get(calcGoldbillUserKey(in.Data))
	if err != nil {
		if err == types.ErrNotFound {
			return nil,g.ErrUserNotExists
		}
		return nil,err
	}
	var user g.GoldbillUser
	types.Decode(value,&user)
	return &user,nil
}

func (c *Goldbill) Query_GetGoldbillUserList(in *types.ReqString) (types.Message,error) {

	value,err := c.GetStateDB().Get(calcGoldbillUserKey(in.Data))
	if err != nil {
		if err == types.ErrNotFound {
			return nil,g.ErrUserNotExists
		}
		return nil,err
	}
	var user g.GoldbillUser
	types.Decode(value,&user)
	return &user,nil
}

