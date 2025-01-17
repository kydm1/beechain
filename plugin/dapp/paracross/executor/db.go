// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"encoding/hex"

	"github.com/33cn/chain33/client"
	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/types"
	pt "github.com/33cn/plugin/plugin/dapp/paracross/types"
)

func getTitle(db dbm.KV, key []byte) (*pt.ParacrossStatus, error) {
	val, err := db.Get(key)
	if err != nil {
		if !isNotFound(err) {
			return nil, err
		}
		// 平行链如果是从其他链上移过来的，  需要增加配置， 对应title的平行链的起始高度
		clog.Info("first time load title", "key", string(key))
		return &pt.ParacrossStatus{Height: -1}, nil
	}

	var title pt.ParacrossStatus
	err = types.Decode(val, &title)
	return &title, err
}

func saveTitle(db dbm.KV, key []byte, title *pt.ParacrossStatus) error {
	val := types.Encode(title)
	return db.Set(key, val)
}

func getTitleHeight(db dbm.KV, key []byte) (*pt.ParacrossHeightStatus, error) {
	val, err := db.Get(key)
	if err != nil {
		// 对应高度第一次提交commit
		if isNotFound(err) {
			clog.Info("paracross.Commit first commit", "key", string(key))
		}
		return nil, err
	}
	var heightStatus pt.ParacrossHeightStatus
	err = types.Decode(val, &heightStatus)
	return &heightStatus, err
}

func saveTitleHeight(db dbm.KV, key []byte, heightStatus types.Message /* heightStatus *types.ParacrossHeightStatus*/) error {
	// use as a types.Message
	val := types.Encode(heightStatus)
	return db.Set(key, val)
}

//GetBlock get block detail by block hash
func GetBlock(api client.QueueProtocolAPI, blockHash []byte) (*types.BlockDetail, error) {
	blockDetails, err := api.GetBlockByHashes(&types.ReqHashes{Hashes: [][]byte{blockHash}})
	if err != nil {
		clog.Error("paracross.Commit getBlockHeader", "db", err,
			"commit tx hash", hex.EncodeToString(blockHash))
		return nil, err
	}
	if len(blockDetails.Items) != 1 {
		clog.Error("paracross.Commit getBlockHeader", "len in not 1", len(blockDetails.Items))
		return nil, pt.ErrParaBlockHashNoMatch
	}
	if blockDetails.Items[0] == nil {
		clog.Error("paracross.Commit getBlockHeader", "commit tx hash net found", hex.EncodeToString(blockHash))
		return nil, pt.ErrParaBlockHashNoMatch
	}
	return blockDetails.Items[0], nil
}

func getBlockHash(api client.QueueProtocolAPI, height int64) (*types.ReplyHash, error) {
	hash, err := api.GetBlockHash(&types.ReqInt{Height: height})
	if err != nil {
		clog.Error("paracross.Commit getBlockHeader", "db", err,
			"commit height", height)
		return nil, err
	}
	return hash, nil
}

func isNotFound(err error) bool {
	if err != nil && (err == dbm.ErrNotFoundInDb || err == types.ErrNotFound) {
		return true
	}
	return false
}

//GetTx get tx by tx hash
func GetTx(api client.QueueProtocolAPI, txHash []byte) (*types.TransactionDetail, error) {
	txs, err := api.GetTransactionByHash(&types.ReqHashes{Hashes: [][]byte{txHash}})
	if err != nil {
		clog.Error("paracross.Commit GetTx", "db", err,
			"commit tx hash", hex.EncodeToString(txHash))
		return nil, err
	}
	if len(txs.Txs) != 1 {
		clog.Error("paracross.Commit GetTx", "len in not 1", len(txs.Txs))
		return nil, pt.ErrParaBlockHashNoMatch
	}
	if txs.Txs == nil {
		clog.Error("paracross.Commit GetTx", "commit tx hash net found", hex.EncodeToString(txHash))
		return nil, pt.ErrParaBlockHashNoMatch
	}
	return txs.Txs[0], nil
}
