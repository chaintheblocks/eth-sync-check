package collector

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/blocknative/sync-check/internal/blockscan"
	"github.com/blocknative/sync-check/internal/util"
)

var Client *ethclient.Client

// Init initializes the Ethereum client
func Init(ipcPath, httpURL string) error {
	var err error = nil
	if ipcPath != "" {
		Client, err = ethclient.Dial(ipcPath)
		if err != nil {
			log.Printf("error: Failed to connect to IPC, falling back to http: %s", err)
		}
	}

	if ipcPath == "" || err != nil {
		Client, err = ethclient.Dial(httpURL)
	}

	if err != nil {
		return err
	}
	return nil
}

// CollectExecutionMetrics returns metrics for the local execution node.
func CollectExecutionMetrics(apiKey string, chainID int) (map[string]uint64, error) {
	// create a map to store the metrics
	metrics := make(map[string]uint64)

	// get the latest block header
	header, err := Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block header: %v", err)
	}

	metrics["executionCurrentBlock"] = header.Number.Uint64()

	// get local HighestBlock from SyncProgress
	syncProgress, err := Client.SyncProgress(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to check sync progress: %v", err)
	}
	var executionLocalHighestBlock uint64
	if syncProgress == nil {
		executionLocalHighestBlock = metrics["executionCurrentBlock"]
	} else {
		executionLocalHighestBlock = syncProgress.HighestBlock
	}
	metrics["executionLocalHighestBlock"] = executionLocalHighestBlock

	// get network HighestBlock from etherscan
	executionNetworkHighestBlock, err := blockscan.EtherscanGetCurrentBlockNumber(apiKey, chainID)
	if err != nil {
		metrics["executionNetworkHighestBlock"] = 0
		err = nil // dont fail completely if etherscan fails
	} else {
		metrics["executionNetworkHighestBlock"] = executionNetworkHighestBlock
	}

	metrics["executionLocalDiff"] = metrics["executionLocalHighestBlock"] - metrics["executionCurrentBlock"]
	if metrics["executionNetworkHighestBlock"] > metrics["executionCurrentBlock"] {
		metrics["executionNetworkDiff"] = metrics["executionNetworkHighestBlock"] - metrics["executionCurrentBlock"]
	} else {
		metrics["executionNetworkDiff"] = 0
	}

	return metrics, err
}

// CollectConsensusMetrics returns metrics for the local consensus node.
func CollectConsensusMetrics(beaconEndpoint string) (map[string]interface{}, error) {
	// send request to get syncing data
	var result struct {
		Data struct {
			HeadSlot     string `json:"head_slot"`
			SyncDistance string `json:"sync_distance"`
			IsSyncing    bool   `json:"is_syncing"`
			IsOptimistic bool   `json:"is_optimistic"`
		} `json:"data"`
	}
	err := util.GetJSON(fmt.Sprintf("%s/eth/v1/node/syncing", beaconEndpoint), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get consensus syncing data: %v", err)
	}

	consensusCurrentSlot, err := strconv.ParseUint(result.Data.HeadSlot, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to get consensus syncing data: %v", err)
	}
	consensusSyncDistance, err := strconv.ParseUint(result.Data.SyncDistance, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to get consensus syncing data: %v", err)
	}
	consensusIsSyncing := result.Data.IsSyncing
	consensusIsOptimistic := result.Data.IsOptimistic

	consensusStatus, err := util.GetStatusCode(fmt.Sprintf("%s/eth/v1/node/health", beaconEndpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to get consensus status: %v", err)
	}

	metrics := map[string]interface{}{
		"consensusCurrentSlot":  consensusCurrentSlot,
		"consensusSyncDistance": consensusSyncDistance,
		"consensusIsSyncing":    consensusIsSyncing,
		"consensusIsOptimistic": consensusIsOptimistic,
		"consensusStatus":       consensusStatus,
	}

	return metrics, nil
}
