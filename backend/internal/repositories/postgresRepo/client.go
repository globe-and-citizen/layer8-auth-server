package postgresRepo

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/backend/internal/models/gormModels"
)

func (r *PostgresRepository) UpdateClient(newClient gormModels.Client) error {
	tx := r.db.Begin()

	result := tx.Model(&gormModels.Client{}).
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

	stats := gormModels.ClientTrafficStatistics{
		ClientId: newClient.ID,
	}

	err := tx.Create(&stats).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("could not create client stats entry: %e", err)
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) GetClientByName(name string) (gormModels.Client, error) {
	var client gormModels.Client
	if err := r.db.Where("name = ?", name).First(&client).Error; err != nil {
		return gormModels.Client{}, err
	}
	return client, nil
}

func (r *PostgresRepository) GetClientByBackendURI(backendURI string) (gormModels.Client, error) {
	var client gormModels.Client
	if err := r.db.Where("backend_uri = ?", backendURI).First(&client).Error; err != nil {
		return gormModels.Client{}, err
	}
	return client, nil
}

func (r *PostgresRepository) IsBackendURIExists(backendURL string) (bool, error) {
	var count int64
	if err := r.db.Model(&gormModels.Client{}).Where("backend_uri = ?", backendURL).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PostgresRepository) GetClientByUsername(username string) (gormModels.Client, error) {
	var client gormModels.Client
	if err := r.db.Where("username = ?", username).First(&client).Error; err != nil {
		return gormModels.Client{}, err
	}
	return client, nil
}

func (r *PostgresRepository) GetClientProfile(username string) (gormModels.Client, error) {
	var client gormModels.Client
	if err := r.db.Where("username = ?", username).First(&client).Error; err != nil {
		return gormModels.Client{}, err
	}
	return client, nil
}

func (r *PostgresRepository) PrecheckClientRegister(client gormModels.Client) error {
	if err := r.db.Create(&client).Error; err != nil {
		return fmt.Errorf("failed to create a new client: %v", err)
	}

	return nil
}

func (r *PostgresRepository) SaveX509Certificate(clientID string, certificate string) error {
	return r.db.Model(&gormModels.Client{}).
		Where("id = ?", clientID).
		Update("x509_certificate_bytes", certificate).
		Error
}

func (r *PostgresRepository) GetClientByID(id string) (gormModels.Client, error) {
	var client gormModels.Client
	err := r.db.Where("id = ?", id).First(&client).Error
	if err != nil {
		return gormModels.Client{}, err
	}

	return client, nil
}
