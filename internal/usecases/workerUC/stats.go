package workerUC

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"math/big"
	"time"
)

func (uc *WorkerUsecase) UpdateUsageBalance(currTime time.Time) error {
	allBalances, err := uc.postgres.GetAllClientBalances(uc.ctx)
	if err != nil {
		return err
	}

	for _, balance := range allBalances {
		timestamp := currTime.UTC()
		consumedBytesFloat, err := uc.influxdb.GetTotalByDateRangeByClient(
			uc.ctx, balance.LastUsageUpdatedAt, timestamp, balance.ClientID,
		)
		if err != nil {
			return fmt.Errorf("failed to get traffic updates for client %s: %w", balance.ClientID, err)
		}

		if consumedBytesFloat == 0 {
			continue
		}

		totalUnbilledBytes, err := utils.DBStringToBigInt(balance.UnbilledUsageBytes)
		if err != nil {
			return fmt.Errorf("failed to convert unbilled usage bytes for client %s: %w", balance.ClientID, err)
		}

		// totalUnbilledBytes = newlyUpdatedConsumption + lastUpdatedUnbilledConsumption
		totalUnbilledBytes.Add(totalUnbilledBytes, big.NewInt(int64(consumedBytesFloat)))

		if utils.StringSliceContains(uc.config.Debug.ClientID, balance.ClientID) {
			uc.logger.Debugf("clientID: %s, newConsumedBytes: %d, lastUnbilledBytes: %s, totalUnbilledBytes: %s",
				balance.ClientID, int64(consumedBytesFloat), balance.UnbilledUsageBytes, totalUnbilledBytes.String())
		}

		// newBilledBalance = quotient(totalUnbilledBytes / BytesPerBillingUnit) * CostPerBillingUnit
		// newUnbilledBytes = remainder(totalUnbilledBytes % BytesPerBillingUnit)
		// DivMod writes the remainder into its last argument, so it must be a fresh big.Int: passing the config
		// value there would overwrite BytesPerBillingUnit and lead to a division by zero on later iterations.
		newBilledBalance, newUnbilledBytes := new(big.Int).DivMod(totalUnbilledBytes, &uc.config.Payment.BytesPerBillingUnit, new(big.Int))
		newBilledBalance.Mul(newBilledBalance, &uc.config.Payment.CostPerBillingUnit)

		if utils.StringSliceContains(uc.config.Debug.ClientID, balance.ClientID) {
			uc.logger.Debugf("clientID: %s, newBilledBalance: %s, newUnbilledBytes: %s", balance.ClientID, newBilledBalance.String(), newUnbilledBytes.String())
		}

		err = uc.postgres.AddClientBalanceDelta(
			uc.ctx,
			balance.ClientID,
			utils.BigIntToDBString(newBilledBalance),
			utils.BigIntToDBString(newUnbilledBytes),
			gormModels.Currency(uc.config.Payment.Currency),
			&timestamp,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
