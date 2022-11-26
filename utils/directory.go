package utils

import (
	"errors"
	"os"
)

func PathExists(path string) (bool, error) {
	f, err := os.Stat(path)
	if err == nil {
		if f.IsDir() {
			return true, nil
		} else {
			return false, errors.New("存在同文件名")
		}
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
