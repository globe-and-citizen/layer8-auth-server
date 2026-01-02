package scram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"globe-and-citizen/layer8/auth-server/backend/pkg/utils"
)

const SaltSize = 16
const NonceSize = 16

func CreateServerRegisterFirstMessage(iterationCount int) ServerRegisterFirstMessage {
	return ServerRegisterFirstMessage{
		Salt:           utils.GenerateRandomSalt(SaltSize),
		IterationCount: iterationCount,
	}
}

func CreateServerLoginFirstMessage(salt string, iterationCount int, clientFirstMsg ClientLoginFirstMessage) ServerLoginFirstMessage {
	serverNonce := clientFirstMsg.ClientNonce + utils.GenerateRandomSalt(NonceSize)

	return ServerLoginFirstMessage{
		Salt:           salt,
		IterationCount: iterationCount,
		Nonce:          clientFirstMsg.ClientNonce + serverNonce,
	}
}

func CreateServerLoginFinalMessage(
	cFinalMsg ClientLoginFinalMessage,
	clientNonce string,
	salt string,
	iterationCount int,
	storedKey string,
	serverKey string,
) (ServerLoginFinalMessage, error) {
	storedKeyBytes, err := hex.DecodeString(storedKey)
	if err != nil {
		return ServerLoginFinalMessage{}, fmt.Errorf("error decoding stored key: %v", err)
	}

	authMessage := fmt.Sprintf("[n=%s,r=%s,s=%s,i=%d,r=%s]", cFinalMsg.Username, clientNonce, salt, iterationCount, cFinalMsg.Nonce)
	//authMessage := fmt.Sprintf("n=%s,r=%s,r=%s,s=%s,i=%d,c=%s,r=%s", cFinalMsg.Username, clientNonce, cFinalMsg.Nonce, salt, iterationCount, cFinalMsg.ChannelBinding, cFinalMsg.Nonce)
	authMessageBytes := []byte(authMessage)

	clientSignatureHMAC := hmac.New(sha256.New, storedKeyBytes)
	clientSignatureHMAC.Write(authMessageBytes)
	clientSignature := clientSignatureHMAC.Sum(nil)

	clientProofBytes, err := hex.DecodeString(cFinalMsg.ClientProof)
	if err != nil {
		return ServerLoginFinalMessage{}, fmt.Errorf("error decoding client proof: %v", err)
	}

	clientKeyBytes, err := utils.XorBytes(clientSignature, clientProofBytes)
	if err != nil {
		return ServerLoginFinalMessage{}, fmt.Errorf("error performing XOR operation: %v", err)
	}

	clientKeyHash := sha256.Sum256(clientKeyBytes)

	clientKeyHashStr := hex.EncodeToString(clientKeyHash[:])
	if clientKeyHashStr != storedKey {
		return ServerLoginFinalMessage{}, fmt.Errorf("server failed to authenticate the user")
	}

	serverKeyBytes, err := hex.DecodeString(serverKey)
	if err != nil {
		return ServerLoginFinalMessage{}, fmt.Errorf("error decoding server key: %v", err)
	}

	serverSignatureHMAC := hmac.New(sha256.New, serverKeyBytes)
	serverSignatureHMAC.Write(authMessageBytes)
	serverSignatureHex := hex.EncodeToString(serverSignatureHMAC.Sum(nil))

	return ServerLoginFinalMessage{
		ServerSignature: serverSignatureHex,
	}, nil
}
