// +build go1.9

package main

import (
	_ "github.com/33cn/chain33/system"
	"github.com/33cn/chain33/util/cli"
	_ "github.com/kydm1/beechain/plugin"
)

func main() {
	cli.RunChain33("beechain")
}
