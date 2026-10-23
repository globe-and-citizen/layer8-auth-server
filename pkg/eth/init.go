package eth

import (
	"context"
	"fmt"
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
	Enabled         bool
	EventABI        *abi.ABI
	ContractAddress common.Address
	EventName       string
	EventID         common.Hash
	FetchTxInput    bool
	handler         EventHandlerFunc[T]
}

func NewEventListener[T any](
	eventABI *abi.ABI,
	contractAddr common.Address,
	eventName string,
	fetchTxInput bool,
	enabled bool,
) EventListener[T] {
	event, ok := eventABI.Events[eventName]
	if !ok {
		panic("event not found in ABI")
	}

	return EventListener[T]{
		Enabled:         enabled,
		EventABI:        eventABI,
		ContractAddress: contractAddr,
		EventName:       eventName,
		EventID:         event.ID,
		FetchTxInput:    fetchTxInput,
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
				saveLastBlock(e.EventName, vLog.BlockNumber)
			}
		}
	RESUBSCRIBE:
	}

}

func (e *EventListener[T]) handleEvent(ctx context.Context, client *ethclient.Client, vLog types.Log) error {
	var eventData EventData[T]

	event, ok := e.EventABI.Events[e.EventName]
	if !ok {
		return fmt.Errorf("event %q not found in ABI", e.EventName)
	}

	if len(vLog.Topics) == 0 || vLog.Topics[0] != event.ID {
		return fmt.Errorf("unexpected event signature")
	}

	// Decode non-indexed parameters from log.Data.
	if err := e.EventABI.UnpackIntoInterface(
		&eventData.Data,
		e.EventName,
		vLog.Data,
	); err != nil {
		return fmt.Errorf("unpack event data: %w", err)
	}

	// Decode indexed parameters from log.Topics.
	indexed := make(abi.Arguments, 0)
	for _, input := range event.Inputs {
		if input.Indexed {
			indexed = append(indexed, input)
		}
	}

	if len(vLog.Topics) != len(indexed)+1 {
		return fmt.Errorf("unexpected number of event topics")
	}

	// ParseTopics writes decoded values into matching struct fields.
	if err := abi.ParseTopics(
		&eventData.Data,
		indexed,
		vLog.Topics[1:],
	); err != nil {
		return fmt.Errorf("unpack indexed topics: %w", err)
	}

	eventData.TxID = vLog.TxHash.Hex()
	// Fetch block header to get the block timestamp (seconds since epoch)
	header, err := client.HeaderByNumber(ctx, new(big.Int).SetUint64(vLog.BlockNumber))
	if err != nil {
		return fmt.Errorf("fetch block header: %w", err)
	}

	eventData.TxTimestamp = time.Unix(int64(header.Time), 0)

	if e.FetchTxInput {
		tx, _, err := client.TransactionByHash(ctx, vLog.TxHash)
		if err != nil {
			return fmt.Errorf("fetch transaction: %w", err)
		}

		eventData.TxInput = tx.Data()
	}

	if err := e.handler(eventData); err != nil {
		return fmt.Errorf("handle event: %w", err)
	}

	return nil
}
