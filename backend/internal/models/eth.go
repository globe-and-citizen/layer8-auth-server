package models

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type TrafficPaidEvent struct {
	ClientID string
	Payer    common.Address
	Amount   *big.Int
}
