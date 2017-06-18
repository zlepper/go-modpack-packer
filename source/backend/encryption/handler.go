package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/zlepper/go-modpack-packer/source/backend/consts"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// NOTE: Don't ever ask me to explain how this works, i barely understand it myself.
// Derived from this answer on SO: http://stackoverflow.com/a/18819040/3950006
var key []byte

func init() {
	// Read bytes from file
	keyfile := path.Join(consts.DataDirectory, "key")
	var err error
	key, err = ioutil.ReadFile(keyfile)
	if err == nil {
		return
	}

	key = make([]byte, 32)
	_, err = rand.Read(key)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	ioutil.WriteFile(keyfile, key, os.ModePerm)
}

func encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func Encrypt(text []byte) []byte {
	b, err := encrypt(key, text)
	if err != nil {
		log.Panic(err)
	}
	return b
}

func EncryptString(text string) string {
	return fmt.Sprintf("%0x", Encrypt([]byte(text)))
}

func decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func Decrypt(text []byte) []byte {
	b, err := decrypt(key, text)
	if err != nil {
		log.Panic(err)
	}
	return b
}

func DecryptString(text string) string {
	by, err := hex.DecodeString(text)
	if err != nil {
		log.Panic(err)
	}
	return string(Decrypt(by))
}
