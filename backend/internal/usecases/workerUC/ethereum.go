package workerUC

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/models"

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

func (uc *WorkerUsecase) handleTrafficPaidEvent(event models.TrafficPaidEvent) {
	// Implement the logic to handle the TrafficPaidEvent
	//log.Println("TrafficPaid")
	//log.Println("ClientID:", event.ClientID)
	//log.Println("Payer:", event.Payer.Hex())
	//log.Println("Amount:", event.Amount.String())
	//log.Println("Tx:", vLog.TxHash.Hex())
	//log.Println("Block:", vLog.BlockNumber)

	log.Debug().Msgf("Handling traffic paid event: %v", event)
	//err := uc.postgres.PayClientTrafficUsage(event.ClientID, float64(event.Amount.Int64()))
	//if err != nil {
	//	log.Error().AnErr("failed to update client paid amount", err)
	//}
}
