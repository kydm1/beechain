// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pbft

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/33cn/chain33/blockchain"
	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/crypto"
	"github.com/33cn/chain33/common/limits"
	"github.com/33cn/chain33/common/log"
	"github.com/33cn/chain33/executor"
	"github.com/33cn/chain33/mempool"
	"github.com/33cn/chain33/p2p"
	"github.com/33cn/chain33/queue"
	"github.com/33cn/chain33/store"
	cty "github.com/33cn/chain33/system/dapp/coins/types"
	"github.com/33cn/chain33/types"
	"github.com/33cn/chain33/wallet"

	_ "github.com/33cn/chain33/system"
	_ "github.com/33cn/plugin/plugin/dapp/init"
	_ "github.com/33cn/plugin/plugin/store/init"
)

var (
	random       *rand.Rand
	transactions []*types.Transaction
	txSize       = 1000
)

func init() {
	err := limits.SetLimits()
	if err != nil {
		panic(err)
	}
	random = rand.New(rand.NewSource(types.Now().UnixNano()))
	log.SetLogLevel("info")
}
func TestPbft(t *testing.T) {
	q, chain, p2pnet, s, mem, exec, cs, wallet := initEnvPbft()
	defer chain.Close()
	defer mem.Close()
	defer p2pnet.Close()
	defer exec.Close()
	defer s.Close()
	defer cs.Close()
	defer q.Close()
	defer wallet.Close()
	time.Sleep(5 * time.Second)

	sendReplyList(q)
	clearTestData()
}

func initEnvPbft() (queue.Queue, *blockchain.BlockChain, *p2p.P2p, queue.Module, queue.Module, *executor.Executor, queue.Module, queue.Module) {
	var q = queue.New("channel")
	flag.Parse()
	cfg, sub := types.InitCfg("chain33.test.toml")
	types.Init(cfg.Title, cfg)
	chain := blockchain.New(cfg.BlockChain)
	chain.SetQueueClient(q.Client())
	mem := mempool.New(cfg.Mempool, nil)
	mem.SetQueueClient(q.Client())
	exec := executor.New(cfg.Exec, sub.Exec)
	exec.SetQueueClient(q.Client())
	types.SetMinFee(0)
	s := store.New(cfg.Store, sub.Store)
	s.SetQueueClient(q.Client())
	cs := NewPbft(cfg.Consensus, sub.Consensus["pbft"])
	cs.SetQueueClient(q.Client())
	p2pnet := p2p.New(cfg.P2P)
	p2pnet.SetQueueClient(q.Client())
	walletm := wallet.New(cfg.Wallet, sub.Wallet)
	walletm.SetQueueClient(q.Client())

	return q, chain, p2pnet, s, mem, exec, cs, walletm

}

func sendReplyList(q queue.Queue) {
	client := q.Client()
	client.Sub("mempool")
	var count int
	for msg := range client.Recv() {
		if msg.Ty == types.EventTxList {
			count++
			createReplyList("test" + strconv.Itoa(count))
			msg.Reply(client.NewMessage("consensus", types.EventReplyTxList,
				&types.ReplyTxList{Txs: transactions}))
			if count == 5 {
				time.Sleep(5 * time.Second)
				break
			}
		}
	}
}

func getprivkey(key string) crypto.PrivKey {
	cr, err := crypto.New(types.GetSignName("", types.SECP256K1))
	if err != nil {
		panic(err)
	}
	bkey, err := common.FromHex(key)
	if err != nil {
		panic(err)
	}
	priv, err := cr.PrivKeyFromBytes(bkey)
	if err != nil {
		panic(err)
	}
	return priv
}

func createReplyList(account string) {
	var result []*types.Transaction
	for j := 0; j < txSize; j++ {
		//tx := &types.Transaction{}
		val := &cty.CoinsAction_Transfer{Transfer: &types.AssetsTransfer{Amount: 10}}
		action := &cty.CoinsAction{Value: val, Ty: cty.CoinsActionTransfer}
		tx := &types.Transaction{Execer: []byte("coins"), Payload: types.Encode(action), Fee: 0}
		tx.To = "14qViLJfdGaP4EeHnDyJbEGQysnCpwn1gZ"

		tx.Nonce = random.Int63()

		tx.Sign(types.SECP256K1, getprivkey("CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944"))
		result = append(result, tx)
	}
	//result = append(result, tx)
	transactions = result
}

func clearTestData() {
	err := os.RemoveAll("datadir")
	if err != nil {
		fmt.Println("delete datadir have a err:", err.Error())
	}
	err = os.RemoveAll("wallet")
	if err != nil {
		fmt.Println("delete wallet have a err:", err.Error())
	}
	fmt.Println("test data clear successfully!")
}
