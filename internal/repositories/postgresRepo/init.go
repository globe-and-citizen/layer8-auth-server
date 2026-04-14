package postgresRepo

import (
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"log"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(config utils.PostgresConfig) IPostgresRepository {
	return &PostgresRepository{db: utils.ConnectDB(config)}
}

func (r *PostgresRepository) Migrate() {
	// Auto migrate tables
	err := r.db.AutoMigrate(
		&gormModels.User{},
		&gormModels.Client{},
		&gormModels.UserMetadata{},
		&gormModels.ClientBalance{},
		&gormModels.ClientPaymentReceipt{},
		&gormModels.EmailVerificationData{},
		&gormModels.PhoneVerificationData{},
		&gormModels.ZkSnarksKeyPair{},
		&gormModels.OAuthAuthorizationCode{},
	)
	if err != nil {
		log.Fatalf("Cannot migrate tables: %v", err)
	}
}

func (r *PostgresRepository) TX() PostgresRepository {
	return PostgresRepository{r.db.Begin()}
}

func (r *PostgresRepository) Commit() error {
	return utils.StackError(r.db.Commit().Error)
}

func (r *PostgresRepository) Rollback() error {
	return utils.StackError(r.db.Rollback().Error)
}
