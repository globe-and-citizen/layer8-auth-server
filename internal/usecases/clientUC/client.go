package clientUC

import (
	"context"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/pkg/scram"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
)

func (uc *ClientUsecase) CheckBackendURI(ctx context.Context, req requestdto.ClientCheckBackendURI) (bool, error) {
	response, err := uc.postgres.IsBackendURIExists(ctx, req.BackendURI)
	if err != nil {
		return false, err
	}

	return response, nil
}

func (uc *ClientUsecase) PrecheckRegister(
	ctx context.Context,
	req requestdto.ClientRegisterPrecheck,
	iterCount int,
) (responsedto.ClientRegisterPrecheck, error) {
	scramMsg := scram.CreateServerRegisterFirstMessage(iterCount)

	client := gormModels.Client{
		Username:            req.Username,
		ScramSalt:           scramMsg.Salt,
		ScramIterationCount: iterCount,
	}

	err := uc.postgres.PrecheckClientRegister(ctx, client)
	if err != nil {
		return responsedto.ClientRegisterPrecheck{}, err
	}

	return responsedto.ClientRegisterPrecheck{
		ServerRegisterFirstMessage: scramMsg,
	}, nil
}

func (uc *ClientUsecase) Register(ctx context.Context, req requestdto.ClientRegister) error {
	clientUUID := utils.GenerateUUID()
	clientSecret := utils.GenerateSecret(consts.SecretSize)
	backendURI, err := utils.GetURLHostPort(req.BackendURI)
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

	return uc.postgres.UpdateClient(ctx, newClient)
}

func (uc *ClientUsecase) PrecheckLogin(ctx context.Context, req requestdto.ClientLoginPrecheck) (responsedto.ClientLoginPrecheck, error) {
	client, err := uc.postgres.GetClientByUsername(ctx, req.Username)
	if err != nil {
		return responsedto.ClientLoginPrecheck{}, err
	}

	scramMsg := scram.CreateServerLoginFirstMessage(client.ScramSalt, client.ScramIterationCount, req.ClientLoginFirstMessage)
	loginPrecheckResp := responsedto.ClientLoginPrecheck{
		ServerLoginFirstMessage: scramMsg,
	}

	return loginPrecheckResp, nil
}

func (uc *ClientUsecase) Login(ctx context.Context, req requestdto.ClientLogin) (responsedto.ClientLogin, error) {
	client, err := uc.postgres.GetClientByUsername(ctx, req.Username)
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

func (uc *ClientUsecase) GetProfile(ctx context.Context, username string) (responsedto.ClientProfile, error) {
	clientData, err := uc.postgres.GetClientByUsername(ctx, username)
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

func (uc *ClientUsecase) GetUnpaidAmount(ctx context.Context, clientID string) (responsedto.ClientGetBalance, error) {
	stats, err := uc.postgres.GetClientBalance(ctx, clientID)
	if err != nil {
		return responsedto.ClientGetBalance{}, err
	}

	balanceWei, err := utils.DBWeiToBigInt(stats.BalanceWei)
	if err != nil {
		return responsedto.ClientGetBalance{}, err
	}

	return responsedto.ClientGetBalance{
		Balance: utils.WeiToEthString(balanceWei, 18),
	}, nil
}

func (uc *ClientUsecase) SaveNTorCertificate(ctx context.Context, clientID string, req requestdto.ClientUploadNTorCertificate) error {
	// todo validate certificate
	return uc.postgres.SaveX509Certificate(ctx, clientID, req.Certificate)
}

func (uc *ClientUsecase) GetNTorCertificate(ctx context.Context, req requestdto.ClientGetNTorCertificate) (*responsedto.ClientGetNTorCertificate, error) {
	client, err := uc.postgres.GetClientByBackendURI(ctx, req.BackendURI)
	if err != nil {
		return nil, err
	}

	return &responsedto.ClientGetNTorCertificate{
		ClientID:    client.ID,
		Certificate: string(client.NTorX509Certificate),
	}, nil
}
