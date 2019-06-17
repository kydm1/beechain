// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"fmt"
)

var (
	tokennoteCreated          = "mavl-tokennote-"
	tokennotePreCreatedOT     = "mavl-create-tokennote-ot-"
	tokennotePreCreatedSTO    = "mavl-create-tokennote-sto-"
	tokennotePreCreatedOTNew  = "mavl-tokennote-create-ot-"
	tokennotePreCreatedSTONew = "mavl-tokennote-create-sto-"
	tokennoteHoldKey     = "mavl-tokennote-hold-"

	tokennotePreCreatedSTONewLocal = "LODB-tokennote-create-sto-"
	tokennoteHoldLocalKey     = "LODB-tokennote-hold-"
	tokennoteCashedLocalKey   = "LODB-tokennote-cashed-"
	tokennoteLocalPre		  = "LODB-tokennote"

	tokennoteContract     = "mavl-tokennote-contract-" //白条合同查询验证  合同不上链
)

func calcTokennoteKey(tokennote string) (key []byte) {
	return []byte(fmt.Sprintf(tokennoteCreated+"%s", tokennote))
}

func calcTokennoteStatusKey(tokennote string, owner string, status int32) []byte {
	return []byte(fmt.Sprintf(tokennotePreCreatedSTO+"%d-%s-%s", status, tokennote, owner))
}

func calcTokennoteAddrNewKeyS(tokennote string, owner string) (key []byte) {
	return []byte(fmt.Sprintf(tokennotePreCreatedOTNew+"%s-%s", owner, tokennote))
}

func calcTokennoteStatusNewKeyS(tokennote string, owner string, status int32) []byte {
	return []byte(fmt.Sprintf(tokennotePreCreatedSTONew+"%d-%s-%s", status, tokennote, owner))
}

func calcTokennoteHoldKey(tokennote string, addr string ,time int64) []byte {
	return []byte(fmt.Sprintf(tokennoteHoldKey+"%s-%s-%d",  tokennote,addr,time))
}

func calcTokennoteHoldKeyNew(tokennote string, addr string) []byte {
	return []byte(fmt.Sprintf(tokennoteHoldKey+"%s-%s",  tokennote,addr))
}

func calcTokennoteLocalHoldKey(tokennote string, addr string ,time int64) []byte {
	return []byte(fmt.Sprintf(tokennoteHoldLocalKey+"%s-%s-%d",  tokennote,addr,time))
}

func calcTokennoteStatusKeyLocal(tokennote string, status int32) []byte {
	return []byte(fmt.Sprintf(tokennotePreCreatedSTONewLocal+"%d-%s", status, tokennote))
}

func calcTokennoteStatusKeyPrefixLocal(status int32) []byte {
	return []byte(fmt.Sprintf(tokennotePreCreatedSTONewLocal+"%d", status))
}

func calcTokennoteStatusTokennoteKeyPrefixLocal(status int32, tokennote string) []byte {
	return []byte(fmt.Sprintf(tokennotePreCreatedSTONewLocal+"%d-%s-", status, tokennote))
}

//存储地址上收币的信息
func calcAddrKey(tokennote string, addr string) []byte {
	return []byte(fmt.Sprintf("LODB-tokennote-%s-Addr:%s", tokennote, addr))
}

//白条合同验证相关信息
func calcTokennoteContractKey(tokennote string) []byte {
	return []byte(fmt.Sprintf(tokennoteContract+"%s-",tokennote))
}

//func calcTokennoteCashedLocal(tokennote string,addr string,to string,heightindex string ) []byte {
//	return []byte(fmt.Sprintf(tokennoteCashedLocalKey+"%s:%s:%s:%s",addr,tokennote,to,heightindex))
//}
