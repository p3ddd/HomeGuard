package wol

import (
	"log/slog"
	"net"
)

// WakeOnLan sends a magic packet to wake up a machine with the given MAC address.
func WakeOnLan(macAddrStr, broadcastAddrStr string) error {
	macAddr, err := net.ParseMAC(macAddrStr)
	if err != nil {
		return err
	}

	broadcastAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(broadcastAddrStr, "9"))
	if err != nil {
		return err
	}

	magicPacket := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	for range 16 {
		magicPacket = append(magicPacket, macAddr...)
	}

	conn, err := net.DialUDP("udp", nil, broadcastAddr)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	_, err = conn.Write(magicPacket)
	if err != nil {
		return err
	}

	return nil
}

func DummyWakeOnLan(macAddr string, broadcastAddrStr string) error {
	slog.Warn("DummyWakeOnLan called", slog.String("macAddr", macAddr), slog.String("broadcastAddr", broadcastAddrStr))
	return nil
}
