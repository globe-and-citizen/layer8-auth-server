package postgresRepo

import (
	"database/sql"
	"fmt"
	"globe-and-citizen/layer8/auth-server/backend/internal/models/gormModels"
	"time"
)

func (r *PostgresRepository) GetClientTrafficStatistics(clientId string) (*gormModels.ClientTrafficStatistics, error) {
	// TODO: is isolation level higher then the default needed?
	var clientStatistics gormModels.ClientTrafficStatistics

	err := r.db.Model(&gormModels.ClientTrafficStatistics{}).
		Where("client_id = ?", clientId).
		First(&clientStatistics).
		Error

	if err != nil {
		return nil, err
	}

	return &clientStatistics, nil
}

func (r *PostgresRepository) AddClientTrafficUsage(clientId string, consumedBytes int, ratePerByte float64, now time.Time) error {
	tx := r.db.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})

	var clientStatistics gormModels.ClientTrafficStatistics
	err := tx.Where("client_id = ?", clientId).
		First(&clientStatistics).
		Error

	if err != nil {
		tx.Rollback()
		return err
	}

	newTrafficBytes := clientStatistics.TotalUsageBytes + consumedBytes
	newUnpaidAmount := float64(clientStatistics.UnpaidAmount) + float64(consumedBytes)*ratePerByte

	err = r.db.Model(&gormModels.ClientTrafficStatistics{}).
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

func (r *PostgresRepository) PayClientTrafficUsage(clientId string, amountPaid float64) error {
	tx := r.db.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})

	var clientStatistics gormModels.ClientTrafficStatistics
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

	err = r.db.Model(&gormModels.ClientTrafficStatistics{}).
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

func (r *PostgresRepository) GetAllClientStatistics() ([]gormModels.ClientTrafficStatistics, error) {
	// TODO: is isolation level higher then the default needed?
	var allClientStatistics []gormModels.ClientTrafficStatistics

	err := r.db.Find(&allClientStatistics).Error
	if err != nil {
		return nil, err
	}

	return allClientStatistics, nil
}
