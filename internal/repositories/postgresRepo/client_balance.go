package postgresRepo

import (
	"database/sql"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"time"
)

func (r *PostgresRepository) GetClientBalance(clientId string) (*gormModels.ClientBalance, error) {
	// TODO: is isolation level higher then the default needed?
	var clientStatistics gormModels.ClientBalance

	err := r.db.Model(&gormModels.ClientBalance{}).
		Where("client_id = ?", clientId).
		First(&clientStatistics).
		Error

	if err != nil {
		return nil, err
	}

	return &clientStatistics, nil
}

func (r *PostgresRepository) UpdateClientBalance(clientId string, newBalance string, status gormModels.AccountStatus, lastUsageUpdatedAt time.Time) error {
	tx := r.db.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})

	err := r.db.Model(&gormModels.ClientBalance{}).
		Where("client_id = ?", clientId).
		Updates(map[string]interface{}{
			"balance_wei":           newBalance,
			"status":                status,
			"last_usage_updated_at": lastUsageUpdatedAt,
		}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) GetAllClientBalances() ([]gormModels.ClientBalance, error) {
	// TODO: is isolation level higher then the default needed?
	var allClientStatistics []gormModels.ClientBalance

	err := r.db.Find(&allClientStatistics).Error
	if err != nil {
		return nil, err
	}

	return allClientStatistics, nil
}

func (r *PostgresRepository) AddClientPaymentReceipt(clientId string, amount string, timestamp time.Time, txHash string) error {
	paymentReceipt := gormModels.ClientPaymentReceipt{
		ClientID:      clientId,
		PaidAmountWei: amount,
		PaidAt:        timestamp,
		TxID:          txHash,
	}

	e := r.db.Create(&paymentReceipt).Error
	if e != nil {
		return e
	}

	return nil
}
