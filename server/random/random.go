package random

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomString() string {
	// Генерируем 32 случайных байта
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	// Кодируем байты в строку Base64
	return base64.URLEncoding.EncodeToString(b)
}
