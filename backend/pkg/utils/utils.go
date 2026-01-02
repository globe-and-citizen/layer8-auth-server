package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
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

func GenerateRandomBase64String(size int) (string, error) {
	buf := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", fmt.Errorf("could not generate random bytes: %s", err)
	}

	return base64.URLEncoding.EncodeToString(buf), nil
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
