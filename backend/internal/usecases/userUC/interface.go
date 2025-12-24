package userUC

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/codeGenRepo"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/emailRepo"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/phoneRepo"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/postgresRepo"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/tokenRepo"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/zkRepo"
)

type IUserUsecase interface { // todo usecase methods should return custom error type that contains http status codes, message and error
	PrecheckRegister(req requestdto.UserRegisterPrecheck, iterCount int) (responsedto.UserRegisterPrecheck, error)
	Register(req requestdto.UserRegister) error
	PrecheckLogin(req requestdto.UserLoginPrecheck) (responsedto.UserLoginPrecheck, error)
	Login(req requestdto.UserLogin) (responsedto.UserLogin, error)
	GetProfile(userID uint) (responsedto.UserProfile, error)
	UpdateUserMetadata(userID uint, req requestdto.UserMetadataUpdate) error
	VerifyEmail(userID uint, userEmail string) error
	CheckEmailVerificationCode(userId uint, code string) error
	SaveProofOfEmailVerification(userID uint, req requestdto.UserCheckEmailVerificationCode) (msg string, err error)
	VerifyPhoneNumber(userID uint) (errMsg string, err error)
	CheckPhoneNumberVerificationCode(userID uint, req requestdto.UserCheckPhoneNumberVerificationCode) (httpStatus int, msg string, err error)
	GenerateAndSaveTelegramSessionIDHash(userID uint) (sessionID []byte, msg string, err error)
	PrecheckResetPassword(req requestdto.UserResetPasswordPrecheck) (responsedto.UserResetPasswordPrecheck, error)
	ResetPassword(request requestdto.UserResetPassword) (httpStatus int, msg string, err error)
}

type UserUsecase struct {
	postgres postgresRepo.IUserRepositories
	token    tokenRepo.ITokenRepository
	email    emailRepo.IEmailRepository
	code     codeGenRepo.ICodeGeneratorRepository
	zk       zkRepo.IZkRepository
	phone    phoneRepo.IPhoneRepository
}

func NewUserUsecase(
	postgres postgresRepo.IUserRepositories,
	token tokenRepo.ITokenRepository,
	email emailRepo.IEmailRepository,
	code codeGenRepo.ICodeGeneratorRepository,
	zk zkRepo.IZkRepository,
	phone phoneRepo.IPhoneRepository,
) IUserUsecase {
	return &UserUsecase{
		postgres: postgres,
		token:    token,
		email:    email,
		code:     code,
		zk:       zk,
		phone:    phone,
	}
}
