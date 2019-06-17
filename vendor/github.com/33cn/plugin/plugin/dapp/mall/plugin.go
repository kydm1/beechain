package mall

import (
	"github.com/33cn/plugin/plugin/dapp/mall/executor"
	"github.com/33cn/plugin/plugin/dapp/mall/types"
	"github.com/33cn/chain33/pluginmgr"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.MallX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      nil,
		RPC:      nil,
	})
}
