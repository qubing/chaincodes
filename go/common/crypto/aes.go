package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"sync"
)

// AESHelper - AESHelper structure definition.
type AESHelper struct {
	syncMutex sync.Mutex
	key       []byte // = []byte("1234567890abcdef")
}

// NewAESHelper - geneate an instance of XAES
func NewAESHelper(key string) (*AESHelper, error) {
	if !(len(key) == 16 || len(key) == 24 || len(key) == 32) {
		return nil, errors.New("key length not supported. (only 16, 24 or 32 is supported)")
	}

	xaes := &AESHelper{}
	xaes.syncMutex.Lock()
	defer xaes.syncMutex.Unlock()
	xaes.key = []byte(key)
	return xaes, nil
}

// Encrypt - encrypt data
func (xaes *AESHelper) Encrypt(data string) (string, error) {
	block, err := aes.NewCipher(xaes.key)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cipher.NewCFBEncrypter(block, iv).XORKeyStream(ciphertext[aes.BlockSize:],
		[]byte(data))
	return hex.EncodeToString(ciphertext), nil
}

// Decrypt - decrypt data
func (xaes *AESHelper) Decrypt(d string) (string, error) {
	ciphertext, err := hex.DecodeString(d)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(xaes.key)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	cipher.NewCFBDecrypter(block, iv).XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}
