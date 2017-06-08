package utils

import (
	"crypto/aes"
	"crypto/rand"
	"io"
	"testing"
)

func TestEnctyptThenDecrypt(t *testing.T) {

	plainText := make([]byte, 1023)
	io.ReadFull(rand.Reader, plainText)

	encryptedBuf := make([]byte, 1100)
	decryptedBuf := make([]byte, 1100)

	encryptor := NewEncryptor("test Key")

	encryptedTextLength := encryptor.CFBEncrypter(plainText, encryptedBuf)
	plainTextLen := encryptor.CFBDecrypter(encryptedBuf[:encryptedTextLength], decryptedBuf)

	if string(decryptedBuf[:plainTextLen]) != string(plainText) {
		t.Error("string not match")
	}

	if encryptedTextLength != plainTextLen+aes.BlockSize {
		t.Error("encrypted size error")
	}
}
