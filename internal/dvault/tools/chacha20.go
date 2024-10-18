package tools

import (
	"crypto/cipher"
	"crypto/rand"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

type ChaCha struct {
	aead cipher.AEAD
}

func NewChaCha(secret []byte) (ChaCha, error) {
	aead, err := chacha20poly1305.New(secret)
	if err != nil {
		return ChaCha{}, err
	}

	return ChaCha{aead: aead}, nil
}

func (a ChaCha) Encrypt(data []byte) ([]byte, error) {
	nonce := make([]byte, a.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := a.aead.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

func (a ChaCha) Decrypt(data []byte) ([]byte, error) {
	decryptedData, err := a.aead.Open(nil, data[:a.aead.NonceSize()], data[a.aead.NonceSize():], nil)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}
