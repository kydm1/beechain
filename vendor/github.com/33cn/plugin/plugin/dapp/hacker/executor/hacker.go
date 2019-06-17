package executor

import (
	log "github.com/33cn/chain33/common/log/log15"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	tty "github.com/33cn/plugin/plugin/dapp/hacker/types"
	"github.com/33cn/chain33/common/address"
)

var clog = log.New("module", "execs.hacker")
var driverName = "hacker"
var adminKey = "f82cde5927ce86aab98f2ba388123b56eb165c76a48666d72ade4369f6af18a1"
var conf       = types.ConfSub(driverName)

func init() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&Hacker{}))
}

func Init(name string, sub []byte) {
	clog.Debug("register hacker execer")
	drivers.Register(GetName(), newHacker, types.GetDappFork(driverName, "Enable"))
}

func GetName() string {
	return newHacker().GetName()
}

type Hacker struct {
	drivers.DriverBase
}

func newHacker() drivers.Driver {
	n := &Hacker{}
	n.SetChild(n)
	n.SetIsFree(true)
	n.SetExecutorType(types.LoadExecutorType(driverName))
	return n
}

func (n *Hacker) GetDriverName() string {
	return driverName
}

func (n *Hacker) CheckTx(tx *types.Transaction, index int) error {
	var payload tty.HackerAction
	err := types.Decode(tx.Payload,&payload)
	if err != nil {
		return tty.ErrDocodeErr
	}
	if payload.Value == nil {
		return tty.ErrEmptyValue
	}
	switch payload.Ty {
	case tty.HackerAddBillAction:
		if payload.GetAddBill() == nil {
			return tty.ErrEmptyValue
		}

	default:
		return tty.ErrWrongActionType
	}
	return nil
}

func (n *Hacker) checkAdmin(pubkey []byte) error {

	addr := address.PubKeyToAddress(pubkey).String()
	if ok := IsSuperManager(addr);!ok {
		return tty.ErrWrongPubkey
	}
	return nil
}


// IsSuperManager is supper manager or not
func IsSuperManager(addr string) bool {
	for _, m := range conf.GStrList("superManager") {
		if addr == m {
			return true
		}
	}
	return false
}
