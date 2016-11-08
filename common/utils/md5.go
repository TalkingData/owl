package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func GetFileMD5(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		return ""
	}
	hash := md5.New()
	io.Copy(hash, file)
	return fmt.Sprintf("%x", hash.Sum([]byte("")))
}

func Md5(c string) string {
	if len(c) == 0 {
		return ""
	}
	hash := md5.New()
	hash.Write([]byte(c))
	return hex.EncodeToString(hash.Sum(nil))
}
