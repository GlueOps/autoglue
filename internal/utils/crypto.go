package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
)

func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

func encryptAESGCM(plain, key []byte) (ciphertext, iv, tag []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	iv, err = randomBytes(12)
	if err != nil {
		return
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	cipherData := aesgcm.Seal(nil, iv, plain, nil)
	ciphertext = cipherData[:len(cipherData)-aesgcm.Overhead()]
	tag = cipherData[len(cipherData)-aesgcm.Overhead():]
	return
}

func decryptAESGCM(ciphertext, key, iv, tag []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	fullCipher := append(ciphertext, tag...)
	return aesgcm.Open(nil, iv, fullCipher, nil)
}

func decodeB64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func encodeB64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
