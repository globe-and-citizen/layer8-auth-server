package workerUC

import (
	"encoding/binary"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/models"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/pkg/eth"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
)

func (uc *WorkerUsecase) ListenToEthereumEvents() {
	// Set up handler for all events
	uc.ethereum.SetAllHandlers(uc.handleTrafficPaidEvent, uc.handleCNCTokenDeposited)

	// Backfill past events
	uc.ethereum.BackfillAll(uc.ctx)

	// Start listening to new events
	uc.ethereum.ListenToAllEvents(uc.ctx)
}

func (uc *WorkerUsecase) handleTrafficPaidEvent(event eth.EventData[models.TrafficPaidEvent]) error {
	uc.logger.Infof("Handling traffic paid event: %+v", event)

	balance, err := uc.postgres.GetClientBalance(uc.ctx, event.Data.ClientID)
	if err != nil {
		uc.logger.Error("failed to get client balance", err)
		return err
	}

	curBalance, err := utils.DBAmountToBigInt(balance.BalanceBaseUnits)
	if err != nil {
		uc.logger.Error("failed to convert unpaid amount", err)
		return err
	}
	curBalance = curBalance.Sub(curBalance, event.Data.Amount)

	var status gormModels.AccountStatus
	err = uc.postgres.UpdateClientBalance(
		uc.ctx,
		event.Data.ClientID,
		utils.BigIntToDBAmount(curBalance),
		gormModels.POL,
		status.GetStatus(curBalance),
		balance.LastUsageUpdatedAt,
	)
	if err != nil {
		return err
	}

	amount := utils.BigIntToDBAmount(event.Data.Amount)
	err = uc.postgres.AddClientPaymentReceipt(
		uc.ctx,
		event.Data.ClientID,
		amount,
		gormModels.POL,
		event.Data.Payer.Hex(),
		event.TxTimestamp,
		event.TxID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (uc *WorkerUsecase) handleCNCTokenDeposited(event eth.EventData[models.CNCTokenDepositedEvent]) error {
	factureID, err := ExtractFactureID(event.TxInput)
	if err != nil {
		return err
	}

	if factureID == "" {
		uc.logger.Warn("Facture ID is empty")
		return nil // Ignore plain deposits
	}

	uc.logger.Infof(
		"FactureID: %s, Depositor: %s, Token: %s, Amount: %s",
		factureID,
		event.Data.Depositor.Hex(),
		event.Data.Token.Hex(),
		event.Data.Amount.String(),
	)

	// TODO: Validate invoice and update payment status.
	balance, err := uc.postgres.GetClientBalance(uc.ctx, factureID)
	if err != nil {
		uc.logger.Error("failed to get client balance", err)
		return err
	}

	curBalance, err := utils.DBAmountToBigInt(balance.BalanceBaseUnits)
	if err != nil {
		uc.logger.Error("failed to convert unpaid amount", err)
		return err
	}
	curBalance = curBalance.Sub(curBalance, event.Data.Amount)

	var status gormModels.AccountStatus
	err = uc.postgres.UpdateClientBalance(
		uc.ctx,
		factureID,
		utils.BigIntToDBAmount(curBalance),
		gormModels.Currency(uc.config.Payment.Currency),
		status.GetStatus(curBalance),
		balance.LastUsageUpdatedAt,
	)
	if err != nil {
		return err
	}

	amount := utils.BigIntToDBAmount(event.Data.Amount)
	err = uc.postgres.AddClientPaymentReceipt(
		uc.ctx,
		factureID,
		amount,
		gormModels.USDC,
		event.Data.Depositor.Hex(),
		event.TxTimestamp,
		event.TxID,
	)
	if err != nil {
		return err
	}

	return nil
}

func ExtractFactureID(input []byte) (string, error) {
	if len(input) <= 70 {
		return "", nil // Plain deposit, not a widget payment
	}

	length := int(binary.BigEndian.Uint16(input[68:70]))

	if length == 0 || len(input) < 70+length {
		return "", fmt.Errorf("invalid facture ID length")
	}

	return string(input[70 : 70+length]), nil
}
