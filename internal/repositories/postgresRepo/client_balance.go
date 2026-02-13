package postgresRepo

import (
	"context"
	"database/sql"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"time"
)

func (r *PostgresRepository) GetClientBalance(ctx context.Context, clientId string) (*gormModels.ClientBalance, error) {
	// TODO: is isolation level higher then the default needed?
	var clientStatistics gormModels.ClientBalance

	err := r.db.WithContext(ctx).Model(&gormModels.ClientBalance{}).
		Where("client_id = ?", clientId).
		First(&clientStatistics).
		Error

	if err != nil {
		return nil, utils.StackError(err)
	}

	return &clientStatistics, nil
}

func (r *PostgresRepository) UpdateClientBalance(
	ctx context.Context,
	clientId string,
	newBalance string,
	status gormModels.AccountStatus,
	lastUsageUpdatedAt time.Time,
) error {
	tx := r.db.WithContext(ctx).Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})

	err := r.db.WithContext(ctx).Model(&gormModels.ClientBalance{}).
		Where("client_id = ?", clientId).
		Updates(map[string]interface{}{
			"balance_wei":           newBalance,
			"status":                status,
			"last_usage_updated_at": lastUsageUpdatedAt,
		}).Error

	if err != nil {
		tx.Rollback()
		return utils.StackError(err)
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) GetAllClientBalances(ctx context.Context) ([]gormModels.ClientBalance, error) {
	// TODO: is isolation level higher then the default needed?
	var allClientStatistics []gormModels.ClientBalance

	err := r.db.WithContext(ctx).Find(&allClientStatistics).Error
	if err != nil {
		return nil, utils.StackError(err)
	}

	return allClientStatistics, nil
}

func (r *PostgresRepository) AddClientPaymentReceipt(
	ctx context.Context,
	clientId string,
	amount string,
	timestamp time.Time,
	txHash string,
) error {
	paymentReceipt := gormModels.ClientPaymentReceipt{
		ClientID:      clientId,
		PaidAmountWei: amount,
		PaidAt:        timestamp,
		TxID:          txHash,
	}

	e := r.db.WithContext(ctx).Create(&paymentReceipt).Error
	if e != nil {
		return utils.StackError(e)
	}

	return nil
}
