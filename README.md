# HomeGuard

[![Go Version](https://img.shields.io/github/go-mod/go-version/p3ddd/HomeGuard)](https://golang.org/)
[![License](https://img.shields.io/github/license/p3ddd/HomeGuard)](LICENSE)

A simple Wake-on-LAN service with HTTP and MQTT support.

[English](README.md) | [中文](README_ZH.md)

---

## Features

- 🚀 Multi-protocol support (HTTP and MQTT)
- 📝 Device management via YAML configuration
- 🔄 Wake by device name or MAC address
- 🌐 Cloud MQTT support (e.g., Bemfa Cloud)
- 🛡️ Graceful shutdown
- ⚡ Lightweight single binary

## Use Cases

HomeGuard is a **single service** that supports **both HTTP and MQTT** simultaneously.

**Scenario 1: Remote wake via Cloud MQTT**
- HomeGuard runs at home and connects to cloud MQTT service (e.g., Bemfa Cloud)
- Send wake commands from anywhere via MQTT client apps
- No need to expose HTTP port or set up VPN

**Scenario 2: Local network via VPN**
- Connect to home network via VPN
- Use HTTP API directly to wake devices
- More direct and lower latency

**You can use both methods at the same time!** HomeGuard can listen on HTTP while connected to cloud MQTT service.

## Quick Start

### Installation

```bash
# Clone
git clone https://github.com/p3ddd/HomeGuard.git
cd HomeGuard

# Build
task build

# Or with go
go build -o homeguard
```

### Configuration

```bash
# Create config
cp config.example.yaml config.yaml
vim config.yaml
```

Example `config.yaml`:

```yaml
devices:
  - name: desktop
    mac: "00:11:22:33:44:55"
    broadcast: "192.168.1.255"

server:
  http:
    addr: ":7092"
  mqtt:
    enabled: false  # Set to true to enable MQTT
    broker: "tcp://mqtt.bemfa.com:9501"
    topic: "homeguard001"
```

### Run

```bash
# Start with config file (recommended)
./homeguard -config config.yaml

# Or use command line flags
./homeguard -http :7092 -mqtt-broker tcp://mqtt.bemfa.com:9501 -mqtt-topic your-topic
```

Enable/disable features in `config.yaml`:
```yaml
server:
  http:
    enabled: true  # HTTP API
  mqtt:
    enabled: true  # Cloud MQTT
```

## Usage

### HTTP API

**Method 1: Wake by device name**

```bash
# GET request
curl "http://localhost:7092/wakeup?device=desktop"

# POST JSON
curl -X POST http://localhost:7092/wakeup \
  -H "Content-Type: application/json" \
  -d '{"device":"desktop"}'
```

**Method 2: Wake by MAC address**

```bash
# GET request
curl "http://localhost:7092/wakeup?mac=00:11:22:33:44:55&broadcast=192.168.1.255"

# POST JSON
curl -X POST http://localhost:7092/wakeup \
  -H "Content-Type: application/json" \
  -d '{"mac":"00:11:22:33:44:55","broadcast":"192.168.1.255"}'
```

### MQTT (Cloud Service)

Connect HomeGuard to cloud MQTT service (e.g., Bemfa Cloud), then publish messages from anywhere:

**Method 1: Wake by device name**

```bash
# Using mosquitto_pub
mosquitto_pub -h mqtt.bemfa.com -p 9501 -t your-topic \
  -m '{"device":"desktop"}'

# Or use any MQTT client app on your phone/computer
```

**Method 2: Wake by MAC address**

```bash
mosquitto_pub -h mqtt.bemfa.com -p 9501 -t your-topic \
  -m '{"mac":"00:11:22:33:44:55","broadcast":"192.168.1.255"}'
```

### Client Tool

```bash
# Build client
go build -o wolctl ./cmd/wolctl/

# Wake by device name
./wolctl -device desktop

# Wake by MAC address
./wolctl -mac 00:11:22:33:44:55 -broadcast 192.168.1.255

# Specify server
./wolctl -server http://192.168.1.100:7092 -device desktop
```

## Task Commands

```bash
task build          # Build binaries
task run            # Run server
task test           # Run tests
task clean          # Clean build files
task fmt            # Format code
task lint           # Run linter
```

## Command Line Options

| Flag | Default | Description |
|------|---------|-------------|
| `-config` | `config.yaml` | Config file path |
| `-http` | `:7092` | HTTP address |
| `-mqtt-broker` | ` ` | MQTT broker (e.g., tcp://mqtt.bemfa.com:9501) |
| `-mqtt-topic` | `homeguard/wakeup` | MQTT topic |
| `-log-level` | `info` | Log level (debug/info/warn/error) |

## Docker

```bash
# Build
docker build -t homeguard:latest .

# Run with config file
docker run --rm --network host -v $(pwd)/config.yaml:/app/config.yaml:ro homeguard:latest

# Docker Compose (recommended)
docker-compose up -d
```

Configure features in `config.yaml` before running.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

## License

MIT License - see [LICENSE](LICENSE)

## Related

- [HomeGuard-rs](https://github.com/p3ddd/HomeGuard-rs) - Rust implementation

---

<div align="center">
Made with ❤️
</div>
