package eth

import (
	"bytes"
	"encoding/json"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const blockFile = "./ethereum_blocks.txt"

func MustLoadABI(path string) abi.ABI {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("abi read error:", err)
	}

	var artifact struct {
		ABI json.RawMessage `json:"abi"`
	}

	if err := json.Unmarshal(data, &artifact); err != nil {
		log.Fatal("artifact unmarshal error:", err)
	}

	parsed, err := abi.JSON(bytes.NewReader(artifact.ABI))
	if err != nil {
		log.Fatal("abi parse error:", err)
	}

	return parsed
}

func LoadLastBlock() uint64 {
	data, err := os.ReadFile(blockFile)
	if err != nil {
		return 0
	}
	n, _ := new(big.Int).SetString(string(data), 10)
	return n.Uint64()
}

func saveLastBlock(block uint64) {
	_ = os.WriteFile(blockFile, []byte(big.NewInt(int64(block)).String()), 0644)
}
