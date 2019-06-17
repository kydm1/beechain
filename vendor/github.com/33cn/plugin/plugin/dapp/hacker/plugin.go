package hacker

import (
	"github.com/33cn/plugin/plugin/dapp/hacker/executor"
	"github.com/33cn/plugin/plugin/dapp/hacker/types"
	"github.com/33cn/chain33/pluginmgr"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.HackerX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      nil,
		RPC:      nil,
	})
}
