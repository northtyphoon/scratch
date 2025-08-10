# Redis Connection Maintainer

A tiny Go tool to open and maintain 1,000–10,000 concurrent TCP/TLS connections to a Redis (cluster) endpoint with periodic PING keepalives.

- No external libraries; raw RESP for minimal overhead.
- Round-robin across provided node addresses.
- Reconnect with jittered backoff.
- Simple metrics in logs.

## Usage

Build and run locally:

```pwsh
# Build
$env:CGO_ENABLED=0; go build -o redis-conn .

# Run 5000 conns across 3 cluster nodes
./redis-conn --addrs "10.0.0.1:6379,10.0.0.2:6379,10.0.0.3:6379" --connections 5000 --open-rate 500
```

Key flags (env var overrides in parentheses):
- `--addrs` (REDIS_ADDRESSES): comma-separated host:port list
- `--connections` (CONNECTIONS): number of concurrent connections
- `--username`/`--password` (REDIS_USERNAME/REDIS_PASSWORD)
- `--db` (REDIS_DB)
- `--tls` (REDIS_TLS=true|false)
- `--ping-interval` (PING_INTERVAL, e.g. 30s)
- `--read-timeout` (READ_TIMEOUT, e.g. 2s)
- `--idle-timeout` (IDLE_TIMEOUT, e.g. 10m)
- `--connect-timeout` (CONNECT_TIMEOUT, e.g. 5s)
- `--open-rate` (OPEN_RATE): new connections per second

Note: Creating 10k sockets may require tuning OS limits (ulimits, sysctl) on the host.

## Container

Build the image:

```pwsh
# Docker build (Linux/amd64)
docker build -t redis-conn:latest .
```

Run the container with 10k connections:

```pwsh
# Example to 3 nodes, 10k connections, TLS disabled
docker run --rm --ulimit nofile=200000:200000 \
  --network host \
  -e CONNECTIONS=10000 \
  -e REDIS_ADDRESSES="10.0.0.1:6379,10.0.0.2:6379,10.0.0.3:6379" \
  -e OPEN_RATE=1000 \
  redis-conn:latest
```

If not using host network, map any needed ports for your environment, and ensure container can reach Redis IPs.

## Notes
- This tool does not perform cluster redirections; it only keeps connections alive with PING.
- For TLS testing with self-signed certs, the client skips verification by default (test only). Harden for production by providing proper certs.
- Adjust `--open-rate` to avoid overwhelming the network during ramp-up.
