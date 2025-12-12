package postgresRepo

import (
	"globe-and-citizen/layer8/auth-server/config"
	"globe-and-citizen/layer8/auth-server/internal/models"
	"globe-and-citizen/layer8/auth-server/utils"
	"log"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(config config.PostgresConfig) IPostgresRepository {
	return &PostgresRepository{db: utils.ConnectDB(config)}
}

func (r *PostgresRepository) Migrate() {
	// Auto migrate tables
	err := r.db.AutoMigrate(
		&models.User{},
		&models.Client{},
		&models.UserMetadata{},
		&models.ClientTrafficStatistics{},
		&models.EmailVerificationData{},
		&models.PhoneNumberVerificationData{},
		&models.ZkSnarksKeyPair{},
	)
	if err != nil {
		log.Fatalf("Cannot migrate tables: %v", err)
	}
}

func (r *PostgresRepository) TX() PostgresRepository {
	return PostgresRepository{r.db.Begin()}
}

func (r *PostgresRepository) Commit() error {
	return r.db.Commit().Error
}

func (r *PostgresRepository) Rollback() error {
	return r.db.Rollback().Error
}
