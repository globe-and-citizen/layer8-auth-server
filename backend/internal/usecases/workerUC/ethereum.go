package workerUC

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/models"
	"globe-and-citizen/layer8/auth-server/backend/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/backend/pkg/eth"
	"globe-and-citizen/layer8/auth-server/backend/pkg/utils"

	"github.com/rs/zerolog/log"
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
	log.Debug().Msgf("Handling traffic paid event: %+v", event)

	balance, err := uc.postgres.GetClientBalance(event.Data.ClientID)
	if err != nil {
		log.Error().Msgf("failed to get client balance: %w", err)
		return err
	}

	curBalance, err := utils.DBWeiToBigInt(balance.BalanceWei)
	if err != nil {
		log.Error().Msgf("failed to convert unpaid amount: %w", err)
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
