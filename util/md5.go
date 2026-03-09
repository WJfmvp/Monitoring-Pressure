package util

import (
	"crypto/md5"
	"fmt"
)

func Md5(data []byte) (result string) {
	mdSSum := md5.Sum(data)
	result = fmt.Sprintf("%x", mdSSum)
	return
}

func GetMD5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}
