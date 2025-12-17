package clientUC

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
)

func (uc ClientUsecase) CheckBackendURI(req requestdto.ClientCheckBackendURI) (bool, error) {
	response, err := uc.postgres.IsBackendURIExists(req.BackendURI)
	if err != nil {
		return false, err
	}

	// todo validate format?
	return response, nil
}

func (uc ClientUsecase) PrecheckRegister(
	req requestdto.ClientRegisterPrecheck,
	iterCount int,
) (responsedto.ClientRegisterPrecheck, error) {
	rmSalt := utils.GenerateRandomSalt(consts.SaltSize)

	client := gormModels.Client{
		Username:       req.Username,
		Salt:           rmSalt,
		IterationCount: iterCount,
	}

	err := uc.postgres.PrecheckClientRegister(client)
	if err != nil {
		return responsedto.ClientRegisterPrecheck{}, err
	}

	return responsedto.ClientRegisterPrecheck{
		Salt:           rmSalt,
		IterationCount: iterCount,
	}, nil
}

func (uc ClientUsecase) Register(req requestdto.ClientRegister) error {
	clientUUID := utils.GenerateUUID()
	clientSecret := utils.GenerateSecret(consts.SecretSize)
	backendURI := utils.RemoveProtocolFromURL(req.BackendURI)

	newClient := gormModels.Client{
		ID:          clientUUID,
		Secret:      clientSecret,
		Name:        req.Name,
		RedirectURI: req.RedirectURI,
		BackendURI:  backendURI,
		Username:    req.Username,
		StoredKey:   req.StoredKey,
		ServerKey:   req.ServerKey,
	}

	return uc.postgres.AddClient(newClient)
}

func (uc ClientUsecase) PrecheckLogin(req requestdto.ClientLoginPrecheck) (responsedto.ClientLoginPrecheck, error) {
	sNonce := utils.GenerateRandomSalt(consts.SaltSize)

	client, err := uc.postgres.GetClientByUsername(req.Username)
	if err != nil {
		return responsedto.ClientLoginPrecheck{}, err
	}

	loginPrecheckResp := responsedto.ClientLoginPrecheck{
		Salt:      client.Salt,
		IterCount: client.IterationCount,
		Nonce:     req.CNonce + sNonce,
	}

	return loginPrecheckResp, nil
}

func (uc ClientUsecase) Login(req requestdto.ClientLogin) (responsedto.ClientLogin, error) {
	client, err := uc.postgres.GetClientByUsername(req.Username)
	if err != nil {
		return responsedto.ClientLogin{}, err
	}

	storedKeyBytes, err := hex.DecodeString(client.StoredKey)
	if err != nil {
		return responsedto.ClientLogin{}, fmt.Errorf("error decoding stored key: %v", err)
	}

	authMessage := fmt.Sprintf("[n=%s,r=%s,s=%s,i=%d,r=%s]", req.Username, req.CNonce, client.Salt, client.IterationCount, req.Nonce)
	authMessageBytes := []byte(authMessage)

	clientSignatureHMAC := hmac.New(sha256.New, storedKeyBytes)
	clientSignatureHMAC.Write(authMessageBytes)
	clientSignature := clientSignatureHMAC.Sum(nil)

	clientProofBytes, err := hex.DecodeString(req.ClientProof)
	if err != nil {
		return responsedto.ClientLogin{}, fmt.Errorf("error decoding client proof: %v", err)
	}

	clientKeyBytes, err := utils.XorBytes(clientSignature, clientProofBytes)
	if err != nil {
		return responsedto.ClientLogin{}, fmt.Errorf("error performing XOR operation: %v", err)
	}

	clientKeyHash := sha256.Sum256(clientKeyBytes)

	clientKeyHashStr := hex.EncodeToString(clientKeyHash[:])
	if clientKeyHashStr != client.StoredKey {
		return responsedto.ClientLogin{}, fmt.Errorf("server failed to authenticate the user")
	}

	serverKeyBytes, err := hex.DecodeString(client.ServerKey)
	if err != nil {
		return responsedto.ClientLogin{}, fmt.Errorf("error decoding server key: %v", err)
	}

	serverSignatureHMAC := hmac.New(sha256.New, serverKeyBytes)
	serverSignatureHMAC.Write(authMessageBytes)
	serverSignatureHex := hex.EncodeToString(serverSignatureHMAC.Sum(nil))

	tokenString, err := uc.token.GenerateClientJWTToken(client)
	if err != nil {
		return responsedto.ClientLogin{}, fmt.Errorf("error generating token: %v", err)
	}

	return responsedto.ClientLogin{
		ServerSignature: serverSignatureHex,
		Token:           tokenString,
	}, nil
}

func (uc ClientUsecase) GetProfile(username string) (responsedto.ClientProfile, error) {
	clientData, err := uc.postgres.GetClientByUsername(username)
	if err != nil {
		return responsedto.ClientProfile{}, err
	}

	clientModel := responsedto.ClientProfile{
		ID:              clientData.ID,
		Secret:          clientData.Secret,
		Name:            clientData.Name,
		RedirectURI:     clientData.RedirectURI,
		BackendURI:      clientData.BackendURI,
		X509Certificate: clientData.X509Certificate,
	}
	return clientModel, nil
}

func (uc ClientUsecase) GetUnpaidAmount(clientID string) (responsedto.ClientUnpaidAmount, error) {
	stats, err := uc.postgres.GetClientTrafficStatistics(clientID)
	if err != nil {
		return responsedto.ClientUnpaidAmount{}, err
	}

	return responsedto.ClientUnpaidAmount{
		UnpaidAmount: stats.UnpaidAmount,
	}, nil
}
