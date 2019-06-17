package traceplatform

import (
	"github.com/33cn/plugin/plugin/dapp/traceplatform/executor"
	"github.com/33cn/plugin/plugin/dapp/traceplatform/types"
	"github.com/33cn/chain33/pluginmgr"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.TraceplatformX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      nil,
		RPC:      nil,
	})
}
