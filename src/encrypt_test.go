package socks5

import "testing"

var text = []byte("hello world!")

func mustPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Errorf("The code did not panic")
	}
}

func TestKeyShouldAtLeast8(t *testing.T) {
	defer mustPanic(t)

	key := []byte("1234567")
	formatKey(key)
}

func TestKeyShouldAtMost32(t *testing.T) {
	defer mustPanic(t)

	key := []byte("012345678901234567890123456789012")
	formatKey(key)
}

func TestKeyShouldComplement(t *testing.T) {
	initDefaultEncryptKey()

	key := []byte("a test key")
	formatedKey := formatKey(key)

	if string(formatedKey) != string(append(key, defaultKey[:32-len(key)]...)) {
		t.Error("key not complement")
	}
}

func TestEnctyptThenDecrypt(t *testing.T) {
	initDefaultEncryptKey()

	encryptedText := encryptAES(text, []byte("12345678"))
	decryptedText := decryptAES(encryptedText, []byte("12345678"))
	if string(decryptedText) != string(text) {
		t.Error(string(decryptedText), "!=", string(text))
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	initDefaultEncryptKey()

	encryptedText := encryptAES(text, []byte("12345678"))
	decryptedText := decryptAES(encryptedText, []byte("abcdefghijklhmln"))

	if string(decryptedText) == string(text) {
		t.Error("wrong key but right result")
	}
}
