package userUC

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	scram2 "globe-and-citizen/layer8/auth-server/pkg/scram"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"net/http"
)

func (uc *UserUsecase) PrecheckRegister(req requestdto.UserRegisterPrecheck, iterCount int) (responsedto.UserRegisterPrecheck, error) {
	registerMsg := scram2.CreateServerRegisterFirstMessage(iterCount)

	user := gormModels.User{
		Username:            req.Username,
		ScramSalt:           registerMsg.Salt,
		ScramIterationCount: iterCount,
		PublicKey:           []byte{},
	}

	err := uc.postgres.PrecheckUserRegister(user)
	if err != nil {
		return responsedto.UserRegisterPrecheck{}, err
	}

	return responsedto.UserRegisterPrecheck{
		ServerRegisterFirstMessage: registerMsg,
	}, nil
}

func (uc *UserUsecase) Register(req requestdto.UserRegister) error {
	newUser := gormModels.User{
		Username:       req.Username,
		PublicKey:      req.PublicKey,
		ScramStoredKey: req.StoredKey,
		ScramServerKey: req.ServerKey,
	}

	return uc.postgres.UpdateUser(newUser)
}

func (uc *UserUsecase) PrecheckLogin(req requestdto.UserLoginPrecheck) (responsedto.UserLoginPrecheck, error) {
	user, err := uc.postgres.GetUserByUsername(req.Username)
	if err != nil {
		return responsedto.UserLoginPrecheck{}, err
	}

	loginPrecheckResp := responsedto.UserLoginPrecheck{
		ServerLoginFirstMessage: scram2.CreateServerLoginFirstMessage(user.ScramSalt, user.ScramIterationCount, req.ClientLoginFirstMessage),
	}

	return loginPrecheckResp, nil
}

func (uc *UserUsecase) Login(req requestdto.UserLogin) (responsedto.UserLogin, error) {
	user, err := uc.postgres.GetUserByUsername(req.Username)
	if err != nil {
		return responsedto.UserLogin{}, err
	}

	tokenString, err := uc.token.GenerateUserJWTToken(user)
	if err != nil {
		return responsedto.UserLogin{}, fmt.Errorf("error generating token: %v", err)
	}

	serverFinalMsg, err := scram2.CreateServerLoginFinalMessage(req.ClientLoginFinalMessage, req.CNonce, user.ScramSalt,
		user.ScramIterationCount, user.ScramStoredKey, user.ScramServerKey)
	if err != nil {
		return responsedto.UserLogin{}, fmt.Errorf("error creating server final message: %v", err)
	}

	return responsedto.UserLogin{
		ServerLoginFinalMessage: serverFinalMsg,
		Token:                   tokenString,
	}, nil
}

func (uc *UserUsecase) GetProfile(userID uint) (responsedto.UserProfile, error) {
	user, metadata, err := uc.postgres.GetUserProfile(userID)
	if err != nil {
		return responsedto.UserProfile{}, err
	}

	return responsedto.UserProfile{
		Username:            user.Username,
		DisplayName:         metadata.DisplayName,
		Bio:                 metadata.Bio,
		Color:               metadata.Color,
		EmailVerified:       metadata.IsEmailVerified,
		PhoneNumberVerified: metadata.IsPhoneNumberVerified,
	}, nil
}

func (uc *UserUsecase) PrecheckResetPassword(req requestdto.UserResetPasswordPrecheck) (responsedto.UserResetPasswordPrecheck, error) {
	user, err := uc.postgres.GetUserByUsername(req.Username)
	if err != nil {
		return responsedto.UserResetPasswordPrecheck{}, err
	}

	return responsedto.UserResetPasswordPrecheck{
		ServerRegisterFirstMessage: scram2.ServerRegisterFirstMessage{
			Salt:           user.ScramSalt,
			IterationCount: user.ScramIterationCount,
		}}, nil
}

func (uc *UserUsecase) ResetPassword(req requestdto.UserResetPassword) (int, string, error) {
	user, err := uc.postgres.GetUserByUsername(req.Username)
	if err != nil {
		return http.StatusNotFound, "User does not exist!", err
	}

	err = utils.ValidateSignature("Sign-in with Layer8", req.Signature, user.PublicKey)
	if err != nil {
		return http.StatusBadRequest, "Signature is invalid!", err
	}

	err = uc.postgres.UpdateUserPassword(user.Username, req.StoredKey, req.ServerKey)
	if err != nil {
		return http.StatusInternalServerError, "Internal error: failed to update user", err
	}

	return http.StatusOK, "", nil
}
