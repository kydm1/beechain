// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tokennote 创建tokennote
package tokennote

import (
	"github.com/33cn/chain33/pluginmgr"
	//_ "github.com/33cn/plugin/plugin/dapp/tokennote/autotest" // register token autotest package
	"github.com/33cn/plugin/plugin/dapp/tokennote/executor"
	"github.com/33cn/plugin/plugin/dapp/tokennote/rpc"
	"github.com/33cn/plugin/plugin/dapp/tokennote/types"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.TokennoteX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		//Cmd:      commands.TokennoteCmd,
		RPC:      rpc.Init,
	})
}
