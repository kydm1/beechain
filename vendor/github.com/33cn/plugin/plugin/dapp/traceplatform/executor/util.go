package executor

import (
	"fmt"
)


var (
	traceplatformPrefix = "mavl-traceplatform-"
	traceplatformHistory = "LODB-traceplatform-nfccodehistory:"
	traceplatformTxsCount = "LODB-traceplatform-txscount:"
)

func traceplatformKeyGood(id string) (key []byte) {
	return []byte(fmt.Sprintf(traceplatformPrefix+"good-"+"%s", id))
}

func traceplatformKeyNFCCode(addr, heightindex string) (key []byte) {
	return []byte(fmt.Sprintf(traceplatformHistory+"%s:%s",addr,heightindex))
}

func traceplatformKeyNFCCodeHistory(nfc string) (key []byte) {
	return []byte(fmt.Sprintf(traceplatformTxsCount+"%s",nfc))
}