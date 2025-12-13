package userUsecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/repositories/postgresRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/tokenRepo"
	"globe-and-citizen/layer8/auth-server/utils"
)

type IUserUseCase interface {
	PrecheckRegister(req requestdto.UserRegisterPrecheck, iterCount int) (responsedto.UserRegisterPrecheck, error)
	Register(req requestdto.UserRegister) error
	PrecheckLogin(req requestdto.UserLoginPrecheck) (responsedto.UserLoginPrecheck, error)
	Login(req requestdto.UserLogin) (responsedto.UserLogin, error)
	UpdateUserMetadata(userID uint, req requestdto.UserMetadataUpdate) error
}

type UserUseCase struct {
	postgres postgresRepo.IUserRepositories
	token    tokenRepo.ITokenRepository
}

func NewUserUseCase(postgres postgresRepo.IUserRepositories, token tokenRepo.ITokenRepository) IUserUseCase {
	return &UserUseCase{postgres: postgres, token: token}
}

func (uc *UserUseCase) PrecheckRegister(req requestdto.UserRegisterPrecheck, iterCount int) (responsedto.UserRegisterPrecheck, error) {
	rmSalt := utils.GenerateRandomSalt(consts.SaltSize)

	err := uc.postgres.PrecheckUserRegister(req, rmSalt, iterCount)
	if err != nil {
		return responsedto.UserRegisterPrecheck{}, err
	}

	return responsedto.UserRegisterPrecheck{
		Salt:           rmSalt,
		IterationCount: iterCount,
	}, nil
}

func (uc *UserUseCase) Register(req requestdto.UserRegister) error {
	return uc.postgres.AddUser(req)
}

func (uc *UserUseCase) PrecheckLogin(req requestdto.UserLoginPrecheck) (responsedto.UserLoginPrecheck, error) {
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

func (uc *UserUseCase) Login(req requestdto.UserLogin) (responsedto.UserLogin, error) {
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
