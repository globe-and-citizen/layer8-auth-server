package workerUC

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"math/big"
	"time"
)

func (uc *WorkerUsecase) UpdateUsageBalance(ratePerByte big.Int, currTime time.Time) error {
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
			return fmt.Errorf("failed to get traffic updates for client %s: %e", balance.ClientID, err)
		}

		if consumedBytesFloat == 0 {
			continue
		}

		newBilled := ratePerByte
		newBilled.Mul(&newBilled, big.NewInt(int64(consumedBytesFloat)))
		err = uc.postgres.AddClientBalanceDelta(
			uc.ctx,
			balance.ClientID,
			utils.BigIntToDBAmount(&newBilled),
			gormModels.Currency(uc.config.Payment.Currency),
			&timestamp,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
