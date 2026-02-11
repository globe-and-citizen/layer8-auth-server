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
	"globe-and-citizen/layer8/auth-server/internal/usecases/ucerror"
)

type IUserUsecase interface { // todo usecase methods should return custom error type that contains http status codes, message and error
	PrecheckRegister(ctx context.Context, req requestdto.UserRegisterPrecheck, iterCount int) (responsedto.UserRegisterPrecheck, *ucerror.UCError)
	Register(ctx context.Context, req requestdto.UserRegister) *ucerror.UCError
	PrecheckLogin(ctx context.Context, req requestdto.UserLoginPrecheck) (responsedto.UserLoginPrecheck, *ucerror.UCError)
	Login(ctx context.Context, req requestdto.UserLogin) (responsedto.UserLogin, *ucerror.UCError)
	GetProfile(ctx context.Context, userID uint) (responsedto.UserProfile, *ucerror.UCError)
	UpdateUserMetadata(ctx context.Context, userID uint, req requestdto.UserMetadataUpdate) *ucerror.UCError
	VerifyEmail(ctx context.Context, userID uint, userEmail string) *ucerror.UCError
	CheckEmailVerificationCode(ctx context.Context, userId uint, code string) *ucerror.UCError
	SaveProofOfEmailVerification(ctx context.Context, userID uint, req requestdto.UserCheckEmailVerificationCode) *ucerror.UCError
	VerifyPhoneNumber(ctx context.Context, userID uint) *ucerror.UCError
	CheckPhoneNumberVerificationCode(ctx context.Context, userID uint, req requestdto.UserCheckPhoneNumberVerificationCode) *ucerror.UCError
	GenerateAndSaveTelegramSessionIDHash(ctx context.Context, userID uint) (sessionID []byte, ucError *ucerror.UCError)
	PrecheckResetPassword(ctx context.Context, req requestdto.UserResetPasswordPrecheck) (responsedto.UserResetPasswordPrecheck, *ucerror.UCError)
	ResetPassword(ctx context.Context, request requestdto.UserResetPassword) *ucerror.UCError
	MdwVerifyUserJWTToken(ctx context.Context, tokenString string) (userID uint, userUsername string, err *ucerror.UCError)
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
