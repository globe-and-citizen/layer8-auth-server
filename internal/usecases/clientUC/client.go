package clientUC

import (
	"context"
	"errors"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/internal/usecases/ucerror"
	"globe-and-citizen/layer8/auth-server/pkg/scram"
	"globe-and-citizen/layer8/auth-server/pkg/utils"

	"gorm.io/gorm"
)

func (uc *ClientUsecase) CheckBackendURI(ctx context.Context, req requestdto.ClientCheckBackendURI) (bool, *ucerror.UCError) {
	response, err := uc.postgres.IsBackendURIExists(ctx, req.BackendURI)
	if err != nil {
		return false, ucerror.New(fmt.Errorf("failed to check backendURI: %w", err), consts.ErrInternalServer)
	}

	return response, nil
}

func (uc *ClientUsecase) PrecheckRegister(
	ctx context.Context,
	req requestdto.ClientRegisterPrecheck,
	iterCount int,
) (responsedto.ClientRegisterPrecheck, *ucerror.UCError) {
	scramMsg := scram.CreateServerRegisterFirstMessage(iterCount)

	client := gormModels.Client{
		Username:            req.Username,
		ScramSalt:           scramMsg.Salt,
		ScramIterationCount: iterCount,
	}

	err := uc.postgres.CreateClient(ctx, client)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return responsedto.ClientRegisterPrecheck{}, ucerror.New(err, consts.ErrDuplicateKey)
		}
		return responsedto.ClientRegisterPrecheck{}, ucerror.New(err, consts.ErrInternalServer)
	}

	return responsedto.ClientRegisterPrecheck{
		ServerRegisterFirstMessage: scramMsg,
	}, nil
}

func (uc *ClientUsecase) Register(ctx context.Context, req requestdto.ClientRegister) *ucerror.UCError {
	clientUUID := utils.GenerateUUID()
	clientSecret := utils.GenerateSecret(consts.SecretSize)
	backendURI, err := utils.GetURLHostPort(req.BackendURI)
	if err != nil {
		return ucerror.New(
			fmt.Errorf("errors extract banckend uri: %w", err),
			fmt.Errorf("%w: invalid backendURI", consts.ErrBadRequest),
		)
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

	err = uc.postgres.UpdateClient(ctx, newClient)
	return ucerror.New(err, consts.ErrInternalServer) // todo handle error properly
}

func (uc *ClientUsecase) PrecheckLogin(
	ctx context.Context,
	req requestdto.ClientLoginPrecheck,
) (responsedto.ClientLoginPrecheck, *ucerror.UCError) {
	client, err := uc.postgres.GetClientByUsername(ctx, req.Username)
	if err != nil {
		return responsedto.ClientLoginPrecheck{}, ucerror.New(err, consts.ErrBadRequest) // todo handle errors
	}

	scramMsg := scram.CreateServerLoginFirstMessage(client.ScramSalt, client.ScramIterationCount, req.ClientLoginFirstMessage)
	loginPrecheckResp := responsedto.ClientLoginPrecheck{
		ServerLoginFirstMessage: scramMsg,
	}

	return loginPrecheckResp, nil
}

func (uc *ClientUsecase) Login(ctx context.Context, req requestdto.ClientLogin) (responsedto.ClientLogin, *ucerror.UCError) {
	client, err := uc.postgres.GetClientByUsername(ctx, req.Username)
	if err != nil {
		return responsedto.ClientLogin{}, ucerror.New(fmt.Errorf("failed to get login client: %w", err), consts.ErrBadRequest) // todo
	}

	scramMsg, err := scram.CreateServerLoginFinalMessage(req.ClientLoginFinalMessage, req.CNonce,
		client.ScramSalt, client.ScramIterationCount, client.ScramStoredKey, client.ScramServerKey)
	if err != nil {
		return responsedto.ClientLogin{}, ucerror.New(fmt.Errorf("error creating final message: %w", err), consts.ErrInternalServer) //todo
	}

	tokenString, err := uc.token.GenerateClientJWTToken(client)
	if err != nil {
		return responsedto.ClientLogin{}, ucerror.New(fmt.Errorf("error generating token: %w", err), consts.ErrInternalServer)
	}

	return responsedto.ClientLogin{
		ServerLoginFinalMessage: scramMsg,
		Token:                   tokenString,
	}, nil
}

func (uc *ClientUsecase) GetProfile(ctx context.Context, username string) (responsedto.ClientProfile, *ucerror.UCError) {
	clientData, err := uc.postgres.GetClientByUsername(ctx, username)
	if err != nil {
		return responsedto.ClientProfile{}, ucerror.New(fmt.Errorf("failed to get client data: %w", err), consts.ErrBadRequest) // todo
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

func (uc *ClientUsecase) GetUnpaidAmount(ctx context.Context, clientID string) (responsedto.ClientGetBalance, *ucerror.UCError) {
	stats, err := uc.postgres.GetClientBalance(ctx, clientID)
	if err != nil {
		return responsedto.ClientGetBalance{}, ucerror.New(fmt.Errorf("failed to get client balance: %w", err), consts.ErrInternalServer) // todo
	}

	balanceWei, err := utils.DBWeiToBigInt(stats.BalanceWei)
	if err != nil {
		return responsedto.ClientGetBalance{}, ucerror.New(fmt.Errorf("failed to convert balance wei: %w", err), consts.ErrInternalServer)
	}

	return responsedto.ClientGetBalance{
		Balance: utils.WeiToEthString(balanceWei, 18),
	}, nil
}

func (uc *ClientUsecase) SaveNTorCertificate(ctx context.Context, clientID string, req requestdto.ClientUploadNTorCertificate) *ucerror.UCError {
	// todo validate certificate
	err := uc.postgres.SaveX509Certificate(ctx, clientID, req.Certificate)
	if err != nil {
		ucerror.New(fmt.Errorf("failed to save ntor certificate: %w", err), consts.ErrInternalServer) // todo
	}
	return nil
}

func (uc *ClientUsecase) GetNTorCertificate(ctx context.Context, req requestdto.ClientGetNTorCertificate) (*responsedto.ClientGetNTorCertificate, *ucerror.UCError) {
	client, err := uc.postgres.GetClientByBackendURI(ctx, req.BackendURI)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ucerror.New(fmt.Errorf("client with backend_url %s: %w", req.BackendURI, err), consts.ErrNotFound)
		}
		return nil, ucerror.New(fmt.Errorf("failed to get client data: %w", err), consts.ErrInternalServer)
	}

	return &responsedto.ClientGetNTorCertificate{
		ClientID:    client.ID,
		Certificate: string(client.NTorX509Certificate), // todo check this format conversion
	}, nil
}
