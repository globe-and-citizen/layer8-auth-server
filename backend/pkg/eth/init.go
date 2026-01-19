package eth

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

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
	handler         func(T)
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

func (e *EventListener[T]) SetHandler(handler func(T)) {
	e.handler = handler
}

func (e *EventListener[T]) Start(ctx context.Context, client *ethclient.Client) {
	if e.handler == nil {
		log.Fatal().Msgf("Handler for %s is nil", e.EventName)
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{e.ContractAddress},
		Topics:    [][]common.Hash{{e.EventID}},
	}

	logs := make(chan types.Log, 100)
	for {
		sub, err := client.SubscribeFilterLogs(ctx, query, logs)
		if err != nil {
			log.Error().Err(err).Msgf("Error subscribed to %s", e.EventName)
			time.Sleep(3 * time.Second)
			continue
		}
		log.Info().Msgf("Subscribing to %s", e.EventName)

		for {
			select {
			case <-ctx.Done():
				sub.Unsubscribe()
				log.Info().Msgf("Unsubscribed from %s", e.EventName)
				return
			case err := <-sub.Err():
				log.Error().Err(err).Msgf("Subscribed to %s", e.EventName)
				sub.Unsubscribe()
				time.Sleep(3 * time.Second)
				goto RESUBSCRIBE
			case vLog := <-logs:
				err = e.handleEvent(vLog)
				if err != nil {
					log.Error().Err(err).Msgf("Error handling ethereum event %s", e.EventName)
				}
				saveLastBlock(vLog.BlockNumber) // todo rethink the logic here
			}
		}
	RESUBSCRIBE:
	}

}

func (e *EventListener[T]) handleEvent(vLog types.Log) error {
	var eventData T

	err := e.ContractABI.UnpackIntoInterface(&eventData, e.EventName, vLog.Data)
	if err != nil {
		return err
	}

	e.handler(eventData)
	return nil
}
