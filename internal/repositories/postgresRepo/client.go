package postgresRepo

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/dto/tmp"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
)

func (r *PostgresRepository) AddClient(req tmp.ClientRegisterDTO, clientUUID string, clientSecret string) error {
	tx := r.db.Begin()

	result := tx.Model(&gormModels.Client{}).
		Where("username = ?", req.Username).
		Updates(map[string]interface{}{
			"name":         req.Name,
			"redirect_uri": req.RedirectURI,
			"backend_uri":  req.BackendURI,
			"id":           clientUUID,
			"secret":       clientSecret,
			"stored_key":   req.StoredKey,
			"server_key":   req.ServerKey,
		})

	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("could not update client: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("no client found with username: %s", req.Username)
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

func (r *PostgresRepository) PrecheckClientRegister(req tmp.ClientRegisterPrecheckDTO, salt string, iterCount int) error {
	client := gormModels.Client{
		Username:       req.Username,
		Salt:           salt,
		IterationCount: iterCount,
	}

	if err := r.db.Create(&client).Error; err != nil {
		return fmt.Errorf("failed to create a new client: %v", err)
	}

	return nil
}
