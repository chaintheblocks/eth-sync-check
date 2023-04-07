package metrics

import (
	"fmt"
	"os"
	"strings"

	"github.com/blocknative/sync-check/internal/collector"
	"github.com/olekukonko/tablewriter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// create gauges for each variable
var (
	executionCurrentBlockGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "sync_execution_current_block",
		Help: "Current block number in the execution node",
	})
	executionLocalHighestBlockGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "sync_execution_local_highest_block",
		Help: "Local highest block number in the execution node",
	})
	executionNetworkHighestBlockGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "sync_execution_network_highest_block",
		Help: "Network highest block number in the execution node",
	})
	executionLocalDiffGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "sync_execution_local_diff",
		Help: "Difference between current block vs node's highest known block",
	})
	executionNetworkDiffGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "sync_execution_network_diff",
		Help: "Difference between current block and etherscan block",
	})
	consensusCurrentSlotGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "sync_consensus_current_slot",
		Help: "Current slot number in the consensus node",
	})
	consensusSyncDistanceGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "sync_consensus_sync_distance",
		Help: "Sync distance in the consensus node",
	})
)

func collectMetrics(apiKey string, chainID int, consensusHTTP string) (map[string]uint64, map[string]interface{}, error) {
	executionMetrics, err := collector.CollectExecutionMetrics(apiKey, chainID)
	if err != nil {
		if strings.Contains(err.Error(), "unsupported chainID") {
		} else {
			return nil, nil, err
		}
	}

	consensusMetrics, err := collector.CollectConsensusMetrics(consensusHTTP)
	if err != nil {
		return nil, nil, err
	}

	return executionMetrics, consensusMetrics, nil
}

func LogMetrics(etherscanAPIKey, consensusHTTP string, chainID int) error {
	executionMetrics, consensusMetrics, err := collectMetrics(etherscanAPIKey, chainID, consensusHTTP)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Current Block", "Local Highest", "Network Highest", "Local Diff", "Network Diff", "CL Slot", "CL Slot Distance", "CL Status", "CL Syncing", "CL Optimistic"})
	table.Append([]string{
		fmt.Sprintf("%d", executionMetrics["executionCurrentBlock"]),
		fmt.Sprintf("%d", executionMetrics["executionLocalHighestBlock"]),
		fmt.Sprintf("%d", executionMetrics["executionNetworkHighestBlock"]),
		fmt.Sprintf("%d", executionMetrics["executionLocalDiff"]),
		fmt.Sprintf("%d", executionMetrics["executionNetworkDiff"]),
		fmt.Sprintf("%d", consensusMetrics["consensusCurrentSlot"]),
		fmt.Sprintf("%d", consensusMetrics["consensusSyncDistance"]),
		fmt.Sprintf("%d", consensusMetrics["consensusStatus"]),
		fmt.Sprintf("%v", consensusMetrics["consensusIsSyncing"]),
		fmt.Sprintf("%v", consensusMetrics["consensusIsOptimistic"]),
	})
	table.Render()

	if err != nil {
		return err
	}

	return nil
}

// Collect metrics and update in prometheus
func UpdatePrometheusMetrics(etherscanAPIKey, consensusHTTP string, chainID int) error {
	executionMetrics, consensusMetrics, err := collectMetrics(etherscanAPIKey, chainID, consensusHTTP)
	if err != nil {
		return err
	}

	executionCurrentBlockGauge.Set(float64(executionMetrics["executionCurrentBlock"]))
	executionLocalHighestBlockGauge.Set(float64(executionMetrics["executionLocalHighestBlock"]))
	executionNetworkHighestBlockGauge.Set(float64(executionMetrics["executionNetworkHighestBlock"]))
	executionLocalDiffGauge.Set(float64(executionMetrics["executionLocalDiff"]))
	executionNetworkDiffGauge.Set(float64(executionMetrics["executionNetworkDiff"]))
	consensusCurrentSlotGauge.Set(float64(consensusMetrics["consensusCurrentSlot"].(uint64)))
	consensusSyncDistanceGauge.Set(float64(consensusMetrics["consensusSyncDistance"].(uint64)))

	return nil
}
