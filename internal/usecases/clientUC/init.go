package clientUC

import (
	"context"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/repositories/influxdbRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/postgresRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/tokenRepo"
)

type IClientUsecase interface {
	CheckBackendURI(ctx context.Context, req requestdto.ClientCheckBackendURI) (bool, error)
	PrecheckRegister(ctx context.Context, req requestdto.ClientRegisterPrecheck, iterCount int) (responsedto.ClientRegisterPrecheck, error)
	Register(ctx context.Context, req requestdto.ClientRegister) error
	PrecheckLogin(ctx context.Context, req requestdto.ClientLoginPrecheck) (responsedto.ClientLoginPrecheck, error)
	Login(ctx context.Context, req requestdto.ClientLogin) (responsedto.ClientLogin, error)
	GetProfile(ctx context.Context, username string) (responsedto.ClientProfile, error)
	GetUsageStatistics(ctx context.Context, clientID string) (responsedto.ClientUsageStatistic, int, string, error)
	GetUnpaidAmount(ctx context.Context, clientID string) (responsedto.ClientGetBalance, error)
	SaveNTorCertificate(ctx context.Context, clientID string, req requestdto.ClientUploadNTorCertificate) error
	VerifyClientJWTToken(ctx context.Context, tokenString string) (clientID string, clientUsername string, err error)
	GetNTorCertificate(ctx context.Context, req requestdto.ClientGetNTorCertificate) (*responsedto.ClientGetNTorCertificate, error)
}

type ClientUsecase struct {
	postgres postgresRepo.IClientRepositories
	token    tokenRepo.ITokenRepository
	influxdb influxdbRepo.IInfluxdbRepository
}

func NewClientUsecase(
	postgres postgresRepo.IClientRepositories,
	token tokenRepo.ITokenRepository,
	influxdb influxdbRepo.IInfluxdbRepository,
) IClientUsecase {
	return &ClientUsecase{
		postgres: postgres,
		token:    token,
		influxdb: influxdb,
	}
}
