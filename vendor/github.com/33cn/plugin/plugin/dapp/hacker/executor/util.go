package executor

import (
	"fmt"
)


var (
	hackerPrefix = "mavl-hacker-"
	hackerHistory = "LODB-hacker-nfccodehistory:"
	hackerTxsCount = "LODB-hacker-txscount:"
)

func hackerKeyGood(id string) (key []byte) {
	return []byte(fmt.Sprintf(hackerPrefix+"good-"+"%s", id))
}

func hackerKeyNFCCode(addr, heightindex string) (key []byte) {
	return []byte(fmt.Sprintf(hackerHistory+"%s:%s",addr,heightindex))
}

func hackerKeyNFCCodeHistory(nfc string) (key []byte) {
	return []byte(fmt.Sprintf(hackerTxsCount+"%s",nfc))
}