package eth

import (
	"context"
	"globe-and-citizen/layer8/auth-server/pkg/log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

const backfillMaxBlockRange uint64 = 10

func Backfill[T any](
	logger log.ILogger,
	ctx context.Context,
	client *ethclient.Client,
	sub EventListener[T],
	start uint64,
) {
	if sub.handler == nil {
		logger.Errorf(nil, "Handler for %s is nil", sub.EventName)
		return
	}

	fromBlock := start
	// Resume from the checkpoint block itself. If the process stopped after
	// committing some events in that block but before persisting checkpoint,
	// replaying the block relies on idempotent receipt handling.

	event := sub.EventABI.Events[sub.EventName]
	latestBlock, err := client.BlockNumber(ctx)
	if err != nil {
		logger.Error("Backfill error", errors.WithStack(err))
		return
	}

	if fromBlock > latestBlock {
		return
	}

	for chunkStart := fromBlock; chunkStart <= latestBlock; chunkStart += backfillMaxBlockRange {
		chunkEnd := chunkStart + backfillMaxBlockRange - 1
		if chunkEnd > latestBlock {
			chunkEnd = latestBlock
		}

		query := ethereum.FilterQuery{
			FromBlock: new(big.Int).SetUint64(chunkStart),
			ToBlock:   new(big.Int).SetUint64(chunkEnd),
			Addresses: []common.Address{sub.ContractAddress},
			Topics:    [][]common.Hash{{event.ID}},
		}

		logs, filterErr := client.FilterLogs(ctx, query)
		if filterErr != nil {
			logger.Error("Backfill error", errors.WithStack(filterErr))
			return
		}

		for _, vLog := range logs {
			err = sub.handleEvent(ctx, client, vLog)
			if err != nil {
				logger.Error("Backfill handle event error", errors.WithStack(err))
				return
			}
			saveLastBlock(sub.EventName, vLog.BlockNumber)
		}

		// Mark chunk as processed even when no matching logs are found.
		saveLastBlock(sub.EventName, chunkEnd)
	}
}
