package dialer

import (
	"net"
	"time"
)

const (
	network = "tcp"
	port    = ":80"
)

type Dialer struct {
	timeout time.Duration
}

func (d *Dialer) Dial(host string) (int, error) {
	start := time.Now()

	conn, err := net.DialTimeout(network, host+port, d.timeout)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	return int(time.Since(start)), nil
}

func NewDialer(timeout time.Duration) *Dialer {
	return &Dialer{timeout}
}
