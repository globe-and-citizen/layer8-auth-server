package postgresRepo

import (
	"fmt"
	gormModels2 "globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"time"
)

func (r *PostgresRepository) UpdateClient(newClient gormModels2.Client) error {
	tx := r.db.Begin()

	result := tx.Model(&gormModels2.Client{}).
		Where("username = ?", newClient.Username).
		Updates(map[string]interface{}{ // why use map here? - because gorm doesn't support updating struct with zero values
			"name":         newClient.Name,
			"redirect_uri": newClient.RedirectURI,
			"backend_uri":  newClient.BackendURI,
			"id":           newClient.ID,
			"secret":       newClient.Secret,
			"stored_key":   newClient.ScramStoredKey,
			"server_key":   newClient.ScramServerKey,
		})

	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("could not update client: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("no client found with username: %s", newClient.Username)
	}

	balance := gormModels2.ClientBalance{
		ClientID:           newClient.ID,
		BalanceWei:         "0",
		Status:             gormModels2.AccountZeroed,
		LastUsageUpdatedAt: time.Now(),
	}

	err := tx.Create(&balance).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("could not create client stats entry: %e", err)
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) GetClientByName(name string) (gormModels2.Client, error) {
	var client gormModels2.Client
	if err := r.db.Where("name = ?", name).First(&client).Error; err != nil {
		return gormModels2.Client{}, err
	}
	return client, nil
}

func (r *PostgresRepository) GetClientByBackendURI(backendURI string) (gormModels2.Client, error) {
	var client gormModels2.Client
	if err := r.db.Where("backend_uri = ?", backendURI).First(&client).Error; err != nil {
		return gormModels2.Client{}, err
	}
	return client, nil
}

func (r *PostgresRepository) IsBackendURIExists(backendURL string) (bool, error) {
	var count int64
	if err := r.db.Model(&gormModels2.Client{}).Where("backend_uri = ?", backendURL).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PostgresRepository) GetClientByUsername(username string) (gormModels2.Client, error) {
	var client gormModels2.Client
	if err := r.db.Where("username = ?", username).First(&client).Error; err != nil {
		return gormModels2.Client{}, err
	}
	return client, nil
}

func (r *PostgresRepository) GetClientProfile(username string) (gormModels2.Client, error) {
	var client gormModels2.Client
	if err := r.db.Where("username = ?", username).First(&client).Error; err != nil {
		return gormModels2.Client{}, err
	}
	return client, nil
}

func (r *PostgresRepository) PrecheckClientRegister(client gormModels2.Client) error {
	if err := r.db.Create(&client).Error; err != nil {
		return fmt.Errorf("failed to create a new client: %v", err)
	}

	return nil
}

func (r *PostgresRepository) SaveX509Certificate(clientID string, certificate string) error {
	return r.db.Model(&gormModels2.Client{}).
		Where("id = ?", clientID).
		Update("x509_certificate_bytes", certificate).
		Error
}

func (r *PostgresRepository) GetClientByID(id string) (gormModels2.Client, error) {
	var client gormModels2.Client
	err := r.db.Where("id = ?", id).First(&client).Error
	if err != nil {
		return gormModels2.Client{}, err
	}

	return client, nil
}
