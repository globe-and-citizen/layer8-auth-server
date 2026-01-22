package workerUC

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"math/big"
	"time"
)

func (uc *WorkerUsecase) UpdateUsageBalance(ratePerByteWei *big.Int, now time.Time) error {
	allBalances, err := uc.postgres.GetAllClientBalances()
	if err != nil {
		return err
	}

	for _, balance := range allBalances {
		timestamp := now.UTC()
		consumedBytesFloat, err := uc.influxdb.GetTotalByDateRangeByClient(
			balance.LastUsageUpdatedAt, timestamp, balance.ClientID,
		)
		if err != nil {
			return fmt.Errorf("failed to get traffic updates for client %s: %e", balance.ClientID, err)
		}

		if consumedBytesFloat == 0 {
			continue
		}

		curBalance, err := utils.DBWeiToBigInt(balance.BalanceWei)
		if err != nil {
			return err
		}

		newBilled := ratePerByteWei.Mul(ratePerByteWei, big.NewInt(int64(consumedBytesFloat)))
		curBalance.Add(curBalance, newBilled)

		var status gormModels.AccountStatus
		err = uc.postgres.UpdateClientBalance(
			balance.ClientID,
			utils.BigIntToDBWei(curBalance),
			status.GetStatus(curBalance),
			timestamp,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
