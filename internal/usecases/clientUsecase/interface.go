package clientUsecase

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/repositories/postgresRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/tokenRepo"
)

type IClientUsecase interface {
	CheckBackendURI(req requestdto.CheckBackendURI) (bool, error)
	PrecheckRegister(req requestdto.ClientRegisterPrecheck, iterCount int) (responsedto.ClientRegisterPrecheck, error)
	Register(req requestdto.ClientRegister) error
	PrecheckLogin(req requestdto.ClientLoginPrecheck) (responsedto.ClientLoginPrecheck, error)
	Login(req requestdto.ClientLogin) (responsedto.ClientLogin, error)
	GetProfile(username string) (responsedto.ClientProfile, error)
}

type ClientUsecase struct {
	postgres postgresRepo.IClientRepositories
	token    tokenRepo.ITokenRepository
}

func NewClientUsecase(
	postgres postgresRepo.IClientRepositories,
	token tokenRepo.ITokenRepository,
) IClientUsecase {
	return &ClientUsecase{
		postgres: postgres,
		token:    token,
	}
}
