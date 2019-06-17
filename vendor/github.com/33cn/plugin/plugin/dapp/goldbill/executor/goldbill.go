package executor

import (
	log "github.com/33cn/chain33/common/log/log15"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	gty "github.com/33cn/plugin/plugin/dapp/goldbill/types"
)

var clog = log.New("module", "execs.goldbill")
var driverName = "goldbill"
var conf       = types.ConfSub(driverName)

func init() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&Goldbill{}))
}

func Init(name string, sub []byte) {
	clog.Debug("register goldbill execer")
	drivers.Register(GetName(), newGoldbill, types.GetDappFork(driverName, "Enable"))
}

func GetName() string {
	return newGoldbill().GetName()
}

type Goldbill struct {
	drivers.DriverBase
}

func newGoldbill() drivers.Driver {
	n := &Goldbill{}
	n.SetChild(n)
	n.SetIsFree(true)
	n.SetExecutorType(types.LoadExecutorType(driverName))
	return n
}

func (n *Goldbill) GetDriverName() string {
	return driverName
}

func (n *Goldbill) CheckTx(tx *types.Transaction, index int) error {
	return nil
}

func Key(str string) (key []byte) {
	key = append(key, []byte("mavl-goldbill-")...)
	key = append(key, str...)
	return key
}


func (n *Goldbill) UpdateUserLocalState(userType gty.GoldbillUserType,add bool) ([]*types.KeyValue,error) {
	var kv[]*types.KeyValue
	key := calcGoldbillUserState()
	value,err := n.GetLocalDB().Get(key)
	if err != nil {
		if err == types.ErrNotFound {
			if userType == gty.GoldbillUserType_UT_USER {
				kv = append(kv,&types.KeyValue{Key:key,Value:types.Encode(&gty.GoldbillUserState{Usernum:int64(1),Adminnum:0})})
			} else if userType == gty.GoldbillUserType_UT_ADMIN {
				kv = append(kv,&types.KeyValue{Key:key,Value:types.Encode(&gty.GoldbillUserState{Usernum:0,Adminnum:int64(1)})})
			}
			return kv,nil
		}
		return nil,err
	}
	var userstate gty.GoldbillUserState
	types.Decode(value,&userstate)
	if userType == gty.GoldbillUserType_UT_USER {
		if add {
			userstate.Usernum ++
		} else {
			userstate.Usernum --
		}

	} else if userType == gty.GoldbillUserType_UT_ADMIN {
		if add {
			userstate.Adminnum ++
		} else {
			userstate.Adminnum --
		}
	}
	kv = append(kv,&types.KeyValue{Key:key,Value:types.Encode(&userstate)})
	return kv,nil
}