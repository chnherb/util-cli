package utils

import (
	"os"
)

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	return !IsDir(path)
}

func Exist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func ExistFile(path string) bool {
	return Exist(path) && IsFile(path)
}

func ExistDir(path string) bool {
	return Exist(path) && IsDir(path)
}

func CreateDirIfNon(dir string) {
	if !ExistDir(dir) {
		_ = os.Mkdir(dir, os.ModePerm)
	}
}
