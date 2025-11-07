package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	version = "1.0.0"
)

var (
	serverURL = flag.String("server", "http://localhost:7092", "HomeGuard server URL")
	device    = flag.String("device", "", "Device name to wake up")
	mac       = flag.String("mac", "", "MAC address to wake up")
	broadcast = flag.String("broadcast", "", "Broadcast address")
	showVer   = flag.Bool("version", false, "Show version information")
)

type WakeUpRequest struct {
	Device    string `json:"device,omitempty"`
	Mac       string `json:"mac,omitempty"`
	Broadcast string `json:"broadcast,omitempty"`
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "wolctl - HomeGuard Wake-on-LAN Client Tool\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  wolctl -device <name>                    Wake up device by name\n")
		fmt.Fprintf(os.Stderr, "  wolctl -mac <MAC> -broadcast <addr>      Wake up device by MAC address\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  wolctl -device desktop\n")
		fmt.Fprintf(os.Stderr, "  wolctl -mac 00:11:22:33:44:55 -broadcast 192.168.1.255\n")
		fmt.Fprintf(os.Stderr, "  wolctl -server http://192.168.1.100:7092 -device laptop\n")
	}

	flag.Parse()

	if *showVer {
		fmt.Printf("wolctl version %s\n", version)
		os.Exit(0)
	}

	// Validate input
	if *device == "" && (*mac == "" || *broadcast == "") {
		fmt.Fprintf(os.Stderr, "Error: Must specify either -device or both -mac and -broadcast\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Prepare request
	req := WakeUpRequest{
		Device:    *device,
		Mac:       *mac,
		Broadcast: *broadcast,
	}

	// Send request
	if err := sendWakeUpRequest(*serverURL, req); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *device != "" {
		fmt.Printf("✓ Successfully sent wake-up request for device: %s\n", *device)
	} else {
		fmt.Printf("✓ Successfully sent wake-up request for MAC: %s\n", *mac)
	}
}

func sendWakeUpRequest(serverURL string, req WakeUpRequest) error {
	// Marshal request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := serverURL + "/wakeup"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned error (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
