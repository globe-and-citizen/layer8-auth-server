package clientUC

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/backend/internal/consts"
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/backend/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/backend/pkg/scram"
	"globe-and-citizen/layer8/auth-server/backend/pkg/utils"
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
	clientUUID := utils.GenerateUUID()
	clientSecret := utils.GenerateSecret(consts.SecretSize)
	backendURI := utils.RemoveProtocolFromURL(req.BackendURI)

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

	return uc.postgres.AddClient(newClient)
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
		NTorCertificate: clientData.NTorX509Certificate,
	}
	return clientModel, nil
}

func (uc *ClientUsecase) GetUnpaidAmount(clientID string) (responsedto.ClientUnpaidAmount, error) {
	stats, err := uc.postgres.GetClientTrafficStatistics(clientID)
	if err != nil {
		return responsedto.ClientUnpaidAmount{}, err
	}

	return responsedto.ClientUnpaidAmount{
		UnpaidAmount: stats.UnpaidAmount,
	}, nil
}

func (uc *ClientUsecase) SaveNTorCertificate(clientID string, req requestdto.ClientUploadNTorCertificate) error {
	// todo validate certificate
	return uc.postgres.SaveX509Certificate(clientID, req.Certificate)
}
