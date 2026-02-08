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

func Backfill[T any](
	logger log.ILogger,
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
		logger.Error("Backfill error", errors.WithStack(err))
		return
	}

	for _, vLog := range logs {
		err = sub.handleEvent(ctx, client, vLog)
		if err != nil {
			logger.Error("Backfill handle event error", errors.WithStack(err))
		}
		saveLastBlock(vLog.BlockNumber)
	}
}
