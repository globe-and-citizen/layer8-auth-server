package postgresRepo

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	gormModels2 "globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"time"
)

type IPostgresRepository interface {
	Migrate()
	IUserRepositories
	IClientRepositories
}

type IUserRepositories interface {
	IUserRepository
	IUserMetadataRepository
	IPhoneNumberVerificationRepository
	IEmailVerificationRepository
	IZKSnarksKeyRepository
}

type IClientRepositories interface {
	IClientRepository
	IClientBalanceRepository
	IClientPaymentReceiptRepository
}

type IUserRepository interface {
	UpdateUser(newUser gormModels2.User) error
	GetUserByID(userId uint) (gormModels2.User, error)
	GetUserByUsername(username string) (gormModels2.User, error)
	GetUserProfile(userID uint) (gormModels2.User, gormModels2.UserMetadata, error)
	PrecheckUserRegister(user gormModels2.User) error
	UpdateUserPassword(username string, storedKey string, serverKey string) error
}

type IClientRepository interface {
	UpdateClient(newClient gormModels2.Client) error
	GetClientByName(name string) (gormModels2.Client, error)
	GetClientByBackendURI(backendURI string) (gormModels2.Client, error)
	IsBackendURIExists(backendURL string) (bool, error)
	GetClientByUsername(username string) (gormModels2.Client, error)
	GetClientProfile(username string) (gormModels2.Client, error)
	PrecheckClientRegister(req gormModels2.Client) error
	SaveX509Certificate(clientID string, certificate string) error
	GetClientByID(id string) (gormModels2.Client, error)
}

type IClientBalanceRepository interface {
	GetClientBalance(clientId string) (*gormModels2.ClientBalance, error)
	UpdateClientBalance(clientId string, newBalance string, status gormModels2.AccountStatus, lastUsageUpdated time.Time) error
	GetAllClientBalances() ([]gormModels2.ClientBalance, error)
}

type IClientPaymentReceiptRepository interface {
	AddClientPaymentReceipt(clientId string, amount string, timestamp time.Time, txHash string) error
}

type IUserMetadataRepository interface {
	GetMetadataByUserID(userID uint) (gormModels2.UserMetadata, error)
	UpdateUserMetadata(userID uint, req requestdto.UserMetadataUpdate) error
}

type IPhoneNumberVerificationRepository interface {
	SavePhoneNumberVerificationData(data gormModels2.PhoneNumberVerificationData) error
	GetPhoneNumberVerificationData(userID uint) (gormModels2.PhoneNumberVerificationData, error)
	SaveProofOfPhoneNumberVerification(
		userID uint,
		phoneNumberVerificationCode string,
		phoneNumberZkProof []byte,
		phoneNumberZkPairID uint,
	) error
}

type IEmailVerificationRepository interface {
	SaveEmailVerificationData(data gormModels2.EmailVerificationData) error
	GetEmailVerificationData(userId uint) (gormModels2.EmailVerificationData, error)
	SaveProofOfEmailVerification(
		userId uint, verificationCode string, emailProof []byte, zkKeyPairId uint,
	) error
	SaveTelegramSessionIDHash(userID uint, sessionID []byte) error
}

type IZKSnarksKeyRepository interface {
	SaveZkSnarksKeyPair(keyPair gormModels2.ZkSnarksKeyPair) (uint, error)
	GetLatestZkSnarksKeys() (gormModels2.ZkSnarksKeyPair, error)
}
