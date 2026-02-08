package workerUC

import (
	"globe-and-citizen/layer8/auth-server/internal/models"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/pkg/eth"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
)

func (uc *WorkerUsecase) ListenToEthereumEvents() {
	// Set up handler for all events
	uc.ethereum.SetAllHandlers(uc.handleTrafficPaidEvent)

	// Backfill past events
	uc.ethereum.BackfillAll(uc.ctx)

	// Start listening to new events
	uc.ethereum.ListenToAllEvents(uc.ctx)
}

func (uc *WorkerUsecase) handleTrafficPaidEvent(event eth.EventData[models.TrafficPaidEvent]) error {
	uc.logger.Infof("Handling traffic paid event: %+v", event)

	balance, err := uc.postgres.GetClientBalance(event.Data.ClientID)
	if err != nil {
		uc.logger.Error("failed to get client balance", err)
		return err
	}

	curBalance, err := utils.DBWeiToBigInt(balance.BalanceWei)
	if err != nil {
		uc.logger.Error("failed to convert unpaid amount", err)
		return err
	}
	curBalance = curBalance.Sub(curBalance, event.Data.Amount)

	var status gormModels.AccountStatus
	err = uc.postgres.UpdateClientBalance(
		event.Data.ClientID,
		utils.BigIntToDBWei(curBalance),
		status.GetStatus(curBalance),
		balance.LastUsageUpdatedAt,
	)
	if err != nil {
		return err
	}

	amount := utils.BigIntToDBWei(event.Data.Amount)
	err = uc.postgres.AddClientPaymentReceipt(event.Data.ClientID, amount, event.TxTimestamp, event.TxID)
	if err != nil {
		return err
	}

	return nil
}
