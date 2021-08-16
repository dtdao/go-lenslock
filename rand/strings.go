package rand

import (
	"crypto/rand"

	"encoding/base64"
)

const RememberTokenBytes = 32

func Bytes(n int) ([]byte, error){
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// String Generate byte slice of size nBytes and then
// return a string that is the base64 url encoded version
// of that byte slice
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err !=nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func RememberToken() (string, error){
	return String(RememberTokenBytes)
}