package helpers

import (
	"crypto/md5"
	"io"
	"log"
	"os"
)

func ComputeMd5(filePath string) ([]byte, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	var result []byte
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return result, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return result, err
	}

	return hash.Sum(result), nil
}
