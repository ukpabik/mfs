package storage

import (
	"crypto/rand"
	"encoding/hex"
)

func generateServerID() string {
	buf := make([]byte, 32)
	rand.Read(buf)
	return hex.EncodeToString(buf)
}
