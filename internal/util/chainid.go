package util

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
)

// GetChainID returns the chain ID of the connected Ethereum client.
func GetChainID(client *ethclient.Client) (int, error) {
	ctx := context.Background()
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return 0, err
	}
	return int(chainID.Int64()), nil
}
