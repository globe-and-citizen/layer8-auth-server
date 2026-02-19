package postgresRepo

import (
	"context"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
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
	UpdateUser(ctx context.Context, newUser gormModels.User) error
	GetUserByID(ctx context.Context, userId uint) (*gormModels.User, error)
	GetUserByUsername(ctx context.Context, username string) (*gormModels.User, error)
	GetUserProfile(ctx context.Context, userID uint) (*gormModels.User, *gormModels.UserMetadata, error)
	CreateUser(ctx context.Context, user gormModels.User) error
	UpdateUserPassword(ctx context.Context, username string, storedKey string, serverKey string) error
	IsUserIDExists(ctx context.Context, id uint) (bool, error)
}

type IClientRepository interface {
	UpdateClient(ctx context.Context, newClient gormModels.Client) error
	GetClientByName(ctx context.Context, name string) (*gormModels.Client, error)
	GetClientByBackendURI(ctx context.Context, backendURI string) (*gormModels.Client, error)
	IsBackendURIExists(ctx context.Context, backendURL string) (bool, error)
	GetClientByUsername(ctx context.Context, username string) (*gormModels.Client, error)
	GetClientProfile(ctx context.Context, username string) (*gormModels.Client, error)
	CreateClient(ctx context.Context, req gormModels.Client) error
	SaveX509Certificate(ctx context.Context, clientID string, certificate string) error
	GetClientByID(ctx context.Context, id string) (*gormModels.Client, error)
	IsClientIDExists(ctx context.Context, id string) (bool, error)
}

type IClientBalanceRepository interface {
	GetClientBalance(ctx context.Context, clientId string) (*gormModels.ClientBalance, error)
	UpdateClientBalance(ctx context.Context, clientId string, newBalance string, status gormModels.AccountStatus, lastUsageUpdated time.Time) error
	GetAllClientBalances(ctx context.Context) ([]gormModels.ClientBalance, error)
}

type IClientPaymentReceiptRepository interface {
	AddClientPaymentReceipt(ctx context.Context, clientId string, amount string, timestamp time.Time, txHash string) error
}

type IUserMetadataRepository interface {
	GetMetadataByUserID(ctx context.Context, userID uint) (*gormModels.UserMetadata, error)
	UpdateUserMetadata(ctx context.Context, userID uint, req requestdto.UserMetadataUpdate) error
}

type IPhoneNumberVerificationRepository interface {
	SavePhoneNumberVerificationData(ctx context.Context, data gormModels.PhoneNumberVerificationData) error
	GetPhoneNumberVerificationData(ctx context.Context, userID uint) (*gormModels.PhoneNumberVerificationData, error)
	SaveProofOfPhoneNumberVerification(
		ctx context.Context,
		userID uint,
		phoneNumberVerificationCode string,
		phoneNumberZkProof []byte,
		phoneNumberZkPairID uint,
	) error
}

type IEmailVerificationRepository interface {
	SaveEmailVerificationData(ctx context.Context, data gormModels.EmailVerificationData) error
	GetEmailVerificationData(ctx context.Context, userId uint) (*gormModels.EmailVerificationData, error)
	SaveProofOfEmailVerification(
		ctx context.Context, userId uint, verificationCode string, emailProof []byte, zkKeyPairId uint,
	) error
	SaveTelegramSessionIDHash(ctx context.Context, userID uint, sessionID []byte) error
}

type IZKSnarksKeyRepository interface {
	SaveZkSnarksKeyPair(keyPair gormModels.ZkSnarksKeyPair) (uint, error)
	GetLatestZkSnarksKeys() (*gormModels.ZkSnarksKeyPair, error)
}
