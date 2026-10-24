package eth

import (
	"bytes"
	"encoding/json"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const TrafficPaidEventBlock = "./traffic_paid_last_block.txt"
const CNCTokenDepositedEventBlock = "./cnc_token_deposited_last_block.txt"

func MustLoadABI(path string) abi.ABI {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("abi read error:", err)
	}

	// Support Hardhat artifacts containing an "abi" field.
	var artifact struct {
		ABI json.RawMessage `json:"abi"`
	}

	if err := json.Unmarshal(data, &artifact); err == nil &&
		len(artifact.ABI) > 0 {
		data = artifact.ABI
	}

	parsed, err := abi.JSON(bytes.NewReader(data))
	if err != nil {
		log.Fatal("abi parse error:", err)
	}

	return parsed
}

func LoadLastBlock(eventName string) uint64 {
	data, err := os.ReadFile(getBlockFile(eventName))
	if err != nil {
		return 0
	}
	n, _ := new(big.Int).SetString(string(data), 10)
	return n.Uint64()
}

func saveLastBlock(eventName string, block uint64) {
	blockFile := getBlockFile(eventName)
	_ = os.WriteFile(blockFile, []byte(big.NewInt(int64(block)).String()), 0644)
}

func getBlockFile(eventName string) string {
	switch eventName {
	case "TrafficPaid":
		return TrafficPaidEventBlock
	case "TokenDeposited":
		return CNCTokenDepositedEventBlock
	default:
		return eventName + "_last_block.txt"
	}
}

type EventHandlerFunc[T any] func(eventData EventData[T]) error

type EventData[T any] struct {
	Data        T
	TxID        string
	LogIndex    uint
	TxTimestamp time.Time
	TxInput     []byte
}
