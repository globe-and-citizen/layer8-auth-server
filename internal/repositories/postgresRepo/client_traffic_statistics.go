package postgresRepo

import (
	"database/sql"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/models"
	"time"
)

func (r *PostgresRepository) GetClientTrafficStatistics(clientId string) (*models.ClientTrafficStatistics, error) {
	// TODO: is isolation level higher then the default needed?
	var clientStatistics models.ClientTrafficStatistics

	err := r.db.Model(&models.ClientTrafficStatistics{}).
		Where("client_id = ?", clientId).
		First(&clientStatistics).
		Error

	if err != nil {
		return nil, err
	}

	return &clientStatistics, nil
}

func (r *PostgresRepository) AddClientTrafficUsage(clientId string, consumedBytes int, now time.Time) error {
	tx := r.db.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})

	var clientStatistics models.ClientTrafficStatistics
	err := tx.Where("client_id = ?", clientId).
		First(&clientStatistics).
		Error

	if err != nil {
		tx.Rollback()
		return err
	}

	newTrafficBytes := clientStatistics.TotalUsageBytes + consumedBytes
	newUnpaidAmount := clientStatistics.UnpaidAmount + consumedBytes*clientStatistics.RatePerByte

	err = r.db.Model(&models.ClientTrafficStatistics{}).
		Where("client_id = ?", clientId).
		Updates(map[string]interface{}{
			"total_usage_bytes":             newTrafficBytes,
			"unpaid_amount":                 newUnpaidAmount,
			"last_traffic_update_timestamp": now,
		}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) PayClientTrafficUsage(clientId string, amountPaid int) error {
	tx := r.db.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})

	var clientStatistics models.ClientTrafficStatistics
	err := tx.Where("client_id = ?", clientId).
		First(&clientStatistics).
		Error

	if err != nil {
		tx.Rollback()
		return err
	}

	if amountPaid < clientStatistics.UnpaidAmount {
		tx.Rollback()
		return fmt.Errorf("full amount must be paid")
	}

	err = r.db.Model(&models.ClientTrafficStatistics{}).
		Where("client_id = ?", clientId).
		Updates(map[string]interface{}{
			"unpaid_amount": clientStatistics.UnpaidAmount - amountPaid,
		}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) GetAllClientStatistics() ([]models.ClientTrafficStatistics, error) {
	// TODO: is isolation level higher then the default needed?
	var allClientStatistics []models.ClientTrafficStatistics

	err := r.db.Find(&allClientStatistics).Error
	if err != nil {
		return nil, err
	}

	return allClientStatistics, nil
}
