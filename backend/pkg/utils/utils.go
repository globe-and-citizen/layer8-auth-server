package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

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

func ValidateSignature(message string, signature []byte, publicKey []byte) error {
	msgHash := crypto.Keccak256([]byte(message))
	verified := crypto.VerifySignature(publicKey, msgHash, signature)

	if !verified {
		return fmt.Errorf("failed to verify the ecdsa signature")
	}

	return nil
}
