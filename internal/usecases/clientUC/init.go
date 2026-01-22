package clientUC

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/repositories/influxdbRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/postgresRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/tokenRepo"
)

type IClientUsecase interface {
	CheckBackendURI(req requestdto.ClientCheckBackendURI) (bool, error)
	PrecheckRegister(req requestdto.ClientRegisterPrecheck, iterCount int) (responsedto.ClientRegisterPrecheck, error)
	Register(req requestdto.ClientRegister) error
	PrecheckLogin(req requestdto.ClientLoginPrecheck) (responsedto.ClientLoginPrecheck, error)
	Login(req requestdto.ClientLogin) (responsedto.ClientLogin, error)
	GetProfile(username string) (responsedto.ClientProfile, error)
	GetUsageStatistics(clientID string) (responsedto.ClientUsageStatistic, int, string, error)
	GetUnpaidAmount(clientID string) (responsedto.ClientGetBalance, error)
	SaveNTorCertificate(clientID string, req requestdto.ClientUploadNTorCertificate) error
	VerifyClientJWTToken(tokenString string) (clientID string, clientUsername string, err error)
	GetNTorCertificate(req requestdto.ClientGetNTorCertificate) (*responsedto.ClientGetNTorCertificate, error)
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
