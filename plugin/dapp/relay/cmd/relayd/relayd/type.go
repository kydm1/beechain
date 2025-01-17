// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package relayd

import (
	ty "github.com/33cn/plugin/plugin/dapp/relay/types"
)

var (
	currentBtcBlockheightKey = []byte("000007bc077cc540821f06280c4c2f04e036931f21ea3b1bf509b972cbbef5ca")
	firstHeaderKey           = []byte("1e1ffda1446d62a96b7ba9347a35ecdb1fad2379f35041d32e82a7b97646b9ff")
	privateKey               = []byte("privateKey-relayd")
)

type latestBlock struct {
	Hash       string   `json:"hash"`
	Time       int64    `json:"time"`
	BlockIndex uint64   `json:"block_index"`
	Height     uint64   `json:"height"`
	TxIndexes  []uint64 `json:"txIndexes"`
}

type header struct {
	Hash         string  `json:"hash"`
	Ver          uint64  `json:"ver"`
	PrevBlock    string  `json:"prev_block"`
	MerkleRoot   string  `json:"mrkl_root"`
	Time         int64   `json:"time"`
	Bits         int64   `json:"bits"`
	Fee          float64 `json:"fee,omitempty"`
	Nonce        int64   `json:"nonce"`
	TxNum        uint64  `json:"n_tx"`
	Size         uint64  `json:"size"`
	BlockIndex   uint64  `json:"block_index"`
	MainChain    bool    `json:"main_chain"`
	Height       uint64  `json:"height"`
	ReceivedTime int64   `json:"received_time"`
	RelayedBy    string  `json:"relayed_by"`
}

func (b *block) BtcHeader() *ty.BtcHeader {
	return &ty.BtcHeader{
		Hash:         b.Hash,
		Height:       b.Height,
		MerkleRoot:   b.MerkleRoot,
		Time:         b.Time,
		Nonce:        b.Nonce,
		Bits:         b.Bits,
		PreviousHash: b.PrevBlock,
		Version:      uint32(b.Ver),
	}
}

type block struct {
	// header
	Hash         string               `json:"hash"`
	Ver          uint64               `json:"ver"`
	PrevBlock    string               `json:"prev_block"`
	MerkleRoot   string               `json:"mrkl_root"`
	Time         int64                `json:"time"`
	Bits         int64                `json:"bits"`
	Nonce        uint64               `json:"nonce"`
	Fee          float64              `json:"fee,omitempty"`
	TxNum        uint64               `json:"n_tx"`
	Size         uint64               `json:"size"`
	BlockIndex   uint64               `json:"block_index"`
	MainChain    bool                 `json:"main_chain"`
	Height       uint64               `json:"height"`
	ReceivedTime int64                `json:"received_time"`
	RelayedBy    string               `json:"relayed_by"`
	Tx           []transactionDetails `json:"tx"`
}

type transactionDetails struct {
	LockTime  int64     `json:"lock_time"`
	Ver       uint64    `json:"ver"`
	Size      uint64    `json:"size"`
	Inputs    []txInput `json:"inputs"`
	Weight    int64     `json:"weight"`
	Time      int64     `json:"time"`
	TxIndex   uint64    `json:"tx_index"`
	VinSz     uint64    `json:"vin_sz"`
	Hash      string    `json:"hash"`
	VoutSz    uint64    `json:"vout_sz"`
	RelayedBy string    `json:"relayed_by"`
	Outs      []txOut   `json:"out"`
}

type txInput struct {
	Sequence int64  `json:"sequence"`
	Witness  string `json:"witness"`
	Script   string `json:"script"`
}

type txOut struct {
	Spent   bool   `json:"spent"`
	TxIndex uint64 `json:"tx_index"`
	Type    int    `json:"type"`
	Address string `json:"addr"`
	Value   uint64 `json:"value"`
	N       uint64 `json:"n"`
	Script  string `json:"script"`
}

type blocks struct {
	Blocks []block `json:"blocks"`
}

type transactionResult struct {
	Ver         uint     `json:"ver"`
	Inputs      []inputs `json:"inputs"`
	Weight      int64    `json:"weight"`
	BlockHeight uint64   `json:"block_height"`
	RelayedBy   string   `json:"relayed_by"`
	Out         []txOut  `json:"out"`
	LockTime    int64    `json:"lock_time"`
	Size        uint64   `json:"size"`
	DoubleSpend bool     `json:"double_spend"`
	Time        int64    `json:"time"`
	TxIndex     uint64   `json:"tx_index"`
	VinSz       uint64   `json:"vin_sz"`
	Hash        string   `json:"hash"`
	VoutSz      uint64   `json:"vout_sz"`
}

type inputs struct {
	Sequence uint   `json:"sequence"`
	Witness  string `json:"witness"`
	PrevOut  txOut  `json:"prev_out"`
	Script   string `json:"script"`
}

func (t *transactionResult) BtcTransaction() *ty.BtcTransaction {
	btcTx := &ty.BtcTransaction{}
	btcTx.Hash = t.Hash
	btcTx.Time = t.Time
	btcTx.BlockHeight = t.BlockHeight

	vin := make([]*ty.Vin, len(t.Inputs))
	for index, in := range t.Inputs {
		var v ty.Vin
		v.Value = in.PrevOut.Value
		v.Address = in.PrevOut.Address
		vin[index] = &v
	}
	btcTx.Vin = vin

	vout := make([]*ty.Vout, len(t.Out))
	for index, in := range t.Out {
		var out ty.Vout
		out.Value = in.Value
		out.Address = in.Address
		vout[index] = &out
		// TODO
		// vout[index].Coinbase
	}
	btcTx.Vout = vout
	return btcTx
}
