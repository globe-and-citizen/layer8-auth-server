package clientUC

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/postgresRepo"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/statsRepo"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/tokenRepo"
)

type IClientUsecase interface {
	CheckBackendURI(req requestdto.ClientCheckBackendURI) (bool, error)
	PrecheckRegister(req requestdto.ClientRegisterPrecheck, iterCount int) (responsedto.ClientRegisterPrecheck, error)
	Register(req requestdto.ClientRegister) error
	PrecheckLogin(req requestdto.ClientLoginPrecheck) (responsedto.ClientLoginPrecheck, error)
	Login(req requestdto.ClientLogin) (responsedto.ClientLogin, error)
	GetProfile(username string) (responsedto.ClientProfile, error)
	GetUsageStatistics(clientID string) (responsedto.ClientUsageStatistic, int, string, error)
	GetUnpaidAmount(clientID string) (responsedto.ClientUnpaidAmount, error)
	SaveNTorCertificate(clientID string, req requestdto.ClientUploadNTorCertificate) error
	VerifyClientJWTToken(tokenString string) (clientID string, clientUsername string, err error)
}

type ClientUsecase struct {
	postgres postgresRepo.IClientRepositories
	token    tokenRepo.ITokenRepository
	stats    statsRepo.IStatisticsRepository
}

func NewClientUsecase(
	postgres postgresRepo.IClientRepositories,
	token tokenRepo.ITokenRepository,
	stats statsRepo.IStatisticsRepository,
) IClientUsecase {
	return &ClientUsecase{
		postgres: postgres,
		token:    token,
		stats:    stats,
	}
}
