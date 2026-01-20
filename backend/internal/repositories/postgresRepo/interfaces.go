package postgresRepo

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/backend/internal/models/gormModels"
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
	UpdateUser(newUser gormModels.User) error
	GetUserByID(userId uint) (gormModels.User, error)
	GetUserByUsername(username string) (gormModels.User, error)
	GetUserProfile(userID uint) (gormModels.User, gormModels.UserMetadata, error)
	PrecheckUserRegister(user gormModels.User) error
	UpdateUserPassword(username string, storedKey string, serverKey string) error
}

type IClientRepository interface {
	UpdateClient(newClient gormModels.Client) error
	GetClientByName(name string) (gormModels.Client, error)
	GetClientByBackendURI(backendURI string) (gormModels.Client, error)
	IsBackendURIExists(backendURL string) (bool, error)
	GetClientByUsername(username string) (gormModels.Client, error)
	GetClientProfile(username string) (gormModels.Client, error)
	PrecheckClientRegister(req gormModels.Client) error
	SaveX509Certificate(clientID string, certificate string) error
	GetClientByID(id string) (gormModels.Client, error)
}

type IClientBalanceRepository interface {
	GetClientBalance(clientId string) (*gormModels.ClientBalance, error)
	UpdateClientBalance(clientId string, newBalance string, status gormModels.AccountStatus, lastUsageUpdated time.Time) error
	GetAllClientBalances() ([]gormModels.ClientBalance, error)
}

type IClientPaymentReceiptRepository interface {
	AddClientPaymentReceipt(clientId string, amount string, timestamp time.Time, txHash string) error
}

type IUserMetadataRepository interface {
	GetMetadataByUserID(userID uint) (gormModels.UserMetadata, error)
	UpdateUserMetadata(userID uint, req requestdto.UserMetadataUpdate) error
}

type IPhoneNumberVerificationRepository interface {
	SavePhoneNumberVerificationData(data gormModels.PhoneNumberVerificationData) error
	GetPhoneNumberVerificationData(userID uint) (gormModels.PhoneNumberVerificationData, error)
	SaveProofOfPhoneNumberVerification(
		userID uint,
		phoneNumberVerificationCode string,
		phoneNumberZkProof []byte,
		phoneNumberZkPairID uint,
	) error
}

type IEmailVerificationRepository interface {
	SaveEmailVerificationData(data gormModels.EmailVerificationData) error
	GetEmailVerificationData(userId uint) (gormModels.EmailVerificationData, error)
	SaveProofOfEmailVerification(
		userId uint, verificationCode string, emailProof []byte, zkKeyPairId uint,
	) error
	SaveTelegramSessionIDHash(userID uint, sessionID []byte) error
}

type IZKSnarksKeyRepository interface {
	SaveZkSnarksKeyPair(keyPair gormModels.ZkSnarksKeyPair) (uint, error)
	GetLatestZkSnarksKeys() (gormModels.ZkSnarksKeyPair, error)
}
