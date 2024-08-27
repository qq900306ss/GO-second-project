package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// 轉小寫並產生MD5值
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tmepStr := h.Sum(nil)
	return hex.EncodeToString(tmepStr)

}

// 大寫
func MD5Encode(data string) string {

	return strings.ToUpper(Md5Encode(data))
}

// 加密操作
func MakePassword(plainpwd, sweet string) string { //傳過來的密碼配上我sweet 再做MD5加密 產生密碼
	return Md5Encode(plainpwd + sweet)
}

// 解密後檢查密碼
func ValidPassword(plainpwd, sweet string, password string) bool { //傳過來的密碼配上我sweet 再做MD5加密 產生密碼
	return Md5Encode(plainpwd+sweet) == password
}
