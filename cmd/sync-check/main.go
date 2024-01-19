package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/blocknative/sync-check/internal/collector"
	"github.com/blocknative/sync-check/internal/metrics"
	"github.com/blocknative/sync-check/internal/util"
)

var rootCmd = &cobra.Command{
	Use:   "sync-check",
	Short: "Sync-check is a tool for checking the synchronization of a blockchain node",
	Run: func(cmd *cobra.Command, args []string) {
		executionIPC := viper.GetString("execution-ipc")
		executionHTTP := viper.GetString("execution-http")
		consensusHTTP := viper.GetString("consensus-http")
		etherscanAPIKey := viper.GetString("etherscan-api-key")
		daemon := viper.GetBool("daemon")
		metricsPort := viper.GetString("metrics-port")

		err := collector.Init(executionIPC, executionHTTP)
		if err != nil {
			log.Fatalf("error: %s", err)
		}

		chainID, err := util.GetChainID(collector.Client)
		if err != nil {
			log.Printf("warn: failed to get chain id: %s", err)
		}

		var isPolygon bool

		if chainID == 137 || chainID == 80001 {
			isPolygon = true
		} else {
			isPolygon = false
		}

		if daemon {
			go func() {
				for {
					err := metrics.UpdatePrometheusMetrics(etherscanAPIKey, consensusHTTP, chainID, isPolygon)
					if err != nil {
						log.Printf("encountered error when collecting metrics: %s", err)
					}
					time.Sleep(time.Second * 2)
				}
			}()
			log.Println(("Starting prometheus server..."))
			prometheus.MustRegister(version.NewCollector("sync_check"))
			http.Handle("/metrics", promhttp.Handler())
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", metricsPort), nil))
		} else {
			metrics.LogMetrics(etherscanAPIKey, consensusHTTP, chainID, isPolygon)
		}
	},
}

func main() {
	// Define command-line flags
	rootCmd.PersistentFlags().String("execution-ipc", "", "Execution IPC endpoint")
	rootCmd.PersistentFlags().String("execution-http", "http://localhost:8545", "Execution HTTP endpoint")
	rootCmd.PersistentFlags().String("consensus-http", "http://localhost:5052", "Consensus HTTP endpoint")
	rootCmd.PersistentFlags().String("etherscan-api-key", "", "Etherscan API key")
	rootCmd.PersistentFlags().Bool("daemon", false, "Continuously collect and expose metrics via Prometheus")
	rootCmd.PersistentFlags().String("metrics-port", "3737", "Prometheus port")

	// Bind environment variables to flags
	viper.BindPFlag("execution-ipc", rootCmd.PersistentFlags().Lookup("execution-ipc"))
	viper.BindPFlag("execution-http", rootCmd.PersistentFlags().Lookup("execution-http"))
	viper.BindPFlag("consensus-http", rootCmd.PersistentFlags().Lookup("consensus-http"))
	viper.BindPFlag("etherscan-api-key", rootCmd.PersistentFlags().Lookup("etherscan-api-key"))
	viper.BindPFlag("daemon", rootCmd.PersistentFlags().Lookup("daemon"))
	viper.BindPFlag("metrics-port", rootCmd.PersistentFlags().Lookup("metrics-port"))

	// Bind environment variables to viper keys
	viper.BindEnv("execution-ipc", "SYNC_CHECK_EXECUTION_IPC")
	viper.BindEnv("execution-http", "SYNC_CHECK_EXECUTION_HTTP")
	viper.BindEnv("consensus-http", "SYNC_CHECK_CONSENSUS_HTTP")
	viper.BindEnv("etherscan-api-key", "SYNC_CHECK_ETHERSCAN_API_KEY")
	viper.BindEnv("daemon", "SYNC_CHECK_DAEMON")
	viper.BindEnv("metrics-port", "SYNC_CHECK_METRICS_PORT")

	// Execute root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
