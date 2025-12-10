package postgresRepo

import "globe-and-citizen/layer8/auth-server/internal/models"

func (r *PostgresRepository) SaveZkSnarksKeyPair(keyPair models.ZkSnarksKeyPair) (uint, error) {
	tx := r.db.Create(&keyPair)

	if tx.Error != nil {
		return 0, tx.Error
	}

	return keyPair.ID, nil
}

func (r *PostgresRepository) GetLatestZkSnarksKeys() (models.ZkSnarksKeyPair, error) {
	var keyPair models.ZkSnarksKeyPair
	err := r.db.Model(&models.ZkSnarksKeyPair{}).Last(&keyPair).Error

	if err != nil {
		return models.ZkSnarksKeyPair{}, err
	}

	return keyPair, nil
}
