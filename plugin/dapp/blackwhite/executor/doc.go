// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

/*
第一个区块链的游戏：黑白配

玩法：

1. gen(s) 算法可以是: s = hash(hash(privkey)+nonce)
2. 发起：创建游戏，设定人数，设定金额，以及自动形成轮次，生成gameid
3. 竞猜:  gameid， 猜想的内容hash(s+黑)，lock 赌注 (match)，如果游戏实际人数小于设定人数则一直等待直到超时，超时时间为24*60*60s,如果超时则退还押金
4. 公布密钥:  公开s，超时时间为5分钟，在超时时间内所有参与人都公布密钥，则自动开奖；如果达到5分钟还有人未公布密钥，则系统自动开奖

约束条件：限制每局最多 20 BTY


status: create -> play -> show -> done(timeout done)


//对外查询接口
//1. 我的所有赌局，按照状态进行分类 （按照地址查询）
//2. 系统所有正在进行的赌局 (按照时间进行排序)
*/
