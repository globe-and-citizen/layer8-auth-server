package postgresRepo

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/tmp"
	"globe-and-citizen/layer8/auth-server/internal/models"
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
	AddUser(req requestdto.UserRegister) error
	FindUserByID(userId uint) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
	GetUserProfile(userID uint) (models.User, models.UserMetadata, error)
	PrecheckUserRegister(req requestdto.UserRegisterPrecheck, salt string, iterCount int) error
	UpdateUserPassword(username string, storedKey string, serverKey string) error
}

type IClientRepository interface {
	AddClient(req tmp.ClientRegisterDTO, clientUUID string, clientSecret string) error
	GetClientByName(name string) (models.Client, error)
	GetClientByBackendURI(backendURI string) (models.Client, error)
	IsBackendURIExists(backendURL string) (bool, error)
	GetClientByUsername(username string) (models.Client, error)
	GetClientProfile(username string) (models.Client, error)
	PrecheckClientRegister(req tmp.ClientRegisterPrecheckDTO, salt string, iterCount int) error
}

type IClientTrafficStatisticsRepository interface {
	GetClientTrafficStatistics(clientId string) (*models.ClientTrafficStatistics, error)
	AddClientTrafficUsage(clientId string, consumedBytes int, now time.Time) error
	PayClientTrafficUsage(clientId string, amountPaid int) error
	GetAllClientStatistics() ([]models.ClientTrafficStatistics, error)
}

type IUserMetadataRepository interface {
	GetMetadataByUserID(userID uint) (models.UserMetadata, error)
	UpdateUserMetadata(userID uint, req requestdto.UserMetadataUpdate) error
}

type IPhoneNumberVerificationRepository interface {
	SavePhoneNumberVerificationData(data models.PhoneNumberVerificationData) error
	GetPhoneNumberVerificationData(userID uint) (models.PhoneNumberVerificationData, error)
	SaveProofOfPhoneNumberVerification(
		userID uint,
		phoneNumberVerificationCode string,
		phoneNumberZkProof []byte,
		phoneNumberZkPairID uint,
	) error
}

type IEmailVerificationRepository interface {
	SaveEmailVerificationData(data models.EmailVerificationData) error
	GetEmailVerificationData(userId uint) (models.EmailVerificationData, error)
	SaveProofOfEmailVerification(
		userId uint, verificationCode string, emailProof []byte, zkKeyPairId uint,
	) error
	SaveTelegramSessionIDHash(userID uint, sessionID []byte) error
}

type IZKSnarksKeyRepository interface {
	SaveZkSnarksKeyPair(keyPair models.ZkSnarksKeyPair) (uint, error)
	GetLatestZkSnarksKeys() (models.ZkSnarksKeyPair, error)
}
