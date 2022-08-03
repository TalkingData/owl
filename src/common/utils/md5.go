package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func GetFileMD5(filename string) (string, error) {
	fp, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = fp.Close()
	}()
	hash := md5.New()
	_, _ = io.Copy(hash, fp)
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func Md5(c string) string {
	if len(c) < 1 {
		return ""
	}
	hash := md5.New()
	hash.Write([]byte(c))
	return hex.EncodeToString(hash.Sum(nil))
}
