package util

import "net"

// GetMacAddress returns the first valid MAC address found on the system.
// If no MAC address is found, it returns "00:00:00:00:00:00".
func GetMacAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "00:00:00:00:00:00"
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue // skip loopback
		}
		if iface.HardwareAddr != nil {
			return iface.HardwareAddr.String()
		}
	}
	return "00:00:00:00:00:00"
}
