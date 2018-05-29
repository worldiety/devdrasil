package rand

import (
	"crypto/rand"
	"encoding/hex"
)

func SecureRandomHex(len int) string {
	tmp := make([]byte, len)
	n, e := rand.Read(tmp)
	if e != nil || n != len {
		panic(e)
	}
	return hex.EncodeToString(tmp)
}
