// package cypher - функции (де)/шифрования
package cypher

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/Azcarot/PasswordStorage/internal/storage"
)

// ShaData - хеширование данных пользователя
func ShaData(result string, key string) string {
	b := []byte(result)
	shakey := []byte(key)
	// создаём новый hash.Hash, вычисляющий контрольную сумму SHA-256
	h := hmac.New(sha256.New, shakey)
	// передаём байты для хеширования
	h.Write(b)
	// вычисляем хеш
	hash := h.Sum(nil)
	sha := base64.URLEncoding.EncodeToString(hash)
	return string(sha)
}

// CypherData - шифрование данных секретом из контекста
func CypherData(ctx context.Context, data string) (string, error) {

	block, err := aes.NewCipher(storage.Secret[:])
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(data))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Dechyper - дешифровка данных секретом из контекста
func Dechypher(ctx context.Context, data string) (string, error) {

	ciphertextBytes, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(storage.Secret[:])

	if len(ciphertextBytes) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext слишком короткий")
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return string(ciphertextBytes), nil
}
