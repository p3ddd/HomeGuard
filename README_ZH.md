# HomeGuard

[![Go Version](https://img.shields.io/github/go-mod/go-version/p3ddd/HomeGuard)](https://golang.org/)
[![License](https://img.shields.io/github/license/p3ddd/HomeGuard)](LICENSE)

简单易用的 Wake-on-LAN 远程唤醒服务，支持 HTTP 和 MQTT 协议。

[English](README.md) | [中文](README_ZH.md)

---

## 特性

- 🚀 多协议支持（HTTP 和 MQTT）
- 📝 设备管理（YAML 配置文件）
- 🔄 支持设备名称或 MAC 地址唤醒
- 🌐 支持云端 MQTT（如巴法云）
- 🛡️ 优雅关闭
- ⚡ 轻量级单二进制文件

## 使用场景

HomeGuard 是一个**单体服务**，可以**同时支持 HTTP 和 MQTT**。

**场景一：通过云端 MQTT 远程唤醒**
- HomeGuard 在家中运行，连接到云端 MQTT 服务（如巴法云）
- 在外网任何地方通过 MQTT 客户端应用发送唤醒指令
- 无需暴露 HTTP 端口或配置 VPN

**场景二：通过 VPN 连接局域网**
- 通过 VPN 连接到家庭局域网
- 直接使用 HTTP API 唤醒设备
- 更直接，延迟更低

**两种方式可以同时使用！** HomeGuard 可以在监听 HTTP 的同时连接云端 MQTT 服务。

## 快速开始

### 安装

```bash
# 克隆
git clone https://github.com/p3ddd/HomeGuard.git
cd HomeGuard

# 编译
task build

# 或使用 go
go build -o homeguard
```

### 配置

```bash
# 创建配置
cp config.example.yaml config.yaml
vim config.yaml
```

配置示例 `config.yaml`：

```yaml
devices:
  - name: desktop
    mac: "00:11:22:33:44:55"
    broadcast: "192.168.1.255"

server:
  http:
    addr: ":7092"
  mqtt:
    enabled: false  # 设置为 true 启用 MQTT
    broker: "tcp://mqtt.bemfa.com:9501"
    topic: "homeguard001"
```

### 运行

```bash
# 使用配置文件启动（推荐）
./homeguard -config config.yaml

# 或使用命令行参数
./homeguard -http :7092 -mqtt-broker tcp://mqtt.bemfa.com:9501 -mqtt-topic your-topic
```

在 `config.yaml` 中启用/禁用功能：
```yaml
server:
  http:
    enabled: true  # HTTP API
  mqtt:
    enabled: true  # 云端 MQTT
```

## 使用方法

### HTTP API

**方式一：通过设备名唤醒**

```bash
# GET 请求
curl "http://localhost:7092/wakeup?device=desktop"

# POST JSON
curl -X POST http://localhost:7092/wakeup \
  -H "Content-Type: application/json" \
  -d '{"device":"desktop"}'
```

**方式二：通过 MAC 地址唤醒**

```bash
# GET 请求
curl "http://localhost:7092/wakeup?mac=00:11:22:33:44:55&broadcast=192.168.1.255"

# POST JSON
curl -X POST http://localhost:7092/wakeup \
  -H "Content-Type: application/json" \
  -d '{"mac":"00:11:22:33:44:55","broadcast":"192.168.1.255"}'
```

### MQTT（云服务）

将 HomeGuard 连接到云端 MQTT 服务（如巴法云），然后在任何地方发布消息：

**方式一：通过设备名唤醒**

```bash
# 使用 mosquitto_pub
mosquitto_pub -h mqtt.bemfa.com -p 9501 -t your-topic \
  -m '{"device":"desktop"}'

# 或使用手机/电脑上的任何 MQTT 客户端应用
```

**方式二：通过 MAC 地址唤醒**

```bash
mosquitto_pub -h mqtt.bemfa.com -p 9501 -t your-topic \
  -m '{"mac":"00:11:22:33:44:55","broadcast":"192.168.1.255"}'
```

### 客户端工具

```bash
# 编译客户端
go build -o wolctl ./cmd/wolctl/

# 通过设备名唤醒
./wolctl -device desktop

# 通过 MAC 地址唤醒
./wolctl -mac 00:11:22:33:44:55 -broadcast 192.168.1.255

# 指定服务器地址
./wolctl -server http://192.168.1.100:7092 -device desktop
```

## Task 命令

```bash
task build          # 编译
task run            # 运行
task test           # 测试
task clean          # 清理
task fmt            # 格式化
task lint           # 代码检查
```

## 命令行参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-config` | `config.yaml` | 配置文件路径 |
| `-http` | `:7092` | HTTP 监听地址 |
| `-mqtt-broker` | ` ` | MQTT Broker 地址（如：tcp://mqtt.bemfa.com:9501） |
| `-mqtt-topic` | `homeguard/wakeup` | MQTT 主题 |
| `-log-level` | `info` | 日志级别（debug/info/warn/error） |

## Docker

```bash
# 构建
docker build -t homeguard:latest .

# 运行（使用配置文件）
docker run --rm --network host -v $(pwd)/config.yaml:/app/config.yaml:ro homeguard:latest

# Docker Compose（推荐）
docker-compose up -d
```

运行前在 `config.yaml` 中配置需要的功能。

## 贡献

参见 [CONTRIBUTING.md](CONTRIBUTING.md)

## 许可证

MIT 许可证 - 详见 [LICENSE](LICENSE)

## 相关项目

- [HomeGuard-rs](https://github.com/p3ddd/HomeGuard-rs) - Rust 实现版本

---

<div align="center">
Made with ❤️
</div>

