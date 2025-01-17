// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/33cn/chain33/types"
	pkt "github.com/33cn/plugin/plugin/dapp/pokerbull/types"
)

func (c *PokerBull) updateIndex(log *pkt.ReceiptPBGame) (kvs []*types.KeyValue) {
	//先保存本次Action产生的索引
	kvs = append(kvs, addPBGameStatusAndPlayer(log.Status, log.PlayerNum, log.Value, log.Index, log.GameId))
	kvs = append(kvs, addPBGameStatusIndexKey(log.Status, log.GameId, log.Index))
	kvs = append(kvs, addPBGameAddrIndexKey(log.Status, log.Addr, log.GameId, log.Index))

	/*
		//状态更新
		if log.Status == pkt.PBGameActionStart {
			kvs = append(kvs, delPBGameStatusAndPlayer(pkt.PBGameActionStart, log.PlayerNum, log.Value, log.PrevIndex))
			kvs = append(kvs, delPBGameStatusIndexKey(pkt.PBGameActionStart, log.PrevIndex))
		}

		if log.Status == pkt.PBGameActionContinue {
			kvs = append(kvs, delPBGameStatusAndPlayer(pkt.PBGameActionStart, log.PlayerNum, log.Value, log.PrevIndex))
			kvs = append(kvs, delPBGameStatusAndPlayer(pkt.PBGameActionContinue, log.PlayerNum, log.Value, log.PrevIndex))
			kvs = append(kvs, delPBGameStatusIndexKey(pkt.PBGameActionStart, log.PrevIndex))
			kvs = append(kvs, delPBGameStatusIndexKey(pkt.PBGameActionContinue, log.PrevIndex))
		}

		if log.Status == pkt.PBGameActionQuit {
			kvs = append(kvs, delPBGameStatusAndPlayer(pkt.PBGameActionStart, log.PlayerNum, log.Value, log.PrevIndex))
			kvs = append(kvs, delPBGameStatusAndPlayer(pkt.PBGameActionContinue, log.PlayerNum, log.Value, log.PrevIndex))
			kvs = append(kvs, delPBGameStatusIndexKey(pkt.PBGameActionStart, log.PrevIndex))
			kvs = append(kvs, delPBGameStatusIndexKey(pkt.PBGameActionContinue, log.PrevIndex))

		}*/

	//结束一局，更新所有玩家地址状态
	if !log.IsWaiting {
		for _, v := range log.Players {
			if v != log.Addr {
				kvs = append(kvs, addPBGameAddrIndexKey(log.Status, v, log.GameId, log.Index))
			}
			kvs = append(kvs, delPBGameAddrIndexKey(v, log.PrevIndex))
		}

		kvs = append(kvs, delPBGameStatusAndPlayer(log.PreStatus, log.PlayerNum, log.Value, log.PrevIndex))
		kvs = append(kvs, delPBGameStatusIndexKey(log.PreStatus, log.PrevIndex))
	}

	return kvs
}

func (c *PokerBull) execLocal(receipt *types.ReceiptData) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	for i := 0; i < len(receipt.Logs); i++ {
		item := receipt.Logs[i]
		if item.Ty == pkt.TyLogPBGameStart || item.Ty == pkt.TyLogPBGameContinue || item.Ty == pkt.TyLogPBGameQuit || item.Ty == pkt.TyLogPBGamePlay {
			var Gamelog pkt.ReceiptPBGame
			err := types.Decode(item.Log, &Gamelog)
			if err != nil {
				panic(err) //数据错误了，已经被修改了
			}
			kv := c.updateIndex(&Gamelog)
			dbSet.KV = append(dbSet.KV, kv...)
		}
	}
	return dbSet, nil
}

// ExecLocal_Start 开始游戏交易local执行
func (c *PokerBull) ExecLocal_Start(payload *pkt.PBGameStart, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return c.execLocal(receiptData)
}

// ExecLocal_Continue 继续游戏交易local执行
func (c *PokerBull) ExecLocal_Continue(payload *pkt.PBGameContinue, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return c.execLocal(receiptData)
}

// ExecLocal_Quit 退出游戏交易local执行
func (c *PokerBull) ExecLocal_Quit(payload *pkt.PBGameQuit, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return c.execLocal(receiptData)
}

// ExecLocal_Play 已匹配游戏交易local执行
func (c *PokerBull) ExecLocal_Play(payload *pkt.PBGamePlay, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return c.execLocal(receiptData)
}
