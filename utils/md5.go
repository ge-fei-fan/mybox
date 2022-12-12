package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

//md5加密
func MD5V(str []byte, b ...byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(b))
}

func MD5(name string, b ...byte) (string, error) {
	fd, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return "", err
	}
	defer fd.Close()
	h := md5.New()
	_, _ = io.Copy(h, fd)
	return hex.EncodeToString(h.Sum(b)), nil
}
