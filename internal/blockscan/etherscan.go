package blockscan

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type response struct {
	Result string `json:"result"`
}

func EtherscanGetCurrentBlockNumber(apiKey string, chainID int) (uint64, error) {
	var baseURL string
	switch chainID {
	case 1:
		baseURL = "https://api.etherscan.io"
	case 5:
		baseURL = "https://api-goerli.etherscan.io"
	case 11155111:
		baseURL = "https://api-sepolia.etherscan.io"
	default:
		return 0, fmt.Errorf("unsupported chainID: %d", chainID)
	}

	url := fmt.Sprintf("%s/api?module=proxy&action=eth_blockNumber&apikey=%s", baseURL, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var r response
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return 0, err
	}

	// strip 0x prefix from result and convert hex to dec
	blockNumber, err := strconv.ParseUint(r.Result[2:], 16, 64)
	if err != nil {
		return 0, err
	}

	return blockNumber, nil
}
