package types

import (
	"github.com/33cn/chain33/types"

	"reflect"
)

var HackerX = "hacker"

func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(HackerX))
	types.RegistorExecutor(HackerX, NewType())
	types.RegisterDappFork(HackerX, "Enable", 0)
}

type HackerType struct {
	types.ExecTypeBase
}

func NewType() *HackerType {
	c := &HackerType{}
	c.SetChild(c)
	return c
}

func (g *HackerType) GetPayload() types.Message {
	return &HackerAction{}
}

func (g *HackerType) GetTypeMap() map[string]int32 {
	return map[string]int32{
		"AddBill":			HackerAddBillAction,
	}
}

func (g *HackerType) GetLogMap() map[int64]*types.LogInfo {
	return map[int64]*types.LogInfo{
		TyLogHackerAddBill:  {Ty: reflect.TypeOf(HackerAddBill{}), Name: "LogCloseGame"},
	}
}
