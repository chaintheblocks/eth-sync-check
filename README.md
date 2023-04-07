# Ethereum Node Sync Checker
This is a command-line tool that allows you to check the synchronization status of your Ethereum node. You can also run it in daemon mode to expose the sync information in Prometheus format.

## Requirements
- Go version 1.19 or higher
- An Ethereum node that is running and accessible from your machine
- (Optional) An Etherscan API key to get the network highest block

## Installation
- Download the [latest release](https://github.com/blocknative/ethereum-sync-check/releases/latest) for your OS/Architecture.
- Unarchive the release e.g. `tar -xvzf sync-check-v0.1.0-linux-amd64.tar.gz`
- (Optional) Move the binary to a location that makes sense e.g. `mv sync-check /usr/local/bin/`

## Usage
```
Usage:
  sync-check [flags]

Flags:
      --consensus-http string      Consensus HTTP endpoint (default "http://localhost:5052")
      --daemon                     Continuously collect and expose metrics via Prometheus
      --etherscan-api-key string   Etherscan API key
      --execution-http string      Execution HTTP endpoint (default "http://localhost:8545")
      --execution-ipc string       Execution IPC endpoint
  -h, --help                       help for sync-check
      --metrics-port string        Prometheus port (default "3737")
```

The value for each flag may also be set via environment variables:
```
SYNC_CHECK_EXECUTION_IPC
SYNC_CHECK_EXECUTION_HTTP
SYNC_CHECK_CONSENSUS_HTTP
SYNC_CHECK_ETHERSCAN_API_KEY
SYNC_CHECK_DAEMON
SYNC_CHECK_METRICS_PORT
```

### Example output
```
+---------------+---------------+-----------------+------------+--------------+---------+------------------+-----------+------------+---------------+
| CURRENT BLOCK | LOCAL HIGHEST | NETWORK HIGHEST | LOCAL DIFF | NETWORK DIFF | CL SLOT | CL SLOT DISTANCE | CL STATUS | CL SYNCING | CL OPTIMISTIC |
+---------------+---------------+-----------------+------------+--------------+---------+------------------+-----------+------------+---------------+
|      16957079 |      16957079 |        16957079 |          0 |            0 | 6130154 |                0 |       200 | false      | false         |
+---------------+---------------+-----------------+------------+--------------+---------+------------------+-----------+------------+---------------+
```

### Daemon mode
When running in daemon mode, this will start a web server that exposes the sync information in Prometheus format. You can then scrape this information with Prometheus and use it for monitoring and alerting.
