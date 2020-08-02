package urlhash

import (
	"encoding/hex"
	"fmt"

	"github.com/cespare/xxhash"
)

func Hash(str string) (string, error) {
	xh := xxhash.New()
	if n, err := xh.Write([]byte(str)); n == 0 || err != nil {
		return "", fmt.Errorf("unable to write bytes to hash: wrote %d: %s", n, err.Error())
	}

	return hex.EncodeToString(xh.Sum(nil)[0:4]), nil
}
