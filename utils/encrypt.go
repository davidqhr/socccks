package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"io"
	"log"
)

type Encryptor struct {
	rawKey string
	key    []byte
}

func NewEncryptor(rawKey string) *Encryptor {
	return &Encryptor{
		rawKey: rawKey,
		key:    formatKey([]byte(rawKey)),
	}
}

// return 32 byte NewCipher key
func formatKey(rawKey []byte) (key []byte) {
	// 16 bytes
	MD5Summary16 := md5.Sum(rawKey)
	MD5Summary := MD5Summary16[:]

	// 16 bytes
	MD5SummarySummary := md5.Sum(MD5Summary)

	return append(MD5Summary, MD5SummarySummary[:]...)
}

func (e *Encryptor) CFBDecrypter(ciphertext []byte, buf []byte) int {
	block, err := aes.NewCipher(e.key)

	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	dataLength := len(ciphertext)
	if dataLength > len(buf) {
		log.Fatalln("CFBDecrypter buf is too small", dataLength, len(buf))
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(buf, ciphertext)

	return dataLength
}

func (e *Encryptor) CFBEncrypter(plaintext []byte, buf []byte) int {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		panic(err)
	}
	textLength := len(plaintext)

	if textLength+aes.BlockSize > len(buf) {
		log.Fatalln("CFBEncrypter buf is too small", textLength+aes.BlockSize, len(buf))
	}

	iv := buf[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(buf[aes.BlockSize:], plaintext)

	return textLength + aes.BlockSize
}
