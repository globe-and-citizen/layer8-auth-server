package postgresRepo

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/models/gormModels"
)

func (r *PostgresRepository) SaveZkSnarksKeyPair(keyPair gormModels.ZkSnarksKeyPair) (uint, error) {
	tx := r.db.Create(&keyPair)

	if tx.Error != nil {
		return 0, tx.Error
	}

	return keyPair.ID, nil
}

func (r *PostgresRepository) GetLatestZkSnarksKeys() (gormModels.ZkSnarksKeyPair, error) {
	var keyPair gormModels.ZkSnarksKeyPair
	err := r.db.Model(&gormModels.ZkSnarksKeyPair{}).Last(&keyPair).Error

	if err != nil {
		return gormModels.ZkSnarksKeyPair{}, err
	}

	return keyPair, nil
}
