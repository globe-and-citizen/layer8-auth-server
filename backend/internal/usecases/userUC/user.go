package userUC

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"globe-and-citizen/layer8/auth-server/backend/internal/consts"
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/backend/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/backend/pkg/utils"
	"net/http"
)

func (uc *UserUsecase) PrecheckRegister(req requestdto.UserRegisterPrecheck, iterCount int) (responsedto.UserRegisterPrecheck, error) {
	rmSalt := utils.GenerateRandomSalt(consts.SaltSize)

	user := gormModels.User{
		Username:       req.Username,
		Salt:           rmSalt,
		IterationCount: iterCount,
		PublicKey:      []byte{},
	}

	err := uc.postgres.PrecheckUserRegister(user)
	if err != nil {
		return responsedto.UserRegisterPrecheck{}, err
	}

	return responsedto.UserRegisterPrecheck{
		Salt:           rmSalt,
		IterationCount: iterCount,
	}, nil
}

func (uc *UserUsecase) Register(req requestdto.UserRegister) error {
	newUser := gormModels.User{
		Username:  req.Username,
		PublicKey: req.PublicKey,
		StoredKey: req.StoredKey,
		ServerKey: req.ServerKey,
	}

	return uc.postgres.AddUser(newUser)
}

func (uc *UserUsecase) PrecheckLogin(req requestdto.UserLoginPrecheck) (responsedto.UserLoginPrecheck, error) {
	sNonce := utils.GenerateRandomSalt(consts.SaltSize)

	user, err := uc.postgres.GetUserByUsername(req.Username)
	if err != nil {
		return responsedto.UserLoginPrecheck{}, err
	}

	loginPrecheckResp := responsedto.UserLoginPrecheck{
		Salt:      user.Salt,
		IterCount: user.IterationCount,
		Nonce:     req.CNonce + sNonce,
	}

	return loginPrecheckResp, nil
}

func (uc *UserUsecase) Login(req requestdto.UserLogin) (responsedto.UserLogin, error) {
	user, err := uc.postgres.GetUserByUsername(req.Username)
	if err != nil {
		return responsedto.UserLogin{}, err
	}

	storedKeyBytes, err := hex.DecodeString(user.StoredKey)
	if err != nil {
		return responsedto.UserLogin{}, fmt.Errorf("error decoding stored key: %v", err)
	}

	authMessage := fmt.Sprintf("[n=%s,r=%s,s=%s,i=%d,r=%s]", req.Username, req.CNonce, user.Salt, user.IterationCount, req.Nonce)
	authMessageBytes := []byte(authMessage)

	clientSignatureHMAC := hmac.New(sha256.New, storedKeyBytes)
	clientSignatureHMAC.Write(authMessageBytes)
	clientSignature := clientSignatureHMAC.Sum(nil)

	clientProofBytes, err := hex.DecodeString(req.ClientProof)
	if err != nil {
		return responsedto.UserLogin{}, fmt.Errorf("error decoding client proof: %v", err)
	}

	clientKeyBytes, err := utils.XorBytes(clientSignature, clientProofBytes)
	if err != nil {
		return responsedto.UserLogin{}, fmt.Errorf("error performing XOR operation: %v", err)
	}

	clientKeyHash := sha256.Sum256(clientKeyBytes)

	clientKeyHashStr := hex.EncodeToString(clientKeyHash[:])
	if clientKeyHashStr != user.StoredKey {
		return responsedto.UserLogin{}, fmt.Errorf("server failed to authenticate the user")
	}

	serverKeyBytes, err := hex.DecodeString(user.ServerKey)
	if err != nil {
		return responsedto.UserLogin{}, fmt.Errorf("error decoding server key: %v", err)
	}

	serverSignatureHMAC := hmac.New(sha256.New, serverKeyBytes)
	serverSignatureHMAC.Write(authMessageBytes)
	serverSignatureHex := hex.EncodeToString(serverSignatureHMAC.Sum(nil))

	tokenString, err := uc.token.GenerateUserJWTToken(user)
	if err != nil {
		return responsedto.UserLogin{}, fmt.Errorf("error generating token: %v", err)
	}

	return responsedto.UserLogin{
		ServerSignature: serverSignatureHex,
		Token:           tokenString,
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
		Salt:           user.Salt,
		IterationCount: user.IterationCount,
	}, nil
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
