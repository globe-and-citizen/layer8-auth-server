package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
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
