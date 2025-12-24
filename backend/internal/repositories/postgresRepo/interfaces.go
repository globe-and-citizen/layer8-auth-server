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
	IClientTrafficStatisticsRepository
}

type IUserRepository interface {
	AddUser(newUser gormModels.User) error
	FindUserByID(userId uint) (gormModels.User, error)
	GetUserByUsername(username string) (gormModels.User, error)
	GetUserProfile(userID uint) (gormModels.User, gormModels.UserMetadata, error)
	PrecheckUserRegister(user gormModels.User) error
	UpdateUserPassword(username string, storedKey string, serverKey string) error
}

type IClientRepository interface {
	AddClient(newClient gormModels.Client) error
	GetClientByName(name string) (gormModels.Client, error)
	GetClientByBackendURI(backendURI string) (gormModels.Client, error)
	IsBackendURIExists(backendURL string) (bool, error)
	GetClientByUsername(username string) (gormModels.Client, error)
	GetClientProfile(username string) (gormModels.Client, error)
	PrecheckClientRegister(req gormModels.Client) error
	SaveX509Certificate(clientID string, certificate string) error
	GetClientByID(id string) (gormModels.Client, error)
}

type IClientTrafficStatisticsRepository interface {
	GetClientTrafficStatistics(clientId string) (*gormModels.ClientTrafficStatistics, error)
	AddClientTrafficUsage(clientId string, consumedBytes int, now time.Time) error
	PayClientTrafficUsage(clientId string, amountPaid int) error
	GetAllClientStatistics() ([]gormModels.ClientTrafficStatistics, error)
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
