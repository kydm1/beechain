package types

import (
	"github.com/33cn/chain33/types"

)

var TraceplatformX = "traceplatform"

func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(TraceplatformX))
	types.RegistorExecutor(TraceplatformX, NewType())
	types.RegisterDappFork(TraceplatformX, "Enable", 0)
}

type TraceplatformType struct {
	types.ExecTypeBase
}

func NewType() *TraceplatformType {
	c := &TraceplatformType{}
	c.SetChild(c)
	return c
}

func (g *TraceplatformType) GetPayload() types.Message {
	return &TraceplatformAction{}
}

func (g *TraceplatformType) GetTypeMap() map[string]int32 {
	return map[string]int32{
		"TraceplatformAddGood":			TraceplatformAddGoodAction,
		"TraceplatformAddGoods": 				TraceplatformAddGoodsAction,
	}
}

func (g *TraceplatformType) GetLogMap() map[int64]*types.LogInfo {
	return map[int64]*types.LogInfo{}
}
