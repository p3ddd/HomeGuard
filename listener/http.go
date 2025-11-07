package listener

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
)

var _ Listener = (*HTTPListener)(nil)

type HTTPListener struct {
	addr   string
	server *http.Server
	mu     sync.Mutex
}

// WakeUpPayload represents the JSON payload for wakeup requests.
type WakeUpPayload struct {
	Device    string `json:"device,omitempty"`
	Mac       string `json:"mac,omitempty"`
	Broadcast string `json:"broadcast,omitempty"`
}

func NewHTTPListener(addr string) *HTTPListener {
	return &HTTPListener{
		addr: addr,
	}
}

func (l *HTTPListener) Name() string {
	return "HTTP"
}

func (l *HTTPListener) logger() *slog.Logger {
	return slog.With("type", l.Name())
}

func (l *HTTPListener) Start(ctx context.Context, wakeUpChan chan<- WakeUpRequest) error {
	mux := http.NewServeMux()

	// Handle wakeup requests
	mux.HandleFunc("/wakeup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request WakeUpRequest
		request.Type = l.Name()

		// Try to parse JSON body first
		if r.Header.Get("Content-Type") == "application/json" {
			var payload WakeUpPayload
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				l.logger().Error("Failed to parse JSON body", "error", err)
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
			request.DeviceName = payload.Device
			request.Mac = payload.Mac
			request.Broadcast = payload.Broadcast
		} else {
			// Parse form data (query parameters or form body)
			if err := r.ParseForm(); err != nil {
				l.logger().Error("Failed to parse form", "error", err)
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}

			request.DeviceName = r.FormValue("device")
			request.Mac = r.FormValue("mac")
			request.Broadcast = r.FormValue("broadcast")
		}

		// Validate request: must have either device name or both mac and broadcast
		if request.DeviceName == "" && (request.Mac == "" || request.Broadcast == "") {
			l.logger().Error("Invalid request: must provide either device name or both mac and broadcast")
			http.Error(w, "Must provide either 'device' or both 'mac' and 'broadcast'", http.StatusBadRequest)
			return
		}

		// Send request to channel
		select {
		case wakeUpChan <- request:
			l.logger().Info("Received wakeup request",
				"device", request.DeviceName,
				"mac", request.Mac,
				"broadcast", request.Broadcast)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Wakeup request received"))
		case <-ctx.Done():
			l.logger().Info("Context canceled")
			http.Error(w, "Server shutting down", http.StatusServiceUnavailable)
		}
	})

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	l.mu.Lock()
	l.server = &http.Server{
		Addr:    l.addr,
		Handler: mux,
	}
	l.mu.Unlock()

	l.logger().Info("Starting HTTP listener", "addr", l.addr)

	// Handle graceful shutdown
	go func() {
		<-ctx.Done()
		l.logger().Info("Shutting down HTTP listener")
		shutdownCtx := context.Background()
		if err := l.server.Shutdown(shutdownCtx); err != nil {
			l.logger().Error("Failed to shutdown HTTP server", "error", err)
		}
	}()

	if err := l.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Stop implements Listener.
func (l *HTTPListener) Stop() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.server != nil {
		ctx := context.Background()
		if err := l.server.Shutdown(ctx); err != nil {
			return err
		}
		l.logger().Info("HTTP listener stopped")
	}

	return nil
}
