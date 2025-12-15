package userUsecase

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
)

type IUserUseCase interface { // todo usecase methods should return custom error type that contains http status codes, message and error
	PrecheckRegister(req requestdto.UserRegisterPrecheck, iterCount int) (responsedto.UserRegisterPrecheck, error)
	Register(req requestdto.UserRegister) error
	PrecheckLogin(req requestdto.UserLoginPrecheck) (responsedto.UserLoginPrecheck, error)
	Login(req requestdto.UserLogin) (responsedto.UserLogin, error)
	UpdateUserMetadata(userID uint, req requestdto.UserMetadataUpdate) error
	VerifyEmail(userID uint, userEmail string) error
	CheckEmailVerificationCode(userId uint, code string) error
	SaveProofOfEmailVerification(userID uint, req requestdto.UserCheckEmailVerificationCode) (msg string, err error)
	VerifyPhoneNumber(userID uint) (errMsg string, err error)
	CheckPhoneNumberVerificationCode(userID uint, req requestdto.UserCheckPhoneNumberVerificationCode) (httpStatus int, msg string, err error)
	GenerateAndSaveTelegramSessionIDHash(userID uint) (sessionID []byte, msg string, err error)
}
