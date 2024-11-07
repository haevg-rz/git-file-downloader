package logic

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

func GetDirFromFilepath(file string) (bool, string) {
	dir := filepath.Dir(file)
	_, err := os.Stat(dir)
	if err == nil {
		return true, dir
	}
	return !os.IsNotExist(err), dir
}

func FileExists(file string) bool {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}

func IsValidPath(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return !os.IsNotExist(err)
}

func IsHashEqual(file, compareHash string) (bool, error) {
	if _, err := os.Stat(file); err != nil {
		return false, err
	}

	f, err := os.Open(file)
	if err != nil {
		return false, err
	}

	defer func() {
		err = f.Close()
	}()

	hash := sha256.New()
	if _, err = io.Copy(hash, f); err != nil {
		return false, err
	}

	return hex.EncodeToString(hash.Sum(nil)) == compareHash, err
}
