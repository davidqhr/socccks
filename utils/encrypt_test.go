package utils

import (
	"crypto/rand"
	"io"
	"testing"
)

func TestEnctyptThenDecrypt(t *testing.T) {

	plainText := make([]byte, 1023)
	io.ReadFull(rand.Reader, plainText)

	encryptor := NewEncryptor("test Key")

	encryptedText := encryptor.CFBEncrypter(plainText)
	decryptedText := encryptor.CFBDecrypter(encryptedText)

	if string(decryptedText) != string(plainText) {
		t.Error(string(decryptedText), "!=", string(plainText))
	}
}
