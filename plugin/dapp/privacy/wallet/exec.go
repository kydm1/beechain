// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wallet

import (
	"github.com/33cn/chain33/types"
	privacytypes "github.com/33cn/plugin/plugin/dapp/privacy/types"
)

func (policy *privacyPolicy) On_ShowPrivacyAccountSpend(req *privacytypes.ReqPrivBal4AddrToken) (types.Message, error) {
	policy.getWalletOperate().GetMutex().Lock()
	defer policy.getWalletOperate().GetMutex().Unlock()
	reply, err := policy.showPrivacyAccountsSpend(req)
	if err != nil {
		bizlog.Error("showPrivacyAccountsSpend", "err", err.Error())
	}
	return reply, err
}

func (policy *privacyPolicy) On_ShowPrivacyKey(req *types.ReqString) (types.Message, error) {
	policy.getWalletOperate().GetMutex().Lock()
	defer policy.getWalletOperate().GetMutex().Unlock()
	reply, err := policy.showPrivacyKeyPair(req)
	if err != nil {
		bizlog.Error("showPrivacyKeyPair", "err", err.Error())
	}
	return reply, err
}

func (policy *privacyPolicy) On_CreateUTXOs(req *privacytypes.ReqCreateUTXOs) (types.Message, error) {
	policy.getWalletOperate().GetMutex().Lock()
	defer policy.getWalletOperate().GetMutex().Unlock()
	reply, err := policy.createUTXOs(req)
	if err != nil {
		bizlog.Error("createUTXOs", "err", err.Error())
	}
	return reply, err
}

func (policy *privacyPolicy) On_CreateTransaction(req *privacytypes.ReqCreatePrivacyTx) (types.Message, error) {
	ok, err := policy.getWalletOperate().CheckWalletStatus()
	if !ok {
		bizlog.Error("createTransaction", "CheckWalletStatus cause error.", err)
		return nil, err
	}
	if ok, err := policy.isRescanUtxosFlagScaning(); ok {
		bizlog.Error("createTransaction", "isRescanUtxosFlagScaning cause error.", err)
		return nil, err
	}
	if !checkAmountValid(req.Amount) {
		err = types.ErrAmount
		bizlog.Error("createTransaction", "isRescanUtxosFlagScaning cause error.", err)
		return nil, err
	}
	policy.getWalletOperate().GetMutex().Lock()
	defer policy.getWalletOperate().GetMutex().Unlock()

	reply, err := policy.createTransaction(req)
	if err != nil {
		bizlog.Error("createTransaction", "err", err.Error())
	}
	return reply, err
}

func (policy *privacyPolicy) On_ShowPrivacyAccountInfo(req *privacytypes.ReqPPrivacyAccount) (types.Message, error) {
	policy.getWalletOperate().GetMutex().Lock()
	defer policy.getWalletOperate().GetMutex().Unlock()
	reply, err := policy.getPrivacyAccountInfo(req)
	if err != nil {
		bizlog.Error("getPrivacyAccountInfo", "err", err.Error())
	}
	return reply, err
}

func (policy *privacyPolicy) On_PrivacyTransactionList(req *privacytypes.ReqPrivacyTransactionList) (types.Message, error) {
	if req.Direction != 0 && req.Direction != 1 {
		bizlog.Error("getPrivacyTransactionList", "invalid direction ", req.Direction)
		return nil, types.ErrInvalidParam
	}
	// convert to sendTx / recvTx
	sendRecvFlag := req.SendRecvFlag + sendTx
	if sendRecvFlag != sendTx && sendRecvFlag != recvTx {
		bizlog.Error("getPrivacyTransactionList", "invalid sendrecvflag ", req.SendRecvFlag)
		return nil, types.ErrInvalidParam
	}
	req.SendRecvFlag = sendRecvFlag

	policy.getWalletOperate().GetMutex().Lock()
	defer policy.getWalletOperate().GetMutex().Unlock()

	reply, err := policy.store.getWalletPrivacyTxDetails(req)
	if err != nil {
		bizlog.Error("getWalletPrivacyTxDetails", "err", err.Error())
	}
	return reply, err
}

func (policy *privacyPolicy) On_RescanUtxos(req *privacytypes.ReqRescanUtxos) (types.Message, error) {
	policy.getWalletOperate().GetMutex().Lock()
	defer policy.getWalletOperate().GetMutex().Unlock()
	reply, err := policy.rescanUTXOs(req)
	if err != nil {
		bizlog.Error("rescanUTXOs", "err", err.Error())
	}
	return reply, err
}

func (policy *privacyPolicy) On_EnablePrivacy(req *privacytypes.ReqEnablePrivacy) (types.Message, error) {
	operater := policy.getWalletOperate()
	operater.GetMutex().Lock()
	defer operater.GetMutex().Unlock()
	reply, err := policy.enablePrivacy(req)
	if err != nil {
		bizlog.Error("enablePrivacy", "err", err.Error())
	}
	return reply, err
}
