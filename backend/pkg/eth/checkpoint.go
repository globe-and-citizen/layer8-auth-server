package eth

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

func Backfill[T any](
	ctx context.Context,
	client *ethclient.Client,
	sub EventListener[T],
	start uint64,
) {
	event := sub.ContractABI.Events[sub.EventName]

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(start)),
		Addresses: []common.Address{sub.ContractAddress},
		Topics:    [][]common.Hash{{event.ID}},
	}

	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		log.Error().Msgf("Backfill error: %v", err)
		return
	}

	for _, vLog := range logs {
		err = sub.handleEvent(ctx, client, vLog)
		if err != nil {
			log.Error().Msgf("Backfill handle event error: %v", err)
		}
		saveLastBlock(vLog.BlockNumber)
	}
}
