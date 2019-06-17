package executor

import (
	log "github.com/33cn/chain33/common/log/log15"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
)

var clog = log.New("module", "execs.gyl")
var driverName = "gyl"

func init() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&Gyl{}))
}

func Init(name string, sub []byte) {
	clog.Debug("register gyl execer")
	drivers.Register(GetName(), newGyl, types.GetDappFork(driverName, "Enable"))
}

func GetName() string {
	return newGyl().GetName()
}

type Gyl struct {
	drivers.DriverBase
}

func newGyl() drivers.Driver {
	n := &Gyl{}
	n.SetChild(n)
	n.SetIsFree(true)
	n.SetExecutorType(types.LoadExecutorType(driverName))
	return n
}

func (n *Gyl) GetDriverName() string {
	return driverName
}

func (n *Gyl) CheckTx(tx *types.Transaction, index int) error {
	return nil
}


