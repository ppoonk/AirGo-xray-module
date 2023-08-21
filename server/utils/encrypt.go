package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strings"
	"time"
)

// BcryptHash 使用 bcrypt 对密码进行加密
func BcryptEncode(password string) string {
	// Go 中的 bcrypt.DefaultCost 是 10
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// BcryptCheck 对比明文密码和数据库的哈希值
func BcryptDecode(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// md5 encode
func Md5Encode(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	res := hex.EncodeToString(hash.Sum(nil))
	//转大写，strings.ToUpper(res)
	return res
}

// sha256 encode
func Sha256Encode(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	res := hex.EncodeToString(hash.Sum(nil))
	return res
}

// 随机数，n为 位数
func RandomString(n int64) string {
	var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	randomStr := make([]rune, n)
	for i := range randomStr {
		randomStr[i] = defaultLetters[rand.Intn(len(defaultLetters))]
	}
	return string(randomStr)
}

// base64 解码
func Base64Decode(s string) (string, error) {
	res, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// 订阅base64 解码
func SubBase64Decode(str string) string {
	i := len(str) % 4
	switch i {
	case 1:
		str = str[:len(str)-1]
	case 2:
		str += "=="
	case 3:
		str += "="
	}
	//str = strings.Split(str, "//")[1]
	var data []byte
	var err error
	if strings.Contains(str, "-") || strings.Contains(str, "_") {
		data, err = base64.URLEncoding.DecodeString(str)

	} else {
		data, err = base64.StdEncoding.DecodeString(str)
		//data, err = base64.RawURLEncoding.DecodeString(str)
	}
	if err != nil {
		fmt.Println(err)
	}
	return string(data)
}
