package executor

import (
	log "github.com/33cn/chain33/common/log/log15"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	tty "github.com/33cn/plugin/plugin/dapp/traceplatform/types"
	"github.com/33cn/chain33/common/address"
)

var clog = log.New("module", "execs.traceplatform")
var driverName = "traceplatform"
var adminKey = "f82cde5927ce86aab98f2ba388123b56eb165c76a48666d72ade4369f6af18a1"
var conf       = types.ConfSub(driverName)

func init() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&Traceplatform{}))
}

func Init(name string, sub []byte) {
	clog.Debug("register traceplatform execer")
	drivers.Register(GetName(), newTraceplatform, types.GetDappFork(driverName, "Enable"))
}

func GetName() string {
	return newTraceplatform().GetName()
}

type Traceplatform struct {
	drivers.DriverBase
}

func newTraceplatform() drivers.Driver {
	n := &Traceplatform{}
	n.SetChild(n)
	n.SetIsFree(true)
	n.SetExecutorType(types.LoadExecutorType(driverName))
	return n
}

func (n *Traceplatform) GetDriverName() string {
	return driverName
}

func (n *Traceplatform) CheckTx(tx *types.Transaction, index int) error {
	var payload tty.TraceplatformAction
	err := types.Decode(tx.Payload,&payload)
	if err != nil {
		return tty.ErrDocodeErr
	}
	if payload.Value == nil {
		return tty.ErrEmptyValue
	}
	switch payload.Ty {
	case tty.TraceplatformAddGoodAction:
		if payload.GetTraceplatformAddGood().Goodinfo == nil || payload.GetTraceplatformAddGood().Goodinfo.NfcCode == "" {
			return tty.ErrEmptyValue
		}
		return n.checkAdmin(tx.Signature.Pubkey)
	case tty.TraceplatformAddGoodsAction:
		if len(payload.GetTraceplatformAddGoods().NfcCode) == 0 {
			return tty.ErrEmptyValue
		}
		//检查是否有重复nfcCode
		m := make(map[string]string,0)
		for _,v := range payload.GetTraceplatformAddGoods().NfcCode {
			if _,ok := m[v];!ok {
				m[v] = v
				continue
			}
			return tty.ErrDupNfcCode
		}
		return n.checkAdmin(tx.Signature.Pubkey)
	default:
		return tty.ErrWrongActionType
	}
	return nil
}

func (n *Traceplatform) checkAdmin(pubkey []byte) error {

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
