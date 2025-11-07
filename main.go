package main

import (
	"HomeGuard/device"
	"HomeGuard/listener"
	"HomeGuard/wol"
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	configPath   = flag.String("config", "config.yaml", "Path to configuration file")
	httpAddr     = flag.String("http", ":7092", "HTTP listener address")
	mqttBroker   = flag.String("mqtt-broker", "", "MQTT broker URL (e.g., tcp://localhost:1883)")
	mqttTopic    = flag.String("mqtt-topic", "homeguard/wakeup", "MQTT topic to subscribe to")
	mqttClientID = flag.String("mqtt-client-id", "", "MQTT client ID (default: auto-generated)")
	mqttUsername = flag.String("mqtt-username", "", "MQTT username")
	mqttPassword = flag.String("mqtt-password", "", "MQTT password")
	mqttQoS      = flag.Uint("mqtt-qos", 1, "MQTT QoS level (0, 1, or 2)")
	logLevel     = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
)

func main() {
	flag.Parse()

	// Setup logging
	var level slog.Level
	switch *logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting HomeGuard WOL Service")

	// Load device configuration
	deviceManager, err := device.NewManager(*configPath)
	if err != nil {
		slog.Error("Failed to load device configuration", "error", err, "path", *configPath)
		slog.Warn("Continuing without device configuration - only direct MAC/broadcast requests will work")
		deviceManager = nil
	} else {
		devices := deviceManager.ListDevices()
		slog.Info("Loaded device configuration", "count", len(devices))
		for _, dev := range devices {
			slog.Debug("Registered device",
				"name", dev.Name,
				"mac", dev.Mac,
				"broadcast", dev.Broadcast,
				"description", dev.Description)
		}
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create channel for wakeup requests
	requestChan := make(chan listener.WakeUpRequest, 100)

	// Start request processor
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		processRequests(ctx, requestChan, deviceManager)
	}()

	// Start listeners
	listeners := make([]listener.Listener, 0)

	// HTTP Listener
	httpListener := listener.NewHTTPListener(*httpAddr)
	listeners = append(listeners, httpListener)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpListener.Start(ctx, requestChan); err != nil {
			slog.Error("HTTP listener error", "error", err)
		}
	}()

	// MQTT Listener (if configured)
	if *mqttBroker != "" {
		mqttConfig := listener.MQTTConfig{
			Broker:   *mqttBroker,
			ClientID: *mqttClientID,
			Topic:    *mqttTopic,
			QoS:      byte(*mqttQoS),
			Username: *mqttUsername,
			Password: *mqttPassword,
		}
		mqttListener := listener.NewMQTTListener(mqttConfig)
		listeners = append(listeners, mqttListener)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := mqttListener.Start(ctx, requestChan); err != nil {
				slog.Error("MQTT listener error", "error", err)
			}
		}()
	} else {
		slog.Info("MQTT listener not configured (use -mqtt-broker flag to enable)")
	}

	slog.Info("HomeGuard WOL Service is running. Press Ctrl+C to stop.")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	slog.Info("Shutdown signal received, stopping services...")

	// Cancel context to stop all listeners
	cancel()

	// Stop all listeners explicitly
	for _, l := range listeners {
		if err := l.Stop(); err != nil {
			slog.Error("Failed to stop listener", "name", l.Name(), "error", err)
		}
	}

	// Close request channel
	close(requestChan)

	// Wait for all goroutines to finish with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		slog.Info("All services stopped gracefully")
	case <-time.After(5 * time.Second):
		slog.Warn("Shutdown timeout exceeded, forcing exit")
	}
}

func processRequests(ctx context.Context, requestChan <-chan listener.WakeUpRequest, deviceManager *device.Manager) {
	slog.Info("Request processor started")
	for {
		select {
		case req, ok := <-requestChan:
			if !ok {
				slog.Info("Request channel closed, stopping processor")
				return
			}
			handleWakeUpRequest(req, deviceManager)
		case <-ctx.Done():
			slog.Info("Request processor context canceled")
			return
		}
	}
}

func handleWakeUpRequest(req listener.WakeUpRequest, deviceManager *device.Manager) {
	var mac, broadcast string

	// If device name is provided, look it up
	if req.DeviceName != "" {
		if deviceManager == nil {
			slog.Error("Device name provided but no device configuration loaded",
				"device", req.DeviceName,
				"type", req.Type)
			return
		}

		dev, err := deviceManager.GetDevice(req.DeviceName)
		if err != nil {
			slog.Error("Failed to get device information",
				"device", req.DeviceName,
				"error", err,
				"type", req.Type)
			return
		}

		mac = dev.Mac
		broadcast = dev.Broadcast
		slog.Info("Resolved device name to MAC address",
			"device", req.DeviceName,
			"mac", mac,
			"broadcast", broadcast,
			"type", req.Type)
	} else {
		// Use provided MAC and broadcast
		mac = req.Mac
		broadcast = req.Broadcast
		slog.Info("Using direct MAC address",
			"mac", mac,
			"broadcast", broadcast,
			"type", req.Type)
	}

	// Send WOL magic packet
	if err := wol.WakeOnLan(mac, broadcast); err != nil {
		slog.Error("Failed to send WOL packet",
			"mac", mac,
			"broadcast", broadcast,
			"error", err,
			"type", req.Type)
		return
	}

	slog.Info("Successfully sent WOL packet",
		"mac", mac,
		"broadcast", broadcast,
		"device", req.DeviceName,
		"type", req.Type)
}
