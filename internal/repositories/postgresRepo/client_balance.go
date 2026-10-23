package postgresRepo

import (
	"context"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"time"

	"gorm.io/gorm"
)

func (r *PostgresRepository) GetClientBalance(ctx context.Context, clientId string) (*gormModels.ClientBalance, error) {
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

func (r *PostgresRepository) AddClientBalanceDelta(
	ctx context.Context,
	clientId string,
	delta string,
	currency gormModels.Currency,
	lastUsageUpdatedAt *time.Time,
) error {
	updates := map[string]interface{}{
		"balance_base_units": gorm.Expr("balance_base_units + (?::numeric)", delta),
		"currency":           currency,
		"status": gorm.Expr(
			`CASE
				WHEN balance_base_units + (?::numeric) < 0 THEN ?
				WHEN balance_base_units + (?::numeric) > 0 THEN ?
				ELSE ?
			END`,
			delta, gormModels.AccountOverpaid,
			delta, gormModels.AccountOwing,
			gormModels.AccountZeroed,
		),
	}
	if lastUsageUpdatedAt != nil {
		updates["last_usage_updated_at"] = *lastUsageUpdatedAt
	}

	result := r.db.WithContext(ctx).Model(&gormModels.ClientBalance{}).
		Where("client_id = ?", clientId).
		Updates(updates)
	if result.Error != nil {
		return utils.StackError(result.Error)
	}
	if result.RowsAffected == 0 {
		return utils.StackError(fmt.Errorf("no client balance found for client_id: %s", clientId))
	}
	return nil
}

func (r *PostgresRepository) GetAllClientBalances(ctx context.Context) ([]gormModels.ClientBalance, error) {
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
	currency gormModels.Currency,
	senderWalletAddress string,
	timestamp time.Time,
	txHash string,
) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return utils.StackError(tx.Error)
	}

	updateResult := tx.Model(&gormModels.ClientBalance{}).
		Where("client_id = ?", clientId).
		Updates(map[string]interface{}{
			"balance_base_units": gorm.Expr("balance_base_units - (?::numeric)", amount),
			"currency":           currency,
			"status": gorm.Expr(
				`CASE
					WHEN balance_base_units - (?::numeric) < 0 THEN ?
					WHEN balance_base_units - (?::numeric) > 0 THEN ?
					ELSE ?
				END`,
				amount, gormModels.AccountOverpaid,
				amount, gormModels.AccountOwing,
				gormModels.AccountZeroed,
			),
		})
	if updateResult.Error != nil {
		tx.Rollback()
		return utils.StackError(updateResult.Error)
	}
	if updateResult.RowsAffected == 0 {
		tx.Rollback()
		return utils.StackError(fmt.Errorf("no client balance found for client_id: %s", clientId))
	}

	paymentReceipt := gormModels.ClientPaymentReceipt{
		ClientID:            clientId,
		PaidAmount:          amount,
		Currency:            currency,
		SenderWalletAddress: senderWalletAddress,
		PaidAt:              timestamp,
		TxID:                txHash,
	}

	e := tx.Create(&paymentReceipt).Error
	if e != nil {
		tx.Rollback()
		return utils.StackError(e)
	}

	if err := tx.Commit().Error; err != nil {
		return utils.StackError(err)
	}

	return nil
}
