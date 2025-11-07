package listener

import (
	"context"
)

// ParseMAC parses s as an IEEE 802 MAC-48, EUI-48, EUI-64, or a 20-octet
// IP over InfiniBand link-layer address using one of the following formats:
//
//	00:00:5e:00:53:01
//	02:00:5e:10:00:00:00:01
//	00:00:00:00:fe:80:00:00:00:00:00:00:02:00:5e:10:00:00:00:01
//	00-00-5e-00-53-01
//	02-00-5e-10-00-00-00-01
//	00-00-00-00-fe-80-00-00-00-00-00-00-02-00-5e-10-˚00-00-00-01
//	0000.5e00.5301
//	0200.5e10.0000.0001
//	0000.0000.fe80.0000.0000.0000.0200.5e10.0000.0001

type WakeUpRequest struct {
	// HardwareAddr net.HardwareAddr
	DeviceName string // Device name for lookup (optional)
	Mac        string // MAC address (required if DeviceName is empty)
	Broadcast  string // Broadcast address (required if DeviceName is empty)
	Type       string // Listener type (HTTP, MQTT, etc.)
}

type Listener interface {
	Name() string
	Start(ctx context.Context, wakeUpChan chan<- WakeUpRequest) error
	Stop() error
}
