package decode

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

func Sum256(v string) string {
	h := sha256.Sum256([]byte(v))
	return hex.EncodeToString(h[:])
}

func Sum512(v string) string {
	h := sha512.Sum512([]byte(v))
	return hex.EncodeToString(h[:])
}
