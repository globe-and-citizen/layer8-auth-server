package userUC

import (
	"context"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/repositories/codeGenRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/emailRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/phoneRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/postgresRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/tokenRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/zkRepo"
)

type IUserUsecase interface { // todo usecase methods should return custom error type that contains http status codes, message and error
	PrecheckRegister(ctx context.Context, req requestdto.UserRegisterPrecheck, iterCount int) (responsedto.UserRegisterPrecheck, error)
	Register(ctx context.Context, req requestdto.UserRegister) error
	PrecheckLogin(ctx context.Context, req requestdto.UserLoginPrecheck) (responsedto.UserLoginPrecheck, error)
	Login(ctx context.Context, req requestdto.UserLogin) (responsedto.UserLogin, error)
	GetProfile(ctx context.Context, userID uint) (responsedto.UserProfile, error)
	UpdateUserMetadata(ctx context.Context, userID uint, req requestdto.UserMetadataUpdate) error
	VerifyEmail(ctx context.Context, userID uint, userEmail string) error
	CheckEmailVerificationCode(ctx context.Context, userId uint, code string) error
	SaveProofOfEmailVerification(ctx context.Context, userID uint, req requestdto.UserCheckEmailVerificationCode) (msg string, err error)
	VerifyPhoneNumber(ctx context.Context, userID uint) (errMsg string, err error)
	CheckPhoneNumberVerificationCode(ctx context.Context, userID uint, req requestdto.UserCheckPhoneNumberVerificationCode) (httpStatus int, msg string, err error)
	GenerateAndSaveTelegramSessionIDHash(ctx context.Context, userID uint) (sessionID []byte, msg string, err error)
	PrecheckResetPassword(ctx context.Context, req requestdto.UserResetPasswordPrecheck) (responsedto.UserResetPasswordPrecheck, error)
	ResetPassword(ctx context.Context, request requestdto.UserResetPassword) (httpStatus int, msg string, err error)
	VerifyUserJWTToken(ctx context.Context, tokenString string) (userID uint, userUsername string, err error)
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
