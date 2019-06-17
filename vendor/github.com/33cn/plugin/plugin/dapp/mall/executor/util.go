package executor

import (
	"fmt"
)


var (
	mallPrefix = "mavl-mall-"
)

func mallKeyUser(id string) (key []byte) {
	return []byte(fmt.Sprintf(mallPrefix+"user-"+"%s", id))
}

func mallKeyPlatform() (key []byte) {
	return []byte(mallPrefix+"platform")
}

func mallKeyGoodInfo(id string) (key []byte) {
	return []byte(fmt.Sprintf(mallPrefix+"good-"+"%s",id))
}

