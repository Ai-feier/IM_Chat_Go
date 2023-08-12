package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// Md5Encoder 小写密文
func Md5Encoder(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

// MD5Encoder 大写密文
func MD5Encoder(s string) string {
	return strings.ToUpper(MD5Encoder(s))
}

// MakePassword 加密
func MakePassword(password, salt string) string {
	return Md5Encoder(password + salt)
}

// ValidPassword 解密
func ValidPassword(newpass, truepass, salt string) bool {
	return MakePassword(newpass, salt) == truepass
}
