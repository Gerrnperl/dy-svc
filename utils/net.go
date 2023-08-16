package utils

import (
	"errors"
	"net"
	"strings"
)

var ip string

// GetLocalIP
//
// returns the IP address of the local machine.
// It uses a UDP connection to a IP address to determine the local IP address.
// The function caches the IP address to avoid repeated network calls.
// It returns the IP address as a string and an error (if any).
func GetLocalIP() (string, error) {
	if ip != "" {
		return ip, nil
	}
	conn, err := net.Dial("udp", "1.1.1.1:53")
	if err != nil {
		return "", errors.New("获取本地IP失败")
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(localAddr.String(), ":")[0]
	return ip, nil
}
