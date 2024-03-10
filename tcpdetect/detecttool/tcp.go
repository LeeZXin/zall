package detecttool

import (
	"fmt"
	"net"
	"time"
)

func CheckTcp(ip string, port int) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), 2*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}
