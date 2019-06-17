package goldbill

import (
	"github.com/33cn/plugin/plugin/dapp/goldbill/executor"
	"github.com/33cn/plugin/plugin/dapp/goldbill/types"
	"github.com/33cn/chain33/pluginmgr"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.GoldbillX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      nil,
		RPC:      nil,
	})
}
