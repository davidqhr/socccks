package socks5

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

var defaultKey = make([]byte, 32)

func initDefaultEncryptKey() {
	rand.Read(defaultKey)
}

// make key valid
func formatKey(rawKey []byte) (key []byte) {
	keyLen := len(rawKey)
	if keyLen < 8 {
		panic("key lenght at least 8")
	} else if keyLen > 32 {
		panic("key lenght at most 32")
	} else if len(rawKey) < 32 {
		key = make([]byte, 32)
		copy(key, rawKey)
		copy(key[keyLen:], defaultKey)
	}
	return
}

func decryptAES(cipherBytes []byte, key []byte) []byte {
	key = formatKey(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	iv := cipherBytes[:aes.BlockSize]
	cipherBytes = cipherBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherBytes, cipherBytes)
	return cipherBytes
}

func encryptAES(rawText []byte, key []byte) []byte {
	key = formatKey(key)

	block, err := aes.NewCipher(key)

	if err != nil {
		panic(err)
	}

	cipherBytes := make([]byte, aes.BlockSize+len(rawText))
	iv := cipherBytes[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherBytes[aes.BlockSize:], rawText)

	return cipherBytes
}
