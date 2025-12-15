package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

func GenerateTelegramSessionID() ([]byte, error) {
	const entropyBytes = 32
	b := make([]byte, entropyBytes)
	if _, err := rand.Read(b); err != nil {
		return []byte{}, fmt.Errorf("read random: %w", err)
	}

	return b, nil
}

func ComputeTelegramSessionIDHash(sessionID []byte) [32]byte {
	return sha256.Sum256(sessionID)
}
