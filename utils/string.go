package utils

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	r "math/rand"
	"time"
)

// md5加密字符串
func Md5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

// 生成随机字符串
// keywords 为可选参数，只取[0]，指定随机数取值范围
func NewRandomString(n int, keywords ...string) string {
	ks := ""
	if len(keywords) > 0 {
		ks = keywords[0]
	}
	alphabets := []byte(ks)

	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	var randby bool
	if num, err := rand.Read(bytes); num != n || err != nil {
		r.Seed(time.Now().UnixNano())
		randby = true
	}
	for i, b := range bytes {
		if len(alphabets) == 0 {
			if randby {
				bytes[i] = alphanum[r.Intn(len(alphanum))]
			} else {
				bytes[i] = alphanum[b%byte(len(alphanum))]
			}
		} else {
			if randby {
				bytes[i] = alphabets[r.Intn(len(alphabets))]
			} else {
				bytes[i] = alphabets[b%byte(len(alphabets))]
			}
		}
	}
	return string(bytes)
}
