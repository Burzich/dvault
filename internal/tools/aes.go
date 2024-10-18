package tools

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type AESEncryptor struct {
	aead cipher.AEAD
}

func NewAESEncryptor(secret []byte) (*AESEncryptor, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &AESEncryptor{
		aead: gcm,
	}, nil
}

func (a AESEncryptor) Encrypt(data []byte) ([]byte, error) {
	nonce := make([]byte, a.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := a.aead.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

func (a AESEncryptor) Decrypt(data []byte) ([]byte, error) {
	decryptedData, err := a.aead.Open(nil, data[:a.aead.NonceSize()], data[a.aead.NonceSize():], nil)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}
