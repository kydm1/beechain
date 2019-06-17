package gyl

import (
	"github.com/33cn/plugin/plugin/dapp/gyl/executor"
	"github.com/33cn/plugin/plugin/dapp/gyl/types"
	"github.com/33cn/chain33/pluginmgr"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.GylX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      nil,
		RPC:      nil,
	})
}
