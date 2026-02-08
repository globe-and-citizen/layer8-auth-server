package eth

import (
	"context"
	"globe-and-citizen/layer8/auth-server/pkg/log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func ConnectToEthereum(wsRPC string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(wsRPC)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func CloseEthereumConnection(client *ethclient.Client) {
	client.Close()
}

// EventListener - Per event listener
type EventListener[T any] struct {
	ContractABI     *abi.ABI
	ContractAddress common.Address
	EventName       string
	EventID         common.Hash
	handler         EventHandlerFunc[T]
}

func NewEventListener[T any](
	contractABI *abi.ABI,
	contractAddr common.Address,
	eventName string,
) EventListener[T] {
	event, ok := contractABI.Events[eventName]
	if !ok {
		panic("event not found in ABI")
	}

	return EventListener[T]{
		ContractABI:     contractABI,
		ContractAddress: contractAddr,
		EventName:       eventName,
		EventID:         event.ID,
	}
}

func (e *EventListener[T]) SetHandler(handler EventHandlerFunc[T]) {
	e.handler = handler
}

func (e *EventListener[T]) Start(ctx context.Context, logger log.ILogger, client *ethclient.Client) {
	if e.handler == nil {
		logger.Errorf(nil, "Handler for %s is nil", e.EventName)
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{e.ContractAddress},
		Topics:    [][]common.Hash{{e.EventID}},
	}

	logs := make(chan types.Log, 100)
	for {
		sub, err := client.SubscribeFilterLogs(ctx, query, logs)
		if err != nil {
			logger.Errorf(err, "Error subscribing to %s", e.EventName)
			time.Sleep(3 * time.Second)
			continue
		}
		logger.Infof("Subscribed to event %s...", e.EventName)

		for {
			select {
			case <-ctx.Done():
				sub.Unsubscribe()
				logger.Infof("Unsubscribed from event %s", e.EventName)
				return
			case err := <-sub.Err():
				logger.Errorf(err, "Error subcribing to %s", e.EventName)
				sub.Unsubscribe()
				time.Sleep(3 * time.Second)
				goto RESUBSCRIBE
			case vLog := <-logs:
				err = e.handleEvent(ctx, client, vLog)
				if err != nil {
					logger.Errorf(err, "Error handling ethereum event %s", e.EventName)
					// todo what to do on handler error?
				}
				saveLastBlock(vLog.BlockNumber)
			}
		}
	RESUBSCRIBE:
	}

}

func (e *EventListener[T]) handleEvent(ctx context.Context, client *ethclient.Client, vLog types.Log) error {
	var eventData EventData[T]

	err := e.ContractABI.UnpackIntoInterface(&eventData.Data, e.EventName, vLog.Data)
	if err != nil {
		return err
	}

	eventData.TxID = vLog.TxHash.Hex()
	// Fetch block header to get the block timestamp (seconds since epoch)
	header, err := client.HeaderByNumber(ctx, new(big.Int).SetUint64(vLog.BlockNumber))
	if err != nil {
		return err
	}
	eventData.TxTimestamp = time.Unix(int64(header.Time), 0)

	err = e.handler(eventData)
	if err != nil {
		return err
	}

	return nil
}
