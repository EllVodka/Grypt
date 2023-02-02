package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/pkg/errors"
)

// Encode slice of byte in string
func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// Decode string to slice of byte
func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

// Encrypt method is to encrypt or hide any classified text
func Encrypt(text, MySecret string, b []byte) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, b)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

// Decrypt method is to extract back the encrypted text
func Decrypt(text, MySecret string, b []byte) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", errors.Wrap(err, "Decrypt() - cipher error")
	}
	cipherText := Decode(text)
	cfb := cipher.NewCFBDecrypter(block, b)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
