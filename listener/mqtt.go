package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var _ Listener = (*MQTTListener)(nil)

// MQTTConfig holds the configuration for MQTT listener.
type MQTTConfig struct {
	Broker   string
	ClientID string
	Topic    string
	QoS      byte
	Username string
	Password string
}

type MQTTListener struct {
	config MQTTConfig
	client mqtt.Client
	mu     sync.Mutex
}

// MQTTPayload represents the MQTT message payload structure.
type MQTTPayload struct {
	Device    string `json:"device,omitempty"`
	Mac       string `json:"mac,omitempty"`
	Broadcast string `json:"broadcast,omitempty"`
}

func NewMQTTListener(config MQTTConfig) *MQTTListener {
	if config.ClientID == "" {
		config.ClientID = fmt.Sprintf("wol-mqtt-%d", time.Now().Unix())
	}
	if config.QoS > 2 {
		config.QoS = 1 // Default to QoS 1
	}
	return &MQTTListener{
		config: config,
	}
}

func (l *MQTTListener) Name() string {
	return "MQTT"
}

func (l *MQTTListener) logger() *slog.Logger {
	return slog.With("type", l.Name())
}

func (l *MQTTListener) Start(ctx context.Context, wakeUpChan chan<- WakeUpRequest) error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(l.config.Broker)
	opts.SetClientID(l.config.ClientID)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetMaxReconnectInterval(1 * time.Minute)

	if l.config.Username != "" {
		opts.SetUsername(l.config.Username)
	}
	if l.config.Password != "" {
		opts.SetPassword(l.config.Password)
	}

	// Connection lost handler
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		l.logger().Error("MQTT connection lost", "error", err)
	})

	// On connect handler
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		l.logger().Info("Connected to MQTT broker", "broker", l.config.Broker)

		// Subscribe to topic
		token := client.Subscribe(l.config.Topic, l.config.QoS, func(client mqtt.Client, msg mqtt.Message) {
			l.handleMessage(ctx, msg, wakeUpChan)
		})

		if token.Wait() && token.Error() != nil {
			l.logger().Error("Failed to subscribe to topic", "topic", l.config.Topic, "error", token.Error())
		} else {
			l.logger().Info("Subscribed to topic", "topic", l.config.Topic, "qos", l.config.QoS)
		}
	})

	// Create client
	l.mu.Lock()
	l.client = mqtt.NewClient(opts)
	l.mu.Unlock()

	// Connect to broker
	l.logger().Info("Connecting to MQTT broker", "broker", l.config.Broker)
	token := l.client.Connect()
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}

	// Wait for context cancellation
	<-ctx.Done()
	l.logger().Info("MQTT listener context canceled")
	return l.Stop()
}

func (l *MQTTListener) handleMessage(ctx context.Context, msg mqtt.Message, wakeUpChan chan<- WakeUpRequest) {
	l.logger().Debug("Received MQTT message", "topic", msg.Topic(), "payload", string(msg.Payload()))

	var payload MQTTPayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		l.logger().Error("Failed to parse MQTT payload", "error", err, "payload", string(msg.Payload()))
		return
	}

	request := WakeUpRequest{
		Type:       l.Name(),
		DeviceName: payload.Device,
		Mac:        payload.Mac,
		Broadcast:  payload.Broadcast,
	}

	// Validate request: must have either device name or both mac and broadcast
	if request.DeviceName == "" && (request.Mac == "" || request.Broadcast == "") {
		l.logger().Error("Invalid MQTT message: must provide either device name or both mac and broadcast",
			"payload", string(msg.Payload()))
		return
	}

	// Send request to channel
	select {
	case wakeUpChan <- request:
		l.logger().Info("Processed MQTT wakeup request",
			"device", request.DeviceName,
			"mac", request.Mac,
			"broadcast", request.Broadcast)
	case <-ctx.Done():
		l.logger().Info("Context canceled while processing message")
	default:
		l.logger().Warn("Channel full, dropping message")
	}
}

// Stop implements Listener.
func (l *MQTTListener) Stop() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.client != nil && l.client.IsConnected() {
		// Unsubscribe from topic
		token := l.client.Unsubscribe(l.config.Topic)
		if token.Wait() && token.Error() != nil {
			l.logger().Error("Failed to unsubscribe from topic", "error", token.Error())
		}

		// Disconnect
		l.client.Disconnect(250)
		l.logger().Info("MQTT listener stopped")
	}

	return nil
}
