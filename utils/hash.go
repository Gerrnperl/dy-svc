package utils

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandString 生成随机字符串
//
//	@params length int 字符串长度
//	@return string
func RandString(length int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	sb := strings.Builder{}
	sb.Grow(length)
	for i := 0; i < length; i++ {
		sb.WriteByte(charset[rnd.Intn(len(charset))])
	}
	return sb.String()
}

// HashWithSalt 生成随机的盐并计算其附加到字符串后的Hash值
//
//	@param text
//	@return hash string Hash值
//	@return salt string 盐
func HashWithSalt(text string) (hash string, salt string) {
	salt = RandString(16)
	hash = Hash(text + salt)
	return
}

// Hash 计算字符串Hash值
//
//	@param text
//	@return string
func Hash(text string) string {
	return HashBytes([]byte(text))
}

// HashBytes 计算byte[]的Hash值
//
//	@param data
//	@return string
func HashBytes(data []byte) string {
	sum := md5.Sum(data)
	return fmt.Sprintf("%x", sum)
}

// CheckHash 检查字符串的Hash值是否与给定的Hash值相同
//
//	@param text string 待检查的字符串
//	@param salt string 盐
//	@param hash string 给定的Hash值
//	@return bool
func CheckHash(text string, salt string, hash string) bool {
	return Hash(text+salt) == hash
}
