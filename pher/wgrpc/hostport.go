package wgrpc

import (
	"net"
	"strconv"
)

type ListenAddress struct {
	Host string
	Port int
}

func HostPort(host string, port int) string {
	return net.JoinHostPort(host, strconv.Itoa(port))
}
