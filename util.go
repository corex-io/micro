package micro

import (
	"fmt"
	"net"
	"os"
	"time"
)

// IP 本地IP地址
var IP string

func init() {
	ip, err := GetLocIp()
	if err != nil {
		fmt.Fprintf(os.Stdout, "get LocIpAddr: %v", err)
	}
	IP = ip.String()
}

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

// GroupList 将长度为max的数据切分开来，每部分batch个，最后不够batch有多少算多少
// var tt = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
// ch := GroupList(len(tt), 3)
// for k := range ch {
//     fmt.Println(k, tt[k[0]:k[1]])
// }
func GroupList(max, batch int) <-chan [2]int {
	ch := make(chan [2]int, max/batch+1)
	defer close(ch)
	for i := 0; i <= max-1; i += batch {
		end := i + batch
		if max < end {
			end = max
		}
		ch <- [2]int{i, end}
	}
	return ch
}
