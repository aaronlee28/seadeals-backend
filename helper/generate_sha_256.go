package helper

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateHMACSHA256(value string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(value))
	buf := h.Sum(nil)
	sign := hex.EncodeToString(buf)
	return sign
}
