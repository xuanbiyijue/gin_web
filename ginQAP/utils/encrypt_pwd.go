package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// secretPassword 登陆密码加密密钥
const secretPassword = "faiz555"

// EncryptPassword 密码加密
func EncryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secretPassword))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
