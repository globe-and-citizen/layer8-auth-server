package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func GenerateRandomSalt(saltSize int) string {
	var salt = make([]byte, saltSize)

	_, err := rand.Read(salt[:])

	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(salt[:])
}

func XorBytes(bytesA, bytesB []byte) ([]byte, error) {
	if len(bytesA) != len(bytesB) {
		return nil, fmt.Errorf("slices must have the same length")
	}

	result := make([]byte, len(bytesA))
	for i := range bytesA {
		result[i] = bytesA[i] ^ bytesB[i]
	}
	return result, nil
}

func GenerateUUID() string {
	newUUID := uuid.New()
	return newUUID.String()
}

func GenerateSecret(secretSize int) string {
	var randomBytes = make([]byte, secretSize)

	_, err := rand.Read(randomBytes[:])
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(randomBytes[:])
}

func RemoveProtocolFromURL(url string) string {
	cleanedURL := strings.Replace(url, "http://", "", -1)
	cleanedURL = strings.Replace(cleanedURL, "https://", "", -1)
	return cleanedURL
}
