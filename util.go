package micro

import (
	"fmt"
	"net"
	"time"
)

// GetLocIp get Local ip
func GetLocIp() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return nil, fmt.Errorf("can not find the Loc ip address: %w", err)
	}

	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP, nil
			}
		}
	}
	return nil, fmt.Errorf("can not find the Loc ip address")
}

// HumanDuration HumanDuration
func HumanDuration(t time.Duration) string {
	m := int(t.Minutes()) % 60
	h := int(t.Hours()) % 24
	d := int(t.Hours()) / 24
	if d > 0 {
		return fmt.Sprintf("%d天%d小时%d分钟", d, h, m)
	} else if h > 0 {
		return fmt.Sprintf("%d小时%d分钟", h, m)
	}
	return fmt.Sprintf("%d分钟", m)

}
