package clientUC

import (
	"context"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/repositories/influxdbRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/postgresRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/tokenRepo"
	"globe-and-citizen/layer8/auth-server/internal/usecases/ucerror"
	"globe-and-citizen/layer8/auth-server/pkg/log"
)

type IClientUsecase interface {
	CheckBackendURI(ctx context.Context, req requestdto.ClientCheckBackendURI) (bool, *ucerror.UCError)
	PrecheckRegister(ctx context.Context, req requestdto.ClientRegisterPrecheck, iterCount int) (*responsedto.ClientRegisterPrecheck, *ucerror.UCError)
	Register(ctx context.Context, req requestdto.ClientRegister) *ucerror.UCError
	PrecheckLogin(ctx context.Context, req requestdto.ClientLoginPrecheck) (*responsedto.ClientLoginPrecheck, *ucerror.UCError)
	Login(ctx context.Context, req requestdto.ClientLogin) (*responsedto.ClientLogin, *ucerror.UCError)
	GetProfile(ctx context.Context, username string) (*responsedto.ClientProfile, *ucerror.UCError)
	GetUsageStatistics(ctx context.Context, clientID string) (*responsedto.ClientUsageStatistic, *ucerror.UCError)
	GetUnpaidAmount(ctx context.Context, clientID string) (*responsedto.ClientGetBalance, *ucerror.UCError)
	SaveNTorCertificate(ctx context.Context, clientID string, req requestdto.ClientUploadNTorCertificate) *ucerror.UCError
	MdwVerifyClientJWTToken(ctx context.Context, tokenString string) (clientID string, clientUsername string, err *ucerror.UCError)
	GetNTorCertificate(ctx context.Context, req requestdto.ClientGetNTorCertificate) (*responsedto.ClientGetNTorCertificate, *ucerror.UCError)
}

type ClientUsecase struct {
	logger   log.ILogger
	postgres postgresRepo.IClientRepositories
	token    tokenRepo.ITokenRepository
	influxdb influxdbRepo.IInfluxdbRepository
}

func NewClientUsecase(
	logger log.ILogger,
	postgres postgresRepo.IClientRepositories,
	token tokenRepo.ITokenRepository,
	influxdb influxdbRepo.IInfluxdbRepository,
) IClientUsecase {
	return &ClientUsecase{
		logger:   logger,
		postgres: postgres,
		token:    token,
		influxdb: influxdb,
	}
}
