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

func DBStringToBigInt(dbValue string) (*big.Int, error) {
	if dbValue == "" || dbValue == "0" {
		return big.NewInt(0), nil
	}

	amount := new(big.Int)
	if _, ok := amount.SetString(dbValue, 10); !ok {
		return nil, fmt.Errorf("cannot convert %s to big int", dbValue)
	}

	return amount, nil
}

// BaseUnitsToCurrency Convert base units amount of a currency to its string representation in that currency.
func BaseUnitsToCurrency(baseUnits *big.Int, currency string) string {
	if baseUnits == nil {
		return "0"
	}

	decimals := 18
	switch strings.ToUpper(strings.TrimSpace(currency)) {
	case "WEI":
		decimals = 0
	case "GWEI":
		decimals = 9
	case "USDC", "USDT":
		decimals = 6
	case "ETH", "WETH", "MATIC", "POL", "BNB":
		decimals = 18
	}

	v := new(big.Int).Set(baseUnits)
	negative := v.Sign() < 0
	if negative {
		v.Abs(v)
	}

	s := v.String()
	if decimals == 0 {
		if negative {
			return "-" + s
		}
		return s
	}

	if len(s) <= decimals {
		s = strings.Repeat("0", decimals-len(s)+1) + s
	}

	split := len(s) - decimals
	intPart := s[:split]
	fracPart := strings.TrimRight(s[split:], "0")

	result := intPart
	if fracPart != "" {
		result += "." + fracPart
	}

	if negative {
		return "-" + result
	}

	return result
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

func BigIntToDBString(amount *big.Int) string {
	if amount == nil {
		return "0"
	}
	return amount.String()
}
