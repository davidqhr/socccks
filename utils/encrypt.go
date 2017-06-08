package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"io"
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

func (e *Encryptor) CFBDecrypter(ciphertext []byte) []byte {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext
}

func (e *Encryptor) CFBEncrypter(plaintext []byte) []byte {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext
}
