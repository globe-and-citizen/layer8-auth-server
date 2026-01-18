package workerUC

import (
	"fmt"
	"time"
)

func (uc *WorkerUsecase) UpdateUsageStatistics(ratePerByte float64, now time.Time) error {
	allClientStatistics, err := uc.postgres.GetAllClientStatistics()
	if err != nil {
		return err
	}

	for _, clientStat := range allClientStatistics {
		consumedBytesFloat, err := uc.influxdb.GetTotalByDateRangeByClient(
			clientStat.LastTrafficUpdateTimestamp, now.UTC(), clientStat.ClientId,
		)
		if err != nil {
			return fmt.Errorf("failed to get traffic updates for client %s: %e", clientStat.ClientId, err)
		}

		if consumedBytesFloat == 0 {
			continue
		}
		consumedBytes := int(consumedBytesFloat)

		err = uc.postgres.AddClientTrafficUsage(clientStat.ClientId, consumedBytes, ratePerByte, now.UTC())
		if err != nil {
			return err
		}
	}

	return nil
}
