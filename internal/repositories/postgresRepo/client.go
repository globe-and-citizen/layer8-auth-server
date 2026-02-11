package postgresRepo

import (
	"context"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"time"
)

func (r *PostgresRepository) UpdateClient(ctx context.Context, newClient gormModels.Client) error {
	tx := r.db.WithContext(ctx).WithContext(ctx).Begin()

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
		return utils.StackError(fmt.Errorf("could not update client: %w", result.Error))
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return utils.StackError(fmt.Errorf("no client found with username: %s", newClient.Username))
	}

	balance := gormModels.ClientBalance{
		ClientID:           newClient.ID,
		BalanceWei:         "0",
		Status:             gormModels.AccountZeroed,
		LastUsageUpdatedAt: time.Now().UTC(),
	}

	err := tx.Create(&balance).Error
	if err != nil {
		tx.Rollback()
		return utils.StackError(fmt.Errorf("could not create client stats entry: %w", err))
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) GetClientByName(ctx context.Context, name string) (gormModels.Client, error) {
	var client gormModels.Client
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&client).Error; err != nil {
		return gormModels.Client{}, utils.StackError(err)
	}
	return client, nil
}

func (r *PostgresRepository) GetClientByBackendURI(ctx context.Context, backendURI string) (gormModels.Client, error) {
	var client gormModels.Client
	if err := r.db.WithContext(ctx).Where("backend_uri = ?", backendURI).First(&client).Error; err != nil {
		return gormModels.Client{}, utils.StackError(err)
	}
	return client, nil
}

func (r *PostgresRepository) IsBackendURIExists(ctx context.Context, backendURL string) (bool, error) {
	var count int64
	if err := r.db.
		WithContext(ctx).
		Model(&gormModels.Client{}).
		Where("backend_uri = ?", backendURL).
		Count(&count).Error; err != nil {
		return false, utils.StackError(err)
	}
	return count > 0, nil
}

func (r *PostgresRepository) GetClientByUsername(ctx context.Context, username string) (gormModels.Client, error) {
	var client gormModels.Client
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&client).Error; err != nil {
		return gormModels.Client{}, utils.StackError(err)
	}
	return client, nil
}

func (r *PostgresRepository) GetClientProfile(ctx context.Context, username string) (gormModels.Client, error) {
	var client gormModels.Client
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&client).Error; err != nil {
		return gormModels.Client{}, utils.StackError(err)
	}
	return client, nil
}

func (r *PostgresRepository) PrecheckClientRegister(ctx context.Context, client gormModels.Client) error {
	if err := r.db.WithContext(ctx).Create(&client).Error; err != nil {
		return utils.StackError(fmt.Errorf("failed to create a new client: %w", err))
	}

	return nil
}

func (r *PostgresRepository) SaveX509Certificate(ctx context.Context, clientID string, certificate string) error {
	return utils.StackError(r.db.WithContext(ctx).Model(&gormModels.Client{}).
		Where("id = ?", clientID).
		Update("x509_certificate_bytes", certificate).
		Error)
}

func (r *PostgresRepository) GetClientByID(ctx context.Context, id string) (gormModels.Client, error) {
	var client gormModels.Client
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&client).Error
	if err != nil {
		return gormModels.Client{}, utils.StackError(err)
	}

	return client, nil
}
