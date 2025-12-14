package userUsecase

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
)

type IUserUseCase interface {
	PrecheckRegister(req requestdto.UserRegisterPrecheck, iterCount int) (responsedto.UserRegisterPrecheck, error)
	Register(req requestdto.UserRegister) error
	PrecheckLogin(req requestdto.UserLoginPrecheck) (responsedto.UserLoginPrecheck, error)
	Login(req requestdto.UserLogin) (responsedto.UserLogin, error)
	UpdateUserMetadata(userID uint, req requestdto.UserMetadataUpdate) error
	VerifyEmail(userID uint, userEmail string) error
	CheckEmailVerificationCode(userId uint, code string) error
	SaveProofOfEmailVerification(userID uint, req requestdto.CheckEmailVerificationCode) (string, error)
}
