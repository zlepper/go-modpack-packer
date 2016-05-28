package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"log"
	"path"
	"os"
	"io/ioutil"
)
var commonIV []byte
var cypher cipher.Block

func init() {
	// Read bytes from file
	keyfile := path.Join(os.Args[1], "key")
	var bytes []byte
	bytes, err := ioutil.ReadFile(keyfile)
	if err == nil {
		getCommonIV(bytes)
		return
	}

	bytes = make([]byte, 32)
	_, err = rand.Read(bytes)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	ioutil.WriteFile(keyfile, bytes, os.ModePerm)
	getCommonIV(bytes)

}

func getCommonIV(key []byte) {
	ivFile := path.Join(os.Args[1], "iv")
	commonIV, err := ioutil.ReadFile(ivFile)
	if err == nil {
		createCypher(key)
		return
	}
	commonIV = make([]byte, 16)
	_, err = rand.Read(commonIV)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	ioutil.WriteFile(ivFile, commonIV, os.ModePerm)
	createCypher(key)
}


func createCypher(key []byte) {
	cypher, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("%v", err)
		return
	}
}

func Encrypt(plain []byte) []byte {
	cfb := cipher.NewCFBEncrypter(cypher, commonIV)
	encryptedText := make([]byte, len(plain))
	cfb.XORKeyStream(encryptedText, plain)
	return encryptedText
}

func Decrypt(encrypted []byte) []byte {
	cfbdec := cipher.NewCFBDecrypter(cypher, commonIV)
	plain := make([]byte, len(encrypted))
	cfbdec.XORKeyStream(plain, encrypted)
	return plain
}
