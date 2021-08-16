package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

// NewHMAC String creates and return hmac object
func NewHMAC(key string) HMAC{
	h := hmac.New(sha256.New, []byte(key))
	return HMAC{
		hmac: h,
	}}

// HMAC string wrapper around hash.Hash
type HMAC struct {
	hmac hash.Hash
}

func (h HMAC) Hash(input string) string {
	h.hmac.Reset()
	h.hmac.Write([]byte(input))

	b := h.hmac.Sum(nil)
	// urlencoding turn byte slice that might not have valid string,
	// and return a valid string and its url safe.
	return base64.URLEncoding.EncodeToString(b)
}

