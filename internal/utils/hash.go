package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func ComputeHash(fr io.Reader) (string, error) {
	hashing := md5.New()

	if _, err := io.Copy(hashing, fr); err != nil {
		return "", err
	}

	hashBytes := hashing.Sum(nil)
	return hex.EncodeToString(hashBytes), nil
}
