package utils

import (
	"fmt"
	"math/big"
	"strings"
)

var weiPerEth = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

func WeiToEth(wei *big.Int) *big.Float {
	f := new(big.Float)
	f.SetInt(wei)

	eth := new(big.Float).Quo(f, big.NewFloat(1e18))
	return eth
}

func DBWeiToBigInt(dbValue string) (*big.Int, error) {
	if dbValue == "" {
		return big.NewInt(0), nil
	}

	wei := new(big.Int)
	if _, ok := wei.SetString(dbValue, 10); !ok {
		return nil, fmt.Errorf("invalid wei value in DB")
	}

	return wei, nil
}

func WeiToEthString(wei *big.Int, decimals int) string {
	if wei == nil {
		return "0"
	}

	f := new(big.Float).SetInt(wei)
	f.Quo(f, new(big.Float).SetInt(weiPerEth))

	return f.Text('f', decimals)
}

func EthStringToWei(eth string) (*big.Int, error) {
	eth = strings.TrimSpace(eth)
	if eth == "" {
		return nil, fmt.Errorf("empty ETH value")
	}

	parts := strings.Split(eth, ".")
	if len(parts) == 2 && len(parts[1]) > 18 {
		return nil, fmt.Errorf("ETH supports max 18 decimals")
	}

	f, ok := new(big.Float).SetString(eth)
	if !ok {
		return nil, fmt.Errorf("invalid ETH value")
	}

	f.Mul(f, new(big.Float).SetInt(weiPerEth))

	wei := new(big.Int)
	f.Int(wei) // truncate (EVM behavior)

	return wei, nil
}

func BigIntToDBWei(wei *big.Int) string {
	if wei == nil {
		return "0"
	}
	return wei.String()
}
