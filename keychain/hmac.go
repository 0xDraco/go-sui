package keychain

import (
	"crypto/hmac"
	"crypto/sha512"
)

func HMACSHA512(key, data []byte) []byte {
	mac := hmac.New(sha512.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}
