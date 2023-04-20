package getfreeport

import (
	"net"
	"strconv"
)

// GetFreePort starts the TCP server for a second and considers this port free,
// after which the TCP server shuts down.
func GetFreePort() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "0", err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "0", err
	}

	defer l.Close()
	port := strconv.Itoa(l.Addr().(*net.TCPAddr).Port)

	return port, nil
}
