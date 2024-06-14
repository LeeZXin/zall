package detecttool

import (
	"net"
	"time"
)

func CheckTcp(ipPort string) error {
	conn, err := net.DialTimeout("tcp", ipPort, 2*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}
