package security

import (
	"crypto/sha256"
	"encoding/hex"
)

type SHA256CodeHasher struct{}

func NewSHA256CodeHasher() *SHA256CodeHasher {
	return &SHA256CodeHasher{}
}

func (h *SHA256CodeHasher) Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
