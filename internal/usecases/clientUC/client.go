package clientUC

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/pkg/scram"
	utils2 "globe-and-citizen/layer8/auth-server/pkg/utils"
)

func (uc *ClientUsecase) CheckBackendURI(req requestdto.ClientCheckBackendURI) (bool, error) {
	response, err := uc.postgres.IsBackendURIExists(req.BackendURI)
	if err != nil {
		return false, err
	}

	return response, nil
}

func (uc *ClientUsecase) PrecheckRegister(
	req requestdto.ClientRegisterPrecheck,
	iterCount int,
) (responsedto.ClientRegisterPrecheck, error) {
	scramMsg := scram.CreateServerRegisterFirstMessage(iterCount)

	client := gormModels.Client{
		Username:            req.Username,
		ScramSalt:           scramMsg.Salt,
		ScramIterationCount: iterCount,
	}

	err := uc.postgres.PrecheckClientRegister(client)
	if err != nil {
		return responsedto.ClientRegisterPrecheck{}, err
	}

	return responsedto.ClientRegisterPrecheck{
		ServerRegisterFirstMessage: scramMsg,
	}, nil
}

func (uc *ClientUsecase) Register(req requestdto.ClientRegister) error {
	clientUUID := utils2.GenerateUUID()
	clientSecret := utils2.GenerateSecret(consts.SecretSize)
	backendURI, err := utils2.GetURLHostPort(req.BackendURI)
	if err != nil {
		return err
	}

	newClient := gormModels.Client{
		ID:             clientUUID,
		Secret:         clientSecret,
		Name:           req.Name,
		RedirectURI:    req.RedirectURI,
		BackendURI:     backendURI,
		Username:       req.Username,
		ScramStoredKey: req.StoredKey,
		ScramServerKey: req.ServerKey,
	}

	return uc.postgres.UpdateClient(newClient)
}

func (uc *ClientUsecase) PrecheckLogin(req requestdto.ClientLoginPrecheck) (responsedto.ClientLoginPrecheck, error) {
	client, err := uc.postgres.GetClientByUsername(req.Username)
	if err != nil {
		return responsedto.ClientLoginPrecheck{}, err
	}

	scramMsg := scram.CreateServerLoginFirstMessage(client.ScramSalt, client.ScramIterationCount, req.ClientLoginFirstMessage)
	loginPrecheckResp := responsedto.ClientLoginPrecheck{
		ServerLoginFirstMessage: scramMsg,
	}

	return loginPrecheckResp, nil
}

func (uc *ClientUsecase) Login(req requestdto.ClientLogin) (responsedto.ClientLogin, error) {
	client, err := uc.postgres.GetClientByUsername(req.Username)
	if err != nil {
		return responsedto.ClientLogin{}, err
	}

	scramMsg, err := scram.CreateServerLoginFinalMessage(req.ClientLoginFinalMessage, req.CNonce,
		client.ScramSalt, client.ScramIterationCount, client.ScramStoredKey, client.ScramServerKey)
	if err != nil {
		return responsedto.ClientLogin{}, fmt.Errorf("error creating final message: %v", err)
	}

	tokenString, err := uc.token.GenerateClientJWTToken(client)
	if err != nil {
		return responsedto.ClientLogin{}, fmt.Errorf("error generating token: %v", err)
	}

	return responsedto.ClientLogin{
		ServerLoginFinalMessage: scramMsg,
		Token:                   tokenString,
	}, nil
}

func (uc *ClientUsecase) GetProfile(username string) (responsedto.ClientProfile, error) {
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
		NTorCertificate: string(clientData.NTorX509Certificate),
	}
	return clientModel, nil
}

func (uc *ClientUsecase) GetUnpaidAmount(clientID string) (responsedto.ClientGetBalance, error) {
	stats, err := uc.postgres.GetClientBalance(clientID)
	if err != nil {
		return responsedto.ClientGetBalance{}, err
	}

	balanceWei, err := utils2.DBWeiToBigInt(stats.BalanceWei)
	if err != nil {
		return responsedto.ClientGetBalance{}, err
	}

	return responsedto.ClientGetBalance{
		Balance: utils2.WeiToEthString(balanceWei, 18),
	}, nil
}

func (uc *ClientUsecase) SaveNTorCertificate(clientID string, req requestdto.ClientUploadNTorCertificate) error {
	// todo validate certificate
	return uc.postgres.SaveX509Certificate(clientID, req.Certificate)
}

func (uc *ClientUsecase) GetNTorCertificate(req requestdto.ClientGetNTorCertificate) (*responsedto.ClientGetNTorCertificate, error) {
	client, err := uc.postgres.GetClientByBackendURI(req.BackendURI)
	if err != nil {
		return nil, err
	}

	return &responsedto.ClientGetNTorCertificate{
		ClientID:    client.ID,
		Certificate: string(client.NTorX509Certificate),
	}, nil
}
